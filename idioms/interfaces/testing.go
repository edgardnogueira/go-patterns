// Package interfaces demonstrates idiomatic Go interface implementation patterns
package interfaces

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Using Interfaces for Testing
// ---------------------------
// Interfaces create natural seams for testing by allowing the substitution
// of real implementations with test doubles (mocks, stubs, fakes, etc.)

// DataStore defines the interface for data storage operations
type DataStore interface {
	Get(id string) ([]byte, error)
	Set(id string, data []byte) error
	Delete(id string) error
}

// UserRepository handles user data operations
type UserRepository interface {
	FindByID(id string) (*User, error)
	Save(user *User) error
	Delete(id string) error
}

// EmailService defines the interface for sending emails
type EmailService interface {
	SendEmail(to, subject, body string) error
}

// User represents a user in the system
type User struct {
	ID       string
	Username string
	Email    string
	Created  time.Time
}

// UserService implements business logic for user operations
type UserService struct {
	repo  UserRepository
	email EmailService
}

// NewUserService creates a new UserService
func NewUserService(repo UserRepository, email EmailService) *UserService {
	return &UserService{
		repo:  repo,
		email: email,
	}
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(id string) (*User, error) {
	return s.repo.FindByID(id)
}

// CreateUser creates a new user and sends a welcome email
func (s *UserService) CreateUser(username, email string) (*User, error) {
	// Validate inputs
	if username == "" || email == "" {
		return nil, errors.New("username and email are required")
	}
	
	// Create user
	user := &User{
		ID:       generateID(),
		Username: username,
		Email:    email,
		Created:  time.Now(),
	}
	
	// Save user
	if err := s.repo.Save(user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}
	
	// Send welcome email
	emailBody := fmt.Sprintf("Welcome, %s! Your account has been created.", username)
	if err := s.email.SendEmail(email, "Welcome to Our Service", emailBody); err != nil {
		// Log but don't fail the operation if email fails
		fmt.Printf("Failed to send welcome email: %v\n", err)
	}
	
	return user, nil
}

// DeleteUser deletes a user and sends a confirmation email
func (s *UserService) DeleteUser(id string) error {
	// Get user first to have their email
	user, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	
	// Delete the user
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	
	// Send confirmation email
	emailBody := fmt.Sprintf("Hello %s, your account has been deleted successfully.", user.Username)
	if err := s.email.SendEmail(user.Email, "Account Deleted", emailBody); err != nil {
		// Log but don't fail the operation if email fails
		fmt.Printf("Failed to send deletion confirmation email: %v\n", err)
	}
	
	return nil
}

// Helper function to generate an ID
func generateID() string {
	return fmt.Sprintf("user-%d", time.Now().UnixNano())
}

// Implementation of the interfaces for production
// ---------------------------------------------

// DatabaseUserRepository is a real implementation of UserRepository
type DatabaseUserRepository struct {
	db DataStore
}

// NewDatabaseUserRepository creates a new DatabaseUserRepository
func NewDatabaseUserRepository(db DataStore) *DatabaseUserRepository {
	return &DatabaseUserRepository{db: db}
}

// FindByID retrieves a user from the database
func (r *DatabaseUserRepository) FindByID(id string) (*User, error) {
	data, err := r.db.Get(id)
	if err != nil {
		return nil, err
	}
	
	// In a real implementation, we would deserialize the data
	// For this example, we'll create a user directly
	parts := strings.Split(string(data), ",")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid user data format")
	}
	
	created, _ := time.Parse(time.RFC3339, parts[2])
	return &User{
		ID:       id,
		Username: parts[0],
		Email:    parts[1],
		Created:  created,
	}, nil
}

// Save stores a user in the database
func (r *DatabaseUserRepository) Save(user *User) error {
	// In a real implementation, we would serialize the user
	// For this example, we'll create a simple string
	data := fmt.Sprintf("%s,%s,%s", 
		user.Username, 
		user.Email, 
		user.Created.Format(time.RFC3339),
	)
	
	return r.db.Set(user.ID, []byte(data))
}

// Delete removes a user from the database
func (r *DatabaseUserRepository) Delete(id string) error {
	return r.db.Delete(id)
}

// SmtpEmailService is a real implementation of EmailService
type SmtpEmailService struct {
	smtpServer string
	port       int
	username   string
	password   string
}

// NewSmtpEmailService creates a new SmtpEmailService
func NewSmtpEmailService(server string, port int, username, password string) *SmtpEmailService {
	return &SmtpEmailService{
		smtpServer: server,
		port:       port,
		username:   username,
		password:   password,
	}
}

// SendEmail sends an email via SMTP
func (s *SmtpEmailService) SendEmail(to, subject, body string) error {
	// In a real implementation, this would connect to an SMTP server
	// For this example, we'll just print the email
	fmt.Printf("Sending email via SMTP:\n")
	fmt.Printf("  Server: %s:%d\n", s.smtpServer, s.port)
	fmt.Printf("  To: %s\n", to)
	fmt.Printf("  Subject: %s\n", subject)
	fmt.Printf("  Body: %s\n", body)
	
	return nil
}

// Test doubles for testing
// ----------------------

// MockUserRepository is a test double for UserRepository
type MockUserRepository struct {
	users       map[string]*User
	findErr     error
	saveErr     error
	deleteErr   error
	saveCount   int
	deleteCount int
	findCount   int
}

// NewMockUserRepository creates a new MockUserRepository
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*User),
	}
}

// SetFindError sets an error to be returned by FindByID
func (m *MockUserRepository) SetFindError(err error) {
	m.findErr = err
}

// SetSaveError sets an error to be returned by Save
func (m *MockUserRepository) SetSaveError(err error) {
	m.saveErr = err
}

// SetDeleteError sets an error to be returned by Delete
func (m *MockUserRepository) SetDeleteError(err error) {
	m.deleteErr = err
}

// AddUser adds a user to the mock repository
func (m *MockUserRepository) AddUser(user *User) {
	m.users[user.ID] = user
}

// GetSaveCount returns the number of times Save was called
func (m *MockUserRepository) GetSaveCount() int {
	return m.saveCount
}

// GetDeleteCount returns the number of times Delete was called
func (m *MockUserRepository) GetDeleteCount() int {
	return m.deleteCount
}

// GetFindCount returns the number of times FindByID was called
func (m *MockUserRepository) GetFindCount() int {
	return m.findCount
}

// FindByID implements UserRepository.FindByID for testing
func (m *MockUserRepository) FindByID(id string) (*User, error) {
	m.findCount++
	
	if m.findErr != nil {
		return nil, m.findErr
	}
	
	user, ok := m.users[id]
	if !ok {
		return nil, fmt.Errorf("user not found: %s", id)
	}
	
	return user, nil
}

// Save implements UserRepository.Save for testing
func (m *MockUserRepository) Save(user *User) error {
	m.saveCount++
	
	if m.saveErr != nil {
		return m.saveErr
	}
	
	m.users[user.ID] = user
	return nil
}

// Delete implements UserRepository.Delete for testing
func (m *MockUserRepository) Delete(id string) error {
	m.deleteCount++
	
	if m.deleteErr != nil {
		return m.deleteErr
	}
	
	if _, ok := m.users[id]; !ok {
		return fmt.Errorf("user not found: %s", id)
	}
	
	delete(m.users, id)
	return nil
}

// MockEmailService is a test double for EmailService
type MockEmailService struct {
	emails    []EmailRecord
	shouldErr bool
}

// EmailRecord represents a sent email
type EmailRecord struct {
	To      string
	Subject string
	Body    string
}

// NewMockEmailService creates a new MockEmailService
func NewMockEmailService() *MockEmailService {
	return &MockEmailService{
		emails: []EmailRecord{},
	}
}

// SetShouldError configures whether the mock should return an error
func (m *MockEmailService) SetShouldError(shouldErr bool) {
	m.shouldErr = shouldErr
}

// GetSentEmails returns all sent emails
func (m *MockEmailService) GetSentEmails() []EmailRecord {
	return m.emails
}

// SendEmail implements EmailService.SendEmail for testing
func (m *MockEmailService) SendEmail(to, subject, body string) error {
	if m.shouldErr {
		return errors.New("simulated email error")
	}
	
	m.emails = append(m.emails, EmailRecord{
		To:      to,
		Subject: subject,
		Body:    body,
	})
	
	return nil
}

// Example of testing with mocks
// ---------------------------

// TestCreateUser demonstrates how to test with mock dependencies
func TestCreateUser() {
	// Create mock dependencies
	mockRepo := NewMockUserRepository()
	mockEmail := NewMockEmailService()
	
	// Create the service with mock dependencies
	service := NewUserService(mockRepo, mockEmail)
	
	// Test successful case
	user, err := service.CreateUser("testuser", "test@example.com")
	
	if err != nil {
		fmt.Println("Test failed - expected no error but got:", err)
		return
	}
	
	// Verify user was saved
	if mockRepo.GetSaveCount() != 1 {
		fmt.Println("Test failed - Save should be called once")
		return
	}
	
	// Verify email was sent
	emails := mockEmail.GetSentEmails()
	if len(emails) != 1 {
		fmt.Println("Test failed - should send exactly one email")
		return
	}
	
	if emails[0].To != "test@example.com" {
		fmt.Println("Test failed - wrong email recipient")
		return
	}
	
	fmt.Println("Create user test passed")
	fmt.Printf("Created user: %+v\n", user)
	fmt.Printf("Email sent: %+v\n", emails[0])
	
	// Test error case
	mockRepo.SetSaveError(errors.New("database connection error"))
	
	_, err = service.CreateUser("anothertestuser", "test2@example.com")
	
	if err == nil {
		fmt.Println("Test failed - expected error but got none")
		return
	}
	
	fmt.Println("Error case test passed")
	fmt.Println("Error:", err)
}

// TestDeleteUser demonstrates more advanced testing with mocks
func TestDeleteUser() {
	// Create mock dependencies
	mockRepo := NewMockUserRepository()
	mockEmail := NewMockEmailService()
	
	// Add a test user to the repo
	testUser := &User{
		ID:       "user-123",
		Username: "testuser",
		Email:    "test@example.com",
		Created:  time.Now(),
	}
	mockRepo.AddUser(testUser)
	
	// Create the service with mock dependencies
	service := NewUserService(mockRepo, mockEmail)
	
	// Test successful case
	err := service.DeleteUser("user-123")
	
	if err != nil {
		fmt.Println("Delete test failed - expected no error but got:", err)
		return
	}
	
	// Verify repository calls
	if mockRepo.GetFindCount() != 1 {
		fmt.Println("Test failed - FindByID should be called once")
		return
	}
	
	if mockRepo.GetDeleteCount() != 1 {
		fmt.Println("Test failed - Delete should be called once")
		return
	}
	
	// Verify email was sent
	emails := mockEmail.GetSentEmails()
	if len(emails) != 1 {
		fmt.Println("Test failed - should send exactly one email")
		return
	}
	
	fmt.Println("Delete user test passed")
	fmt.Printf("Email sent: %+v\n", emails[0])
	
	// Test user-not-found case
	err = service.DeleteUser("nonexistent-user")
	
	if err == nil {
		fmt.Println("Test failed - expected error but got none")
		return
	}
	
	fmt.Println("User not found test passed")
	fmt.Println("Error:", err)
}

// InterfaceTestingDemo demonstrates using interfaces for testing
func InterfaceTestingDemo() {
	fmt.Println("============================================")
	fmt.Println("Interface Testing Pattern Demo")
	fmt.Println("============================================")
	
	// Run tests
	fmt.Println("Running CreateUser test:")
	TestCreateUser()
	
	fmt.Println("\nRunning DeleteUser test:")
	TestDeleteUser()
	
	// Show real implementations
	fmt.Println("\nReal implementations would be used in production:")
	
	// Create real dependencies
	// db := some real database implementation
	// repo := NewDatabaseUserRepository(db)
	// email := NewSmtpEmailService("smtp.example.com", 587, "user", "pass")
	// service := NewUserService(repo, email)
	
	fmt.Println("Database repository would connect to a real database")
	fmt.Println("SMTP service would connect to a real email server")
	
	fmt.Println("\nBenefits of using interfaces for testing:")
	fmt.Println("1. Isolate the code being tested from external dependencies")
	fmt.Println("2. Control the behavior of dependencies for testing edge cases")
	fmt.Println("3. Test business logic without network calls, database connections, etc.")
	fmt.Println("4. Mock expensive or complex operations for faster tests")
	fmt.Println("5. Verify interactions with dependencies")
	fmt.Println("============================================")
}
