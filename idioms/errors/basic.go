// Package errors demonstrates idiomatic Go error handling patterns
package errors

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// Simple error creation and checking
// -----------------------------------------------------

// BasicErrorCreation demonstrates the basic ways to create errors in Go
func BasicErrorCreation() {
	// Method 1: Using errors.New for simple static error messages
	err1 := errors.New("something went wrong")
	fmt.Println("errors.New:", err1)

	// Method 2: Using fmt.Errorf for formatted error messages
	value := 42
	err2 := fmt.Errorf("invalid value: %d", value)
	fmt.Println("fmt.Errorf:", err2)
}

// BasicErrorHandling demonstrates how to handle errors in Go
func BasicErrorHandling() error {
	// Attempt to open a non-existent file
	_, err := os.Open("non-existent-file.txt")
	if err != nil {
		// Most common Go error handling pattern:
		// Check immediately after the function call
		fmt.Println("Error occurred:", err)
		
		// Return early with the error
		return err
	}
	
	// Continue with normal flow if no error
	return nil
}

// BasicErrorChecking shows different ways to check errors
func BasicErrorChecking() {
	// Example function call that might return an error
	err := simulateError(true)
	
	// Pattern 1: Basic error checking
	if err != nil {
		fmt.Println("Error occurred:", err)
		return
	}
	
	// Pattern 2: Checking for specific error conditions
	file, err := os.Open("some-file.txt")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("File does not exist:", err)
		} else if os.IsPermission(err) {
			fmt.Println("Permission denied:", err)
		} else {
			fmt.Println("Unknown error:", err)
		}
		return
	}
	defer file.Close()
	
	// Pattern 3: Checking for EOF (common pattern in loops)
	data := make([]byte, 100)
	for {
		_, err := file.Read(data)
		if err == io.EOF {
			// End of file reached, not an error condition
			break
		}
		if err != nil {
			// Some other error occurred
			fmt.Println("Error reading file:", err)
			return
		}
		// Process data...
	}
}

// Returning errors from functions
// -----------------------------------------------------

// Division demonstrates returning errors from functions
func Division(a, b float64) (float64, error) {
	if b == 0 {
		// Return meaningful error for expected error condition
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// Typical error handling pattern using multiple return values
func CalculateAndPrint(a, b float64) {
	result, err := Division(a, b)
	if err != nil {
		fmt.Println("Error in calculation:", err)
		return
	}
	fmt.Printf("%f / %f = %f\n", a, b, result)
}

// Helper function for examples
func simulateError(shouldError bool) error {
	if shouldError {
		return errors.New("simulated error")
	}
	return nil
}
