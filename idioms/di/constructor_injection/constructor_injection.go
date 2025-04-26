// Package constructor_injection demonstrates dependency injection via constructors.
package constructor_injection

import (
	"fmt"
	"time"
)

// Logger defines a simple logging interface
type Logger interface {
	Log(message string)
}

// ConsoleLogger is a concrete implementation that logs to console
type ConsoleLogger struct{}

func (l *ConsoleLogger) Log(message string) {
	fmt.Printf("[%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), message)
}

// FileLogger simulates logging to a file
type FileLogger struct {
	FilePath string
}

func (l *FileLogger) Log(message string) {
	// In a real implementation, this would write to a file
	fmt.Printf("[FILE:%s] %s\n", l.FilePath, message)
}

// UserService depends on a logger
type UserService struct {
	logger Logger
}

// NewUserService creates a new UserService with the provided logger
// This is constructor injection - the dependency is provided at creation time
func NewUserService(logger Logger) *UserService {
	return &UserService{
		logger: logger,
	}
}

// CreateUser simulates creating a user
func (s *UserService) CreateUser(name string) {
	s.logger.Log(fmt.Sprintf("Creating user: %s", name))
	// User creation logic would go here
	s.logger.Log(fmt.Sprintf("User created: %s", name))
}

// EmailService depends on multiple services
type EmailService struct {
	logger      Logger
	emailSender EmailSender
	config      *Config
}

// EmailSender defines an interface for sending emails
type EmailSender interface {
	Send(to, subject, body string) error
}

// Config holds service configuration
type Config struct {
	SenderEmail string
	RetryCount  int
}

// NewEmailService creates a new EmailService with all its dependencies
// This demonstrates constructor injection with multiple dependencies
func NewEmailService(logger Logger, emailSender EmailSender, config *Config) *EmailService {
	return &EmailService{
		logger:      logger,
		emailSender: emailSender,
		config:      config,
	}
}

// SendWelcomeEmail sends a welcome email to a user
func (s *EmailService) SendWelcomeEmail(userEmail, userName string) error {
	s.logger.Log(fmt.Sprintf("Sending welcome email to %s", userEmail))
	
	subject := "Welcome to our service!"
	body := fmt.Sprintf("Hello %s, welcome aboard!", userName)
	
	err := s.emailSender.Send(userEmail, subject, body)
	if err != nil {
		s.logger.Log(fmt.Sprintf("Failed to send welcome email: %v", err))
		return err
	}
	
	s.logger.Log("Welcome email sent successfully")
	return nil
}
