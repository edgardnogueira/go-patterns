package wire_di

import (
	"fmt"
	"time"
)

// NotificationService provides notification functionality
type NotificationService interface {
	SendNotification(userID, message string) error
}

// EmailNotificationService is an email implementation of NotificationService
type EmailNotificationService struct {
	apiClient *APIClient
	logger    *Logger
	config    *Config
}

// NewNotificationService creates a new notification service
func NewNotificationService(client *APIClient, logger *Logger, config *Config) NotificationService {
	logger.Log("Creating email notification service")
	return &EmailNotificationService{
		apiClient: client,
		logger:    logger,
		config:    config,
	}
}

// SendNotification sends a notification to a user
func (s *EmailNotificationService) SendNotification(userID, message string) error {
	s.logger.Log(fmt.Sprintf("Sending notification to user %s: %s", userID, message))
	
	// In a real implementation, this would use an email service
	// For this example, we'll use the API client
	payload := map[string]string{
		"user_id": userID,
		"message": message,
	}
	
	return s.apiClient.Call("/notifications/email", payload)
}

// UserService provides user-related functionality
type UserService struct {
	userRepository UserRepository
	logger         *Logger
	notification   NotificationService
}

// NewUserService creates a new user service
func NewUserService(
	userRepo UserRepository,
	logger *Logger,
	notification NotificationService,
) *UserService {
	logger.Log("Creating user service")
	return &UserService{
		userRepository: userRepo,
		logger:         logger,
		notification:   notification,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(username, email string) (*User, error) {
	s.logger.Log(fmt.Sprintf("Creating user: %s (%s)", username, email))
	
	// Generate a unique ID (in a real app, this would be more robust)
	id := fmt.Sprintf("user_%d", time.Now().UnixNano())
	
	// Create and save the user
	user := &User{
		ID:       id,
		Username: username,
		Email:    email,
	}
	
	if err := s.userRepository.Save(user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}
	
	// Send welcome notification
	if err := s.notification.SendNotification(id, "Welcome to our service!"); err != nil {
		s.logger.Log(fmt.Sprintf("Failed to send welcome notification: %v", err))
		// Continue despite notification failure
	}
	
	return user, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(id string) (*User, error) {
	s.logger.Log(fmt.Sprintf("Getting user with ID: %s", id))
	return s.userRepository.FindByID(id)
}

// MessageService provides message-related functionality
type MessageService struct {
	messageRepository MessageRepository
	userRepository    UserRepository
	logger            *Logger
	notification      NotificationService
}

// NewMessageService creates a new message service
func NewMessageService(
	messageRepo MessageRepository,
	userRepo UserRepository,
	logger *Logger,
	notification NotificationService,
) *MessageService {
	logger.Log("Creating message service")
	return &MessageService{
		messageRepository: messageRepo,
		userRepository:    userRepo,
		logger:            logger,
		notification:      notification,
	}
}

// CreateMessage creates a new message
func (s *MessageService) CreateMessage(content, userID string) (*Message, error) {
	s.logger.Log(fmt.Sprintf("Creating message for user %s: %s", userID, content))
	
	// Verify user exists
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user: %w", err)
	}
	
	// Generate a unique ID
	id := fmt.Sprintf("msg_%d", time.Now().UnixNano())
	
	// Create and save the message
	message := &Message{
		ID:      id,
		Content: content,
		UserID:  user.ID,
	}
	
	if err := s.messageRepository.Save(message); err != nil {
		return nil, fmt.Errorf("failed to save message: %w", err)
	}
	
	// Notify user about the new message
	notificationText := fmt.Sprintf("You have created a new message: %s", content)
	if err := s.notification.SendNotification(userID, notificationText); err != nil {
		s.logger.Log(fmt.Sprintf("Failed to send notification: %v", err))
		// Continue despite notification failure
	}
	
	return message, nil
}

// GetUserMessages gets all messages for a user
func (s *MessageService) GetUserMessages(userID string) ([]*Message, error) {
	s.logger.Log(fmt.Sprintf("Getting messages for user: %s", userID))
	
	// Verify user exists
	_, err := s.userRepository.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user: %w", err)
	}
	
	// Get messages
	return s.messageRepository.FindByUserID(userID)
}

// Application is the main application struct that coordinates all services
type Application struct {
	UserService    *UserService
	MessageService *MessageService
	Config         *Config
	Logger         *Logger
	db             *DatabaseConnection
}

// NewApplication creates a new application instance with all dependencies
func NewApplication(
	userService *UserService,
	messageService *MessageService,
	config *Config,
	logger *Logger,
	db *DatabaseConnection,
) *Application {
	logger.Log("Creating application")
	return &Application{
		UserService:    userService,
		MessageService: messageService,
		Config:         config,
		Logger:         logger,
		db:             db,
	}
}

// Run runs the application
func (a *Application) Run() error {
	a.Logger.Log(fmt.Sprintf("Running application in %s environment", a.Config.Environment))
	
	// Create a test user
	user, err := a.UserService.CreateUser("testuser", "test@example.com")
	if err != nil {
		return fmt.Errorf("failed to create test user: %w", err)
	}
	
	// Create a test message
	_, err = a.MessageService.CreateMessage("Hello, Wire DI!", user.ID)
	if err != nil {
		return fmt.Errorf("failed to create test message: %w", err)
	}
	
	// Get user messages
	messages, err := a.MessageService.GetUserMessages(user.ID)
	if err != nil {
		return fmt.Errorf("failed to get user messages: %w", err)
	}
	
	a.Logger.Log(fmt.Sprintf("Found %d messages for user %s", len(messages), user.ID))
	for _, msg := range messages {
		a.Logger.Log(fmt.Sprintf("Message: %s", msg.Content))
	}
	
	return nil
}

// Close closes the application and releases resources
func (a *Application) Close() error {
	a.Logger.Log("Shutting down application")
	return a.db.Close()
}
