package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/edgardnogueira/go-patterns/project_structures/event_driven/internal/events/types"
	"github.com/streadway/amqp"
)

// RabbitMQEventPublisher implements the EventPublisher interface using RabbitMQ
type RabbitMQEventPublisher struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	exchange   string
}

// NewEventPublisher creates a new RabbitMQEventPublisher
func NewEventPublisher(connectionString string) (types.EventPublisher, error) {
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Declare the exchange
	exchange := "events"
	err = ch.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	return &RabbitMQEventPublisher{
		connection: conn,
		channel:    ch,
		exchange:   exchange,
	}, nil
}

// PublishEvent publishes an event to RabbitMQ
func (p *RabbitMQEventPublisher) PublishEvent(event types.Event) error {
	// Marshal the event to JSON
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error marshaling event: %w", err)
	}

	// Create routing key from event type
	routingKey := event.GetEventType()

	// Publish the message
	err = p.channel.Publish(
		p.exchange, // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Headers: amqp.Table{
				"event_id":        event.GetEventID(),
				"aggregate_id":    event.GetAggregateID(),
				"aggregate_type":  event.GetAggregateType(),
				"event_type":      event.GetEventType(),
				"timestamp":       event.GetTimestamp(),
			},
		},
	)

	if err != nil {
		return fmt.Errorf("error publishing event: %w", err)
	}

	log.Printf("Published event %s with routing key %s", event.GetEventID(), routingKey)
	return nil
}

// Close closes the RabbitMQ connection and channel
func (p *RabbitMQEventPublisher) Close() error {
	// Close the channel
	if err := p.channel.Close(); err != nil {
		log.Printf("Error closing channel: %v", err)
	}

	// Close the connection
	if err := p.connection.Close(); err != nil {
		return fmt.Errorf("error closing connection: %w", err)
	}

	return nil
}
