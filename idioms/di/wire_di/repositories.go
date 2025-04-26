package wire_di

import (
	"fmt"
	"sync"
)

// UserRepository defines methods for user data operations
type UserRepository interface {
	FindByID(id string) (*User, error)
	Save(user *User) error
}

// MessageRepository defines methods for message data operations
type MessageRepository interface {
	FindByID(id string) (*Message, error)
	FindByUserID(userID string) ([]*Message, error)
	Save(message *Message) error
}

// DBUserRepository is a database implementation of UserRepository
type DBUserRepository struct {
	db     *DatabaseConnection
	logger *Logger
	users  map[string]*User // Simulated storage
	mu     sync.RWMutex
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *DatabaseConnection, logger *Logger) UserRepository {
	logger.Log("Creating user repository")
	return &DBUserRepository{
		db:     db,
		logger: logger,
		users:  make(map[string]*User),
	}
}

// FindByID finds a user by ID
func (r *DBUserRepository) FindByID(id string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	r.logger.Log(fmt.Sprintf("Finding user with ID: %s", id))
	
	user, ok := r.users[id]
	if !ok {
		return nil, fmt.Errorf("user not found: %s", id)
	}
	
	return user, nil
}

// Save saves a user
func (r *DBUserRepository) Save(user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.logger.Log(fmt.Sprintf("Saving user: %s (%s)", user.Username, user.ID))
	
	// In a real implementation, this would save to a database
	r.users[user.ID] = user
	return nil
}

// DBMessageRepository is a database implementation of MessageRepository
type DBMessageRepository struct {
	db       *DatabaseConnection
	logger   *Logger
	messages map[string]*Message // Simulated storage
	userMsgs map[string][]string // User ID -> Message IDs
	mu       sync.RWMutex
}

// NewMessageRepository creates a new message repository
func NewMessageRepository(db *DatabaseConnection, logger *Logger) MessageRepository {
	logger.Log("Creating message repository")
	return &DBMessageRepository{
		db:       db,
		logger:   logger,
		messages: make(map[string]*Message),
		userMsgs: make(map[string][]string),
	}
}

// FindByID finds a message by ID
func (r *DBMessageRepository) FindByID(id string) (*Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	r.logger.Log(fmt.Sprintf("Finding message with ID: %s", id))
	
	message, ok := r.messages[id]
	if !ok {
		return nil, fmt.Errorf("message not found: %s", id)
	}
	
	return message, nil
}

// FindByUserID finds all messages for a user
func (r *DBMessageRepository) FindByUserID(userID string) ([]*Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	r.logger.Log(fmt.Sprintf("Finding messages for user: %s", userID))
	
	messageIDs, ok := r.userMsgs[userID]
	if !ok {
		return []*Message{}, nil
	}
	
	messages := make([]*Message, 0, len(messageIDs))
	for _, id := range messageIDs {
		if msg, ok := r.messages[id]; ok {
			messages = append(messages, msg)
		}
	}
	
	return messages, nil
}

// Save saves a message
func (r *DBMessageRepository) Save(message *Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.logger.Log(fmt.Sprintf("Saving message: %s for user %s", message.ID, message.UserID))
	
	// In a real implementation, this would save to a database
	r.messages[message.ID] = message
	
	// Update user-message index
	r.userMsgs[message.UserID] = append(r.userMsgs[message.UserID], message.ID)
	
	return nil
}
