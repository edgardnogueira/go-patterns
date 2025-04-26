package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/edgardnogueira/go-patterns/project_structures/clean_architecture/internal/delivery/http"
	"github.com/edgardnogueira/go-patterns/project_structures/clean_architecture/internal/domain/repositories"
	"github.com/edgardnogueira/go-patterns/project_structures/clean_architecture/internal/infrastructure/db/memory"
	"github.com/edgardnogueira/go-patterns/project_structures/clean_architecture/internal/usecase"
)

func main() {
	// Initialize repositories (using in-memory implementation for demo)
	var userRepo repositories.UserRepository = memory.NewUserMemoryRepository()
	var taskRepo repositories.TaskRepository = memory.NewTaskMemoryRepository()

	// Initialize services
	userService := usecase.NewUserService(userRepo, taskRepo)
	taskService := usecase.NewTaskService(taskRepo)

	// Initialize HTTP server
	server := http.NewServer()

	// Initialize and register handlers
	userHandler := http.NewUserHandler(userService)
	userHandler.RegisterRoutes(server.Router())

	taskHandler := http.NewTaskHandler(taskService)
	taskHandler.RegisterRoutes(server.Router())

	// Create channel for shutdown signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		addr := ":8080"
		log.Printf("Starting API server on %s", addr)
		if err := server.Start(addr); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-stop
	log.Println("Shutting down server...")

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")
}
