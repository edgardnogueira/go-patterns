package driven

import (
	"context"
	
	"github.com/edgardnogueira/go-patterns/project_structures/hexagonal/internal/domain/model"
)

// NotificationType represents different types of notifications
type NotificationType string

const (
	NotificationTypeOrderCreated     NotificationType = "order_created"
	NotificationTypeOrderProcessing  NotificationType = "order_processing"
	NotificationTypeOrderShipped     NotificationType = "order_shipped"
	NotificationTypeOrderDelivered   NotificationType = "order_delivered"
	NotificationTypeOrderCancelled   NotificationType = "order_cancelled"
)

// NotificationService is a secondary port (driven port) that defines
// how the domain interacts with external notification systems
type NotificationService interface {
	// NotifyOrderStatus sends a notification about an order status change
	NotifyOrderStatus(ctx context.Context, order *model.Order, notificationType NotificationType) error
	
	// NotifyCustomer sends a notification to a customer
	NotifyCustomer(ctx context.Context, customerID string, message string, metadata map[string]string) error
}
