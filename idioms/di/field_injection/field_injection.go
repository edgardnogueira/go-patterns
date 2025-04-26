// Package field_injection demonstrates dependency injection via fields after construction.
package field_injection

import (
	"fmt"
	"time"
)

// Logger defines a simple logging interface
type Logger interface {
	Log(message string)
}

// ConsoleLogger is a concrete implementation that logs to console
type ConsoleLogger struct{}

func (l *ConsoleLogger) Log(message string) {
	fmt.Printf("[%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), message)
}

// Config holds application configuration
type Config struct {
	DatabaseURL    string
	APIKey         string
	TimeoutSeconds int
}

// UserRepository handles user data operations
type UserRepository struct {
	// Dependencies set via field injection
	Logger     Logger
	Config     *Config
	Metrics    MetricsCollector // Optional dependency
	IsReadOnly bool             // Simple configuration field
}

// MetricsCollector tracks performance metrics
type MetricsCollector interface {
	RecordOperation(name string, durationMs int)
}

// GetUser fetches a user by ID
func (r *UserRepository) GetUser(id int) (string, error) {
	startTime := time.Now()
	
	// Log the operation
	if r.Logger != nil {
		r.Logger.Log(fmt.Sprintf("Getting user with ID: %d", id))
	}
	
	// Simulate database operation
	time.Sleep(10 * time.Millisecond)
	
	// Record metrics if available
	if r.Metrics != nil {
		duration := time.Since(startTime).Milliseconds()
		r.Metrics.RecordOperation("GetUser", int(duration))
	}
	
	// Return mock data
	return fmt.Sprintf("User%d", id), nil
}

// ApplicationBuilder uses field injection to construct complex objects with many dependencies
type ApplicationBuilder struct {
	app *Application
}

// Application is a complex object with many dependencies
type Application struct {
	Name           string
	Logger         Logger
	Repository     *UserRepository
	Config         *Config
	AdminUsers     []string
	EnabledFeatures map[string]bool
}

// NewApplicationBuilder creates a new application builder
func NewApplicationBuilder() *ApplicationBuilder {
	return &ApplicationBuilder{
		app: &Application{
			EnabledFeatures: make(map[string]bool),
		},
	}
}

// WithName sets the application name
func (b *ApplicationBuilder) WithName(name string) *ApplicationBuilder {
	b.app.Name = name
	return b
}

// WithLogger sets the logger
func (b *ApplicationBuilder) WithLogger(logger Logger) *ApplicationBuilder {
	b.app.Logger = logger
	return b
}

// WithRepository sets the user repository
func (b *ApplicationBuilder) WithRepository(repo *UserRepository) *ApplicationBuilder {
	b.app.Repository = repo
	return b
}

// WithConfig sets the configuration
func (b *ApplicationBuilder) WithConfig(config *Config) *ApplicationBuilder {
	b.app.Config = config
	return b
}

// WithAdminUsers sets the admin users
func (b *ApplicationBuilder) WithAdminUsers(users ...string) *ApplicationBuilder {
	b.app.AdminUsers = users
	return b
}

// WithEnabledFeature enables a feature
func (b *ApplicationBuilder) WithEnabledFeature(featureName string) *ApplicationBuilder {
	b.app.EnabledFeatures[featureName] = true
	return b
}

// Build returns the configured application
func (b *ApplicationBuilder) Build() *Application {
	// Validate required fields
	if b.app.Logger == nil {
		b.app.Logger = &ConsoleLogger{} // Default logger
	}
	
	return b.app
}

// Run starts the application
func (a *Application) Run() {
	a.Logger.Log(fmt.Sprintf("Starting application: %s", a.Name))
	
	// Application logic would go here
	
	a.Logger.Log("Application running")
}
