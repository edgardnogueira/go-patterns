package mapper

import (
	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/data/entity"
	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/domain/service"
)

// NotificationMapper converts between domain Notification models and NotificationEntity data objects
type NotificationMapper struct{}

// NewNotificationMapper creates a new instance of NotificationMapper
func NewNotificationMapper() *NotificationMapper {
	return &NotificationMapper{}
}

// ToEntity converts a domain Notification model to a NotificationEntity
func (m *NotificationMapper) ToEntity(notification *service.Notification) *entity.NotificationEntity {
	return &entity.NotificationEntity{
		ID:         notification.ID,
		UserID:     notification.UserID,
		Type:       string(notification.Type),
		Title:      notification.Title,
		Content:    notification.Content,
		ResourceID: notification.ResourceID,
		Read:       notification.Read,
		CreatedAt:  notification.CreatedAt,
	}
}

// ToDomain converts a NotificationEntity to a domain Notification model
func (m *NotificationMapper) ToDomain(entity *entity.NotificationEntity) *service.Notification {
	return &service.Notification{
		ID:         entity.ID,
		UserID:     entity.UserID,
		Type:       service.NotificationType(entity.Type),
		Title:      entity.Title,
		Content:    entity.Content,
		ResourceID: entity.ResourceID,
		Read:       entity.Read,
		CreatedAt:  entity.CreatedAt,
	}
}

// ToDomainList converts a slice of NotificationEntity to a slice of domain Notification models
func (m *NotificationMapper) ToDomainList(entities []*entity.NotificationEntity) []*service.Notification {
	notifications := make([]*service.Notification, len(entities))
	for i, entity := range entities {
		notifications[i] = m.ToDomain(entity)
	}
	return notifications
}
