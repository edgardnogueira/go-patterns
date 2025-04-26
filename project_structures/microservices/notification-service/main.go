package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/edgardnogueira/go-patterns/project_structures/microservices/pkg/messaging"
	"github.com/edgardnogueira/go-patterns/project_structures/microservices/pkg/observability"
)

const serviceName = "notification-service"

var (
	port    = flag.Int("port", 8083, "Port for the Notification Service")
	natsURL = flag.String("nats-url", "nats://nats:4222", "NATS server URL")
	smtpHost = flag.String("smtp-host", "mailhog", "SMTP server host")
	smtpPort = flag.Int("smtp-port", 1025, "SMTP server port")
)

func main() {
	flag.Parse()

	// Override from environment if provided
	if envPort := os.Getenv("PORT"); envPort != "" {
		fmt.Sscanf(envPort, "%d", port)
	}
	if envNATS := os.Getenv("NATS_URL"); envNATS != "" {
		*natsURL = envNATS
	}
	if envHost := os.Getenv("SMTP_HOST"); envHost != "" {
		*smtpHost = envHost
	}
	if envPort := os.Getenv("SMTP_PORT"); envPort != "" {
		fmt.Sscanf(envPort, "%d", smtpPort)
	}

	// Initialize logger
	logger := observability.NewLogger(serviceName)
	logger.Info("Starting Notification Service", map[string]interface{}{
		"port":      *port,
		"smtp_host": *smtpHost,
		"smtp_port": *smtpPort,
	})

	// Initialize metrics
	metricsServer := observability.NewMetricsServer(serviceName, *port+100)
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

	// Initialize message consumer
	consumer, err := messaging.NewNatsClient(*natsURL)
	if err != nil {
		logger.Fatal("Failed to connect to NATS", err)
	}
	defer consumer.Close()

	// Create email service
	emailService := NewEmailService(*smtpHost, *smtpPort, logger)

	// Subscribe to events
	setupSubscriptions(consumer, emailService, logger)

	// Create and start HTTP server for health checks
	server := setupHTTPServer(*port)
	go func() {
		logger.Info("HTTP server listening", map[string]interface{}{
			"port": *port,
		})
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start HTTP server", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Gracefully shutdown the server
	logger.Info("Shutting down notification service...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", err)
	}

	logger.Info("Service exited properly")
}

// EmailService handles sending emails to users
type EmailService struct {
	smtpHost string
	smtpPort int
	logger   *observability.Logger
}

// NewEmailService creates a new email service
func NewEmailService(smtpHost string, smtpPort int, logger *observability.Logger) *EmailService {
	return &EmailService{
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		logger:   logger,
	}
}

// SendOrderConfirmation sends an order confirmation email
func (s *EmailService) SendOrderConfirmation(orderID, customerID, email string) error {
	// In a real implementation, this would connect to an SMTP server
	// and send an actual email. For demo purposes, we'll just log it.
	s.logger.Info("Sending order confirmation email", map[string]interface{}{
		"order_id":    orderID,
		"customer_id": customerID,
		"email":       email,
	})
	return nil
}

// SendOrderShippedNotification sends an order shipped notification
func (s *EmailService) SendOrderShippedNotification(orderID, customerID, email string) error {
	s.logger.Info("Sending order shipped notification", map[string]interface{}{
		"order_id":    orderID,
		"customer_id": customerID,
		"email":       email,
	})
	return nil
}

// SendOrderCancelledNotification sends an order cancellation notification
func (s *EmailService) SendOrderCancelledNotification(orderID, customerID, email string) error {
	s.logger.Info("Sending order cancellation notification", map[string]interface{}{
		"order_id":    orderID,
		"customer_id": customerID,
		"email":       email,
	})
	return nil
}

func setupSubscriptions(consumer messaging.Consumer, emailService *EmailService, logger *observability.Logger) {
	// Subscribe to OrderCreated events
	err := consumer.Subscribe(messaging.OrderCreated, func(event messaging.Event) {
		processOrderCreated(event, emailService, logger)
	})
	if err != nil {
		logger.Fatal("Failed to subscribe to OrderCreated events", err)
	}

	// Subscribe to OrderShipped events
	err = consumer.Subscribe(messaging.OrderShipped, func(event messaging.Event) {
		processOrderShipped(event, emailService, logger)
	})
	if err != nil {
		logger.Fatal("Failed to subscribe to OrderShipped events", err)
	}

	// Subscribe to OrderCancelled events
	err = consumer.Subscribe(messaging.OrderCancelled, func(event messaging.Event) {
		processOrderCancelled(event, emailService, logger)
	})
	if err != nil {
		logger.Fatal("Failed to subscribe to OrderCancelled events", err)
	}

	logger.Info("Subscribed to events")
}

func processOrderCreated(event messaging.Event, emailService *EmailService, logger *observability.Logger) {
	// Extract order data from event
	var order struct {
		ID         string `json:"id"`
		CustomerID string `json:"customer_id"`
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

	// In a real application, we would look up the customer's email
	// For demo purposes, we'll use a fake email
	customerEmail := fmt.Sprintf("%s@example.com", order.CustomerID)

	// Send confirmation email
	if err := emailService.SendOrderConfirmation(order.ID, order.CustomerID, customerEmail); err != nil {
		logger.Error("Failed to send order confirmation", err, map[string]interface{}{
			"order_id":    order.ID,
			"customer_id": order.CustomerID,
		})
	}
}

func processOrderShipped(event messaging.Event, emailService *EmailService, logger *observability.Logger) {
	// Extract order data from event
	var order struct {
		ID         string `json:"id"`
		CustomerID string `json:"customer_id"`
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

	// In a real application, we would look up the customer's email
	customerEmail := fmt.Sprintf("%s@example.com", order.CustomerID)

	// Send shipped notification
	if err := emailService.SendOrderShippedNotification(order.ID, order.CustomerID, customerEmail); err != nil {
		logger.Error("Failed to send order shipped notification", err, map[string]interface{}{
			"order_id":    order.ID,
			"customer_id": order.CustomerID,
		})
	}
}

func processOrderCancelled(event messaging.Event, emailService *EmailService, logger *observability.Logger) {
	// Extract order data from event
	var order struct {
		ID         string `json:"id"`
		CustomerID string `json:"customer_id"`
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

	// In a real application, we would look up the customer's email
	customerEmail := fmt.Sprintf("%s@example.com", order.CustomerID)

	// Send cancellation notification
	if err := emailService.SendOrderCancelledNotification(order.ID, order.CustomerID, customerEmail); err != nil {
		logger.Error("Failed to send order cancellation notification", err, map[string]interface{}{
			"order_id":    order.ID,
			"customer_id": order.CustomerID,
		})
	}
}

func setupHTTPServer(port int) *http.Server {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","service":"notification-service"}`))
	})

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return server
}
