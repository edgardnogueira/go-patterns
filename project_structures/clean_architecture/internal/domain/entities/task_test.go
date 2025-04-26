package entities

import (
	"testing"
	"time"
)

func TestNewTask(t *testing.T) {
	tests := []struct {
		name     string
		taskType string
		data     []byte
		wantErr  error
	}{
		{
			name:     "Valid task",
			taskType: "test_task",
			data:     []byte(`{"test":"data"}`),
			wantErr:  nil,
		},
		{
			name:     "Empty task type",
			taskType: "",
			data:     []byte(`{"test":"data"}`),
			wantErr:  ErrInvalidTaskType,
		},
		{
			name:     "Empty task data",
			taskType: "test_task",
			data:     []byte{},
			wantErr:  ErrEmptyTaskData,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := NewTask(tt.taskType, tt.data)
			
			// Check error
			if err != tt.wantErr {
				t.Errorf("NewTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			// If we expect an error, no need to check the task
			if tt.wantErr != nil {
				if task != nil {
					t.Errorf("NewTask() returned non-nil task when error occurred")
				}
				return
			}
			
			// Check task fields
			if task.Type != tt.taskType {
				t.Errorf("Task.Type = %v, want %v", task.Type, tt.taskType)
			}
			if string(task.Data) != string(tt.data) {
				t.Errorf("Task.Data = %v, want %v", string(task.Data), string(tt.data))
			}
			if task.Status != TaskStatusPending {
				t.Errorf("Task.Status = %v, want %v", task.Status, TaskStatusPending)
			}
			if task.ID == "" {
				t.Error("Task.ID is empty, expected UUID")
			}
			
			// Timestamps should be set
			if task.CreatedAt.IsZero() {
				t.Error("Task.CreatedAt is zero, expected timestamp")
			}
			if task.UpdatedAt.IsZero() {
				t.Error("Task.UpdatedAt is zero, expected timestamp")
			}
			
			// These should be nil for a new task
			if task.StartedAt != nil {
				t.Error("Task.StartedAt is not nil for new task")
			}
			if task.CompletedAt != nil {
				t.Error("Task.CompletedAt is not nil for new task")
			}
		})
	}
}

func TestTaskMarkAsProcessing(t *testing.T) {
	// Create a test task
	task, err := NewTask("test_task", []byte(`{"test":"data"}`))
	if err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}
	
	// Store original timestamps for comparison
	originalCreatedAt := task.CreatedAt
	originalUpdatedAt := task.UpdatedAt
	
	// Wait a moment to ensure updated timestamp will be different
	time.Sleep(1 * time.Millisecond)
	
	// Mark as processing
	task.MarkAsProcessing()
	
	// Check status was updated
	if task.Status != TaskStatusProcessing {
		t.Errorf("Task.Status = %v, want %v", task.Status, TaskStatusProcessing)
	}
	
	// CreatedAt should remain the same
	if task.CreatedAt != originalCreatedAt {
		t.Errorf("Task.CreatedAt changed after update, should remain constant")
	}
	
	// UpdatedAt should be updated
	if task.UpdatedAt == originalUpdatedAt {
		t.Errorf("Task.UpdatedAt did not change after update")
	}
	
	// StartedAt should be set
	if task.StartedAt == nil {
		t.Error("Task.StartedAt is nil, expected timestamp")
	}
	
	// CompletedAt should still be nil
	if task.CompletedAt != nil {
		t.Error("Task.CompletedAt is not nil after marking as processing")
	}
}

func TestTaskComplete(t *testing.T) {
	// Create a test task
	task, err := NewTask("test_task", []byte(`{"test":"data"}`))
	if err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}
	
	// Mark as processing first
	task.MarkAsProcessing()
	
	// Store timestamps for comparison
	originalStartedAt := task.StartedAt
	originalUpdatedAt := task.UpdatedAt
	
	// Wait a moment to ensure updated timestamp will be different
	time.Sleep(1 * time.Millisecond)
	
	// Complete the task
	task.Complete()
	
	// Check status was updated
	if task.Status != TaskStatusCompleted {
		t.Errorf("Task.Status = %v, want %v", task.Status, TaskStatusCompleted)
	}
	
	// UpdatedAt should be updated
	if task.UpdatedAt == originalUpdatedAt {
		t.Errorf("Task.UpdatedAt did not change after completion")
	}
	
	// StartedAt should remain the same
	if task.StartedAt != originalStartedAt {
		t.Errorf("Task.StartedAt changed after completion, should remain constant")
	}
	
	// CompletedAt should be set
	if task.CompletedAt == nil {
		t.Error("Task.CompletedAt is nil, expected timestamp")
	}
}

func TestTaskFail(t *testing.T) {
	// Create a test task
	task, err := NewTask("test_task", []byte(`{"test":"data"}`))
	if err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}
	
	// Mark as processing first
	task.MarkAsProcessing()
	
	// Store timestamps for comparison
	originalStartedAt := task.StartedAt
	originalUpdatedAt := task.UpdatedAt
	
	// Wait a moment to ensure updated timestamp will be different
	time.Sleep(1 * time.Millisecond)
	
	// Fail the task with an error message
	errorMsg := "Something went wrong"
	task.Fail(errorMsg)
	
	// Check status was updated
	if task.Status != TaskStatusFailed {
		t.Errorf("Task.Status = %v, want %v", task.Status, TaskStatusFailed)
	}
	
	// Check error message was set
	if task.Error != errorMsg {
		t.Errorf("Task.Error = %v, want %v", task.Error, errorMsg)
	}
	
	// UpdatedAt should be updated
	if task.UpdatedAt == originalUpdatedAt {
		t.Errorf("Task.UpdatedAt did not change after failure")
	}
	
	// StartedAt should remain the same
	if task.StartedAt != originalStartedAt {
		t.Errorf("Task.StartedAt changed after failure, should remain constant")
	}
	
	// CompletedAt should be set
	if task.CompletedAt == nil {
		t.Error("Task.CompletedAt is nil, expected timestamp")
	}
}

func TestTaskIsActive(t *testing.T) {
	// Create a test task
	task, err := NewTask("test_task", []byte(`{"test":"data"}`))
	if err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}
	
	// New task should be active (pending)
	if !task.IsActive() {
		t.Error("New task (pending) should be active")
	}
	
	// Processing task should be active
	task.MarkAsProcessing()
	if !task.IsActive() {
		t.Error("Processing task should be active")
	}
	
	// Completed task should not be active
	task, err = NewTask("test_task", []byte(`{"test":"data"}`))
	if err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}
	task.Complete()
	if task.IsActive() {
		t.Error("Completed task should not be active")
	}
	
	// Failed task should not be active
	task, err = NewTask("test_task", []byte(`{"test":"data"}`))
	if err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}
	task.Fail("Error message")
	if task.IsActive() {
		t.Error("Failed task should not be active")
	}
}

func TestTaskDuration(t *testing.T) {
	// Create a test task
	task, err := NewTask("test_task", []byte(`{"test":"data"}`))
	if err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}
	
	// New task should have zero duration
	if task.Duration() != 0 {
		t.Errorf("New task duration = %v, want 0", task.Duration())
	}
	
	// Start the task
	task.MarkAsProcessing()
	
	// Wait a bit to accumulate some duration
	time.Sleep(10 * time.Millisecond)
	
	// Duration should be non-zero for a started task
	if task.Duration() <= 0 {
		t.Error("Started task should have non-zero duration")
	}
	
	// Complete the task and check final duration
	time.Sleep(10 * time.Millisecond)
	task.Complete()
	
	finalDuration := task.Duration()
	if finalDuration <= 0 {
		t.Error("Completed task should have positive duration")
	}
	
	// Duration should not change after completion
	time.Sleep(10 * time.Millisecond)
	if task.Duration() != finalDuration {
		t.Error("Task duration should not change after completion")
	}
}
