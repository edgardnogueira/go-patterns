package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/edgardnogueira/go-patterns/project_structures/hexagonal/internal/domain/model"
	"github.com/edgardnogueira/go-patterns/project_structures/hexagonal/internal/ports/driven"
)

// mockOrderRepository is a mock implementation of the OrderRepository interface for testing
type mockOrderRepository struct {
	orders map[string]*model.Order
}

func newMockOrderRepository() *mockOrderRepository {
	return &mockOrderRepository{
		orders: make(map[string]*model.Order),
	}
}

func (m *mockOrderRepository) Save(ctx context.Context, order *model.Order) error {
	if order == nil {
		return errors.New("cannot save nil order")
	}
	m.orders[order.ID] = order
	return nil
}

func (m *mockOrderRepository) FindByID(ctx context.Context, orderID string) (*model.Order, error) {
	order, exists := m.orders[orderID]
	if !exists {
		return nil, errors.New("order not found")
	}
	return order, nil
}

func (m *mockOrderRepository) FindByCustomer(ctx context.Context, customerID string) ([]*model.Order, error) {
	var result []*model.Order
	for _, order := range m.orders {
		if order.CustomerID == customerID {
			result = append(result, order)
		}
	}
	return result, nil
}

func (m *mockOrderRepository) FindByStatus(ctx context.Context, status model.OrderStatus) ([]*model.Order, error) {
	var result []*model.Order
	for _, order := range m.orders {
		if order.Status == status {
			result = append(result, order)
		}
	}
	return result, nil
}

func (m *mockOrderRepository) Delete(ctx context.Context, orderID string) error {
	if _, exists := m.orders[orderID]; !exists {
		return errors.New("order not found")
	}
	delete(m.orders, orderID)
	return nil
}

// mockNotificationService is a mock implementation of the NotificationService interface for testing
type mockNotificationService struct {
	notifications []struct {
		order           *model.Order
		notificationType driven.NotificationType
	}
	customerNotifications []struct {
		customerID string
		message    string
		metadata   map[string]string
	}
}

func newMockNotificationService() *mockNotificationService {
	return &mockNotificationService{}
}

func (m *mockNotificationService) NotifyOrderStatus(ctx context.Context, order *model.Order, notificationType driven.NotificationType) error {
	if order == nil {
		return errors.New("cannot notify about nil order")
	}
	m.notifications = append(m.notifications, struct {
		order           *model.Order
		notificationType driven.NotificationType
	}{order, notificationType})
	return nil
}

func (m *mockNotificationService) NotifyCustomer(ctx context.Context, customerID string, message string, metadata map[string]string) error {
	m.customerNotifications = append(m.customerNotifications, struct {
		customerID string
		message    string
		metadata   map[string]string
	}{customerID, message, metadata})
	return nil
}

func TestOrderService_CreateOrder(t *testing.T) {
	// Create context for all tests
	ctx := context.Background()
	
	// Create mocks
	mockRepo := newMockOrderRepository()
	mockNotifier := newMockNotificationService()
	
	// Create the service
	service := NewOrderService(mockRepo, mockNotifier)
	
	// Test data
	customerID := "test-customer-1"
	items := []model.OrderItem{
		{
			ProductID:  "product-1",
			Name:       "Test Product",
			Quantity:   2,
			UnitPrice:  19.99,
		},
	}
	shippingAddress := "123 Test St, City, Country"
	
	// Test creating an order
	t.Run("Create order successfully", func(t *testing.T) {
		order, err := service.CreateOrder(ctx, customerID, items, shippingAddress)
		
		// Check no error occurred
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		
		// Check order attributes
		if order.CustomerID != customerID {
			t.Errorf("Expected customer ID %s, got %s", customerID, order.CustomerID)
		}
		
		if len(order.Items) != len(items) {
			t.Errorf("Expected %d items, got %d", len(items), len(order.Items))
		}
		
		if order.Status != model.OrderStatusCreated {
			t.Errorf("Expected status %s, got %s", model.OrderStatusCreated, order.Status)
		}
		
		// Check that order was saved to repository
		savedOrder, err := mockRepo.FindByID(ctx, order.ID)
		if err != nil {
			t.Fatalf("Failed to retrieve saved order: %v", err)
		}
		
		if savedOrder.ID != order.ID {
			t.Errorf("Expected saved order ID %s, got %s", order.ID, savedOrder.ID)
		}
		
		// Check that notification was sent
		if len(mockNotifier.notifications) != 1 {
			t.Errorf("Expected 1 notification, got %d", len(mockNotifier.notifications))
		} else {
			notification := mockNotifier.notifications[0]
			if notification.order.ID != order.ID {
				t.Errorf("Expected notification for order %s, got %s", order.ID, notification.order.ID)
			}
			if notification.notificationType != driven.NotificationTypeOrderCreated {
				t.Errorf("Expected notification type %s, got %s", driven.NotificationTypeOrderCreated, notification.notificationType)
			}
		}
	})
}

func TestOrderService_UpdateOrderStatus(t *testing.T) {
	// Create context for all tests
	ctx := context.Background()
	
	// Create test order
	order := &model.Order{
		ID:              "test-order-1",
		CustomerID:      "test-customer-1",
		Items: []model.OrderItem{
			{
				ProductID:  "product-1",
				Name:       "Test Product",
				Quantity:   2,
				UnitPrice:  19.99,
			},
		},
		Status:          model.OrderStatusCreated,
		TotalAmount:     39.98, // 2 * 19.99
		CreatedAt:       time.Now().Add(-1 * time.Hour),
		UpdatedAt:       time.Now().Add(-1 * time.Hour),
		ShippingAddress: "123 Test St, City, Country",
	}
	
	// Test updating order status
	t.Run("Update status successfully", func(t *testing.T) {
		// Create mocks
		mockRepo := newMockOrderRepository()
		mockNotifier := newMockNotificationService()
		
		// Save test order to repository
		mockRepo.Save(ctx, order)
		
		// Create the service
		service := NewOrderService(mockRepo, mockNotifier)
		
		// Update the order status
		err := service.UpdateOrderStatus(ctx, order.ID, model.OrderStatusProcessing)
		
		// Check no error occurred
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		
		// Check that order was updated in repository
		updatedOrder, err := mockRepo.FindByID(ctx, order.ID)
		if err != nil {
			t.Fatalf("Failed to retrieve updated order: %v", err)
		}
		
		if updatedOrder.Status != model.OrderStatusProcessing {
			t.Errorf("Expected status %s, got %s", model.OrderStatusProcessing, updatedOrder.Status)
		}
		
		// Check that notification was sent
		if len(mockNotifier.notifications) != 1 {
			t.Errorf("Expected 1 notification, got %d", len(mockNotifier.notifications))
		} else {
			notification := mockNotifier.notifications[0]
			if notification.order.ID != order.ID {
				t.Errorf("Expected notification for order %s, got %s", order.ID, notification.order.ID)
			}
			if notification.notificationType != driven.NotificationTypeOrderProcessing {
				t.Errorf("Expected notification type %s, got %s", driven.NotificationTypeOrderProcessing, notification.notificationType)
			}
		}
	})
	
	t.Run("Update with invalid status transition", func(t *testing.T) {
		// Create mocks
		mockRepo := newMockOrderRepository()
		mockNotifier := newMockNotificationService()
		
		// Create an order in "delivered" status (terminal state)
		deliveredOrder := *order
		deliveredOrder.Status = model.OrderStatusDelivered
		mockRepo.Save(ctx, &deliveredOrder)
		
		// Create the service
		service := NewOrderService(mockRepo, mockNotifier)
		
		// Try to update the order status back to "processing"
		err := service.UpdateOrderStatus(ctx, deliveredOrder.ID, model.OrderStatusProcessing)
		
		// Should return an error
		if err == nil {
			t.Error("Expected error for invalid status transition, got nil")
		}
		
		// Check that no notification was sent
		if len(mockNotifier.notifications) != 0 {
			t.Errorf("Expected no notifications, got %d", len(mockNotifier.notifications))
		}
	})
}

func TestOrderService_ListOrders(t *testing.T) {
	// Create context for all tests
	ctx := context.Background()
	
	// Create test orders
	order1 := &model.Order{
		ID:              "test-order-1",
		CustomerID:      "test-customer-1",
		Status:          model.OrderStatusCreated,
		ShippingAddress: "Address 1",
	}
	
	order2 := &model.Order{
		ID:              "test-order-2",
		CustomerID:      "test-customer-1",
		Status:          model.OrderStatusProcessing,
		ShippingAddress: "Address 2",
	}
	
	order3 := &model.Order{
		ID:              "test-order-3",
		CustomerID:      "test-customer-2",
		Status:          model.OrderStatusCreated,
		ShippingAddress: "Address 3",
	}
	
	// Create mocks
	mockRepo := newMockOrderRepository()
	mockNotifier := newMockNotificationService()
	
	// Save test orders to repository
	mockRepo.Save(ctx, order1)
	mockRepo.Save(ctx, order2)
	mockRepo.Save(ctx, order3)
	
	// Create the service
	service := NewOrderService(mockRepo, mockNotifier)
	
	// Test listing orders for a customer
	t.Run("List orders for customer", func(t *testing.T) {
		// List orders for customer 1
		orders, err := service.ListOrders(ctx, "test-customer-1")
		
		// Check no error occurred
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		
		// Check correct number of orders returned
		if len(orders) != 2 {
			t.Fatalf("Expected 2 orders, got %d", len(orders))
		}
		
		// Check all orders belong to the correct customer
		for _, order := range orders {
			if order.CustomerID != "test-customer-1" {
				t.Errorf("Expected order for customer test-customer-1, got %s", order.CustomerID)
			}
		}
		
		// List orders for customer 2
		orders, err = service.ListOrders(ctx, "test-customer-2")
		
		// Check no error occurred
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		
		// Check correct number of orders returned
		if len(orders) != 1 {
			t.Fatalf("Expected 1 order, got %d", len(orders))
		}
		
		// Check order belongs to the correct customer
		if orders[0].CustomerID != "test-customer-2" {
			t.Errorf("Expected order for customer test-customer-2, got %s", orders[0].CustomerID)
		}
	})
}
