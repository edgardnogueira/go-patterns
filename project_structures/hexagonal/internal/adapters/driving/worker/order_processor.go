package worker

import (
	"context"
	"fmt"
	"log"
	"time"
	
	"github.com/edgardnogueira/go-patterns/project_structures/hexagonal/internal/domain/model"
	"github.com/edgardnogueira/go-patterns/project_structures/hexagonal/internal/domain/service"
)

// OrderProcessor is a worker that processes orders
type OrderProcessor struct {
	orderService *service.OrderService
	logger       *log.Logger
	interval     time.Duration
	stopCh       chan struct{}
}

// NewOrderProcessor creates a new order processor worker
func NewOrderProcessor(orderService *service.OrderService, logger *log.Logger, interval time.Duration) *OrderProcessor {
	if logger == nil {
		logger = log.Default()
	}
	
	if interval <= 0 {
		interval = 30 * time.Second
	}
	
	return &OrderProcessor{
		orderService: orderService,
		logger:       logger,
		interval:     interval,
		stopCh:       make(chan struct{}),
	}
}

// Start begins processing orders at regular intervals
func (p *OrderProcessor) Start(ctx context.Context) {
	p.logger.Println("Starting order processor worker...")
	
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			p.processPendingOrders(ctx)
		case <-p.stopCh:
			p.logger.Println("Order processor worker stopped")
			return
		case <-ctx.Done():
			p.logger.Println("Order processor worker context cancelled")
			return
		}
	}
}

// Stop halts the order processor
func (p *OrderProcessor) Stop() {
	p.logger.Println("Stopping order processor worker...")
	close(p.stopCh)
}

// ProcessPendingOrders processes all orders with "created" status
func (p *OrderProcessor) processPendingOrders(ctx context.Context) {
	p.logger.Println("Processing pending orders...")
	
	// Find all orders with "created" status
	orders, err := p.orderService.FindOrdersByStatus(ctx, model.OrderStatusCreated)
	if err != nil {
		p.logger.Printf("Error finding pending orders: %v\n", err)
		return
	}
	
	p.logger.Printf("Found %d pending orders to process\n", len(orders))
	
	// Process each order
	for _, order := range orders {
		if err := p.processOrder(ctx, order.ID); err != nil {
			p.logger.Printf("Error processing order %s: %v\n", order.ID, err)
		} else {
			p.logger.Printf("Successfully processed order %s\n", order.ID)
		}
		
		// Add a small delay between processing orders
		time.Sleep(100 * time.Millisecond)
	}
}

// ProcessOrder processes a specific order (moves it to "processing" status)
func (p *OrderProcessor) processOrder(ctx context.Context, orderID string) error {
	p.logger.Printf("Processing order %s...\n", orderID)
	
	// Retrieve the order
	order, err := p.orderService.GetOrder(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}
	
	// Skip orders that are not in "created" status
	if order.Status != model.OrderStatusCreated {
		return fmt.Errorf("order %s is not in 'created' status (current status: %s)", orderID, order.Status)
	}
	
	// Update the status to "processing"
	if err := p.orderService.UpdateOrderStatus(ctx, orderID, model.OrderStatusProcessing); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	
	// Simulate processing time
	time.Sleep(500 * time.Millisecond)
	
	return nil
}
