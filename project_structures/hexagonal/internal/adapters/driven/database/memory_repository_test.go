package database

import (
	"context"
	"testing"
	"time"
	
	"github.com/edgardnogueira/go-patterns/project_structures/hexagonal/internal/domain/model"
)

func TestMemoryOrderRepository(t *testing.T) {
	// Create context for all tests
	ctx := context.Background()
	
	// Create a sample order for testing
	createSampleOrder := func() *model.Order {
		return &model.Order{
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
	}
	
	t.Run("Save and FindByID", func(t *testing.T) {
		// Create a new repository
		repo := NewMemoryOrderRepository()
		
		// Create a sample order
		originalOrder := createSampleOrder()
		
		// Save the order
		err := repo.Save(ctx, originalOrder)
		if err != nil {
			t.Fatalf("Error saving order: %v", err)
		}
		
		// Retrieve the order
		retrievedOrder, err := repo.FindByID(ctx, originalOrder.ID)
		if err != nil {
			t.Fatalf("Error retrieving order: %v", err)
		}
		
		// Check that the retrieved order matches the original
		if retrievedOrder.ID != originalOrder.ID {
			t.Errorf("Expected ID %s, got %s", originalOrder.ID, retrievedOrder.ID)
		}
		if retrievedOrder.CustomerID != originalOrder.CustomerID {
			t.Errorf("Expected CustomerID %s, got %s", originalOrder.CustomerID, retrievedOrder.CustomerID)
		}
		if retrievedOrder.Status != originalOrder.Status {
			t.Errorf("Expected Status %s, got %s", originalOrder.Status, retrievedOrder.Status)
		}
		if retrievedOrder.TotalAmount != originalOrder.TotalAmount {
			t.Errorf("Expected TotalAmount %.2f, got %.2f", originalOrder.TotalAmount, retrievedOrder.TotalAmount)
		}
	})
	
	t.Run("Save and FindByCustomer", func(t *testing.T) {
		// Create a new repository
		repo := NewMemoryOrderRepository()
		
		// Create two orders for the same customer and one for another customer
		order1 := createSampleOrder()
		
		order2 := createSampleOrder()
		order2.ID = "test-order-2"
		
		order3 := createSampleOrder()
		order3.ID = "test-order-3"
		order3.CustomerID = "test-customer-2"
		
		// Save all orders
		repo.Save(ctx, order1)
		repo.Save(ctx, order2)
		repo.Save(ctx, order3)
		
		// Retrieve orders for customer 1
		orders, err := repo.FindByCustomer(ctx, "test-customer-1")
		if err != nil {
			t.Fatalf("Error retrieving orders by customer: %v", err)
		}
		
		// Check that we got the correct number of orders
		if len(orders) != 2 {
			t.Errorf("Expected 2 orders for customer, got %d", len(orders))
		}
		
		// Verify customer IDs
		for _, order := range orders {
			if order.CustomerID != "test-customer-1" {
				t.Errorf("Expected CustomerID test-customer-1, got %s", order.CustomerID)
			}
		}
		
		// Check for customer 2
		orders, err = repo.FindByCustomer(ctx, "test-customer-2")
		if err != nil {
			t.Fatalf("Error retrieving orders by customer: %v", err)
		}
		
		if len(orders) != 1 {
			t.Errorf("Expected 1 order for customer-2, got %d", len(orders))
		}
	})
	
	t.Run("Save and FindByStatus", func(t *testing.T) {
		// Create a new repository
		repo := NewMemoryOrderRepository()
		
		// Create orders with different statuses
		order1 := createSampleOrder() // Created status
		
		order2 := createSampleOrder()
		order2.ID = "test-order-2"
		order2.Status = model.OrderStatusProcessing
		
		order3 := createSampleOrder()
		order3.ID = "test-order-3"
		order3.Status = model.OrderStatusProcessing
		
		// Save all orders
		repo.Save(ctx, order1)
		repo.Save(ctx, order2)
		repo.Save(ctx, order3)
		
		// Retrieve orders with "created" status
		orders, err := repo.FindByStatus(ctx, model.OrderStatusCreated)
		if err != nil {
			t.Fatalf("Error retrieving orders by status: %v", err)
		}
		
		// Check that we got the correct number of orders
		if len(orders) != 1 {
			t.Errorf("Expected 1 order with 'created' status, got %d", len(orders))
		}
		
		// Retrieve orders with "processing" status
		orders, err = repo.FindByStatus(ctx, model.OrderStatusProcessing)
		if err != nil {
			t.Fatalf("Error retrieving orders by status: %v", err)
		}
		
		// Check that we got the correct number of orders
		if len(orders) != 2 {
			t.Errorf("Expected 2 orders with 'processing' status, got %d", len(orders))
		}
	})
	
	t.Run("Update existing order", func(t *testing.T) {
		// Create a new repository
		repo := NewMemoryOrderRepository()
		
		// Create and save an order
		order := createSampleOrder()
		repo.Save(ctx, order)
		
		// Modify the order and save again
		order.Status = model.OrderStatusProcessing
		order.UpdatedAt = time.Now()
		
		err := repo.Save(ctx, order)
		if err != nil {
			t.Fatalf("Error updating order: %v", err)
		}
		
		// Retrieve the updated order
		updatedOrder, err := repo.FindByID(ctx, order.ID)
		if err != nil {
			t.Fatalf("Error retrieving updated order: %v", err)
		}
		
		// Check that the status was updated
		if updatedOrder.Status != model.OrderStatusProcessing {
			t.Errorf("Expected updated status %s, got %s", model.OrderStatusProcessing, updatedOrder.Status)
		}
	})
	
	t.Run("Delete order", func(t *testing.T) {
		// Create a new repository
		repo := NewMemoryOrderRepository()
		
		// Create and save an order
		order := createSampleOrder()
		repo.Save(ctx, order)
		
		// Delete the order
		err := repo.Delete(ctx, order.ID)
		if err != nil {
			t.Fatalf("Error deleting order: %v", err)
		}
		
		// Try to retrieve the deleted order
		_, err = repo.FindByID(ctx, order.ID)
		if err == nil {
			t.Error("Expected error when retrieving deleted order, got nil")
		}
	})
	
	t.Run("Deep copy prevents external modifications", func(t *testing.T) {
		// Create a new repository
		repo := NewMemoryOrderRepository()
		
		// Create and save an order
		originalOrder := createSampleOrder()
		repo.Save(ctx, originalOrder)
		
		// Retrieve the order
		retrievedOrder, _ := repo.FindByID(ctx, originalOrder.ID)
		
		// Modify the retrieved order
		retrievedOrder.Status = model.OrderStatusProcessing
		retrievedOrder.Items[0].Quantity = 99
		
		// Retrieve the order again
		secondRetrieval, _ := repo.FindByID(ctx, originalOrder.ID)
		
		// Check that the second retrieval still has the original values
		if secondRetrieval.Status != model.OrderStatusCreated {
			t.Errorf("Expected original status %s, got %s", model.OrderStatusCreated, secondRetrieval.Status)
		}
		
		if secondRetrieval.Items[0].Quantity != 2 {
			t.Errorf("Expected original quantity 2, got %d", secondRetrieval.Items[0].Quantity)
		}
	})
}
