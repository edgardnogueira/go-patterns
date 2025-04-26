package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// TaskStatus represents the current status of a task
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusProcessing TaskStatus = "processing"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusFailed     TaskStatus = "failed"
)

var (
	ErrInvalidTaskType = errors.New("invalid task type")
	ErrEmptyTaskData   = errors.New("task data cannot be empty")
)

// Task represents a background task in the domain
type Task struct {
	ID          string      `json:"id"`
	Type        string      `json:"type"`
	Data        []byte      `json:"data"`
	Status      TaskStatus  `json:"status"`
	Error       string      `json:"error,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	StartedAt   *time.Time  `json:"started_at,omitempty"`
	CompletedAt *time.Time  `json:"completed_at,omitempty"`
}

// NewTask creates a new task with validation
func NewTask(taskType string, data []byte) (*Task, error) {
	if taskType == "" {
		return nil, ErrInvalidTaskType
	}

	if len(data) == 0 {
		return nil, ErrEmptyTaskData
	}

	now := time.Now()
	return &Task{
		ID:        uuid.New().String(),
		Type:      taskType,
		Data:      data,
		Status:    TaskStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// MarkAsProcessing marks the task as processing
func (t *Task) MarkAsProcessing() {
	now := time.Now()
	t.Status = TaskStatusProcessing
	t.StartedAt = &now
	t.UpdatedAt = now
}

// Complete marks the task as completed
func (t *Task) Complete() {
	now := time.Now()
	t.Status = TaskStatusCompleted
	t.CompletedAt = &now
	t.UpdatedAt = now
}

// Fail marks the task as failed with an error message
func (t *Task) Fail(errorMsg string) {
	now := time.Now()
	t.Status = TaskStatusFailed
	t.Error = errorMsg
	t.CompletedAt = &now
	t.UpdatedAt = now
}

// IsActive returns true if the task is pending or processing
func (t *Task) IsActive() bool {
	return t.Status == TaskStatusPending || t.Status == TaskStatusProcessing
}

// Duration returns the duration of the task processing if available
func (t *Task) Duration() time.Duration {
	if t.StartedAt == nil {
		return 0
	}

	if t.CompletedAt == nil {
		return time.Since(*t.StartedAt)
	}

	return t.CompletedAt.Sub(*t.StartedAt)
}
