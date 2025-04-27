// Package workerpool implements the Worker Pool concurrency pattern
package workerpool

import (
	"context"
	"sync"
)

// Task represents a unit of work to be processed
type Task interface {
	Execute() interface{}
}

// Result represents the output from a task execution
type Result struct {
	Value    interface{}
	TaskID   interface{}
	WorkerID int
	Error    error
}

// Worker represents a worker that processes tasks
type Worker struct {
	ID         int
	taskQueue  <-chan Task
	resultChan chan<- Result
	wg         *sync.WaitGroup
	quit       chan struct{}
}

// NewWorker creates a new worker
func NewWorker(id int, taskQueue <-chan Task, resultChan chan<- Result, wg *sync.WaitGroup) *Worker {
	return &Worker{
		ID:         id,
		taskQueue:  taskQueue,
		resultChan: resultChan,
		wg:         wg,
		quit:       make(chan struct{}),
	}
}

// Start begins the worker's processing loop
func (w *Worker) Start() {
	go func() {
		defer w.wg.Done()
		for {
			select {
			case task, ok := <-w.taskQueue:
				if !ok {
					return
				}
				result := task.Execute()
				w.resultChan <- Result{
					Value:    result,
					WorkerID: w.ID,
				}
			case <-w.quit:
				return
			}
		}
	}()
}

// Stop signals the worker to stop processing
func (w *Worker) Stop() {
	close(w.quit)
}

// Pool represents a pool of workers
type Pool struct {
	Tasks       chan Task
	Results     chan Result
	Workers     []*Worker
	workerCount int
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewPool creates a new worker pool with the specified number of workers
func NewPool(workerCount int) *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &Pool{
		Tasks:       make(chan Task),
		Results:     make(chan Result),
		workerCount: workerCount,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Start initializes and starts the worker pool
func (p *Pool) Start() {
	// Create workers
	p.Workers = make([]*Worker, p.workerCount)
	p.wg.Add(p.workerCount)
	
	for i := 0; i < p.workerCount; i++ {
		worker := NewWorker(i, p.Tasks, p.Results, &p.wg)
		p.Workers[i] = worker
		worker.Start()
	}
	
	// Start a goroutine to close Results channel when all workers are done
	go func() {
		p.wg.Wait()
		close(p.Results)
	}()
}

// Submit adds a task to the pool
func (p *Pool) Submit(task Task) {
	select {
	case <-p.ctx.Done():
		return // Don't submit if the pool is shutting down
	case p.Tasks <- task:
		// Task submitted
	}
}

// Stop gracefully shuts down the worker pool
func (p *Pool) Stop() {
	p.cancel() // Signal cancellation
	close(p.Tasks) // Close tasks channel to signal workers to stop
	
	// Stop all workers
	for _, worker := range p.Workers {
		worker.Stop()
	}
	
	p.wg.Wait() // Wait for all workers to finish
}

// WithContext creates a worker pool with a custom context
func WithContext(ctx context.Context, workerCount int) *Pool {
	childCtx, cancel := context.WithCancel(ctx)
	
	return &Pool{
		Tasks:       make(chan Task),
		Results:     make(chan Result),
		workerCount: workerCount,
		ctx:         childCtx,
		cancel:      cancel,
	}
}
