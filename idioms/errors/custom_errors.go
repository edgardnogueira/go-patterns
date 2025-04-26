package errors

import (
	"fmt"
	"time"
)

// Custom error types
// -----------------------------------------------------

// ValidationError represents an error that occurs during data validation
type ValidationError struct {
	Field   string // The field that failed validation
	Message string // The validation error message
}

// Error implements the error interface for ValidationError
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed on field %s: %s", e.Field, e.Message)
}

// IsValidationError checks if an error is a ValidationError
func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}

// NewValidationError creates a new ValidationError
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// ResourceError represents an error related to a specific resource
type ResourceError struct {
	ResourceID string    // ID of the resource that caused the error
	Operation  string    // Operation being performed (e.g., "create", "read", "update", "delete")
	Message    string    // Error message
	Timestamp  time.Time // When the error occurred
	Cause      error     // Underlying error (if any)
}

// Error implements the error interface for ResourceError
func (e *ResourceError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s operation failed on resource %s at %s: %s (caused by: %v)",
			e.Operation, e.ResourceID, e.Timestamp.Format(time.RFC3339), e.Message, e.Cause)
	}
	return fmt.Sprintf("%s operation failed on resource %s at %s: %s",
		e.Operation, e.ResourceID, e.Timestamp.Format(time.RFC3339), e.Message)
}

// NewResourceError creates a new ResourceError
func NewResourceError(resourceID, operation, message string, cause error) *ResourceError {
	return &ResourceError{
		ResourceID: resourceID,
		Operation:  operation,
		Message:    message,
		Timestamp:  time.Now(),
		Cause:      cause,
	}
}

// Unwrap returns the underlying cause of the error
func (e *ResourceError) Unwrap() error {
	return e.Cause
}

// Example usage of custom error types
// -----------------------------------------------------

// ValidateUser demonstrates using custom validation errors
func ValidateUser(name, email string) error {
	if name == "" {
		return NewValidationError("name", "name cannot be empty")
	}
	
	if email == "" {
		return NewValidationError("email", "email cannot be empty")
	}
	
	// Very basic email validation for example purposes
	if !containsChar(email, '@') {
		return NewValidationError("email", "email must contain @ symbol")
	}
	
	return nil
}

// FetchResource demonstrates using resource errors with context
func FetchResource(id string) (*Resource, error) {
	// Simulate a database lookup
	resource, err := simulateDBLookup(id)
	if err != nil {
		// Wrap the underlying error with additional context
		return nil, NewResourceError(id, "fetch", "database lookup failed", err)
	}
	
	return resource, nil
}

// Helper function for validation examples
func containsChar(s string, c byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return true
		}
	}
	return false
}

// Resource is a simple example struct
type Resource struct {
	ID   string
	Name string
}

// Simulate a database lookup (for example purposes)
func simulateDBLookup(id string) (*Resource, error) {
	if id == "" {
		return nil, fmt.Errorf("invalid id")
	}
	if id == "not-found" {
		return nil, fmt.Errorf("resource not found")
	}
	
	// Successful case
	return &Resource{ID: id, Name: "Example Resource"}, nil
}
