package main

import (
	"errors"
	"testing"
)

// MockNotifier is a test implementation of NotificationSender
type MockNotifier struct {
	SentMessages  []string
	SentRecipients []string
	ShouldFail    bool
	ErrorToReturn error
}

// Send records the message and recipient for testing
func (m *MockNotifier) Send(message, recipient string) error {
	if m.ShouldFail {
		return m.ErrorToReturn
	}
	m.SentMessages = append(m.SentMessages, message)
	m.SentRecipients = append(m.SentRecipients, recipient)
	return nil
}

// TestNotificationService tests the NotificationService with different notifiers
func TestNotificationService(t *testing.T) {
	// Create a notification service
	notificationService := NewNotificationService()
	
	// Register a mock notifier
	mockEmailNotifier := &MockNotifier{}
	notificationService.RegisterNotifier("email", mockEmailNotifier)
	
	// Test sending a notification
	t.Run("sending email notification", func(t *testing.T) {
		err := notificationService.Notify("Test message", "test@example.com", "email")
		
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		// Verify the message was sent
		if len(mockEmailNotifier.SentMessages) != 1 {
			t.Errorf("Expected 1 message to be sent, got %d", len(mockEmailNotifier.SentMessages))
		}
		
		if mockEmailNotifier.SentMessages[0] != "Test message" {
			t.Errorf("Expected message 'Test message', got '%s'", mockEmailNotifier.SentMessages[0])
		}
		
		if mockEmailNotifier.SentRecipients[0] != "test@example.com" {
			t.Errorf("Expected recipient 'test@example.com', got '%s'", mockEmailNotifier.SentRecipients[0])
		}
	})
	
	// Test with an unregistered notifier
	t.Run("unregistered notifier", func(t *testing.T) {
		err := notificationService.Notify("Test message", "test", "sms")
		
		if err == nil {
			t.Error("Expected error for unregistered notifier, got nil")
		}
	})
	
	// Test with a failing notifier
	t.Run("failing notifier", func(t *testing.T) {
		mockFailingNotifier := &MockNotifier{
			ShouldFail:    true,
			ErrorToReturn: errors.New("send failed"),
		}
		
		notificationService.RegisterNotifier("failing", mockFailingNotifier)
		
		err := notificationService.Notify("Test message", "test", "failing")
		
		if err == nil {
			t.Error("Expected error from failing notifier, got nil")
		}
	})
	
	// Test adding multiple notifiers
	t.Run("multiple notifiers", func(t *testing.T) {
		mockSMSNotifier := &MockNotifier{}
		mockPushNotifier := &MockNotifier{}
		
		notificationService.RegisterNotifier("sms", mockSMSNotifier)
		notificationService.RegisterNotifier("push", mockPushNotifier)
		
		// Send SMS
		err := notificationService.Notify("SMS message", "+1234567890", "sms")
		if err != nil {
			t.Errorf("Expected no error for SMS, got %v", err)
		}
		
		// Send push
		err = notificationService.Notify("Push message", "device123", "push")
		if err != nil {
			t.Errorf("Expected no error for push, got %v", err)
		}
		
		// Verify SMS was sent
		if len(mockSMSNotifier.SentMessages) != 1 {
			t.Errorf("Expected 1 SMS message, got %d", len(mockSMSNotifier.SentMessages))
		}
		
		// Verify push was sent
		if len(mockPushNotifier.SentMessages) != 1 {
			t.Errorf("Expected 1 push message, got %d", len(mockPushNotifier.SentMessages))
		}
	})
}

// TestUserService tests the UserService with injected dependencies
func TestUserService(t *testing.T) {
	// Create mock notifiers
	mockEmailNotifier := &MockNotifier{}
	mockSMSNotifier := &MockNotifier{}
	
	// Create notification service and register notifiers
	notificationService := NewNotificationService()
	notificationService.RegisterNotifier("email", mockEmailNotifier)
	notificationService.RegisterNotifier("sms", mockSMSNotifier)
	
	// Create user service with the notification service
	userService := NewUserService(notificationService)
	
	// Create a test user
	user := User{
		Name:     "Test User",
		Email:    "test@example.com",
		Phone:    "+1234567890",
		DeviceID: "device123",
	}
	
	// Test registering a user
	t.Run("register user", func(t *testing.T) {
		err := userService.RegisterUser(user)
		
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		// Verify an email was sent
		if len(mockEmailNotifier.SentMessages) != 1 {
			t.Errorf("Expected 1 email, got %d", len(mockEmailNotifier.SentMessages))
		}
		
		// Verify the email was sent to the right address
		if len(mockEmailNotifier.SentRecipients) != 1 || mockEmailNotifier.SentRecipients[0] != user.Email {
			t.Errorf("Expected email to %s, got %v", user.Email, mockEmailNotifier.SentRecipients)
		}
	})
	
	// Reset the mock
	mockEmailNotifier.SentMessages = nil
	mockEmailNotifier.SentRecipients = nil
	
	// Test resetting a password
	t.Run("reset password", func(t *testing.T) {
		err := userService.ResetPassword(user)
		
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		// Verify an SMS was sent
		if len(mockSMSNotifier.SentMessages) != 1 {
			t.Errorf("Expected 1 SMS, got %d", len(mockSMSNotifier.SentMessages))
		}
		
		// Verify the SMS was sent to the right number
		if len(mockSMSNotifier.SentRecipients) != 1 || mockSMSNotifier.SentRecipients[0] != user.Phone {
			t.Errorf("Expected SMS to %s, got %v", user.Phone, mockSMSNotifier.SentRecipients)
		}
	})
}

// TestDependencyInversion tests the dependency inversion principles in action
func TestDependencyInversion(t *testing.T) {
	t.Run("dependency inversion with different implementations", func(t *testing.T) {
		// Create a notification service
		notificationService := NewNotificationService()
		
		// Register different implementations of the NotificationSender interface
		notificationService.RegisterNotifier("email", EmailNotifier{})
		notificationService.RegisterNotifier("sms", SMSNotifier{})
		notificationService.RegisterNotifier("push", PushNotifier{})
		notificationService.RegisterNotifier("slack", SlackNotifier{})
		
		// Create a custom implementation at runtime
		customNotifier := &MockNotifier{}
		notificationService.RegisterNotifier("custom", customNotifier)
		
		// Create user service that depends on the abstraction
		userService := NewUserService(notificationService)
		
		// The key point: UserService depends on the NotificationService interface,
		// not on concrete implementations of notifiers
		
		// This demonstrates DIP:
		// 1. High-level module (UserService) depends on abstraction (NotificationService)
		// 2. Low-level modules (EmailNotifier, etc.) depend on abstraction (NotificationSender)
		// 3. Both high and low level modules depend on abstractions
		
		// The benefit: We can easily swap implementations without changing code
		
		// This test doesn't verify specific behaviors but demonstrates the principle
	})
}

// Test with a completely different NotificationSender implementation
type LogNotifier struct {
	LoggedMessages []string
}

func (l *LogNotifier) Send(message, recipient string) error {
	l.LoggedMessages = append(l.LoggedMessages, message)
	return nil
}

func TestNewNotifierType(t *testing.T) {
	t.Run("adding new notifier type without modifying code", func(t *testing.T) {
		// Create services
		notificationService := NewNotificationService()
		logNotifier := &LogNotifier{}
		
		// Register our new notifier type
		notificationService.RegisterNotifier("log", logNotifier)
		
		// Create user service
		userService := NewUserService(notificationService)
		
		// Create a test user
		user := User{
			Name:  "Test User",
			Email: "log@example.com",
		}
		
		// Send notification via a method that normally uses email
		err := userService.notificationService.Notify("Log message", "log-recipient", "log")
		
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		// Verify the message was logged
		if len(logNotifier.LoggedMessages) != 1 {
			t.Errorf("Expected 1 logged message, got %d", len(logNotifier.LoggedMessages))
		}
		
		// The key insight: We added a completely new notification method
		// without modifying any existing code in UserService or NotificationService
	})
}
