package main

import (
	"fmt"
)

// NotificationSender is an interface that defines the contract for sending notifications
// This is the abstraction that both high-level and low-level modules depend on
type NotificationSender interface {
	Send(message, recipient string) error
}

// EmailNotifier implements the NotificationSender interface for email
type EmailNotifier struct{}

// Send sends an email notification
func (en EmailNotifier) Send(message, recipient string) error {
	fmt.Printf("Sending email to %s: %s\n", recipient, message)
	// In a real application, this would connect to an email service
	return nil
}

// SMSNotifier implements the NotificationSender interface for SMS
type SMSNotifier struct{}

// Send sends an SMS notification
func (sn SMSNotifier) Send(message, recipient string) error {
	fmt.Printf("Sending SMS to %s: %s\n", recipient, message)
	// In a real application, this would connect to an SMS service
	return nil
}

// PushNotifier implements the NotificationSender interface for push notifications
type PushNotifier struct{}

// Send sends a push notification
func (pn PushNotifier) Send(message, recipient string) error {
	fmt.Printf("Sending push notification to %s: %s\n", recipient, message)
	// In a real application, this would connect to a push notification service
	return nil
}

// SlackNotifier is a new implementation of NotificationSender for Slack messages
// It can be added without changing any existing code
type SlackNotifier struct{}

// Send sends a Slack notification
func (sn SlackNotifier) Send(message, recipient string) error {
	fmt.Printf("Sending Slack message to %s: %s\n", recipient, message)
	// In a real application, this would connect to Slack's API
	return nil
}

// NotificationService depends on the NotificationSender interface, not concrete implementations
// This follows DIP because the high-level module depends on an abstraction
type NotificationService struct {
	notifiers map[string]NotificationSender
}

// NewNotificationService creates a notification service with the provided notifiers
func NewNotificationService() *NotificationService {
	return &NotificationService{
		notifiers: make(map[string]NotificationSender),
	}
}

// RegisterNotifier registers a notification sender for a specific type
func (ns *NotificationService) RegisterNotifier(notificationType string, notifier NotificationSender) {
	ns.notifiers[notificationType] = notifier
}

// Notify sends a notification using the appropriate notifier
func (ns *NotificationService) Notify(message, recipient, notificationType string) error {
	notifier, exists := ns.notifiers[notificationType]
	if !exists {
		return fmt.Errorf("no notifier registered for type: %s", notificationType)
	}
	return notifier.Send(message, recipient)
}

// User represents a user in the system
type User struct {
	Name     string
	Email    string
	Phone    string
	DeviceID string
	SlackID  string
}

// UserService depends on the NotificationService interface
// The dependency is injected rather than created internally
type UserService struct {
	notificationService *NotificationService
}

// NewUserService creates a new user service with the provided notification service
func NewUserService(notificationService *NotificationService) *UserService {
	return &UserService{
		notificationService: notificationService,
	}
}

// RegisterUser registers a new user and sends a welcome notification
func (us *UserService) RegisterUser(user User) error {
	// In a real application, this would save the user to a database
	fmt.Printf("Registering user: %s\n", user.Name)
	
	// Send welcome notification via email
	message := fmt.Sprintf("Welcome, %s! Thank you for registering.", user.Name)
	err := us.notificationService.Notify(message, user.Email, "email")
	if err != nil {
		return fmt.Errorf("failed to send welcome notification: %w", err)
	}
	
	return nil
}

// ResetPassword resets a user's password and sends a notification
func (us *UserService) ResetPassword(user User) error {
	// In a real application, this would reset the user's password
	fmt.Printf("Resetting password for user: %s\n", user.Name)
	
	// Send password reset notification via SMS
	message := "Your password has been reset. If you didn't request this, please contact support."
	err := us.notificationService.Notify(message, user.Phone, "sms")
	if err != nil {
		return fmt.Errorf("failed to send password reset notification: %w", err)
	}
	
	return nil
}

// SendActivityAlert sends an activity alert to the user via push notification
func (us *UserService) SendActivityAlert(user User, activity string) error {
	message := fmt.Sprintf("Activity alert: %s", activity)
	err := us.notificationService.Notify(message, user.DeviceID, "push")
	if err != nil {
		return fmt.Errorf("failed to send activity alert: %w", err)
	}
	return nil
}

// SendTeamNotification sends a team notification via Slack
func (us *UserService) SendTeamNotification(user User, notification string) error {
	message := fmt.Sprintf("Team notification: %s", notification)
	err := us.notificationService.Notify(message, user.SlackID, "slack")
	if err != nil {
		return fmt.Errorf("failed to send team notification: %w", err)
	}
	return nil
}

// MockNotifier is a mock implementation of NotificationSender for testing
type MockNotifier struct {
	messages []string
	recipients []string
}

// Send records the message and recipient for testing verification
func (mn *MockNotifier) Send(message, recipient string) error {
	mn.messages = append(mn.messages, message)
	mn.recipients = append(mn.recipients, recipient)
	return nil
}

// This function demonstrates the notification service after applying DIP
func demonstrateNotificationAfterDIP() {
	// Create notification service and register notifiers
	notificationService := NewNotificationService()
	notificationService.RegisterNotifier("email", EmailNotifier{})
	notificationService.RegisterNotifier("sms", SMSNotifier{})
	notificationService.RegisterNotifier("push", PushNotifier{})
	notificationService.RegisterNotifier("slack", SlackNotifier{})
	
	// Create user service with the notification service injected
	userService := NewUserService(notificationService)
	
	// Create a user
	user := User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Phone:    "+1234567890",
		DeviceID: "device123",
		SlackID:  "U123456",
	}
	
	// Register user
	err := userService.RegisterUser(user)
	if err != nil {
		fmt.Println("Error registering user:", err)
	}
	
	// Reset password
	err = userService.ResetPassword(user)
	if err != nil {
		fmt.Println("Error resetting password:", err)
	}
	
	// Send activity alert
	err = userService.SendActivityAlert(user, "New login from unknown device")
	if err != nil {
		fmt.Println("Error sending activity alert:", err)
	}
	
	// Send team notification
	err = userService.SendTeamNotification(user, "Team meeting at 3 PM")
	if err != nil {
		fmt.Println("Error sending team notification:", err)
	}
	
	// Demonstration of testing with a mock
	fmt.Println("\nDemonstrating testing with a mock notifier:")
	mockNotifier := &MockNotifier{}
	testNotificationService := NewNotificationService()
	testNotificationService.RegisterNotifier("email", mockNotifier)
	
	testUserService := NewUserService(testNotificationService)
	_ = testUserService.RegisterUser(user)
	
	fmt.Printf("Mock recorded message: %s\n", mockNotifier.messages[0])
	fmt.Printf("Mock recorded recipient: %s\n", mockNotifier.recipients[0])
	
	fmt.Println("\nBenefits of this approach:")
	fmt.Println("1. High-level modules (UserService) depend on abstractions (NotificationService)")
	fmt.Println("2. Low-level modules (EmailNotifier, etc.) depend on abstractions (NotificationSender)")
	fmt.Println("3. We can add new notification types without modifying existing code")
	fmt.Println("4. Dependencies are injected, making testing much easier")
	fmt.Println("5. Modules are loosely coupled and can be developed independently")
}
