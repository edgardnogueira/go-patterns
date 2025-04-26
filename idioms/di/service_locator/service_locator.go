// Package service_locator demonstrates the service locator pattern for dependency management.
package service_locator

import (
	"fmt"
	"sync"
)

// ServiceLocator provides a central registry for services/dependencies
type ServiceLocator struct {
	mu       sync.RWMutex
	services map[string]interface{}
}

// NewServiceLocator creates a new service locator
func NewServiceLocator() *ServiceLocator {
	return &ServiceLocator{
		services: make(map[string]interface{}),
	}
}

// Register adds a service to the locator
func (sl *ServiceLocator) Register(name string, service interface{}) {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	sl.services[name] = service
}

// Get retrieves a service by name
func (sl *ServiceLocator) Get(name string) (interface{}, error) {
	sl.mu.RLock()
	defer sl.mu.RUnlock()
	
	service, exists := sl.services[name]
	if !exists {
		return nil, fmt.Errorf("service not found: %s", name)
	}
	
	return service, nil
}

// GetTyped retrieves a service by name and casts it to the expected type
func (sl *ServiceLocator) GetTyped(name string, dest interface{}) error {
	sl.mu.RLock()
	defer sl.mu.RUnlock()
	
	service, exists := sl.services[name]
	if !exists {
		return fmt.Errorf("service not found: %s", name)
	}
	
	// Use type assertions to cast to the pointer of the expected type
	switch ptr := dest.(type) {
	case *Logger:
		if logger, ok := service.(Logger); ok {
			*ptr = logger
			return nil
		}
	case *UserRepository:
		if repo, ok := service.(UserRepository); ok {
			*ptr = repo
			return nil
		}
	case *NotificationService:
		if notifier, ok := service.(NotificationService); ok {
			*ptr = notifier
			return nil
		}
	// Add other types as needed
	}
	
	return fmt.Errorf("service '%s' cannot be cast to requested type", name)
}

// HasService checks if a service exists
func (sl *ServiceLocator) HasService(name string) bool {
	sl.mu.RLock()
	defer sl.mu.RUnlock()
	_, exists := sl.services[name]
	return exists
}

// Remove removes a service from the locator
func (sl *ServiceLocator) Remove(name string) {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	delete(sl.services, name)
}

//
// Example service interfaces and implementations
//

// Logger defines a logging interface
type Logger interface {
	Log(message string)
}

// ConsoleLogger logs to console
type ConsoleLogger struct{}

func (l *ConsoleLogger) Log(message string) {
	fmt.Println("[LOG]", message)
}

// UserRepository provides user data access
type UserRepository interface {
	FindByID(id string) (User, error)
	Save(user User) error
}

// User represents a user entity
type User struct {
	ID    string
	Name  string
	Email string
}

// InMemoryUserRepository implements UserRepository with in-memory storage
type InMemoryUserRepository struct {
	users map[string]User
	mu    sync.RWMutex
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]User),
	}
}

func (r *InMemoryUserRepository) FindByID(id string) (User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	user, exists := r.users[id]
	if !exists {
		return User{}, fmt.Errorf("user not found: %s", id)
	}
	
	return user, nil
}

func (r *InMemoryUserRepository) Save(user User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.users[user.ID] = user
	return nil
}

// NotificationService handles sending notifications
type NotificationService interface {
	SendNotification(userID, message string) error
}

// EmailNotificationService implements NotificationService via email
type EmailNotificationService struct {
	Logger Logger
}

func (s *EmailNotificationService) SendNotification(userID, message string) error {
	if s.Logger != nil {
		s.Logger.Log(fmt.Sprintf("Sending email notification to user %s: %s", userID, message))
	}
	// Email sending logic would go here
	return nil
}

// UserService uses the service locator to access its dependencies
type UserService struct {
	locator *ServiceLocator
}

// NewUserService creates a new user service with the given service locator
func NewUserService(locator *ServiceLocator) *UserService {
	return &UserService{
		locator: locator,
	}
}

// CreateUser creates a new user and sends a welcome notification
func (s *UserService) CreateUser(id, name, email string) error {
	// Get logger
	loggerService, err := s.locator.Get("logger")
	if err != nil {
		return fmt.Errorf("failed to get logger: %w", err)
	}
	logger, ok := loggerService.(Logger)
	if !ok {
		return fmt.Errorf("invalid logger service type")
	}
	
	logger.Log(fmt.Sprintf("Creating user: %s", name))
	
	// Get user repository
	repoService, err := s.locator.Get("userRepository")
	if err != nil {
		return fmt.Errorf("failed to get user repository: %w", err)
	}
	repo, ok := repoService.(UserRepository)
	if !ok {
		return fmt.Errorf("invalid user repository service type")
	}
	
	// Create and save user
	user := User{
		ID:    id,
		Name:  name,
		Email: email,
	}
	
	if err := repo.Save(user); err != nil {
		logger.Log(fmt.Sprintf("Failed to save user: %v", err))
		return err
	}
	
	// Get notification service
	notificationService, err := s.locator.Get("notificationService")
	if err != nil {
		logger.Log(fmt.Sprintf("Failed to get notification service: %v", err))
		// Continue despite missing notification service
	} else {
		notifier, ok := notificationService.(NotificationService)
		if !ok {
			logger.Log("Invalid notification service type")
		} else {
			// Send welcome notification
			if err := notifier.SendNotification(id, "Welcome to our service!"); err != nil {
				logger.Log(fmt.Sprintf("Failed to send welcome notification: %v", err))
			}
		}
	}
	
	logger.Log(fmt.Sprintf("User created successfully: %s", name))
	return nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(id string) (User, error) {
	// Get user repository
	repoService, err := s.locator.Get("userRepository")
	if err != nil {
		return User{}, fmt.Errorf("failed to get user repository: %w", err)
	}
	repo, ok := repoService.(UserRepository)
	if !ok {
		return User{}, fmt.Errorf("invalid user repository service type")
	}
	
	// Find user
	return repo.FindByID(id)
}
