package eventsourcing_test

import (
	"testing"
	"time"

	"github.com/edgardnogueira/go-patterns/project_structures/event_driven/internal/events/types"
	"github.com/edgardnogueira/go-patterns/project_structures/event_driven/pkg/eventsourcing"
	"github.com/google/uuid"
)

// MockEvent implements the Event interface for testing
type MockEvent struct {
	types.BaseEvent
	Data string
}

// TestAggregateBasics tests basic aggregate functionality
func TestAggregateBasics(t *testing.T) {
	// Create a new aggregate
	aggID := uuid.New().String()
	agg := eventsourcing.NewAggregate(aggID, "test")

	// Check initial state
	if agg.ID != aggID {
		t.Errorf("Expected aggregate ID %s, got %s", aggID, agg.ID)
	}

	if agg.Type != "test" {
		t.Errorf("Expected aggregate type 'test', got %s", agg.Type)
	}

	if agg.Version != 0 {
		t.Errorf("Expected initial version 0, got %d", agg.Version)
	}

	if len(agg.Changes) != 0 {
		t.Errorf("Expected empty changes, got %d", len(agg.Changes))
	}
}

// TestAggregateApplyEvents tests applying events to an aggregate
func TestAggregateApplyEvents(t *testing.T) {
	// Create a new aggregate
	aggID := uuid.New().String()
	agg := eventsourcing.NewAggregate(aggID, "test")

	// Create mock events
	events := []types.Event{
		MockEvent{
			BaseEvent: types.BaseEvent{
				EventID:       uuid.New().String(),
				AggregateID:   aggID,
				AggregateType: "test",
				EventType:     "mock.created",
				Timestamp:     time.Now(),
				Version:       1,
			},
			Data: "event1",
		},
		MockEvent{
			BaseEvent: types.BaseEvent{
				EventID:       uuid.New().String(),
				AggregateID:   aggID,
				AggregateType: "test",
				EventType:     "mock.updated",
				Timestamp:     time.Now(),
				Version:       2,
			},
			Data: "event2",
		},
	}

	// Register event appliers
	applied := false
	agg.RegisterEventApplier("mock.created", func(event types.Event) {
		applied = true
	})

	// Apply events
	agg.ApplyEvents(events)

	// Check state after applying events
	if agg.Version != 2 {
		t.Errorf("Expected version 2 after applying events, got %d", agg.Version)
	}

	if !applied {
		t.Error("Event applier was not called")
	}
}

// TestProductAggregate tests the product aggregate specifically
func TestProductAggregate(t *testing.T) {
	// Create a new product aggregate
	productID := uuid.New().String()
	product := eventsourcing.NewProductAggregate(productID)

	// Check initial state
	if product.ID != productID {
		t.Errorf("Expected product ID %s, got %s", productID, product.ID)
	}

	if product.Stock != 0 {
		t.Errorf("Expected initial stock 0, got %d", product.Stock)
	}

	// Create a product
	err := product.CreateProduct("Test Product", "Test Description", 19.99, 10)
	if err != nil {
		t.Errorf("Unexpected error creating product: %v", err)
	}

	// Check state after creation
	if product.Name != "Test Product" {
		t.Errorf("Expected name 'Test Product', got %s", product.Name)
	}

	if product.Stock != 10 {
		t.Errorf("Expected stock 10, got %d", product.Stock)
	}

	if len(product.GetUncommittedChanges()) != 1 {
		t.Errorf("Expected 1 uncommitted change, got %d", len(product.GetUncommittedChanges()))
	}

	// Reduce stock
	err = product.ReduceStock(3)
	if err != nil {
		t.Errorf("Unexpected error reducing stock: %v", err)
	}

	// Check state after stock reduction
	if product.Stock != 7 {
		t.Errorf("Expected stock 7, got %d", product.Stock)
	}

	if len(product.GetUncommittedChanges()) != 2 {
		t.Errorf("Expected 2 uncommitted changes, got %d", len(product.GetUncommittedChanges()))
	}

	// Add stock
	err = product.AddStock(5)
	if err != nil {
		t.Errorf("Unexpected error adding stock: %v", err)
	}

	// Check state after adding stock
	if product.Stock != 12 {
		t.Errorf("Expected stock 12, got %d", product.Stock)
	}

	if len(product.GetUncommittedChanges()) != 3 {
		t.Errorf("Expected 3 uncommitted changes, got %d", len(product.GetUncommittedChanges()))
	}

	// Clear uncommitted changes
	product.ClearUncommittedChanges()

	if len(product.GetUncommittedChanges()) != 0 {
		t.Errorf("Expected 0 uncommitted changes after clearing, got %d", len(product.GetUncommittedChanges()))
	}
}
