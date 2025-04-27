// Package concurrency demonstrates Go's idiomatic concurrency patterns
package concurrency

import (
	"fmt"
	"os"
	"strings"
)

// PrintPatterns displays all available concurrency patterns
func PrintPatterns() {
	fmt.Println("Go Concurrency Patterns")
	fmt.Println("=======================")
	fmt.Println()
	
	patterns := []struct {
		Name        string
		Description string
		Path        string
	}{
		{
			Name:        "Worker Pool",
			Description: "Pool of worker goroutines that process tasks from a shared queue",
			Path:        "workerpool",
		},
		{
			Name:        "Fan-Out/Fan-In",
			Description: "Distributing work across multiple goroutines and collecting results",
			Path:        "fanout_fanin",
		},
		{
			Name:        "Pipeline",
			Description: "Series of stages connected by channels, each stage processing data",
			Path:        "pipeline",
		},
		{
			Name:        "Context Cancellation",
			Description: "Using context to handle cancellation across goroutines",
			Path:        "context",
		},
		{
			Name:        "Rate Limiting",
			Description: "Controlling throughput of operations in concurrent systems",
			Path:        "rate_limiting",
		},
		{
			Name:        "Pub/Sub Pattern",
			Description: "Publishing events to multiple subscribers",
			Path:        "pubsub",
		},
		{
			Name:        "Semaphore Pattern",
			Description: "Limiting concurrent access to a resource",
			Path:        "semaphore",
		},
		{
			Name:        "Futures/Promises",
			Description: "Asynchronous result handling in Go style",
			Path:        "future",
		},
		{
			Name:        "Channel Generators",
			Description: "Functions that yield a stream of values via channels",
			Path:        "generator",
		},
		{
			Name:        "Heartbeat Pattern",
			Description: "Monitoring long-running processes and detecting failures",
			Path:        "heartbeat",
		},
	}
	
	for i, pattern := range patterns {
		fmt.Printf("%d. %s\n", i+1, pattern.Name)
		fmt.Printf("   Description: %s\n", pattern.Description)
		fmt.Printf("   Directory: ./idioms/concurrency/%s\n", pattern.Path)
		fmt.Println()
	}
}

func main() {
	// Get the pattern name from command line arguments
	if len(os.Args) > 1 {
		patternName := strings.ToLower(os.Args[1])
		
		switch patternName {
		case "workerpool", "worker", "pool":
			fmt.Println("Running Worker Pool Pattern Example")
			fmt.Println("See ./workerpool/example/main.go for implementation")
			// Here we could directly call the example if it was implemented
			
		case "fanout", "fanin", "fan-out", "fan-in":
			fmt.Println("Running Fan-Out/Fan-In Pattern Example")
			fmt.Println("See ./fanout_fanin/example/main.go for implementation")
			// Here we could directly call the example if it was implemented
			
		default:
			fmt.Printf("Pattern '%s' not recognized or not yet implemented\n", patternName)
		}
		
		return
	}
	
	// If no arguments provided, print all available patterns
	PrintPatterns()
}
