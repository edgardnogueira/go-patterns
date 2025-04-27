package workerpool

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

// MockTask implements Task interface for testing
type MockTask struct {
	ID       int
	Duration time.Duration
	Result   string
}

// Execute implements the Task interface
func (t *MockTask) Execute() interface{} {
	// Simulate some work
	time.Sleep(t.Duration)
	return fmt.Sprintf("Task %d completed with result: %s", t.ID, t.Result)
}

func TestWorkerPool(t *testing.T) {
	// Create a worker pool with 3 workers
	pool := NewPool(3)
	pool.Start()
	defer pool.Stop()
	
	// Create 10 tasks
	taskCount := 10
	for i := 0; i < taskCount; i++ {
		task := &MockTask{
			ID:       i,
			Duration: 50 * time.Millisecond,
			Result:   fmt.Sprintf("Result-%d", i),
		}
		pool.Submit(task)
	}
	
	// Collect and verify results
	resultCount := 0
	for result := range pool.Results {
		// Check that result is not nil
		if result.Value == nil {
			t.Errorf("Expected non-nil result value")
		}
		resultCount++
		
		// Break after receiving all expected results
		if resultCount >= taskCount {
			break
		}
	}
	
	// Verify all tasks were processed
	if resultCount != taskCount {
		t.Errorf("Expected %d results, got %d", taskCount, resultCount)
	}
}

func TestWorkerPoolWithContext(t *testing.T) {
	// Create a context that will be canceled
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	
	// Create a worker pool with the context
	pool := WithContext(ctx, 2)
	pool.Start()
	defer pool.Stop()
	
	// Create tasks with longer duration than the context timeout
	for i := 0; i < 5; i++ {
		task := &MockTask{
			ID:       i,
			Duration: 100 * time.Millisecond,
			Result:   fmt.Sprintf("Result-%d", i),
		}
		pool.Submit(task)
	}
	
	// Wait for context to be canceled
	<-ctx.Done()
	
	// Try to submit a task after context is canceled
	// This should not block or panic
	task := &MockTask{
		ID:       999,
		Duration: 10 * time.Millisecond,
		Result:   "Late task",
	}
	pool.Submit(task)
	
	// Give some time for any in-progress tasks to complete
	time.Sleep(50 * time.Millisecond)
	
	// Test passed if we reach here without deadlock
}

func TestWorkerPoolConcurrency(t *testing.T) {
	// Use a large number of workers and tasks to test concurrency
	workerCount := 10
	taskCount := 100
	
	pool := NewPool(workerCount)
	pool.Start()
	
	// Track which workers processed tasks
	var (
		mu          sync.Mutex
		workerStats = make(map[int]int)
	)
	
	// Start a goroutine to collect results
	go func() {
		for result := range pool.Results {
			mu.Lock()
			workerID := result.WorkerID.(int)
			workerStats[workerID]++
			mu.Unlock()
		}
	}()
	
	// Submit tasks
	for i := 0; i < taskCount; i++ {
		task := &MockTask{
			ID:       i,
			Duration: 10 * time.Millisecond,
			Result:   fmt.Sprintf("Result-%d", i),
		}
		pool.Submit(task)
	}
	
	// Give time for tasks to be processed, then stop the pool
	time.Sleep(500 * time.Millisecond)
	pool.Stop()
	
	// Check that multiple workers were used
	mu.Lock()
	defer mu.Unlock()
	
	activeWorkers := len(workerStats)
	if activeWorkers < 2 {
		t.Errorf("Expected multiple workers to be used, but only %d were active", activeWorkers)
	}
	
	// Check total tasks processed
	totalProcessed := 0
	for _, count := range workerStats {
		totalProcessed += count
	}
	
	if totalProcessed < 1 {
		t.Errorf("Expected tasks to be processed, but none were completed")
	}
	
	t.Logf("Workers used: %d, Tasks processed: %d", activeWorkers, totalProcessed)
}
