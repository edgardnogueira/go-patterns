package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/edgardnogueira/go-patterns/behavioral/command"
)

func main() {
	fmt.Println("Smart Home Automation System")
	fmt.Println("============================")

	// Create devices
	livingRoomLight := command.NewLight("Living Room")
	kitchenLight := command.NewLight("Kitchen")
	bedroomLight := command.NewLight("Bedroom")
	thermostat := command.NewThermostat("Main Floor")
	audioSystem := command.NewAudioSystem("Living Room")
	garageDoor := command.NewGarageDoor("Main")

	// Create a remote control with 7 slots
	remote := command.NewRemoteControl(7)

	// Configure remote buttons
	remote.SetCommand(0, command.NewLightOnCommand(livingRoomLight), command.NewLightOffCommand(livingRoomLight))
	remote.SetCommand(1, command.NewLightOnCommand(kitchenLight), command.NewLightOffCommand(kitchenLight))
	remote.SetCommand(2, command.NewLightOnCommand(bedroomLight), command.NewLightOffCommand(bedroomLight))
	remote.SetCommand(3, command.NewThermostatSetCommand(thermostat, 72), command.NewThermostatSetCommand(thermostat, 68))
	remote.SetCommand(4, command.NewThermostatModeCommand(thermostat, "heat"), command.NewThermostatModeCommand(thermostat, "cool"))
	remote.SetCommand(5, command.NewAudioSystemOnCommand(audioSystem), command.NewAudioSystemOffCommand(audioSystem))
	remote.SetCommand(6, command.NewGarageDoorOpenCommand(garageDoor), command.NewGarageDoorCloseCommand(garageDoor))

	// Create scene macros
	morningScene := createMorningScene(livingRoomLight, kitchenLight, thermostat)
	eveningScene := createEveningScene(livingRoomLight, kitchenLight, bedroomLight, thermostat, audioSystem)
	awayScene := createAwayScene(livingRoomLight, kitchenLight, bedroomLight, thermostat, audioSystem, garageDoor)

	// Create a command queue for scheduling
	queue := command.NewCommandQueue()

	// Menu loop
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("\nSmart Home Control Menu:")
		fmt.Println("1. Show Device Status")
		fmt.Println("2. Use Remote Control")
		fmt.Println("3. Activate Scene")
		fmt.Println("4. Schedule Command")
		fmt.Println("5. Execute Scheduled Commands")
		fmt.Println("6. Undo Last Command")
		fmt.Println("7. Exit")
		fmt.Print("\nEnter choice: ")

		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			showStatus(livingRoomLight, kitchenLight, bedroomLight, thermostat, audioSystem, garageDoor)
		case "2":
			useRemoteControl(remote, scanner)
		case "3":
			activateScene(morningScene, eveningScene, awayScene, scanner)
		case "4":
			scheduleCommand(queue, livingRoomLight, kitchenLight, scanner)
		case "5":
			executeScheduled(queue)
		case "6":
			undoLastCommand(remote)
		case "7":
			fmt.Println("Exiting Smart Home Control System. Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func showStatus(
	livingRoomLight, kitchenLight, bedroomLight *command.Light,
	thermostat *command.Thermostat,
	audioSystem *command.AudioSystem,
	garageDoor *command.GarageDoor,
) {
	fmt.Println("\n=== Device Status ===")
	fmt.Println(livingRoomLight.GetStatus())
	fmt.Println(kitchenLight.GetStatus())
	fmt.Println(bedroomLight.GetStatus())
	fmt.Println(thermostat.GetStatus())
	fmt.Println(audioSystem.GetStatus())
	fmt.Println(garageDoor.GetStatus())
}

func useRemoteControl(remote *command.RemoteControl, scanner *bufio.Scanner) {
	fmt.Println("\n=== Remote Control ===")
	fmt.Println("Buttons:")
	fmt.Println("[0] Living Room Light")
	fmt.Println("[1] Kitchen Light")
	fmt.Println("[2] Bedroom Light")
	fmt.Println("[3] Thermostat Temperature (On=72°, Off=68°)")
	fmt.Println("[4] Thermostat Mode (On=heat, Off=cool)")
	fmt.Println("[5] Audio System")
	fmt.Println("[6] Garage Door")

	fmt.Print("\nEnter button number (0-6): ")
	scanner.Scan()
	buttonStr := scanner.Text()
	button, err := strconv.Atoi(buttonStr)
	if err != nil || button < 0 || button > 6 {
		fmt.Println("Invalid button number")
		return
	}

	fmt.Print("Press (on/off): ")
	scanner.Scan()
	action := strings.ToLower(scanner.Text())

	switch action {
	case "on":
		if err := remote.PressOn(button); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	case "off":
		if err := remote.PressOff(button); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	default:
		fmt.Println("Invalid action. Use 'on' or 'off'.")
	}
}

func createMorningScene(livingRoomLight, kitchenLight *command.Light, thermostat *command.Thermostat) *command.MacroCommand {
	morningScene := command.NewMacroCommand("Morning")
	morningScene.AddCommand(command.NewLightOnCommand(livingRoomLight))
	morningScene.AddCommand(command.NewLightDimCommand(livingRoomLight, 80))
	morningScene.AddCommand(command.NewLightOnCommand(kitchenLight))
	morningScene.AddCommand(command.NewLightDimCommand(kitchenLight, 100))
	morningScene.AddCommand(command.NewThermostatModeCommand(thermostat, "heat"))
	morningScene.AddCommand(command.NewThermostatSetCommand(thermostat, 70))
	return morningScene
}

func createEveningScene(livingRoomLight, kitchenLight, bedroomLight *command.Light, thermostat *command.Thermostat, audioSystem *command.AudioSystem) *command.MacroCommand {
	eveningScene := command.NewMacroCommand("Evening")
	eveningScene.AddCommand(command.NewLightOnCommand(livingRoomLight))
	eveningScene.AddCommand(command.NewLightDimCommand(livingRoomLight, 30))
	eveningScene.AddCommand(command.NewLightOnCommand(kitchenLight))
	eveningScene.AddCommand(command.NewLightDimCommand(kitchenLight, 50))
	eveningScene.AddCommand(command.NewLightOnCommand(bedroomLight))
	eveningScene.AddCommand(command.NewLightDimCommand(bedroomLight, 40))
	eveningScene.AddCommand(command.NewThermostatModeCommand(thermostat, "heat"))
	eveningScene.AddCommand(command.NewThermostatSetCommand(thermostat, 72))
	eveningScene.AddCommand(command.NewAudioSystemOnCommand(audioSystem))
	eveningScene.AddCommand(command.NewAudioPlayCommand(audioSystem, "Evening Jazz Playlist"))
	return eveningScene
}

func createAwayScene(
	livingRoomLight, kitchenLight, bedroomLight *command.Light,
	thermostat *command.Thermostat,
	audioSystem *command.AudioSystem,
	garageDoor *command.GarageDoor,
) *command.MacroCommand {
	awayScene := command.NewMacroCommand("Away")
	awayScene.AddCommand(command.NewLightOffCommand(livingRoomLight))
	awayScene.AddCommand(command.NewLightOffCommand(kitchenLight))
	awayScene.AddCommand(command.NewLightOffCommand(bedroomLight))
	awayScene.AddCommand(command.NewThermostatModeCommand(thermostat, "eco"))
	awayScene.AddCommand(command.NewThermostatSetCommand(thermostat, 65))
	awayScene.AddCommand(command.NewAudioSystemOffCommand(audioSystem))
	awayScene.AddCommand(command.NewGarageDoorCloseCommand(garageDoor))
	return awayScene
}

func activateScene(morningScene, eveningScene, awayScene *command.MacroCommand, scanner *bufio.Scanner) {
	fmt.Println("\n=== Activate Scene ===")
	fmt.Println("Available scenes:")
	fmt.Println("1. Morning Scene")
	fmt.Println("2. Evening Scene")
	fmt.Println("3. Away Scene")

	fmt.Print("\nSelect scene (1-3): ")
	scanner.Scan()
	sceneChoice := scanner.Text()

	var scene *command.MacroCommand
	switch sceneChoice {
	case "1":
		scene = morningScene
	case "2":
		scene = eveningScene
	case "3":
		scene = awayScene
	default:
		fmt.Println("Invalid scene selection")
		return
	}

	fmt.Printf("Activating %s scene...\n", scene.String())
	if err := scene.Execute(); err != nil {
		fmt.Printf("Error activating scene: %v\n", err)
	} else {
		fmt.Println("Scene activated successfully!")
	}
}

func scheduleCommand(queue *command.CommandQueue, livingRoomLight, kitchenLight *command.Light, scanner *bufio.Scanner) {
	fmt.Println("\n=== Schedule Command ===")
	fmt.Println("Available commands:")
	fmt.Println("1. Turn Living Room Light On")
	fmt.Println("2. Turn Living Room Light Off")
	fmt.Println("3. Turn Kitchen Light On")
	fmt.Println("4. Turn Kitchen Light Off")

	fmt.Print("\nSelect command (1-4): ")
	scanner.Scan()
	cmdChoice := scanner.Text()

	var cmd command.Command
	switch cmdChoice {
	case "1":
		cmd = command.NewLightOnCommand(livingRoomLight)
	case "2":
		cmd = command.NewLightOffCommand(livingRoomLight)
	case "3":
		cmd = command.NewLightOnCommand(kitchenLight)
	case "4":
		cmd = command.NewLightOffCommand(kitchenLight)
	default:
		fmt.Println("Invalid command selection")
		return
	}

	fmt.Print("Enter delay in seconds: ")
	scanner.Scan()
	delayStr := scanner.Text()
	delay, err := strconv.Atoi(delayStr)
	if err != nil || delay < 0 {
		fmt.Println("Invalid delay time")
		return
	}

	executionTime := time.Now().Add(time.Duration(delay) * time.Second)
	queue.AddScheduledCommand(cmd, executionTime)
	fmt.Printf("Command scheduled for execution in %d seconds\n", delay)
}

func executeScheduled(queue *command.CommandQueue) {
	fmt.Println("\n=== Executing Scheduled Commands ===")
	executed, err := queue.ExecuteDue()
	if err != nil {
		fmt.Printf("Error executing commands: %v\n", err)
		return
	}
	fmt.Printf("Executed %d command(s)\n", executed)
	fmt.Printf("%d command(s) remaining in queue\n", queue.Size())
}

func undoLastCommand(remote *command.RemoteControl) {
	fmt.Println("\n=== Undoing Last Command ===")
	if err := remote.Undo(); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Last command undone successfully!")
	}
}
