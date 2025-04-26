package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/edgardnogueira/go-patterns/project_structures/clean_architecture/internal/domain/entities"
)

// TaskMemoryRepository is an in-memory implementation of the TaskRepository interface
type TaskMemoryRepository struct {
	tasks map[string]*entities.Task
	mutex sync.RWMutex
}

// NewTaskMemoryRepository creates a new instance of TaskMemoryRepository
func NewTaskMemoryRepository() *TaskMemoryRepository {
	return &TaskMemoryRepository{
		tasks: make(map[string]*entities.Task),
	}
}

// Create adds a new task to the in-memory storage
func (r *TaskMemoryRepository) Create(_ context.Context, task *entities.Task) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.tasks[task.ID]; exists {
		return errors.New("task already exists")
	}

	// Create a deep copy of the task to avoid external modifications
	taskCopy := *task
	r.tasks[task.ID] = &taskCopy

	return nil
}

// GetByID retrieves a task by ID from the in-memory storage
func (r *TaskMemoryRepository) GetByID(_ context.Context, id string) (*entities.Task, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return nil, nil
	}

	// Return a copy to prevent external modifications
	taskCopy := *task
	return &taskCopy, nil
}

// Update updates an existing task in the in-memory storage
func (r *TaskMemoryRepository) Update(_ context.Context, task *entities.Task) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.tasks[task.ID]; !exists {
		return errors.New("task not found")
	}

	// Create a deep copy of the task to avoid external modifications
	taskCopy := *task
	r.tasks[task.ID] = &taskCopy

	return nil
}

// Delete removes a task from the in-memory storage
func (r *TaskMemoryRepository) Delete(_ context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.tasks[id]; !exists {
		return errors.New("task not found")
	}

	delete(r.tasks, id)
	return nil
}

// ListPending retrieves pending tasks from the in-memory storage
func (r *TaskMemoryRepository) ListPending(_ context.Context, limit int) ([]*entities.Task, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var pendingTasks []*entities.Task
	for _, task := range r.tasks {
		if task.Status == entities.TaskStatusPending {
			// Create a copy of each task
			taskCopy := *task
			pendingTasks = append(pendingTasks, &taskCopy)
		}

		if limit > 0 && len(pendingTasks) >= limit {
			break
		}
	}

	return pendingTasks, nil
}

// ListByStatus retrieves tasks by status with pagination from the in-memory storage
func (r *TaskMemoryRepository) ListByStatus(_ context.Context, status entities.TaskStatus, limit, offset int) ([]*entities.Task, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var filteredTasks []*entities.Task
	for _, task := range r.tasks {
		if task.Status == status {
			// Create a copy of each task
			taskCopy := *task
			filteredTasks = append(filteredTasks, &taskCopy)
		}
	}

	// Apply pagination
	if offset >= len(filteredTasks) {
		return []*entities.Task{}, nil
	}

	end := offset + limit
	if end > len(filteredTasks) {
		end = len(filteredTasks)
	}

	return filteredTasks[offset:end], nil
}

// GetNextPendingTask retrieves and locks the next pending task for processing
func (r *TaskMemoryRepository) GetNextPendingTask(_ context.Context) (*entities.Task, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, task := range r.tasks {
		if task.Status == entities.TaskStatusPending {
			// Mark as processing and update timestamp
			task.MarkAsProcessing()

			// Return a copy to prevent external modifications
			taskCopy := *task
			return &taskCopy, nil
		}
	}

	return nil, nil
}

// Count returns the total number of tasks by status in the in-memory storage
func (r *TaskMemoryRepository) Count(_ context.Context, status entities.TaskStatus) (int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	count := 0
	for _, task := range r.tasks {
		if task.Status == status {
			count++
		}
	}

	return count, nil
}
