package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// User represents a user in the system
// This struct only holds data and has no responsibilities
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserValidator is responsible only for validating user data
// SRP: This type has a single responsibility - validation
type UserValidator struct{}

// ValidateUser validates user data
func (v *UserValidator) ValidateUser(user User) error {
	if user.Name == "" {
		return fmt.Errorf("user name cannot be empty")
	}
	if user.Email == "" {
		return fmt.Errorf("user email cannot be empty")
	}
	if len(user.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	return nil
}

// UserRepository is responsible only for user data persistence
// SRP: This type has a single responsibility - data storage and retrieval
type UserRepository struct {
	filename string
}

// NewUserRepository creates a new user repository
func NewUserRepository(filename string) *UserRepository {
	return &UserRepository{
		filename: filename,
	}
}

// SaveUsers saves users to a file
func (r *UserRepository) SaveUsers(users []User) error {
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(r.filename, data, 0644)
}

// LoadUsers loads users from a file
func (r *UserRepository) LoadUsers() ([]User, error) {
	var users []User
	data, err := ioutil.ReadFile(r.filename)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist yet, that's fine for a new system
			return []User{}, nil
		}
		return nil, err
	}

	err = json.Unmarshal(data, &users)
	return users, err
}

// EmailService is responsible only for sending emails
// SRP: This type has a single responsibility - email communication
type EmailService struct{}

// SendWelcomeEmail sends a welcome email to a user
func (e *EmailService) SendWelcomeEmail(user User) error {
	// In a real application, this would connect to an email service
	fmt.Printf("Sending welcome email to %s at %s\n", user.Name, user.Email)
	// Simulate email sending
	return nil
}

// Logger is responsible only for logging
// SRP: This type has a single responsibility - logging
type Logger struct{}

// Log logs a message
func (l *Logger) Log(message string) {
	fmt.Printf("[LOG] %s\n", message)
}

// LogError logs an error
func (l *Logger) LogError(message string, err error) {
	fmt.Printf("[ERROR] %s: %v\n", message, err)
}

// UserService coordinates the different components but doesn't implement their functionality
// This is a facade that uses the specialized components
type UserService struct {
	users      []User
	validator  *UserValidator
	repository *UserRepository
	emailer    *EmailService
	logger     *Logger
}

// NewUserService creates a new user service with all dependencies
func NewUserService(repository *UserRepository, validator *UserValidator, emailer *EmailService, logger *Logger) *UserService {
	return &UserService{
		users:      []User{},
		validator:  validator,
		repository: repository,
		emailer:    emailer,
		logger:     logger,
	}
}

// LoadUsers loads all users from storage
func (s *UserService) LoadUsers() error {
	users, err := s.repository.LoadUsers()
	if err != nil {
		s.logger.LogError("Failed to load users", err)
		return err
	}
	s.users = users
	s.logger.Log(fmt.Sprintf("Loaded %d users", len(users)))
	return nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(user User) error {
	// Validate user
	if err := s.validator.ValidateUser(user); err != nil {
		s.logger.LogError("Validation error", err)
		return err
	}

	// Add user to the list
	s.users = append(s.users, user)

	// Save to repository
	if err := s.repository.SaveUsers(s.users); err != nil {
		s.logger.LogError("Failed to save users", err)
		return fmt.Errorf("failed to save user: %w", err)
	}

	// Send welcome email
	if err := s.emailer.SendWelcomeEmail(user); err != nil {
		s.logger.LogError("Failed to send welcome email", err)
		// Continue even if email fails
	}

	s.logger.Log(fmt.Sprintf("User created successfully: %s", user.Name))
	return nil
}

// GetUserByID finds a user by their ID
func (s *UserService) GetUserByID(id int) (User, error) {
	for _, user := range s.users {
		if user.ID == id {
			return user, nil
		}
	}
	return User{}, fmt.Errorf("user with ID %d not found", id)
}

// This function demonstrates using the UserService with proper SRP
func demonstrateUserManagerAfterSRP() {
	// Create all the necessary components with single responsibilities
	logger := &Logger{}
	validator := &UserValidator{}
	repository := NewUserRepository("users.json")
	emailer := &EmailService{}
	
	// Create the service that coordinates these components
	userService := NewUserService(repository, validator, emailer, logger)

	// Try to load existing users
	if err := userService.LoadUsers(); err != nil {
		logger.LogError("Error loading users", err)
	}

	// Create a new user
	newUser := User{
		ID:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	if err := userService.CreateUser(newUser); err != nil {
		logger.LogError("Error creating user", err)
		return
	}

	fmt.Println("User created successfully!")

	// Try to retrieve the user
	user, err := userService.GetUserByID(1)
	if err != nil {
		logger.LogError("Error retrieving user", err)
		return
	}

	fmt.Printf("Retrieved user: %s (%s)\n", user.Name, user.Email)
}
