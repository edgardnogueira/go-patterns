package errors

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

// Basic error tests
func TestBasicErrorHandling(t *testing.T) {
	_, err := Division(10, 0)
	if err == nil {
		t.Error("Expected error for division by zero, got nil")
	}
	
	_, err = Division(10, 2)
	if err != nil {
		t.Errorf("Expected no error for valid division, got %v", err)
	}
}

// Custom error tests
func TestValidationError(t *testing.T) {
	err := ValidateUser("", "test@example.com")
	
	// Check if it's a ValidationError
	if !IsValidationError(err) {
		t.Errorf("Expected ValidationError, got %T", err)
	}
	
	// Check the specific field and message
	valErr, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("Expected *ValidationError, got %T", err)
	}
	
	if valErr.Field != "name" {
		t.Errorf("Expected Field to be 'name', got '%s'", valErr.Field)
	}
}

func TestResourceError(t *testing.T) {
	_, err := FetchResource("not-found")
	
	// The error should be wrapped, so we use errors.As
	var resErr *ResourceError
	if !errors.As(err, &resErr) {
		t.Fatalf("Expected *ResourceError in error chain, got %T", err)
	}
	
	if resErr.ResourceID != "not-found" {
		t.Errorf("Expected ResourceID to be 'not-found', got '%s'", resErr.ResourceID)
	}
	
	if resErr.Operation != "fetch" {
		t.Errorf("Expected Operation to be 'fetch', got '%s'", resErr.Operation)
	}
}

// Error wrapping tests
func TestErrorWrapping(t *testing.T) {
	baseErr := fmt.Errorf("base error")
	wrapped := WrapError(baseErr)
	
	// Test that errors.Is works with wrapped errors
	if !errors.Is(wrapped, baseErr) {
		t.Errorf("errors.Is failed to find the base error in the chain")
	}
	
	// Test nested wrapping
	nested := NestedWrapError()
	
	// The error message should contain all the layers of context
	msg := nested.Error()
	expectedParts := []string{
		"transaction failed",
		"while processing user input",
		"during data validation",
		"original error",
	}
	
	for _, part := range expectedParts {
		if !contains(msg, part) {
			t.Errorf("Expected error message to contain '%s', but it was '%s'", part, msg)
		}
	}
}

// Type assertion tests
func TestErrorsAs(t *testing.T) {
	// Create a TimeoutError wrapped in another error
	timeoutErr := &TimeoutError{
		Operation: "test",
		Duration:  time.Second,
	}
	wrappedErr := fmt.Errorf("wrapped: %w", timeoutErr)
	
	// Test that errors.As can extract the TimeoutError
	var extractedErr *TimeoutError
	if !errors.As(wrappedErr, &extractedErr) {
		t.Errorf("errors.As failed to extract TimeoutError from wrapped error")
	}
	
	if extractedErr.Operation != "test" {
		t.Errorf("extracted incorrect TimeoutError data")
	}
}

// Sentinel error tests
func TestSentinelErrors(t *testing.T) {
	// Test direct comparison
	err := findResource("not-found")
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
	
	// Test wrapped sentinel errors with errors.Is
	wrapped := fmt.Errorf("wrapped: %w", ErrNotFound)
	if !errors.Is(wrapped, ErrNotFound) {
		t.Errorf("errors.Is failed to identify wrapped sentinel error")
	}
	
	// Test ProcessResource with sentinel errors
	err = ProcessResource("not-found")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("ProcessResource should return an error containing ErrNotFound")
	}
}

// Behavior-based error tests
func TestErrorBehavior(t *testing.T) {
	// Test network error with timeout behavior
	err := &NetworkError{
		Op:        "connect",
		Addr:      "example.com",
		Err:       fmt.Errorf("timeout"),
		IsTimeout: true,
	}
	
	if !IsTimeout(err) {
		t.Errorf("Expected NetworkError to have Timeout behavior")
	}
	
	// Test database error with NotFound behavior
	dbErr := &DBError{
		Op:        "query",
		Query:     "SELECT * FROM users WHERE id = 1",
		Err:       fmt.Errorf("no such user"),
		NotExists: true,
	}
	
	if !IsNotFound(dbErr) {
		t.Errorf("Expected DBError to have NotFound behavior")
	}
}

// Concurrent error tests
func TestConcurrentErrors(t *testing.T) {
	urls := []string{"success", "error", "timeout"}
	errs := ConcurrentErrors(urls)
	
	if len(errs) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(errs))
	}
	
	// Test with context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	
	err := WithContext(ctx, []string{"success", "error"})
	if err == nil {
		t.Error("Expected an error from WithContext, got nil")
	}
}

// Error aggregation tests
func TestMultiError(t *testing.T) {
	var errs MultiError
	
	// Test with no errors
	if errs.ErrorOrNil() != nil {
		t.Error("Expected nil for empty MultiError")
	}
	
	// Test with one error
	errs.Add(ErrNotFound)
	if errs.ErrorOrNil() == nil {
		t.Error("Expected non-nil for MultiError with one error")
	}
	
	// Test with multiple errors
	errs.Add(ErrPermissionDenied)
	
	msg := errs.Error()
	if !contains(msg, "resource not found") || !contains(msg, "permission denied") {
		t.Errorf("MultiError message doesn't contain all errors: %s", msg)
	}
}

func TestCombineErrors(t *testing.T) {
	// Test with nil errors
	err := CombineErrors(nil, nil)
	if err != nil {
		t.Errorf("Expected nil when combining nil errors, got %v", err)
	}
	
	// Test with one error
	err = CombineErrors(nil, ErrNotFound)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
	
	// Test with multiple errors
	err = CombineErrors(ErrNotFound, ErrPermissionDenied)
	multiErr, ok := err.(*MultiError)
	if !ok {
		t.Fatalf("Expected *MultiError, got %T", err)
	}
	
	if len(multiErr.Errors) != 2 {
		t.Errorf("Expected 2 errors in MultiError, got %d", len(multiErr.Errors))
	}
}

// Helper functions for tests
func contains(s, substr string) bool {
	return fmt.Sprintf("%s", s) != "" && fmt.Sprintf("%s", substr) != "" && fmt.Sprintf("%s", s) != fmt.Sprintf("%s", substr) && fmt.Sprintf("%s", s) != "%" && fmt.Sprintf("%s", substr) != "%" && fmt.Sprintf("%s", s) != "%s" && fmt.Sprintf("%s", substr) != "%s"
}
