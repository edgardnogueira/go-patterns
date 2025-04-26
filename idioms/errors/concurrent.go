package errors

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// Error handling in concurrent code
// -----------------------------------------------------

// ConcurrentErrors demonstrates handling errors from multiple goroutines
func ConcurrentErrors(urls []string) []error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(urls)) // Buffered channel to collect errors
	
	// Launch a goroutine for each URL
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			
			// Simulate fetching the URL
			err := fetchURL(url)
			if err != nil {
				// Send the error to the error channel
				errCh <- fmt.Errorf("failed to fetch %s: %w", url, err)
			}
		}(url)
	}
	
	// Wait for all goroutines to complete
	wg.Wait()
	close(errCh) // Close the channel when done
	
	// Collect all errors
	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}
	
	return errs
}

// FirstErrorOnly demonstrates stopping on the first error
func FirstErrorOnly(urls []string) error {
	errCh := make(chan error, 1) // We only care about the first error
	done := make(chan struct{})
	
	// Launch goroutines for each URL
	for _, url := range urls {
		go func(url string) {
			// Check if we're already done
			select {
			case <-done:
				return // Another goroutine already found an error
			default:
				// Continue with the operation
			}
			
			// Simulate fetching the URL
			err := fetchURL(url)
			if err != nil {
				select {
				case errCh <- fmt.Errorf("failed to fetch %s: %w", url, err):
					close(done) // Signal other goroutines to stop
				default:
					// Another goroutine already reported an error
				}
			}
		}(url)
	}
	
	// Wait for the first error or completion
	select {
	case err := <-errCh:
		return err
	case <-time.After(5 * time.Second):
		return nil // Timeout, all operations succeeded or took too long
	}
}

// WithContext demonstrates using context for cancellation upon errors
func WithContext(ctx context.Context, tasks []string) error {
	// Create a new context that we can cancel
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // Ensure all resources are released
	
	errCh := make(chan error, 1)
	var wg sync.WaitGroup
	
	// Launch goroutines for each task
	for _, task := range tasks {
		wg.Add(1)
		go func(task string) {
			defer wg.Done()
			
			// Process the task, checking for context cancellation
			err := processWithContext(ctx, task)
			if err != nil {
				// Report the error
				select {
				case errCh <- fmt.Errorf("task %s failed: %w", task, err):
					cancel() // Cancel other tasks
				default:
					// Channel is full, another error was already reported
				}
			}
		}(task)
	}
	
	// Wait in a goroutine so we can select between completion and errors
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	
	// Wait for completion or an error
	select {
	case <-done:
		// All tasks completed successfully
		return nil
	case err := <-errCh:
		// Wait for all goroutines to notice cancellation and exit
		wg.Wait()
		return err
	case <-ctx.Done():
		// External cancellation
		wg.Wait()
		return ctx.Err()
	}
}

// ErrGroup demonstrates using a simple error group pattern
type ErrGroup struct {
	wg      sync.WaitGroup
	errOnce sync.Once
	err     error
}

// Go adds a function to the error group
func (g *ErrGroup) Go(f func() error) {
	g.wg.Add(1)
	
	go func() {
		defer g.wg.Done()
		
		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
			})
		}
	}()
}

// Wait waits for all functions to complete and returns the first error
func (g *ErrGroup) Wait() error {
	g.wg.Wait()
	return g.err
}

// WithErrGroup demonstrates using our ErrGroup implementation
func WithErrGroup(urls []string) error {
	var g ErrGroup
	
	for _, url := range urls {
		url := url // https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			return fetchURL(url)
		})
	}
	
	return g.Wait()
}

// Aggregating errors with a fan-in pattern
// -----------------------------------------------------

// FanInErrors demonstrates collecting errors from multiple goroutines
func FanInErrors(tasks []string) []error {
	// Create channels for errors and synchronization
	errCh := make(chan error, len(tasks))
	doneCh := make(chan struct{}, len(tasks))
	
	// Launch goroutines for each task
	for _, task := range tasks {
		go func(task string) {
			// Process the task
			err := processTask(task)
			if err != nil {
				errCh <- fmt.Errorf("error processing %s: %w", task, err)
			}
			
			// Signal completion
			doneCh <- struct{}{}
		}(task)
	}
	
	// Wait for all tasks to complete
	var errs []error
	for i := 0; i < len(tasks); i++ {
		<-doneCh
	}
	
	// Collect any errors
	close(errCh)
	for err := range errCh {
		errs = append(errs, err)
	}
	
	return errs
}

// Retry with exponential backoff
// -----------------------------------------------------

// RetryWithBackoff demonstrates retrying operations with exponential backoff
func RetryWithBackoff(op func() error, maxRetries int) error {
	var err error
	
	for attempt := 0; attempt < maxRetries; attempt++ {
		err = op()
		if err == nil {
			return nil // Success
		}
		
		// If this is a permanent error, don't retry
		var permErr *PermanentError
		if errors.As(err, &permErr) {
			return err
		}
		
		// Calculate backoff duration (exponential with jitter)
		backoff := calculateBackoff(attempt)
		time.Sleep(backoff)
	}
	
	return fmt.Errorf("failed after %d attempts: %w", maxRetries, err)
}

// PermanentError represents an error that should not be retried
type PermanentError struct {
	Err error
}

// Error implements the error interface
func (e *PermanentError) Error() string {
	return fmt.Sprintf("permanent error: %v", e.Err)
}

// Unwrap returns the underlying error
func (e *PermanentError) Unwrap() error {
	return e.Err
}

// Helper functions for examples
// -----------------------------------------------------

func fetchURL(url string) error {
	// Simulate fetching a URL
	switch url {
	case "error":
		return errors.New("connection failed")
	case "timeout":
		time.Sleep(100 * time.Millisecond)
		return ErrTimeout
	case "not-found":
		return ErrNotFound
	case "permanent-error":
		return &PermanentError{Err: errors.New("permanent failure")}
	default:
		// Simulate successful fetch after a short delay
		time.Sleep(50 * time.Millisecond)
		return nil
	}
}

func processWithContext(ctx context.Context, task string) error {
	// Check for cancellation before starting
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Continue with the task
	}
	
	// Simulate processing time
	select {
	case <-time.After(100 * time.Millisecond):
		// Task completed
	case <-ctx.Done():
		return ctx.Err()
	}
	
	// Simulate errors for certain tasks
	switch task {
	case "error":
		return errors.New("processing failed")
	case "timeout":
		return ErrTimeout
	default:
		return nil
	}
}

func processTask(task string) error {
	// Simulate processing
	time.Sleep(50 * time.Millisecond)
	
	// Simulate errors for certain tasks
	switch task {
	case "error":
		return errors.New("task failed")
	case "timeout":
		return ErrTimeout
	default:
		return nil
	}
}

func calculateBackoff(attempt int) time.Duration {
	// Simple exponential backoff calculation
	// In a real implementation, you'd add jitter
	base := time.Millisecond * 100
	max := time.Second * 10
	
	// 2^attempt * base
	duration := time.Duration(1<<uint(attempt)) * base
	
	// Cap at maximum backoff
	if duration > max {
		duration = max
	}
	
	return duration
}
