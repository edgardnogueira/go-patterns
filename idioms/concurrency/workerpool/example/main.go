package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/edgardnogueira/go-patterns/idioms/concurrency/workerpool"
)

// ImageProcessingTask represents a task that simulates image processing
type ImageProcessingTask struct {
	ID           int
	ImageSize    int // Simulated image size in MB
	ProcessingOp string
}

// Execute implements the Task interface and processes the "image"
func (t *ImageProcessingTask) Execute() interface{} {
	fmt.Printf("Worker processing image #%d (%d MB) - %s...\n", 
		t.ID, t.ImageSize, t.ProcessingOp)
	
	// Simulate the actual work based on image size
	processingTime := time.Duration(t.ImageSize * 100) * time.Millisecond
	time.Sleep(processingTime)
	
	// Simulate occasional errors
	if rand.Intn(10) == 0 {
		return fmt.Errorf("failed to process image #%d: corrupted data", t.ID)
	}
	
	return fmt.Sprintf("Image #%d successfully processed (%s)", t.ID, t.ProcessingOp)
}

// Define a set of image processing operations
var processingOperations = []string{
	"Resize",
	"Crop",
	"Filter",
	"Rotate",
	"Compress",
	"Convert",
	"Color adjustment",
	"Noise reduction",
}

func main() {
	fmt.Println("=== Worker Pool Pattern - Image Processing Service ===")
	
	// Parse worker count from args or use default
	workerCount := 4
	if len(os.Args) > 1 {
		if count, err := strconv.Atoi(os.Args[1]); err == nil && count > 0 {
			workerCount = count
		}
	}
	
	fmt.Printf("Starting image processing service with %d workers\n", workerCount)
	
	// Create and start worker pool
	pool := workerpool.NewPool(workerCount)
	pool.Start()
	defer pool.Stop()
	
	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	// Start a goroutine to collect and process results
	go processResults(pool.Results)
	
	// Create a ticker for generating new tasks
	taskTicker := time.NewTicker(500 * time.Millisecond)
	defer taskTicker.Stop()
	
	// Keep track of submitted tasks
	taskCounter := 0
	
	fmt.Println("\nSubmitting random image processing tasks...")
	fmt.Println("Press Ctrl+C to stop the service")
	
	// Main loop for submitting tasks
	for {
		select {
		case <-taskTicker.C:
			taskCounter++
			
			// Create a new image processing task with random parameters
			task := &ImageProcessingTask{
				ID:           taskCounter,
				ImageSize:    rand.Intn(10) + 1, // 1-10 MB
				ProcessingOp: processingOperations[rand.Intn(len(processingOperations))],
			}
			
			fmt.Printf("Submitting task #%d: Process %d MB image with operation '%s'\n", 
				task.ID, task.ImageSize, task.ProcessingOp)
			
			pool.Submit(task)
			
		case <-sigChan:
			fmt.Println("\nReceived shutdown signal. Stopping service...")
			return
		}
	}
}

// processResults handles the results from the worker pool
func processResults(results <-chan workerpool.Result) {
	var (
		successCount int
		errorCount   int
	)
	
	for result := range results {
		switch res := result.Value.(type) {
		case string:
			// Success case
			fmt.Printf("✅ Worker %d: %s\n", result.WorkerID, res)
			successCount++
			
		case error:
			// Error case
			fmt.Printf("❌ Worker %d: ERROR - %s\n", result.WorkerID, res.Error())
			errorCount++
			
		default:
			fmt.Printf("⚠️ Worker %d: Unknown result type: %v\n", result.WorkerID, result.Value)
		}
		
		// Print stats periodically
		if (successCount+errorCount)%10 == 0 {
			fmt.Printf("\n--- Statistics ---\n")
			fmt.Printf("Tasks processed: %d (Success: %d, Errors: %d)\n", 
				successCount+errorCount, successCount, errorCount)
			fmt.Printf("-----------------\n\n")
		}
	}
}
