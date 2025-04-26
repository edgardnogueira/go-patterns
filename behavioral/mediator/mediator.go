// Package mediator implements the Mediator design pattern.
//
// The Mediator pattern defines an object that encapsulates how a set of objects interact.
// It promotes loose coupling by keeping objects from referring to each other explicitly,
// and it lets you vary their interaction independently.
package mediator

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

// MessageType defines the types of messages that can be sent in the system
type MessageType int

const (
	// LandingRequest is sent when an aircraft wants to land
	LandingRequest MessageType = iota
	// TakeoffRequest is sent when an aircraft wants to take off
	TakeoffRequest
	// EmergencyMessage is sent when an aircraft has an emergency
	EmergencyMessage
	// PositionUpdate is sent periodically to update aircraft position
	PositionUpdate
	// ControlMessage is sent from the control tower to aircraft
	ControlMessage
)

// Message represents a communication between aircraft and control tower
type Message struct {
	// Type of the message
	Type MessageType
	// From is the sender's ID
	From string
	// To is the recipient's ID (empty if broadcast)
	To string
	// Content is the message body
	Content string
	// Priority of the message (higher number = higher priority)
	Priority int
	// Timestamp when the message was created
	Timestamp time.Time
}

// String returns a string representation of a message
func (m Message) String() string {
	var typeStr string
	switch m.Type {
	case LandingRequest:
		typeStr = "LANDING REQUEST"
	case TakeoffRequest:
		typeStr = "TAKEOFF REQUEST"
	case EmergencyMessage:
		typeStr = "EMERGENCY"
	case PositionUpdate:
		typeStr = "POSITION UPDATE"
	case ControlMessage:
		typeStr = "CONTROL MESSAGE"
	default:
		typeStr = "UNKNOWN"
	}

	to := "ALL"
	if m.To != "" {
		to = m.To
	}

	return fmt.Sprintf("[%s] %s -> %s (Priority: %d): %s",
		typeStr, m.From, to, m.Priority, m.Content)
}

// Colleague is an interface that all participants in the mediation must implement
type Colleague interface {
	// GetID returns the unique identifier for this colleague
	GetID() string
	// GetStatus returns the current status information
	GetStatus() string
	// SendMessage sends a message to the mediator
	SendMessage(msg Message)
	// ReceiveMessage receives a message from the mediator
	ReceiveMessage(msg Message)
	// SetMediator associates a mediator with this colleague
	SetMediator(mediator Mediator)
}

// Mediator is the interface that defines how colleagues interact
type Mediator interface {
	// Register adds a colleague to the mediation
	Register(colleague Colleague)
	// Unregister removes a colleague from the mediation
	Unregister(colleague Colleague)
	// Send delivers a message to the appropriate recipient(s)
	Send(msg Message)
	// Broadcast sends a message to all registered colleagues
	Broadcast(from string, msgType MessageType, content string, priority int)
	// GetColleague returns a colleague by ID
	GetColleague(id string) (Colleague, bool)
}

// Position represents coordinates in 2D space
type Position struct {
	X int
	Y int
}

// String returns a string representation of a position
func (p Position) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

// Aircraft is the base struct for all aircraft types
type Aircraft struct {
	// ID is the unique identifier for the aircraft (e.g., flight number)
	ID string
	// Type is the aircraft model
	Type string
	// Position is the current coordinates
	Position Position
	// Altitude in feet
	Altitude int
	// IsFlying indicates if the aircraft is in the air
	IsFlying bool
	// InEmergency indicates if the aircraft is experiencing an emergency
	InEmergency bool
	// mediator is the control tower this aircraft communicates with
	mediator Mediator
	// messageLog stores messages received by this aircraft
	messageLog []Message
	// status is the current operational status
	status string
}

// GetID returns the aircraft's identifier
func (a *Aircraft) GetID() string {
	return a.ID
}

// GetStatus returns the aircraft's current status
func (a *Aircraft) GetStatus() string {
	return a.status
}

// SetStatus updates the aircraft's status
func (a *Aircraft) SetStatus(status string) {
	a.status = status
}

// SetMediator associates a mediator with this aircraft
func (a *Aircraft) SetMediator(mediator Mediator) {
	a.mediator = mediator
}

// SendMessage sends a message to the mediator
func (a *Aircraft) SendMessage(msg Message) {
	if a.mediator != nil {
		msg.From = a.ID
		msg.Timestamp = time.Now()
		a.mediator.Send(msg)
	}
}

// ReceiveMessage processes a message from the mediator
func (a *Aircraft) ReceiveMessage(msg Message) {
	// Log the message
	a.messageLog = append(a.messageLog, msg)

	// Process based on message type
	switch msg.Type {
	case ControlMessage:
		if strings.Contains(strings.ToLower(msg.Content), "landing clearance granted") {
			a.SetStatus("Preparing for landing")
		} else if strings.Contains(strings.ToLower(msg.Content), "takeoff clearance granted") {
			a.SetStatus("Taking off")
			a.IsFlying = true
		}
	}
}

// GetMessageLog returns the message history for this aircraft
func (a *Aircraft) GetMessageLog() []Message {
	return a.messageLog
}

// RequestLanding sends a landing request to the control tower
func (a *Aircraft) RequestLanding() {
	if !a.IsFlying {
		return
	}
	priority := 5 // Normal priority
	if a.InEmergency {
		priority = 10 // Emergency priority
	}
	a.SendMessage(Message{
		Type:     LandingRequest,
		Content:  fmt.Sprintf("Aircraft %s requesting landing at altitude %d, position %s", a.ID, a.Altitude, a.Position),
		Priority: priority,
	})
	a.SetStatus("Awaiting landing clearance")
}

// RequestTakeoff sends a takeoff request to the control tower
func (a *Aircraft) RequestTakeoff() {
	if a.IsFlying {
		return
	}
	a.SendMessage(Message{
		Type:     TakeoffRequest,
		Content:  fmt.Sprintf("Aircraft %s requesting takeoff from position %s", a.ID, a.Position),
		Priority: 5,
	})
	a.SetStatus("Awaiting takeoff clearance")
}

// ReportEmergency sends an emergency message to the control tower
func (a *Aircraft) ReportEmergency(details string) {
	a.InEmergency = true
	a.SendMessage(Message{
		Type:     EmergencyMessage,
		Content:  fmt.Sprintf("EMERGENCY: %s at altitude %d, position %s: %s", a.ID, a.Altitude, a.Position, details),
		Priority: 10,
	})
	a.SetStatus("Emergency: " + details)
}

// UpdatePosition sends the current position to the control tower
func (a *Aircraft) UpdatePosition(x, y, altitude int) {
	a.Position.X = x
	a.Position.Y = y
	a.Altitude = altitude
	a.SendMessage(Message{
		Type:     PositionUpdate,
		Content:  fmt.Sprintf("Aircraft %s at position %s, altitude %d", a.ID, a.Position, a.Altitude),
		Priority: 3,
	})
}

// Land completes the landing procedure
func (a *Aircraft) Land() {
	if !a.IsFlying {
		return
	}
	a.IsFlying = false
	a.Altitude = 0
	a.SetStatus("Landed")
	a.SendMessage(Message{
		Type:     PositionUpdate,
		Content:  fmt.Sprintf("Aircraft %s has landed at position %s", a.ID, a.Position),
		Priority: 4,
	})
}

// TakeOff completes the takeoff procedure
func (a *Aircraft) TakeOff(targetAltitude int) {
	if a.IsFlying {
		return
	}
	a.IsFlying = true
	a.Altitude = targetAltitude
	a.SetStatus("In air")
	a.SendMessage(Message{
		Type:     PositionUpdate,
		Content:  fmt.Sprintf("Aircraft %s has taken off, climbing to %d feet", a.ID, targetAltitude),
		Priority: 4,
	})
}

// PassengerAircraft represents a commercial passenger plane
type PassengerAircraft struct {
	Aircraft
	// PassengerCount is the number of passengers on board
	PassengerCount int
	// Airline is the operating airline name
	Airline string
}

// NewPassengerAircraft creates a new passenger aircraft
func NewPassengerAircraft(id, airline string, passengerCount int) *PassengerAircraft {
	return &PassengerAircraft{
		Aircraft: Aircraft{
			ID:          id,
			Type:        "Passenger",
			IsFlying:    false,
			InEmergency: false,
			status:      "Idle",
		},
		PassengerCount: passengerCount,
		Airline:        airline,
	}
}

// CargoAircraft represents a cargo transport plane
type CargoAircraft struct {
	Aircraft
	// CargoWeight in kg
	CargoWeight float64
	// Company is the cargo company operating the aircraft
	Company string
}

// NewCargoAircraft creates a new cargo aircraft
func NewCargoAircraft(id, company string, cargoWeight float64) *CargoAircraft {
	return &CargoAircraft{
		Aircraft: Aircraft{
			ID:          id,
			Type:        "Cargo",
			IsFlying:    false,
			InEmergency: false,
			status:      "Idle",
		},
		CargoWeight: cargoWeight,
		Company:     company,
	}
}

// PrivateAircraft represents a private jet or small plane
type PrivateAircraft struct {
	Aircraft
	// Owner is the owner's name
	Owner string
}

// NewPrivateAircraft creates a new private aircraft
func NewPrivateAircraft(id, owner string) *PrivateAircraft {
	return &PrivateAircraft{
		Aircraft: Aircraft{
			ID:          id,
			Type:        "Private",
			IsFlying:    false,
			InEmergency: false,
			status:      "Idle",
		},
		Owner: owner,
	}
}

// MilitaryAircraft represents a military aircraft with special protocols
type MilitaryAircraft struct {
	Aircraft
	// Mission is the current mission identifier
	Mission string
	// Branch is the military branch (e.g., "Air Force")
	Branch string
}

// NewMilitaryAircraft creates a new military aircraft
func NewMilitaryAircraft(id, branch, mission string) *MilitaryAircraft {
	return &MilitaryAircraft{
		Aircraft: Aircraft{
			ID:          id,
			Type:        "Military",
			IsFlying:    false,
			InEmergency: false,
			status:      "Idle",
		},
		Mission: mission,
		Branch:  branch,
	}
}

// ReceiveMessage extends the base ReceiveMessage with military protocols
func (m *MilitaryAircraft) ReceiveMessage(msg Message) {
	m.Aircraft.ReceiveMessage(msg)
	// Military aircraft might have special handling for certain messages
	if msg.Type == EmergencyMessage {
		m.SetStatus("Alert: Supporting emergency situation")
	}
}

// AirTrafficControl is a concrete mediator that coordinates aircraft
type AirTrafficControl struct {
	// Name of the airport or control region
	Name string
	// colleagues maps aircraft ID to the aircraft
	colleagues map[string]Colleague
	// messageLog stores all messages processed by the control tower
	messageLog []Message
	// mutex protects concurrent access to the mediator
	mutex sync.RWMutex
}

// NewAirTrafficControl creates a new control tower
func NewAirTrafficControl(name string) *AirTrafficControl {
	return &AirTrafficControl{
		Name:       name,
		colleagues: make(map[string]Colleague),
		messageLog: []Message{},
	}
}

// Register adds an aircraft to the control system
func (atc *AirTrafficControl) Register(colleague Colleague) {
	atc.mutex.Lock()
	defer atc.mutex.Unlock()

	if colleague == nil {
		return
	}

	id := colleague.GetID()
	if _, exists := atc.colleagues[id]; !exists {
		atc.colleagues[id] = colleague
		colleague.SetMediator(atc)
		log.Printf("%s: Registered new aircraft %s\n", atc.Name, id)
	}
}

// Unregister removes an aircraft from the control system
func (atc *AirTrafficControl) Unregister(colleague Colleague) {
	atc.mutex.Lock()
	defer atc.mutex.Unlock()

	if colleague == nil {
		return
	}

	id := colleague.GetID()
	if _, exists := atc.colleagues[id]; exists {
		delete(atc.colleagues, id)
		log.Printf("%s: Unregistered aircraft %s\n", atc.Name, id)
	}
}

// Send processes and delivers a message
func (atc *AirTrafficControl) Send(msg Message) {
	atc.mutex.Lock()
	defer atc.mutex.Unlock()

	// Log the message
	log.Printf("%s: %s\n", atc.Name, msg.String())
	atc.messageLog = append(atc.messageLog, msg)

	// Process the message based on type
	switch msg.Type {
	case LandingRequest:
		// Handle landing request
		atc.handleLandingRequest(msg)
	case TakeoffRequest:
		// Handle takeoff request
		atc.handleTakeoffRequest(msg)
	case EmergencyMessage:
		// Handle emergency with highest priority
		atc.handleEmergency(msg)
	case PositionUpdate:
		// Just store position updates in the log
	default:
		// Forward message to specific recipient if specified
		if msg.To != "" {
			if recipient, ok := atc.colleagues[msg.To]; ok {
				recipient.ReceiveMessage(msg)
			}
		}
	}
}

// handleLandingRequest processes a landing request
func (atc *AirTrafficControl) handleLandingRequest(msg Message) {
	// In a real system, this would check runway availability, weather, etc.
	response := Message{
		Type:      ControlMessage,
		From:      atc.Name,
		To:        msg.From,
		Content:   fmt.Sprintf("Landing clearance granted for %s. Proceed to runway 27.", msg.From),
		Priority:  msg.Priority,
		Timestamp: time.Now(),
	}

	if recipient, ok := atc.colleagues[msg.From]; ok {
		recipient.ReceiveMessage(response)
	}
}

// handleTakeoffRequest processes a takeoff request
func (atc *AirTrafficControl) handleTakeoffRequest(msg Message) {
	// In a real system, this would check runway availability, weather, etc.
	response := Message{
		Type:      ControlMessage,
		From:      atc.Name,
		To:        msg.From,
		Content:   fmt.Sprintf("Takeoff clearance granted for %s. Proceed to runway 09.", msg.From),
		Priority:  msg.Priority,
		Timestamp: time.Now(),
	}

	if recipient, ok := atc.colleagues[msg.From]; ok {
		recipient.ReceiveMessage(response)
	}
}

// handleEmergency processes an emergency message and notifies all aircraft
func (atc *AirTrafficControl) handleEmergency(msg Message) {
	// Alert all aircraft about the emergency
	alert := Message{
		Type:      ControlMessage,
		From:      atc.Name,
		Content:   fmt.Sprintf("EMERGENCY ALERT: %s has declared an emergency. All aircraft maintain positions.", msg.From),
		Priority:  10, // Highest priority
		Timestamp: time.Now(),
	}

	for id, colleague := range atc.colleagues {
		if id != msg.From { // Don't send to the aircraft in emergency
			colleague.ReceiveMessage(alert)
		}
	}

	// Send specific instructions to the aircraft in emergency
	response := Message{
		Type:      ControlMessage,
		From:      atc.Name,
		To:        msg.From,
		Content:   "Emergency acknowledged. You have priority clearance. All runways being cleared.",
		Priority:  10,
		Timestamp: time.Now(),
	}

	if recipient, ok := atc.colleagues[msg.From]; ok {
		recipient.ReceiveMessage(response)
	}
}

// Broadcast sends a message to all registered colleagues
func (atc *AirTrafficControl) Broadcast(from string, msgType MessageType, content string, priority int) {
	atc.mutex.RLock()
	defer atc.mutex.RUnlock()

	msg := Message{
		Type:      msgType,
		From:      from,
		Content:   content,
		Priority:  priority,
		Timestamp: time.Now(),
	}

	// Log the broadcast message
	log.Printf("%s: BROADCAST: %s\n", atc.Name, msg.String())
	atc.messageLog = append(atc.messageLog, msg)

	// Send to all colleagues
	for _, colleague := range atc.colleagues {
		colleague.ReceiveMessage(msg)
	}
}

// GetColleague returns a colleague by ID
func (atc *AirTrafficControl) GetColleague(id string) (Colleague, bool) {
	atc.mutex.RLock()
	defer atc.mutex.RUnlock()

	colleague, exists := atc.colleagues[id]
	return colleague, exists
}

// GetMessageLog returns the control tower's message log
func (atc *AirTrafficControl) GetMessageLog() []Message {
	atc.mutex.RLock()
	defer atc.mutex.RUnlock()

	// Return a copy to avoid concurrent access issues
	logCopy := make([]Message, len(atc.messageLog))
	copy(logCopy, atc.messageLog)
	return logCopy
}

// GetRegisteredAircraft returns a list of all registered aircraft IDs
func (atc *AirTrafficControl) GetRegisteredAircraft() []string {
	atc.mutex.RLock()
	defer atc.mutex.RUnlock()

	var aircraftIDs []string
	for id := range atc.colleagues {
		aircraftIDs = append(aircraftIDs, id)
	}
	return aircraftIDs
}
