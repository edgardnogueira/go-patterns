package driving

import (
	"context"
	
	"github.com/edgardnogueira/go-patterns/project_structures/hexagonal/internal/domain/model"
)

// OrderHandler is a primary port (driving port) that defines how external actors
// can interact with the domain for order-related operations
type OrderHandler interface {
	// CreateOrder creates a new order
	CreateOrder(ctx context.Context, customerID string, items []model.OrderItem, shippingAddress string) (*model.Order, error)
	
	// GetOrder retrieves an order by ID
	GetOrder(ctx context.Context, orderID string) (*model.Order, error)
	
	// UpdateOrderStatus updates the status of an order
	UpdateOrderStatus(ctx context.Context, orderID string, newStatus model.OrderStatus) error
	
	// ListOrders retrieves all orders for a specific customer
	ListOrders(ctx context.Context, customerID string) ([]*model.Order, error)
	
	// AddItemToOrder adds an item to an existing order
	AddItemToOrder(ctx context.Context, orderID string, item model.OrderItem) error
	
	// RemoveItemFromOrder removes an item from an existing order
	RemoveItemFromOrder(ctx context.Context, orderID string, productID string) error
}

// OrderProcessor is a primary port for processing orders asynchronously
type OrderProcessor interface {
	// ProcessPendingOrders processes all pending orders
	ProcessPendingOrders(ctx context.Context) error
	
	// ProcessOrder processes a specific order
	ProcessOrder(ctx context.Context, orderID string) error
}
