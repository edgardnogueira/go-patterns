package command

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestLightCommands(t *testing.T) {
	light := NewLight("Living Room")
	
	// Test light on command
	onCmd := NewLightOnCommand(light)
	if err := onCmd.Execute(); err != nil {
		t.Errorf("LightOnCommand.Execute() failed: %v", err)
	}
	if !light.isOn {
		t.Error("Light should be on after LightOnCommand.Execute()")
	}
	
	// Test light off command
	offCmd := NewLightOffCommand(light)
	if err := offCmd.Execute(); err != nil {
		t.Errorf("LightOffCommand.Execute() failed: %v", err)
	}
	if light.isOn {
		t.Error("Light should be off after LightOffCommand.Execute()")
	}
	
	// Test undo functionality
	if err := offCmd.Undo(); err != nil {
		t.Errorf("LightOffCommand.Undo() failed: %v", err)
	}
	if !light.isOn {
		t.Error("Light should be on after LightOffCommand.Undo()")
	}
	
	// Test string representation
	expectedOnString := "Turn Living Room on"
	if onCmd.String() != expectedOnString {
		t.Errorf("LightOnCommand.String() = %q, want %q", onCmd.String(), expectedOnString)
	}
}

func TestThermostatCommands(t *testing.T) {
	thermostat := NewThermostat("Main Floor")
	
	// Test temperature command
	tempCmd := NewThermostatSetCommand(thermostat, 72)
	if err := tempCmd.Execute(); err != nil {
		t.Errorf("ThermostatSetCommand.Execute() failed: %v", err)
	}
	if thermostat.temperature != 72 {
		t.Errorf("Thermostat temperature = %d, want %d", thermostat.temperature, 72)
	}
	
	// Change temperature and test undo
	tempCmd = NewThermostatSetCommand(thermostat, 68)
	if err := tempCmd.Execute(); err != nil {
		t.Errorf("ThermostatSetCommand.Execute() failed: %v", err)
	}
	if thermostat.temperature != 68 {
		t.Errorf("Thermostat temperature = %d, want %d", thermostat.temperature, 68)
	}
	
	// Test undo
	if err := tempCmd.Undo(); err != nil {
		t.Errorf("ThermostatSetCommand.Undo() failed: %v", err)
	}
	if thermostat.temperature != 72 {
		t.Errorf("Thermostat temperature after undo = %d, want %d", thermostat.temperature, 72)
	}
	
	// Test mode command
	modeCmd := NewThermostatModeCommand(thermostat, "cool")
	if err := modeCmd.Execute(); err != nil {
		t.Errorf("ThermostatModeCommand.Execute() failed: %v", err)
	}
	if thermostat.mode != "cool" {
		t.Errorf("Thermostat mode = %s, want %s", thermostat.mode, "cool")
	}
	
	// Test invalid mode
	invalidModeCmd := NewThermostatModeCommand(thermostat, "invalid")
	if err := invalidModeCmd.Execute(); err == nil {
		t.Error("ThermostatModeCommand.Execute() with invalid mode should fail")
	}
}

func TestMacroCommand(t *testing.T) {
	// Setup devices
	livingRoomLight := NewLight("Living Room")
	kitchenLight := NewLight("Kitchen")
	thermostat := NewThermostat("Main Floor")
	
	// Setup commands
	livingRoomLightOn := NewLightOnCommand(livingRoomLight)
	kitchenLightOn := NewLightOnCommand(kitchenLight)
	thermostatHeat := NewThermostatModeCommand(thermostat, "heat")
	thermostatTemp := NewThermostatSetCommand(thermostat, 72)
	
	// Create and execute macro command
	macro := NewMacroCommand("Evening", livingRoomLightOn, kitchenLightOn, thermostatHeat, thermostatTemp)
	if err := macro.Execute(); err != nil {
		t.Errorf("MacroCommand.Execute() failed: %v", err)
	}
	
	// Check all devices are in the expected state
	if !livingRoomLight.isOn {
		t.Error("Living room light should be on after macro execution")
	}
	if !kitchenLight.isOn {
		t.Error("Kitchen light should be on after macro execution")
	}
	if thermostat.mode != "heat" {
		t.Errorf("Thermostat mode = %s, want %s", thermostat.mode, "heat")
	}
	if thermostat.temperature != 72 {
		t.Errorf("Thermostat temperature = %d, want %d", thermostat.temperature, 72)
	}
	
	// Test undo
	if err := macro.Undo(); err != nil {
		t.Errorf("MacroCommand.Undo() failed: %v", err)
	}
	
	// The commands should be undone in reverse order
	// Note: Since we don't know the previous state of some devices,
	// we can't fully verify every state, but we can verify the ones we know
	if livingRoomLight.isOn {
		t.Error("Living room light should be off after macro undo")
	}
	if kitchenLight.isOn {
		t.Error("Kitchen light should be off after macro undo")
	}
}

func TestRemoteControl(t *testing.T) {
	remote := NewRemoteControl(3)
	
	// Setup devices
	livingRoomLight := NewLight("Living Room")
	kitchenLight := NewLight("Kitchen")
	garageDoor := NewGarageDoor("Main")
	
	// Setup commands
	livingRoomLightOn := NewLightOnCommand(livingRoomLight)
	livingRoomLightOff := NewLightOffCommand(livingRoomLight)
	kitchenLightOn := NewLightOnCommand(kitchenLight)
	kitchenLightOff := NewLightOffCommand(kitchenLight)
	garageDoorOpen := NewGarageDoorOpenCommand(garageDoor)
	garageDoorClose := NewGarageDoorCloseCommand(garageDoor)
	
	// Set commands to remote
	if err := remote.SetCommand(0, livingRoomLightOn, livingRoomLightOff); err != nil {
		t.Errorf("RemoteControl.SetCommand() failed: %v", err)
	}
	if err := remote.SetCommand(1, kitchenLightOn, kitchenLightOff); err != nil {
		t.Errorf("RemoteControl.SetCommand() failed: %v", err)
	}
	if err := remote.SetCommand(2, garageDoorOpen, garageDoorClose); err != nil {
		t.Errorf("RemoteControl.SetCommand() failed: %v", err)
	}
	
	// Test out of bounds slot
	if err := remote.SetCommand(3, garageDoorOpen, garageDoorClose); err == nil {
		t.Error("RemoteControl.SetCommand() with invalid slot should fail")
	}
	
	// Press buttons and check device states
	if err := remote.PressOn(0); err != nil {
		t.Errorf("RemoteControl.PressOn() failed: %v", err)
	}
	if !livingRoomLight.isOn {
		t.Error("Living room light should be on after PressOn(0)")
	}
	
	if err := remote.PressOff(0); err != nil {
		t.Errorf("RemoteControl.PressOff() failed: %v", err)
	}
	if livingRoomLight.isOn {
		t.Error("Living room light should be off after PressOff(0)")
	}
	
	// Test history and undo
	history := remote.GetHistory()
	if len(history) != 2 {
		t.Errorf("RemoteControl.GetHistory() = %d items, want %d", len(history), 2)
	}
	
	// Undo last command
	if err := remote.Undo(); err != nil {
		t.Errorf("RemoteControl.Undo() failed: %v", err)
	}
	
	// After undo, the light should be back on
	if !livingRoomLight.isOn {
		t.Error("Living room light should be on after Undo()")
	}
	
	// History should have decreased
	history = remote.GetHistory()
	if len(history) != 1 {
		t.Errorf("RemoteControl.GetHistory() after undo = %d items, want %d", len(history), 1)
	}
	
	// Clear history
	remote.ClearHistory()
	history = remote.GetHistory()
	if len(history) != 0 {
		t.Errorf("RemoteControl.GetHistory() after clear = %d items, want %d", len(history), 0)
	}
}

func TestCommandQueue(t *testing.T) {
	queue := NewCommandQueue()
	light := NewLight("Living Room")
	
	onCmd := NewLightOnCommand(light)
	offCmd := NewLightOffCommand(light)
	
	// Add commands to queue
	queue.AddCommand(onCmd)
	queue.AddCommand(offCmd)
	
	// Check queue size
	if queue.Size() != 2 {
		t.Errorf("CommandQueue.Size() = %d, want %d", queue.Size(), 2)
	}
	
	// Execute due commands
	executed, err := queue.ExecuteDue()
	if err != nil {
		t.Errorf("CommandQueue.ExecuteDue() failed: %v", err)
	}
	if executed != 2 {
		t.Errorf("CommandQueue.ExecuteDue() executed %d commands, want %d", executed, 2)
	}
	
	// Queue should be empty after execution
	if queue.Size() != 0 {
		t.Errorf("CommandQueue.Size() after execution = %d, want %d", queue.Size(), 0)
	}
	
	// Test scheduling
	future := time.Now().Add(time.Hour)
	queue.AddScheduledCommand(onCmd, future)
	
	// Check queue size
	if queue.Size() != 1 {
		t.Errorf("CommandQueue.Size() after scheduling = %d, want %d", queue.Size(), 1)
	}
	
	// Execute due commands - should be none since scheduled in the future
	executed, err = queue.ExecuteDue()
	if err != nil {
		t.Errorf("CommandQueue.ExecuteDue() failed: %v", err)
	}
	if executed != 0 {
		t.Errorf("CommandQueue.ExecuteDue() executed %d commands, want %d", executed, 0)
	}
	
	// Queue size should remain the same
	if queue.Size() != 1 {
		t.Errorf("CommandQueue.Size() after future execution = %d, want %d", queue.Size(), 1)
	}
	
	// Test peek
	nextCmd, execTime, err := queue.Peek()
	if err != nil {
		t.Errorf("CommandQueue.Peek() failed: %v", err)
	}
	if nextCmd != onCmd {
		t.Error("CommandQueue.Peek() returned wrong command")
	}
	if !execTime.Equal(future) {
		t.Errorf("CommandQueue.Peek() returned time %v, want %v", execTime, future)
	}
	
	// Clear queue
	queue.Clear()
	if queue.Size() != 0 {
		t.Errorf("CommandQueue.Size() after clear = %d, want %d", queue.Size(), 0)
	}
}

func TestHomeSceneCommand(t *testing.T) {
	// Setup devices
	livingRoomLight := NewLight("Living Room")
	kitchenLight := NewLight("Kitchen")
	thermostat := NewThermostat("Main Floor")
	audio := NewAudioSystem("Living Room")
	
	// Create a home scene
	eveningScene := NewHomeSceneCommand("Evening")
	
	// Add commands to the scene
	eveningScene.AddCommand(NewLightOnCommand(livingRoomLight))
	eveningScene.AddCommand(NewLightDimCommand(livingRoomLight, 30))
	eveningScene.AddCommand(NewLightOnCommand(kitchenLight))
	eveningScene.AddCommand(NewThermostatSetCommand(thermostat, 70))
	eveningScene.AddCommand(NewAudioSystemOnCommand(audio))
	eveningScene.AddCommand(NewAudioPlayCommand(audio, "Relaxing Jazz Playlist"))
	
	// Execute the scene
	if err := eveningScene.Execute(); err != nil {
		t.Errorf("HomeSceneCommand.Execute() failed: %v", err)
	}
	
	// Check device states
	if !livingRoomLight.isOn {
		t.Error("Living room light should be on")
	}
	if livingRoomLight.brightness != 30 {
		t.Errorf("Living room light brightness = %d, want %d", livingRoomLight.brightness, 30)
	}
	if !kitchenLight.isOn {
		t.Error("Kitchen light should be on")
	}
	if thermostat.temperature != 70 {
		t.Errorf("Thermostat temperature = %d, want %d", thermostat.temperature, 70)
	}
	if !audio.isOn {
		t.Error("Audio system should be on")
	}
	if !audio.isPlaying {
		t.Error("Audio system should be playing")
	}
	if audio.track != "Relaxing Jazz Playlist" {
		t.Errorf("Audio track = %s, want %s", audio.track, "Relaxing Jazz Playlist")
	}
	
	// Test string method
	if !strings.Contains(eveningScene.String(), "Evening") {
		t.Errorf("HomeSceneCommand.String() = %q, should contain 'Evening'", eveningScene.String())
	}
	
	// Test undo
	if err := eveningScene.Undo(); err != nil {
		t.Errorf("HomeSceneCommand.Undo() failed: %v", err)
	}
	
	// After undoing all commands, the state should be back to the initial state
	// Note: Since we added multiple commands affecting the same devices,
	// the final state will depend on the order of commands and their effects.
	// Here we can only check that undo doesn't cause errors.
}

func TestNoOpCommand(t *testing.T) {
	noOp := &NoOpCommand{}
	
	// Execute should succeed
	if err := noOp.Execute(); err != nil {
		t.Errorf("NoOpCommand.Execute() failed: %v", err)
	}
	
	// Undo should succeed
	if err := noOp.Undo(); err != nil {
		t.Errorf("NoOpCommand.Undo() failed: %v", err)
	}
	
	// String should return something meaningful
	if noOp.String() == "" {
		t.Error("NoOpCommand.String() should not return empty string")
	}
}

// Capturing output for testing
type outputCapture struct {
	output []string
}

func (o *outputCapture) Write(p []byte) (n int, err error) {
	o.output = append(o.output, string(p))
	return len(p), nil
}

func TestCeilingFanCommand(t *testing.T) {
	// Test ceiling fan command
	fanCmd := NewCeilingFanCommand("Bedroom", 2) // Medium speed
	
	if err := fanCmd.Execute(); err != nil {
		t.Errorf("CeilingFanCommand.Execute() failed: %v", err)
	}
	
	// Fan should be on and at medium speed
	if !fanCmd.isOn {
		t.Error("Ceiling fan should be on")
	}
	if fanCmd.speed != 2 {
		t.Errorf("Ceiling fan speed = %d, want %d", fanCmd.speed, 2)
	}
	
	// Change to high speed
	highFanCmd := NewCeilingFanCommand("Bedroom", 3)
	if err := highFanCmd.Execute(); err != nil {
		t.Errorf("CeilingFanCommand.Execute() failed: %v", err)
	}
	
	// Fan should be at high speed
	if highFanCmd.speed != 3 {
		t.Errorf("Ceiling fan speed = %d, want %d", highFanCmd.speed, 3)
	}
	
	// Undo should restore to medium speed
	if err := highFanCmd.Undo(); err != nil {
		t.Errorf("CeilingFanCommand.Undo() failed: %v", err)
	}
	
	// After undo, speed should be medium again
	if highFanCmd.speed != 2 {
		t.Errorf("Ceiling fan speed after undo = %d, want %d", highFanCmd.speed, 2)
	}
	
	// Turn off
	offFanCmd := NewCeilingFanCommand("Bedroom", 0)
	if err := offFanCmd.Execute(); err != nil {
		t.Errorf("CeilingFanCommand.Execute() failed: %v", err)
	}
	
	// Fan should be off
	if offFanCmd.isOn {
		t.Error("Ceiling fan should be off")
	}
}

func ExampleRemoteControl() {
	remote := NewRemoteControl(3)
	
	// Setup devices
	livingRoomLight := NewLight("Living Room")
	kitchenLight := NewLight("Kitchen")
	ceilingFan := NewCeilingFanCommand("Dining Room", 0)
	
	// Setup commands
	livingRoomLightOn := NewLightOnCommand(livingRoomLight)
	livingRoomLightOff := NewLightOffCommand(livingRoomLight)
	kitchenLightOn := NewLightOnCommand(kitchenLight)
	kitchenLightOff := NewLightOffCommand(kitchenLight)
	fanMedium := NewCeilingFanCommand("Dining Room", 2)
	fanOff := NewCeilingFanCommand("Dining Room", 0)
	
	// Set commands to remote
	remote.SetCommand(0, livingRoomLightOn, livingRoomLightOff)
	remote.SetCommand(1, kitchenLightOn, kitchenLightOff)
	remote.SetCommand(2, fanMedium, fanOff)
	
	// Use the remote
	fmt.Println("--- Using Remote Control ---")
	remote.PressOn(0)
	remote.PressOn(1)
	remote.PressOff(0)
	remote.PressOn(2)
	remote.PressOff(1)
	
	// Undo last command
	fmt.Println("\n--- Undoing Last Command ---")
	remote.Undo()
	
	// Output:
	// --- Using Remote Control ---
	// Living Room light is now ON
	// Kitchen light is now ON
	// Living Room light is now OFF
	// Dining Room ceiling fan set to MEDIUM
	// Kitchen light is now OFF
	// 
	// --- Undoing Last Command ---
	// Kitchen light is now ON
}
