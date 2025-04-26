package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/edgardnogueira/go-patterns/project_structures/clean_architecture/internal/domain/entities"
	"github.com/edgardnogueira/go-patterns/project_structures/clean_architecture/internal/infrastructure/providers"
	"github.com/edgardnogueira/go-patterns/project_structures/clean_architecture/internal/usecase"
)

// TaskProcessor processes tasks asynchronously
type TaskProcessor struct {
	taskService  *usecase.TaskService
	emailProvider providers.EmailProvider
	pollInterval time.Duration
	running      bool
	stopCh       chan struct{}
}

// NewTaskProcessor creates a new TaskProcessor
func NewTaskProcessor(
	taskService *usecase.TaskService,
	emailProvider providers.EmailProvider,
	pollInterval time.Duration,
) *TaskProcessor {
	if pollInterval == 0 {
		pollInterval = 5 * time.Second
	}

	return &TaskProcessor{
		taskService:  taskService,
		emailProvider: emailProvider,
		pollInterval: pollInterval,
		stopCh:       make(chan struct{}),
	}
}

// Start begins the task processing loop
func (p *TaskProcessor) Start(ctx context.Context) {
	if p.running {
		return
	}
	p.running = true

	go func() {
		log.Println("Task processor started")
		ticker := time.NewTicker(p.pollInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				p.processNextTask(ctx)
			case <-p.stopCh:
				log.Println("Task processor stopped")
				return
			case <-ctx.Done():
				log.Println("Task processor stopped due to context cancellation")
				return
			}
		}
	}()
}

// Stop stops the task processing loop
func (p *TaskProcessor) Stop() {
	if !p.running {
		return
	}
	p.running = false
	p.stopCh <- struct{}{}
}

// processNextTask attempts to process the next pending task
func (p *TaskProcessor) processNextTask(ctx context.Context) {
	task, err := p.taskService.GetNextPendingTask(ctx)
	if err != nil {
		if err != usecase.ErrNoTaskAvailable {
			log.Printf("Error fetching next task: %v\n", err)
		}
		return
	}

	log.Printf("Processing task ID: %s, Type: %s\n", task.ID, task.Type)

	// Process task based on its type
	var processingErr error
	switch task.Type {
	case usecase.UserTaskTypeWelcomeEmail:
		processingErr = p.processWelcomeEmailTask(ctx, task)
	case usecase.UserTaskTypeProfileUpdate:
		processingErr = p.processProfileUpdateTask(ctx, task)
	default:
		processingErr = fmt.Errorf("unknown task type: %s", task.Type)
	}

	// Update task status based on processing result
	if processingErr != nil {
		log.Printf("Task processing error: %v\n", processingErr)
		_, err = p.taskService.FailTask(ctx, task.ID, processingErr.Error())
	} else {
		_, err = p.taskService.CompleteTask(ctx, task.ID)
	}

	if err != nil {
		log.Printf("Error updating task status: %v\n", err)
	}
}

// processWelcomeEmailTask handles sending welcome emails
func (p *TaskProcessor) processWelcomeEmailTask(ctx context.Context, task *entities.Task) error {
	var data struct {
		UserID    string `json:"user_id"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	if err := json.Unmarshal(task.Data, &data); err != nil {
		return fmt.Errorf("failed to unmarshal task data: %w", err)
	}

	name := data.FirstName
	if name == "" {
		name = data.LastName
	}
	if name == "" {
		name = "there"
	}

	// Format the email content
	subject := "Welcome to Our Platform!"
	body := fmt.Sprintf(providers.WelcomeEmailTemplate, name)

	// Send the email using the provider
	if err := p.emailProvider.SendEmail(ctx, data.Email, subject, body); err != nil {
		return fmt.Errorf("failed to send welcome email: %w", err)
	}

	return nil
}

// processProfileUpdateTask handles profile update notifications
func (p *TaskProcessor) processProfileUpdateTask(ctx context.Context, task *entities.Task) error {
	var data struct {
		UserID    string `json:"user_id"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	if err := json.Unmarshal(task.Data, &data); err != nil {
		return fmt.Errorf("failed to unmarshal task data: %w", err)
	}

	name := data.FirstName
	if name == "" {
		name = data.LastName
	}
	if name == "" {
		name = data.Username
	}
	if name == "" {
		name = "there"
	}

	// Format the email content
	subject := "Your Profile Has Been Updated"
	body := fmt.Sprintf(providers.ProfileUpdateEmailTemplate, name)

	// Send the email using the provider
	if err := p.emailProvider.SendEmail(ctx, data.Email, subject, body); err != nil {
		return fmt.Errorf("failed to send profile update email: %w", err)
	}

	return nil
}
