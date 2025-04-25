// Package state implements the State design pattern in Go.
//
// The State pattern allows an object to alter its behavior when its internal state changes.
// This implementation demonstrates a package delivery system with different states 
// (ordered, processing, shipped, delivered, returned, canceled) that change the
// behavior of the package object.
package state

import (
	"errors"
	"fmt"
	"time"
)

// Event represents a state transition event
type Event struct {
	From      string
	To        string
	Timestamp time.Time
	Details   string
}

// TransitionHandler is a function that gets called when a state transition occurs
type TransitionHandler func(e Event)

// PackageState is the interface that defines the behavior for all package states
type PackageState interface {
	// Name returns the name of the state
	Name() string

	// Process attempts to process the package in the current state
	Process() error

	// Ship attempts to ship the package in the current state
	Ship() error

	// Deliver attempts to deliver the package in the current state
	Deliver() error

	// Return attempts to return the package in the current state
	Return() error

	// Cancel attempts to cancel the order in the current state
	Cancel() error

	// AllowedTransitions returns the list of states this state can transition to
	AllowedTransitions() []string

	// Enter is called when the package enters this state
	Enter(*Package)

	// Exit is called when the package exits this state
	Exit(*Package)
}

// Package is the context that maintains the current state and delegates
// state-specific behavior to the current state object
type Package struct {
	// ID is the package identifier
	ID string

	// Description is the package description
	Description string

	// CurrentState holds the current state of the package
	CurrentState PackageState

	// TransitionHandlers holds registered handlers for state transitions
	TransitionHandlers []TransitionHandler

	// History holds the state transition history
	History []Event

	// Metadata holds additional package information
	Metadata map[string]interface{}

	// CreatedAt is the package creation timestamp
	CreatedAt time.Time

	// LastUpdatedAt is the timestamp of the last state change
	LastUpdatedAt time.Time
}

// NewPackage creates a new package with the initial ordered state
func NewPackage(id, description string) *Package {
	now := time.Now()
	p := &Package{
		ID:                id,
		Description:       description,
		TransitionHandlers: []TransitionHandler{},
		History:           []Event{},
		Metadata:          make(map[string]interface{}),
		CreatedAt:         now,
		LastUpdatedAt:     now,
	}

	// Initialize with OrderedState (will be defined later)
	// The actual state assignment happens when implementing OrderedState
	return p
}

// SetState changes the current state of the package
func (p *Package) SetState(newState PackageState) {
	if p.CurrentState != nil {
		event := Event{
			From:      p.CurrentState.Name(),
			To:        newState.Name(),
			Timestamp: time.Now(),
			Details:   fmt.Sprintf("State changed from %s to %s", p.CurrentState.Name(), newState.Name()),
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
	}

	// Set new state and call enter
	p.CurrentState = newState
	p.CurrentState.Enter(p)
}

// AddTransitionHandler registers a handler for state transitions
func (p *Package) AddTransitionHandler(handler TransitionHandler) {
	p.TransitionHandlers = append(p.TransitionHandlers, handler)
}

// GetStateHistory returns the state transition history
func (p *Package) GetStateHistory() []Event {
	return p.History
}

// Process delegates the process action to the current state
func (p *Package) Process() error {
	if p.CurrentState == nil {
		return errors.New("package has no state")
	}
	return p.CurrentState.Process()
}

// Ship delegates the ship action to the current state
func (p *Package) Ship() error {
	if p.CurrentState == nil {
		return errors.New("package has no state")
	}
	return p.CurrentState.Ship()
}

// Deliver delegates the deliver action to the current state
func (p *Package) Deliver() error {
	if p.CurrentState == nil {
		return errors.New("package has no state")
	}
	return p.CurrentState.Deliver()
}

// Return delegates the return action to the current state
func (p *Package) Return() error {
	if p.CurrentState == nil {
		return errors.New("package has no state")
	}
	return p.CurrentState.Return()
}

// Cancel delegates the cancel action to the current state
func (p *Package) Cancel() error {
	if p.CurrentState == nil {
		return errors.New("package has no state")
	}
	return p.CurrentState.Cancel()
}

// GetState returns the current state name
func (p *Package) GetState() string {
	if p.CurrentState == nil {
		return "undefined"
	}
	return p.CurrentState.Name()
}

// GetAllowedTransitions returns the list of allowed state transitions
func (p *Package) GetAllowedTransitions() []string {
	if p.CurrentState == nil {
		return []string{}
	}
	return p.CurrentState.AllowedTransitions()
}
