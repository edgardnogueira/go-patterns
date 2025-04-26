package domain

import (
	"testing"
	"time"
)

func TestCreateOrder(t *testing.T) {
	// Set up test cases
	testCases := []struct {
		name       string
		customerID string
		items      []OrderItem
		wantErr    bool
	}{
		{
			name:       "Valid order",
			customerID: "customer123",
			items: []OrderItem{
				{ProductID: "prod1", Quantity: 2, Price: 10.0},
				{ProductID: "prod2", Quantity: 1, Price: 15.0},
			},
			wantErr: false,
		},
		{
			name:       "Empty customer ID",
			customerID: "",
			items: []OrderItem{
				{ProductID: "prod1", Quantity: 2, Price: 10.0},
			},
			wantErr: true,
		},
		{
			name:       "No items",
			customerID: "customer123",
			items:      []OrderItem{},
			wantErr:    true,
		},
		{
			name:       "Item with zero quantity",
			customerID: "customer123",
			items: []OrderItem{
				{ProductID: "prod1", Quantity: 0, Price: 10.0},
			},
			wantErr: true,
		},
		{
			name:       "Item with negative quantity",
			customerID: "customer123",
			items: []OrderItem{
				{ProductID: "prod1", Quantity: -1, Price: 10.0},
			},
			wantErr: true,
		},
		{
			name:       "Item with zero price",
			customerID: "customer123",
			items: []OrderItem{
				{ProductID: "prod1", Quantity: 2, Price: 0.0},
			},
			wantErr: true,
		},
		{
			name:       "Item with negative price",
			customerID: "customer123",
			items: []OrderItem{
				{ProductID: "prod1", Quantity: 2, Price: -10.0},
			},
			wantErr: true,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			order, err := CreateOrder(tc.customerID, tc.items)

			// Check error
			if tc.wantErr && err == nil {
				t.Errorf("CreateOrder() expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("CreateOrder() unexpected error: %v", err)
			}

			// If we don't expect an error, check the order properties
			if !tc.wantErr {
				// Check if order ID is generated
				if order.ID == "" {
					t.Errorf("CreateOrder() order ID is empty")
				}

				// Check customer ID
				if order.CustomerID != tc.customerID {
					t.Errorf("CreateOrder() customer ID = %v, want %v", order.CustomerID, tc.customerID)
				}

				// Check status
				if order.Status != OrderStatusCreated {
					t.Errorf("CreateOrder() status = %v, want %v", order.Status, OrderStatusCreated)
				}

				// Check total amount
				var expectedTotal float64
				for _, item := range tc.items {
					expectedTotal += float64(item.Quantity) * item.Price
				}
				if order.TotalAmount != expectedTotal {
					t.Errorf("CreateOrder() total amount = %v, want %v", order.TotalAmount, expectedTotal)
				}

				// Check items length
				if len(order.Items) != len(tc.items) {
					t.Errorf("CreateOrder() items length = %v, want %v", len(order.Items), len(tc.items))
				}

				// Check if timestamps are set
				now := time.Now()
				if order.CreatedAt.After(now) || order.CreatedAt.Before(now.Add(-10*time.Second)) {
					t.Errorf("CreateOrder() createdAt not in expected range: %v", order.CreatedAt)
				}
				if order.UpdatedAt.After(now) || order.UpdatedAt.Before(now.Add(-10*time.Second)) {
					t.Errorf("CreateOrder() updatedAt not in expected range: %v", order.UpdatedAt)
				}
			}
		})
	}
}

func TestOrder_ConfirmOrder(t *testing.T) {
	// Create a valid order for testing
	items := []OrderItem{
		{ProductID: "prod1", Quantity: 2, Price: 10.0},
	}
	order, _ := CreateOrder("customer123", items)

	// Test confirming a valid order
	err := order.ConfirmOrder()
	if err != nil {
		t.Errorf("ConfirmOrder() unexpected error: %v", err)
	}
	if order.Status != OrderStatusConfirmed {
		t.Errorf("ConfirmOrder() status = %v, want %v", order.Status, OrderStatusConfirmed)
	}

	// Test confirming an already confirmed order
	err = order.ConfirmOrder()
	if err == nil {
		t.Errorf("ConfirmOrder() expected error for already confirmed order, got nil")
	}

	// Test other state transitions
	order.Status = OrderStatusCreated
	_ = order.ConfirmOrder()
	_ = order.ShipOrder()
	err = order.ConfirmOrder() // Try to confirm a shipped order
	if err == nil {
		t.Errorf("ConfirmOrder() expected error for shipped order, got nil")
	}
}

func TestOrder_CancelOrder(t *testing.T) {
	// Create a valid order for testing
	items := []OrderItem{
		{ProductID: "prod1", Quantity: 2, Price: 10.0},
	}
	order, _ := CreateOrder("customer123", items)

	// Test cancelling a newly created order
	err := order.CancelOrder()
	if err != nil {
		t.Errorf("CancelOrder() unexpected error: %v", err)
	}
	if order.Status != OrderStatusCancelled {
		t.Errorf("CancelOrder() status = %v, want %v", order.Status, OrderStatusCancelled)
	}

	// Test cancelling an already cancelled order
	order, _ = CreateOrder("customer123", items)
	order.Status = OrderStatusCancelled
	err = order.CancelOrder()
	if err != nil {
		t.Errorf("CancelOrder() unexpected error for already cancelled order: %v", err)
	}

	// Test cancelling a confirmed order
	order, _ = CreateOrder("customer123", items)
	order.Status = OrderStatusConfirmed
	err = order.CancelOrder()
	if err != nil {
		t.Errorf("CancelOrder() unexpected error for confirmed order: %v", err)
	}

	// Test cancelling a shipped order
	order, _ = CreateOrder("customer123", items)
	order.Status = OrderStatusShipped
	err = order.CancelOrder()
	if err == nil {
		t.Errorf("CancelOrder() expected error for shipped order, got nil")
	}

	// Test cancelling a delivered order
	order, _ = CreateOrder("customer123", items)
	order.Status = OrderStatusDelivered
	err = order.CancelOrder()
	if err == nil {
		t.Errorf("CancelOrder() expected error for delivered order, got nil")
	}
}

func TestOrder_AddItem(t *testing.T) {
	// Create a valid order for testing
	items := []OrderItem{
		{ProductID: "prod1", Quantity: 2, Price: 10.0},
	}
	order, _ := CreateOrder("customer123", items)
	initialTotal := order.TotalAmount

	// Test adding a new item
	err := order.AddItem("prod2", 3, 5.0)
	if err != nil {
		t.Errorf("AddItem() unexpected error: %v", err)
	}
	if len(order.Items) != 2 {
		t.Errorf("AddItem() items length = %v, want %v", len(order.Items), 2)
	}
	expectedTotal := initialTotal + (3 * 5.0)
	if order.TotalAmount != expectedTotal {
		t.Errorf("AddItem() total amount = %v, want %v", order.TotalAmount, expectedTotal)
	}

	// Test adding an existing product
	initialTotal = order.TotalAmount
	initialQuantity := order.Items[1].Quantity
	err = order.AddItem("prod2", 2, 5.0)
	if err != nil {
		t.Errorf("AddItem() unexpected error when adding existing product: %v", err)
	}
	if len(order.Items) != 2 {
		t.Errorf("AddItem() items length = %v, want %v", len(order.Items), 2)
	}
	if order.Items[1].Quantity != initialQuantity+2 {
		t.Errorf("AddItem() quantity = %v, want %v", order.Items[1].Quantity, initialQuantity+2)
	}
	expectedTotal = initialTotal + (2 * 5.0)
	if order.TotalAmount != expectedTotal {
		t.Errorf("AddItem() total amount = %v, want %v", order.TotalAmount, expectedTotal)
	}

	// Test adding item to a confirmed order
	order.Status = OrderStatusConfirmed
	err = order.AddItem("prod3", 1, 20.0)
	if err == nil {
		t.Errorf("AddItem() expected error for confirmed order, got nil")
	}
}
