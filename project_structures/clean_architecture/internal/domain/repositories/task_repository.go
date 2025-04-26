package repositories

import (
	"context"

	"github.com/edgardnogueira/go-patterns/project_structures/clean_architecture/internal/domain/entities"
)

// TaskRepository defines the interface for task persistence operations
type TaskRepository interface {
	// Create saves a new task
	Create(ctx context.Context, task *entities.Task) error

	// GetByID retrieves a task by its ID
	GetByID(ctx context.Context, id string) (*entities.Task, error)

	// Update updates an existing task
	Update(ctx context.Context, task *entities.Task) error

	// Delete removes a task by its ID
	Delete(ctx context.Context, id string) error

	// ListPending retrieves pending tasks with optional limit
	ListPending(ctx context.Context, limit int) ([]*entities.Task, error)

	// ListByStatus retrieves tasks by status with optional pagination
	ListByStatus(ctx context.Context, status entities.TaskStatus, limit, offset int) ([]*entities.Task, error)

	// GetNextPendingTask retrieves and locks the next pending task for processing
	GetNextPendingTask(ctx context.Context) (*entities.Task, error)

	// Count returns the total number of tasks by status
	Count(ctx context.Context, status entities.TaskStatus) (int, error)
}
