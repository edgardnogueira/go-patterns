package events

import (
	"context"
)

// EventHandler is a function that processes domain events
type EventHandler func(ctx context.Context, event DomainEvent) error

// EventBus is an interface for publishing and subscribing to domain events
type EventBus interface {
	// Publish publishes an event to all subscribers
	Publish(ctx context.Context, event DomainEvent) error
	
	// Subscribe registers a handler for a specific event type
	Subscribe(eventType string, handler EventHandler) error
	
	// Unsubscribe removes a handler for a specific event type
	Unsubscribe(eventType string, handler EventHandler) error
}
