package field_injection

import (
	"strings"
	"sync"
	"testing"
)

// MockLogger implements Logger interface for testing
type MockLogger struct {
	logs []string
	mu   sync.Mutex
}

func (m *MockLogger) Log(message string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.logs = append(m.logs, message)
}

func (m *MockLogger) GetLogs() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.logs
}

// MockMetricsCollector implements MetricsCollector for testing
type MockMetricsCollector struct {
	operations map[string][]int
	mu         sync.Mutex
}

func NewMockMetricsCollector() *MockMetricsCollector {
	return &MockMetricsCollector{
		operations: make(map[string][]int),
	}
}

func (m *MockMetricsCollector) RecordOperation(name string, durationMs int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.operations[name] = append(m.operations[name], durationMs)
}

func (m *MockMetricsCollector) GetOperationCount(name string) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.operations[name])
}

func TestUserRepositoryFieldInjection(t *testing.T) {
	// Create a repository with nil dependencies
	emptyRepo := &UserRepository{}
	
	// Repository should work even with nil dependencies
	_, err := emptyRepo.GetUser(1)
	if err != nil {
		t.Errorf("Expected no error with nil dependencies, got %v", err)
	}
	
	// Create dependencies
	logger := &MockLogger{}
	metrics := NewMockMetricsCollector()
	config := &Config{
		DatabaseURL:    "postgresql://localhost:5432/mydb",
		APIKey:         "secret-key",
		TimeoutSeconds: 30,
	}
	
	// Inject dependencies via fields after creation
	repo := &UserRepository{}
	repo.Logger = logger
	repo.Config = config
	repo.Metrics = metrics
	repo.IsReadOnly = true
	
	// Use the repository
	_, err = repo.GetUser(42)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Verify logger was called
	logs := logger.GetLogs()
	if len(logs) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(logs))
	}
	
	if !strings.Contains(logs[0], "Getting user with ID: 42") {
		t.Errorf("Log message incorrect, got: %s", logs[0])
	}
	
	// Verify metrics were recorded
	opCount := metrics.GetOperationCount("GetUser")
	if opCount != 1 {
		t.Errorf("Expected 1 recorded operation, got %d", opCount)
	}
}

func TestApplicationBuilder(t *testing.T) {
	// Create dependencies
	logger := &MockLogger{}
	repo := &UserRepository{}
	config := &Config{
		DatabaseURL:    "postgresql://localhost:5432/mydb",
		APIKey:         "secret-key",
		TimeoutSeconds: 30,
	}
	
	// Use the builder pattern with field injection
	app := NewApplicationBuilder().
		WithName("TestApp").
		WithLogger(logger).
		WithRepository(repo).
		WithConfig(config).
		WithAdminUsers("admin1", "admin2").
		WithEnabledFeature("reporting").
		WithEnabledFeature("analytics").
		Build()
	
	// Verify the application was built correctly
	if app.Name != "TestApp" {
		t.Errorf("Expected name TestApp, got %s", app.Name)
	}
	
	if app.Logger != logger {
		t.Error("Logger was not correctly injected")
	}
	
	if app.Repository != repo {
		t.Error("Repository was not correctly injected")
	}
	
	if app.Config != config {
		t.Error("Config was not correctly injected")
	}
	
	if len(app.AdminUsers) != 2 {
		t.Errorf("Expected 2 admin users, got %d", len(app.AdminUsers))
	}
	
	if !app.EnabledFeatures["reporting"] || !app.EnabledFeatures["analytics"] {
		t.Error("Features were not correctly enabled")
	}
	
	// Test the application runs
	app.Run()
	
	// Verify logs from running the app
	logs := logger.GetLogs()
	if len(logs) != 2 {
		t.Errorf("Expected 2 log entries, got %d", len(logs))
	}
	
	if !strings.Contains(logs[0], "Starting application: TestApp") {
		t.Errorf("First log message incorrect, got: %s", logs[0])
	}
	
	if !strings.Contains(logs[1], "Application running") {
		t.Errorf("Second log message incorrect, got: %s", logs[1])
	}
}

func TestBuilderWithDefaultLogger(t *testing.T) {
	// Build without specifying logger
	app := NewApplicationBuilder().
		WithName("DefaultLoggerApp").
		Build()
	
	// Verify default logger was provided
	if app.Logger == nil {
		t.Error("Default logger was not provided")
	}
	
	// Run the app with default logger (just to ensure no panics)
	app.Run()
}
