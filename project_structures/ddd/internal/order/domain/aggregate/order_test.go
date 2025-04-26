package aggregate

import (
	"testing"
	"time"
	
	"github.com/edgardnogueira/go-patterns/project_structures/ddd/internal/order/domain/entity"
	"github.com/edgardnogueira/go-patterns/project_structures/ddd/internal/order/domain/valueobject"
	"github.com/google/uuid"
)

func TestNewOrder(t *testing.T) {
	// Test valid order creation
	id := uuid.New().String()
	customerID := uuid.New().String()
	
	order, err := NewOrder(id, customerID)
	
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if order.ID() != id {
		t.Errorf("Expected order ID %s, got %s", id, order.ID())
	}
	
	if order.CustomerID() != customerID {
		t.Errorf("Expected customer ID %s, got %s", customerID, order.CustomerID())
	}
	
	if order.Status() != OrderStatusCreated {
		t.Errorf("Expected order status %s, got %s", OrderStatusCreated, order.Status())
	}
	
	if len(order.Items()) != 0 {
		t.Errorf("Expected empty items, got %d items", len(order.Items()))
	}
	
	if len(order.Events()) != 1 {
		t.Errorf("Expected 1 event, got %d events", len(order.Events()))
	}
	
	// Test invalid order creation - empty ID
	_, err = NewOrder("", customerID)
	if err == nil {
		t.Error("Expected error for empty ID, got nil")
	}
	
	// Test invalid order creation - empty customer ID
	_, err = NewOrder(id, "")
	if err == nil {
		t.Error("Expected error for empty customer ID, got nil")
	}
}

func TestAddItem(t *testing.T) {
	// Create a valid order
	orderID := uuid.New().String()
	customerID := uuid.New().String()
	
	order, err := NewOrder(orderID, customerID)
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}
	
	// Create a valid item
	itemID := uuid.New().String()
	productID := uuid.New().String()
	quantity := 2
	unitPrice := valueobject.MustNewMoney(1000, "USD") // $10.00
	description := "Test Product"
	
	item, err := entity.NewOrderItem(itemID, productID, quantity, unitPrice, description)
	if err != nil {
		t.Fatalf("Failed to create order item: %v", err)
	}
	
	// Add the item to the order
	err = order.AddItem(item)
	if err != nil {
		t.Errorf("Failed to add item to order: %v", err)
	}
	
	// Verify the item was added
	items := order.Items()
	if len(items) != 1 {
		t.Fatalf("Expected 1 item, got %d items", len(items))
	}
	
	if items[0].ID() != itemID {
		t.Errorf("Expected item ID %s, got %s", itemID, items[0].ID())
	}
	
	// Verify an event was added
	events := order.Events()
	if len(events) != 2 { // 1 for order creation, 1 for item addition
		t.Errorf("Expected 2 events, got %d events", len(events))
	}
	
	// Try to add the same item again (should fail)
	err = order.AddItem(item)
	if err == nil {
		t.Error("Expected error when adding duplicate item, got nil")
	}
}

func TestCalculateTotal(t *testing.T) {
	// Create a valid order
	orderID := uuid.New().String()
	customerID := uuid.New().String()
	
	order, err := NewOrder(orderID, customerID)
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}
	
	// Calculate total for empty order
	total, err := order.CalculateTotal()
	if err != nil {
		t.Errorf("Failed to calculate total for empty order: %v", err)
	}
	
	if total.Amount() != 0 {
		t.Errorf("Expected total 0, got %d", total.Amount())
	}
	
	// Add two items
	item1, err := entity.NewOrderItem(
		uuid.New().String(),
		uuid.New().String(),
		2,
		valueobject.MustNewMoney(1000, "USD"), // $10.00
		"Item 1",
	)
	if err != nil {
		t.Fatalf("Failed to create order item 1: %v", err)
	}
	
	item2, err := entity.NewOrderItem(
		uuid.New().String(),
		uuid.New().String(),
		1,
		valueobject.MustNewMoney(2000, "USD"), // $20.00
		"Item 2",
	)
	if err != nil {
		t.Fatalf("Failed to create order item 2: %v", err)
	}
	
	order.AddItem(item1)
	order.AddItem(item2)
	
	// Calculate total with items
	total, err = order.CalculateTotal()
	if err != nil {
		t.Errorf("Failed to calculate total: %v", err)
	}
	
	// Expected total: (2 * $10.00) + (1 * $20.00) = $40.00 = 4000 cents
	if total.Amount() != 4000 {
		t.Errorf("Expected total 4000 cents, got %d cents", total.Amount())
	}
}

func TestOrderStatusTransitions(t *testing.T) {
	// Create a valid order
	orderID := uuid.New().String()
	customerID := uuid.New().String()
	
	order, err := NewOrder(orderID, customerID)
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}
	
	// Add an item to the order
	item, err := entity.NewOrderItem(
		uuid.New().String(),
		uuid.New().String(),
		1,
		valueobject.MustNewMoney(1000, "USD"),
		"Test Item",
	)
	if err != nil {
		t.Fatalf("Failed to create order item: %v", err)
	}
	
	err = order.AddItem(item)
	if err != nil {
		t.Fatalf("Failed to add item to order: %v", err)
	}
	
	// Test status transitions
	
	// CREATED -> PAID (valid)
	err = order.MarkAsPaid()
	if err != nil {
		t.Errorf("Failed to mark order as paid: %v", err)
	}
	
	if order.Status() != OrderStatusPaid {
		t.Errorf("Expected status %s, got %s", OrderStatusPaid, order.Status())
	}
	
	// PAID -> SHIPPED (valid)
	err = order.MarkAsShipped()
	if err != nil {
		t.Errorf("Failed to mark order as shipped: %v", err)
	}
	
	if order.Status() != OrderStatusShipped {
		t.Errorf("Expected status %s, got %s", OrderStatusShipped, order.Status())
	}
	
	// SHIPPED -> DELIVERED (valid)
	err = order.MarkAsDelivered()
	if err != nil {
		t.Errorf("Failed to mark order as delivered: %v", err)
	}
	
	if order.Status() != OrderStatusDelivered {
		t.Errorf("Expected status %s, got %s", OrderStatusDelivered, order.Status())
	}
	
	// DELIVERED -> CANCELLED (invalid)
	err = order.Cancel()
	if err == nil {
		t.Error("Expected error when cancelling delivered order, got nil")
	}
}

func TestRemoveItem(t *testing.T) {
	// Create a valid order
	orderID := uuid.New().String()
	customerID := uuid.New().String()
	
	order, err := NewOrder(orderID, customerID)
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}
	
	// Create two items
	item1ID := uuid.New().String()
	item1, err := entity.NewOrderItem(
		item1ID,
		uuid.New().String(),
		2,
		valueobject.MustNewMoney(1000, "USD"),
		"Item 1",
	)
	if err != nil {
		t.Fatalf("Failed to create order item 1: %v", err)
	}
	
	item2ID := uuid.New().String()
	item2, err := entity.NewOrderItem(
		item2ID,
		uuid.New().String(),
		1,
		valueobject.MustNewMoney(2000, "USD"),
		"Item 2",
	)
	if err != nil {
		t.Fatalf("Failed to create order item 2: %v", err)
	}
	
	// Add both items
	order.AddItem(item1)
	order.AddItem(item2)
	
	// Verify both items were added
	if len(order.Items()) != 2 {
		t.Errorf("Expected 2 items, got %d items", len(order.Items()))
	}
	
	// Remove the first item
	err = order.RemoveItem(item1ID)
	if err != nil {
		t.Errorf("Failed to remove item: %v", err)
	}
	
	// Verify only the second item remains
	items := order.Items()
	if len(items) != 1 {
		t.Errorf("Expected 1 item after removal, got %d items", len(items))
	}
	
	if items[0].ID() != item2ID {
		t.Errorf("Expected remaining item ID %s, got %s", item2ID, items[0].ID())
	}
	
	// Try to remove a non-existent item
	err = order.RemoveItem(uuid.New().String())
	if err == nil {
		t.Error("Expected error when removing non-existent item, got nil")
	}
	
	// Try to remove an item after order is paid
	order.MarkAsPaid()
	err = order.RemoveItem(item2ID)
	if err == nil {
		t.Error("Expected error when removing item from paid order, got nil")
	}
}

func TestClearEvents(t *testing.T) {
	// Create a valid order
	orderID := uuid.New().String()
	customerID := uuid.New().String()
	
	order, err := NewOrder(orderID, customerID)
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}
	
	// Verify initial event (order created)
	if len(order.Events()) != 1 {
		t.Errorf("Expected 1 event initially, got %d events", len(order.Events()))
	}
	
	// Clear events
	order.ClearEvents()
	
	// Verify events were cleared
	if len(order.Events()) != 0 {
		t.Errorf("Expected 0 events after clearing, got %d events", len(order.Events()))
	}
	
	// Generate new events
	item, err := entity.NewOrderItem(
		uuid.New().String(),
		uuid.New().String(),
		1,
		valueobject.MustNewMoney(1000, "USD"),
		"Test Item",
	)
	if err != nil {
		t.Fatalf("Failed to create order item: %v", err)
	}
	
	order.AddItem(item)
	
	// Verify new events were added
	if len(order.Events()) != 1 {
		t.Errorf("Expected 1 event after adding item, got %d events", len(order.Events()))
	}
}
