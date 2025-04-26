package domain

import (
	"errors"
	"time"

	"github.com/edgardnogueira/go-patterns/project_structures/microservices/pkg/common"
)

// Product represents a product in the inventory
type Product struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	SKU         string    `json:"sku" db:"sku"`
	Price       float64   `json:"price" db:"price"`
	Quantity    int       `json:"quantity" db:"quantity"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// InventoryReservation represents a reservation of inventory
type InventoryReservation struct {
	ID         string    `json:"id" db:"id"`
	ProductID  string    `json:"product_id" db:"product_id"`
	OrderID    string    `json:"order_id" db:"order_id"`
	Quantity   int       `json:"quantity" db:"quantity"`
	ReservedAt time.Time `json:"reserved_at" db:"reserved_at"`
	ExpiresAt  time.Time `json:"expires_at" db:"expires_at"`
	Status     string    `json:"status" db:"status"`
}

// NewProduct creates a new product
func NewProduct(name, description, sku string, price float64, quantity int) (*Product, error) {
	if name == "" {
		return nil, errors.New("product name is required")
	}
	
	if sku == "" {
		return nil, errors.New("product SKU is required")
	}
	
	if price <= 0 {
		return nil, errors.New("product price must be positive")
	}
	
	if quantity < 0 {
		return nil, errors.New("product quantity cannot be negative")
	}
	
	now := time.Now()
	return &Product{
		ID:          common.GenerateUUID(),
		Name:        name,
		Description: description,
		SKU:         sku,
		Price:       price,
		Quantity:    quantity,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// UpdateStock updates the product's stock quantity
func (p *Product) UpdateStock(quantity int) error {
	if quantity < 0 {
		return errors.New("product quantity cannot be negative")
	}
	
	p.Quantity = quantity
	p.UpdatedAt = time.Now()
	return nil
}

// ReserveStock reserves stock for an order
func (p *Product) ReserveStock(orderID string, quantity int) (*InventoryReservation, error) {
	if quantity <= 0 {
		return nil, errors.New("reservation quantity must be positive")
	}
	
	if quantity > p.Quantity {
		return nil, errors.New("insufficient inventory")
	}
	
	// Reduce available quantity
	p.Quantity -= quantity
	p.UpdatedAt = time.Now()
	
	// Create reservation
	now := time.Now()
	reservation := &InventoryReservation{
		ID:         common.GenerateUUID(),
		ProductID:  p.ID,
		OrderID:    orderID,
		Quantity:   quantity,
		ReservedAt: now,
		ExpiresAt:  now.Add(24 * time.Hour), // Expires in 24 hours
		Status:     "active",
	}
	
	return reservation, nil
}

// ReleaseStock releases previously reserved stock
func (p *Product) ReleaseStock(quantity int) error {
	p.Quantity += quantity
	p.UpdatedAt = time.Now()
	return nil
}

// IsAvailable checks if there is enough stock available
func (p *Product) IsAvailable(quantity int) bool {
	return p.Quantity >= quantity
}
