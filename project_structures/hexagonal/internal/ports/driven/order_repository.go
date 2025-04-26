package driven

import (
	"context"
	
	"github.com/edgardnogueira/go-patterns/project_structures/hexagonal/internal/domain/model"
)

// OrderRepository is a secondary port (driven port) that defines 
// what the domain needs from a persistence mechanism
type OrderRepository interface {
	// Save stores or updates an order
	Save(ctx context.Context, order *model.Order) error
	
	// FindByID retrieves an order by its ID
	FindByID(ctx context.Context, orderID string) (*model.Order, error)
	
	// FindByCustomer retrieves all orders for a customer
	FindByCustomer(ctx context.Context, customerID string) ([]*model.Order, error)
	
	// FindByStatus retrieves all orders with a specific status
	FindByStatus(ctx context.Context, status model.OrderStatus) ([]*model.Order, error)
	
	// Delete removes an order
	Delete(ctx context.Context, orderID string) error
}
