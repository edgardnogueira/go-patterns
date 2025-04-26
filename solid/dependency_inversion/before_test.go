package main

import (
	"testing"
)

// Test the NotificationService without DIP
func TestNotificationServiceBeforeDIP(t *testing.T) {
	// Create a notification service - notice there's no way to inject dependencies
	notificationService := NewNotificationService()
	
	// Test sending an email notification
	t.Run("sending email notification", func(t *testing.T) {
		// We can only test with the concrete implementations that are hard-coded
		// in the NotificationService
		err := notificationService.Notify("Test email", "test@example.com", "email")
		
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		// We have no way to verify if the message was actually sent or how it was processed
		// because the EmailSender is tightly coupled inside the NotificationService
		
		t.Log("Notice: We can't verify that the email was actually sent")
		t.Log("The email sender is hard-coded in the NotificationService")
	})
	
	// Test sending an SMS notification
	t.Run("sending SMS notification", func(t *testing.T) {
		err := notificationService.Notify("Test SMS", "+1234567890", "sms")
		
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		// Again, we can't verify the actual behavior
		
		t.Log("Notice: We can't verify that the SMS was actually sent")
		t.Log("The SMS sender is hard-coded in the NotificationService")
	})
	
	// Test invalid notification type
	t.Run("invalid notification type", func(t *testing.T) {
		err := notificationService.Notify("Test message", "recipient", "invalid")
		
		if err == nil {
			t.Error("Expected error for invalid notification type, got nil")
		}
	})
}

// Test the UserService without DIP
func TestUserServiceBeforeDIP(t *testing.T) {
	// Create a user service - notice the NotificationService is created internally
	userService := NewUserService()
	
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
		
		// We have no way to verify that the notification was sent,
		// because the notification service is tightly coupled
		
		t.Log("Notice: We can't verify that the welcome email was sent")
		t.Log("The NotificationService is created inside UserService")
	})
	
	// Test resetting a password
	t.Run("reset password", func(t *testing.T) {
		err := userService.ResetPassword(user)
		
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		// Again, we can't verify the actual behavior
		
		t.Log("Notice: We can't verify that the password reset SMS was sent")
		t.Log("The NotificationService is created inside UserService")
	})
}

// Try to test with mock objects - shows why it's difficult without DIP
func TestMockingDifficulties(t *testing.T) {
	t.Run("cannot easily mock dependencies", func(t *testing.T) {
		// We cannot easily mock the NotificationService because:
		// 1. It creates its own EmailSender, SMSSender, and PushNotifier
		// 2. UserService creates its own NotificationService
		
		// To test the UserService with a mock notification service, we would need to:
		// 1. Create a test-specific version of UserService that accepts a NotificationService
		// 2. Or use global variables or function pointers (which makes code harder to understand)
		// 3. Or use complex mocking frameworks or runtime reflection
		
		// Instead of actually testing, this test just documents the limitation
		
		t.Log("Without DIP, we can't easily inject mock dependencies")
		t.Log("This makes testing much more difficult, as we can't verify behavior")
		t.Log("or isolate the unit under test from its dependencies")
		
		// The test does nothing but document the issue
		t.Skip("This test just documents mocking limitations without DIP")
	})
}

// TestAddingNewNotificationTypeWithoutDIP demonstrates why adding new notification types is difficult
func TestAddingNewNotificationTypeWithoutDIP(t *testing.T) {
	t.Run("adding new notification type requires modifying code", func(t *testing.T) {
		// To add a new notification type (e.g., Slack), we would need to:
		// 1. Create a new SlackSender struct with a SendSlack method
		// 2. Add a new field for SlackSender in NotificationService
		// 3. Add a new case in the switch statement in Notify method
		// 4. Initialize the SlackSender in NewNotificationService
		
		// This violates the Open/Closed Principle (part of SOLID)
		
		t.Log("Without DIP, adding a new notification type requires:")
		t.Log("1. Adding a new sender struct")
		t.Log("2. Adding a new field to NotificationService")
		t.Log("3. Modifying the Notify method's switch statement")
		t.Log("4. Modifying the NewNotificationService constructor")
		
		// This test does nothing but document the issue
		t.Skip("This test just documents extensibility limitations without DIP")
	})
}

func TestDIPViolationIssues(t *testing.T) {
	t.Run("highlight issues without DIP", func(t *testing.T) {
		// This "test" highlights the issues with DIP violations
		
		// ISSUE 1: Difficult to test
		t.Log("ISSUE 1: Difficult to test")
		t.Log("- NotificationService creates its own EmailSender, SMSSender, etc.")
		t.Log("- UserService creates its own NotificationService")
		t.Log("- No way to inject mock implementations for testing")
		
		// ISSUE 2: Tight coupling
		t.Log("ISSUE 2: Tight coupling")
		t.Log("- High-level modules depend directly on low-level modules")
		t.Log("- Changes to low-level modules may affect high-level modules")
		t.Log("- Changes to low-level modules may require recompiling high-level modules")
		
		// ISSUE 3: Poor extensibility
		t.Log("ISSUE 3: Poor extensibility")
		t.Log("- Adding new notification types requires modifying NotificationService")
		t.Log("- No way to add new implementations without changing existing code")
		
		// ISSUE 4: No reusability
		t.Log("ISSUE 4: No reusability")
		t.Log("- NotificationService can't be used with different sender implementations")
		t.Log("- UserService can't be used with different notification services")
		
		// This is not a real test, just documentation
		t.Skip("This is not a real test, just documentation of DIP violation issues")
	})
}
