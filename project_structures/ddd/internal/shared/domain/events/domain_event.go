package events

import (
	"time"
)

// DomainEvent is the base interface for all domain events
type DomainEvent interface {
	// EventID returns the unique identifier of the event
	EventID() string
	
	// EventType returns the type of the event
	EventType() string
	
	// AggregateID returns the identifier of the aggregate that produced the event
	AggregateID() string
	
	// AggregateType returns the type of the aggregate that produced the event
	AggregateType() string
	
	// OccurredAt returns when the event occurred
	OccurredAt() time.Time
	
	// Payload returns the event payload as a map
	Payload() map[string]interface{}
}

// BaseDomainEvent implements the common functionality for all domain events
type BaseDomainEvent struct {
	eventID      string
	eventType    string
	aggregateID  string
	aggregateType string
	occurredAt   time.Time
	payload      map[string]interface{}
}

// NewBaseDomainEvent creates a new BaseDomainEvent
func NewBaseDomainEvent(
	eventID string,
	eventType string,
	aggregateID string,
	aggregateType string,
	payload map[string]interface{},
) BaseDomainEvent {
	return BaseDomainEvent{
		eventID:       eventID,
		eventType:     eventType,
		aggregateID:   aggregateID,
		aggregateType: aggregateType,
		occurredAt:    time.Now(),
		payload:       payload,
	}
}

// EventID returns the unique identifier of the event
func (e BaseDomainEvent) EventID() string {
	return e.eventID
}

// EventType returns the type of the event
func (e BaseDomainEvent) EventType() string {
	return e.eventType
}

// AggregateID returns the identifier of the aggregate that produced the event
func (e BaseDomainEvent) AggregateID() string {
	return e.aggregateID
}

// AggregateType returns the type of the aggregate that produced the event
func (e BaseDomainEvent) AggregateType() string {
	return e.aggregateType
}

// OccurredAt returns when the event occurred
func (e BaseDomainEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// Payload returns the event payload as a map
func (e BaseDomainEvent) Payload() map[string]interface{} {
	return e.payload
}
