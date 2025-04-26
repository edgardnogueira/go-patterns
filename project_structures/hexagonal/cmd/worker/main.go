package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	"github.com/edgardnogueira/go-patterns/project_structures/hexagonal/internal/adapters/driven/database"
	"github.com/edgardnogueira/go-patterns/project_structures/hexagonal/internal/adapters/driven/external"
	"github.com/edgardnogueira/go-patterns/project_structures/hexagonal/internal/adapters/driving/worker"
	"github.com/edgardnogueira/go-patterns/project_structures/hexagonal/internal/domain/service"
)

func main() {
	// Setup logger
	logger := log.New(os.Stdout, "[WORKER] ", log.LstdFlags)
	logger.Println("Starting order processing worker...")
	
	// Setup context that will be canceled on shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// Setup ports & adapters
	orderRepo := database.NewMemoryOrderRepository()
	notificationService := external.NewLogNotificationService(logger)
	orderService := service.NewOrderService(orderRepo, notificationService)
	
	// Create and start the worker
	processor := worker.NewOrderProcessor(orderService, logger, 10*time.Second)
	
	// Start the worker in a separate goroutine
	go processor.Start(ctx)
	
	// Add some example data if needed (for demonstration)
	addExampleData(ctx, orderService, logger)
	
	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	
	logger.Println("Shutting down worker...")
	
	// Stop the worker
	processor.Stop()
	
	// Cancel the context to signal shutdown to all components
	cancel()
	
	logger.Println("Worker gracefully stopped")
}

// addExampleData creates some example orders for demonstration purposes
func addExampleData(ctx context.Context, orderService *service.OrderService, logger *log.Logger) {
	logger.Println("Adding example data...")
	
	// Example items
	items1 := []model.OrderItem{
		{
			ProductID:  "prod-001",
			Name:       "Smartphone",
			Quantity:   1,
			UnitPrice:  799.99,
		},
		{
			ProductID:  "prod-002",
			Name:       "Phone Case",
			Quantity:   1,
			UnitPrice:  24.99,
		},
	}
	
	items2 := []model.OrderItem{
		{
			ProductID:  "prod-003",
			Name:       "Laptop",
			Quantity:   1,
			UnitPrice:  1299.99,
		},
		{
			ProductID:  "prod-004",
			Name:       "Laptop Bag",
			Quantity:   1,
			UnitPrice:  59.99,
		},
		{
			ProductID:  "prod-005",
			Name:       "Wireless Mouse",
			Quantity:   1,
			UnitPrice:  29.99,
		},
	}
	
	// Create example orders
	_, err := orderService.CreateOrder(ctx, "customer-001", items1, "123 Main St, City, Country")
	if err != nil {
		logger.Printf("Error creating example order 1: %v\n", err)
	}
	
	_, err = orderService.CreateOrder(ctx, "customer-002", items2, "456 Oak Ave, Town, Country")
	if err != nil {
		logger.Printf("Error creating example order 2: %v\n", err)
	}
	
	logger.Println("Example data added successfully")
}
