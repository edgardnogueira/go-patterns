package errors

import (
	"fmt"
	"io"
	"time"
)

// Behavior-based error checking
// -----------------------------------------------------

// Error interfaces define behaviors that errors can implement
// These interfaces allow code to check what an error can do
// rather than what type it is

// Temporary indicates whether an error is temporary and may be retried
type Temporary interface {
	Temporary() bool
}

// Timeout indicates whether an error is related to a timeout
type Timeout interface {
	Timeout() bool
}

// NotFound indicates whether an error is because a resource wasn't found
type NotFound interface {
	NotFound() bool
}

// Unauthorized indicates whether an error is due to lack of authorization
type Unauthorized interface {
	Unauthorized() bool
}

// Behavioral errors
// -----------------------------------------------------

// NetworkError represents an error that occurred during network operations
type NetworkError struct {
	Op        string
	Addr      string
	Err       error
	IsTemp    bool
	IsTimeout bool
}

// Error implements the error interface
func (e *NetworkError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s %s: %v", e.Op, e.Addr, e.Err)
	}
	return fmt.Sprintf("%s %s", e.Op, e.Addr)
}

// Unwrap returns the underlying error
func (e *NetworkError) Unwrap() error {
	return e.Err
}

// Temporary implements the Temporary interface
func (e *NetworkError) Temporary() bool {
	return e.IsTemp
}

// Timeout implements the Timeout interface
func (e *NetworkError) Timeout() bool {
	return e.IsTimeout
}

// DBError represents an error that occurred during database operations
type DBError struct {
	Op        string
	Query     string
	Err       error
	NotExists bool
	NoAccess  bool
}

// Error implements the error interface
func (e *DBError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s query %q: %v", e.Op, e.Query, e.Err)
	}
	return fmt.Sprintf("%s query %q", e.Op, e.Query)
}

// Unwrap returns the underlying error
func (e *DBError) Unwrap() error {
	return e.Err
}

// NotFound implements the NotFound interface
func (e *DBError) NotFound() bool {
	return e.NotExists
}

// Unauthorized implements the Unauthorized interface
func (e *DBError) Unauthorized() bool {
	return e.NoAccess
}

// Behavior-based error handling
// -----------------------------------------------------

// IsTemporary checks if an error has the Temporary behavior
func IsTemporary(err error) bool {
	temp, ok := err.(Temporary)
	return ok && temp.Temporary()
}

// IsTimeout checks if an error has the Timeout behavior
func IsTimeout(err error) bool {
	timeout, ok := err.(Timeout)
	return ok && timeout.Timeout()
}

// IsNotFound checks if an error has the NotFound behavior
func IsNotFound(err error) bool {
	notFound, ok := err.(NotFound)
	return ok && notFound.NotFound()
}

// IsUnauthorized checks if an error has the Unauthorized behavior
func IsUnauthorized(err error) bool {
	unauthorized, ok := err.(Unauthorized)
	return ok && unauthorized.Unauthorized()
}

// Retry demonstrates using behavioral checks to implement retry logic
func Retry(op func() error, maxRetries int, shouldRetry func(error) bool) error {
	var err error
	
	for attempt := 0; attempt < maxRetries; attempt++ {
		err = op()
		if err == nil {
			return nil // Success
		}
		
		// Check if the error is retryable using the provided function
		if !shouldRetry(err) {
			return err // Non-retryable error, stop immediately
		}
		
		// Simple backoff - in a real implementation, use exponential backoff with jitter
		time.Sleep(time.Duration(attempt+1) * 100 * time.Millisecond)
	}
	
	return fmt.Errorf("failed after %d attempts: %w", maxRetries, err)
}

// RetryTemporary demonstrates using the Temporary behavior for retries
func RetryTemporary(op func() error, maxRetries int) error {
	return Retry(op, maxRetries, IsTemporary)
}

// Practical examples
// -----------------------------------------------------

// HandleConnection demonstrates behavior-based error handling for network operations
func HandleConnection(address string) error {
	// Try to connect to the address
	err := connectTo(address)
	if err != nil {
		// Check for different behaviors
		if IsTimeout(err) {
			return fmt.Errorf("connection timed out, check your network: %w", err)
		}
		
		if IsTemporary(err) {
			return fmt.Errorf("temporary connection issue, please try again: %w", err)
		}
		
		// Generic error case
		return fmt.Errorf("connection failed: %w", err)
	}
	
	return nil
}

// QueryDatabase demonstrates behavior-based error handling for database operations
func QueryDatabase(query string) error {
	// Try to execute the query
	err := executeQuery(query)
	if err != nil {
		// Check for different behaviors
		if IsNotFound(err) {
			return fmt.Errorf("no records found for query: %w", err)
		}
		
		if IsUnauthorized(err) {
			return fmt.Errorf("you don't have permission to execute this query: %w", err)
		}
		
		// Generic error case
		return fmt.Errorf("query execution failed: %w", err)
	}
	
	return nil
}

// Standard library examples
// -----------------------------------------------------

// ProcessReader demonstrates how the standard library uses behavior checking
func ProcessReader(r io.Reader) error {
	data := make([]byte, 100)
	for {
		_, err := r.Read(data)
		if err == io.EOF {
			// End of file is not an error, just a signal that we're done
			break
		}
		if err != nil {
			// A generic error occurred
			return fmt.Errorf("read error: %w", err)
		}
		
		// Process data...
	}
	
	return nil
}

// Helper functions for examples
// -----------------------------------------------------

func connectTo(address string) error {
	// Simulate different network errors
	switch address {
	case "timeout.example.com":
		return &NetworkError{
			Op:        "connect",
			Addr:      address,
			Err:       fmt.Errorf("connection timed out"),
			IsTimeout: true,
		}
	case "temp.example.com":
		return &NetworkError{
			Op:     "connect",
			Addr:   address,
			Err:    fmt.Errorf("connection reset"),
			IsTemp: true,
		}
	case "error.example.com":
		return &NetworkError{
			Op:   "connect",
			Addr: address,
			Err:  fmt.Errorf("connection refused"),
		}
	default:
		return nil
	}
}

func executeQuery(query string) error {
	// Simulate different database errors
	switch query {
	case "SELECT * FROM nonexistent_table":
		return &DBError{
			Op:        "query",
			Query:     query,
			Err:       fmt.Errorf("table not found"),
			NotExists: true,
		}
	case "UPDATE restricted_table SET value = 1":
		return &DBError{
			Op:       "query",
			Query:    query,
			Err:      fmt.Errorf("insufficient privileges"),
			NoAccess: true,
		}
	case "INVALID SQL":
		return &DBError{
			Op:    "query",
			Query: query,
			Err:   fmt.Errorf("syntax error"),
		}
	default:
		return nil
	}
}

// Comparison with try-catch
// -----------------------------------------------------

/*
Example try-catch in another language (e.g., Java):

try {
    connection = connect(address);
} catch (TimeoutException e) {
    log.error("Connection timed out");
    throw new ServiceUnavailableException("Network issue", e);
} catch (TemporaryNetworkException e) {
    log.warn("Temporary network issue");
    return retryLater();
} catch (Exception e) {
    log.error("Connection failed", e);
    throw new ServiceUnavailableException("Network issue", e);
}

Go equivalent using behavior:

err := connectTo(address)
if err != nil {
    if IsTimeout(err) {
        log.Printf("Connection timed out: %v", err)
        return fmt.Errorf("service unavailable: %w", err)
    }
    
    if IsTemporary(err) {
        log.Printf("Temporary network issue: %v", err)
        return retryLater()
    }
    
    log.Printf("Connection failed: %v", err)
    return fmt.Errorf("service unavailable: %w", err)
}
*/
