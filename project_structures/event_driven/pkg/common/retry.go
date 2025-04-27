package common

import (
	"fmt"
	"log"
	"time"
)

// RetryOptions configures the retry strategy
type RetryOptions struct {
	MaxRetries  int
	InitialWait time.Duration
	MaxWait     time.Duration
	Factor      float64
}

// DefaultRetryOptions returns the default retry options
func DefaultRetryOptions() RetryOptions {
	return RetryOptions{
		MaxRetries:  5,
		InitialWait: 100 * time.Millisecond,
		MaxWait:     10 * time.Second,
		Factor:      2.0,
	}
}

// Retry executes a function with exponential backoff retry
func Retry(operation func() error, opts RetryOptions) error {
	var err error
	wait := opts.InitialWait

	for i := 0; i <= opts.MaxRetries; i++ {
		// First attempt or retry
		err = operation()
		
		// If successful or max retries reached
		if err == nil || i == opts.MaxRetries {
			break
		}
		
		// Calculate next wait time with exponential backoff
		nextWait := time.Duration(float64(wait) * opts.Factor)
		if nextWait > opts.MaxWait {
			nextWait = opts.MaxWait
		}
		
		log.Printf("Retry %d/%d after error: %v. Waiting %v before next attempt", 
			i+1, opts.MaxRetries, err, wait)
		
		// Wait before next retry
		time.Sleep(wait)
		wait = nextWait
	}

	return err
}

// RetryWithContext executes a function with exponential backoff retry and context support
func RetryWithContext(ctx context.Context, operation func() error, opts RetryOptions) error {
	var err error
	wait := opts.InitialWait

	for i := 0; i <= opts.MaxRetries; i++ {
		// First attempt or retry
		err = operation()
		
		// If successful or max retries reached
		if err == nil || i == opts.MaxRetries {
			break
		}
		
		// Calculate next wait time with exponential backoff
		nextWait := time.Duration(float64(wait) * opts.Factor)
		if nextWait > opts.MaxWait {
			nextWait = opts.MaxWait
		}
		
		log.Printf("Retry %d/%d after error: %v. Waiting %v before next attempt", 
			i+1, opts.MaxRetries, err, wait)
		
		// Create a timer for the wait
		timer := time.NewTimer(wait)
		select {
		case <-timer.C:
			// Proceed with the next retry
		case <-ctx.Done():
			// Context canceled or timed out
			timer.Stop()
			return fmt.Errorf("operation canceled: %w", ctx.Err())
		}
		
		wait = nextWait
	}

	return err
}
