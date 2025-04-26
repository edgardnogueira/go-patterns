package main

import (
	"fmt"
	"strings"
)

// NotificationService depends directly on concrete implementations
// This violates the Dependency Inversion Principle because high-level modules
// should not depend on low-level modules but both should depend on abstractions
type NotificationService struct {
	// Direct dependency on concrete implementations
	emailSender    EmailSender
	smsSender      SMSSender
	pushNotifier   PushNotifier
}

// NewNotificationService creates a notification service
func NewNotificationService() *NotificationService {
	return &NotificationService{
		emailSender:    EmailSender{},
		smsSender:      SMSSender{},
		pushNotifier:   PushNotifier{},
	}
}

// Notify sends a notification based on the provided type
// This method has to be changed if we want to add a new notification method
func (ns *NotificationService) Notify(message, recipient, notificationType string) error {
	switch strings.ToLower(notificationType) {
	case "email":
		return ns.emailSender.SendEmail(message, recipient)
	case "sms":
		return ns.smsSender.SendSMS(message, recipient)
	case "push":
		return ns.pushNotifier.SendPushNotification(message, recipient)
	default:
		return fmt.Errorf("unknown notification type: %s", notificationType)
	}
}

// EmailSender is a low-level module that sends emails
type EmailSender struct{}

// SendEmail sends an email
func (es EmailSender) SendEmail(message, recipient string) error {
	fmt.Printf("Sending email to %s: %s\n", recipient, message)
	// In a real application, this would connect to an email service
	return nil
}

// SMSSender is a low-level module that sends SMS messages
type SMSSender struct{}

// SendSMS sends an SMS
func (ss SMSSender) SendSMS(message, recipient string) error {
	fmt.Printf("Sending SMS to %s: %s\n", recipient, message)
	// In a real application, this would connect to an SMS service
	return nil
}

// PushNotifier is a low-level module that sends push notifications
type PushNotifier struct{}

// SendPushNotification sends a push notification
func (pn PushNotifier) SendPushNotification(message, recipient string) error {
	fmt.Printf("Sending push notification to %s: %s\n", recipient, message)
	// In a real application, this would connect to a push notification service
	return nil
}

// User represents a user in the system
type User struct {
	Name  string
	Email string
	Phone string
	DeviceID string
}

// UserService depends directly on NotificationService
// This tightly couples the UserService to NotificationService
type UserService struct {
	notificationService *NotificationService
}

// NewUserService creates a new user service
func NewUserService() *UserService {
	return &UserService{
		notificationService: NewNotificationService(),
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
		return fmt.Errorf("failed to send welcome email: %w", err)
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
		return fmt.Errorf("failed to send password reset SMS: %w", err)
	}
	
	return nil
}

// This function demonstrates the notification service before applying DIP
func demonstrateNotificationBeforeDIP() {
	userService := NewUserService()
	
	// Register a new user
	user := User{
		Name:  "John Doe",
		Email: "john@example.com",
		Phone: "+1234567890",
		DeviceID: "device123",
	}
	
	err := userService.RegisterUser(user)
	if err != nil {
		fmt.Println("Error registering user:", err)
	}
	
	err = userService.ResetPassword(user)
	if err != nil {
		fmt.Println("Error resetting password:", err)
	}
	
	fmt.Println("\nThe problem with this approach is:")
	fmt.Println("1. NotificationService depends directly on concrete implementations (EmailSender, SMSSender, etc.)")
	fmt.Println("2. UserService depends directly on NotificationService")
	fmt.Println("3. Adding a new notification method requires modifying existing code")
	fmt.Println("4. Testing is difficult because we can't easily substitute implementations")
}
