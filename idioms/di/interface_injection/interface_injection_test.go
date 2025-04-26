package interface_injection

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// MockStorage implements Storage interface for testing
type MockStorage struct {
	data         map[string][]byte
	saveCount    int
	loadCount    int
	deleteCount  int
	shouldFail   bool
	failureError error
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		data:         make(map[string][]byte),
		shouldFail:   false,
		failureError: nil,
	}
}

func (s *MockStorage) Save(key string, data []byte) error {
	s.saveCount++
	if s.shouldFail {
		return s.failureError
	}
	s.data[key] = data
	return nil
}

func (s *MockStorage) Load(key string) ([]byte, error) {
	s.loadCount++
	if s.shouldFail {
		return nil, s.failureError
	}
	data, ok := s.data[key]
	if !ok {
		return nil, fmt.Errorf("key not found: %s", key)
	}
	return data, nil
}

func (s *MockStorage) Delete(key string) error {
	s.deleteCount++
	if s.shouldFail {
		return s.failureError
	}
	delete(s.data, key)
	return nil
}

// MockNotifier implements Notifier interface for testing
type MockNotifier struct {
	messages     []string
	shouldFail   bool
	failureError error
}

func (n *MockNotifier) Notify(message string) error {
	n.messages = append(n.messages, message)
	if n.shouldFail {
		return n.failureError
	}
	return nil
}

// MockLogger implements Logger interface for testing
type MockLogger struct {
	logs []struct {
		Level   string
		Message string
	}
}

func (l *MockLogger) Log(level string, message string) {
	l.logs = append(l.logs, struct {
		Level   string
		Message string
	}{Level: level, Message: message})
}

func TestUserServiceOperations(t *testing.T) {
	// Create mock dependencies
	mockStorage := NewMockStorage()
	mockNotifier := &MockNotifier{}
	mockLogger := &MockLogger{}
	appVersion := "1.0.0"
	
	// Create service with injected dependencies
	userService := NewUserService(mockStorage, mockNotifier, mockLogger, appVersion)
	
	// Test CreateUser
	userData := []byte(`{"name":"John Doe","email":"john@example.com"}`)
	err := userService.CreateUser("user1", userData)
	
	// Verify behavior
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if mockStorage.saveCount != 1 {
		t.Errorf("Expected 1 save operation, got %d", mockStorage.saveCount)
	}
	
	if len(mockNotifier.messages) != 1 {
		t.Errorf("Expected 1 notification, got %d", len(mockNotifier.messages))
	}
	
	if !strings.Contains(mockNotifier.messages[0], "user1") {
		t.Errorf("Notification doesn't contain user ID, got: %s", mockNotifier.messages[0])
	}
	
	// Test GetUser
	retrievedData, err := userService.GetUser("user1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if !bytes.Equal(retrievedData, userData) {
		t.Errorf("Expected %s, got %s", userData, retrievedData)
	}
	
	if mockStorage.loadCount != 1 {
		t.Errorf("Expected 1 load operation, got %d", mockStorage.loadCount)
	}
	
	// Test DeleteUser
	err = userService.DeleteUser("user1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if mockStorage.deleteCount != 1 {
		t.Errorf("Expected 1 delete operation, got %d", mockStorage.deleteCount)
	}
	
	if len(mockNotifier.messages) != 2 {
		t.Errorf("Expected 2 notifications, got %d", len(mockNotifier.messages))
	}
	
	// Verify logging
	if len(mockLogger.logs) != 6 { // 2 logs per operation
		t.Errorf("Expected 6 log entries, got %d", len(mockLogger.logs))
	}
}

func TestUserServiceWithErrorHandling(t *testing.T) {
	// Create mock dependencies with failures
	mockStorage := NewMockStorage()
	mockStorage.shouldFail = true
	mockStorage.failureError = fmt.Errorf("storage error")
	
	mockNotifier := &MockNotifier{}
	mockLogger := &MockLogger{}
	
	// Create service with injected dependencies
	userService := NewUserService(mockStorage, mockNotifier, mockLogger, "1.0.0")
	
	// Test with failing storage
	err := userService.CreateUser("user2", []byte(`{"name":"Jane Doe"}`))
	if err == nil {
		t.Error("Expected error from failing storage, got nil")
	}
	
	if !strings.Contains(err.Error(), "storage error") {
		t.Errorf("Expected storage error, got: %v", err)
	}
	
	// Test notification failure (should not cause operation to fail)
	mockStorage.shouldFail = false
	mockNotifier.shouldFail = true
	mockNotifier.failureError = fmt.Errorf("notification error")
	
	err = userService.CreateUser("user3", []byte(`{"name":"Bob Smith"}`))
	if err != nil {
		t.Errorf("Expected success despite notification failure, got: %v", err)
	}
	
	// Verify logging of notification failure
	var foundWarning bool
	for _, log := range mockLogger.logs {
		if log.Level == "WARN" && strings.Contains(log.Message, "Failed to send notification") {
			foundWarning = true
			break
		}
	}
	
	if !foundWarning {
		t.Error("Expected warning log for notification failure")
	}
}

func TestUserServiceWithDifferentImplementations(t *testing.T) {
	// Test with in-memory storage
	memoryStorage := NewMemoryStorage()
	emailNotifier := &EmailNotifier{SMTPServer: "smtp.example.com", FromEmail: "no-reply@example.com"}
	consoleLogger := &ConsoleLogger{}
	
	memoryService := NewUserService(memoryStorage, emailNotifier, consoleLogger, "1.0.0")
	
	// Test with file storage (just verifying it compiles and runs)
	fileStorage := &FileStorage{BasePath: "/tmp/users"}
	smsNotifier := &SMSNotifier{AccountID: "1234", APIKey: "secret"}
	
	fileService := NewUserService(fileStorage, smsNotifier, consoleLogger, "1.0.0")
	
	// Just demonstrate that we can switch implementations
	// In real tests, we would verify the behavior differences
	userData := []byte(`{"name":"Test User"}`)
	
	// These operations should not fail, but output will differ based on the implementation
	_ = memoryService.CreateUser("memory-user", userData)
	_, _ = memoryService.GetUser("memory-user")
	
	_ = fileService.CreateUser("file-user", userData)
	_, _ = fileService.GetUser("file-user")
}
