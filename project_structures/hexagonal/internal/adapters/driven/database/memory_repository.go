package database

import (
	"context"
	"errors"
	"sync"
	
	"github.com/edgardnogueira/go-patterns/project_structures/hexagonal/internal/domain/model"
)

// MemoryOrderRepository is an implementation of the OrderRepository interface
// that stores data in memory
type MemoryOrderRepository struct {
	orders map[string]*model.Order
	mutex  sync.RWMutex
}

// NewMemoryOrderRepository creates a new in-memory order repository
func NewMemoryOrderRepository() *MemoryOrderRepository {
	return &MemoryOrderRepository{
		orders: make(map[string]*model.Order),
	}
}

// Save stores or updates an order
func (r *MemoryOrderRepository) Save(ctx context.Context, order *model.Order) error {
	if order == nil {
		return errors.New("cannot save nil order")
	}
	
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// Create a deep copy to prevent external modifications
	r.orders[order.ID] = copyOrder(order)
	
	return nil
}

// FindByID retrieves an order by its ID
func (r *MemoryOrderRepository) FindByID(ctx context.Context, orderID string) (*model.Order, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	order, exists := r.orders[orderID]
	if !exists {
		return nil, errors.New("order not found")
	}
	
	// Return a copy to prevent external modifications
	return copyOrder(order), nil
}

// FindByCustomer retrieves all orders for a customer
func (r *MemoryOrderRepository) FindByCustomer(ctx context.Context, customerID string) ([]*model.Order, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	var result []*model.Order
	
	for _, order := range r.orders {
		if order.CustomerID == customerID {
			// Add a copy to the result
			result = append(result, copyOrder(order))
		}
	}
	
	return result, nil
}

// FindByStatus retrieves all orders with a specific status
func (r *MemoryOrderRepository) FindByStatus(ctx context.Context, status model.OrderStatus) ([]*model.Order, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	var result []*model.Order
	
	for _, order := range r.orders {
		if order.Status == status {
			// Add a copy to the result
			result = append(result, copyOrder(order))
		}
	}
	
	return result, nil
}

// Delete removes an order
func (r *MemoryOrderRepository) Delete(ctx context.Context, orderID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.orders[orderID]; !exists {
		return errors.New("order not found")
	}
	
	delete(r.orders, orderID)
	
	return nil
}

// copyOrder creates a deep copy of an order to prevent external modifications
func copyOrder(order *model.Order) *model.Order {
	if order == nil {
		return nil
	}
	
	// Copy items
	itemsCopy := make([]model.OrderItem, len(order.Items))
	copy(itemsCopy, order.Items)
	
	// Create a new order with the same values
	return &model.Order{
		ID:              order.ID,
		CustomerID:      order.CustomerID,
		Items:           itemsCopy,
		Status:          order.Status,
		TotalAmount:     order.TotalAmount,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
		ShippingAddress: order.ShippingAddress,
	}
}
