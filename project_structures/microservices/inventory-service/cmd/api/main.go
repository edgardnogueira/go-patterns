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

	"github.com/edgardnogueira/go-patterns/project_structures/microservices/inventory-service/internal/handlers"
	"github.com/edgardnogueira/go-patterns/project_structures/microservices/inventory-service/internal/repositories"
	"github.com/edgardnogueira/go-patterns/project_structures/microservices/pkg/messaging"
	"github.com/edgardnogueira/go-patterns/project_structures/microservices/pkg/observability"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const serviceName = "inventory-service"

var (
	port            = flag.Int("port", 8082, "Port for the Inventory Service API")
	postgresConnStr = flag.String("postgres-dsn", "postgres://postgres:postgres@postgres:5432/inventory?sslmode=disable", "PostgreSQL connection string")
	natsURL         = flag.String("nats-url", "nats://nats:4222", "NATS server URL")
)

func main() {
	flag.Parse()

	// Override from environment if provided
	if envPort := os.Getenv("PORT"); envPort != "" {
		fmt.Sscanf(envPort, "%d", port)
	}
	if envDSN := os.Getenv("POSTGRES_DSN"); envDSN != "" {
		*postgresConnStr = envDSN
	}
	if envNATS := os.Getenv("NATS_URL"); envNATS != "" {
		*natsURL = envNATS
	}

	// Initialize logger
	logger := observability.NewLogger(serviceName)
	logger.Info("Starting Inventory Service API", map[string]interface{}{
		"port": *port,
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
	repo, err := repositories.NewPostgresInventoryRepository(*postgresConnStr)
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
	router := setupRouter(logger, repo, publisher)

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
		logger.Info("Inventory Service API listening", map[string]interface{}{
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

func setupRouter(logger *observability.Logger, repo repositories.InventoryRepository, publisher messaging.Publisher) http.Handler {
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
		w.Write([]byte(`{"status":"ok","service":"inventory-service"}`))
	})

	// API Routes
	r.Route("/api/inventory", func(r chi.Router) {
		inventoryHandler := handlers.NewInventoryHandler(repo, publisher, logger)
		r.Get("/", inventoryHandler.ListProducts)
		r.Post("/", inventoryHandler.CreateProduct)
		r.Get("/{id}", inventoryHandler.GetProduct)
		r.Put("/{id}", inventoryHandler.UpdateProduct)
		r.Delete("/{id}", inventoryHandler.DeleteProduct)
		r.Post("/reserve", inventoryHandler.ReserveStock)
	})

	// Metrics
	r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, fmt.Sprintf("http://localhost:%d/metrics", *port+100), http.StatusTemporaryRedirect)
	})

	return r
}
