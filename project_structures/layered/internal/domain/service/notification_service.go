package service

import (
	"context"
	"time"
)

// NotificationType defines the type of notification
type NotificationType string

const (
	// NotificationTypeNewPost indicates a new post notification
	NotificationTypeNewPost NotificationType = "new_post"
	
	// NotificationTypeComment indicates a comment notification
	NotificationTypeComment NotificationType = "comment"
	
	// NotificationTypeSubscription indicates a subscription notification
	NotificationTypeSubscription NotificationType = "subscription"
)

// Notification represents a notification object in the domain layer
type Notification struct {
	ID          int64
	UserID      int64
	Type        NotificationType
	Title       string
	Content     string
	ResourceID  int64  // ID of the related entity (e.g., post ID)
	Read        bool
	CreatedAt   time.Time
}

// NotificationRepository defines the interface for notification data access
type NotificationRepository interface {
	FindByID(ctx context.Context, id int64) (*Notification, error)
	FindByUserID(ctx context.Context, userID int64) ([]*Notification, error)
	FindUnreadByUserID(ctx context.Context, userID int64) ([]*Notification, error)
	Save(ctx context.Context, notification *Notification) error
	Update(ctx context.Context, notification *Notification) error
	Delete(ctx context.Context, id int64) error
	MarkAsRead(ctx context.Context, id int64) error
	MarkAllAsRead(ctx context.Context, userID int64) error
}

// NotificationService defines the business operations for notifications
type NotificationService interface {
	GetNotification(ctx context.Context, id int64) (*Notification, error)
	GetUserNotifications(ctx context.Context, userID int64) ([]*Notification, error)
	GetUnreadNotifications(ctx context.Context, userID int64) ([]*Notification, error)
	CreateNotification(ctx context.Context, notification *Notification) error
	MarkAsRead(ctx context.Context, id int64) error
	MarkAllAsRead(ctx context.Context, userID int64) error
	DeleteNotification(ctx context.Context, id int64) error
	
	// Create specific notification types
	NotifyNewPost(ctx context.Context, userID, postID int64, postTitle string) error
	NotifyComment(ctx context.Context, userID, postID int64, commenterName string) error
	NotifySubscription(ctx context.Context, userID, authorID int64, authorName string) error
}

// domainNotificationService implements NotificationService interface
type domainNotificationService struct {
	notificationRepo NotificationRepository
}

// NewNotificationService creates a new instance of the notification service
func NewNotificationService(repo NotificationRepository) NotificationService {
	return &domainNotificationService{
		notificationRepo: repo,
	}
}

// GetNotification retrieves a notification by ID
func (s *domainNotificationService) GetNotification(ctx context.Context, id int64) (*Notification, error) {
	return s.notificationRepo.FindByID(ctx, id)
}

// GetUserNotifications retrieves all notifications for a user
func (s *domainNotificationService) GetUserNotifications(ctx context.Context, userID int64) ([]*Notification, error) {
	return s.notificationRepo.FindByUserID(ctx, userID)
}

// GetUnreadNotifications retrieves all unread notifications for a user
func (s *domainNotificationService) GetUnreadNotifications(ctx context.Context, userID int64) ([]*Notification, error) {
	return s.notificationRepo.FindUnreadByUserID(ctx, userID)
}

// CreateNotification creates a new notification
func (s *domainNotificationService) CreateNotification(ctx context.Context, notification *Notification) error {
	notification.CreatedAt = time.Now()
	notification.Read = false
	return s.notificationRepo.Save(ctx, notification)
}

// MarkAsRead marks a notification as read
func (s *domainNotificationService) MarkAsRead(ctx context.Context, id int64) error {
	return s.notificationRepo.MarkAsRead(ctx, id)
}

// MarkAllAsRead marks all notifications as read for a user
func (s *domainNotificationService) MarkAllAsRead(ctx context.Context, userID int64) error {
	return s.notificationRepo.MarkAllAsRead(ctx, userID)
}

// DeleteNotification deletes a notification
func (s *domainNotificationService) DeleteNotification(ctx context.Context, id int64) error {
	return s.notificationRepo.Delete(ctx, id)
}

// NotifyNewPost creates a notification for a new post
func (s *domainNotificationService) NotifyNewPost(ctx context.Context, userID, postID int64, postTitle string) error {
	notification := &Notification{
		UserID:     userID,
		Type:       NotificationTypeNewPost,
		Title:      "New Post Available",
		Content:    "A new post has been published: " + postTitle,
		ResourceID: postID,
		CreatedAt:  time.Now(),
		Read:       false,
	}
	
	return s.notificationRepo.Save(ctx, notification)
}

// NotifyComment creates a notification for a new comment
func (s *domainNotificationService) NotifyComment(ctx context.Context, userID, postID int64, commenterName string) error {
	notification := &Notification{
		UserID:     userID,
		Type:       NotificationTypeComment,
		Title:      "New Comment",
		Content:    commenterName + " commented on your post.",
		ResourceID: postID,
		CreatedAt:  time.Now(),
		Read:       false,
	}
	
	return s.notificationRepo.Save(ctx, notification)
}

// NotifySubscription creates a notification for a new subscription
func (s *domainNotificationService) NotifySubscription(ctx context.Context, userID, authorID int64, authorName string) error {
	notification := &Notification{
		UserID:     userID,
		Type:       NotificationTypeSubscription,
		Title:      "New Subscription",
		Content:    "You are now subscribed to " + authorName + "'s posts.",
		ResourceID: authorID,
		CreatedAt:  time.Now(),
		Read:       false,
	}
	
	return s.notificationRepo.Save(ctx, notification)
}
