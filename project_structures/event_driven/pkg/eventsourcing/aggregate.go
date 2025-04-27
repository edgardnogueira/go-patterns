package eventsourcing

import (
	"errors"
	"time"

	"github.com/edgardnogueira/go-patterns/project_structures/event_driven/internal/events/types"
	"github.com/google/uuid"
)

// EventApplier is an interface for applying events to an aggregate
type EventApplier interface {
	ApplyEvent(event types.Event)
}

// Aggregate is the base struct for all aggregates
type Aggregate struct {
	ID            string
	Type          string
	Version       int
	Changes       []types.Event
	EventAppliers map[string]func(types.Event)
}

// NewAggregate creates a new aggregate
func NewAggregate(id, aggregateType string) *Aggregate {
	return &Aggregate{
		ID:            id,
		Type:          aggregateType,
		Version:       0,
		Changes:       []types.Event{},
		EventAppliers: make(map[string]func(types.Event)),
	}
}

// RegisterEventApplier registers a function to apply events of a specific type
func (a *Aggregate) RegisterEventApplier(eventType string, applier func(types.Event)) {
	a.EventAppliers[eventType] = applier
}

// ApplyEvent applies an event to the aggregate
func (a *Aggregate) ApplyEvent(event types.Event) {
	// Apply the event using the registered applier
	applier, exists := a.EventAppliers[event.GetEventType()]
	if exists {
		applier(event)
	}

	// Increment version
	a.Version++
}

// ApplyEvents applies multiple events to the aggregate
func (a *Aggregate) ApplyEvents(events []types.Event) {
	for _, event := range events {
		a.ApplyEvent(event)
	}
}

// AddEvent adds a new event to the aggregate's changes
func (a *Aggregate) AddEvent(eventType string, data map[string]interface{}) {
	// Create base event
	baseEvent := types.BaseEvent{
		EventID:       uuid.New().String(),
		AggregateID:   a.ID,
		AggregateType: a.Type,
		EventType:     eventType,
		Timestamp:     time.Now(),
		Version:       a.Version + 1,
	}

	// Create specific event type based on eventType
	var event types.Event
	switch eventType {
	case "product.created":
		event = types.ProductCreatedEvent{
			BaseEvent:   baseEvent,
			ID:          data["id"].(string),
			Name:        data["name"].(string),
			Description: data["description"].(string),
			Price:       data["price"].(float64),
			Stock:       data["stock"].(int),
		}
	case "product.stock.reduced":
		event = types.ProductStockReducedEvent{
			BaseEvent:  baseEvent,
			ProductID:  data["product_id"].(string),
			Quantity:   data["quantity"].(int),
			NewStock:   data["new_stock"].(int),
		}
	case "product.stock.added":
		event = types.ProductStockAddedEvent{
			BaseEvent:  baseEvent,
			ProductID:  data["product_id"].(string),
			Quantity:   data["quantity"].(int),
			NewStock:   data["new_stock"].(int),
		}
	case "product.stock.out":
		event = types.ProductOutOfStockEvent{
			BaseEvent:  baseEvent,
			ProductID:  data["product_id"].(string),
		}
	case "product.stock.low":
		event = types.ProductLowStockEvent{
			BaseEvent:  baseEvent,
			ProductID:  data["product_id"].(string),
			Stock:      data["stock"].(int),
			Threshold:  data["threshold"].(int),
		}
	default:
		// Unknown event type - should not happen
		return
	}

	// Apply the event to update the aggregate state
	a.ApplyEvent(event)

	// Add to changes
	a.Changes = append(a.Changes, event)
}

// GetUncommittedChanges returns all uncommitted changes
func (a *Aggregate) GetUncommittedChanges() []types.Event {
	return a.Changes
}

// ClearUncommittedChanges clears all uncommitted changes
func (a *Aggregate) ClearUncommittedChanges() {
	a.Changes = []types.Event{}
}
