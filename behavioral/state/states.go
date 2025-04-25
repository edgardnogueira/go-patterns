package state

import (
	"errors"
	"fmt"
)

// Common error messages
var (
	ErrInvalidTransition = errors.New("invalid state transition")
	ErrAlreadyInState    = errors.New("package already in this state")
	ErrOrderCanceled     = errors.New("order has been canceled")
)

// BaseState provides common functionality for all states
type BaseState struct {
	StateName string
}

// Name returns the name of the state
func (s *BaseState) Name() string {
	return s.StateName
}

// Default implementations that will be overridden by specific states as needed
func (s *BaseState) Process() error {
	return ErrInvalidTransition
}

func (s *BaseState) Ship() error {
	return ErrInvalidTransition
}

func (s *BaseState) Deliver() error {
	return ErrInvalidTransition
}

func (s *BaseState) Return() error {
	return ErrInvalidTransition
}

func (s *BaseState) Cancel() error {
	return ErrInvalidTransition
}

func (s *BaseState) AllowedTransitions() []string {
	return []string{}
}

func (s *BaseState) Enter(p *Package) {
	// Default implementation does nothing
}

func (s *BaseState) Exit(p *Package) {
	// Default implementation does nothing
}

// OrderedState represents the initial state of a package after it's ordered
type OrderedState struct {
	BaseState
}

// NewOrderedState creates a new OrderedState
func NewOrderedState() *OrderedState {
	return &OrderedState{BaseState{StateName: "Ordered"}}
}

// Process transitions from Ordered to Processing
func (s *OrderedState) Process() error {
	return nil // Valid transition, will be handled by the context
}

// Cancel transitions from Ordered to Canceled
func (s *OrderedState) Cancel() error {
	return nil // Valid transition, will be handled by the context
}

// AllowedTransitions returns the states that OrderedState can transition to
func (s *OrderedState) AllowedTransitions() []string {
	return []string{"Processing", "Canceled"}
}

// Enter is called when a package enters the Ordered state
func (s *OrderedState) Enter(p *Package) {
	p.Metadata["ordered_time"] = p.LastUpdatedAt
}

// ProcessingState represents a package being processed in the warehouse
type ProcessingState struct {
	BaseState
}

// NewProcessingState creates a new ProcessingState
func NewProcessingState() *ProcessingState {
	return &ProcessingState{BaseState{StateName: "Processing"}}
}

// Ship transitions from Processing to Shipped
func (s *ProcessingState) Ship() error {
	return nil // Valid transition, will be handled by the context
}

// Cancel transitions from Processing to Canceled
func (s *ProcessingState) Cancel() error {
	return nil // Valid transition, will be handled by the context
}

// AllowedTransitions returns the states that ProcessingState can transition to
func (s *ProcessingState) AllowedTransitions() []string {
	return []string{"Shipped", "Canceled"}
}

// Enter is called when a package enters the Processing state
func (s *ProcessingState) Enter(p *Package) {
	p.Metadata["processing_time"] = p.LastUpdatedAt
}

// ShippedState represents a package in transit
type ShippedState struct {
	BaseState
}

// NewShippedState creates a new ShippedState
func NewShippedState() *ShippedState {
	return &ShippedState{BaseState{StateName: "Shipped"}}
}

// Deliver transitions from Shipped to Delivered
func (s *ShippedState) Deliver() error {
	return nil // Valid transition, will be handled by the context
}

// Return transitions from Shipped to Returned
func (s *ShippedState) Return() error {
	return nil // Valid transition, will be handled by the context
}

// AllowedTransitions returns the states that ShippedState can transition to
func (s *ShippedState) AllowedTransitions() []string {
	return []string{"Delivered", "Returned"}
}

// Enter is called when a package enters the Shipped state
func (s *ShippedState) Enter(p *Package) {
	p.Metadata["shipped_time"] = p.LastUpdatedAt
}

// DeliveredState represents a package that has been delivered
type DeliveredState struct {
	BaseState
}

// NewDeliveredState creates a new DeliveredState
func NewDeliveredState() *DeliveredState {
	return &DeliveredState{BaseState{StateName: "Delivered"}}
}

// Return transitions from Delivered to Returned
func (s *DeliveredState) Return() error {
	return nil // Valid transition, will be handled by the context
}

// AllowedTransitions returns the states that DeliveredState can transition to
func (s *DeliveredState) AllowedTransitions() []string {
	return []string{"Returned"}
}

// Enter is called when a package enters the Delivered state
func (s *DeliveredState) Enter(p *Package) {
	p.Metadata["delivered_time"] = p.LastUpdatedAt
}

// ReturnedState represents a package that is being returned
type ReturnedState struct {
	BaseState
}

// NewReturnedState creates a new ReturnedState
func NewReturnedState() *ReturnedState {
	return &ReturnedState{BaseState{StateName: "Returned"}}
}

// Process transitions from Returned to Processing (for reprocessing)
func (s *ReturnedState) Process() error {
	return nil // Valid transition, will be handled by the context
}

// AllowedTransitions returns the states that ReturnedState can transition to
func (s *ReturnedState) AllowedTransitions() []string {
	return []string{"Processing"}
}

// Enter is called when a package enters the Returned state
func (s *ReturnedState) Enter(p *Package) {
	p.Metadata["returned_time"] = p.LastUpdatedAt
}

// CanceledState represents a canceled order
type CanceledState struct {
	BaseState
}

// NewCanceledState creates a new CanceledState
func NewCanceledState() *CanceledState {
	return &CanceledState{BaseState{StateName: "Canceled"}}
}

// AllowedTransitions returns the states that CanceledState can transition to
// Canceled is a terminal state, so no transitions are allowed
func (s *CanceledState) AllowedTransitions() []string {
	return []string{}
}

// Enter is called when a package enters the Canceled state
func (s *CanceledState) Enter(p *Package) {
	p.Metadata["canceled_time"] = p.LastUpdatedAt
	p.Metadata["canceled_reason"] = "User canceled order" // Default reason
}

// Override all actions to return ErrOrderCanceled
func (s *CanceledState) Process() error {
	return ErrOrderCanceled
}

func (s *CanceledState) Ship() error {
	return ErrOrderCanceled
}

func (s *CanceledState) Deliver() error {
	return ErrOrderCanceled
}

func (s *CanceledState) Return() error {
	return ErrOrderCanceled
}

func (s *CanceledState) Cancel() error {
	return ErrAlreadyInState
}

// StateFactory creates the appropriate state based on the state name
func StateFactory(stateName string) (PackageState, error) {
	switch stateName {
	case "Ordered":
		return NewOrderedState(), nil
	case "Processing":
		return NewProcessingState(), nil
	case "Shipped":
		return NewShippedState(), nil
	case "Delivered":
		return NewDeliveredState(), nil
	case "Returned":
		return NewReturnedState(), nil
	case "Canceled":
		return NewCanceledState(), nil
	default:
		return nil, fmt.Errorf("unknown state: %s", stateName)
	}
}

// InitializePackage sets the initial state of a newly created package
func InitializePackage(p *Package) {
	// Set initial state to Ordered
	initialState := NewOrderedState()
	p.CurrentState = initialState
	
	// Record the initial state in history
	event := Event{
		From:      "",
		To:        initialState.Name(),
		Timestamp: p.CreatedAt,
		Details:   "Package created in Ordered state",
	}
	p.History = append(p.History, event)
	
	// Call Enter on the initial state
	initialState.Enter(p)
}
