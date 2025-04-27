package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/edgardnogueira/go-patterns/project_structures/event_driven/internal/events/publisher"
	"github.com/edgardnogueira/go-patterns/project_structures/event_driven/internal/infrastructure/config"
	"github.com/edgardnogueira/go-patterns/project_structures/event_driven/internal/infrastructure/eventstore/postgres"
	"github.com/edgardnogueira/go-patterns/project_structures/event_driven/internal/infrastructure/messaging/rabbitmq"
)

func main() {
	log.Println("Starting Event Publisher...")

	// Load configuration
	cfg, err := config.LoadConfig("config/publisher.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize event store
	eventStore, err := postgres.NewEventStore(cfg.Database.ConnectionString)
	if err != nil {
		log.Fatalf("Failed to initialize event store: %v", err)
	}
	defer eventStore.Close()

	// Initialize event publisher
	eventPublisher, err := rabbitmq.NewEventPublisher(cfg.MessageBroker.ConnectionString)
	if err != nil {
		log.Fatalf("Failed to initialize event publisher: %v", err)
	}
	defer eventPublisher.Close()

	// Initialize and start the outbox processor
	outboxProcessor := publisher.NewOutboxProcessor(eventStore, eventPublisher, cfg.Publisher.BatchSize)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the outbox processor in a goroutine
	go func() {
		outboxProcessor.Start(ctx, time.Duration(cfg.Publisher.IntervalSeconds)*time.Second)
	}()

	log.Println("Event Publisher is running")

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	log.Println("Shutting down...")
	// Cancel context to stop the outbox processor
	cancel()

	// Give some time for the outbox processor to finish ongoing work
	time.Sleep(2 * time.Second)
	log.Println("Event Publisher stopped gracefully")
}
