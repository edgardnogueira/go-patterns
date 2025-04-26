package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	
	"github.com/edgardnogueira/go-patterns/project_structures/hexagonal/internal/domain/model"
	"github.com/edgardnogueira/go-patterns/project_structures/hexagonal/internal/ports/driven"
)

// OrderService implements the business logic for order operations
type OrderService struct {
	orderRepo         driven.OrderRepository
	notificationService driven.NotificationService
}

// NewOrderService creates a new OrderService with required dependencies
func NewOrderService(
	orderRepo driven.OrderRepository,
	notificationService driven.NotificationService,
) *OrderService {
	return &OrderService{
		orderRepo:           orderRepo,
		notificationService: notificationService,
	}
}

// CreateOrder creates a new order and saves it to the repository
func (s *OrderService) CreateOrder(
	ctx context.Context,
	customerID string,
	items []model.OrderItem,
	shippingAddress string,
) (*model.Order, error) {
	// Generate a new order ID
	orderID := uuid.New().String()
	
	// Create the order
	order, err := model.NewOrder(orderID, customerID, items, shippingAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	
	// Save to repository
	if err := s.orderRepo.Save(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}
	
	// Send notification
	if err := s.notificationService.NotifyOrderStatus(
		ctx, 
		order, 
		driven.NotificationTypeOrderCreated,
	); err != nil {
		// Log the error but don't fail the operation
		fmt.Printf("Failed to send notification: %v\n", err)
	}
	
	return order, nil
}

// GetOrder retrieves an order by ID
func (s *OrderService) GetOrder(ctx context.Context, orderID string) (*model.Order, error) {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	
	return order, nil
}

// UpdateOrderStatus updates the status of an order
func (s *OrderService) UpdateOrderStatus(
	ctx context.Context,
	orderID string,
	newStatus model.OrderStatus,
) error {
	// Retrieve the order
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}
	
	// Update the status
	if err := order.UpdateStatus(newStatus); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	
	// Save the changes
	if err := s.orderRepo.Save(ctx, order); err != nil {
		return fmt.Errorf("failed to save order: %w", err)
	}
	
	// Map order status to notification type
	var notificationType driven.NotificationType
	switch newStatus {
	case model.OrderStatusProcessing:
		notificationType = driven.NotificationTypeOrderProcessing
	case model.OrderStatusShipped:
		notificationType = driven.NotificationTypeOrderShipped
	case model.OrderStatusDelivered:
		notificationType = driven.NotificationTypeOrderDelivered
	case model.OrderStatusCancelled:
		notificationType = driven.NotificationTypeOrderCancelled
	default:
		notificationType = driven.NotificationTypeOrderCreated
	}
	
	// Send notification
	if err := s.notificationService.NotifyOrderStatus(ctx, order, notificationType); err != nil {
		// Log the error but don't fail the operation
		fmt.Printf("Failed to send notification: %v\n", err)
	}
	
	return nil
}

// ListOrders retrieves all orders for a specific customer
func (s *OrderService) ListOrders(ctx context.Context, customerID string) ([]*model.Order, error) {
	orders, err := s.orderRepo.FindByCustomer(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}
	
	return orders, nil
}

// AddItemToOrder adds an item to an existing order
func (s *OrderService) AddItemToOrder(
	ctx context.Context,
	orderID string,
	item model.OrderItem,
) error {
	// Retrieve the order
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}
	
	// Add the item
	if err := order.AddItem(item); err != nil {
		return fmt.Errorf("failed to add item to order: %w", err)
	}
	
	// Save the changes
	if err := s.orderRepo.Save(ctx, order); err != nil {
		return fmt.Errorf("failed to save order: %w", err)
	}
	
	return nil
}

// RemoveItemFromOrder removes an item from an existing order
func (s *OrderService) RemoveItemFromOrder(
	ctx context.Context,
	orderID string,
	productID string,
) error {
	// Retrieve the order
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}
	
	// Remove the item
	if err := order.RemoveItem(productID); err != nil {
		return fmt.Errorf("failed to remove item from order: %w", err)
	}
	
	// Save the changes
	if err := s.orderRepo.Save(ctx, order); err != nil {
		return fmt.Errorf("failed to save order: %w", err)
	}
	
	return nil
}

// ProcessOrder processes an order by updating its status to "processing"
func (s *OrderService) ProcessOrder(ctx context.Context, orderID string) error {
	return s.UpdateOrderStatus(ctx, orderID, model.OrderStatusProcessing)
}

// FindOrdersByStatus returns all orders with a specific status
func (s *OrderService) FindOrdersByStatus(ctx context.Context, status model.OrderStatus) ([]*model.Order, error) {
	return s.orderRepo.FindByStatus(ctx, status)
}
