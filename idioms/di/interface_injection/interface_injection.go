// Package interface_injection demonstrates dependency injection using interfaces.
package interface_injection

import (
	"fmt"
	"time"
)

// Storage defines an interface for data persistence
type Storage interface {
	Save(key string, data []byte) error
	Load(key string) ([]byte, error)
	Delete(key string) error
}

// MemoryStorage implements the Storage interface with in-memory storage
type MemoryStorage struct {
	data map[string][]byte
}

// NewMemoryStorage creates a new in-memory storage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string][]byte),
	}
}

func (s *MemoryStorage) Save(key string, data []byte) error {
	s.data[key] = data
	return nil
}

func (s *MemoryStorage) Load(key string) ([]byte, error) {
	data, ok := s.data[key]
	if !ok {
		return nil, fmt.Errorf("key not found: %s", key)
	}
	return data, nil
}

func (s *MemoryStorage) Delete(key string) error {
	delete(s.data, key)
	return nil
}

// FileStorage implements the Storage interface with file-based storage
type FileStorage struct {
	BasePath string
}

func (s *FileStorage) Save(key string, data []byte) error {
	// Simplified for example purposes - would actually write to file
	fmt.Printf("Writing %d bytes to file %s/%s\n", len(data), s.BasePath, key)
	return nil
}

func (s *FileStorage) Load(key string) ([]byte, error) {
	// Simplified for example purposes - would actually read from file
	fmt.Printf("Reading from file %s/%s\n", s.BasePath, key)
	return []byte("file content"), nil
}

func (s *FileStorage) Delete(key string) error {
	// Simplified for example purposes - would actually delete file
	fmt.Printf("Deleting file %s/%s\n", s.BasePath, key)
	return nil
}

// Notifier defines an interface for sending notifications
type Notifier interface {
	Notify(message string) error
}

// EmailNotifier implements the Notifier interface via email
type EmailNotifier struct {
	SMTPServer string
	FromEmail  string
}

func (n *EmailNotifier) Notify(message string) error {
	// Simplified for example purposes
	fmt.Printf("Sending email via %s: %s\n", n.SMTPServer, message)
	return nil
}

// SMSNotifier implements the Notifier interface via SMS
type SMSNotifier struct {
	AccountID string
	APIKey    string
}

func (n *SMSNotifier) Notify(message string) error {
	// Simplified for example purposes
	fmt.Printf("Sending SMS via account %s: %s\n", n.AccountID, message)
	return nil
}

// Logger defines an interface for logging
type Logger interface {
	Log(level string, message string)
}

// ConsoleLogger implements the Logger interface for console output
type ConsoleLogger struct{}

func (l *ConsoleLogger) Log(level string, message string) {
	fmt.Printf("[%s] %s: %s\n", time.Now().Format("2006-01-02 15:04:05"), level, message)
}

// UserService depends on multiple interfaces
type UserService struct {
	storage    Storage
	notifier   Notifier
	logger     Logger
	appVersion string
}

// NewUserService creates a new UserService with injected dependencies
func NewUserService(storage Storage, notifier Notifier, logger Logger, appVersion string) *UserService {
	return &UserService{
		storage:    storage,
		notifier:   notifier,
		logger:     logger,
		appVersion: appVersion,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(id string, userData []byte) error {
	s.logger.Log("INFO", fmt.Sprintf("Creating user with ID: %s", id))
	
	// Store user data
	err := s.storage.Save(fmt.Sprintf("user:%s", id), userData)
	if err != nil {
		s.logger.Log("ERROR", fmt.Sprintf("Failed to create user: %v", err))
		return fmt.Errorf("storage error: %w", err)
	}
	
	// Send notification
	err = s.notifier.Notify(fmt.Sprintf("New user created: %s", id))
	if err != nil {
		s.logger.Log("WARN", fmt.Sprintf("Failed to send notification: %v", err))
		// Continue despite notification failure
	}
	
	s.logger.Log("INFO", fmt.Sprintf("User created successfully: %s", id))
	return nil
}

// GetUser retrieves a user
func (s *UserService) GetUser(id string) ([]byte, error) {
	s.logger.Log("INFO", fmt.Sprintf("Retrieving user with ID: %s", id))
	
	// Retrieve user data
	userData, err := s.storage.Load(fmt.Sprintf("user:%s", id))
	if err != nil {
		s.logger.Log("ERROR", fmt.Sprintf("Failed to retrieve user: %v", err))
		return nil, fmt.Errorf("user not found: %w", err)
	}
	
	return userData, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(id string) error {
	s.logger.Log("INFO", fmt.Sprintf("Deleting user with ID: %s", id))
	
	// Delete user data
	err := s.storage.Delete(fmt.Sprintf("user:%s", id))
	if err != nil {
		s.logger.Log("ERROR", fmt.Sprintf("Failed to delete user: %v", err))
		return fmt.Errorf("delete error: %w", err)
	}
	
	// Send notification
	err = s.notifier.Notify(fmt.Sprintf("User deleted: %s", id))
	if err != nil {
		s.logger.Log("WARN", fmt.Sprintf("Failed to send notification: %v", err))
		// Continue despite notification failure
	}
	
	s.logger.Log("INFO", fmt.Sprintf("User deleted successfully: %s", id))
	return nil
}
