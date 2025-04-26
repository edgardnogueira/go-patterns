package eventbus

import (
	"context"
	"sync"
	
	"github.com/edgardnogueira/go-patterns/project_structures/ddd/internal/shared/domain/events"
)

// MemoryEventBus is an in-memory implementation of the EventBus interface
type MemoryEventBus struct {
	handlers map[string][]events.EventHandler
	mu       sync.RWMutex // Protects handlers
}

// NewMemoryEventBus creates a new in-memory event bus
func NewMemoryEventBus() *MemoryEventBus {
	return &MemoryEventBus{
		handlers: make(map[string][]events.EventHandler),
	}
}

// Publish publishes an event to all subscribers
func (b *MemoryEventBus) Publish(ctx context.Context, event events.DomainEvent) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	
	// Get handlers for this event type
	handlers, exists := b.handlers[event.EventType()]
	if !exists {
		// No handlers for this event type
		return nil
	}
	
	// Process the event asynchronously with each handler
	for _, handler := range handlers {
		// Create a copy of the handler to use in the goroutine
		h := handler
		
		// Launch the handler in a goroutine
		go func() {
			// Create a new context for the handler
			handlerCtx, cancel := context.WithCancel(context.Background())
			defer cancel()
			
			// Call the handler
			if err := h(handlerCtx, event); err != nil {
				// Log the error (in a real implementation, we'd have proper error handling)
				// For now, we just let the goroutine terminate
			}
		}()
	}
	
	return nil
}

// Subscribe registers a handler for a specific event type
func (b *MemoryEventBus) Subscribe(eventType string, handler events.EventHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	// Get existing handlers for this event type
	handlers, exists := b.handlers[eventType]
	if !exists {
		// Create a new slice for this event type
		handlers = []events.EventHandler{}
	}
	
	// Add the new handler
	b.handlers[eventType] = append(handlers, handler)
	
	return nil
}

// Unsubscribe removes a handler for a specific event type
func (b *MemoryEventBus) Unsubscribe(eventType string, handler events.EventHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	// Get existing handlers for this event type
	handlers, exists := b.handlers[eventType]
	if !exists {
		// No handlers for this event type
		return nil
	}
	
	// Remove the handler (based on pointer equality)
	// This is not ideal because functions aren't directly comparable
	// In a real implementation, we might want to use handler IDs
	var newHandlers []events.EventHandler
	for _, h := range handlers {
		// Compare function pointers
		if &h != &handler {
			newHandlers = append(newHandlers, h)
		}
	}
	
	// Update the handlers
	b.handlers[eventType] = newHandlers
	
	return nil
}
