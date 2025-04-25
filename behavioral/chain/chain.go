// Package chain implements the Chain of Responsibility design pattern in Go.
//
// The Chain of Responsibility pattern passes a request along a chain of handlers.
// Upon receiving a request, each handler decides either to process the request
// or to pass it to the next handler in the chain.
// This implementation demonstrates a support ticket system where tickets are
// processed based on their priority and type through different support levels.
package chain

import (
	"fmt"
	"time"
)

// Priority defines the urgency level of a support ticket
type Priority int

// Priority levels for support tickets
const (
	Low Priority = iota
	Medium
	High
	Critical
)

// String returns the string representation of a priority
func (p Priority) String() string {
	switch p {
	case Low:
		return "Low"
	case Medium:
		return "Medium"
	case High:
		return "High"
	case Critical:
		return "Critical"
	default:
		return "Unknown"
	}
}

// TicketType defines the category of a support ticket
type TicketType string

// Ticket types
const (
	General      TicketType = "General"
	Technical    TicketType = "Technical"
	Billing      TicketType = "Billing"
	FeatureRequest TicketType = "Feature Request"
	Bug          TicketType = "Bug"
	Security     TicketType = "Security"
	Complaint    TicketType = "Complaint"
)

// SupportTicket represents a customer support request
type SupportTicket struct {
	ID          string
	Type        TicketType
	Priority    Priority
	Subject     string
	Description string
	CreatedAt   time.Time
	ResolvedBy  string
	Resolution  string
	IsResolved  bool
	Metadata    map[string]interface{}
}

// NewSupportTicket creates a new support ticket with the given details
func NewSupportTicket(id string, ticketType TicketType, priority Priority, subject, description string) *SupportTicket {
	return &SupportTicket{
		ID:          id,
		Type:        ticketType,
		Priority:    priority,
		Subject:     subject,
		Description: description,
		CreatedAt:   time.Now(),
		IsResolved:  false,
		Metadata:    make(map[string]interface{}),
	}
}

// SetResolution marks the ticket as resolved with the given resolution
func (t *SupportTicket) SetResolution(handler string, resolution string) {
	t.IsResolved = true
	t.ResolvedBy = handler
	t.Resolution = resolution
}

// String returns a string representation of the ticket
func (t *SupportTicket) String() string {
	status := "Unresolved"
	if t.IsResolved {
		status = fmt.Sprintf("Resolved by %s", t.ResolvedBy)
	}
	
	return fmt.Sprintf("Ticket #%s [%s, %s] - %s\nStatus: %s\nDescription: %s", 
		t.ID, t.Type, t.Priority, t.Subject, status, t.Description)
}

// Handler defines the interface for processing support tickets
type Handler interface {
	// SetNext sets the next handler in the chain
	SetNext(handler Handler) Handler
	
	// GetNext returns the next handler in the chain
	GetNext() Handler
	
	// Handle processes the support ticket or passes it to the next handler
	Handle(ticket *SupportTicket) bool
	
	// Name returns the name of this handler
	Name() string
}

// BaseHandler provides common functionality for all handlers
type BaseHandler struct {
	next Handler
	name string
}

// SetNext sets the next handler in the chain
func (h *BaseHandler) SetNext(next Handler) Handler {
	h.next = next
	return next
}

// GetNext returns the next handler in the chain
func (h *BaseHandler) GetNext() Handler {
	return h.next
}

// Handle passes the request to the next handler
func (h *BaseHandler) Handle(ticket *SupportTicket) bool {
	if h.next != nil {
		return h.next.Handle(ticket)
	}
	return false
}

// Name returns the name of the handler
func (h *BaseHandler) Name() string {
	return h.name
}

// Chain represents a chain of responsibility for handling support tickets
type Chain struct {
	head Handler
	tail Handler
}

// NewChain creates a new chain with the provided handler as head
func NewChain(head Handler) *Chain {
	return &Chain{
		head: head,
		tail: head,
	}
}

// AddHandler adds a new handler to the end of the chain
func (c *Chain) AddHandler(handler Handler) {
	if c.head == nil {
		c.head = handler
		c.tail = handler
		return
	}
	
	c.tail.SetNext(handler)
	c.tail = handler
}

// InsertHandler inserts a handler after the handler with the specified name
func (c *Chain) InsertHandler(after string, newHandler Handler) bool {
	if c.head == nil {
		c.head = newHandler
		c.tail = newHandler
		return true
	}
	
	current := c.head
	for current != nil {
		if current.Name() == after {
			newHandler.SetNext(current.GetNext())
			current.SetNext(newHandler)
			
			// If we're inserting after the tail, update the tail
			if current == c.tail {
				c.tail = newHandler
			}
			
			return true
		}
		current = current.GetNext()
	}
	
	return false
}

// RemoveHandler removes a handler from the chain
func (c *Chain) RemoveHandler(name string) bool {
	if c.head == nil {
		return false
	}
	
	// Special case: removing the head
	if c.head.Name() == name {
		nextHandler := c.head.GetNext()
		c.head = nextHandler
		
		// If we just removed the only handler, update tail
		if c.tail.Name() == name {
			c.tail = nil
		}
		
		return true
	}
	
	// Find the handler before the one we want to remove
	current := c.head
	for current != nil && current.GetNext() != nil {
		if current.GetNext().Name() == name {
			nextHandler := current.GetNext().GetNext()
			current.SetNext(nextHandler)
			
			// If we're removing the tail, update it
			if c.tail.Name() == name {
				c.tail = current
			}
			
			return true
		}
		current = current.GetNext()
	}
	
	return false
}

// Process passes a ticket through the chain of handlers
func (c *Chain) Process(ticket *SupportTicket) bool {
	if c.head == nil {
		return false
	}
	
	return c.head.Handle(ticket)
}
