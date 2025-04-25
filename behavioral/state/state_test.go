package state

import (
	"strings"
	"testing"
	"time"
)

// TestPackageCreation tests the creation of a new package
func TestPackageCreation(t *testing.T) {
	p := NewPackage("PKG123", "Test Package")
	InitializePackage(p)

	// Check initial state
	if p.CurrentState == nil {
		t.Fatal("Package should have an initial state")
	}

	if p.CurrentState.Name() != "Ordered" {
		t.Errorf("Package should be in Ordered state, got %s", p.CurrentState.Name())
	}

	// Check package properties
	if p.ID != "PKG123" {
		t.Errorf("Package ID should be PKG123, got %s", p.ID)
	}

	if p.Description != "Test Package" {
		t.Errorf("Package description should be 'Test Package', got %s", p.Description)
	}

	// Check that history is initialized
	if len(p.History) != 1 {
		t.Errorf("Package should have 1 history entry, got %d", len(p.History))
	}

	// Check that the timestamps are set
	if p.CreatedAt.IsZero() {
		t.Error("Package creation time should be set")
	}

	if p.LastUpdatedAt.IsZero() {
		t.Error("Package last updated time should be set")
	}
}

// TestOrderLifecycle tests a typical package lifecycle
func TestOrderLifecycle(t *testing.T) {
	p := NewPackage("PKG456", "Lifecycle Test")
	InitializePackage(p)

	// Test initial state
	if p.GetState() != "Ordered" {
		t.Fatalf("Package should start in Ordered state, got %s", p.GetState())
	}

	// Test valid state transitions
	// Ordered -> Processing
	err := p.HandleProcess()
	if err != nil {
		t.Fatalf("Failed to process package: %v", err)
	}
	if p.GetState() != "Processing" {
		t.Fatalf("Package should be in Processing state, got %s", p.GetState())
	}

	// Processing -> Shipped
	err = p.HandleShip()
	if err != nil {
		t.Fatalf("Failed to ship package: %v", err)
	}
	if p.GetState() != "Shipped" {
		t.Fatalf("Package should be in Shipped state, got %s", p.GetState())
	}

	// Shipped -> Delivered
	err = p.HandleDeliver()
	if err != nil {
		t.Fatalf("Failed to deliver package: %v", err)
	}
	if p.GetState() != "Delivered" {
		t.Fatalf("Package should be in Delivered state, got %s", p.GetState())
	}

	// Check history entries
	if len(p.History) != 4 { // Initial + 3 transitions
		t.Errorf("Package should have 4 history entries, got %d", len(p.History))
	}

	// Verify the last transition
	lastTransition := p.History[len(p.History)-1]
	if lastTransition.From != "Shipped" || lastTransition.To != "Delivered" {
		t.Errorf("Last transition should be Shipped -> Delivered, got %s -> %s", 
			lastTransition.From, lastTransition.To)
	}
}

// TestInvalidTransitions tests that invalid state transitions are rejected
func TestInvalidTransitions(t *testing.T) {
	p := NewPackage("PKG789", "Invalid Transitions Test")
	InitializePackage(p)

	// Attempt invalid transitions
	
	// Cannot deliver a package that's just been ordered
	err := p.HandleDeliver()
	if err == nil {
		t.Error("Should not be able to deliver a package that's just been ordered")
	}

	// Cannot directly go from Ordered to Delivered
	err = p.TransitionTo("Delivered")
	if err == nil {
		t.Error("Should not be able to transition directly from Ordered to Delivered")
	}

	// Try a valid transition followed by an invalid one
	err = p.HandleProcess() // Ordered -> Processing
	if err != nil {
		t.Fatalf("Failed to process package: %v", err)
	}

	// Cannot deliver a package that's being processed
	err = p.HandleDeliver()
	if err == nil {
		t.Error("Should not be able to deliver a package that's being processed")
	}
}

// TestCancellation tests the cancellation workflow
func TestCancellation(t *testing.T) {
	p := NewPackage("PKG101", "Cancellation Test")
	InitializePackage(p)

	// Cancel from Ordered state
	err := p.HandleCancel()
	if err != nil {
		t.Fatalf("Failed to cancel package: %v", err)
	}
	if p.GetState() != "Canceled" {
		t.Fatalf("Package should be in Canceled state, got %s", p.GetState())
	}

	// Attempt operations on canceled package
	err = p.HandleProcess()
	if err == nil {
		t.Error("Should not be able to process a canceled package")
	}
	if err != ErrOrderCanceled {
		t.Errorf("Expected ErrOrderCanceled, got %v", err)
	}

	// Test cancellation from Processing state
	p2 := NewPackage("PKG102", "Cancellation Test 2")
	InitializePackage(p2)
	
	err = p2.HandleProcess() // Ordered -> Processing
	if err != nil {
		t.Fatalf("Failed to process package: %v", err)
	}
	
	err = p2.HandleCancel() // Processing -> Canceled
	if err != nil {
		t.Fatalf("Failed to cancel package: %v", err)
	}
	if p2.GetState() != "Canceled" {
		t.Fatalf("Package should be in Canceled state, got %s", p2.GetState())
	}

	// Cannot cancel an already canceled package
	err = p2.HandleCancel()
	if err == nil {
		t.Error("Should not be able to cancel an already canceled package")
	}
}

// TestReturnFlow tests the return workflow
func TestReturnFlow(t *testing.T) {
	p := NewPackage("PKG103", "Return Test")
	InitializePackage(p)
	
	// Progress package to shipped state
	err := p.HandleProcess() // Ordered -> Processing
	if err != nil {
		t.Fatalf("Failed to process package: %v", err)
	}
	
	err = p.HandleShip() // Processing -> Shipped
	if err != nil {
		t.Fatalf("Failed to ship package: %v", err)
	}
	
	// Return the package
	err = p.HandleReturn() // Shipped -> Returned
	if err != nil {
		t.Fatalf("Failed to return package: %v", err)
	}
	if p.GetState() != "Returned" {
		t.Fatalf("Package should be in Returned state, got %s", p.GetState())
	}
	
	// Process the returned package
	err = p.HandleProcess() // Returned -> Processing
	if err != nil {
		t.Fatalf("Failed to process returned package: %v", err)
	}
	if p.GetState() != "Processing" {
		t.Fatalf("Package should be in Processing state, got %s", p.GetState())
	}
}

// TestEventHandlers tests that event handlers are called on state transitions
func TestEventHandlers(t *testing.T) {
	p := NewPackage("PKG104", "Event Handler Test")
	InitializePackage(p)
	
	// Track handler calls
	handlerCalls := 0
	
	// Add a handler
	p.AddTransitionHandler(func(e Event) {
		handlerCalls++
	})
	
	// Perform a transition
	err := p.HandleProcess() // Ordered -> Processing
	if err != nil {
		t.Fatalf("Failed to process package: %v", err)
	}
	
	// Check that the handler was called
	if handlerCalls != 1 {
		t.Errorf("Handler should have been called once, was called %d times", handlerCalls)
	}
	
	// Perform another transition
	err = p.HandleShip() // Processing -> Shipped
	if err != nil {
		t.Fatalf("Failed to ship package: %v", err)
	}
	
	// Check that the handler was called again
	if handlerCalls != 2 {
		t.Errorf("Handler should have been called twice, was called %d times", handlerCalls)
	}
}

// TestLoggingHandler tests the LoggingHandler implementation
func TestLoggingHandler(t *testing.T) {
	p := NewPackage("PKG105", "Logging Handler Test")
	InitializePackage(p)
	
	// Collect log messages
	logMessages := []string{}
	
	// Add a logging handler
	p.AddTransitionHandler(LoggingHandler(func(msg string) {
		logMessages = append(logMessages, msg)
	}))
	
	// Perform a transition
	err := p.HandleProcess() // Ordered -> Processing
	if err != nil {
		t.Fatalf("Failed to process package: %v", err)
	}
	
	// Check that we got a log message
	if len(logMessages) != 1 {
		t.Fatalf("Expected 1 log message, got %d", len(logMessages))
	}
	
	// Check the log message content
	if !strings.Contains(logMessages[0], "Ordered -> Processing") {
		t.Errorf("Log message should contain 'Ordered -> Processing', got: %s", logMessages[0])
	}
}

// TestForceTransition tests the force transition functionality
func TestForceTransition(t *testing.T) {
	p := NewPackage("PKG106", "Force Transition Test")
	InitializePackage(p)
	
	// Try to force an invalid transition
	err := p.ForceTransitionTo("Delivered")
	if err != nil {
		t.Fatalf("ForceTransitionTo should succeed even for invalid transitions: %v", err)
	}
	
	// Check that the transition happened
	if p.GetState() != "Delivered" {
		t.Fatalf("Package should be in Delivered state, got %s", p.GetState())
	}
	
	// Check that the history entry indicates this was forced
	lastEntry := p.History[len(p.History)-1]
	if !strings.Contains(lastEntry.Details, "FORCED") {
		t.Errorf("History entry should indicate this was a forced transition")
	}
}

// TestTimeoutTransition tests the automatic timeout transition
func TestTimeoutTransition(t *testing.T) {
	p := NewPackage("PKG107", "Timeout Test")
	InitializePackage(p)
	
	// Progress to processing
	err := p.HandleProcess() // Ordered -> Processing
	if err != nil {
		t.Fatalf("Failed to process package: %v", err)
	}
	
	// Set a timeout to ship after a very short time
	timeout := 50 * time.Millisecond
	timer := TimeoutTransition(p, "Shipped", timeout)
	defer timer.Stop() // Cleanup
	
	// Wait for the timeout
	time.Sleep(timeout * 2)
	
	// Check that the transition happened
	if p.GetState() != "Shipped" {
		t.Fatalf("Package should have automatically transitioned to Shipped, got %s", p.GetState())
	}
	
	// Check that the history entry indicates this was automatic
	lastEntry := p.History[len(p.History)-1]
	if !strings.Contains(lastEntry.Details, "Automatic") {
		t.Errorf("History entry should indicate this was an automatic transition")
	}
}

// TestMetadataTracking tests that metadata is tracked properly
func TestMetadataTracking(t *testing.T) {
	p := NewPackage("PKG108", "Metadata Test")
	InitializePackage(p)
	
	// Check that initial metadata is set
	if _, found := p.Metadata["ordered_time"]; !found {
		t.Error("ordered_time should be set in metadata")
	}
	
	// Process the package
	err := p.HandleProcess()
	if err != nil {
		t.Fatalf("Failed to process package: %v", err)
	}
	
	// Check that processing metadata is set
	if _, found := p.Metadata["processing_time"]; !found {
		t.Error("processing_time should be set in metadata")
	}
	
	// Ship the package
	err = p.HandleShip()
	if err != nil {
		t.Fatalf("Failed to ship package: %v", err)
	}
	
	// Check that shipping metadata is set
	if _, found := p.Metadata["shipped_time"]; !found {
		t.Error("shipped_time should be set in metadata")
	}
}
