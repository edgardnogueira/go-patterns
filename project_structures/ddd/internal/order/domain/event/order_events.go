package event

import (
	"fmt"
	"github.com/edgardnogueira/go-patterns/project_structures/ddd/internal/shared/domain/events"
	"github.com/google/uuid"
)

// Event types
const (
	OrderCreatedEventType    = "order.created"
	OrderPaidEventType       = "order.paid"
	OrderShippedEventType    = "order.shipped"
	OrderDeliveredEventType  = "order.delivered"
	OrderCancelledEventType  = "order.cancelled"
	OrderItemAddedEventType  = "order.item_added"
	OrderItemRemovedEventType = "order.item_removed"
)

// OrderCreatedEvent represents the event when an order is created
type OrderCreatedEvent struct {
	events.BaseDomainEvent
}

// NewOrderCreatedEvent creates a new OrderCreatedEvent
func NewOrderCreatedEvent(orderID, customerID string) OrderCreatedEvent {
	payload := map[string]interface{}{
		"customer_id": customerID,
	}
	
	base := events.NewBaseDomainEvent(
		uuid.New().String(),
		OrderCreatedEventType,
		orderID,
		"Order",
		payload,
	)
	
	return OrderCreatedEvent{BaseDomainEvent: base}
}

// OrderPaidEvent represents the event when an order is paid
type OrderPaidEvent struct {
	events.BaseDomainEvent
}

// NewOrderPaidEvent creates a new OrderPaidEvent
func NewOrderPaidEvent(orderID string) OrderPaidEvent {
	base := events.NewBaseDomainEvent(
		uuid.New().String(),
		OrderPaidEventType,
		orderID,
		"Order",
		map[string]interface{}{},
	)
	
	return OrderPaidEvent{BaseDomainEvent: base}
}

// OrderShippedEvent represents the event when an order is shipped
type OrderShippedEvent struct {
	events.BaseDomainEvent
}

// NewOrderShippedEvent creates a new OrderShippedEvent
func NewOrderShippedEvent(orderID string) OrderShippedEvent {
	base := events.NewBaseDomainEvent(
		uuid.New().String(),
		OrderShippedEventType,
		orderID,
		"Order",
		map[string]interface{}{},
	)
	
	return OrderShippedEvent{BaseDomainEvent: base}
}

// OrderDeliveredEvent represents the event when an order is delivered
type OrderDeliveredEvent struct {
	events.BaseDomainEvent
}

// NewOrderDeliveredEvent creates a new OrderDeliveredEvent
func NewOrderDeliveredEvent(orderID string) OrderDeliveredEvent {
	base := events.NewBaseDomainEvent(
		uuid.New().String(),
		OrderDeliveredEventType,
		orderID,
		"Order",
		map[string]interface{}{},
	)
	
	return OrderDeliveredEvent{BaseDomainEvent: base}
}

// OrderCancelledEvent represents the event when an order is cancelled
type OrderCancelledEvent struct {
	events.BaseDomainEvent
}

// NewOrderCancelledEvent creates a new OrderCancelledEvent
func NewOrderCancelledEvent(orderID string) OrderCancelledEvent {
	base := events.NewBaseDomainEvent(
		uuid.New().String(),
		OrderCancelledEventType,
		orderID,
		"Order",
		map[string]interface{}{},
	)
	
	return OrderCancelledEvent{BaseDomainEvent: base}
}

// OrderItemAddedEvent represents the event when an item is added to an order
type OrderItemAddedEvent struct {
	events.BaseDomainEvent
}

// NewOrderItemAddedEvent creates a new OrderItemAddedEvent
func NewOrderItemAddedEvent(orderID, itemID, productID string, quantity int) OrderItemAddedEvent {
	payload := map[string]interface{}{
		"item_id":    itemID,
		"product_id": productID,
		"quantity":   quantity,
	}
	
	base := events.NewBaseDomainEvent(
		uuid.New().String(),
		OrderItemAddedEventType,
		orderID,
		"Order",
		payload,
	)
	
	return OrderItemAddedEvent{BaseDomainEvent: base}
}

// OrderItemRemovedEvent represents the event when an item is removed from an order
type OrderItemRemovedEvent struct {
	events.BaseDomainEvent
}

// NewOrderItemRemovedEvent creates a new OrderItemRemovedEvent
func NewOrderItemRemovedEvent(orderID, itemID string) OrderItemRemovedEvent {
	payload := map[string]interface{}{
		"item_id": itemID,
	}
	
	base := events.NewBaseDomainEvent(
		uuid.New().String(),
		OrderItemRemovedEventType,
		orderID,
		"Order",
		payload,
	)
	
	return OrderItemRemovedEvent{BaseDomainEvent: base}
}

// String returns a string representation of the event
func (e OrderCreatedEvent) String() string {
	return fmt.Sprintf("OrderCreatedEvent{OrderID: %s, CustomerID: %s}", 
		e.AggregateID(), e.Payload()["customer_id"])
}

// String returns a string representation of the event
func (e OrderPaidEvent) String() string {
	return fmt.Sprintf("OrderPaidEvent{OrderID: %s}", e.AggregateID())
}

// String returns a string representation of the event
func (e OrderShippedEvent) String() string {
	return fmt.Sprintf("OrderShippedEvent{OrderID: %s}", e.AggregateID())
}

// String returns a string representation of the event
func (e OrderDeliveredEvent) String() string {
	return fmt.Sprintf("OrderDeliveredEvent{OrderID: %s}", e.AggregateID())
}

// String returns a string representation of the event
func (e OrderCancelledEvent) String() string {
	return fmt.Sprintf("OrderCancelledEvent{OrderID: %s}", e.AggregateID())
}

// String returns a string representation of the event
func (e OrderItemAddedEvent) String() string {
	return fmt.Sprintf("OrderItemAddedEvent{OrderID: %s, ItemID: %s, ProductID: %s, Quantity: %d}", 
		e.AggregateID(), e.Payload()["item_id"], e.Payload()["product_id"], e.Payload()["quantity"])
}

// String returns a string representation of the event
func (e OrderItemRemovedEvent) String() string {
	return fmt.Sprintf("OrderItemRemovedEvent{OrderID: %s, ItemID: %s}", 
		e.AggregateID(), e.Payload()["item_id"])
}
