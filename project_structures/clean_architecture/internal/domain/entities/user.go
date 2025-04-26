package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidUserID       = errors.New("invalid user ID")
	ErrEmptyUsername       = errors.New("username cannot be empty")
	ErrInvalidEmailFormat  = errors.New("invalid email format")
	ErrPasswordTooShort    = errors.New("password must be at least 8 characters")
)

// User represents a user entity in the domain
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Never expose password in JSON
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUser creates a new user entity with validation
func NewUser(username, email, password, firstName, lastName string) (*User, error) {
	if username == "" {
		return nil, ErrEmptyUsername
	}

	if !isValidEmail(email) {
		return nil, ErrInvalidEmailFormat
	}

	if len(password) < 8 {
		return nil, ErrPasswordTooShort
	}

	now := time.Now()
	return &User{
		ID:        uuid.New().String(),
		Username:  username,
		Email:     email,
		Password:  password, // In a real app, this would be hashed
		FirstName: firstName,
		LastName:  lastName,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Update updates user properties with validation
func (u *User) Update(username, email, firstName, lastName string) error {
	if username != "" {
		u.Username = username
	}

	if email != "" {
		if !isValidEmail(email) {
			return ErrInvalidEmailFormat
		}
		u.Email = email
	}

	if firstName != "" {
		u.FirstName = firstName
	}

	if lastName != "" {
		u.LastName = lastName
	}

	u.UpdatedAt = time.Now()
	return nil
}

// UpdatePassword updates the user's password with validation
func (u *User) UpdatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}
	
	u.Password = password // In a real app, this would be hashed
	u.UpdatedAt = time.Now()
	return nil
}

// FullName returns the user's full name
func (u *User) FullName() string {
	if u.FirstName == "" && u.LastName == "" {
		return ""
	}
	
	if u.FirstName == "" {
		return u.LastName
	}
	
	if u.LastName == "" {
		return u.FirstName
	}
	
	return u.FirstName + " " + u.LastName
}

// isValidEmail validates email format
func isValidEmail(email string) bool {
	// This is a simple validation for demonstration purposes
	// In a real app, you would use a more comprehensive validation
	return email != "" && len(email) > 5 && contains(email, "@")
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr
}
