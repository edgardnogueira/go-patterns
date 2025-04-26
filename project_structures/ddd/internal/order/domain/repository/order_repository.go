package repository

import (
	"context"
	
	"github.com/edgardnogueira/go-patterns/project_structures/ddd/internal/order/domain/aggregate"
)

// OrderRepository defines the interface for order persistence operations
type OrderRepository interface {
	// FindByID finds an order by its ID
	FindByID(ctx context.Context, id string) (*aggregate.Order, error)
	
	// FindByCustomerID finds all orders for a customer
	FindByCustomerID(ctx context.Context, customerID string) ([]*aggregate.Order, error)
	
	// Save persists an order
	Save(ctx context.Context, order *aggregate.Order) error
	
	// Update updates an existing order
	Update(ctx context.Context, order *aggregate.Order) error
	
	// Delete removes an order
	Delete(ctx context.Context, id string) error
}
