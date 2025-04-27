package types

import (
	"time"
)

// Event is the base interface for all events
type Event interface {
	GetEventID() string
	GetAggregateID() string
	GetAggregateType() string
	GetEventType() string
	GetTimestamp() time.Time
}

// BaseEvent provides common fields for all events
type BaseEvent struct {
	EventID       string    `json:"event_id"`
	AggregateID   string    `json:"aggregate_id"`
	AggregateType string    `json:"aggregate_type"`
	EventType     string    `json:"event_type"`
	Timestamp     time.Time `json:"timestamp"`
	Version       int       `json:"version"`
}

// GetEventID returns the event ID
func (e BaseEvent) GetEventID() string {
	return e.EventID
}

// GetAggregateID returns the aggregate ID
func (e BaseEvent) GetAggregateID() string {
	return e.AggregateID
}

// GetAggregateType returns the aggregate type
func (e BaseEvent) GetAggregateType() string {
	return e.AggregateType
}

// GetEventType returns the event type
func (e BaseEvent) GetEventType() string {
	return e.EventType
}

// GetTimestamp returns the event timestamp
func (e BaseEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

// ProductCreatedEvent represents an event that occurs when a new product is created
type ProductCreatedEvent struct {
	BaseEvent
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

// ProductStockReducedEvent represents an event that occurs when product stock is reduced
type ProductStockReducedEvent struct {
	BaseEvent
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	NewStock  int    `json:"new_stock"`
}

// ProductStockAddedEvent represents an event that occurs when product stock is added
type ProductStockAddedEvent struct {
	BaseEvent
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	NewStock  int    `json:"new_stock"`
}

// ProductOutOfStockEvent represents an event that occurs when a product goes out of stock
type ProductOutOfStockEvent struct {
	BaseEvent
	ProductID string `json:"product_id"`
}

// ProductLowStockEvent represents an event that occurs when a product has low stock
type ProductLowStockEvent struct {
	BaseEvent
	ProductID string `json:"product_id"`
	Stock     int    `json:"stock"`
	Threshold int    `json:"threshold"`
}
