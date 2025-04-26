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
	"github.com/go-chi/cors"

	appservice "github.com/edgardnogueira/go-patterns/project_structures/layered/internal/application/service"
	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/data/repository"
	domainservice "github.com/edgardnogueira/go-patterns/project_structures/layered/internal/domain/service"
	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/infrastructure/config"
	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/infrastructure/logger"
	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/infrastructure/providers/database"
	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/presentation/api/handlers"
	appmiddleware "github.com/edgardnogueira/go-patterns/project_structures/layered/internal/presentation/api/middleware"
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

	// Initialize HTTP handlers
	postHandler := handlers.NewPostHandler(postAppService, authorAppService, appLogger)
	authorHandler := handlers.NewAuthorHandler(authorAppService, appLogger)

	// Initialize router
	r := chi.NewRouter()

	// Add middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(appmiddleware.RequestLoggerMiddleware(appLogger))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Register routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Register handlers
		postHandler.RegisterRoutes(r)
		authorHandler.RegisterRoutes(r)
	})

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:           addr,
		Handler:        r,
		ReadTimeout:    time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(cfg.Server.WriteTimeout) * time.Second,
		MaxHeaderBytes: cfg.Server.MaxHeaderBytes,
	}

	// Start server in a goroutine
	go func() {
		appLogger.Infof("Starting server on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Fatalf("Server forced to shutdown: %v", err)
	}

	appLogger.Info("Server exited gracefully")
}
