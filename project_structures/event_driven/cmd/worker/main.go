package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/edgardnogueira/go-patterns/project_structures/event_driven/internal/events/consumer"
	"github.com/edgardnogueira/go-patterns/project_structures/event_driven/internal/infrastructure/config"
	"github.com/edgardnogueira/go-patterns/project_structures/event_driven/internal/infrastructure/messaging/rabbitmq"
	"github.com/edgardnogueira/go-patterns/project_structures/event_driven/internal/projections/builder"
)

func main() {
	log.Println("Starting Worker...")

	// Load configuration
	cfg, err := config.LoadConfig("config/worker.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize event consumer
	eventConsumer, err := rabbitmq.NewEventConsumer(cfg.MessageBroker.ConnectionString)
	if err != nil {
		log.Fatalf("Failed to initialize event consumer: %v", err)
	}
	defer eventConsumer.Close()

	// Initialize projection builders
	productProjectionBuilder := builder.NewProductProjectionBuilder(cfg.Database.ConnectionString)

	// Initialize event handlers
	productStockHandler := consumer.NewProductStockHandler(productProjectionBuilder)

	// Register event handlers with the consumer
	eventConsumer.Subscribe("product.stock.reduced", productStockHandler.HandleStockReduced)
	eventConsumer.Subscribe("product.stock.added", productStockHandler.HandleStockAdded)
	eventConsumer.Subscribe("product.created", productStockHandler.HandleProductCreated)

	// Start the event consumer in a goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := eventConsumer.Start(ctx); err != nil {
			log.Fatalf("Failed to start event consumer: %v", err)
		}
	}()

	log.Println("Worker is running")

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	log.Println("Shutting down...")
	
	// Cancel context to stop the event consumer
	cancel()

	// Give some time for the event consumer to finish processing
	time.Sleep(2 * time.Second)
	log.Println("Worker stopped gracefully")
}
