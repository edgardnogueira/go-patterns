package fanout_fanin

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync/atomic"
	"testing"
	"time"
)

// Test FanOut function with successful processing
func TestFanOut_Success(t *testing.T) {
	ctx := context.Background()
	
	// Create input channel
	inputs := make(chan int, 10)
	for i := 1; i <= 10; i++ {
		inputs <- i
	}
	close(inputs)
	
	// Define worker function that squares its input
	squareWorker := func(ctx context.Context, x int) (int, error) {
		return x * x, nil
	}
	
	// Fan out with 3 workers
	resultChan, errChan := FanOut(ctx, inputs, squareWorker, 3)
	
	// Collect results
	var results []int
	for result := range resultChan {
		results = append(results, result)
	}
	
	// Check for errors
	for err := range errChan {
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	}
	
	// Verify results (order may vary)
	if len(results) != 10 {
		t.Errorf("Expected 10 results, got %d", len(results))
	}
	
	// Sort results for deterministic comparison
	sort.Ints(results)
	
	// Expected: squares of numbers 1-10
	expected := []int{1, 4, 9, 16, 25, 36, 49, 64, 81, 100}
	
	for i, v := range expected {
		if results[i] != v {
			t.Errorf("Result at index %d: expected %d, got %d", i, v, results[i])
		}
	}
}

// Test FanOut function with errors
func TestFanOut_Error(t *testing.T) {
	ctx := context.Background()
	
	// Create input channel
	inputs := make(chan int, 5)
	for i := -2; i <= 2; i++ {
		inputs <- i
	}
	close(inputs)
	
	// Define worker function that returns error for negative numbers
	worker := func(ctx context.Context, x int) (int, error) {
		if x < 0 {
			return 0, fmt.Errorf("negative number not allowed: %d", x)
		}
		return x * 10, nil
	}
	
	// Fan out with 2 workers
	resultChan, errChan := FanOut(ctx, inputs, worker, 2)
	
	// Collect results and errors
	var results []int
	var errs []error
	
	for result := range resultChan {
		results = append(results, result)
	}
	
	for err := range errChan {
		if err != nil {
			errs = append(errs, err)
		}
	}
	
	// Verify expected results and errors
	if len(results) != 3 {
		t.Errorf("Expected 3 successful results, got %d", len(results))
	}
	
	if len(errs) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(errs))
	}
	
	// Sort results for deterministic comparison
	sort.Ints(results)
	expected := []int{0, 10, 20}
	
	for i, v := range expected {
		if i < len(results) && results[i] != v {
			t.Errorf("Result at index %d: expected %d, got %d", i, v, results[i])
		}
	}
}

// Test context cancellation
func TestFanOut_ContextCancellation(t *testing.T) {
	// Create a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	
	// Create a buffered input channel
	inputs := make(chan int, 100)
	for i := 1; i <= 100; i++ {
		inputs <- i
	}
	close(inputs)
	
	// Create a worker that sleeps
	worker := func(ctx context.Context, x int) (int, error) {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-time.After(50 * time.Millisecond):
			return x * 2, nil
		}
	}
	
	// Fan out with 5 workers
	resultChan, errChan := FanOut(ctx, inputs, worker, 5)
	
	// Count results and cancel early
	var resultCount atomic.Int32
	go func() {
		for range resultChan {
			resultCount.Add(1)
			if resultCount.Load() >= 10 {
				cancel() // Cancel after receiving 10 results
				break
			}
		}
	}()
	
	// Count errors
	var errorCount atomic.Int32
	for err := range errChan {
		if err != nil {
			errorCount.Add(1)
		}
	}
	
	finalResultCount := resultCount.Load()
	if finalResultCount < 10 {
		t.Errorf("Expected at least 10 results before cancellation, got %d", finalResultCount)
	}
	
	// We can't predict exact error count, but we expect some due to cancellation
	if errorCount.Load() == 0 {
		t.Errorf("Expected some errors after cancellation, got none")
	}
}

// Test the FanIn function
func TestFanIn(t *testing.T) {
	ctx := context.Background()
	
	// Create 3 input channels
	ch1 := make(chan string, 3)
	ch2 := make(chan string, 3)
	ch3 := make(chan string, 3)
	
	// Send values to channels
	ch1 <- "A1"
	ch1 <- "A2"
	ch1 <- "A3"
	close(ch1)
	
	ch2 <- "B1"
	ch2 <- "B2"
	close(ch2)
	
	ch3 <- "C1"
	close(ch3)
	
	// Fan in the channels
	combined := FanIn(ctx, ch1, ch2, ch3)
	
	// Collect results
	var results []string
	for result := range combined {
		results = append(results, result)
	}
	
	// Verify count (order is non-deterministic)
	if len(results) != 6 {
		t.Errorf("Expected 6 results, got %d", len(results))
	}
	
	// Verify that all values are present
	expectedItems := map[string]bool{
		"A1": false, "A2": false, "A3": false,
		"B1": false, "B2": false, "C1": false,
	}
	
	for _, item := range results {
		expectedItems[item] = true
	}
	
	for item, found := range expectedItems {
		if !found {
			t.Errorf("Expected item %s not found in results", item)
		}
	}
}

// Test the ProcessAll convenience function
func TestProcessAll(t *testing.T) {
	ctx := context.Background()
	
	// Create input data with some values that will cause errors
	inputs := []int{-5, -2, 0, 3, 7, 12}
	
	// Create work function
	workFn := func(ctx context.Context, x int) (string, error) {
		if x < 0 {
			return "", errors.New("negative value")
		}
		if x == 0 {
			return "", errors.New("zero value")
		}
		return fmt.Sprintf("Processed %d", x), nil
	}
	
	// Process all inputs
	results, errs := ProcessAll(ctx, inputs, workFn, 3)
	
	// Verify results
	if len(results) != 3 {
		t.Errorf("Expected 3 successful results, got %d", len(results))
	}
	
	if len(errs) != 3 {
		t.Errorf("Expected 3 errors, got %d", len(errs))
	}
	
	// Expected successful results (order may vary)
	expectedResults := map[string]bool{
		"Processed 3":  false,
		"Processed 7":  false,
		"Processed 12": false,
	}
	
	for _, result := range results {
		expectedResults[result] = true
	}
	
	for result, found := range expectedResults {
		if !found {
			t.Errorf("Expected result %s not found", result)
		}
	}
}
