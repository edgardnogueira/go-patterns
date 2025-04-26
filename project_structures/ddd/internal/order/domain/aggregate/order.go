package aggregate

import (
	"errors"
	"time"

	"github.com/edgardnogueira/go-patterns/project_structures/ddd/internal/order/domain/entity"
	"github.com/edgardnogueira/go-patterns/project_structures/ddd/internal/order/domain/valueobject"
	"github.com/edgardnogueira/go-patterns/project_structures/ddd/internal/shared/domain/events"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusCreated     OrderStatus = "CREATED"
	OrderStatusPaid        OrderStatus = "PAID"
	OrderStatusShipped     OrderStatus = "SHIPPED"
	OrderStatusDelivered   OrderStatus = "DELIVERED"
	OrderStatusCancelled   OrderStatus = "CANCELLED"
)

// Order is an aggregate root that represents a customer order
type Order struct {
	id          string
	customerID  string
	items       []*entity.OrderItem
	status      OrderStatus
	createdAt   time.Time
	updatedAt   time.Time
	events      []events.DomainEvent
}

// NewOrder creates a new Order aggregate
func NewOrder(id string, customerID string) (*Order, error) {
	if id == "" {
		return nil, errors.New("order ID cannot be empty")
	}
	
	if customerID == "" {
		return nil, errors.New("customer ID cannot be empty")
	}
	
	now := time.Now()
	
	order := &Order{
		id:         id,
		customerID: customerID,
		items:      []*entity.OrderItem{},
		status:     OrderStatusCreated,
		createdAt:  now,
		updatedAt:  now,
		events:     []events.DomainEvent{},
	}
	
	// Add domain event
	order.AddEvent(events.NewOrderCreatedEvent(id, customerID))
	
	return order, nil
}

// ID returns the order ID
func (o *Order) ID() string {
	return o.id
}

// CustomerID returns the customer ID
func (o *Order) CustomerID() string {
	return o.customerID
}

// Status returns the order status
func (o *Order) Status() OrderStatus {
	return o.status
}

// CreatedAt returns the creation time
func (o *Order) CreatedAt() time.Time {
	return o.createdAt
}

// UpdatedAt returns the last update time
func (o *Order) UpdatedAt() time.Time {
	return o.updatedAt
}

// Items returns a copy of the order items slice
func (o *Order) Items() []*entity.OrderItem {
	// Return a copy to prevent external modification
	copyItems := make([]*entity.OrderItem, len(o.items))
	copy(copyItems, o.items)
	return copyItems
}

// AddItem adds an item to the order
func (o *Order) AddItem(item *entity.OrderItem) error {
	// Validate order state
	if o.status != OrderStatusCreated {
		return errors.New("cannot add items to non-CREATED orders")
	}
	
	// Check if item already exists
	for _, existingItem := range o.items {
		if existingItem.ID() == item.ID() {
			return errors.New("item with this ID already exists in the order")
		}
	}
	
	o.items = append(o.items, item)
	o.updateLastModified()
	
	// Add domain event
	o.AddEvent(events.NewOrderItemAddedEvent(o.id, item.ID(), item.ProductID(), item.Quantity()))
	
	return nil
}

// RemoveItem removes an item from the order
func (o *Order) RemoveItem(itemID string) error {
	// Validate order state
	if o.status != OrderStatusCreated {
		return errors.New("cannot remove items from non-CREATED orders")
	}
	
	for i, item := range o.items {
		if item.ID() == itemID {
			// Remove item from slice
			o.items = append(o.items[:i], o.items[i+1:]...)
			o.updateLastModified()
			
			// Add domain event
			o.AddEvent(events.NewOrderItemRemovedEvent(o.id, itemID))
			
			return nil
		}
	}
	
	return errors.New("item not found in order")
}

// CalculateTotal calculates the total price of the order
func (o *Order) CalculateTotal() (valueobject.Money, error) {
	if len(o.items) == 0 {
		return valueobject.MustNewMoney(0, "USD"), nil
	}
	
	// Start with the first item's total
	firstItem := o.items[0]
	total := firstItem.TotalPrice()
	
	// Add the total for each additional item
	for _, item := range o.items[1:] {
		itemTotal := item.TotalPrice()
		var err error
		total, err = total.Add(itemTotal)
		if err != nil {
			return valueobject.Money{}, err
		}
	}
	
	return total, nil
}

// MarkAsPaid marks the order as paid
func (o *Order) MarkAsPaid() error {
	if o.status != OrderStatusCreated {
		return errors.New("only CREATED orders can be marked as PAID")
	}
	
	if len(o.items) == 0 {
		return errors.New("cannot mark empty orders as PAID")
	}
	
	o.status = OrderStatusPaid
	o.updateLastModified()
	
	// Add domain event
	o.AddEvent(events.NewOrderPaidEvent(o.id))
	
	return nil
}

// MarkAsShipped marks the order as shipped
func (o *Order) MarkAsShipped() error {
	if o.status != OrderStatusPaid {
		return errors.New("only PAID orders can be marked as SHIPPED")
	}
	
	o.status = OrderStatusShipped
	o.updateLastModified()
	
	// Add domain event
	o.AddEvent(events.NewOrderShippedEvent(o.id))
	
	return nil
}

// MarkAsDelivered marks the order as delivered
func (o *Order) MarkAsDelivered() error {
	if o.status != OrderStatusShipped {
		return errors.New("only SHIPPED orders can be marked as DELIVERED")
	}
	
	o.status = OrderStatusDelivered
	o.updateLastModified()
	
	// Add domain event
	o.AddEvent(events.NewOrderDeliveredEvent(o.id))
	
	return nil
}

// Cancel cancels the order
func (o *Order) Cancel() error {
	if o.status == OrderStatusShipped || o.status == OrderStatusDelivered {
		return errors.New("cannot cancel SHIPPED or DELIVERED orders")
	}
	
	if o.status == OrderStatusCancelled {
		return errors.New("order is already cancelled")
	}
	
	o.status = OrderStatusCancelled
	o.updateLastModified()
	
	// Add domain event
	o.AddEvent(events.NewOrderCancelledEvent(o.id))
	
	return nil
}

// Events returns all accumulated domain events
func (o *Order) Events() []events.DomainEvent {
	// Return a copy to prevent external modification
	copyEvents := make([]events.DomainEvent, len(o.events))
	copy(copyEvents, o.events)
	return copyEvents
}

// AddEvent adds a domain event to the aggregate
func (o *Order) AddEvent(event events.DomainEvent) {
	o.events = append(o.events, event)
}

// ClearEvents clears all accumulated domain events
func (o *Order) ClearEvents() {
	o.events = []events.DomainEvent{}
}

// updateLastModified updates the last modified timestamp
func (o *Order) updateLastModified() {
	o.updatedAt = time.Now()
}
