package rabbitmq

import (
	"context"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// EventHandler is a function that processes an event
type EventHandler func([]byte) error

// RabbitMQEventConsumer implements event consumption from RabbitMQ
type RabbitMQEventConsumer struct {
	connection    *amqp.Connection
	channel       *amqp.Channel
	exchange      string
	queueName     string
	eventHandlers map[string]EventHandler
}

// NewEventConsumer creates a new RabbitMQEventConsumer
func NewEventConsumer(connectionString string) (*RabbitMQEventConsumer, error) {
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

	// Declare a queue
	queue, err := ch.QueueDeclare(
		"",    // name - empty for auto-generation
		false, // durable
		true,  // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return &RabbitMQEventConsumer{
		connection:    conn,
		channel:       ch,
		exchange:      exchange,
		queueName:     queue.Name,
		eventHandlers: make(map[string]EventHandler),
	}, nil
}

// Subscribe adds a handler for a specific event type
func (c *RabbitMQEventConsumer) Subscribe(eventType string, handler EventHandler) error {
	// Bind the queue to the exchange with the event type as routing key
	err := c.channel.QueueBind(
		c.queueName, // queue name
		eventType,   // routing key
		c.exchange,  // exchange
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue to exchange: %w", err)
	}

	// Store the handler
	c.eventHandlers[eventType] = handler
	log.Printf("Subscribed to event type: %s", eventType)

	return nil
}

// Start begins consuming events
func (c *RabbitMQEventConsumer) Start(ctx context.Context) error {
	// Start consuming from the queue
	msgs, err := c.channel.Consume(
		c.queueName, // queue
		"",         // consumer
		false,       // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	log.Println("Started consuming events")

	// Process messages
	for {
		select {
		case <-ctx.Done():
			log.Println("Context canceled, stopping consumer")
			return nil

		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("message channel closed")
			}

			// Get event type from routing key
			eventType := msg.RoutingKey

			// Find handler for this event type
			handler, exists := c.eventHandlers[eventType]
			if !exists {
				log.Printf("No handler registered for event type: %s, acknowledging", eventType)
				msg.Ack(false)
				continue
			}

			// Process the event
			log.Printf("Processing event type: %s", eventType)
			err := handler(msg.Body)
			if err != nil {
				log.Printf("Error processing event: %v, rejecting", err)
				// Reject and requeue the message
				msg.Reject(true)
				continue
			}

			// Acknowledge the message
			msg.Ack(false)
			log.Printf("Successfully processed event type: %s", eventType)
		}
	}
}

// Close closes the RabbitMQ connection and channel
func (c *RabbitMQEventConsumer) Close() error {
	// Close the channel
	if err := c.channel.Close(); err != nil {
		log.Printf("Error closing channel: %v", err)
	}

	// Close the connection
	if err := c.connection.Close(); err != nil {
		return fmt.Errorf("error closing connection: %w", err)
	}

	return nil
}
