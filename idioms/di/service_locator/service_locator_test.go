package service_locator

import (
	"testing"
)

// TestServiceLocatorBasics tests the basic operations of the service locator
func TestServiceLocatorBasics(t *testing.T) {
	locator := NewServiceLocator()
	
	// Register a service
	logger := &ConsoleLogger{}
	locator.Register("logger", logger)
	
	// Check service exists
	if !locator.HasService("logger") {
		t.Error("Expected service 'logger' to exist")
	}
	
	// Get service
	service, err := locator.Get("logger")
	if err != nil {
		t.Errorf("Error getting service: %v", err)
	}
	
	// Check service type
	retrievedLogger, ok := service.(Logger)
	if !ok {
		t.Error("Retrieved service is not a Logger")
	}
	
	// Test the service
	retrievedLogger.Log("Test message")
	
	// Get non-existent service
	_, err = locator.Get("nonexistent")
	if err == nil {
		t.Error("Expected error when getting non-existent service")
	}
	
	// Remove service
	locator.Remove("logger")
	if locator.HasService("logger") {
		t.Error("Expected service 'logger' to be removed")
	}
}

// TestServiceLocatorGetTyped tests the typed service retrieval
func TestServiceLocatorGetTyped(t *testing.T) {
	locator := NewServiceLocator()
	
	// Register services
	logger := &ConsoleLogger{}
	repo := NewInMemoryUserRepository()
	
	locator.Register("logger", logger)
	locator.Register("userRepository", repo)
	
	// Get typed service
	var retrievedLogger Logger
	err := locator.GetTyped("logger", &retrievedLogger)
	if err != nil {
		t.Errorf("Error getting typed service: %v", err)
	}
	
	// Test service
	retrievedLogger.Log("Test typed service")
	
	// Test wrong type
	var wrongType NotificationService
	err = locator.GetTyped("logger", &wrongType)
	if err == nil {
		t.Error("Expected error when getting service with wrong type")
	}
}

// TestUserServiceWithServiceLocator tests the UserService with service locator
func TestUserServiceWithServiceLocator(t *testing.T) {
	// Create and configure service locator
	locator := NewServiceLocator()
	
	// Register required services
	logger := &ConsoleLogger{}
	repo := NewInMemoryUserRepository()
	notifier := &EmailNotificationService{Logger: logger}
	
	locator.Register("logger", logger)
	locator.Register("userRepository", repo)
	locator.Register("notificationService", notifier)
	
	// Create service with locator
	userService := NewUserService(locator)
	
	// Test creating a user
	err := userService.CreateUser("user1", "John Doe", "john@example.com")
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}
	
	// Test retrieving the user
	user, err := userService.GetUser("user1")
	if err != nil {
		t.Errorf("Error getting user: %v", err)
	}
	
	if user.Name != "John Doe" || user.Email != "john@example.com" {
		t.Errorf("User data mismatch, got: %+v", user)
	}
	
	// Test missing dependency
	badLocator := NewServiceLocator()
	badLocator.Register("logger", logger)
	// Intentionally not registering user repository
	
	badUserService := NewUserService(badLocator)
	err = badUserService.CreateUser("user2", "Jane Doe", "jane@example.com")
	if err == nil {
		t.Error("Expected error when user repository is missing")
	}
}

// MockLogger is a test implementation of Logger
type MockLogger struct {
	Messages []string
}

func (l *MockLogger) Log(message string) {
	l.Messages = append(l.Messages, message)
}

// TestServiceLocatorWithMocks tests the service locator with mock objects
func TestServiceLocatorWithMocks(t *testing.T) {
	// Create mock objects
	mockLogger := &MockLogger{Messages: []string{}}
	mockRepo := NewInMemoryUserRepository()
	
	// Create and configure service locator
	locator := NewServiceLocator()
	locator.Register("logger", mockLogger)
	locator.Register("userRepository", mockRepo)
	
	// Create service with locator
	userService := NewUserService(locator)
	
	// Test creating a user
	err := userService.CreateUser("user3", "Mock User", "mock@example.com")
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}
	
	// Verify mock was called
	if len(mockLogger.Messages) == 0 {
		t.Error("Expected logger to be called")
	}
	
	// Verify user was saved to repository
	user, err := mockRepo.FindByID("user3")
	if err != nil {
		t.Errorf("Error finding saved user: %v", err)
	}
	
	if user.Name != "Mock User" || user.Email != "mock@example.com" {
		t.Errorf("User data mismatch, got: %+v", user)
	}
}
