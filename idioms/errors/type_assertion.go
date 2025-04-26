package errors

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"time"
)

// Type assertion with errors
// -----------------------------------------------------

// TimeoutError represents an error that occurred because an operation timed out
type TimeoutError struct {
	Operation string
	Duration  time.Duration
}

// Error implements the error interface for TimeoutError
func (e *TimeoutError) Error() string {
	return fmt.Sprintf("%s timed out after %v", e.Operation, e.Duration)
}

// Timeout is a method that can be used to identify this error type
func (e *TimeoutError) Timeout() bool {
	return true
}

// NotFoundError represents an error when a resource is not found
type NotFoundError struct {
	ResourceType string
	ID           string
}

// Error implements the error interface for NotFoundError
func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %s not found", e.ResourceType, e.ID)
}

// NotFound is a method that can be used to identify this error type
func (e *NotFoundError) NotFound() bool {
	return true
}

// TypeAssertionExample demonstrates traditional type assertion with errors
func TypeAssertionExample(err error) {
	// Traditional type assertion approach (pre-Go 1.13)
	if timeoutErr, ok := err.(*TimeoutError); ok {
		fmt.Printf("Timeout error detected: %s timed out after %v\n", 
			timeoutErr.Operation, timeoutErr.Duration)
		return
	}
	
	if notFoundErr, ok := err.(*NotFoundError); ok {
		fmt.Printf("Not found error detected: %s with ID %s\n", 
			notFoundErr.ResourceType, notFoundErr.ID)
		return
	}
	
	fmt.Println("Unknown error type:", err)
}

// Using errors.As for type assertion (Go 1.13+)
// -----------------------------------------------------

// ErrorsAsExample demonstrates using errors.As for cleaner type assertion
func ErrorsAsExample(err error) {
	// Using errors.As for cleaner type assertion with wrapped errors
	var timeoutErr *TimeoutError
	if errors.As(err, &timeoutErr) {
		fmt.Printf("Timeout error detected using errors.As: %s timed out after %v\n", 
			timeoutErr.Operation, timeoutErr.Duration)
		return
	}
	
	var notFoundErr *NotFoundError
	if errors.As(err, &notFoundErr) {
		fmt.Printf("Not found error detected using errors.As: %s with ID %s\n", 
			notFoundErr.ResourceType, notFoundErr.ID)
		return
	}
	
	// Check for standard library error types too
	var pathErr *fs.PathError
	if errors.As(err, &pathErr) {
		fmt.Printf("Path error detected: operation %s on path %s failed: %v\n",
			pathErr.Op, pathErr.Path, pathErr.Err)
		return
	}
	
	fmt.Println("Unknown error type:", err)
}

// Practical examples
// -----------------------------------------------------

// FindUserByID simulates a database lookup that might fail in various ways
func FindUserByID(id string) (*User, error) {
	// Simulate various error conditions
	if id == "" {
		return nil, fmt.Errorf("invalid user ID: %w", &ValidationError{
			Field:   "id",
			Message: "ID cannot be empty",
		})
	}

	if id == "timeout" {
		return nil, fmt.Errorf("database error: %w", &TimeoutError{
			Operation: "user lookup",
			Duration:  time.Second * 5,
		})
	}

	if id == "not-found" {
		return nil, fmt.Errorf("user lookup failed: %w", &NotFoundError{
			ResourceType: "user",
			ID:           id,
		})
	}

	// Successful case
	return &User{ID: id, Name: "Example User"}, nil
}

// User is a simple example struct
type User struct {
	ID   string
	Name string
}

// HandleUserLookup demonstrates handling multiple error types
func HandleUserLookup(id string) {
	user, err := FindUserByID(id)
	if err != nil {
		// Check for specific error types using errors.As
		var valErr *ValidationError
		if errors.As(err, &valErr) {
			fmt.Printf("Validation error: %s - %s\n", valErr.Field, valErr.Message)
			return
		}
		
		var timeoutErr *TimeoutError
		if errors.As(err, &timeoutErr) {
			fmt.Printf("The operation timed out after %v. Please try again later.\n", 
				timeoutErr.Duration)
			return
		}
		
		var notFoundErr *NotFoundError
		if errors.As(err, &notFoundErr) {
			fmt.Printf("Sorry, we couldn't find the %s with ID %s\n", 
				notFoundErr.ResourceType, notFoundErr.ID)
			return
		}
		
		// Generic error handler for unknown error types
		fmt.Printf("Unexpected error occurred: %v\n", err)
		return
	}
	
	// Process the user
	fmt.Printf("Found user: %s (ID: %s)\n", user.Name, user.ID)
}

// Checking for error behavior rather than concrete types
// -----------------------------------------------------

// TimeoutChecker defines behavior for errors that represent timeouts
type TimeoutChecker interface {
	Timeout() bool
}

// NotFoundChecker defines behavior for errors that represent "not found" conditions
type NotFoundChecker interface {
	NotFound() bool
}

// ValidationChecker defines behavior for validation errors
type ValidationChecker interface {
	ValidationError() bool
}

// CheckErrorBehavior demonstrates checking for error behavior
func CheckErrorBehavior(err error) {
	// Unwrap all errors in the chain
	for err != nil {
		// Check for timeout behavior
		if timeoutErr, ok := err.(TimeoutChecker); ok && timeoutErr.Timeout() {
			fmt.Println("Timeout detected through behavior")
			break
		}
		
		// Check for not found behavior
		if notFoundErr, ok := err.(NotFoundChecker); ok && notFoundErr.NotFound() {
			fmt.Println("Not found condition detected through behavior")
			break
		}
		
		// Try to unwrap and continue checking
		unwrappedErr := errors.Unwrap(err)
		if unwrappedErr == err || unwrappedErr == nil {
			break
		}
		err = unwrappedErr
	}
}

// OpenFileWithBehaviorCheck demonstrates checking for error behavior with os package
func OpenFileWithBehaviorCheck(filename string) (*os.File, error) {
	file, err := os.Open(filename)
	if err != nil {
		// Use behavior checking functions from os package
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file %s does not exist: %w", filename, err)
		}
		
		if os.IsPermission(err) {
			return nil, fmt.Errorf("permission denied for file %s: %w", filename, err)
		}
		
		if os.IsTimeout(err) {
			return nil, fmt.Errorf("timeout opening file %s: %w", filename, err)
		}
		
		// Generic error case
		return nil, fmt.Errorf("unexpected error opening %s: %w", filename, err)
	}
	
	return file, nil
}
