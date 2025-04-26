// Package functional_injection demonstrates dependency injection using higher-order functions.
package functional_injection

import (
	"fmt"
	"net/http"
	"time"
)

// Logger defines a simple logging interface
type Logger interface {
	Log(message string)
}

// SimpleLogger is a basic logger implementation
type SimpleLogger struct{}

func (l *SimpleLogger) Log(message string) {
	fmt.Printf("[%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), message)
}

// AuthChecker verifies if a user is authorized
type AuthChecker func(userID string, action string) bool

// Cache provides data caching functionality
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, expiration time.Duration)
}

// MetricsRecorder records performance metrics
type MetricsRecorder func(operation string, duration time.Duration)

// Database represents a database connection
type Database interface {
	Query(query string, args ...interface{}) ([]map[string]interface{}, error)
	Execute(statement string, args ...interface{}) (int64, error)
}

// WithLogging adds logging to a function
func WithLogging(fn func(string) (string, error), logger Logger) func(string) (string, error) {
	return func(input string) (string, error) {
		logger.Log(fmt.Sprintf("Function called with input: %s", input))
		
		startTime := time.Now()
		result, err := fn(input)
		duration := time.Since(startTime)
		
		if err != nil {
			logger.Log(fmt.Sprintf("Function failed after %v: %v", duration, err))
		} else {
			logger.Log(fmt.Sprintf("Function succeeded in %v with result: %s", duration, result))
		}
		
		return result, err
	}
}

// WithCaching adds caching to a function
func WithCaching(fn func(string) (string, error), cache Cache) func(string) (string, error) {
	return func(input string) (string, error) {
		// Check cache first
		cacheKey := fmt.Sprintf("result:%s", input)
		if cachedValue, found := cache.Get(cacheKey); found {
			return cachedValue.(string), nil
		}
		
		// Call the original function
		result, err := fn(input)
		if err == nil {
			// Cache the result on success
			cache.Set(cacheKey, result, 5*time.Minute)
		}
		
		return result, err
	}
}

// WithRetry adds retry logic to a function
func WithRetry(fn func(string) (string, error), maxRetries int, delay time.Duration) func(string) (string, error) {
	return func(input string) (string, error) {
		var lastErr error
		
		for attempt := 0; attempt <= maxRetries; attempt++ {
			if attempt > 0 {
				time.Sleep(delay)
			}
			
			result, err := fn(input)
			if err == nil {
				return result, nil
			}
			
			lastErr = err
		}
		
		return "", fmt.Errorf("function failed after %d retries: %w", maxRetries, lastErr)
	}
}

// WithAuth adds authorization checking to an HTTP handler
func WithAuth(handler http.HandlerFunc, authChecker AuthChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("User-ID")
		action := r.URL.Path
		
		if userID == "" {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}
		
		if !authChecker(userID, action) {
			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		}
		
		handler(w, r)
	}
}

// WithMetrics adds performance metrics to an HTTP handler
func WithMetrics(handler http.HandlerFunc, recorder MetricsRecorder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		handler(w, r)
		duration := time.Since(startTime)
		recorder(r.URL.Path, duration)
	}
}

// WithTransaction adds database transaction handling to a function
func WithTransaction(fn func(db Database) error, db Database, logger Logger) func() error {
	return func() error {
		// In a real implementation, this would start and manage a transaction
		logger.Log("Starting database transaction")
		
		err := fn(db)
		
		if err != nil {
			logger.Log(fmt.Sprintf("Transaction failed: %v", err))
			// In a real implementation, this would roll back the transaction
			logger.Log("Rolling back transaction")
			return err
		}
		
		// In a real implementation, this would commit the transaction
		logger.Log("Committing transaction")
		return nil
	}
}

// Example usage functions
func ProcessData(input string) (string, error) {
	return fmt.Sprintf("Processed: %s", input), nil
}

// CreateUserHandler is an HTTP handler for user creation
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "User created")
}

// SaveUserData simulates saving user data to a database
func SaveUserData(db Database) error {
	_, err := db.Execute("INSERT INTO users (name, email) VALUES (?, ?)", "John Doe", "john@example.com")
	return err
}

// CreateUserServiceWithDI demonstrates creating a service with multiple functional middlewares
func CreateUserServiceWithDI(logger Logger, cache Cache, authChecker AuthChecker, metricsRecorder MetricsRecorder) http.HandlerFunc {
	// Start with the basic handler
	handler := CreateUserHandler
	
	// Add middleware layers (dependencies) in the desired order
	handler = WithAuth(handler, authChecker)
	handler = WithMetrics(handler, metricsRecorder)
	
	return handler
}
