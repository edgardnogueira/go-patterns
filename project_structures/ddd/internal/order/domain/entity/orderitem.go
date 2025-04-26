package entity

import (
	"errors"
	"github.com/edgardnogueira/go-patterns/project_structures/ddd/internal/order/domain/valueobject"
)

// OrderItem represents an item in an order
type OrderItem struct {
	id          string
	productID   string
	quantity    int
	unitPrice   valueobject.Money
	description string
}

// NewOrderItem creates a new OrderItem
func NewOrderItem(id string, productID string, quantity int, unitPrice valueobject.Money, description string) (*OrderItem, error) {
	if id == "" {
		return nil, errors.New("order item ID cannot be empty")
	}
	
	if productID == "" {
		return nil, errors.New("product ID cannot be empty")
	}
	
	if quantity <= 0 {
		return nil, errors.New("quantity must be positive")
	}
	
	return &OrderItem{
		id:          id,
		productID:   productID,
		quantity:    quantity,
		unitPrice:   unitPrice,
		description: description,
	}, nil
}

// ID returns the order item ID
func (oi *OrderItem) ID() string {
	return oi.id
}

// ProductID returns the product ID
func (oi *OrderItem) ProductID() string {
	return oi.productID
}

// Quantity returns the quantity
func (oi *OrderItem) Quantity() int {
	return oi.quantity
}

// UnitPrice returns the unit price
func (oi *OrderItem) UnitPrice() valueobject.Money {
	return oi.unitPrice
}

// Description returns the description
func (oi *OrderItem) Description() string {
	return oi.description
}

// TotalPrice calculates the total price for this order item
func (oi *OrderItem) TotalPrice() valueobject.Money {
	return oi.unitPrice.Multiply(oi.quantity)
}

// UpdateQuantity updates the item quantity
func (oi *OrderItem) UpdateQuantity(newQuantity int) error {
	if newQuantity <= 0 {
		return errors.New("quantity must be positive")
	}
	oi.quantity = newQuantity
	return nil
}
