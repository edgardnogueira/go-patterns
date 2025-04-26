package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

// EventType represents the type of event being published
type EventType string

// Event types used in the system
const (
	OrderCreated      EventType = "order.created"
	OrderConfirmed    EventType = "order.confirmed"
	OrderShipped      EventType = "order.shipped"
	InventoryReserved EventType = "inventory.reserved"
	InventoryReleased EventType = "inventory.released"
)

// Event represents a message published to the message bus
type Event struct {
	ID        string      `json:"id"`
	Type      EventType   `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// Publisher defines the interface for publishing events
type Publisher interface {
	Publish(eventType EventType, data interface{}) error
	Close() error
}

// Consumer defines the interface for consuming events
type Consumer interface {
	Subscribe(eventType EventType, handler func(event Event)) error
	Close() error
}

// NatsClient implements both Publisher and Consumer interfaces
type NatsClient struct {
	conn      *nats.Conn
	jetStream nats.JetStreamContext
}

// NewNatsClient creates a new NATS client
func NewNatsClient(url string) (*NatsClient, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	js, err := conn.JetStream()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	// Create the stream if it doesn't exist
	_, err = js.StreamInfo("EVENTS")
	if err != nil {
		// Create the stream
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     "EVENTS",
			Subjects: []string{"events.>"},
			Storage:  nats.FileStorage,
		})
		if err != nil {
			conn.Close()
			return nil, fmt.Errorf("failed to create stream: %w", err)
		}
	}

	return &NatsClient{
		conn:      conn,
		jetStream: js,
	}, nil
}

// Publish publishes an event to NATS
func (c *NatsClient) Publish(eventType EventType, data interface{}) error {
	event := Event{
		ID:        generateID(), // Implement your own ID generation function
		Type:      eventType,
		Data:      data,
		Timestamp: time.Now(),
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	subject := fmt.Sprintf("events.%s", event.Type)
	_, err = c.jetStream.Publish(subject, payload)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

// Subscribe subscribes to events of a specific type
func (c *NatsClient) Subscribe(eventType EventType, handler func(event Event)) error {
	subject := fmt.Sprintf("events.%s", eventType)
	
	// Create a durable consumer
	_, err := c.jetStream.Subscribe(subject, func(msg *nats.Msg) {
		var event Event
		err := json.Unmarshal(msg.Data, &event)
		if err != nil {
			log.Printf("Failed to unmarshal event: %v", err)
			return
		}

		handler(event)
		msg.Ack()
	}, nats.Durable(fmt.Sprintf("%s-consumer", eventType)))

	if err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}

	return nil
}

// Close closes the NATS connection
func (c *NatsClient) Close() error {
	c.conn.Close()
	return nil
}

// generateID generates a unique ID for events
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// CircuitBreaker implements a simple circuit breaker pattern
type CircuitBreaker struct {
	timeout     time.Duration
	maxFailures int
	failures    int
	lastFailure time.Time
	state       string // "closed", "open", "half-open"
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(timeout time.Duration, maxFailures int) *CircuitBreaker {
	return &CircuitBreaker{
		timeout:     timeout,
		maxFailures: maxFailures,
		failures:    0,
		state:       "closed",
	}
}

// Execute executes a function with circuit breaker protection
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
	if cb.state == "open" {
		// Check if timeout has elapsed to transition to half-open
		if time.Since(cb.lastFailure) > cb.timeout {
			cb.state = "half-open"
		} else {
			return fmt.Errorf("circuit breaker open")
		}
	}

	err := fn()
	if err != nil {
		cb.failures++
		cb.lastFailure = time.Now()

		if cb.failures >= cb.maxFailures || cb.state == "half-open" {
			cb.state = "open"
		}
		return err
	}

	// Success, reset if in half-open state
	if cb.state == "half-open" {
		cb.failures = 0
		cb.state = "closed"
	}
	return nil
}
