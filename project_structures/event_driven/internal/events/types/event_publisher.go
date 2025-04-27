package types

// EventPublisher defines the interface for publishing events
type EventPublisher interface {
	// PublishEvent publishes an event to the message broker
	PublishEvent(event Event) error
	
	// Close closes the connection to the message broker
	Close() error
}
