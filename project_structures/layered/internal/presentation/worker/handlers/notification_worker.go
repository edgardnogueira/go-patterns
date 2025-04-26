package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/application/service"
	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/infrastructure/config"
	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/infrastructure/logger"
	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/infrastructure/providers/external"
)

// NotificationType represents the type of notification task
type NotificationType string

const (
	// NotifyNewPost is a notification for a new post
	NotifyNewPost NotificationType = "new_post"
	
	// NotifyComment is a notification for a new comment
	NotifyComment NotificationType = "comment"
	
	// NotifySubscription is a notification for a new subscription
	NotifySubscription NotificationType = "subscription"
	
	// SendEmailNotification is a task to send email notifications
	SendEmailNotification NotificationType = "send_email"
)

// NotificationTask represents a task to be processed by the notification worker
type NotificationTask struct {
	Type       NotificationType     `json:"type"`
	Data       map[string]interface{} `json:"data"`
	CreatedAt  time.Time            `json:"created_at"`
	RetryCount int                  `json:"retry_count"`
}

// NotificationWorker handles background notification processing
type NotificationWorker struct {
	notificationService *service.NotificationAppService
	emailProvider       *external.EmailProvider
	logger              *logger.Logger
	config              *config.WorkerConfig
	taskQueue           chan NotificationTask
	shutdown            chan struct{}
}

// NewNotificationWorker creates a new notification worker
func NewNotificationWorker(
	notificationService *service.NotificationAppService,
	emailProvider *external.EmailProvider,
	log *logger.Logger,
	cfg *config.WorkerConfig,
) *NotificationWorker {
	return &NotificationWorker{
		notificationService: notificationService,
		emailProvider:       emailProvider,
		logger:              log,
		config:              cfg,
		taskQueue:           make(chan NotificationTask, cfg.QueueSize),
		shutdown:            make(chan struct{}),
	}
}

// Start starts the notification worker
func (w *NotificationWorker) Start() {
	w.logger.Info("Starting notification worker")
	
	// Start worker goroutines
	for i := 0; i < w.config.Concurrency; i++ {
		go w.processTasksWorker(i)
	}
}

// Stop stops the notification worker
func (w *NotificationWorker) Stop() {
	w.logger.Info("Stopping notification worker")
	close(w.shutdown)
}

// EnqueueTask adds a task to the queue
func (w *NotificationWorker) EnqueueTask(task NotificationTask) error {
	select {
	case w.taskQueue <- task:
		w.logger.WithFields(map[string]interface{}{
			"type": task.Type,
			"data": task.Data,
		}).Debug("Task enqueued")
		return nil
	default:
		return fmt.Errorf("task queue is full")
	}
}

// processTasksWorker processes tasks from the queue
func (w *NotificationWorker) processTasksWorker(workerID int) {
	w.logger.WithField("workerID", workerID).Info("Worker started")
	
	for {
		select {
		case task := <-w.taskQueue:
			w.processTask(task, workerID)
		case <-w.shutdown:
			w.logger.WithField("workerID", workerID).Info("Worker shutting down")
			return
		}
	}
}

// processTask processes a notification task
func (w *NotificationWorker) processTask(task NotificationTask, workerID int) {
	ctx := context.Background()
	logger := w.logger.WithFields(map[string]interface{}{
		"type":      task.Type,
		"workerID":  workerID,
		"createdAt": task.CreatedAt.Format(time.RFC3339),
	})
	
	logger.Info("Processing task")
	
	var err error
	
	switch task.Type {
	case NotifyNewPost:
		err = w.handleNewPostNotification(ctx, task)
	case NotifyComment:
		err = w.handleCommentNotification(ctx, task)
	case NotifySubscription:
		err = w.handleSubscriptionNotification(ctx, task)
	case SendEmailNotification:
		err = w.handleEmailNotification(task)
	default:
		logger.WithField("type", task.Type).Warn("Unknown task type")
		return
	}
	
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to process task")
		
		// Retry logic
		if task.RetryCount < 3 {
			task.RetryCount++
			time.Sleep(time.Duration(task.RetryCount) * time.Second)
			w.EnqueueTask(task)
			logger.WithField("retryCount", task.RetryCount).Info("Task re-enqueued for retry")
		} else {
			logger.WithField("retryCount", task.RetryCount).Error("Max retries reached, task failed")
		}
	} else {
		logger.Info("Task processed successfully")
	}
}

// handleNewPostNotification handles notifications for new posts
func (w *NotificationWorker) handleNewPostNotification(ctx context.Context, task NotificationTask) error {
	// Extract task data
	postID, ok := task.Data["postID"].(float64)
	if !ok {
		return fmt.Errorf("invalid postID")
	}
	
	subscribersData, ok := task.Data["subscribers"].([]interface{})
	if !ok {
		return fmt.Errorf("invalid subscribers data")
	}
	
	// Convert subscribers to int64 slice
	subscribers := make([]int64, len(subscribersData))
	for i, v := range subscribersData {
		if userID, ok := v.(float64); ok {
			subscribers[i] = int64(userID)
		}
	}
	
	// Notify subscribers
	return w.notificationService.NotifyNewPost(ctx, int64(postID), subscribers)
}

// handleCommentNotification handles notifications for new comments
func (w *NotificationWorker) handleCommentNotification(ctx context.Context, task NotificationTask) error {
	// Extract task data
	postID, ok := task.Data["postID"].(float64)
	if !ok {
		return fmt.Errorf("invalid postID")
	}
	
	commenterName, ok := task.Data["commenterName"].(string)
	if !ok {
		return fmt.Errorf("invalid commenterName")
	}
	
	// Notify post author
	return w.notificationService.NotifyComment(ctx, int64(postID), commenterName)
}

// handleSubscriptionNotification handles notifications for new subscriptions
func (w *NotificationWorker) handleSubscriptionNotification(ctx context.Context, task NotificationTask) error {
	// Extract task data
	userID, ok := task.Data["userID"].(float64)
	if !ok {
		return fmt.Errorf("invalid userID")
	}
	
	authorID, ok := task.Data["authorID"].(float64)
	if !ok {
		return fmt.Errorf("invalid authorID")
	}
	
	// Notify user about their subscription
	return w.notificationService.NotifySubscription(ctx, int64(userID), int64(authorID))
}

// handleEmailNotification handles sending email notifications
func (w *NotificationWorker) handleEmailNotification(task NotificationTask) error {
	// Extract task data
	to, ok := task.Data["to"].(string)
	if !ok {
		return fmt.Errorf("invalid to field")
	}
	
	subject, ok := task.Data["subject"].(string)
	if !ok {
		return fmt.Errorf("invalid subject")
	}
	
	body, ok := task.Data["body"].(string)
	if !ok {
		return fmt.Errorf("invalid body")
	}
	
	// Send email
	return w.emailProvider.SendEmail(to, subject, body)
}

// SerializeTask serializes a notification task to JSON
func SerializeTask(task NotificationTask) (string, error) {
	data, err := json.Marshal(task)
	if err != nil {
		return "", fmt.Errorf("failed to serialize task: %w", err)
	}
	return string(data), nil
}

// DeserializeTask deserializes a notification task from JSON
func DeserializeTask(data string) (NotificationTask, error) {
	var task NotificationTask
	err := json.Unmarshal([]byte(data), &task)
	if err != nil {
		return task, fmt.Errorf("failed to deserialize task: %w", err)
	}
	return task, nil
}
