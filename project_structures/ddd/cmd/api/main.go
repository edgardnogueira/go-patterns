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
)

func main() {
	log.Println("Starting API server...")

	// Setup server
	server := setupServer()

	// Start server
	go func() {
		log.Printf("Server listening on %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %v\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gracefully shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

	log.Println("Server exited properly")
}

func setupServer() *http.Server {
	// Create a new router
	mux := http.NewServeMux()

	// Register API routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok","message":"DDD E-commerce API"}`)
	})

	mux.HandleFunc("/orders", handleOrders)
	mux.HandleFunc("/inventory", handleInventory)
	mux.HandleFunc("/users", handleUsers)

	// Create and configure the server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return server
}

func handleOrders(w http.ResponseWriter, r *http.Request) {
	// This is just a placeholder. In a real implementation, we would:
	// 1. Parse the request (e.g., JSON)
	// 2. Call the appropriate application service
	// 3. Return the response
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"ok","message":"Order endpoints"}`)
}

func handleInventory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"ok","message":"Inventory endpoints"}`)
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"ok","message":"User endpoints"}`)
}
