package dto

import (
	"time"

	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/domain/service"
)

// NotificationResponse represents the response body for a notification
type NotificationResponse struct {
	ID         int64     `json:"id"`
	Type       string    `json:"type"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	ResourceID int64     `json:"resourceId"`
	Read       bool      `json:"read"`
	CreatedAt  time.Time `json:"createdAt"`
}

// NotificationListResponse represents a list of notifications with count
type NotificationListResponse struct {
	Notifications []*NotificationResponse `json:"notifications"`
	Count         int                     `json:"count"`
	Unread        int                     `json:"unread"`
}

// CreateNotificationRequest represents the request body for creating a notification
type CreateNotificationRequest struct {
	UserID     int64  `json:"userId" validate:"required,gt=0"`
	Type       string `json:"type" validate:"required"`
	Title      string `json:"title" validate:"required"`
	Content    string `json:"content" validate:"required"`
	ResourceID int64  `json:"resourceId" validate:"required,gt=0"`
}

// ToResponse converts a domain Notification model to a NotificationResponse DTO
func ToNotificationResponse(notification *service.Notification) *NotificationResponse {
	return &NotificationResponse{
		ID:         notification.ID,
		Type:       string(notification.Type),
		Title:      notification.Title,
		Content:    notification.Content,
		ResourceID: notification.ResourceID,
		Read:       notification.Read,
		CreatedAt:  notification.CreatedAt,
	}
}

// ToNotificationResponses converts a slice of domain Notification models to NotificationResponse DTOs
func ToNotificationResponses(notifications []*service.Notification) []*NotificationResponse {
	responses := make([]*NotificationResponse, len(notifications))
	for i, notification := range notifications {
		responses[i] = ToNotificationResponse(notification)
	}
	return responses
}

// ToDomain converts a CreateNotificationRequest to a domain Notification model
func (r *CreateNotificationRequest) ToDomain() *service.Notification {
	return &service.Notification{
		UserID:     r.UserID,
		Type:       service.NotificationType(r.Type),
		Title:      r.Title,
		Content:    r.Content,
		ResourceID: r.ResourceID,
		Read:       false,
		CreatedAt:  time.Now(),
	}
}
