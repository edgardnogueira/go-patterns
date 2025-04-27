package eventsourcing

import (
	"errors"
	"fmt"

	"github.com/edgardnogueira/go-patterns/project_structures/event_driven/internal/events/types"
)

// ProductAggregate represents a product aggregate
type ProductAggregate struct {
	*Aggregate
	Name        string
	Description string
	Price       float64
	Stock       int
}

// NewProductAggregate creates a new ProductAggregate
func NewProductAggregate(id string) *ProductAggregate {
	aggregate := NewAggregate(id, "product")
	product := &ProductAggregate{
		Aggregate: aggregate,
	}

	// Register event appliers
	product.RegisterEventApplier("product.created", product.applyProductCreated)
	product.RegisterEventApplier("product.stock.reduced", product.applyStockReduced)
	product.RegisterEventApplier("product.stock.added", product.applyStockAdded)

	return product
}

// LoadFromHistory loads the aggregate state from its event history
func (p *ProductAggregate) LoadFromHistory(events []types.Event) {
	p.ApplyEvents(events)
}

// applyProductCreated applies a ProductCreatedEvent
func (p *ProductAggregate) applyProductCreated(event types.Event) {
	if productCreatedEvent, ok := event.(types.ProductCreatedEvent); ok {
		p.Name = productCreatedEvent.Name
		p.Description = productCreatedEvent.Description
		p.Price = productCreatedEvent.Price
		p.Stock = productCreatedEvent.Stock
	}
}

// applyStockReduced applies a ProductStockReducedEvent
func (p *ProductAggregate) applyStockReduced(event types.Event) {
	if stockReducedEvent, ok := event.(types.ProductStockReducedEvent); ok {
		p.Stock = stockReducedEvent.NewStock
	}
}

// applyStockAdded applies a ProductStockAddedEvent
func (p *ProductAggregate) applyStockAdded(event types.Event) {
	if stockAddedEvent, ok := event.(types.ProductStockAddedEvent); ok {
		p.Stock = stockAddedEvent.NewStock
	}
}

// CreateProduct creates a new product
func (p *ProductAggregate) CreateProduct(name, description string, price float64, initialStock int) error {
	// Validate
	if p.Version > 0 {
		return errors.New("product already exists")
	}

	if name == "" {
		return errors.New("product name cannot be empty")
	}

	if price < 0 {
		return errors.New("price cannot be negative")
	}

	if initialStock < 0 {
		return errors.New("initial stock cannot be negative")
	}

	// Add event
	p.AddEvent("product.created", map[string]interface{}{
		"id":          p.ID,
		"name":        name,
		"description": description,
		"price":       price,
		"stock":       initialStock,
	})

	return nil
}

// ReduceStock reduces the product stock
func (p *ProductAggregate) ReduceStock(quantity int) error {
	// Validate
	if p.Version == 0 {
		return errors.New("product does not exist")
	}

	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	if p.Stock < quantity {
		return fmt.Errorf("insufficient stock: have %d, requested %d", p.Stock, quantity)
	}

	// Calculate new stock
	newStock := p.Stock - quantity

	// Add event
	p.AddEvent("product.stock.reduced", map[string]interface{}{
		"product_id": p.ID,
		"quantity":   quantity,
		"new_stock":  newStock,
	})

	// Add low stock event if needed
	if newStock > 0 && newStock <= 5 {
		p.AddEvent("product.stock.low", map[string]interface{}{
			"product_id": p.ID,
			"stock":      newStock,
			"threshold":  5,
		})
	}

	// Add out of stock event if needed
	if newStock == 0 {
		p.AddEvent("product.stock.out", map[string]interface{}{
			"product_id": p.ID,
		})
	}

	return nil
}

// AddStock adds to the product stock
func (p *ProductAggregate) AddStock(quantity int) error {
	// Validate
	if p.Version == 0 {
		return errors.New("product does not exist")
	}

	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	// Calculate new stock
	newStock := p.Stock + quantity

	// Add event
	p.AddEvent("product.stock.added", map[string]interface{}{
		"product_id": p.ID,
		"quantity":   quantity,
		"new_stock":  newStock,
	})

	return nil
}
