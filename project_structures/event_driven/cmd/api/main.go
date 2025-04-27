package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/edgardnogueira/go-patterns/project_structures/event_driven/internal/commands/handler"
	"github.com/edgardnogueira/go-patterns/project_structures/event_driven/internal/infrastructure/config"
	"github.com/edgardnogueira/go-patterns/project_structures/event_driven/internal/infrastructure/eventstore/postgres"
	"github.com/edgardnogueira/go-patterns/project_structures/event_driven/internal/infrastructure/messaging/rabbitmq"
	"github.com/edgardnogueira/go-patterns/project_structures/event_driven/pkg/cqrs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	log.Println("Starting API Server...")

	// Load configuration
	cfg, err := config.LoadConfig("config/api.yaml")
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

	// Initialize command bus
	commandBus := cqrs.NewCommandBus()

	// Register command handlers
	productCommandHandler := handler.NewProductCommandHandler(eventStore, eventPublisher)
	commandBus.Register("ReduceProductStock", productCommandHandler.HandleReduceProductStock)
	commandBus.Register("AddProductStock", productCommandHandler.HandleAddProductStock)
	commandBus.Register("CreateProduct", productCommandHandler.HandleCreateProduct)

	// Initialize HTTP router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/products", func(r chi.Router) {
			r.Post("/", productCommandHandler.HTTPHandleCreateProduct)
			r.Post("/{id}/reduce-stock", productCommandHandler.HTTPHandleReduceStock)
			r.Post("/{id}/add-stock", productCommandHandler.HTTPHandleAddStock)
		})
	})

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}

	// Handle graceful shutdown
	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("Shutting down server...")

		// Create a deadline for server shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Error during server shutdown: %v\n", err)
		}

		close(stopped)
	}()

	// Start the server
	log.Printf("Server is running on port %d", cfg.Server.Port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server error: %v", err)
	}

	<-stopped
	log.Println("Server stopped gracefully")
}
