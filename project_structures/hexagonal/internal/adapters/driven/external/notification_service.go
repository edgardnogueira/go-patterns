package external

import (
	"context"
	"fmt"
	"log"
	"time"
	
	"github.com/edgardnogueira/go-patterns/project_structures/hexagonal/internal/domain/model"
	"github.com/edgardnogueira/go-patterns/project_structures/hexagonal/internal/ports/driven"
)

// LogNotificationService is a simple implementation of the NotificationService
// that logs notifications instead of sending them to external systems
type LogNotificationService struct {
	logger *log.Logger
}

// NewLogNotificationService creates a new notification service that logs messages
func NewLogNotificationService(logger *log.Logger) *LogNotificationService {
	if logger == nil {
		logger = log.Default()
	}
	
	return &LogNotificationService{
		logger: logger,
	}
}

// NotifyOrderStatus logs a notification about an order status change
func (s *LogNotificationService) NotifyOrderStatus(
	ctx context.Context,
	order *model.Order,
	notificationType driven.NotificationType,
) error {
	if order == nil {
		return fmt.Errorf("cannot notify about nil order")
	}
	
	message := fmt.Sprintf(
		"[%s] Notification: Order %s (customer: %s) status changed to %s",
		time.Now().Format(time.RFC3339),
		order.ID,
		order.CustomerID,
		order.Status,
	)
	
	s.logger.Println(message)
	
	return nil
}

// NotifyCustomer logs a notification to a customer
func (s *LogNotificationService) NotifyCustomer(
	ctx context.Context,
	customerID string,
	message string,
	metadata map[string]string,
) error {
	logMessage := fmt.Sprintf(
		"[%s] Customer Notification: To %s - %s",
		time.Now().Format(time.RFC3339),
		customerID,
		message,
	)
	
	if len(metadata) > 0 {
		logMessage += fmt.Sprintf(" - Metadata: %v", metadata)
	}
	
	s.logger.Println(logMessage)
	
	return nil
}
