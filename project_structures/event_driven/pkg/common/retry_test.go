package common_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/edgardnogueira/go-patterns/project_structures/event_driven/pkg/common"
)

// TestRetrySuccessOnFirstAttempt tests that the retry function succeeds on the first attempt
func TestRetrySuccessOnFirstAttempt(t *testing.T) {
	attempts := 0
	operation := func() error {
		attempts++
		return nil // Succeed immediately
	}

	opts := common.RetryOptions{
		MaxRetries:  3,
		InitialWait: 1 * time.Millisecond,
		MaxWait:     10 * time.Millisecond,
		Factor:      2.0,
	}

	err := common.Retry(operation, opts)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if attempts != 1 {
		t.Errorf("Expected 1 attempt, got %d", attempts)
	}
}

// TestRetryEventualSuccess tests that the retry function succeeds after several attempts
func TestRetryEventualSuccess(t *testing.T) {
	attempts := 0
	operation := func() error {
		attempts++
		if attempts < 3 {
			return errors.New("temporary error")
		}
		return nil // Succeed on third attempt
	}

	opts := common.RetryOptions{
		MaxRetries:  5,
		InitialWait: 1 * time.Millisecond,
		MaxWait:     10 * time.Millisecond,
		Factor:      2.0,
	}

	err := common.Retry(operation, opts)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

// TestRetryMaxAttemptsExceeded tests that the retry function fails after max attempts
func TestRetryMaxAttemptsExceeded(t *testing.T) {
	attempts := 0
	expectedError := errors.New("persistent error")

	operation := func() error {
		attempts++
		return expectedError // Always fail
	}

	opts := common.RetryOptions{
		MaxRetries:  3,
		InitialWait: 1 * time.Millisecond,
		MaxWait:     10 * time.Millisecond,
		Factor:      2.0,
	}

	err := common.Retry(operation, opts)

	if err == nil {
		t.Error("Expected an error, got nil")
	}

	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}

	if attempts != opts.MaxRetries+1 { // Initial attempt + retries
		t.Errorf("Expected %d attempts, got %d", opts.MaxRetries+1, attempts)
	}
}

// TestRetryWithContextCancellation tests that the retry function respects context cancellation
func TestRetryWithContextCancellation(t *testing.T) {
	attempts := 0
	operation := func() error {
		attempts++
		return errors.New("temporary error")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	opts := common.RetryOptions{
		MaxRetries:  5,
		InitialWait: 20 * time.Millisecond, // Longer than context timeout
		MaxWait:     100 * time.Millisecond,
		Factor:      2.0,
	}

	err := common.RetryWithContext(ctx, operation, opts)

	if err == nil {
		t.Error("Expected an error due to context cancellation, got nil")
	}

	if attempts != 1 {
		t.Errorf("Expected 1 attempt before context cancellation, got %d", attempts)
	}
}
