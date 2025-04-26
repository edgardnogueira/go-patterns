package main

import (
	"bufio"
	"fmt"
	"github.com/edgardnogueira/go-patterns/behavioral/mediator"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// AirportSimulator extends the mediator pattern to simulate a busy airport
type AirportSimulator struct {
	// The control tower (mediator)
	controlTower *mediator.AirTrafficControl
	// Map of aircraft by ID for easy lookup
	aircraft map[string]mediator.Colleague
	// Running flag to control simulation
	running bool
	// Current simulation time (minutes)
	time int
	// Weather conditions
	weather string
	// Random source for simulation events
	random *rand.Rand
}

// NewAirportSimulator creates a new airport simulator
func NewAirportSimulator(airportName string) *AirportSimulator {
	return &AirportSimulator{
		controlTower: mediator.NewAirTrafficControl(airportName + " Tower"),
		aircraft:     make(map[string]mediator.Colleague),
		running:      false,
		time:         0,
		weather:      "Clear",
		random:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// AddAircraft adds a new aircraft to the simulation
func (sim *AirportSimulator) AddAircraft(aircraft mediator.Colleague) {
	id := aircraft.GetID()
	sim.aircraft[id] = aircraft
	sim.controlTower.Register(aircraft)
	fmt.Printf("Added aircraft %s to simulation\n", id)
}

// Start begins the simulation
func (sim *AirportSimulator) Start() {
	if sim.running {
		return
	}
	
	sim.running = true
	fmt.Printf("\nStarting airport simulation at %s...\n", sim.controlTower.Name)
	
	// Main simulation loop
	go func() {
		for sim.running {
			sim.time++
			
			// Run simulation steps
			sim.updateWeather()
			sim.generateRandomEvents()
			
			// Print simulation time every 5 minutes
			if sim.time%5 == 0 {
				fmt.Printf("\n-- Simulation time: %d minutes | Weather: %s --\n", 
					sim.time, sim.weather)
				
				// Print aircraft status periodically
				if len(sim.aircraft) > 0 && sim.time%10 == 0 {
					fmt.Println("Current aircraft status:")
					for id, ac := range sim.aircraft {
						fmt.Printf("  %s: %s\n", id, ac.GetStatus())
					}
					fmt.Println()
				}
			}
			
			// Sleep to control simulation speed
			time.Sleep(1 * time.Second) // 1 second = 1 simulation minute
		}
	}()
}

// Stop ends the simulation
func (sim *AirportSimulator) Stop() {
	sim.running = false
	fmt.Println("\nStopping airport simulation...")
	time.Sleep(1 * time.Second)
	fmt.Printf("Simulation ended after %d minutes\n", sim.time)
}

// updateWeather randomly changes weather conditions
func (sim *AirportSimulator) updateWeather() {
	// Only change weather occasionally (1% chance per minute)
	if sim.random.Intn(100) == 0 {
		weatherTypes := []string{
			"Clear", "Cloudy", "Rain", "Heavy Rain", "Fog", "Snow", "Windy", "Stormy",
		}
		
		newWeather := weatherTypes[sim.random.Intn(len(weatherTypes))]
		if newWeather != sim.weather {
			sim.weather = newWeather
			// Broadcast weather change to all aircraft
			sim.controlTower.Broadcast(
				sim.controlTower.Name,
				mediator.ControlMessage,
				fmt.Sprintf("Weather update: Conditions changing to %s", sim.weather),
				6,
			)
		}
	}
}

// generateRandomEvents creates random aviation events
func (sim *AirportSimulator) generateRandomEvents() {
	// Only generate events occasionally
	if sim.random.Intn(5) != 0 {
		return
	}
	
	// Select a random event type
	eventType := sim.random.Intn(10)
	
	switch eventType {
	case 0:
		// Incoming aircraft
		sim.createRandomAircraft(true)
	case 1:
		// Departing aircraft
		sim.createRandomAircraft(false)
	case 2:
		// Position updates
		sim.updateRandomAircraftPosition()
	case 3:
		// Landing request
		sim.requestRandomLanding()
	case 4:
		// Takeoff request
		sim.requestRandomTakeoff()
	case 5:
		// Emergency (rare)
		if sim.random.Intn(10) == 0 {
			sim.createRandomEmergency()
		}
	}
}

// createRandomAircraft generates a new random aircraft
func (sim *AirportSimulator) createRandomAircraft(arriving bool) {
	// Aircraft types
	aircraftTypes := []string{"Passenger", "Cargo", "Private", "Military"}
	aircraftType := aircraftTypes[sim.random.Intn(len(aircraftTypes))]
	
	// Generate ID
	var id string
	var aircraft mediator.Colleague
	
	switch aircraftType {
	case "Passenger":
		airlines := []string{"United", "Delta", "American", "Southwest", "JetBlue"}
		airline := airlines[sim.random.Intn(len(airlines))]
		id = fmt.Sprintf("%s%d", airline[:2], 100+sim.random.Intn(900))
		passengers := 50 + sim.random.Intn(300)
		aircraft = mediator.NewPassengerAircraft(id, airline+" Airlines", passengers)
	
	case "Cargo":
		companies := []string{"FedEx", "UPS", "DHL", "Amazon", "Atlas"}
		company := companies[sim.random.Intn(len(companies))]
		id = fmt.Sprintf("%s%d", company[:2], 100+sim.random.Intn(900))
		weight := 1000.0 + sim.random.Float64()*20000.0
		aircraft = mediator.NewCargoAircraft(id, company+" Cargo", weight)
	
	case "Private":
		owners := []string{"Smith", "Johnson", "Williams", "Brown", "Jones"}
		owner := owners[sim.random.Intn(len(owners))]
		id = fmt.Sprintf("N%d%s", 100+sim.random.Intn(900), owner[:1])
		aircraft = mediator.NewPrivateAircraft(id, owner)
	
	case "Military":
		branches := []string{"Air Force", "Navy", "Army", "Coast Guard"}
		branch := branches[sim.random.Intn(len(branches))]
		missions := []string{"Training", "Transport", "Patrol", "Exercise"}
		mission := missions[sim.random.Intn(len(missions))]
		id = fmt.Sprintf("MIL%d", 100+sim.random.Intn(900))
		aircraft = mediator.NewMilitaryAircraft(id, branch, mission)
	}
	
	// Set initial state
	if ac, ok := aircraft.(*mediator.Aircraft); ok {
		// Random position
		x := sim.random.Intn(100)
		y := sim.random.Intn(100)
		
		if arriving {
			ac.IsFlying = true
			altitude := 5000 + sim.random.Intn(25000)
			ac.UpdatePosition(x, y, altitude)
			ac.SetStatus("Approaching airport")
			fmt.Printf("New aircraft %s (%s) approaching\n", id, aircraftType)
		} else {
			ac.IsFlying = false
			ac.UpdatePosition(x, y, 0)
			ac.SetStatus("Ready for departure")
			fmt.Printf("New aircraft %s (%s) ready for departure\n", id, aircraftType)
		}
	}
	
	sim.AddAircraft(aircraft)
}

// updateRandomAircraftPosition updates position of a random aircraft
func (sim *AirportSimulator) updateRandomAircraftPosition() {
	if len(sim.aircraft) == 0 {
		return
	}
	
	// Select random aircraft
	ids := make([]string, 0, len(sim.aircraft))
	for id := range sim.aircraft {
		ids = append(ids, id)
	}
	
	id := ids[sim.random.Intn(len(ids))]
	ac := sim.aircraft[id]
	
	// Update position if it's an Aircraft type
	if aircraft, ok := getAircraftFromColleague(ac); ok {
		x := sim.random.Intn(100)
		y := sim.random.Intn(100)
		
		var altitude int
		if aircraft.IsFlying {
			altitude = 1000 + sim.random.Intn(30000)
		} else {
			altitude = 0
		}
		
		aircraft.UpdatePosition(x, y, altitude)
	}
}

// requestRandomLanding makes a random aircraft request landing
func (sim *AirportSimulator) requestRandomLanding() {
	// Find flying aircraft
	var flyingAircraft []*mediator.Aircraft
	
	for _, ac := range sim.aircraft {
		if aircraft, ok := getAircraftFromColleague(ac); ok {
			if aircraft.IsFlying {
				flyingAircraft = append(flyingAircraft, aircraft)
			}
		}
	}
	
	// Request landing if there are flying aircraft
	if len(flyingAircraft) > 0 {
		aircraft := flyingAircraft[sim.random.Intn(len(flyingAircraft))]
		fmt.Printf("Aircraft %s requesting landing\n", aircraft.ID)
		aircraft.RequestLanding()
	}
}

// requestRandomTakeoff makes a random aircraft request takeoff
func (sim *AirportSimulator) requestRandomTakeoff() {
	// Find grounded aircraft
	var groundedAircraft []*mediator.Aircraft
	
	for _, ac := range sim.aircraft {
		if aircraft, ok := getAircraftFromColleague(ac); ok {
			if !aircraft.IsFlying {
				groundedAircraft = append(groundedAircraft, aircraft)
			}
		}
	}
	
	// Request takeoff if there are grounded aircraft
	if len(groundedAircraft) > 0 {
		aircraft := groundedAircraft[sim.random.Intn(len(groundedAircraft))]
		fmt.Printf("Aircraft %s requesting takeoff\n", aircraft.ID)
		aircraft.RequestTakeoff()
	}
}

// createRandomEmergency generates a random emergency situation
func (sim *AirportSimulator) createRandomEmergency() {
	if len(sim.aircraft) == 0 {
		return
	}
	
	// Select random aircraft
	ids := make([]string, 0, len(sim.aircraft))
	for id := range sim.aircraft {
		ids = append(ids, id)
	}
	
	id := ids[sim.random.Intn(len(ids))]
	ac := sim.aircraft[id]
	
	// Create emergency if it's an Aircraft type
	if aircraft, ok := getAircraftFromColleague(ac); ok {
		emergencies := []string{
			"Engine failure", "Hydraulic system failure", "Low fuel", 
			"Medical emergency", "Pressurization problem", "Bird strike",
			"Navigation system failure", "Weather-related issue",
		}
		
		emergency := emergencies[sim.random.Intn(len(emergencies))]
		fmt.Printf("EMERGENCY: Aircraft %s reporting %s\n", aircraft.ID, emergency)
		aircraft.ReportEmergency(emergency)
	}
}

// Helper function to get the aircraft from a colleague interface
func getAircraftFromColleague(colleague mediator.Colleague) (*mediator.Aircraft, bool) {
	// Try to cast to each aircraft type
	if aircraft, ok := colleague.(*mediator.PassengerAircraft); ok {
		return &aircraft.Aircraft, true
	} else if aircraft, ok := colleague.(*mediator.CargoAircraft); ok {
		return &aircraft.Aircraft, true
	} else if aircraft, ok := colleague.(*mediator.PrivateAircraft); ok {
		return &aircraft.Aircraft, true
	} else if aircraft, ok := colleague.(*mediator.MilitaryAircraft); ok {
		return &aircraft.Aircraft, true
	}
	
	return nil, false
}

// Interactive menu for the simulator
func (sim *AirportSimulator) InteractiveMenu() {
	reader := bufio.NewReader(os.Stdin)
	
	for {
		fmt.Println("\n=========================================")
		fmt.Println("AIRPORT SIMULATOR - INTERACTIVE MENU")
		fmt.Println("=========================================")
		fmt.Println("1. Start simulation")
		fmt.Println("2. Stop simulation")
		fmt.Println("3. Add aircraft")
		fmt.Println("4. List all aircraft")
		fmt.Println("5. Show control tower logs")
		fmt.Println("6. Create emergency")
		fmt.Println("7. Change weather")
		fmt.Println("8. Exit")
		fmt.Print("\nEnter your choice: ")
		
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		
		switch text {
		case "1":
			sim.Start()
		case "2":
			sim.Stop()
		case "3":
			sim.menuAddAircraft(reader)
		case "4":
			sim.listAllAircraft()
		case "5":
			sim.showControlTowerLogs()
		case "6":
			sim.menuCreateEmergency(reader)
		case "7":
			sim.menuChangeWeather(reader)
		case "8":
			fmt.Println("Exiting simulator...")
			sim.Stop()
			return
		default:
			fmt.Println("Invalid option, try again.")
		}
	}
}

// Helper function for adding aircraft through the menu
func (sim *AirportSimulator) menuAddAircraft(reader *bufio.Reader) {
	fmt.Println("\nADD AIRCRAFT")
	fmt.Println("===========")
	fmt.Println("Types: 1=Passenger, 2=Cargo, 3=Private, 4=Military")
	
	fmt.Print("Enter aircraft type (1-4): ")
	typeStr, _ := reader.ReadString('\n')
	typeStr = strings.TrimSpace(typeStr)
	typeNum, _ := strconv.Atoi(typeStr)
	
	fmt.Print("Enter aircraft ID: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)
	
	fmt.Print("Is aircraft already flying? (y/n): ")
	flying, _ := reader.ReadString('\n')
	flying = strings.TrimSpace(flying)
	isFlying := strings.ToLower(flying) == "y"
	
	var aircraft mediator.Colleague
	
	switch typeNum {
	case 1:
		fmt.Print("Enter airline name: ")
		airline, _ := reader.ReadString('\n')
		airline = strings.TrimSpace(airline)
		
		fmt.Print("Enter passenger count: ")
		countStr, _ := reader.ReadString('\n')
		countStr = strings.TrimSpace(countStr)
		count, _ := strconv.Atoi(countStr)
		
		aircraft = mediator.NewPassengerAircraft(id, airline, count)
	case 2:
		fmt.Print("Enter company name: ")
		company, _ := reader.ReadString('\n')
		company = strings.TrimSpace(company)
		
		fmt.Print("Enter cargo weight (kg): ")
		weightStr, _ := reader.ReadString('\n')
		weightStr = strings.TrimSpace(weightStr)
		weight, _ := strconv.ParseFloat(weightStr, 64)
		
		aircraft = mediator.NewCargoAircraft(id, company, weight)
	case 3:
		fmt.Print("Enter owner name: ")
		owner, _ := reader.ReadString('\n')
		owner = strings.TrimSpace(owner)
		
		aircraft = mediator.NewPrivateAircraft(id, owner)
	case 4:
		fmt.Print("Enter military branch: ")
		branch, _ := reader.ReadString('\n')
		branch = strings.TrimSpace(branch)
		
		fmt.Print("Enter mission: ")
		mission, _ := reader.ReadString('\n')
		mission = strings.TrimSpace(mission)
		
		aircraft = mediator.NewMilitaryAircraft(id, branch, mission)
	default:
		fmt.Println("Invalid aircraft type.")
		return
	}
	
	// Set flying status if it's an Aircraft type
	if ac, ok := getAircraftFromColleague(aircraft); ok {
		ac.IsFlying = isFlying
		if isFlying {
			ac.Altitude = 10000
			ac.SetStatus("In flight")
		} else {
			ac.Altitude = 0
			ac.SetStatus("On ground")
		}
	}
	
	sim.AddAircraft(aircraft)
	fmt.Printf("Aircraft %s added successfully\n", id)
}

// Helper function to list all aircraft
func (sim *AirportSimulator) listAllAircraft() {
	fmt.Println("\nALL AIRCRAFT")
	fmt.Println("============")
	
	if len(sim.aircraft) == 0 {
		fmt.Println("No aircraft in simulation.")
		return
	}
	
	for id, ac := range sim.aircraft {
		aircraft, ok := getAircraftFromColleague(ac)
		if ok {
			flying := "On ground"
			if aircraft.IsFlying {
				flying = fmt.Sprintf("Flying at %d feet", aircraft.Altitude)
			}
			
			fmt.Printf("%s: %s, Status: %s, Position: %s, %s\n",
				id, aircraft.Type, aircraft.GetStatus(), aircraft.Position.String(), flying)
		} else {
			fmt.Printf("%s: Unknown aircraft type\n", id)
		}
	}
}

// Helper function to show control tower logs
func (sim *AirportSimulator) showControlTowerLogs() {
	fmt.Println("\nCONTROL TOWER LOGS")
	fmt.Println("==================")
	
	logs := sim.controlTower.GetMessageLog()
	if len(logs) == 0 {
		fmt.Println("No logs available.")
		return
	}
	
	// Show last 10 logs or all if fewer than 10
	start := 0
	if len(logs) > 10 {
		start = len(logs) - 10
	}
	
	for i := start; i < len(logs); i++ {
		fmt.Printf("[%s] %s\n", 
			logs[i].Timestamp.Format("15:04:05"), logs[i].String())
	}
}

// Helper function to create an emergency through the menu
func (sim *AirportSimulator) menuCreateEmergency(reader *bufio.Reader) {
	fmt.Println("\nCREATE EMERGENCY")
	fmt.Println("===============")
	
	if len(sim.aircraft) == 0 {
		fmt.Println("No aircraft in simulation.")
		return
	}
	
	fmt.Println("Available aircraft:")
	for id := range sim.aircraft {
		fmt.Printf("- %s\n", id)
	}
	
	fmt.Print("Enter aircraft ID: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)
	
	ac, exists := sim.aircraft[id]
	if !exists {
		fmt.Println("Aircraft not found.")
		return
	}
	
	fmt.Print("Enter emergency details: ")
	details, _ := reader.ReadString('\n')
	details = strings.TrimSpace(details)
	
	if aircraft, ok := getAircraftFromColleague(ac); ok {
		aircraft.ReportEmergency(details)
		fmt.Printf("Emergency created for aircraft %s\n", id)
	} else {
		fmt.Println("Failed to create emergency.")
	}
}

// Helper function to change weather through the menu
func (sim *AirportSimulator) menuChangeWeather(reader *bufio.Reader) {
	fmt.Println("\nCHANGE WEATHER")
	fmt.Println("==============")
	fmt.Println("Current weather: " + sim.weather)
	fmt.Println("Options: Clear, Cloudy, Rain, Heavy Rain, Fog, Snow, Windy, Stormy")
	
	fmt.Print("Enter new weather: ")
	weather, _ := reader.ReadString('\n')
	weather = strings.TrimSpace(weather)
	
	if weather != "" {
		oldWeather := sim.weather
		sim.weather = weather
		
		// Broadcast weather change
		sim.controlTower.Broadcast(
			sim.controlTower.Name,
			mediator.ControlMessage,
			fmt.Sprintf("Weather update: Conditions changing from %s to %s", oldWeather, sim.weather),
			6,
		)
		
		fmt.Printf("Weather changed to: %s\n", sim.weather)
	}
}

func main() {
	// Create the simulator
	simulator := NewAirportSimulator("International")
	
	// Add some initial aircraft
	simulator.AddAircraft(mediator.NewPassengerAircraft("AA123", "American Airlines", 220))
	simulator.AddAircraft(mediator.NewCargoAircraft("FX456", "FedEx", 15000.0))
	simulator.AddAircraft(mediator.NewPrivateAircraft("N789J", "John Smith"))
	
	// Start interactive menu
	simulator.InteractiveMenu()
}
