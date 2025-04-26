package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/edgardnogueira/go-patterns/project_structures/microservices/pkg/observability"
)

const serviceName = "api-gateway"

var (
	port              = flag.Int("port", 8080, "Port for the API Gateway")
	orderServiceURL   = flag.String("order-service", "http://order-service:8081", "URL of the Order Service")
	inventoryServiceURL = flag.String("inventory-service", "http://inventory-service:8082", "URL of the Inventory Service")
)

func main() {
	flag.Parse()

	// Override from environment if provided
	if envPort := os.Getenv("PORT"); envPort != "" {
		fmt.Sscanf(envPort, "%d", port)
	}
	if envURL := os.Getenv("ORDER_SERVICE_URL"); envURL != "" {
		*orderServiceURL = envURL
	}
	if envURL := os.Getenv("INVENTORY_SERVICE_URL"); envURL != "" {
		*inventoryServiceURL = envURL
	}

	// Initialize logger
	logger := observability.NewLogger(serviceName)
	logger.Info("Starting API Gateway", map[string]interface{}{
		"port":                *port,
		"order_service_url":   *orderServiceURL,
		"inventory_service_url": *inventoryServiceURL,
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

	// Create router
	router := setupRouter(logger, *orderServiceURL, *inventoryServiceURL)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("API Gateway listening", map[string]interface{}{
			"port": *port,
		})
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Gracefully shutdown the server
	logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", err)
	}

	logger.Info("Server exited properly")
}
