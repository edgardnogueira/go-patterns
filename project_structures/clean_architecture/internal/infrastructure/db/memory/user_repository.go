package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/edgardnogueira/go-patterns/project_structures/clean_architecture/internal/domain/entities"
)

// UserMemoryRepository is an in-memory implementation of the UserRepository interface
type UserMemoryRepository struct {
	users map[string]*entities.User
	mutex sync.RWMutex
}

// NewUserMemoryRepository creates a new instance of UserMemoryRepository
func NewUserMemoryRepository() *UserMemoryRepository {
	return &UserMemoryRepository{
		users: make(map[string]*entities.User),
	}
}

// Create adds a new user to the in-memory storage
func (r *UserMemoryRepository) Create(_ context.Context, user *entities.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.users[user.ID]; exists {
		return errors.New("user already exists")
	}

	// Create a deep copy of the user to avoid external modifications
	userCopy := *user
	r.users[user.ID] = &userCopy

	return nil
}

// GetByID retrieves a user by ID from the in-memory storage
func (r *UserMemoryRepository) GetByID(_ context.Context, id string) (*entities.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, nil
	}

	// Return a copy to prevent external modifications
	userCopy := *user
	return &userCopy, nil
}

// GetByUsername retrieves a user by username from the in-memory storage
func (r *UserMemoryRepository) GetByUsername(_ context.Context, username string) (*entities.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.Username == username {
			// Return a copy to prevent external modifications
			userCopy := *user
			return &userCopy, nil
		}
	}

	return nil, nil
}

// GetByEmail retrieves a user by email from the in-memory storage
func (r *UserMemoryRepository) GetByEmail(_ context.Context, email string) (*entities.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			// Return a copy to prevent external modifications
			userCopy := *user
			return &userCopy, nil
		}
	}

	return nil, nil
}

// Update updates an existing user in the in-memory storage
func (r *UserMemoryRepository) Update(_ context.Context, user *entities.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return errors.New("user not found")
	}

	// Create a deep copy of the user to avoid external modifications
	userCopy := *user
	r.users[user.ID] = &userCopy

	return nil
}

// Delete removes a user from the in-memory storage
func (r *UserMemoryRepository) Delete(_ context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.users[id]; !exists {
		return errors.New("user not found")
	}

	delete(r.users, id)
	return nil
}

// List retrieves users with pagination from the in-memory storage
func (r *UserMemoryRepository) List(_ context.Context, limit, offset int) ([]*entities.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Convert map to slice
	users := make([]*entities.User, 0, len(r.users))
	for _, user := range r.users {
		// Create a copy of each user
		userCopy := *user
		users = append(users, &userCopy)
	}

	// Apply pagination
	if offset >= len(users) {
		return []*entities.User{}, nil
	}

	end := offset + limit
	if end > len(users) {
		end = len(users)
	}

	return users[offset:end], nil
}

// Count returns the total number of users in the in-memory storage
func (r *UserMemoryRepository) Count(_ context.Context) (int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return len(r.users), nil
}
