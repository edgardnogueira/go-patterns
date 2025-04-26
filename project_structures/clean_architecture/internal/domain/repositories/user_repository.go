package repositories

import (
	"context"

	"github.com/edgardnogueira/go-patterns/project_structures/clean_architecture/internal/domain/entities"
)

// UserRepository defines the interface for user persistence operations
type UserRepository interface {
	// Create saves a new user
	Create(ctx context.Context, user *entities.User) error

	// GetByID retrieves a user by their ID
	GetByID(ctx context.Context, id string) (*entities.User, error)

	// GetByUsername retrieves a user by their username
	GetByUsername(ctx context.Context, username string) (*entities.User, error)

	// GetByEmail retrieves a user by their email
	GetByEmail(ctx context.Context, email string) (*entities.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *entities.User) error

	// Delete removes a user by their ID
	Delete(ctx context.Context, id string) error

	// List retrieves users with optional pagination
	List(ctx context.Context, limit, offset int) ([]*entities.User, error)

	// Count returns the total number of users
	Count(ctx context.Context) (int, error)
}
