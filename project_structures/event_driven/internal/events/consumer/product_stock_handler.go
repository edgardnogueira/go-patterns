package consumer

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/edgardnogueira/go-patterns/project_structures/event_driven/internal/events/types"
	"github.com/edgardnogueira/go-patterns/project_structures/event_driven/internal/projections/builder"
)

// ProductStockHandler handles product stock related events
type ProductStockHandler struct {
	projectionBuilder *builder.ProductProjectionBuilder
}

// NewProductStockHandler creates a new ProductStockHandler
func NewProductStockHandler(projectionBuilder *builder.ProductProjectionBuilder) *ProductStockHandler {
	return &ProductStockHandler{
		projectionBuilder: projectionBuilder,
	}
}

// HandleProductCreated handles the ProductCreatedEvent
func (h *ProductStockHandler) HandleProductCreated(eventData []byte) error {
	var event types.ProductCreatedEvent
	if err := json.Unmarshal(eventData, &event); err != nil {
		return fmt.Errorf("error unmarshaling product created event: %w", err)
	}

	log.Printf("Processing product created event for product ID %s\n", event.ID)

	// Update projection
	if err := h.projectionBuilder.CreateProduct(event); err != nil {
		return fmt.Errorf("error updating product projection: %w", err)
	}

	log.Printf("Successfully processed product created event for product ID %s\n", event.ID)
	return nil
}

// HandleStockReduced handles the ProductStockReducedEvent
func (h *ProductStockHandler) HandleStockReduced(eventData []byte) error {
	var event types.ProductStockReducedEvent
	if err := json.Unmarshal(eventData, &event); err != nil {
		return fmt.Errorf("error unmarshaling stock reduced event: %w", err)
	}

	log.Printf("Processing stock reduced event for product ID %s\n", event.ProductID)

	// Update projection
	if err := h.projectionBuilder.UpdateProductStock(event.ProductID, event.NewStock); err != nil {
		return fmt.Errorf("error updating product stock projection: %w", err)
	}

	log.Printf("Successfully processed stock reduced event for product ID %s\n", event.ProductID)
	return nil
}

// HandleStockAdded handles the ProductStockAddedEvent
func (h *ProductStockHandler) HandleStockAdded(eventData []byte) error {
	var event types.ProductStockAddedEvent
	if err := json.Unmarshal(eventData, &event); err != nil {
		return fmt.Errorf("error unmarshaling stock added event: %w", err)
	}

	log.Printf("Processing stock added event for product ID %s\n", event.ProductID)

	// Update projection
	if err := h.projectionBuilder.UpdateProductStock(event.ProductID, event.NewStock); err != nil {
		return fmt.Errorf("error updating product stock projection: %w", err)
	}

	log.Printf("Successfully processed stock added event for product ID %s\n", event.ProductID)
	return nil
}

// HandleOutOfStock handles the ProductOutOfStockEvent
func (h *ProductStockHandler) HandleOutOfStock(eventData []byte) error {
	var event types.ProductOutOfStockEvent
	if err := json.Unmarshal(eventData, &event); err != nil {
		return fmt.Errorf("error unmarshaling out of stock event: %w", err)
	}

	log.Printf("Processing out of stock event for product ID %s\n", event.ProductID)

	// Here we could trigger notifications or other business processes
	log.Printf("Product %s is out of stock! Notification would be sent here.\n", event.ProductID)

	return nil
}

// HandleLowStock handles the ProductLowStockEvent
func (h *ProductStockHandler) HandleLowStock(eventData []byte) error {
	var event types.ProductLowStockEvent
	if err := json.Unmarshal(eventData, &event); err != nil {
		return fmt.Errorf("error unmarshaling low stock event: %w", err)
	}

	log.Printf("Processing low stock event for product ID %s\n", event.ProductID)

	// Here we could trigger reorder processes or notifications
	log.Printf("Product %s has low stock (%d)! Reorder notification would be sent here.\n", 
		event.ProductID, event.Stock)

	return nil
}
