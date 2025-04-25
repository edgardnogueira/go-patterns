package main

import (
	"fmt"
	"github.com/edgardnogueira/go-patterns/behavioral/state"
	"os"
	"time"
)

func main() {
	fmt.Println("State Pattern Example - Package Delivery System")
	fmt.Println("===============================================")
	fmt.Println()

	// Create a new package
	fmt.Println("Creating a new package...")
	pkg := state.NewPackage("PKG12345", "Smartphone")
	state.InitializePackage(pkg)

	// Add a logging handler to see state transitions
	pkg.AddTransitionHandler(state.LoggingHandler(func(message string) {
		fmt.Println(message)
	}))

	// Add a notification handler
	pkg.AddTransitionHandler(state.NotificationHandler(func(message, details string) {
		fmt.Printf("NOTIFICATION: %s\n  Details: %s\n", message, details)
	}))

	// Display initial state
	fmt.Printf("\nPackage ID: %s\n", pkg.ID)
	fmt.Printf("Description: %s\n", pkg.Description)
	fmt.Printf("Current State: %s\n", pkg.GetState())
	fmt.Printf("Created At: %s\n", pkg.CreatedAt.Format(time.RFC3339))
	fmt.Printf("Allowed Transitions: %v\n\n", pkg.GetAllowedTransitions())

	// Process the package
	fmt.Println("1. Processing the package...")
	if err := pkg.HandleProcess(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Current State: %s\n", pkg.GetState())
	fmt.Printf("Allowed Transitions: %v\n\n", pkg.GetAllowedTransitions())

	// Ship the package
	fmt.Println("2. Shipping the package...")
	if err := pkg.HandleShip(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Current State: %s\n", pkg.GetState())
	fmt.Printf("Allowed Transitions: %v\n\n", pkg.GetAllowedTransitions())

	// Try an invalid transition
	fmt.Println("3. Trying to process the package again (invalid transition)...")
	err := pkg.HandleProcess()
	fmt.Printf("Result: %v\n\n", err)

	// Deliver the package
	fmt.Println("4. Delivering the package...")
	if err := pkg.HandleDeliver(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Current State: %s\n", pkg.GetState())
	fmt.Printf("Allowed Transitions: %v\n\n", pkg.GetAllowedTransitions())

	// Create a second package to demonstrate different flows
	fmt.Println("Creating a second package to demonstrate cancellation...")
	pkg2 := state.NewPackage("PKG67890", "Headphones")
	state.InitializePackage(pkg2)

	// Add a logging handler
	pkg2.AddTransitionHandler(state.LoggingHandler(func(message string) {
		fmt.Println(message)
	}))

	fmt.Printf("\nPackage ID: %s\n", pkg2.ID)
	fmt.Printf("Current State: %s\n", pkg2.GetState())
	
	// Process the package
	fmt.Println("1. Processing the package...")
	if err := pkg2.HandleProcess(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Current State: %s\n", pkg2.GetState())

	// Cancel the order
	fmt.Println("2. Canceling the order...")
	if err := pkg2.HandleCancel(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Current State: %s\n", pkg2.GetState())

	// Try to ship the canceled package
	fmt.Println("3. Trying to ship the canceled package...")
	err = pkg2.HandleShip()
	fmt.Printf("Result: %v\n\n", err)

	// Create a third package to demonstrate automatic transitions
	fmt.Println("Creating a third package to demonstrate automatic transitions...")
	pkg3 := state.NewPackage("PKG24680", "Book")
	state.InitializePackage(pkg3)

	// Add a logging handler
	pkg3.AddTransitionHandler(state.LoggingHandler(func(message string) {
		fmt.Println(message)
	}))

	fmt.Printf("\nPackage ID: %s\n", pkg3.ID)
	fmt.Printf("Current State: %s\n", pkg3.GetState())
	
	// Process the package
	fmt.Println("1. Processing the package...")
	if err := pkg3.HandleProcess(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Current State: %s\n", pkg3.GetState())

	// Set up automatic transition after 2 seconds
	fmt.Println("\n2. Setting up automatic transition to Shipped state in 2 seconds...")
	timer := state.TimeoutTransition(pkg3, "Shipped", 2*time.Second)
	defer timer.Stop()

	fmt.Println("Waiting for automatic transition...")
	time.Sleep(3 * time.Second)

	fmt.Printf("Current State: %s\n\n", pkg3.GetState())

	// Display the state history for the first package
	fmt.Println("State History for package", pkg.ID)
	fmt.Println("----------------------------")
	for i, event := range pkg.GetStateHistory() {
		fmt.Printf("%d. [%s] %s -> %s: %s\n",
			i+1,
			event.Timestamp.Format(time.RFC3339),
			event.From,
			event.To,
			event.Details)
	}

	fmt.Println("\nState Pattern Benefits:")
	fmt.Println("1. Each state encapsulates its own behavior")
	fmt.Println("2. State transitions are explicitly defined and validated")
	fmt.Println("3. Adding new states doesn't affect existing code")
	fmt.Println("4. The context (Package) delegates behavior to the current state")
	fmt.Println("5. The pattern makes state transitions explicit in the code")
}
