package usecase

import (
	"context"
	"errors"

	"github.com/edgardnogueira/go-patterns/project_structures/clean_architecture/internal/domain/entities"
	"github.com/edgardnogueira/go-patterns/project_structures/clean_architecture/internal/domain/repositories"
)

// TaskService error definitions
var (
	ErrTaskNotFound = errors.New("task not found")
	ErrNoTaskAvailable = errors.New("no tasks available for processing")
)

// TaskService implements task-related use cases
type TaskService struct {
	taskRepo repositories.TaskRepository
}

// NewTaskService creates a new instance of TaskService
func NewTaskService(taskRepo repositories.TaskRepository) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
	}
}

// GetTaskByID retrieves a task by its ID
func (s *TaskService) GetTaskByID(ctx context.Context, id string) (*entities.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

// GetNextPendingTask retrieves the next pending task for processing
func (s *TaskService) GetNextPendingTask(ctx context.Context) (*entities.Task, error) {
	task, err := s.taskRepo.GetNextPendingTask(ctx)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrNoTaskAvailable
	}
	return task, nil
}

// ListTasksByStatus retrieves tasks by status with pagination
func (s *TaskService) ListTasksByStatus(ctx context.Context, status entities.TaskStatus, limit, offset int) ([]*entities.Task, int, error) {
	tasks, err := s.taskRepo.ListByStatus(ctx, status, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.taskRepo.Count(ctx, status)
	if err != nil {
		return nil, 0, err
	}

	return tasks, count, nil
}

// StartTask marks a task as processing
func (s *TaskService) StartTask(ctx context.Context, id string) (*entities.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}

	task.MarkAsProcessing()
	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

// CompleteTask marks a task as completed
func (s *TaskService) CompleteTask(ctx context.Context, id string) (*entities.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}

	task.Complete()
	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

// FailTask marks a task as failed with an error message
func (s *TaskService) FailTask(ctx context.Context, id string, errorMsg string) (*entities.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}

	task.Fail(errorMsg)
	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

// DeleteTask deletes a task by ID
func (s *TaskService) DeleteTask(ctx context.Context, id string) error {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if task == nil {
		return ErrTaskNotFound
	}

	return s.taskRepo.Delete(ctx, id)
}

// GetPendingTasksCount returns the count of pending tasks
func (s *TaskService) GetPendingTasksCount(ctx context.Context) (int, error) {
	return s.taskRepo.Count(ctx, entities.TaskStatusPending)
}

// ProcessTask is a method that would be used by workers to process a task
// It's a placeholder and would be implemented with specific task handling logic
func (s *TaskService) ProcessTask(ctx context.Context, task *entities.Task) error {
	// This would handle different task types with different processing logic
	// For simplicity, we're just updating the status to completed
	task.Complete()
	return s.taskRepo.Update(ctx, task)
}
