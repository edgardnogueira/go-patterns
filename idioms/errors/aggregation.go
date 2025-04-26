package errors

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

// Error aggregation
// -----------------------------------------------------

// MultiError represents a collection of errors
type MultiError struct {
	Errors []error
}

// Error implements the error interface, returning all error messages concatenated
func (m *MultiError) Error() string {
	if len(m.Errors) == 0 {
		return ""
	}
	
	if len(m.Errors) == 1 {
		return m.Errors[0].Error()
	}
	
	var messages []string
	for i, err := range m.Errors {
		messages = append(messages, fmt.Sprintf("[%d] %v", i+1, err))
	}
	
	return fmt.Sprintf("%d errors occurred:\n%s", len(m.Errors), strings.Join(messages, "\n"))
}

// Add adds an error to the MultiError
func (m *MultiError) Add(err error) {
	if err != nil {
		m.Errors = append(m.Errors, err)
	}
}

// ErrorOrNil returns nil if there are no errors, or the MultiError otherwise
func (m *MultiError) ErrorOrNil() error {
	if len(m.Errors) == 0 {
		return nil
	}
	
	return m
}

// Simple error aggregation
// -----------------------------------------------------

// CombineErrors combines multiple errors into a single error
func CombineErrors(errs ...error) error {
	// Filter out nil errors
	var nonNil []error
	for _, err := range errs {
		if err != nil {
			nonNil = append(nonNil, err)
		}
	}
	
	// Return nil if there are no non-nil errors
	if len(nonNil) == 0 {
		return nil
	}
	
	// Return the single error if there's only one
	if len(nonNil) == 1 {
		return nonNil[0]
	}
	
	// Create a MultiError for multiple errors
	return &MultiError{Errors: nonNil}
}

// Thread-safe error aggregation
// -----------------------------------------------------

// ErrorCollector collects errors in a thread-safe manner
type ErrorCollector struct {
	errors []error
	mu     sync.Mutex
}

// Add adds an error to the collector
func (c *ErrorCollector) Add(err error) {
	if err == nil {
		return
	}
	
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.errors = append(c.errors, err)
}

// Error implements the error interface
func (c *ErrorCollector) Error() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if len(c.errors) == 0 {
		return ""
	}
	
	if len(c.errors) == 1 {
		return c.errors[0].Error()
	}
	
	var messages []string
	for i, err := range c.errors {
		messages = append(messages, fmt.Sprintf("[%d] %v", i+1, err))
	}
	
	return fmt.Sprintf("%d errors occurred:\n%s", len(c.errors), strings.Join(messages, "\n"))
}

// ErrorOrNil returns nil if there are no errors, or the ErrorCollector otherwise
func (c *ErrorCollector) ErrorOrNil() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if len(c.errors) == 0 {
		return nil
	}
	
	return c
}

// Errors returns a copy of the collected errors
func (c *ErrorCollector) Errors() []error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Make a copy to prevent external modification
	result := make([]error, len(c.errors))
	copy(result, c.errors)
	
	return result
}

// Go1.20 errors.Join
// -----------------------------------------------------

// JoinErrors joins multiple errors into a single error.
// This is similar to errors.Join in Go 1.20+
func JoinErrors(errs ...error) error {
	// Filter out nil errors
	var nonNil []error
	for _, err := range errs {
		if err != nil {
			nonNil = append(nonNil, err)
		}
	}
	
	// Return nil if there are no non-nil errors
	if len(nonNil) == 0 {
		return nil
	}
	
	// In Go 1.20+, you would use: return errors.Join(nonNil...)
	// Our implementation uses MultiError
	return &MultiError{Errors: nonNil}
}

// Error wrapping with multiple errors
// -----------------------------------------------------

// WrapMultipleErrors wraps multiple errors with context
func WrapMultipleErrors(message string, errs ...error) error {
	// Join the errors
	err := JoinErrors(errs...)
	if err == nil {
		return nil
	}
	
	// Wrap with context
	return fmt.Errorf("%s: %w", message, err)
}

// Practical examples
// -----------------------------------------------------

// ValidateUserData demonstrates collecting multiple validation errors
func ValidateUserData(user *User) error {
	var errs MultiError
	
	// Check various fields
	if user.Name == "" {
		errs.Add(errors.New("name cannot be empty"))
	}
	
	if len(user.Name) < 3 {
		errs.Add(errors.New("name must be at least 3 characters"))
	}
	
	if !isValidEmail(user.Email) {
		errs.Add(errors.New("invalid email format"))
	}
	
	if user.Age < 18 {
		errs.Add(errors.New("must be at least 18 years old"))
	}
	
	return errs.ErrorOrNil()
}

// ProcessMultipleItems demonstrates collecting errors from multiple operations
func ProcessMultipleItems(items []string) error {
	var collector ErrorCollector
	var wg sync.WaitGroup
	
	for _, item := range items {
		wg.Add(1)
		
		go func(item string) {
			defer wg.Done()
			
			// Process the item
			err := processItem(item)
			if err != nil {
				collector.Add(fmt.Errorf("failed to process %s: %w", item, err))
			}
		}(item)
	}
	
	// Wait for all goroutines to complete
	wg.Wait()
	
	// Return error or nil
	return collector.ErrorOrNil()
}

// BatchOperations demonstrates an all-or-nothing approach with error collection
func BatchOperations(operations []func() error) error {
	var errs MultiError
	
	// Attempt all operations
	for i, op := range operations {
		if err := op(); err != nil {
			errs.Add(fmt.Errorf("operation %d failed: %w", i+1, err))
		}
	}
	
	// If any operation failed, return all errors
	if err := errs.ErrorOrNil(); err != nil {
		// In a real implementation, you might also roll back successful operations
		return fmt.Errorf("batch failed: %w", err)
	}
	
	return nil
}

// Helper functions for examples
// -----------------------------------------------------

// Include the email validation function from earlier
func isValidEmail(email string) bool {
	// Basic validation for example purposes
	return email != "" && containsChar(email, '@')
}

// Add a helper function for ProcessMultipleItems
func processItem(item string) error {
	// Simulate processing
	switch item {
	case "invalid":
		return errors.New("invalid item")
	case "not-found":
		return ErrNotFound
	case "permission":
		return ErrPermissionDenied
	default:
		return nil
	}
}

// User is extended for validation examples
type UserWithAge struct {
	User
	Age int
}
