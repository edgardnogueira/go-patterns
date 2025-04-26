package errors

import (
	"errors"
	"fmt"
	"io"
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
	
	// Make a copy to avoid race conditions
	result := make([]error, len(c.errors))
	copy(result, c.errors)
	
	return result
}

// Unwrap implementation for Go 1.20+ error joining
// -----------------------------------------------------

// Unwrap for MultiError returns the slice of contained errors
// This is compatible with Go 1.20's errors.Join approach
func (m *MultiError) Unwrap() []error {
	return m.Errors
}

// Unwrap for ErrorCollector returns the slice of contained errors
func (c *ErrorCollector) Unwrap() []error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Make a copy to avoid race conditions
	result := make([]error, len(c.errors))
	copy(result, c.errors)
	
	return result
}

// Using Go 1.20's errors.Join (if available)
// -----------------------------------------------------

// JoinErrors is a wrapper around errors.Join (for Go 1.20+)
// For earlier Go versions, it falls back to our implementation
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
	
	// Return the single error if there's only one
	if len(nonNil) == 1 {
		return nonNil[0]
	}
	
	// Use errors.Join if running on Go 1.20+
	return errors.Join(nonNil...)
}

// BatchOperations executes multiple operations and aggregates errors
func BatchOperations(operations []func() error) error {
	var errs MultiError
	
	for i, op := range operations {
		if err := op(); err != nil {
			errs.Add(fmt.Errorf("operation %d failed: %w", i+1, err))
		}
	}
	
	return errs.ErrorOrNil()
}

// Practical examples
// -----------------------------------------------------

// ValidateUserData demonstrates using MultiError for validation
func ValidateUserData(user *User) error {
	var errs MultiError
	
	if user == nil {
		return errors.New("user cannot be nil")
	}
	
	if user.ID == "" {
		errs.Add(errors.New("user ID cannot be empty"))
	}
	
	if user.Name == "" {
		errs.Add(errors.New("user name cannot be empty"))
	}
	
	if !isValidEmail(user.Email) {
		errs.Add(errors.New("invalid email format"))
	}
	
	return errs.ErrorOrNil()
}

// ProcessBatch demonstrates collecting errors from multiple operations
func ProcessBatch(items []string) error {
	var collector ErrorCollector
	var wg sync.WaitGroup
	
	for _, item := range items {
		wg.Add(1)
		go func(item string) {
			defer wg.Done()
			
			err := processItem(item)
			if err != nil {
				collector.Add(fmt.Errorf("failed to process %q: %w", item, err))
			}
		}(item)
	}
	
	wg.Wait()
	return collector.ErrorOrNil()
}

// CleanupResources demonstrates aggregating errors from cleanup operations
func CleanupResources(resources []io.Closer) error {
	var errs []error
	
	for _, r := range resources {
		err := r.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}
	
	return CombineErrors(errs...)
}

// ApplyMigrations demonstrates sequential error collection
func ApplyMigrations(migrations []Migration) error {
	var errs MultiError
	
	for i, migration := range migrations {
		if err := migration.Apply(); err != nil {
			errs.Add(fmt.Errorf("migration %d failed: %w", i+1, err))
			
			// Stop at the first error for migrations
			return errs.ErrorOrNil()
		}
	}
	
	return nil
}

// Helper types and functions for examples
// -----------------------------------------------------

// Migration is a simple example struct for ApplyMigrations
type Migration struct {
	Name string
	SQL  string
}

// Apply simulates applying a migration
func (m Migration) Apply() error {
	// Simulate errors for certain SQL
	if m.SQL == "INVALID" {
		return errors.New("syntax error in migration")
	}
	return nil
}

// Email validation helper
func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// Process an item (used in examples)
func processItem(item string) error {
	// Simulate errors for certain items
	switch item {
	case "error":
		return errors.New("processing failed")
	case "invalid":
		return errors.New("invalid item")
	default:
		return nil
	}
}

// User is a simple struct for user data validation examples
type User struct {
	ID    string
	Name  string
	Email string
}
