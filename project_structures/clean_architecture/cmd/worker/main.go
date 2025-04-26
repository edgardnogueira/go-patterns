package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/edgardnogueira/go-patterns/project_structures/clean_architecture/internal/delivery/worker"
	"github.com/edgardnogueira/go-patterns/project_structures/clean_architecture/internal/domain/repositories"
	"github.com/edgardnogueira/go-patterns/project_structures/clean_architecture/internal/infrastructure/db/memory"
	"github.com/edgardnogueira/go-patterns/project_structures/clean_architecture/internal/infrastructure/providers"
	"github.com/edgardnogueira/go-patterns/project_structures/clean_architecture/internal/usecase"
)

func main() {
	// Initialize repositories (using in-memory implementation for demo)
	var taskRepo repositories.TaskRepository = memory.NewTaskMemoryRepository()

	// Initialize providers
	emailProvider := providers.NewMockEmailProvider()

	// Initialize services
	taskService := usecase.NewTaskService(taskRepo)

	// Initialize task processor
	// Poll interval of 5 seconds
	processor := worker.NewTaskProcessor(taskService, emailProvider, 5*time.Second)

	// Create context that can be canceled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the processor
	processor.Start(ctx)
	log.Println("Worker started and processing tasks")

	// Create channel for shutdown signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Wait for shutdown signal
	<-stop
	log.Println("Shutting down worker...")

	// Stop the processor
	processor.Stop()

	log.Println("Worker gracefully stopped")
}
