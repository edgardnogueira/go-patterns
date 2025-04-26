package external

import (
	"errors"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/infrastructure/logger"
)

// Common errors
var (
	ErrInvalidEmail     = errors.New("invalid email address")
	ErrEmptySubject     = errors.New("empty subject")
	ErrEmptyBody        = errors.New("empty body")
	ErrEmailSendFailure = errors.New("failed to send email")
)

// EmailProvider represents an email service provider
type EmailProvider struct {
	host     string
	port     int
	username string
	password string
	from     string
	logger   *logger.Logger
}

// EmailConfig represents the email provider configuration
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// NewEmailProvider creates a new email provider
func NewEmailProvider(config EmailConfig, log *logger.Logger) *EmailProvider {
	return &EmailProvider{
		host:     config.Host,
		port:     config.Port,
		username: config.Username,
		password: config.Password,
		from:     config.From,
		logger:   log,
	}
}

// SendEmail sends an email to the specified recipient
func (p *EmailProvider) SendEmail(to, subject, body string) error {
	if !isValidEmail(to) {
		return ErrInvalidEmail
	}
	
	if subject == "" {
		return ErrEmptySubject
	}
	
	if body == "" {
		return ErrEmptyBody
	}
	
	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", p.from, to, subject, body))
	
	auth := smtp.PlainAuth("", p.username, p.password, p.host)
	addr := fmt.Sprintf("%s:%d", p.host, p.port)
	
	err := smtp.SendMail(addr, auth, p.from, []string{to}, msg)
	if err != nil {
		p.logger.WithFields(map[string]interface{}{
			"to":      to,
			"subject": subject,
			"error":   err.Error(),
		}).Error("Failed to send email")
		return fmt.Errorf("%w: %v", ErrEmailSendFailure, err)
	}
	
	p.logger.WithFields(map[string]interface{}{
		"to":      to,
		"subject": subject,
	}).Info("Email sent successfully")
	
	return nil
}

// SendEmailToMultiple sends an email to multiple recipients
func (p *EmailProvider) SendEmailToMultiple(to []string, subject, body string) error {
	for _, recipient := range to {
		if !isValidEmail(recipient) {
			p.logger.WithField("email", recipient).Warn("Invalid email address")
			continue
		}
	}
	
	if subject == "" {
		return ErrEmptySubject
	}
	
	if body == "" {
		return ErrEmptyBody
	}
	
	toHeader := strings.Join(to, ", ")
	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", p.from, toHeader, subject, body))
	
	auth := smtp.PlainAuth("", p.username, p.password, p.host)
	addr := fmt.Sprintf("%s:%d", p.host, p.port)
	
	err := smtp.SendMail(addr, auth, p.from, to, msg)
	if err != nil {
		p.logger.WithFields(map[string]interface{}{
			"to":      to,
			"subject": subject,
			"error":   err.Error(),
		}).Error("Failed to send email to multiple recipients")
		return fmt.Errorf("%w: %v", ErrEmailSendFailure, err)
	}
	
	p.logger.WithFields(map[string]interface{}{
		"to":      to,
		"subject": subject,
	}).Info("Email sent successfully to multiple recipients")
	
	return nil
}

// SendTemplatedEmail sends an email using a template
func (p *EmailProvider) SendTemplatedEmail(to, subject, templateName string, data map[string]interface{}) error {
	// In a real application, this would use a templating engine like text/template
	// For simplicity, we'll just use a placeholder
	body := fmt.Sprintf("Template: %s\nData: %v", templateName, data)
	
	return p.SendEmail(to, subject, body)
}

// isValidEmail performs a basic validation of email format
func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}
