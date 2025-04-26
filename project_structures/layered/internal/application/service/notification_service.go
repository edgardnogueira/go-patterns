package service

import (
	"context"
	"errors"
	"fmt"

	domainservice "github.com/edgardnogueira/go-patterns/project_structures/layered/internal/domain/service"
)

// Application-level errors for notification service
var (
	ErrNotificationNotFound = errors.New("notification not found")
	ErrInvalidNotification  = errors.New("invalid notification data")
)

// NotificationAppService orchestrates notification-related use cases
type NotificationAppService struct {
	notificationService domainservice.NotificationService
	authorService       domainservice.AuthorService
	postService         domainservice.PostService
}

// NewNotificationAppService creates a new instance of NotificationAppService
func NewNotificationAppService(
	notificationService domainservice.NotificationService,
	authorService domainservice.AuthorService,
	postService domainservice.PostService,
) *NotificationAppService {
	return &NotificationAppService{
		notificationService: notificationService,
		authorService:       authorService,
		postService:         postService,
	}
}

// GetNotification retrieves a notification by ID
func (s *NotificationAppService) GetNotification(ctx context.Context, id int64) (*domainservice.Notification, error) {
	notification, err := s.notificationService.GetNotification(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}
	
	return notification, nil
}

// GetUserNotifications retrieves all notifications for a user
func (s *NotificationAppService) GetUserNotifications(ctx context.Context, userID int64) ([]*domainservice.Notification, error) {
	notifications, err := s.notificationService.GetUserNotifications(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user notifications: %w", err)
	}
	
	return notifications, nil
}

// GetUnreadNotifications retrieves all unread notifications for a user
func (s *NotificationAppService) GetUnreadNotifications(ctx context.Context, userID int64) ([]*domainservice.Notification, error) {
	notifications, err := s.notificationService.GetUnreadNotifications(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get unread notifications: %w", err)
	}
	
	return notifications, nil
}

// CreateNotification creates a new notification
func (s *NotificationAppService) CreateNotification(ctx context.Context, notification *domainservice.Notification) error {
	if err := s.notificationService.CreateNotification(ctx, notification); err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}
	
	return nil
}

// MarkAsRead marks a notification as read
func (s *NotificationAppService) MarkAsRead(ctx context.Context, id int64) error {
	// Check if notification exists
	_, err := s.notificationService.GetNotification(ctx, id)
	if err != nil {
		return ErrNotificationNotFound
	}
	
	if err := s.notificationService.MarkAsRead(ctx, id); err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}
	
	return nil
}

// MarkAllAsRead marks all notifications as read for a user
func (s *NotificationAppService) MarkAllAsRead(ctx context.Context, userID int64) error {
	if err := s.notificationService.MarkAllAsRead(ctx, userID); err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}
	
	return nil
}

// DeleteNotification deletes a notification
func (s *NotificationAppService) DeleteNotification(ctx context.Context, id int64) error {
	// Check if notification exists
	_, err := s.notificationService.GetNotification(ctx, id)
	if err != nil {
		return ErrNotificationNotFound
	}
	
	if err := s.notificationService.DeleteNotification(ctx, id); err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	
	return nil
}

// NotifyNewPost creates a notification for a new post for all subscribers
func (s *NotificationAppService) NotifyNewPost(ctx context.Context, postID int64, subscribers []int64) error {
	// Get post details
	post, err := s.postService.GetPostByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("failed to get post details: %w", err)
	}
	
	// Notify each subscriber
	for _, userID := range subscribers {
		err := s.notificationService.NotifyNewPost(ctx, userID, postID, post.Title)
		if err != nil {
			// Continue notifying other users even if one fails
			continue
		}
	}
	
	return nil
}

// NotifyComment creates a notification for a new comment
func (s *NotificationAppService) NotifyComment(ctx context.Context, postID int64, commenterName string) error {
	// Get post details to find its author
	post, err := s.postService.GetPostByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("failed to get post details: %w", err)
	}
	
	// Notify post author about the comment
	if err := s.notificationService.NotifyComment(ctx, post.AuthorID, postID, commenterName); err != nil {
		return fmt.Errorf("failed to notify about comment: %w", err)
	}
	
	return nil
}

// NotifySubscription creates a notification for a new subscription
func (s *NotificationAppService) NotifySubscription(ctx context.Context, userID, authorID int64) error {
	// Get author details
	author, err := s.authorService.GetAuthorByID(ctx, authorID)
	if err != nil {
		return fmt.Errorf("failed to get author details: %w", err)
	}
	
	// Notify user about their subscription
	if err := s.notificationService.NotifySubscription(ctx, userID, authorID, author.Name); err != nil {
		return fmt.Errorf("failed to notify about subscription: %w", err)
	}
	
	return nil
}

// GetNotificationsWithCount retrieves notifications with count info
func (s *NotificationAppService) GetNotificationsWithCount(ctx context.Context, userID int64) ([]*domainservice.Notification, int, int, error) {
	// Get all notifications
	allNotifications, err := s.notificationService.GetUserNotifications(ctx, userID)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get user notifications: %w", err)
	}
	
	// Get unread notifications
	unreadNotifications, err := s.notificationService.GetUnreadNotifications(ctx, userID)
	if err != nil {
		// If we can't get unread count, still return all notifications
		return allNotifications, len(allNotifications), 0, nil
	}
	
	return allNotifications, len(allNotifications), len(unreadNotifications), nil
}
