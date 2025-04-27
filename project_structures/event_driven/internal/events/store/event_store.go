package store

import (
	"time"

	"github.com/edgardnogueira/go-patterns/project_structures/event_driven/internal/events/types"
)

// EventStore defines the interface for an event store
type EventStore interface {
	// SaveEvent saves an event to the event store
	SaveEvent(event types.Event) error

	// GetEvents retrieves events for a specific aggregate
	GetEvents(aggregateType, aggregateID string) ([]types.Event, error)

	// GetUnpublishedEvents retrieves unpublished events since a specific time
	GetUnpublishedEvents(since time.Time, limit int) ([]types.Event, error)

	// MarkEventAsPublished marks an event as published
	MarkEventAsPublished(eventID string) error

	// Close closes the connection to the event store
	Close() error
}
