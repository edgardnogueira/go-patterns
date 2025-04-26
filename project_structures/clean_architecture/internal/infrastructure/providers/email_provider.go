package providers

import (
	"context"
	"fmt"
	"time"
)

// EmailProvider defines the interface for sending emails
type EmailProvider interface {
	SendEmail(ctx context.Context, to, subject, body string) error
}

// MockEmailProvider is a mock implementation of EmailProvider for demonstration
type MockEmailProvider struct {
	// In a real application, this might include client configuration
}

// NewMockEmailProvider creates a new instance of MockEmailProvider
func NewMockEmailProvider() *MockEmailProvider {
	return &MockEmailProvider{}
}

// SendEmail simulates sending an email by logging it
func (p *MockEmailProvider) SendEmail(ctx context.Context, to, subject, body string) error {
	// Simulate network latency
	time.Sleep(100 * time.Millisecond)

	// In a real app, this would connect to an email service
	// For demonstration, we just log the email
	fmt.Printf("[EMAIL] To: %s, Subject: %s, Body: %s\n", to, subject, body)
	return nil
}

// Email templates
const (
	WelcomeEmailTemplate = `
Hello %s,

Welcome to our platform! We're excited to have you join us.

Here are a few things you can do to get started:
1. Complete your profile
2. Explore our features
3. Connect with other users

If you have any questions, feel free to reply to this email.

Best regards,
The Team
`

	ProfileUpdateEmailTemplate = `
Hello %s,

Your profile has been successfully updated.

If you did not make these changes, please contact our support team immediately.

Best regards,
The Team
`
)
