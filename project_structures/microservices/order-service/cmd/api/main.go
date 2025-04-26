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

	"github.com/edgardnogueira/go-patterns/project_structures/microservices/order-service/internal/handlers"
	"github.com/edgardnogueira/go-patterns/project_structures/microservices/order-service/internal/repositories"
	"github.com/edgardnogueira/go-patterns/project_structures/microservices/pkg/messaging"
	"github.com/edgardnogueira/go-patterns/project_structures/microservices/pkg/observability"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const serviceName = "order-service"

var (
	port            = flag.Int("port", 8081, "Port for the Order Service API")
	inventoryURL    = flag.String("inventory-url", "http://inventory-service:8082", "URL for the Inventory Service")
	postgresConnStr = flag.String("postgres-dsn", "postgres://postgres:postgres@postgres:5432/orders?sslmode=disable", "PostgreSQL connection string")
	natsURL         = flag.String("nats-url", "nats://nats:4222", "NATS server URL")
)

func main() {
	flag.Parse()

	// Override from environment if provided
	if envPort := os.Getenv("PORT"); envPort != "" {
		fmt.Sscanf(envPort, "%d", port)
	}
	if envURL := os.Getenv("INVENTORY_SERVICE_URL"); envURL != "" {
		*inventoryURL = envURL
	}
	if envDSN := os.Getenv("POSTGRES_DSN"); envDSN != "" {
		*postgresConnStr = envDSN
	}
	if envNATS := os.Getenv("NATS_URL"); envNATS != "" {
		*natsURL = envNATS
	}

	// Initialize logger
	logger := observability.NewLogger(serviceName)
	logger.Info("Starting Order Service API", map[string]interface{}{
		"port":          *port,
		"inventory_url": *inventoryURL,
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

	// Initialize database repository
	repo, err := repositories.NewPostgresOrderRepository(*postgresConnStr)
	if err != nil {
		logger.Fatal("Failed to connect to database", err)
	}
	defer repo.Close()

	// Initialize message publisher
	publisher, err := messaging.NewNatsClient(*natsURL)
	if err != nil {
		logger.Fatal("Failed to connect to NATS", err)
	}
	defer publisher.Close()

	// Create router
	router := setupRouter(logger, repo, publisher, *inventoryURL)

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
		logger.Info("Order Service API listening", map[string]interface{}{
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

func setupRouter(logger *observability.Logger, repo repositories.OrderRepository, publisher messaging.Publisher, inventoryURL string) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(observability.TraceMiddleware(serviceName))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","service":"order-service"}`))
	})

	// API Routes
	r.Route("/api/orders", func(r chi.Router) {
		orderHandler := handlers.NewOrderHandler(repo, publisher, inventoryURL, logger)
		r.Get("/", orderHandler.ListOrders)
		r.Post("/", orderHandler.CreateOrder)
		r.Get("/{id}", orderHandler.GetOrder)
		r.Put("/{id}", orderHandler.UpdateOrder)
		r.Delete("/{id}", orderHandler.CancelOrder)
	})

	// Metrics
	r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, fmt.Sprintf("http://localhost:%d/metrics", *port+100), http.StatusTemporaryRedirect)
	})

	return r
}
