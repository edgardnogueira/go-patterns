package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/edgardnogueira/go-patterns/project_structures/microservices/inventory-service/internal/domain"
	"github.com/edgardnogueira/go-patterns/project_structures/microservices/inventory-service/internal/repositories"
	"github.com/edgardnogueira/go-patterns/project_structures/microservices/pkg/messaging"
	"github.com/edgardnogueira/go-patterns/project_structures/microservices/pkg/observability"
)

const serviceName = "inventory-service-worker"

var (
	postgresConnStr = flag.String("postgres-dsn", "postgres://postgres:postgres@postgres:5432/inventory?sslmode=disable", "PostgreSQL connection string")
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
	logger.Info("Starting Inventory Service Worker")

	// Initialize metrics
	metricsServer := observability.NewMetricsServer(serviceName, 8092)
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
	repo, err := repositories.NewPostgresInventoryRepository(*postgresConnStr)
	if err != nil {
		logger.Fatal("Failed to connect to database", err)
	}
	defer repo.Close()

	// Initialize message consumer and publisher
	client, err := messaging.NewNatsClient(*natsURL)
	if err != nil {
		logger.Fatal("Failed to connect to NATS", err)
	}
	defer client.Close()

	// Subscribe to order-related events
	setupSubscriptions(client, repo, logger)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Worker shutting down...")
}

func setupSubscriptions(client *messaging.NatsClient, repo repositories.InventoryRepository, logger *observability.Logger) {
	// Subscribe to OrderCreated events
	err := client.Subscribe(messaging.OrderCreated, func(event messaging.Event) {
		processOrderCreated(event, repo, client, logger)
	})
	if err != nil {
		logger.Fatal("Failed to subscribe to OrderCreated events", err)
	}

	// Subscribe to OrderCancelled events
	err = client.Subscribe(messaging.OrderCancelled, func(event messaging.Event) {
		processOrderCancelled(event, repo, client, logger)
	})
	if err != nil {
		logger.Fatal("Failed to subscribe to OrderCancelled events", err)
	}

	logger.Info("Subscribed to events")
}

func processOrderCreated(event messaging.Event, repo repositories.InventoryRepository, publisher messaging.Publisher, logger *observability.Logger) {
	ctx := context.Background()
	
	// Extract the order data
	var order struct {
		ID    string `json:"id"`
		Items []struct {
			ProductID string `json:"product_id"`
			Quantity  int    `json:"quantity"`
		} `json:"items"`
	}
	
	dataBytes, err := json.Marshal(event.Data)
	if err != nil {
		logger.Error("Failed to marshal event data", err)
		return
	}
	
	if err := json.Unmarshal(dataBytes, &order); err != nil {
		logger.Error("Failed to unmarshal event data", err)
		return
	}
	
	if order.ID == "" || len(order.Items) == 0 {
		logger.Error("Invalid order data in event", nil)
		return
	}
	
	// Process each item in the order
	for _, item := range order.Items {
		// Get product
		product, err := repo.GetProductByID(ctx, item.ProductID)
		if err != nil {
			logger.Error("Failed to get product", err, map[string]interface{}{
				"product_id": item.ProductID,
				"order_id":   order.ID,
			})
			
			// Release any reservations already made for this order
			releaseReservationsForOrder(ctx, order.ID, repo, publisher, logger)
			return
		}
		
		// Check if there's enough inventory
		if !product.IsAvailable(item.Quantity) {
			logger.Error("Insufficient inventory", nil, map[string]interface{}{
				"product_id": item.ProductID,
				"order_id":   order.ID,
				"required":   item.Quantity,
				"available":  product.Quantity,
			})
			
			// Release any reservations already made for this order
			releaseReservationsForOrder(ctx, order.ID, repo, publisher, logger)
			return
		}
		
		// Reserve stock
		reservation, err := product.ReserveStock(order.ID, item.Quantity)
		if err != nil {
			logger.Error("Failed to reserve stock", err, map[string]interface{}{
				"product_id": item.ProductID,
				"order_id":   order.ID,
			})
			
			// Release any reservations already made for this order
			releaseReservationsForOrder(ctx, order.ID, repo, publisher, logger)
			return
		}
		
		// Save reservation
		if err := repo.CreateReservation(ctx, reservation); err != nil {
			logger.Error("Failed to save reservation", err, map[string]interface{}{
				"product_id": item.ProductID,
				"order_id":   order.ID,
			})
			
			// Release any reservations already made for this order
			releaseReservationsForOrder(ctx, order.ID, repo, publisher, logger)
			return
		}
		
		// Update product stock
		if err := repo.UpdateProduct(ctx, product); err != nil {
			logger.Error("Failed to update product stock", err, map[string]interface{}{
				"product_id": item.ProductID,
				"order_id":   order.ID,
			})
			
			// Release any reservations already made for this order
			releaseReservationsForOrder(ctx, order.ID, repo, publisher, logger)
			return
		}
	}
	
	// All items processed successfully
	// Publish inventory reserved event
	err = publisher.Publish(messaging.InventoryReserved, map[string]interface{}{
		"order_id": order.ID,
	})
	if err != nil {
		logger.Error("Failed to publish inventory reserved event", err, map[string]interface{}{
			"order_id": order.ID,
		})
		// Continue even if event publishing fails
	}
	
	logger.Info("Inventory reserved for order", map[string]interface{}{
		"order_id": order.ID,
	})
}

func processOrderCancelled(event messaging.Event, repo repositories.InventoryRepository, publisher messaging.Publisher, logger *observability.Logger) {
	ctx := context.Background()
	
	// Extract the order ID
	var order struct {
		ID string `json:"id"`
	}
	
	dataBytes, err := json.Marshal(event.Data)
	if err != nil {
		logger.Error("Failed to marshal event data", err)
		return
	}
	
	if err := json.Unmarshal(dataBytes, &order); err != nil {
		logger.Error("Failed to unmarshal event data", err)
		return
	}
	
	if order.ID == "" {
		logger.Error("Missing order ID in event", nil)
		return
	}
	
	// Release all reservations for this order
	if err := releaseReservationsForOrder(ctx, order.ID, repo, publisher, logger); err != nil {
		logger.Error("Failed to release reservations for order", err, map[string]interface{}{
			"order_id": order.ID,
		})
		return
	}
	
	logger.Info("Inventory released for cancelled order", map[string]interface{}{
		"order_id": order.ID,
	})
}

func releaseReservationsForOrder(ctx context.Context, orderID string, repo repositories.InventoryRepository, publisher messaging.Publisher, logger *observability.Logger) error {
	// Get all reservations for this order
	reservations, err := repo.GetReservationsByOrderID(ctx, orderID)
	if err != nil {
		return err
	}
	
	// No reservations found
	if len(reservations) == 0 {
		return nil
	}
	
	// Release each reservation
	for _, reservation := range reservations {
		// Get product
		product, err := repo.GetProductByID(ctx, reservation.ProductID)
		if err != nil {
			logger.Error("Failed to get product", err, map[string]interface{}{
				"product_id": reservation.ProductID,
				"order_id":   orderID,
			})
			continue
		}
		
		// Release stock
		if err := product.ReleaseStock(reservation.Quantity); err != nil {
			logger.Error("Failed to release stock", err, map[string]interface{}{
				"product_id": reservation.ProductID,
				"order_id":   orderID,
			})
			continue
		}
		
		// Update product stock
		if err := repo.UpdateProduct(ctx, product); err != nil {
			logger.Error("Failed to update product stock", err, map[string]interface{}{
				"product_id": reservation.ProductID,
				"order_id":   orderID,
			})
			continue
		}
		
		// Update reservation status
		if err := repo.UpdateReservationStatus(ctx, reservation.ID, "released"); err != nil {
			logger.Error("Failed to update reservation status", err, map[string]interface{}{
				"reservation_id": reservation.ID,
				"order_id":       orderID,
			})
			continue
		}
	}
	
	// Publish inventory released event
	err = publisher.Publish(messaging.InventoryReleased, map[string]interface{}{
		"order_id": orderID,
	})
	if err != nil {
		logger.Error("Failed to publish inventory released event", err, map[string]interface{}{
			"order_id": orderID,
		})
		// Continue even if event publishing fails
	}
	
	return nil
}
