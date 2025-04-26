package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/edgardnogueira/go-patterns/idioms/errors"
)

func main() {
	fmt.Println("=== Go Error Handling Patterns ===")
	fmt.Println()

	// Basic error handling
	fmt.Println("=== Basic Error Handling ===")
	errors.BasicErrorCreation()
	err := errors.BasicErrorHandling()
	fmt.Printf("Basic error handling result: %v\n", err)
	
	// Error handling with functions
	fmt.Println()
	fmt.Println("=== Error Handling with Functions ===")
	errors.CalculateAndPrint(10, 2)  // Should succeed
	errors.CalculateAndPrint(10, 0)  // Should fail with division by zero
	
	// Custom error types
	fmt.Println()
	fmt.Println("=== Custom Error Types ===")
	err = errors.ValidateUser("", "test@example.com")
	fmt.Printf("Validation error: %v\n", err)
	
	_, err = errors.FetchResource("not-found")
	fmt.Printf("Resource error: %v\n", err)
	
	// Error wrapping and unwrapping
	fmt.Println()
	fmt.Println("=== Error Wrapping and Unwrapping ===")
	nestedErr := errors.NestedWrapError()
	fmt.Printf("Nested error: %v\n", nestedErr)
	
	fmt.Println("Unwrapping the nested error:")
	errors.UnwrapExample(nestedErr)
	
	fmt.Println("Using errors.Is to check for specific errors:")
	errors.IsExample()
	
	fmt.Println("Using errors.As for type assertion:")
	errors.AsExample()
	
	// Type assertions with errors
	fmt.Println()
	fmt.Println("=== Type Assertion with Errors ===")
	timeoutErr := &errors.TimeoutError{
		Operation: "database query",
		Duration:  time.Second * 5,
	}
	wrappedTimeoutErr := fmt.Errorf("error occurred: %w", timeoutErr)
	
	fmt.Println("Traditional type assertion:")
	errors.TypeAssertionExample(timeoutErr)
	
	fmt.Println("Using errors.As:")
	errors.ErrorsAsExample(wrappedTimeoutErr)
	
	fmt.Println("Handling different error types in a user lookup:")
	errors.HandleUserLookup("not-found")
	
	// Sentinel errors
	fmt.Println()
	fmt.Println("=== Sentinel Errors ===")
	err = errors.CheckUserAccess("unknown", "resource1")
	fmt.Printf("Access check result: %v\n", err)
	
	_, err = errors.GetResource("not-found")
	fmt.Printf("Get resource result: %v\n", err)
	
	err = errors.ProcessResource("not-found")
	fmt.Printf("Process resource result: %v\n", err)
	
	// Concurrent error handling
	fmt.Println()
	fmt.Println("=== Concurrent Error Handling ===")
	urls := []string{"example.com", "error", "timeout", "not-found"}
	errs := errors.ConcurrentErrors(urls)
	fmt.Println("Results from concurrent operations:")
	for i, err := range errs {
		fmt.Printf("  Error %d: %v\n", i+1, err)
	}
	
	fmt.Println("First error only approach:")
	err = errors.FirstErrorOnly(urls)
	fmt.Printf("  First error: %v\n", err)
	
	fmt.Println("Using context for cancellation:")
	ctx := context.Background()
	tasks := []string{"task1", "error", "task3"}
	err = errors.WithContext(ctx, tasks)
	fmt.Printf("  Context result: %v\n", err)
	
	// Behavior-based error checking
	fmt.Println()
	fmt.Println("=== Behavior-Based Error Checking ===")
	err = errors.HandleConnection("timeout.example.com")
	fmt.Printf("Connection result: %v\n", err)
	
	err = errors.QueryDatabase("SELECT * FROM nonexistent_table")
	fmt.Printf("Query result: %v\n", err)
	
	// Error aggregation
	fmt.Println()
	fmt.Println("=== Error Aggregation ===")
	var multiErr errors.MultiError
	multiErr.Add(errors.ErrNotFound)
	multiErr.Add(errors.ErrPermissionDenied)
	fmt.Printf("Multiple errors: %v\n", multiErr.ErrorOrNil())
	
	// Combine different errors
	combinedErr := errors.CombineErrors(errors.ErrNotFound, nil, errors.ErrTimeout)
	fmt.Printf("Combined errors: %v\n", combinedErr)
	
	// Example of validating user data with multiple potential errors
	user := &errors.User{
		ID:    "123",
		Name:  "A",  // Too short
		Email: "invalid-email", // Invalid format
	}
	err = errors.ValidateUserData(user)
	fmt.Printf("User validation errors: %v\n", err)
	
	// Batch operations
	fmt.Println()
	fmt.Println("=== Batch Operations ===")
	operations := []func() error{
		func() error { return nil },  // Successful operation
		func() error { return errors.ErrNotFound },  // Failed operation
		func() error { return errors.ErrPermissionDenied },  // Failed operation
	}
	err = errors.BatchOperations(operations)
	fmt.Printf("Batch operations result: %v\n", err)
	
	fmt.Println()
	fmt.Println("=== Error Handling Patterns Demonstration Complete ===")
}
