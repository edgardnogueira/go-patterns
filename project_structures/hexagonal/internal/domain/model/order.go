package model

import (
	"errors"
	"time"
)

// OrderStatus represents the possible states of an order
type OrderStatus string

const (
	OrderStatusCreated     OrderStatus = "created"
	OrderStatusProcessing  OrderStatus = "processing"
	OrderStatusShipped     OrderStatus = "shipped"
	OrderStatusDelivered   OrderStatus = "delivered"
	OrderStatusCancelled   OrderStatus = "cancelled"
)

// OrderItem represents a single item in an order
type OrderItem struct {
	ProductID  string  `json:"product_id"`
	Name       string  `json:"name"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
}

// Order represents a customer order
type Order struct {
	ID            string       `json:"id"`
	CustomerID    string       `json:"customer_id"`
	Items         []OrderItem  `json:"items"`
	Status        OrderStatus  `json:"status"`
	TotalAmount   float64      `json:"total_amount"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	ShippingAddress string     `json:"shipping_address"`
}

// NewOrder creates a new order with the given details
func NewOrder(id, customerID string, items []OrderItem, shippingAddress string) (*Order, error) {
	if id == "" {
		return nil, errors.New("order ID cannot be empty")
	}
	
	if customerID == "" {
		return nil, errors.New("customer ID cannot be empty")
	}
	
	if len(items) == 0 {
		return nil, errors.New("order must contain at least one item")
	}
	
	if shippingAddress == "" {
		return nil, errors.New("shipping address cannot be empty")
	}
	
	now := time.Now()
	totalAmount := calculateTotalAmount(items)
	
	return &Order{
		ID:              id,
		CustomerID:      customerID,
		Items:           items,
		Status:          OrderStatusCreated,
		TotalAmount:     totalAmount,
		CreatedAt:       now,
		UpdatedAt:       now,
		ShippingAddress: shippingAddress,
	}, nil
}

// UpdateStatus changes the order status and updates the UpdatedAt timestamp
func (o *Order) UpdateStatus(status OrderStatus) error {
	// Validate status transitions
	if !isValidStatusTransition(o.Status, status) {
		return errors.New("invalid status transition")
	}
	
	o.Status = status
	o.UpdatedAt = time.Now()
	return nil
}

// AddItem adds an item to the order if the order is still in 'created' status
func (o *Order) AddItem(item OrderItem) error {
	if o.Status != OrderStatusCreated {
		return errors.New("cannot add items to an order that is already being processed")
	}
	
	if item.ProductID == "" || item.Quantity <= 0 || item.UnitPrice <= 0 {
		return errors.New("invalid order item")
	}
	
	// Check if item with same product already exists
	for i, existingItem := range o.Items {
		if existingItem.ProductID == item.ProductID {
			// Update quantity instead of adding new item
			o.Items[i].Quantity += item.Quantity
			o.TotalAmount = calculateTotalAmount(o.Items)
			o.UpdatedAt = time.Now()
			return nil
		}
	}
	
	// Add new item
	o.Items = append(o.Items, item)
	o.TotalAmount = calculateTotalAmount(o.Items)
	o.UpdatedAt = time.Now()
	return nil
}

// RemoveItem removes an item from the order if the order is still in 'created' status
func (o *Order) RemoveItem(productID string) error {
	if o.Status != OrderStatusCreated {
		return errors.New("cannot remove items from an order that is already being processed")
	}
	
	found := false
	items := make([]OrderItem, 0, len(o.Items)-1)
	
	for _, item := range o.Items {
		if item.ProductID != productID {
			items = append(items, item)
		} else {
			found = true
		}
	}
	
	if !found {
		return errors.New("item not found in order")
	}
	
	if len(items) == 0 {
		return errors.New("cannot remove last item from order")
	}
	
	o.Items = items
	o.TotalAmount = calculateTotalAmount(o.Items)
	o.UpdatedAt = time.Now()
	return nil
}

// calculateTotalAmount calculates the total amount for all items in the order
func calculateTotalAmount(items []OrderItem) float64 {
	var total float64
	for _, item := range items {
		total += float64(item.Quantity) * item.UnitPrice
	}
	return total
}

// isValidStatusTransition validates if a status transition is allowed
func isValidStatusTransition(current, new OrderStatus) bool {
	switch current {
	case OrderStatusCreated:
		return new == OrderStatusProcessing || new == OrderStatusCancelled
	case OrderStatusProcessing:
		return new == OrderStatusShipped || new == OrderStatusCancelled
	case OrderStatusShipped:
		return new == OrderStatusDelivered || new == OrderStatusCancelled
	case OrderStatusDelivered, OrderStatusCancelled:
		return false // Terminal states
	default:
		return false
	}
}
