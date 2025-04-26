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
	
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	
	"github.com/edgardnogueira/go-patterns/project_structures/hexagonal/internal/adapters/driven/database"
	"github.com/edgardnogueira/go-patterns/project_structures/hexagonal/internal/adapters/driven/external"
	httpAdapter "github.com/edgardnogueira/go-patterns/project_structures/hexagonal/internal/adapters/driving/http"
	"github.com/edgardnogueira/go-patterns/project_structures/hexagonal/internal/domain/service"
)

func main() {
	// Setup logger
	logger := log.New(os.Stdout, "[API] ", log.LstdFlags)
	logger.Println("Starting API server...")
	
	// Setup ports & adapters
	orderRepo := database.NewMemoryOrderRepository()
	notificationService := external.NewLogNotificationService(logger)
	orderService := service.NewOrderService(orderRepo, notificationService)
	orderHandler := httpAdapter.NewOrderHandler(orderService)
	
	// Setup router with middleware
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	
	// Register API routes
	r.Route("/api", func(r chi.Router) {
		// Register order routes
		orderHandler.RegisterRoutes(r)
	})
	
	// Add health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	
	// Create server
	addr := ":8080"
	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	// Start server in a goroutine
	go func() {
		logger.Printf("Server listening on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Error starting server: %v", err)
		}
	}()
	
	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	
	logger.Println("Shutting down server...")
	
	// Create a deadline to wait for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}
	
	logger.Println("Server gracefully stopped")
}
