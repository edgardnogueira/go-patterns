// Package fanout_fanin implements the Fan-Out/Fan-In concurrency pattern
package fanout_fanin

import (
	"context"
	"sync"
)

// WorkFunc is a function that processes input data and returns a result
type WorkFunc[I, O any] func(ctx context.Context, input I) (O, error)

// FanOut distributes work across multiple goroutines and returns results
// through a single output channel (Fan-Out, then Fan-In)
func FanOut[I, O any](
	ctx context.Context,
	inputs <-chan I,
	workFn WorkFunc[I, O],
	workerCount int,
) (<-chan O, <-chan error) {
	// Create output channels
	results := make(chan O)
	errs := make(chan error)
	
	// Create a wait group to track when all goroutines are done
	var wg sync.WaitGroup
	
	// Launch multiple workers (fan-out)
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		
		go func(workerId int) {
			defer wg.Done()
			
			// Process inputs until the channel is closed or context is canceled
			for input := range inputs {
				// Check if context is canceled
				select {
				case <-ctx.Done():
					return
				default:
					// Process the input
					result, err := workFn(ctx, input)
					if err != nil {
						select {
						case errs <- err:
							// Error sent
						case <-ctx.Done():
							return
						}
					} else {
						select {
						case results <- result:
							// Result sent
						case <-ctx.Done():
							return
						}
					}
				}
			}
		}(i)
	}
	
	// Close channels when all workers are done (fan-in)
	go func() {
		wg.Wait()
		close(results)
		close(errs)
	}()
	
	return results, errs
}

// FanIn combines multiple input channels into a single output channel
func FanIn[T any](ctx context.Context, channels ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	multiplexed := make(chan T)
	
	// Start a goroutine for each input channel
	multiplex := func(ch <-chan T) {
		defer wg.Done()
		for item := range ch {
			select {
			case <-ctx.Done():
				return
			case multiplexed <- item:
			}
		}
	}
	
	// Add waitgroup count and start goroutine for each channel
	wg.Add(len(channels))
	for _, ch := range channels {
		go multiplex(ch)
	}
	
	// Close multiplexed channel when all input channels are done
	go func() {
		wg.Wait()
		close(multiplexed)
	}()
	
	return multiplexed
}

// Distribute takes a slice of inputs and distributes them across multiple goroutines,
// then collects results into a single channel
func Distribute[I, O any](
	ctx context.Context,
	inputs []I,
	workFn WorkFunc[I, O],
	workerCount int,
) (<-chan O, <-chan error) {
	// Create a channel to send inputs to workers
	inputChan := make(chan I)
	
	// Start a goroutine to send inputs
	go func() {
		defer close(inputChan)
		for _, input := range inputs {
			select {
			case <-ctx.Done():
				return
			case inputChan <- input:
				// Input sent
			}
		}
	}()
	
	// Fan out the work and collect results
	return FanOut(ctx, inputChan, workFn, workerCount)
}

// ProcessAll takes a slice of inputs, processes them concurrently, and returns
// all successful results and errors
func ProcessAll[I, O any](
	ctx context.Context,
	inputs []I,
	workFn WorkFunc[I, O],
	workerCount int,
) ([]O, []error) {
	// Create a context with cancellation
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	
	// Distribute the work
	resultChan, errChan := Distribute(ctx, inputs, workFn, workerCount)
	
	// Collect results and errors
	var results []O
	var errors []error
	
	// Process result channel
	for result := range resultChan {
		results = append(results, result)
	}
	
	// Process error channel
	for err := range errChan {
		if err != nil {
			errors = append(errors, err)
		}
	}
	
	return results, errors
}
