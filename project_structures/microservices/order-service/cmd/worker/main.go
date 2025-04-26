package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/edgardnogueira/go-patterns/project_structures/microservices/order-service/internal/domain"
	"github.com/edgardnogueira/go-patterns/project_structures/microservices/order-service/internal/repositories"
	"github.com/edgardnogueira/go-patterns/project_structures/microservices/pkg/messaging"
	"github.com/edgardnogueira/go-patterns/project_structures/microservices/pkg/observability"
)

const serviceName = "order-service-worker"

var (
	postgresConnStr = flag.String("postgres-dsn", "postgres://postgres:postgres@postgres:5432/orders?sslmode=disable", "PostgreSQL connection string")
	natsURL         = flag.String("nats-url", "nats://nats:4222", "NATS server URL")
)

func main() {
	flag.Parse()

	// Override from environment if provided
	if envDSN := os.Getenv("POSTGRES_DSN"); envDSN != "" {
		*postgresConnStr = envDSN
	}
	if envNATS := os.Getenv("NATS_URL"); envNATS != "" {
		*natsURL = envNATS
	}

	// Initialize logger
	logger := observability.NewLogger(serviceName)
	logger.Info("Starting Order Service Worker")

	// Initialize metrics
	metricsServer := observability.NewMetricsServer(serviceName, 8091)
	if err := metricsServer.Start(); err != nil {
		logger.Fatal("Failed to start metrics server", err)
	}
	defer metricsServer.Stop()

	// Initialize tracing
	tp, err := observability.NewTracingProvider(serviceName)
	if err != nil {
		logger.Fatal("Failed to initialize tracing", err)
	}
	defer tp.Shutdown(context.Background())

	// Initialize database repository
	repo, err := repositories.NewPostgresOrderRepository(*postgresConnStr)
	if err != nil {
		logger.Fatal("Failed to connect to database", err)
	}
	defer repo.Close()

	// Initialize message consumer
	consumer, err := messaging.NewNatsClient(*natsURL)
	if err != nil {
		logger.Fatal("Failed to connect to NATS", err)
	}
	defer consumer.Close()

	// Subscribe to order-related events
	setupSubscriptions(consumer, repo, logger)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Worker shutting down...")
}

func setupSubscriptions(consumer messaging.Consumer, repo repositories.OrderRepository, logger *observability.Logger) {
	// Subscribe to inventory events
	err := consumer.Subscribe(messaging.InventoryReserved, func(event messaging.Event) {
		processInventoryReserved(event, repo, logger)
	})
	if err != nil {
		logger.Fatal("Failed to subscribe to InventoryReserved events", err)
	}

	err = consumer.Subscribe(messaging.InventoryReleased, func(event messaging.Event) {
		processInventoryReleased(event, repo, logger)
	})
	if err != nil {
		logger.Fatal("Failed to subscribe to InventoryReleased events", err)
	}

	// Subscribe to OrderCreated events for demo purposes
	// In a real system, the worker might process these differently
	err = consumer.Subscribe(messaging.OrderCreated, func(event messaging.Event) {
		logger.Info("Received OrderCreated event", map[string]interface{}{
			"event_id": event.ID,
		})
	})
	if err != nil {
		logger.Fatal("Failed to subscribe to OrderCreated events", err)
	}

	logger.Info("Subscribed to events")
}

func processInventoryReserved(event messaging.Event, repo repositories.OrderRepository, logger *observability.Logger) {
	ctx := context.Background()
	
	// Extract the order ID from the event data
	// In a real system, the event data structure would be more defined
	var data struct {
		OrderID string `json:"order_id"`
	}
	
	dataBytes, err := json.Marshal(event.Data)
	if err != nil {
		logger.Error("Failed to marshal event data", err)
		return
	}
	
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		logger.Error("Failed to unmarshal event data", err)
		return
	}
	
	if data.OrderID == "" {
		logger.Error("Missing order ID in event", nil)
		return
	}
	
	// Get order
	order, err := repo.GetOrderByID(ctx, data.OrderID)
	if err != nil {
		logger.Error("Failed to get order", err, map[string]interface{}{
			"order_id": data.OrderID,
		})
		return
	}
	
	// Confirm the order
	if err := order.ConfirmOrder(); err != nil {
		logger.Error("Failed to confirm order", err, map[string]interface{}{
			"order_id": data.OrderID,
		})
		return
	}
	
	// Update the order in the database
	if err := repo.UpdateOrder(ctx, order); err != nil {
		logger.Error("Failed to update order", err, map[string]interface{}{
			"order_id": data.OrderID,
		})
		return
	}
	
	logger.Info("Order confirmed after inventory reserved", map[string]interface{}{
		"order_id": data.OrderID,
	})
}

func processInventoryReleased(event messaging.Event, repo repositories.OrderRepository, logger *observability.Logger) {
	ctx := context.Background()
	
	// Extract the order ID from the event data
	var data struct {
		OrderID string `json:"order_id"`
	}
	
	dataBytes, err := json.Marshal(event.Data)
	if err != nil {
		logger.Error("Failed to marshal event data", err)
		return
	}
	
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		logger.Error("Failed to unmarshal event data", err)
		return
	}
	
	if data.OrderID == "" {
		logger.Error("Missing order ID in event", nil)
		return
	}
	
	// Get order
	order, err := repo.GetOrderByID(ctx, data.OrderID)
	if err != nil {
		logger.Error("Failed to get order", err, map[string]interface{}{
			"order_id": data.OrderID,
		})
		return
	}
	
	// Cancel the order
	if err := order.CancelOrder(); err != nil {
		logger.Error("Failed to cancel order", err, map[string]interface{}{
			"order_id": data.OrderID,
		})
		return
	}
	
	// Update the order in the database
	if err := repo.UpdateOrder(ctx, order); err != nil {
		logger.Error("Failed to update order", err, map[string]interface{}{
			"order_id": data.OrderID,
		})
		return
	}
	
	logger.Info("Order cancelled after inventory released", map[string]interface{}{
		"order_id": data.OrderID,
	})
}
