package constructor_injection

import (
	"fmt"
	"strings"
	"testing"
)

// MockLogger implements Logger interface for testing
type MockLogger struct {
	Logs []string
}

func (m *MockLogger) Log(message string) {
	m.Logs = append(m.Logs, message)
}

// MockEmailSender implements EmailSender for testing
type MockEmailSender struct {
	SentEmails []struct {
		To      string
		Subject string
		Body    string
	}
	ShouldFail bool
}

func (m *MockEmailSender) Send(to, subject, body string) error {
	if m.ShouldFail {
		return fmt.Errorf("failed to send email")
	}
	
	m.SentEmails = append(m.SentEmails, struct {
		To      string
		Subject string
		Body    string
	}{To: to, Subject: subject, Body: body})
	
	return nil
}

func TestUserServiceCreation(t *testing.T) {
	// Create a mock logger
	mockLogger := &MockLogger{}
	
	// Inject the mock logger into UserService
	userService := NewUserService(mockLogger)
	
	// Test the service
	userService.CreateUser("Alice")
	
	// Verify logs
	if len(mockLogger.Logs) != 2 {
		t.Errorf("Expected 2 log entries, got %d", len(mockLogger.Logs))
	}
	
	if !strings.Contains(mockLogger.Logs[0], "Creating user: Alice") {
		t.Errorf("Expected log to contain 'Creating user: Alice', got %s", mockLogger.Logs[0])
	}
	
	if !strings.Contains(mockLogger.Logs[1], "User created: Alice") {
		t.Errorf("Expected log to contain 'User created: Alice', got %s", mockLogger.Logs[1])
	}
}

func TestEmailServiceSuccess(t *testing.T) {
	// Create mock dependencies
	mockLogger := &MockLogger{}
	mockEmailSender := &MockEmailSender{}
	config := &Config{
		SenderEmail: "noreply@example.com",
		RetryCount:  3,
	}
	
	// Inject dependencies into EmailService
	emailService := NewEmailService(mockLogger, mockEmailSender, config)
	
	// Test the service
	err := emailService.SendWelcomeEmail("user@example.com", "Bob")
	
	// Verify behavior
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if len(mockEmailSender.SentEmails) != 1 {
		t.Errorf("Expected 1 email to be sent, got %d", len(mockEmailSender.SentEmails))
	}
	
	sentEmail := mockEmailSender.SentEmails[0]
	if sentEmail.To != "user@example.com" {
		t.Errorf("Expected email to be sent to user@example.com, got %s", sentEmail.To)
	}
	
	if !strings.Contains(sentEmail.Body, "Hello Bob") {
		t.Errorf("Expected email body to contain 'Hello Bob', got %s", sentEmail.Body)
	}
	
	if len(mockLogger.Logs) != 2 {
		t.Errorf("Expected 2 log entries, got %d", len(mockLogger.Logs))
	}
}

func TestEmailServiceFailure(t *testing.T) {
	// Create mock dependencies
	mockLogger := &MockLogger{}
	mockEmailSender := &MockEmailSender{ShouldFail: true}
	config := &Config{
		SenderEmail: "noreply@example.com",
		RetryCount:  3,
	}
	
	// Inject dependencies into EmailService
	emailService := NewEmailService(mockLogger, mockEmailSender, config)
	
	// Test the service with failing email sender
	err := emailService.SendWelcomeEmail("user@example.com", "Charlie")
	
	// Verify behavior
	if err == nil {
		t.Error("Expected an error, got nil")
	}
	
	// Verify logs
	if len(mockLogger.Logs) != 2 {
		t.Errorf("Expected 2 log entries, got %d", len(mockLogger.Logs))
	}
	
	if !strings.Contains(mockLogger.Logs[1], "Failed to send") {
		t.Errorf("Expected failure log, got %s", mockLogger.Logs[1])
	}
}

func TestUsingDifferentLoggers(t *testing.T) {
	// Create two different logger implementations
	consoleLogger := &ConsoleLogger{}
	fileLogger := &FileLogger{FilePath: "users.log"}
	
	// Create two user services with different loggers
	consoleUserService := NewUserService(consoleLogger)
	fileUserService := NewUserService(fileLogger)
	
	// This just demonstrates we can swap implementations
	// Real test would capture output and verify it
	consoleUserService.CreateUser("User1")
	fileUserService.CreateUser("User2")
	
	// Test passes if code executes without panic
	// In a real test, we'd verify the output format differences
}
