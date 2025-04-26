package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	appservice "github.com/edgardnogueira/go-patterns/project_structures/layered/internal/application/service"
	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/data/repository"
	domainservice "github.com/edgardnogueira/go-patterns/project_structures/layered/internal/domain/service"
	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/infrastructure/config"
	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/infrastructure/logger"
	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/infrastructure/providers/database"
	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/infrastructure/providers/external"
	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/presentation/worker/handlers"
)

func main() {
	// Initialize configuration
	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	appLogger, err := logger.NewLogger(cfg.Logging)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Initialize database
	dbProvider, err := database.NewSqliteProvider(cfg.Database.FilePath)
	if err != nil {
		appLogger.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbProvider.Close()

	// Initialize database schema
	if err := dbProvider.Initialize(); err != nil {
		appLogger.Fatalf("Failed to initialize database schema: %v", err)
	}

	// Get database connection
	db := dbProvider.GetDB()

	// Initialize repositories
	postRepo := repository.NewSqlitePostRepository(db)
	authorRepo := repository.NewSqliteAuthorRepository(db)
	notificationRepo := repository.NewSqliteNotificationRepository(db)

	// Initialize domain services
	postService := domainservice.NewPostService(postRepo)
	authorService := domainservice.NewAuthorService(authorRepo)
	notificationService := domainservice.NewNotificationService(notificationRepo)

	// Initialize application services
	postAppService := appservice.NewPostAppService(postService, authorService)
	authorAppService := appservice.NewAuthorAppService(authorService, postService)
	notificationAppService := appservice.NewNotificationAppService(notificationService, authorService, postService)

	// Initialize email provider
	emailConfig := external.EmailConfig{
		Host:     "smtp.example.com",  // In a real app, this would come from config
		Port:     587,
		Username: "user@example.com",
		Password: "password",
		From:     "no-reply@example.com",
	}
	emailProvider := external.NewEmailProvider(emailConfig, appLogger)

	// Initialize notification worker
	notificationWorker := handlers.NewNotificationWorker(
		notificationAppService,
		emailProvider,
		appLogger,
		&cfg.Worker,
	)

	// Start worker
	notificationWorker.Start()
	appLogger.Info("Worker started successfully")

	// Create example task (for demonstration)
	createExampleTask(notificationWorker, appLogger)

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down worker...")
	notificationWorker.Stop()
	appLogger.Info("Worker stopped successfully")
}

// createExampleTask creates an example notification task (for demonstration)
func createExampleTask(worker *handlers.NotificationWorker, log *logger.Logger) {
	// Email notification task
	emailTask := handlers.NotificationTask{
		Type: handlers.SendEmailNotification,
		Data: map[string]interface{}{
			"to":      "user@example.com",
			"subject": "Welcome to Our Blog",
			"body":    "Thank you for signing up!",
		},
		CreatedAt:  handlers.Now(),
		RetryCount: 0,
	}

	// Enqueue task
	if err := worker.EnqueueTask(emailTask); err != nil {
		log.WithField("error", err.Error()).Error("Failed to enqueue example task")
		return
	}

	log.Info("Example task enqueued successfully")
}

// Now returns the current time
func Now() handlers.NotificationTask {
	return handlers.NotificationTask{}
}
