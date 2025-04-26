package errors

import (
	"errors"
	"fmt"
)

// Sentinel errors (predefined error values)
// -----------------------------------------------------

// Sentinel errors for expected error conditions
var (
	// ErrNotFound indicates a resource couldn't be found
	ErrNotFound = errors.New("resource not found")

	// ErrPermissionDenied indicates lack of permissions for an operation
	ErrPermissionDenied = errors.New("permission denied")

	// ErrInvalidInput indicates invalid input data
	ErrInvalidInput = errors.New("invalid input")

	// ErrTimeout indicates an operation timed out
	ErrTimeout = errors.New("operation timed out")

	// ErrAlreadyExists indicates an attempt to create a duplicate resource
	ErrAlreadyExists = errors.New("resource already exists")

	// ErrUnavailable indicates a service or resource is currently unavailable
	ErrUnavailable = errors.New("service unavailable")
)

// DatabaseError represents various database-related error codes
type DatabaseErrorCode int

// Define specific database error codes
const (
	DbErrUnknown DatabaseErrorCode = iota
	DbErrConnection
	DbErrDuplicate
	DbErrTransaction
	DbErrConstraint
)

// DatabaseSpecificErrors maps error codes to sentinel errors
var DatabaseSpecificErrors = map[DatabaseErrorCode]error{
	DbErrConnection:  errors.New("database connection error"),
	DbErrDuplicate:   errors.New("duplicate key in database"),
	DbErrTransaction: errors.New("transaction error"),
	DbErrConstraint:  errors.New("constraint violation"),
}

// Using sentinel errors
// -----------------------------------------------------

// CheckUserAccess demonstrates using sentinel errors for access control
func CheckUserAccess(userID string, resourceID string) error {
	// Check if user exists
	exists, err := userExists(userID)
	if err != nil {
		return fmt.Errorf("error checking user: %w", err)
	}
	if !exists {
		// Return the sentinel error directly
		return ErrNotFound
	}

	// Check if user has access
	hasAccess, err := hasPermission(userID, resourceID)
	if err != nil {
		return fmt.Errorf("error checking permissions: %w", err)
	}
	if !hasAccess {
		// Return the sentinel error directly
		return ErrPermissionDenied
	}

	return nil
}

// GetResource demonstrates comparing against sentinel errors
func GetResource(id string) (*Resource, error) {
	resource, err := findResource(id)
	
	// Compare directly to sentinel error
	if err == ErrNotFound {
		// Handle the specific not found case
		return nil, fmt.Errorf("resource %s does not exist: %w", id, err)
	}
	
	if err != nil {
		// Handle other errors
		return nil, fmt.Errorf("failed to retrieve resource %s: %w", id, err)
	}
	
	return resource, nil
}

// DatabaseOperation demonstrates using domain-specific sentinel errors
func DatabaseOperation(id string) error {
	// Simulate a database operation
	err := simulateDbOperation(id)
	
	// Check for known database errors
	if err == DatabaseSpecificErrors[DbErrConnection] {
		return fmt.Errorf("database unavailable, try again later: %w", err)
	}
	
	if err == DatabaseSpecificErrors[DbErrDuplicate] {
		return fmt.Errorf("item %s already exists: %w", id, err)
	}
	
	if err != nil {
		return fmt.Errorf("database operation failed: %w", err)
	}
	
	return nil
}

// With errors.Is (Go 1.13+)
// -----------------------------------------------------

// ProcessResource demonstrates using errors.Is with sentinel errors
func ProcessResource(id string) error {
	// Try to update the resource
	err := updateResource(id)
	if err != nil {
		// Use errors.Is to check for sentinel errors in a wrapped error chain
		if errors.Is(err, ErrNotFound) {
			return fmt.Errorf("cannot process non-existent resource %s: %w", id, err)
		}
		
		if errors.Is(err, ErrPermissionDenied) {
			return fmt.Errorf("you don't have permission to process %s: %w", id, err)
		}
		
		if errors.Is(err, ErrTimeout) {
			return fmt.Errorf("processing timed out for %s, please try again: %w", id, err)
		}
		
		// Generic error case
		return fmt.Errorf("failed to process resource %s: %w", id, err)
	}
	
	return nil
}

// Multiple error checking
// -----------------------------------------------------

// CreateUser demonstrates checking multiple potential error conditions
func CreateUser(username, email string) (*User, error) {
	// Validate input
	if username == "" || email == "" {
		return nil, ErrInvalidInput
	}
	
	// Check if username already exists
	exists, err := checkUsernameExists(username)
	if err != nil {
		return nil, fmt.Errorf("error checking username: %w", err)
	}
	if exists {
		return nil, ErrAlreadyExists
	}
	
	// Try to create the user
	user, err := saveUser(username, email)
	if err != nil {
		// Check for specific database errors
		if errors.Is(err, DatabaseSpecificErrors[DbErrConnection]) {
			return nil, ErrUnavailable
		}
		
		// Fall back to a generic error
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	
	return user, nil
}

// Helper functions for examples
// -----------------------------------------------------

func userExists(userID string) (bool, error) {
	if userID == "invalid" {
		return false, errors.New("database error")
	}
	return userID != "unknown", nil
}

func hasPermission(userID, resourceID string) (bool, error) {
	if userID == "invalid" || resourceID == "invalid" {
		return false, errors.New("error checking permissions")
	}
	return userID == "admin", nil
}

func findResource(id string) (*Resource, error) {
	switch id {
	case "not-found":
		return nil, ErrNotFound
	case "error":
		return nil, errors.New("unexpected error")
	default:
		return &Resource{ID: id, Name: "Example Resource"}, nil
	}
}

func simulateDbOperation(id string) error {
	switch id {
	case "connection-error":
		return DatabaseSpecificErrors[DbErrConnection]
	case "duplicate":
		return DatabaseSpecificErrors[DbErrDuplicate]
	case "transaction-error":
		return DatabaseSpecificErrors[DbErrTransaction]
	case "error":
		return errors.New("unknown database error")
	default:
		return nil
	}
}

func updateResource(id string) error {
	switch id {
	case "not-found":
		return fmt.Errorf("update failed: %w", ErrNotFound)
	case "no-permission":
		return fmt.Errorf("update failed: %w", ErrPermissionDenied)
	case "timeout":
		return fmt.Errorf("update failed: %w", ErrTimeout)
	case "error":
		return errors.New("unknown error during update")
	default:
		return nil
	}
}

func checkUsernameExists(username string) (bool, error) {
	switch username {
	case "existing":
		return true, nil
	case "error":
		return false, errors.New("database error")
	default:
		return false, nil
	}
}

func saveUser(username, email string) (*User, error) {
	switch username {
	case "db-error":
		return nil, DatabaseSpecificErrors[DbErrConnection]
	case "error":
		return nil, errors.New("unknown error")
	default:
		return &User{ID: "new-id", Name: username}, nil
	}
}
