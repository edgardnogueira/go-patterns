package repository

import (
	"context"
	"errors"
	"sync"
	
	"github.com/edgardnogueira/go-patterns/project_structures/ddd/internal/order/domain/aggregate"
	domainRepository "github.com/edgardnogueira/go-patterns/project_structures/ddd/internal/order/domain/repository"
)

// MemoryOrderRepository is an in-memory implementation of OrderRepository
type MemoryOrderRepository struct {
	orders map[string]*aggregate.Order
	mu     sync.RWMutex // Protects orders
}

// NewMemoryOrderRepository creates a new MemoryOrderRepository
func NewMemoryOrderRepository() domainRepository.OrderRepository {
	return &MemoryOrderRepository{
		orders: make(map[string]*aggregate.Order),
	}
}

// FindByID finds an order by its ID
func (r *MemoryOrderRepository) FindByID(ctx context.Context, id string) (*aggregate.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	order, exists := r.orders[id]
	if !exists {
		return nil, nil
	}
	
	return order, nil
}

// FindByCustomerID finds all orders for a customer
func (r *MemoryOrderRepository) FindByCustomerID(ctx context.Context, customerID string) ([]*aggregate.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var customerOrders []*aggregate.Order
	
	for _, order := range r.orders {
		if order.CustomerID() == customerID {
			customerOrders = append(customerOrders, order)
		}
	}
	
	return customerOrders, nil
}

// Save persists an order
func (r *MemoryOrderRepository) Save(ctx context.Context, order *aggregate.Order) error {
	if order == nil {
		return errors.New("cannot save nil order")
	}
	
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Check if order already exists
	if _, exists := r.orders[order.ID()]; exists {
		return errors.New("order already exists")
	}
	
	// Store the order
	r.orders[order.ID()] = order
	
	return nil
}

// Update updates an existing order
func (r *MemoryOrderRepository) Update(ctx context.Context, order *aggregate.Order) error {
	if order == nil {
		return errors.New("cannot update nil order")
	}
	
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Check if order exists
	if _, exists := r.orders[order.ID()]; !exists {
		return errors.New("order not found")
	}
	
	// Update the order
	r.orders[order.ID()] = order
	
	return nil
}

// Delete removes an order
func (r *MemoryOrderRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Check if order exists
	if _, exists := r.orders[id]; !exists {
		return errors.New("order not found")
	}
	
	// Delete the order
	delete(r.orders, id)
	
	return nil
}
