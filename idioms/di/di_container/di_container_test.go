package di_container

import (
	"fmt"
	"testing"
)

// TestBasicRegistrationAndResolution tests registering and resolving basic services
func TestBasicRegistrationAndResolution(t *testing.T) {
	container := NewContainer()
	
	// Register a concrete service
	logger := &ConsoleLogger{}
	container.Register(logger)
	
	// Resolve the service by concrete type
	var resolvedLogger *ConsoleLogger
	err := container.Resolve(&resolvedLogger)
	if err != nil {
		t.Errorf("Error resolving concrete service: %v", err)
	}
	
	if resolvedLogger == nil {
		t.Error("Resolved logger is nil")
	}
}

// TestInterfaceRegistrationAndResolution tests registering and resolving interface types
func TestInterfaceRegistrationAndResolution(t *testing.T) {
	container := NewContainer()
	
	// Register a service for an interface
	var loggerInterface Logger
	logger := &ConsoleLogger{}
	
	err := container.RegisterInstance(&loggerInterface, logger)
	if err != nil {
		t.Errorf("Error registering instance: %v", err)
	}
	
	// Resolve the service by interface type
	var resolvedLogger Logger
	err = container.Resolve(&resolvedLogger)
	if err != nil {
		t.Errorf("Error resolving interface: %v", err)
	}
	
	if resolvedLogger == nil {
		t.Error("Resolved logger is nil")
	}
}

// TestFactoryRegistration tests registering and resolving factory functions
func TestFactoryRegistration(t *testing.T) {
	container := NewContainer()
	
	// Register a factory function
	var loggerInterface Logger
	err := container.RegisterFactory(&loggerInterface, LoggerFactory)
	if err != nil {
		t.Errorf("Error registering factory: %v", err)
	}
	
	// Resolve the service
	var resolvedLogger Logger
	err = container.Resolve(&resolvedLogger)
	if err != nil {
		t.Errorf("Error resolving factory service: %v", err)
	}
	
	if resolvedLogger == nil {
		t.Error("Resolved logger is nil")
	}
	
	// Factory should be called every time
	var anotherLogger Logger
	err = container.Resolve(&anotherLogger)
	if err != nil {
		t.Errorf("Error resolving factory service again: %v", err)
	}
	
	// Should be different instances
	if fmt.Sprintf("%p", resolvedLogger) == fmt.Sprintf("%p", anotherLogger) {
		t.Error("Factory should return new instances each time")
	}
}

// TestSingletonRegistration tests registering and resolving singleton services
func TestSingletonRegistration(t *testing.T) {
	container := NewContainer()
	
	// Counter to track constructor calls
	constructorCalls := 0
	
	// Register a singleton constructor
	var loggerInterface Logger
	err := container.RegisterSingleton(&loggerInterface, func() Logger {
		constructorCalls++
		return &ConsoleLogger{}
	})
	if err != nil {
		t.Errorf("Error registering singleton: %v", err)
	}
	
	// Resolve the service multiple times
	var logger1 Logger
	var logger2 Logger
	
	err = container.Resolve(&logger1)
	if err != nil {
		t.Errorf("Error resolving singleton: %v", err)
	}
	
	err = container.Resolve(&logger2)
	if err != nil {
		t.Errorf("Error resolving singleton again: %v", err)
	}
	
	// Constructor should be called only once
	if constructorCalls != 1 {
		t.Errorf("Expected constructor to be called once, got: %d", constructorCalls)
	}
	
	// Both references should be the same instance
	if fmt.Sprintf("%p", logger1) != fmt.Sprintf("%p", logger2) {
		t.Error("Singleton should return the same instance each time")
	}
}

// TestAutoWire tests automatic dependency injection into struct fields
func TestAutoWire(t *testing.T) {
	container := NewContainer()
	
	// Register services
	var loggerInterface Logger
	var repoInterface Repository
	
	container.RegisterInstance(&loggerInterface, &ConsoleLogger{})
	container.RegisterInstance(&repoInterface, NewMemoryRepository())
	container.Register(&ServiceConfig{Timeout: 30, BaseURL: "http://example.com"})
	
	// Create service struct with empty fields
	service := &UserService{}
	
	// Auto-wire dependencies
	err := container.AutoWire(service)
	if err != nil {
		t.Errorf("Error auto-wiring: %v", err)
	}
	
	// Check that fields were injected
	if service.Logger == nil {
		t.Error("Logger was not injected")
	}
	
	if service.Repo == nil {
		t.Error("Repository was not injected")
	}
	
	if service.Config == nil {
		t.Error("Config was not injected")
	}
	
	// Test actual functionality
	err = service.CreateUser("test-id", "Test User")
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}
	
	user, err := service.GetUser("test-id")
	if err != nil {
		t.Errorf("Error retrieving user: %v", err)
	}
	
	if user == nil {
		t.Error("User not found")
	}
}

// MockLogger is a test double for Logger
type MockLogger struct {
	LoggedMessages []string
}

func (m *MockLogger) Log(message string) {
	m.LoggedMessages = append(m.LoggedMessages, message)
}

// TestWithMocks tests dependency injection with mocks
func TestWithMocks(t *testing.T) {
	container := NewContainer()
	
	// Create mocks
	mockLogger := &MockLogger{LoggedMessages: []string{}}
	mockRepo := NewMemoryRepository()
	
	// Register mocks
	var loggerInterface Logger
	var repoInterface Repository
	
	container.RegisterInstance(&loggerInterface, mockLogger)
	container.RegisterInstance(&repoInterface, mockRepo)
	container.Register(&ServiceConfig{Timeout: 30, BaseURL: "http://example.com"})
	
	// Create service with mocks
	service := &UserService{}
	container.AutoWire(service)
	
	// Test functionality
	service.CreateUser("test1", "Test User")
	
	// Verify mock was called
	if len(mockLogger.LoggedMessages) == 0 {
		t.Error("Mock logger was not called")
	}
	
	// Verify repository was called
	user, _ := mockRepo.FindByID("test1")
	if user == nil {
		t.Error("User was not saved to repository")
	}
}

// TestInvalidRegistrations tests error handling for invalid registrations
func TestInvalidRegistrations(t *testing.T) {
	container := NewContainer()
	
	// Try to register non-interface type
	var notAnInterface string
	err := container.RegisterInstance(&notAnInterface, "value")
	if err == nil {
		t.Error("Expected error when registering non-interface type")
	}
	
	// Try to register implementation that doesn't implement interface
	var loggerInterface Logger
	err = container.RegisterInstance(&loggerInterface, "not a logger")
	if err == nil {
		t.Error("Expected error when registering implementation that doesn't implement interface")
	}
	
	// Try to register invalid factory
	err = container.RegisterFactory(&loggerInterface, "not a function")
	if err == nil {
		t.Error("Expected error when registering non-function as factory")
	}
}

// TestResolveErrors tests error handling for resolve failures
func TestResolveErrors(t *testing.T) {
	container := NewContainer()
	
	// Try to resolve unregistered service
	var logger Logger
	err := container.Resolve(&logger)
	if err == nil {
		t.Error("Expected error when resolving unregistered service")
	}
	
	// Try to resolve into non-pointer
	var notAPointer Logger
	err = container.Resolve(notAPointer)
	if err == nil {
		t.Error("Expected error when resolving into non-pointer")
	}
	
	// Try to resolve nil
	err = container.Resolve(nil)
	if err == nil {
		t.Error("Expected error when resolving nil")
	}
}
