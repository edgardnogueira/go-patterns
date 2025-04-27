package handler

import (
	"encoding/json"
	"fmt"
	"github.com/edgardnogueira/go-patterns/project_structures/event_driven/internal/commands/model"
	"github.com/edgardnogueira/go-patterns/project_structures/event_driven/internal/events/store"
	"github.com/edgardnogueira/go-patterns/project_structures/event_driven/internal/events/types"
	"github.com/google/uuid"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"net/http"
	"time"
)

// ProductCommandHandler handles commands related to products
type ProductCommandHandler struct {
	eventStore    store.EventStore
	eventPublisher types.EventPublisher
}

// NewProductCommandHandler creates a new ProductCommandHandler
func NewProductCommandHandler(eventStore store.EventStore, eventPublisher types.EventPublisher) *ProductCommandHandler {
	return &ProductCommandHandler{
		eventStore:    eventStore,
		eventPublisher: eventPublisher,
	}
}

// HTTPHandleCreateProduct handles HTTP requests to create a product
func (h *ProductCommandHandler) HTTPHandleCreateProduct(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	var cmd model.CreateProductCommand
	if err := json.Unmarshal(body, &cmd); err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	// Set command metadata
	cmd.CommandID = uuid.New().String()
	cmd.Timestamp = time.Now()

	result, err := h.HandleCreateProduct(cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// HTTPHandleReduceStock handles HTTP requests to reduce product stock
func (h *ProductCommandHandler) HTTPHandleReduceStock(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	var cmd model.ReduceProductStockCommand
	if err := json.Unmarshal(body, &cmd); err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	// Get product ID from URL
	productID := chi.URLParam(r, "id")
	cmd.ProductID = productID

	// Set command metadata
	cmd.CommandID = uuid.New().String()
	cmd.Timestamp = time.Now()

	result, err := h.HandleReduceProductStock(cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// HTTPHandleAddStock handles HTTP requests to add product stock
func (h *ProductCommandHandler) HTTPHandleAddStock(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	var cmd model.AddProductStockCommand
	if err := json.Unmarshal(body, &cmd); err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	// Get product ID from URL
	productID := chi.URLParam(r, "id")
	cmd.ProductID = productID

	// Set command metadata
	cmd.CommandID = uuid.New().String()
	cmd.Timestamp = time.Now()

	result, err := h.HandleAddProductStock(cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// HandleCreateProduct handles the CreateProductCommand
func (h *ProductCommandHandler) HandleCreateProduct(cmd model.CreateProductCommand) (interface{}, error) {
	// Validate command
	if err := cmd.Validate(); err != nil {
		return nil, err
	}

	// Check if product already exists
	events, err := h.eventStore.GetEvents("product", cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("error checking product existence: %w", err)
	}

	if len(events) > 0 {
		return nil, fmt.Errorf("product with ID %s already exists", cmd.ID)
	}

	// Create event
	eventID := uuid.New().String()
	event := types.ProductCreatedEvent{
		BaseEvent: types.BaseEvent{
			EventID:       eventID,
			AggregateID:   cmd.ID,
			AggregateType: "product",
			EventType:     "product.created",
			Timestamp:     time.Now(),
			Version:       1,
		},
		ID:          cmd.ID,
		Name:        cmd.Name,
		Description: cmd.Description,
		Price:       cmd.Price,
		Stock:       cmd.InitialStock,
	}

	// Store event
	if err := h.eventStore.SaveEvent(event); err != nil {
		return nil, fmt.Errorf("error saving product created event: %w", err)
	}

	// Return result
	return map[string]interface{}{
		"id":          cmd.ID,
		"name":        cmd.Name,
		"description": cmd.Description,
		"price":       cmd.Price,
		"stock":       cmd.InitialStock,
	}, nil
}

// HandleReduceProductStock handles the ReduceProductStockCommand
func (h *ProductCommandHandler) HandleReduceProductStock(cmd model.ReduceProductStockCommand) (interface{}, error) {
	// Validate command
	if err := cmd.Validate(); err != nil {
		return nil, err
	}

	// Get product events
	events, err := h.eventStore.GetEvents("product", cmd.ProductID)
	if err != nil {
		return nil, fmt.Errorf("error getting product events: %w", err)
	}

	if len(events) == 0 {
		return nil, fmt.Errorf("product with ID %s not found", cmd.ProductID)
	}

	// Reconstruct product state
	var product struct {
		ID    string
		Stock int
	}

	product.ID = cmd.ProductID
	product.Stock = 0

	// Apply events to reconstruct current state
	for _, event := range events {
		switch e := event.(type) {
		case types.ProductCreatedEvent:
			product.Stock = e.Stock
		case types.ProductStockReducedEvent:
			product.Stock = e.NewStock
		case types.ProductStockAddedEvent:
			product.Stock = e.NewStock
		}
	}

	// Check if there's enough stock
	if product.Stock < cmd.Quantity {
		return nil, fmt.Errorf("insufficient stock: have %d, requested %d", product.Stock, cmd.Quantity)
	}

	// Calculate new stock level
	newStock := product.Stock - cmd.Quantity

	// Create event
	eventID := uuid.New().String()
	version := len(events) + 1

	event := types.ProductStockReducedEvent{
		BaseEvent: types.BaseEvent{
			EventID:       eventID,
			AggregateID:   cmd.ProductID,
			AggregateType: "product",
			EventType:     "product.stock.reduced",
			Timestamp:     time.Now(),
			Version:       version,
		},
		ProductID: cmd.ProductID,
		Quantity:  cmd.Quantity,
		NewStock:  newStock,
	}

	// Store event
	if err := h.eventStore.SaveEvent(event); err != nil {
		return nil, fmt.Errorf("error saving stock reduced event: %w", err)
	}

	// Create and store low stock event if needed (threshold of 5)
	if newStock > 0 && newStock <= 5 {
		lowStockEvent := types.ProductLowStockEvent{
			BaseEvent: types.BaseEvent{
				EventID:       uuid.New().String(),
				AggregateID:   cmd.ProductID,
				AggregateType: "product",
				EventType:     "product.stock.low",
				Timestamp:     time.Now(),
				Version:       version + 1,
			},
			ProductID: cmd.ProductID,
			Stock:     newStock,
			Threshold: 5,
		}
		
		if err := h.eventStore.SaveEvent(lowStockEvent); err != nil {
			// Just log this error, don't fail the main operation
			fmt.Printf("error saving low stock event: %v\n", err)
		}
	}

	// Create and store out of stock event if needed
	if newStock == 0 {
		outOfStockEvent := types.ProductOutOfStockEvent{
			BaseEvent: types.BaseEvent{
				EventID:       uuid.New().String(),
				AggregateID:   cmd.ProductID,
				AggregateType: "product",
				EventType:     "product.stock.out",
				Timestamp:     time.Now(),
				Version:       version + 1,
			},
			ProductID: cmd.ProductID,
		}
		
		if err := h.eventStore.SaveEvent(outOfStockEvent); err != nil {
			// Just log this error, don't fail the main operation
			fmt.Printf("error saving out of stock event: %v\n", err)
		}
	}

	// Return result
	return map[string]interface{}{
		"product_id": cmd.ProductID,
		"quantity":   cmd.Quantity,
		"new_stock":  newStock,
	}, nil
}

// HandleAddProductStock handles the AddProductStockCommand
func (h *ProductCommandHandler) HandleAddProductStock(cmd model.AddProductStockCommand) (interface{}, error) {
	// Validate command
	if err := cmd.Validate(); err != nil {
		return nil, err
	}

	// Get product events
	events, err := h.eventStore.GetEvents("product", cmd.ProductID)
	if err != nil {
		return nil, fmt.Errorf("error getting product events: %w", err)
	}

	if len(events) == 0 {
		return nil, fmt.Errorf("product with ID %s not found", cmd.ProductID)
	}

	// Reconstruct product state
	var product struct {
		ID    string
		Stock int
	}

	product.ID = cmd.ProductID
	product.Stock = 0

	// Apply events to reconstruct current state
	for _, event := range events {
		switch e := event.(type) {
		case types.ProductCreatedEvent:
			product.Stock = e.Stock
		case types.ProductStockReducedEvent:
			product.Stock = e.NewStock
		case types.ProductStockAddedEvent:
			product.Stock = e.NewStock
		}
	}

	// Calculate new stock level
	newStock := product.Stock + cmd.Quantity

	// Create event
	eventID := uuid.New().String()
	version := len(events) + 1

	event := types.ProductStockAddedEvent{
		BaseEvent: types.BaseEvent{
			EventID:       eventID,
			AggregateID:   cmd.ProductID,
			AggregateType: "product",
			EventType:     "product.stock.added",
			Timestamp:     time.Now(),
			Version:       version,
		},
		ProductID: cmd.ProductID,
		Quantity:  cmd.Quantity,
		NewStock:  newStock,
	}

	// Store event
	if err := h.eventStore.SaveEvent(event); err != nil {
		return nil, fmt.Errorf("error saving stock added event: %w", err)
	}

	// Return result
	return map[string]interface{}{
		"product_id": cmd.ProductID,
		"quantity":   cmd.Quantity,
		"new_stock":  newStock,
	}, nil
}
