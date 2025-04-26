package model

import (
	"testing"
	"time"
)

func TestNewOrder(t *testing.T) {
	// Setup valid order parameters
	id := "order-123"
	customerID := "customer-456"
	items := []OrderItem{
		{
			ProductID:  "product-001",
			Name:       "Test Product",
			Quantity:   2,
			UnitPrice:  19.99,
		},
	}
	shippingAddress := "123 Test St, City, Country"
	
	// Test successful order creation
	t.Run("Valid order creation", func(t *testing.T) {
		order, err := NewOrder(id, customerID, items, shippingAddress)
		
		// Check no error occurred
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		
		// Check order attributes
		if order.ID != id {
			t.Errorf("Expected order ID %s, got %s", id, order.ID)
		}
		if order.CustomerID != customerID {
			t.Errorf("Expected customer ID %s, got %s", customerID, order.CustomerID)
		}
		if len(order.Items) != len(items) {
			t.Errorf("Expected %d items, got %d", len(items), len(order.Items))
		}
		if order.Status != OrderStatusCreated {
			t.Errorf("Expected status %s, got %s", OrderStatusCreated, order.Status)
		}
		
		// Check total amount calculation
		expectedTotal := 39.98 // 2 * 19.99
		if order.TotalAmount != expectedTotal {
			t.Errorf("Expected total amount %.2f, got %.2f", expectedTotal, order.TotalAmount)
		}
		
		// Check timestamps
		if order.CreatedAt.IsZero() {
			t.Error("Expected CreatedAt to be set")
		}
		if order.UpdatedAt.IsZero() {
			t.Error("Expected UpdatedAt to be set")
		}
	})
	
	// Test validation errors
	t.Run("Empty order ID", func(t *testing.T) {
		_, err := NewOrder("", customerID, items, shippingAddress)
		if err == nil {
			t.Error("Expected error for empty order ID, got nil")
		}
	})
	
	t.Run("Empty customer ID", func(t *testing.T) {
		_, err := NewOrder(id, "", items, shippingAddress)
		if err == nil {
			t.Error("Expected error for empty customer ID, got nil")
		}
	})
	
	t.Run("Empty items", func(t *testing.T) {
		_, err := NewOrder(id, customerID, []OrderItem{}, shippingAddress)
		if err == nil {
			t.Error("Expected error for empty items, got nil")
		}
	})
	
	t.Run("Empty shipping address", func(t *testing.T) {
		_, err := NewOrder(id, customerID, items, "")
		if err == nil {
			t.Error("Expected error for empty shipping address, got nil")
		}
	})
}

func TestOrder_UpdateStatus(t *testing.T) {
	// Create a test order
	order := &Order{
		ID:              "order-123",
		CustomerID:      "customer-456",
		Items: []OrderItem{
			{
				ProductID:  "product-001",
				Name:       "Test Product",
				Quantity:   1,
				UnitPrice:  29.99,
			},
		},
		Status:          OrderStatusCreated,
		TotalAmount:     29.99,
		CreatedAt:       time.Now().Add(-1 * time.Hour), // 1 hour ago
		UpdatedAt:       time.Now().Add(-1 * time.Hour),
		ShippingAddress: "123 Test St, City, Country",
	}
	
	// Test valid status transitions
	t.Run("Valid transition: created -> processing", func(t *testing.T) {
		initialUpdatedAt := order.UpdatedAt
		err := order.UpdateStatus(OrderStatusProcessing)
		
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if order.Status != OrderStatusProcessing {
			t.Errorf("Expected status %s, got %s", OrderStatusProcessing, order.Status)
		}
		if !order.UpdatedAt.After(initialUpdatedAt) {
			t.Error("Expected UpdatedAt to be updated")
		}
	})
	
	t.Run("Valid transition: processing -> shipped", func(t *testing.T) {
		initialUpdatedAt := order.UpdatedAt
		err := order.UpdateStatus(OrderStatusShipped)
		
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if order.Status != OrderStatusShipped {
			t.Errorf("Expected status %s, got %s", OrderStatusShipped, order.Status)
		}
		if !order.UpdatedAt.After(initialUpdatedAt) {
			t.Error("Expected UpdatedAt to be updated")
		}
	})
	
	t.Run("Valid transition: shipped -> delivered", func(t *testing.T) {
		initialUpdatedAt := order.UpdatedAt
		err := order.UpdateStatus(OrderStatusDelivered)
		
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if order.Status != OrderStatusDelivered {
			t.Errorf("Expected status %s, got %s", OrderStatusDelivered, order.Status)
		}
		if !order.UpdatedAt.After(initialUpdatedAt) {
			t.Error("Expected UpdatedAt to be updated")
		}
	})
	
	// Test invalid status transitions
	t.Run("Invalid transition: delivered -> processing", func(t *testing.T) {
		err := order.UpdateStatus(OrderStatusProcessing)
		
		if err == nil {
			t.Error("Expected error for invalid transition, got nil")
		}
		if order.Status != OrderStatusDelivered {
			t.Errorf("Expected status to remain %s, got %s", OrderStatusDelivered, order.Status)
		}
	})
}

func TestOrder_AddItem(t *testing.T) {
	// Create a test order
	order := &Order{
		ID:              "order-123",
		CustomerID:      "customer-456",
		Items: []OrderItem{
			{
				ProductID:  "product-001",
				Name:       "Test Product 1",
				Quantity:   1,
				UnitPrice:  29.99,
			},
		},
		Status:          OrderStatusCreated,
		TotalAmount:     29.99,
		CreatedAt:       time.Now().Add(-1 * time.Hour),
		UpdatedAt:       time.Now().Add(-1 * time.Hour),
		ShippingAddress: "123 Test St, City, Country",
	}
	
	// Test adding a new item
	t.Run("Add new item", func(t *testing.T) {
		newItem := OrderItem{
			ProductID:  "product-002",
			Name:       "Test Product 2",
			Quantity:   2,
			UnitPrice:  15.50,
		}
		
		initialUpdatedAt := order.UpdatedAt
		err := order.AddItem(newItem)
		
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		
		// Check item was added
		if len(order.Items) != 2 {
			t.Fatalf("Expected 2 items, got %d", len(order.Items))
		}
		
		// Check total amount was updated
		expectedTotal := 29.99 + (2 * 15.50) // Original + new items
		if order.TotalAmount != expectedTotal {
			t.Errorf("Expected total amount %.2f, got %.2f", expectedTotal, order.TotalAmount)
		}
		
		// Check UpdatedAt was updated
		if !order.UpdatedAt.After(initialUpdatedAt) {
			t.Error("Expected UpdatedAt to be updated")
		}
	})
	
	// Test updating existing item quantity
	t.Run("Update existing item quantity", func(t *testing.T) {
		existingItem := OrderItem{
			ProductID:  "product-001",
			Name:       "Test Product 1",
			Quantity:   3,
			UnitPrice:  29.99,
		}
		
		initialUpdatedAt := order.UpdatedAt
		err := order.AddItem(existingItem)
		
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		
		// Check item count remains the same
		if len(order.Items) != 2 {
			t.Fatalf("Expected item count to remain 2, got %d", len(order.Items))
		}
		
		// Check item quantity was updated
		if order.Items[0].Quantity != 4 { // 1 original + 3 new
			t.Errorf("Expected quantity 4, got %d", order.Items[0].Quantity)
		}
		
		// Check total amount was updated correctly
		expectedTotal := (4 * 29.99) + (2 * 15.50)
		if order.TotalAmount != expectedTotal {
			t.Errorf("Expected total amount %.2f, got %.2f", expectedTotal, order.TotalAmount)
		}
		
		// Check UpdatedAt was updated
		if !order.UpdatedAt.After(initialUpdatedAt) {
			t.Error("Expected UpdatedAt to be updated")
		}
	})
	
	// Test adding item to an order that is already in progress
	t.Run("Add item to processing order", func(t *testing.T) {
		// Change order status to processing
		order.Status = OrderStatusProcessing
		
		newItem := OrderItem{
			ProductID:  "product-003",
			Name:       "Test Product 3",
			Quantity:   1,
			UnitPrice:  9.99,
		}
		
		err := order.AddItem(newItem)
		
		// Should return an error
		if err == nil {
			t.Error("Expected error when adding item to processing order, got nil")
		}
		
		// Check item count remains the same
		if len(order.Items) != 2 {
			t.Errorf("Expected item count to remain 2, got %d", len(order.Items))
		}
	})
}
