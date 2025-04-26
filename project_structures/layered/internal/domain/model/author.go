package model

import (
	"errors"
	"regexp"
	"time"
)

// Common domain errors for Author entity
var (
	ErrEmptyName  = errors.New("author name cannot be empty")
	ErrInvalidEmail = errors.New("invalid email format")
)

// Author represents an author entity in the domain layer
type Author struct {
	ID        int64
	Name      string
	Email     string
	Bio       string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// emailRegex is a simple regex to validate email format
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// NewAuthor creates a new author with validation
func NewAuthor(name, email, bio string) (*Author, error) {
	if name == "" {
		return nil, ErrEmptyName
	}
	
	if !emailRegex.MatchString(email) {
		return nil, ErrInvalidEmail
	}
	
	now := time.Now()
	
	return &Author{
		Name:      name,
		Email:     email,
		Bio:       bio,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Validate ensures the author is in a valid state
func (a *Author) Validate() error {
	if a.Name == "" {
		return ErrEmptyName
	}
	
	if !emailRegex.MatchString(a.Email) {
		return ErrInvalidEmail
	}
	
	return nil
}

// Update updates the author information with validation
func (a *Author) Update(name, email, bio string) error {
	if name == "" {
		return ErrEmptyName
	}
	
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}
	
	a.Name = name
	a.Email = email
	a.Bio = bio
	a.UpdatedAt = time.Now()
	
	return nil
}

// UpdateBio updates only the author's bio
func (a *Author) UpdateBio(bio string) {
	a.Bio = bio
	a.UpdatedAt = time.Now()
}
