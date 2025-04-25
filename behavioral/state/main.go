package main

import (
	"fmt"
	"github.com/edgardnogueira/go-patterns/behavioral/state"
	"time"
)

// Simple demonstration of the State Pattern
func main() {
	fmt.Println("State Pattern - Package Delivery System")
	fmt.Println("======================================")
	fmt.Println()
	fmt.Println("This example demonstrates the State pattern using a package delivery system")
	fmt.Println("where a package transitions through different states (ordered, processing,")
	fmt.Println("shipped, delivered, returned, canceled) that change its behavior.")
	fmt.Println()

	// Create a new package
	pkg := state.NewPackage("PKG123", "Electronics")
	state.InitializePackage(pkg)
	
	// Add a logging handler
	pkg.AddTransitionHandler(func(e state.Event) {
		fmt.Printf("Event: %s -> %s at %s\n", 
			e.From, e.To, e.Timestamp.Format(time.RFC3339))
	})

	// Display package info
	fmt.Printf("Package: %s - %s\n", pkg.ID, pkg.Description)
	fmt.Printf("Current state: %s\n", pkg.GetState())
	fmt.Printf("Allowed transitions: %v\n\n", pkg.GetAllowedTransitions())

	// Demonstrate the state transitions
	fmt.Println("Processing the package...")
	err := pkg.HandleProcess()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Current state: %s\n", pkg.GetState())
		fmt.Printf("Allowed transitions: %v\n\n", pkg.GetAllowedTransitions())
	}

	fmt.Println("Shipping the package...")
	err = pkg.HandleShip()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Current state: %s\n", pkg.GetState())
		fmt.Printf("Allowed transitions: %v\n\n", pkg.GetAllowedTransitions())
	}

	fmt.Println("Delivering the package...")
	err = pkg.HandleDeliver()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Current state: %s\n", pkg.GetState())
		fmt.Printf("Allowed transitions: %v\n\n", pkg.GetAllowedTransitions())
	}

	// Demonstrate invalid operations
	fmt.Println("Attempting to ship an already delivered package...")
	err = pkg.HandleShip()
	fmt.Printf("Result: %v\n\n", err)

	// Create another package to demonstrate a different flow
	fmt.Println("Creating a new package to demonstrate cancellation...")
	pkg2 := state.NewPackage("PKG456", "Books")
	state.InitializePackage(pkg2)
	
	fmt.Printf("Package: %s - %s\n", pkg2.ID, pkg2.Description)
	fmt.Printf("Current state: %s\n", pkg2.GetState())
	
	fmt.Println("Canceling the order...")
	err = pkg2.HandleCancel()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Current state: %s\n", pkg2.GetState())
		fmt.Printf("Allowed transitions: %v\n\n", pkg2.GetAllowedTransitions())
	}
	
	fmt.Println("Attempting to process a canceled order...")
	err = pkg2.HandleProcess()
	fmt.Printf("Result: %v\n\n", err)

	// Explain the pattern
	fmt.Println("State Pattern Key Points:")
	fmt.Println("------------------------")
	fmt.Println("1. The State pattern allows an object to alter its behavior when its internal state changes.")
	fmt.Println("2. Each state is represented by a separate class implementing a common interface.")
	fmt.Println("3. The context delegates state-specific behavior to the current state object.")
	fmt.Println("4. State transitions can be controlled by the states themselves or by the context.")
	fmt.Println("5. The pattern eliminates the need for large conditional statements based on the object's state.")
	fmt.Println()
	
	fmt.Println("For more details and a more comprehensive example, see the example/ directory.")
}
