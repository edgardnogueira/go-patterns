package errors

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// Error wrapping and unwrapping (requires Go 1.13+)
// -----------------------------------------------------

// WrapError demonstrates how to wrap errors using fmt.Errorf with %w
func WrapError(err error) error {
	// Simple error wrapping with context
	if err != nil {
		return fmt.Errorf("during processing: %w", err)
	}
	return nil
}

// NestedWrapError demonstrates multiple levels of error wrapping
func NestedWrapError() error {
	// Start with a basic error
	baseErr := errors.New("original error")
	
	// Wrap it with first layer of context
	firstWrap := fmt.Errorf("during data validation: %w", baseErr)
	
	// Wrap again with second layer of context
	secondWrap := fmt.Errorf("while processing user input: %w", firstWrap)
	
	// Wrap with a final layer
	return fmt.Errorf("transaction failed: %w", secondWrap)
}

// LoadConfig demonstrates practical use of error wrapping
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		// Wrap the file error with context about what we were trying to do
		return nil, fmt.Errorf("failed to load config from %s: %w", filename, err)
	}
	
	config, err := parseConfig(data)
	if err != nil {
		// Wrap any parsing errors
		return nil, fmt.Errorf("invalid config format in %s: %w", filename, err)
	}
	
	return config, nil
}

// Config is a simple example struct
type Config struct {
	Value string
}

// Helper function for config example
func parseConfig(data []byte) (*Config, error) {
	if len(data) == 0 {
		return nil, errors.New("empty config data")
	}
	return &Config{Value: string(data)}, nil
}

// Error unwrapping (errors.Unwrap, errors.Is, errors.As)
// -----------------------------------------------------

// UnwrapExample demonstrates using errors.Unwrap to manually unwrap errors
func UnwrapExample(err error) {
	fmt.Println("Original error:", err)
	
	// Unwrap once
	unwrapped := errors.Unwrap(err)
	if unwrapped != nil {
		fmt.Println("After first unwrap:", unwrapped)
		
		// Unwrap again if possible
		deeperUnwrapped := errors.Unwrap(unwrapped)
		if deeperUnwrapped != nil {
			fmt.Println("After second unwrap:", deeperUnwrapped)
		}
	}
}

// IsExample demonstrates using errors.Is to check for specific errors in a chain
func IsExample() {
	// Try to open a non-existent file
	_, err := os.Open("non-existent-file.txt")
	if err != nil {
		// Wrap the error
		wrappedErr := fmt.Errorf("could not open config file: %w", err)
		
		// Wrap it again
		doubleWrapped := fmt.Errorf("configuration error: %w", wrappedErr)
		
		// Check if the error chain contains os.ErrNotExist
		if errors.Is(doubleWrapped, os.ErrNotExist) {
			fmt.Println("The file does not exist, even though we checked the wrapped error")
		}
		
		// Using the traditional method would fail:
		if doubleWrapped == os.ErrNotExist { // This will be false
			fmt.Println("This won't be printed because direct comparison doesn't work for wrapped errors")
		}
	}
}

// AsExample demonstrates using errors.As for type assertion on wrapped errors
func AsExample() {
	// Call a function that returns a wrapped PathError
	err := findConfigFile("config.txt")
	if err != nil {
		// Use errors.As to extract the specific error type from the chain
		var pathErr *fs.PathError
		if errors.As(err, &pathErr) {
			fmt.Printf("Path error details - Op: %s, Path: %s, Err: %v\n", 
				pathErr.Op, pathErr.Path, pathErr.Err)
		} else {
			fmt.Println("Not a path error:", err)
		}
	}
}

// Helper function for AsExample
func findConfigFile(name string) error {
	// This will return a *fs.PathError error
	_, err := os.Open(name)
	if err != nil {
		// Wrap the error
		return fmt.Errorf("config error: %w", err)
	}
	return nil
}

// TraverseErrorChain walks through the entire chain of wrapped errors
func TraverseErrorChain(err error) {
	fmt.Println("Error chain:")
	
	// Starting with the outermost error
	currentErr := err
	level := 1
	
	// Continue until we can't unwrap anymore
	for currentErr != nil {
		fmt.Printf("Level %d: %v\n", level, currentErr)
		
		// Move to the next error in the chain
		currentErr = errors.Unwrap(currentErr)
		level++
	}
}

// FindAllFilesWithRetry demonstrates practical error handling with wrapping
func FindAllFilesWithRetry(dir string, maxRetries int) ([]string, error) {
	var files []string
	var lastErr error
	
	// Try multiple times
	for attempt := 0; attempt < maxRetries; attempt++ {
		files, lastErr = findAllFiles(dir)
		if lastErr == nil {
			return files, nil
		}
		
		// Check if this is a permission error (which won't be fixed by retrying)
		var pathErr *fs.PathError
		if errors.As(lastErr, &pathErr) && errors.Is(pathErr.Err, os.ErrPermission) {
			return nil, fmt.Errorf("permission denied to access %s: %w", dir, lastErr)
		}
		
		// For other errors, we'll retry
	}
	
	// If we get here, all attempts failed
	return nil, fmt.Errorf("failed to list files after %d attempts: %w", maxRetries, lastErr)
}

// Helper function for FindAllFilesWithRetry
func findAllFiles(dir string) ([]string, error) {
	var files []string
	
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("error scanning directory %s: %w", dir, err)
	}
	
	return files, nil
}
