package main

import (
	"os"
	"testing"
)

// MockLogger is a test implementation of the logger
type MockLogger struct {
	LogMessages   []string
	ErrorMessages []string
}

// Log records log messages
func (l *MockLogger) Log(message string) {
	l.LogMessages = append(l.LogMessages, message)
}

// LogError records error messages
func (l *MockLogger) LogError(message string, err error) {
	l.ErrorMessages = append(l.ErrorMessages, message)
}

// MockEmailService is a test implementation of the email service
type MockEmailService struct {
	SentEmails     []User
	ShouldFail     bool
	ErrorToReturn  error
}

// SendWelcomeEmail records the email sending event
func (e *MockEmailService) SendWelcomeEmail(user User) error {
	if e.ShouldFail {
		return e.ErrorToReturn
	}
	e.SentEmails = append(e.SentEmails, user)
	return nil
}

// MockRepository is a test implementation of the repository
type MockRepository struct {
	SavedUsers [][]User
	ShouldFail bool
	ErrorToReturn error
}

// SaveUsers records the save operation
func (r *MockRepository) SaveUsers(users []User) error {
	if r.ShouldFail {
		return r.ErrorToReturn
	}
	usersCopy := make([]User, len(users))
	copy(usersCopy, users)
	r.SavedUsers = append(r.SavedUsers, usersCopy)
	return nil
}

// LoadUsers returns an empty list or an error if set to fail
func (r *MockRepository) LoadUsers() ([]User, error) {
	if r.ShouldFail {
		return nil, r.ErrorToReturn
	}
	if len(r.SavedUsers) == 0 {
		return []User{}, nil
	}
	return r.SavedUsers[len(r.SavedUsers)-1], nil
}

func TestUserServiceCreateUser(t *testing.T) {
	// Setup
	mockLogger := &MockLogger{}
	mockEmailer := &MockEmailService{}
	mockRepo := &MockRepository{}
	validator := &UserValidator{}
	
	userService := NewUserService(mockRepo, validator, mockEmailer, mockLogger)

	validUser := User{
		ID:       1,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	// Test successful user creation
	t.Run("successful user creation", func(t *testing.T) {
		err := userService.CreateUser(validUser)
		
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		// Check that user was saved
		if len(mockRepo.SavedUsers) != 1 {
			t.Errorf("Expected 1 save operation, got %d", len(mockRepo.SavedUsers))
		}
		
		// Check that email was sent
		if len(mockEmailer.SentEmails) != 1 {
			t.Errorf("Expected 1 email to be sent, got %d", len(mockEmailer.SentEmails))
		}
		
		// Check that success was logged
		if len(mockLogger.LogMessages) != 1 {
			t.Errorf("Expected 1 log message, got %d", len(mockLogger.LogMessages))
		}
	})

	// Reset mocks
	mockLogger = &MockLogger{}
	mockEmailer = &MockEmailService{ShouldFail: true, ErrorToReturn: os.ErrNotExist}
	mockRepo = &MockRepository{}
	userService = NewUserService(mockRepo, validator, mockEmailer, mockLogger)

	// Test email sending failure
	t.Run("email sending failure", func(t *testing.T) {
		err := userService.CreateUser(validUser)
		
		// Operation should still succeed even if email fails
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		// Check that user was saved
		if len(mockRepo.SavedUsers) != 1 {
			t.Errorf("Expected 1 save operation, got %d", len(mockRepo.SavedUsers))
		}
		
		// Check that error was logged
		if len(mockLogger.ErrorMessages) != 1 {
			t.Errorf("Expected 1 error message, got %d", len(mockLogger.ErrorMessages))
		}
	})

	// Reset mocks
	mockLogger = &MockLogger{}
	mockEmailer = &MockEmailService{}
	mockRepo = &MockRepository{ShouldFail: true, ErrorToReturn: os.ErrPermission}
	userService = NewUserService(mockRepo, validator, mockEmailer, mockLogger)

	// Test repository failure
	t.Run("repository failure", func(t *testing.T) {
		err := userService.CreateUser(validUser)
		
		// Operation should fail
		if err == nil {
			t.Error("Expected an error, got nil")
		}
		
		// Check that no email was sent
		if len(mockEmailer.SentEmails) != 0 {
			t.Errorf("Expected 0 emails to be sent, got %d", len(mockEmailer.SentEmails))
		}
		
		// Check that error was logged
		if len(mockLogger.ErrorMessages) != 1 {
			t.Errorf("Expected 1 error message, got %d", len(mockLogger.ErrorMessages))
		}
	})

	// Test validation failure
	t.Run("validation failure", func(t *testing.T) {
		invalidUser := User{
			ID:       2,
			Name:     "", // Name is required
			Email:    "test@example.com",
			Password: "password123",
		}

		err := userService.CreateUser(invalidUser)
		
		// Operation should fail
		if err == nil {
			t.Error("Expected a validation error, got nil")
		}
		
		// Check that no user was saved
		if len(mockRepo.SavedUsers) != 0 {
			t.Errorf("Expected 0 save operations, got %d", len(mockRepo.SavedUsers))
		}
		
		// Check that no email was sent
		if len(mockEmailer.SentEmails) != 0 {
			t.Errorf("Expected 0 emails to be sent, got %d", len(mockEmailer.SentEmails))
		}
		
		// Check that error was logged
		if len(mockLogger.ErrorMessages) != 1 {
			t.Errorf("Expected 1 error message, got %d", len(mockLogger.ErrorMessages))
		}
	})
}

func TestUserServiceLoadUsers(t *testing.T) {
	// Setup
	mockLogger := &MockLogger{}
	mockEmailer := &MockEmailService{}
	mockRepo := &MockRepository{}
	validator := &UserValidator{}
	
	userService := NewUserService(mockRepo, validator, mockEmailer, mockLogger)

	// Prepare test data
	testUsers := []User{
		{ID: 1, Name: "User 1", Email: "user1@example.com", Password: "password123"},
		{ID: 2, Name: "User 2", Email: "user2@example.com", Password: "password456"},
	}
	mockRepo.SavedUsers = append(mockRepo.SavedUsers, testUsers)

	// Test successful loading
	t.Run("successful load", func(t *testing.T) {
		err := userService.LoadUsers()
		
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		// Check that users were loaded
		if len(userService.users) != 2 {
			t.Errorf("Expected 2 users to be loaded, got %d", len(userService.users))
		}
		
		// Check that success was logged
		if len(mockLogger.LogMessages) != 1 {
			t.Errorf("Expected 1 log message, got %d", len(mockLogger.LogMessages))
		}
	})

	// Reset mocks
	mockLogger = &MockLogger{}
	mockRepo = &MockRepository{ShouldFail: true, ErrorToReturn: os.ErrNotExist}
	userService = NewUserService(mockRepo, validator, mockEmailer, mockLogger)

	// Test load failure
	t.Run("load failure", func(t *testing.T) {
		err := userService.LoadUsers()
		
		// Operation should fail
		if err == nil {
			t.Error("Expected an error, got nil")
		}
		
		// Check that error was logged
		if len(mockLogger.ErrorMessages) != 1 {
			t.Errorf("Expected 1 error message, got %d", len(mockLogger.ErrorMessages))
		}
	})
}
