package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/edgardnogueira/go-patterns/project_structures/event_driven/internal/events/store"
	"github.com/edgardnogueira/go-patterns/project_structures/event_driven/internal/events/types"
	_ "github.com/lib/pq"
)

// PostgresEventStore implements the EventStore interface using PostgreSQL
type PostgresEventStore struct {
	db *sql.DB
}

// NewEventStore creates a new PostgresEventStore
func NewEventStore(connectionString string) (store.EventStore, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Ensure tables exist
	if err = createTablesIfNotExist(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &PostgresEventStore{db: db}, nil
}

// createTablesIfNotExist creates the necessary tables if they don't exist
func createTablesIfNotExist(db *sql.DB) error {
	// Create events table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS events (
			event_id VARCHAR(36) PRIMARY KEY,
			aggregate_id VARCHAR(36) NOT NULL,
			aggregate_type VARCHAR(50) NOT NULL,
			event_type VARCHAR(100) NOT NULL,
			version INT NOT NULL,
			payload JSONB NOT NULL,
			timestamp TIMESTAMP NOT NULL,
			published BOOLEAN DEFAULT FALSE,
			CREATE INDEX IF NOT EXISTS idx_events_aggregate ON events(aggregate_type, aggregate_id),
			CREATE INDEX IF NOT EXISTS idx_events_published ON events(published, timestamp)
		)
	`)

	return err
}

// SaveEvent stores an event in the database
func (s *PostgresEventStore) SaveEvent(event types.Event) error {
	// Serialize the event to JSON
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Insert into database
	_, err = s.db.Exec(
		`INSERT INTO events 
		(event_id, aggregate_id, aggregate_type, event_type, version, payload, timestamp, published) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, false)`,
		event.GetEventID(),
		event.GetAggregateID(),
		event.GetAggregateType(),
		event.GetEventType(),
		event.GetVersion(),
		payload,
		event.GetTimestamp(),
	)

	if err != nil {
		return fmt.Errorf("failed to insert event: %w", err)
	}

	return nil
}

// GetEvents retrieves all events for a specific aggregate
func (s *PostgresEventStore) GetEvents(aggregateType, aggregateID string) ([]types.Event, error) {
	rows, err := s.db.Query(
		`SELECT event_type, payload FROM events 
		WHERE aggregate_type = $1 AND aggregate_id = $2 
		ORDER BY version ASC`,
		aggregateType, aggregateID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	var events []types.Event
	for rows.Next() {
		var eventType string
		var payload []byte

		if err := rows.Scan(&eventType, &payload); err != nil {
			return nil, fmt.Errorf("failed to scan event row: %w", err)
		}

		// Deserialize the event based on its type
		event, err := deserializeEvent(eventType, payload)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize event: %w", err)
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating event rows: %w", err)
	}

	return events, nil
}

// GetUnpublishedEvents retrieves unpublished events since a specific time
func (s *PostgresEventStore) GetUnpublishedEvents(since time.Time, limit int) ([]types.Event, error) {
	rows, err := s.db.Query(
		`SELECT event_type, payload FROM events 
		WHERE published = false AND timestamp > $1 
		ORDER BY timestamp ASC 
		LIMIT $2`,
		since, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query unpublished events: %w", err)
	}
	defer rows.Close()

	var events []types.Event
	for rows.Next() {
		var eventType string
		var payload []byte

		if err := rows.Scan(&eventType, &payload); err != nil {
			return nil, fmt.Errorf("failed to scan event row: %w", err)
		}

		// Deserialize the event based on its type
		event, err := deserializeEvent(eventType, payload)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize event: %w", err)
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating event rows: %w", err)
	}

	return events, nil
}

// MarkEventAsPublished marks an event as published
func (s *PostgresEventStore) MarkEventAsPublished(eventID string) error {
	result, err := s.db.Exec(
		"UPDATE events SET published = true WHERE event_id = $1",
		eventID,
	)
	if err != nil {
		return fmt.Errorf("failed to mark event as published: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("event not found")
	}

	return nil
}

// Close closes the database connection
func (s *PostgresEventStore) Close() error {
	return s.db.Close()
}

// deserializeEvent deserializes an event based on its type
func deserializeEvent(eventType string, payload []byte) (types.Event, error) {
	var event types.Event

	switch eventType {
	case "product.created":
		var e types.ProductCreatedEvent
		if err := json.Unmarshal(payload, &e); err != nil {
			return nil, err
		}
		event = e

	case "product.stock.reduced":
		var e types.ProductStockReducedEvent
		if err := json.Unmarshal(payload, &e); err != nil {
			return nil, err
		}
		event = e

	case "product.stock.added":
		var e types.ProductStockAddedEvent
		if err := json.Unmarshal(payload, &e); err != nil {
			return nil, err
		}
		event = e

	case "product.stock.out":
		var e types.ProductOutOfStockEvent
		if err := json.Unmarshal(payload, &e); err != nil {
			return nil, err
		}
		event = e

	case "product.stock.low":
		var e types.ProductLowStockEvent
		if err := json.Unmarshal(payload, &e); err != nil {
			return nil, err
		}
		event = e

	default:
		return nil, fmt.Errorf("unknown event type: %s", eventType)
	}

	return event, nil
}
