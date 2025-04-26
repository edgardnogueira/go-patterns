package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// UserManager is a type that handles multiple user-related responsibilities
// This violates the Single Responsibility Principle by doing too much:
// - User data validation
// - User persistence
// - Email notifications
// - Logging
type UserManager struct {
	users []User
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// NewUserManager creates a new user manager
func NewUserManager() *UserManager {
	return &UserManager{
		users: []User{},
	}
}

// ValidateUser checks if the user data is valid
func (um *UserManager) ValidateUser(user User) error {
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

// CreateUser adds a new user and handles all related operations
func (um *UserManager) CreateUser(user User) error {
	// Validate user
	if err := um.ValidateUser(user); err != nil {
		log.Printf("Validation error: %v", err)
		return err
	}

	// Add user to the list
	um.users = append(um.users, user)

	// Save to file
	if err := um.SaveUsersToFile("users.json"); err != nil {
		log.Printf("Failed to save users: %v", err)
		return fmt.Errorf("failed to save user: %w", err)
	}

	// Send welcome email
	if err := um.SendWelcomeEmail(user); err != nil {
		log.Printf("Failed to send welcome email: %v", err)
		// Continue even if email fails
	}

	log.Printf("User created successfully: %s", user.Name)
	return nil
}

// SaveUsersToFile saves all users to a JSON file
func (um *UserManager) SaveUsersToFile(filename string) error {
	data, err := json.MarshalIndent(um.users, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

// LoadUsersFromFile loads users from a JSON file
func (um *UserManager) LoadUsersFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist yet, that's fine for a new system
			return nil
		}
		return err
	}

	return json.Unmarshal(data, &um.users)
}

// SendWelcomeEmail sends a welcome email to the new user
func (um *UserManager) SendWelcomeEmail(user User) error {
	// In a real application, this would connect to an email service
	fmt.Printf("Sending welcome email to %s at %s\n", user.Name, user.Email)
	// Simulate email sending
	return nil
}

// GetUserByID finds a user by their ID
func (um *UserManager) GetUserByID(id int) (User, error) {
	for _, user := range um.users {
		if user.ID == id {
			return user, nil
		}
	}
	return User{}, fmt.Errorf("user with ID %d not found", id)
}

// This function demonstrates using the UserManager
// It violates SRP because it's handling too many responsibilities
func demonstrateUserManagerBeforeSRP() {
	userManager := NewUserManager()

	// Try to load existing users
	if err := userManager.LoadUsersFromFile("users.json"); err != nil {
		fmt.Println("Error loading users:", err)
	}

	// Create a new user
	newUser := User{
		ID:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	if err := userManager.CreateUser(newUser); err != nil {
		fmt.Println("Error creating user:", err)
		return
	}

	fmt.Println("User created successfully!")

	// Try to retrieve the user
	user, err := userManager.GetUserByID(1)
	if err != nil {
		fmt.Println("Error retrieving user:", err)
		return
	}

	fmt.Printf("Retrieved user: %s (%s)\n", user.Name, user.Email)
}
