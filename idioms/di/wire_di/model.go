// Package wire_di demonstrates dependency injection using Google's Wire tool.
package wire_di

import "fmt"

// Message represents a message entity
type Message struct {
	ID      string
	Content string
	UserID  string
}

// User represents a user entity
type User struct {
	ID       string
	Username string
	Email    string
}

// Config represents application configuration
type Config struct {
	DatabaseURL  string
	APIKey       string
	Environment  string
	LogLevel     string
	MaxRetries   int
	FeatureFlags map[string]bool
}

// NewConfig creates a new configuration with default values
func NewConfig() *Config {
	return &Config{
		DatabaseURL:  "postgres://localhost:5432/appdb",
		APIKey:       "default-api-key",
		Environment:  "development",
		LogLevel:     "info",
		MaxRetries:   3,
		FeatureFlags: map[string]bool{"feature_1": true, "feature_2": false},
	}
}

// DatabaseConnection represents a database connection
type DatabaseConnection struct {
	URL string
}

// NewDatabaseConnection creates a new database connection
func NewDatabaseConnection(config *Config) (*DatabaseConnection, error) {
	// In a real application, this would establish a connection to the database
	// For this example, we'll just use the URL from the config
	return &DatabaseConnection{
		URL: config.DatabaseURL,
	}, nil
}

// Close closes the database connection
func (db *DatabaseConnection) Close() error {
	fmt.Println("Closing database connection:", db.URL)
	return nil
}

// Logger provides logging functionality
type Logger struct {
	Level string
}

// NewLogger creates a new logger
func NewLogger(config *Config) *Logger {
	return &Logger{
		Level: config.LogLevel,
	}
}

// Log logs a message
func (l *Logger) Log(message string) {
	fmt.Printf("[%s] %s\n", l.Level, message)
}

// APIClient provides a client for external API calls
type APIClient struct {
	APIKey string
	logger *Logger
}

// NewAPIClient creates a new API client
func NewAPIClient(config *Config, logger *Logger) *APIClient {
	return &APIClient{
		APIKey: config.APIKey,
		logger: logger,
	}
}

// Call makes an API call
func (a *APIClient) Call(endpoint string, payload interface{}) error {
	a.logger.Log(fmt.Sprintf("Making API call to %s with API key %s", endpoint, a.APIKey))
	// In a real implementation, this would make an actual API call
	return nil
}
