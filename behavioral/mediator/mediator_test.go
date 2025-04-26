package mediator

import (
	"strings"
	"testing"
	"time"
)

// TestAircraftRegistration tests the registration and unregistration of aircraft with the control tower
func TestAircraftRegistration(t *testing.T) {
	// Create a control tower
	controlTower := NewAirTrafficControl("Test Tower")
	
	// Create aircraft
	passenger := NewPassengerAircraft("FL101", "Test Airways", 150)
	cargo := NewCargoAircraft("CG202", "Test Cargo", 5000.0)
	
	// Register aircraft
	controlTower.Register(passenger)
	controlTower.Register(cargo)
	
	// Check that they were registered
	registeredAircraft := controlTower.GetRegisteredAircraft()
	if len(registeredAircraft) != 2 {
		t.Errorf("Expected 2 registered aircraft, got %d", len(registeredAircraft))
	}
	
	// Check that we can retrieve an aircraft
	if _, found := controlTower.GetColleague("FL101"); !found {
		t.Errorf("Expected to find aircraft with ID FL101")
	}
	
	// Unregister an aircraft
	controlTower.Unregister(passenger)
	
	// Check that it was unregistered
	registeredAircraft = controlTower.GetRegisteredAircraft()
	if len(registeredAircraft) != 1 {
		t.Errorf("Expected 1 registered aircraft after unregistration, got %d", len(registeredAircraft))
	}
	
	// Try to unregister again (should have no effect)
	controlTower.Unregister(passenger)
	
	// Check count didn't change
	registeredAircraft = controlTower.GetRegisteredAircraft()
	if len(registeredAircraft) != 1 {
		t.Errorf("Expected 1 registered aircraft after repeated unregistration, got %d", len(registeredAircraft))
	}
}

// TestLandingRequest tests the landing request handling
func TestLandingRequest(t *testing.T) {
	// Setup
	controlTower := NewAirTrafficControl("Test Tower")
	passenger := NewPassengerAircraft("FL101", "Test Airways", 150)
	passenger.IsFlying = true
	passenger.Altitude = 5000
	
	// Register with control tower
	controlTower.Register(passenger)
	
	// Request landing
	passenger.RequestLanding()
	
	// Give some time for message processing
	time.Sleep(10 * time.Millisecond)
	
	// Verify the aircraft received a response
	messageLog := passenger.GetMessageLog()
	if len(messageLog) == 0 {
		t.Errorf("Expected at least one message in the aircraft's log")
	}
	
	// Check for landing clearance
	foundClearance := false
	for _, msg := range messageLog {
		if msg.Type == ControlMessage && strings.Contains(msg.Content, "Landing clearance granted") {
			foundClearance = true
			break
		}
	}
	
	if !foundClearance {
		t.Errorf("Expected to find landing clearance message")
	}
	
	// Check that the aircraft status was updated
	if passenger.GetStatus() != "Preparing for landing" {
		t.Errorf("Expected aircraft status to be 'Preparing for landing', got '%s'", passenger.GetStatus())
	}
}

// TestTakeoffRequest tests the takeoff request handling
func TestTakeoffRequest(t *testing.T) {
	// Setup
	controlTower := NewAirTrafficControl("Test Tower")
	cargo := NewCargoAircraft("CG202", "Test Cargo", 5000.0)
	
	// Register with control tower
	controlTower.Register(cargo)
	
	// Request takeoff
	cargo.RequestTakeoff()
	
	// Give some time for message processing
	time.Sleep(10 * time.Millisecond)
	
	// Verify the aircraft received a response
	messageLog := cargo.GetMessageLog()
	if len(messageLog) == 0 {
		t.Errorf("Expected at least one message in the aircraft's log")
	}
	
	// Check for takeoff clearance
	foundClearance := false
	for _, msg := range messageLog {
		if msg.Type == ControlMessage && strings.Contains(msg.Content, "Takeoff clearance granted") {
			foundClearance = true
			break
		}
	}
	
	if !foundClearance {
		t.Errorf("Expected to find takeoff clearance message")
	}
	
	// Check that the aircraft status was updated
	if cargo.GetStatus() != "Taking off" {
		t.Errorf("Expected aircraft status to be 'Taking off', got '%s'", cargo.GetStatus())
	}
	
	// Check that the aircraft is now flying
	if !cargo.IsFlying {
		t.Errorf("Expected aircraft to be flying after takeoff clearance")
	}
}

// TestEmergencyHandling tests emergency reporting and handling
func TestEmergencyHandling(t *testing.T) {
	// Setup
	controlTower := NewAirTrafficControl("Test Tower")
	aircraft1 := NewPassengerAircraft("FL101", "Test Airways", 150)
	aircraft2 := NewCargoAircraft("CG202", "Test Cargo", 5000.0)
	aircraft3 := NewMilitaryAircraft("MIL303", "Air Force", "Test Mission")
	
	// Register all aircraft
	controlTower.Register(aircraft1)
	controlTower.Register(aircraft2)
	controlTower.Register(aircraft3)
	
	// Report an emergency from aircraft1
	aircraft1.ReportEmergency("Engine failure")
	
	// Give some time for message processing
	time.Sleep(10 * time.Millisecond)
	
	// Check that the aircraft in emergency received an acknowledgment
	messageLog1 := aircraft1.GetMessageLog()
	foundAcknowledgment := false
	for _, msg := range messageLog1 {
		if msg.Type == ControlMessage && strings.Contains(msg.Content, "Emergency acknowledged") {
			foundAcknowledgment = true
			break
		}
	}
	
	if !foundAcknowledgment {
		t.Errorf("Expected emergency acknowledgment message for the reporting aircraft")
	}
	
	// Check that other aircraft received an alert
	messageLog2 := aircraft2.GetMessageLog()
	messageLog3 := aircraft3.GetMessageLog()
	
	foundAlert2 := false
	foundAlert3 := false
	
	for _, msg := range messageLog2 {
		if msg.Type == ControlMessage && strings.Contains(msg.Content, "EMERGENCY ALERT") {
			foundAlert2 = true
			break
		}
	}
	
	for _, msg := range messageLog3 {
		if msg.Type == ControlMessage && strings.Contains(msg.Content, "EMERGENCY ALERT") {
			foundAlert3 = true
			break
		}
	}
	
	if !foundAlert2 || !foundAlert3 {
		t.Errorf("Expected emergency alert messages for other aircraft")
	}
	
	// Check that the military aircraft's status was updated (special handling)
	if aircraft3.GetStatus() != "Alert: Supporting emergency situation" {
		t.Errorf("Expected military aircraft status to be updated for emergency, got '%s'", aircraft3.GetStatus())
	}
}

// TestBroadcastMessage tests the broadcast functionality
func TestBroadcastMessage(t *testing.T) {
	// Setup
	controlTower := NewAirTrafficControl("Test Tower")
	aircraft1 := NewPassengerAircraft("FL101", "Test Airways", 150)
	aircraft2 := NewCargoAircraft("CG202", "Test Cargo", 5000.0)
	
	// Register aircraft
	controlTower.Register(aircraft1)
	controlTower.Register(aircraft2)
	
	// Broadcast a message
	controlTower.Broadcast("Test Tower", ControlMessage, "Weather alert: Thunderstorm approaching", 8)
	
	// Give some time for message processing
	time.Sleep(10 * time.Millisecond)
	
	// Check that all aircraft received the broadcast
	messageLog1 := aircraft1.GetMessageLog()
	messageLog2 := aircraft2.GetMessageLog()
	
	foundBroadcast1 := false
	foundBroadcast2 := false
	
	for _, msg := range messageLog1 {
		if msg.Type == ControlMessage && strings.Contains(msg.Content, "Weather alert") {
			foundBroadcast1 = true
			break
		}
	}
	
	for _, msg := range messageLog2 {
		if msg.Type == ControlMessage && strings.Contains(msg.Content, "Weather alert") {
			foundBroadcast2 = true
			break
		}
	}
	
	if !foundBroadcast1 || !foundBroadcast2 {
		t.Errorf("Expected broadcast message to be received by all aircraft")
	}
}

// TestPositionUpdates tests position reporting
func TestPositionUpdates(t *testing.T) {
	// Setup
	controlTower := NewAirTrafficControl("Test Tower")
	aircraft := NewPrivateAircraft("PV404", "Test Owner")
	
	// Register aircraft
	controlTower.Register(aircraft)
	
	// Update position
	aircraft.UpdatePosition(100, 200, 3000)
	
	// Give some time for message processing
	time.Sleep(10 * time.Millisecond)
	
	// Check the tower's message log
	messageLog := controlTower.GetMessageLog()
	
	foundPositionUpdate := false
	for _, msg := range messageLog {
		if msg.Type == PositionUpdate && msg.From == "PV404" {
			foundPositionUpdate = true
			if !strings.Contains(msg.Content, "position (100, 200)") {
				t.Errorf("Position update content doesn't match expected coordinates")
			}
			if !strings.Contains(msg.Content, "altitude 3000") {
				t.Errorf("Position update content doesn't include correct altitude")
			}
			break
		}
	}
	
	if !foundPositionUpdate {
		t.Errorf("Expected position update in the control tower's log")
	}
}

// TestAircraftTypes tests different aircraft type behaviors
func TestAircraftTypes(t *testing.T) {
	// Test passenger aircraft
	passenger := NewPassengerAircraft("FL101", "Test Airways", 150)
	if passenger.Type != "Passenger" || passenger.Airline != "Test Airways" || passenger.PassengerCount != 150 {
		t.Errorf("Passenger aircraft initialization failed")
	}
	
	// Test cargo aircraft
	cargo := NewCargoAircraft("CG202", "Test Cargo", 5000.0)
	if cargo.Type != "Cargo" || cargo.Company != "Test Cargo" || cargo.CargoWeight != 5000.0 {
		t.Errorf("Cargo aircraft initialization failed")
	}
	
	// Test private aircraft
	private := NewPrivateAircraft("PV404", "Test Owner")
	if private.Type != "Private" || private.Owner != "Test Owner" {
		t.Errorf("Private aircraft initialization failed")
	}
	
	// Test military aircraft
	military := NewMilitaryAircraft("MIL303", "Air Force", "Test Mission")
	if military.Type != "Military" || military.Branch != "Air Force" || military.Mission != "Test Mission" {
		t.Errorf("Military aircraft initialization failed")
	}
}

// TestMessagePriority tests message prioritization
func TestMessagePriority(t *testing.T) {
	// Create message types with different priorities
	emergencyMsg := Message{
		Type:     EmergencyMessage,
		From:     "FL101",
		Content:  "Emergency!",
		Priority: 10,
	}
	
	landingMsg := Message{
		Type:     LandingRequest,
		From:     "CG202",
		Content:  "Request landing",
		Priority: 5,
	}
	
	positionMsg := Message{
		Type:     PositionUpdate,
		From:     "PV404",
		Content:  "Position update",
		Priority: 3,
	}
	
	// Check that priorities are set correctly
	if emergencyMsg.Priority <= landingMsg.Priority {
		t.Errorf("Emergency message should have higher priority than landing request")
	}
	
	if landingMsg.Priority <= positionMsg.Priority {
		t.Errorf("Landing request should have higher priority than position update")
	}
}

// TestLandTakeoffOperations tests the complete land and takeoff operations
func TestLandTakeoffOperations(t *testing.T) {
	// Setup
	controlTower := NewAirTrafficControl("Test Tower")
	aircraft := NewPassengerAircraft("FL101", "Test Airways", 150)
	
	// Register with control tower
	controlTower.Register(aircraft)
	
	// Initial state
	if aircraft.IsFlying {
		t.Errorf("New aircraft should not be flying initially")
	}
	
	// Take off
	aircraft.TakeOff(10000)
	
	// Verify state
	if !aircraft.IsFlying {
		t.Errorf("Aircraft should be flying after takeoff")
	}
	if aircraft.Altitude != 10000 {
		t.Errorf("Aircraft altitude should be 10000, got %d", aircraft.Altitude)
	}
	if aircraft.GetStatus() != "In air" {
		t.Errorf("Aircraft status should be 'In air', got '%s'", aircraft.GetStatus())
	}
	
	// Land
	aircraft.Land()
	
	// Verify state
	if aircraft.IsFlying {
		t.Errorf("Aircraft should not be flying after landing")
	}
	if aircraft.Altitude != 0 {
		t.Errorf("Aircraft altitude should be 0 after landing, got %d", aircraft.Altitude)
	}
	if aircraft.GetStatus() != "Landed" {
		t.Errorf("Aircraft status should be 'Landed', got '%s'", aircraft.GetStatus())
	}
}

// TestNilHandling tests handling of nil values
func TestNilHandling(t *testing.T) {
	// Create control tower
	controlTower := NewAirTrafficControl("Test Tower")
	
	// Try to register a nil colleague (should not panic)
	controlTower.Register(nil)
	
	// Try to unregister a nil colleague (should not panic)
	controlTower.Unregister(nil)
	
	// Try to get a non-existent colleague
	_, found := controlTower.GetColleague("NONEXISTENT")
	if found {
		t.Errorf("GetColleague should return false for non-existent ID")
	}
}

// TestMessageLogging tests that messages are properly logged
func TestMessageLogging(t *testing.T) {
	// Setup
	controlTower := NewAirTrafficControl("Test Tower")
	aircraft := NewPassengerAircraft("FL101", "Test Airways", 150)
	
	// Register with control tower
	controlTower.Register(aircraft)
	
	// Generate some messages
	aircraft.RequestTakeoff()
	aircraft.UpdatePosition(100, 200, 5000)
	aircraft.RequestLanding()
	
	// Give some time for message processing
	time.Sleep(10 * time.Millisecond)
	
	// Check control tower log
	towerLog := controlTower.GetMessageLog()
	if len(towerLog) < 3 {
		t.Errorf("Expected at least 3 messages in tower log, got %d", len(towerLog))
	}
	
	// Check aircraft log
	aircraftLog := aircraft.GetMessageLog()
	if len(aircraftLog) < 1 {
		t.Errorf("Expected at least 1 message in aircraft log, got %d", len(aircraftLog))
	}
}
