package publisher

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/edgardnogueira/go-patterns/project_structures/event_driven/internal/events/store"
	"github.com/edgardnogueira/go-patterns/project_structures/event_driven/internal/events/types"
)

// OutboxProcessor processes events from the outbox table and publishes them
type OutboxProcessor struct {
	eventStore    store.EventStore
	publisher     types.EventPublisher
	batchSize     int
	lastProcessed time.Time
}

// NewOutboxProcessor creates a new OutboxProcessor
func NewOutboxProcessor(eventStore store.EventStore, publisher types.EventPublisher, batchSize int) *OutboxProcessor {
	return &OutboxProcessor{
		eventStore:    eventStore,
		publisher:     publisher,
		batchSize:     batchSize,
		lastProcessed: time.Now().Add(-24 * time.Hour), // Start by processing events from the last 24 hours
	}
}

// Start begins the outbox processing loop
func (p *OutboxProcessor) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Println("Outbox processor started")

	// Process immediately on start
	p.processOutbox()

	for {
		select {
		case <-ticker.C:
			p.processOutbox()
		case <-ctx.Done():
			log.Println("Outbox processor stopping...")
			return
		}
	}
}

// processOutbox processes a batch of events from the outbox
func (p *OutboxProcessor) processOutbox() {
	// Get unpublished events from the event store
	events, err := p.eventStore.GetUnpublishedEvents(p.lastProcessed, p.batchSize)
	if err != nil {
		log.Printf("Error getting unpublished events: %v\n", err)
		return
	}

	if len(events) == 0 {
		return // Nothing to process
	}

	log.Printf("Processing %d events from outbox\n", len(events))

	// Update last processed time
	p.lastProcessed = events[len(events)-1].GetTimestamp()

	// Process each event
	for _, event := range events {
		if err := p.processEvent(event); err != nil {
			log.Printf("Error processing event %s: %v\n", event.GetEventID(), err)
			// Continue with other events even if one fails
		}
	}
}

// processEvent publishes a single event and marks it as published if successful
func (p *OutboxProcessor) processEvent(event types.Event) error {
	// Publish event
	err := p.publisher.PublishEvent(event)
	if err != nil {
		return fmt.Errorf("error publishing event: %w", err)
	}

	// Mark event as published
	if err := p.eventStore.MarkEventAsPublished(event.GetEventID()); err != nil {
		return fmt.Errorf("error marking event as published: %w", err)
	}

	log.Printf("Published event %s of type %s\n", event.GetEventID(), event.GetEventType())
	return nil
}
