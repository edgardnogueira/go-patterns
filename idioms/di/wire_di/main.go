package wire_di

import (
	"fmt"
	"os"
)

// The main function demonstrates how to use the Wire-generated
// InitializeApplication function to create an application with all
// dependencies properly wired together.
func Example_wireDemo() {
	// Initialize the application using the Wire-generated function
	app, err := InitializeApplication()
	if err != nil {
		fmt.Printf("Failed to initialize application: %v\n", err)
		os.Exit(1)
	}
	
	// Ensure resources are properly released
	defer func() {
		if err := app.Close(); err != nil {
			fmt.Printf("Error closing application: %v\n", err)
		}
	}()
	
	// Run the application
	if err := app.Run(); err != nil {
		fmt.Printf("Application error: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("Application completed successfully")
}

// This demonstrates how the application would be initialized
// without using Wire, by manually wiring all dependencies.
// This is shown for comparison purposes.
func Example_manualWiring() {
	// Create basic dependencies
	config := NewConfig()
	logger := NewLogger(config)
	
	// Create database connection
	db, err := NewDatabaseConnection(config)
	if err != nil {
		fmt.Printf("Failed to create database connection: %v\n", err)
		os.Exit(1)
	}
	
	// Create API client
	apiClient := NewAPIClient(config, logger)
	
	// Create repositories
	userRepo := NewUserRepository(db, logger)
	messageRepo := NewMessageRepository(db, logger)
	
	// Create notification service
	notificationService := NewNotificationService(apiClient, logger, config)
	
	// Create services
	userService := NewUserService(userRepo, logger, notificationService)
	messageService := NewMessageService(messageRepo, userRepo, logger, notificationService)
	
	// Create application
	app := NewApplication(userService, messageService, config, logger, db)
	
	// Ensure resources are properly released
	defer func() {
		if err := app.Close(); err != nil {
			fmt.Printf("Error closing application: %v\n", err)
		}
	}()
	
	// Run the application
	if err := app.Run(); err != nil {
		fmt.Printf("Application error: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("Application completed successfully")
}

// This function shows how to override dependencies for testing
func Example_testOverrides() {
	// Create configuration with test values
	testConfig := NewConfig()
	testConfig.Environment = "testing"
	testConfig.DatabaseURL = "inmemory:test"
	
	logger := NewLogger(testConfig)
	
	// Create database connection
	db, err := NewDatabaseConnection(testConfig)
	if err != nil {
		fmt.Printf("Failed to create database connection: %v\n", err)
		os.Exit(1)
	}
	
	// Create API client
	apiClient := NewAPIClient(testConfig, logger)
	
	// Create repositories
	userRepo := NewUserRepository(db, logger)
	messageRepo := NewMessageRepository(db, logger)
	
	// Create a mock notification service that doesn't send real notifications
	mockNotificationService := &MockNotificationService{Logger: logger}
	
	// Create services with the mock
	userService := NewUserService(userRepo, logger, mockNotificationService)
	messageService := NewMessageService(messageRepo, userRepo, logger, mockNotificationService)
	
	// Create application
	app := NewApplication(userService, messageService, testConfig, logger, db)
	
	// Ensure resources are properly released
	defer func() {
		if err := app.Close(); err != nil {
			fmt.Printf("Error closing application: %v\n", err)
		}
	}()
	
	// Run the application
	if err := app.Run(); err != nil {
		fmt.Printf("Application error: %v\n", err)
		os.Exit(1)
	}
	
	// Check if notifications were "sent"
	fmt.Printf("Notifications sent: %d\n", mockNotificationService.NotificationCount)
}

// MockNotificationService is a mock implementation for testing
type MockNotificationService struct {
	Logger            *Logger
	NotificationCount int
}

func (s *MockNotificationService) SendNotification(userID, message string) error {
	s.Logger.Log(fmt.Sprintf("[MOCK] Would send to %s: %s", userID, message))
	s.NotificationCount++
	return nil
}
