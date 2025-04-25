package state

import (
	"errors"
	"fmt"
	"time"
)

// TransitionValidator validates if a state transition is allowed
type TransitionValidator func(from, to string) bool

// DefaultTransitionValidator validates transitions based on the allowed transitions
// defined in each state
func DefaultTransitionValidator(from, to string) bool {
	// Create a temporary instance of the 'from' state to check allowed transitions
	fromState, err := StateFactory(from)
	if err != nil {
		return false
	}

	// Check if the 'to' state is in the list of allowed transitions
	for _, allowedState := range fromState.AllowedTransitions() {
		if allowedState == to {
			return true
		}
	}
	return false
}

// StateTransitioner provides methods to transition between states
type StateTransitioner interface {
	// TransitionTo attempts to transition to the specified state
	TransitionTo(stateName string) error
	
	// TransitionToWithDetails transitions with custom details
	TransitionToWithDetails(stateName, details string) error
	
	// ForceTransitionTo transitions without validation
	ForceTransitionTo(stateName string) error
}

// Enable package to transition between states
func (p *Package) TransitionTo(stateName string) error {
	return p.TransitionToWithDetails(stateName, "")
}

// TransitionToWithDetails transitions with custom details
func (p *Package) TransitionToWithDetails(stateName, details string) error {
	if p.CurrentState == nil {
		return errors.New("package has no current state")
	}

	// Don't transition if we're already in this state
	if p.CurrentState.Name() == stateName {
		return nil
	}

	// Validate the transition
	if !DefaultTransitionValidator(p.CurrentState.Name(), stateName) {
		return fmt.Errorf("invalid transition from %s to %s", p.CurrentState.Name(), stateName)
	}

	// Get the new state
	newState, err := StateFactory(stateName)
	if err != nil {
		return err
	}

	// Set the event details
	eventDetails := details
	if eventDetails == "" {
		eventDetails = fmt.Sprintf("State changed from %s to %s", p.CurrentState.Name(), stateName)
	}

	// Create transition event
	event := Event{
		From:      p.CurrentState.Name(),
		To:        stateName,
		Timestamp: time.Now(),
		Details:   eventDetails,
	}

	// Call exit on current state
	p.CurrentState.Exit(p)

	// Notify handlers about the transition
	for _, handler := range p.TransitionHandlers {
		handler(event)
	}

	// Add to history
	p.History = append(p.History, event)
	p.LastUpdatedAt = event.Timestamp

	// Set new state and call enter
	p.CurrentState = newState
	p.CurrentState.Enter(p)

	return nil
}

// ForceTransitionTo transitions without validation
func (p *Package) ForceTransitionTo(stateName string) error {
	newState, err := StateFactory(stateName)
	if err != nil {
		return err
	}

	// Create transition event
	event := Event{
		From:      p.CurrentState.Name(),
		To:        stateName,
		Timestamp: time.Now(),
		Details:   fmt.Sprintf("FORCED state change from %s to %s", p.CurrentState.Name(), stateName),
	}

	// Call exit on current state
	if p.CurrentState != nil {
		p.CurrentState.Exit(p)
	}

	// Notify handlers about the transition
	for _, handler := range p.TransitionHandlers {
		handler(event)
	}

	// Add to history
	p.History = append(p.History, event)
	p.LastUpdatedAt = event.Timestamp

	// Set new state and call enter
	p.CurrentState = newState
	p.CurrentState.Enter(p)

	return nil
}

// HandleProcess processes the current state and transitions to the next state if appropriate
func (p *Package) HandleProcess() error {
	if p.CurrentState == nil {
		return errors.New("package has no state")
	}

	// Process the current state
	err := p.CurrentState.Process()
	if err != nil {
		return err
	}

	// If we're in Ordered state, move to Processing
	if p.CurrentState.Name() == "Ordered" {
		return p.TransitionTo("Processing")
	}

	// If we're in Returned state, move to Processing
	if p.CurrentState.Name() == "Returned" {
		return p.TransitionTo("Processing")
	}

	return nil
}

// HandleShip ships the package and transitions to the next state if appropriate
func (p *Package) HandleShip() error {
	if p.CurrentState == nil {
		return errors.New("package has no state")
	}

	// Ship the current state
	err := p.CurrentState.Ship()
	if err != nil {
		return err
	}

	// If we're in Processing state, move to Shipped
	if p.CurrentState.Name() == "Processing" {
		return p.TransitionTo("Shipped")
	}

	return nil
}

// HandleDeliver delivers the package and transitions to the next state if appropriate
func (p *Package) HandleDeliver() error {
	if p.CurrentState == nil {
		return errors.New("package has no state")
	}

	// Deliver the current state
	err := p.CurrentState.Deliver()
	if err != nil {
		return err
	}

	// If we're in Shipped state, move to Delivered
	if p.CurrentState.Name() == "Shipped" {
		return p.TransitionTo("Delivered")
	}

	return nil
}

// HandleReturn returns the package and transitions to the next state if appropriate
func (p *Package) HandleReturn() error {
	if p.CurrentState == nil {
		return errors.New("package has no state")
	}

	// Return the current state
	err := p.CurrentState.Return()
	if err != nil {
		return err
	}

	// If we're in Shipped or Delivered state, move to Returned
	if p.CurrentState.Name() == "Shipped" || p.CurrentState.Name() == "Delivered" {
		return p.TransitionTo("Returned")
	}

	return nil
}

// HandleCancel cancels the order and transitions to the next state if appropriate
func (p *Package) HandleCancel() error {
	if p.CurrentState == nil {
		return errors.New("package has no state")
	}

	// Cancel the current state
	err := p.CurrentState.Cancel()
	if err != nil {
		return err
	}

	// If we're in Ordered or Processing state, move to Canceled
	if p.CurrentState.Name() == "Ordered" || p.CurrentState.Name() == "Processing" {
		return p.TransitionTo("Canceled")
	}

	return nil
}

// Common event handlers

// LoggingHandler creates a handler that logs state transitions
func LoggingHandler(logger func(string)) TransitionHandler {
	return func(e Event) {
		logger(fmt.Sprintf("[%s] Package state change: %s -> %s (%s)",
			e.Timestamp.Format(time.RFC3339),
			e.From, e.To, e.Details))
	}
}

// NotificationHandler creates a handler that sends notifications
func NotificationHandler(notify func(string, string)) TransitionHandler {
	return func(e Event) {
		message := fmt.Sprintf("Package state changed from %s to %s", e.From, e.To)
		details := e.Details
		notify(message, details)
	}
}

// TimeoutTransition schedules an automatic state transition after a specified duration
func TimeoutTransition(p *Package, targetState string, duration time.Duration) *time.Timer {
	timer := time.NewTimer(duration)
	
	go func() {
		<-timer.C
		// When the timer fires, attempt the transition
		err := p.TransitionToWithDetails(targetState, 
			fmt.Sprintf("Automatic transition to %s after timeout of %v", targetState, duration))
		
		if err != nil {
			// Handle error (could add a logger or error handler in the package)
			// For now, just annotate the metadata
			p.Metadata["timeout_error"] = fmt.Sprintf("Failed to transition to %s: %v", targetState, err)
		}
	}()
	
	return timer
}
