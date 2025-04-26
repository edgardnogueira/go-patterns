package main

import (
	"fmt"
	"github.com/edgardnogueira/go-patterns/behavioral/mediator"
	"time"
)

func main() {
	fmt.Println("Air Traffic Control - Mediator Pattern Demo")
	fmt.Println("===========================================")
	fmt.Println()

	// Create the air traffic control tower (Mediator)
	controlTower := mediator.NewAirTrafficControl("Central Tower")
	fmt.Printf("Created Air Traffic Control Tower: %s\n", controlTower.Name)
	fmt.Println()

	// Create various aircraft (Colleagues)
	aircraft1 := mediator.NewPassengerAircraft("UA123", "United Airlines", 240)
	aircraft2 := mediator.NewCargoAircraft("FX456", "FedEx", 12500.75)
	aircraft3 := mediator.NewPrivateAircraft("N789PJ", "John Smith")
	aircraft4 := mediator.NewMilitaryAircraft("AF001", "Air Force", "Training")

	// Register aircraft with the control tower
	fmt.Println("Registering aircraft with the control tower...")
	controlTower.Register(aircraft1)
	controlTower.Register(aircraft2)
	controlTower.Register(aircraft3)
	controlTower.Register(aircraft4)

	// Display registered aircraft
	registeredAircraft := controlTower.GetRegisteredAircraft()
	fmt.Printf("Registered aircraft (%d): %v\n", len(registeredAircraft), registeredAircraft)
	fmt.Println()

	fmt.Println("Starting simulation...")
	fmt.Println()

	// Scenario 1: Normal operations
	fmt.Println("SCENARIO 1: Normal Operations")
	fmt.Println("-----------------------------")

	// Position updates
	fmt.Println("Aircraft reporting positions...")
	aircraft1.UpdatePosition(10, 20, 10000)
	aircraft2.UpdatePosition(30, 40, 15000)
	aircraft3.UpdatePosition(50, 60, 8000)
	aircraft4.UpdatePosition(70, 80, 20000)
	time.Sleep(100 * time.Millisecond) // Allow time for processing
	fmt.Println()

	// Takeoff request
	fmt.Println("Private aircraft requesting takeoff...")
	aircraft3.RequestTakeoff()
	time.Sleep(100 * time.Millisecond) // Allow time for processing
	fmt.Printf("Private aircraft status: %s\n", aircraft3.GetStatus())
	fmt.Println()

	// Landing request
	fmt.Println("Passenger aircraft requesting landing...")
	aircraft1.IsFlying = true // Set as flying first
	aircraft1.RequestLanding()
	time.Sleep(100 * time.Millisecond) // Allow time for processing
	fmt.Printf("Passenger aircraft status: %s\n", aircraft1.GetStatus())
	fmt.Println()

	// Scenario 2: Emergency situation
	fmt.Println("SCENARIO 2: Emergency Situation")
	fmt.Println("------------------------------")

	fmt.Println("Cargo aircraft reporting emergency...")
	aircraft2.ReportEmergency("Engine failure")
	time.Sleep(100 * time.Millisecond) // Allow time for processing

	fmt.Println("Status of all aircraft after emergency:")
	fmt.Printf("Passenger: %s\n", aircraft1.GetStatus())
	fmt.Printf("Cargo: %s\n", aircraft2.GetStatus())
	fmt.Printf("Private: %s\n", aircraft3.GetStatus())
	fmt.Printf("Military: %s\n", aircraft4.GetStatus())
	fmt.Println()

	// Scenario 3: Weather broadcast
	fmt.Println("SCENARIO 3: Weather Broadcast")
	fmt.Println("----------------------------")

	fmt.Println("Control tower broadcasting weather alert...")
	controlTower.Broadcast(controlTower.Name, mediator.ControlMessage, 
		"Weather Alert: Heavy fog expected in 30 minutes, visibility will be reduced.", 7)
	time.Sleep(100 * time.Millisecond) // Allow time for processing
	fmt.Println()

	// Display message logs
	fmt.Println("LOGS")
	fmt.Println("----")
	
	// Show a sample of the most recent control tower log messages (last 3)
	towerLog := controlTower.GetMessageLog()
	fmt.Printf("Control Tower Log (last 3 of %d messages):\n", len(towerLog))
	if len(towerLog) > 0 {
		startIdx := len(towerLog) - 3
		if startIdx < 0 {
			startIdx = 0
		}
		for i := startIdx; i < len(towerLog); i++ {
			fmt.Printf("  %s\n", towerLog[i].String())
		}
	}
	fmt.Println()
	
	// Show a sample of aircraft message logs
	fmt.Printf("Passenger Aircraft Log (last 2 of %d messages):\n", len(aircraft1.GetMessageLog()))
	messages := aircraft1.GetMessageLog()
	if len(messages) > 0 {
		startIdx := len(messages) - 2
		if startIdx < 0 {
			startIdx = 0
		}
		for i := startIdx; i < len(messages); i++ {
			fmt.Printf("  %s\n", messages[i].String())
		}
	}
	fmt.Println()

	fmt.Println("Mediator Pattern Demo Complete")
}
