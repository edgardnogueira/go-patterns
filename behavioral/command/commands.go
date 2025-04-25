package command

import "fmt"

// LightOnCommand turns a light on
type LightOnCommand struct {
	light *Light
}

// NewLightOnCommand creates a new LightOnCommand
func NewLightOnCommand(light *Light) *LightOnCommand {
	return &LightOnCommand{light: light}
}

// Execute turns the light on
func (c *LightOnCommand) Execute() error {
	c.light.TurnOn()
	return nil
}

// Undo turns the light off
func (c *LightOnCommand) Undo() error {
	c.light.TurnOff()
	return nil
}

// String returns a description of the command
func (c *LightOnCommand) String() string {
	return fmt.Sprintf("Turn %s on", c.light.name)
}

// LightOffCommand turns a light off
type LightOffCommand struct {
	light *Light
}

// NewLightOffCommand creates a new LightOffCommand
func NewLightOffCommand(light *Light) *LightOffCommand {
	return &LightOffCommand{light: light}
}

// Execute turns the light off
func (c *LightOffCommand) Execute() error {
	c.light.TurnOff()
	return nil
}

// Undo turns the light on
func (c *LightOffCommand) Undo() error {
	c.light.TurnOn()
	return nil
}

// String returns a description of the command
func (c *LightOffCommand) String() string {
	return fmt.Sprintf("Turn %s off", c.light.name)
}

// LightDimCommand sets the brightness of a light
type LightDimCommand struct {
	light           *Light
	level           int
	previousLevel   int
}

// NewLightDimCommand creates a new LightDimCommand
func NewLightDimCommand(light *Light, level int) *LightDimCommand {
	return &LightDimCommand{
		light: light,
		level: level,
	}
}

// Execute sets the light brightness to the specified level
func (c *LightDimCommand) Execute() error {
	c.previousLevel = c.light.brightness
	c.light.SetBrightness(c.level)
	return nil
}

// Undo restores the light brightness to its previous level
func (c *LightDimCommand) Undo() error {
	c.light.SetBrightness(c.previousLevel)
	return nil
}

// String returns a description of the command
func (c *LightDimCommand) String() string {
	return fmt.Sprintf("Set %s brightness to %d%%", c.light.name, c.level)
}

// ThermostatSetCommand sets the temperature of a thermostat
type ThermostatSetCommand struct {
	thermostat      *Thermostat
	temperature     int
	previousTemp    int
}

// NewThermostatSetCommand creates a new ThermostatSetCommand
func NewThermostatSetCommand(thermostat *Thermostat, temperature int) *ThermostatSetCommand {
	return &ThermostatSetCommand{
		thermostat:  thermostat,
		temperature: temperature,
	}
}

// Execute sets the thermostat temperature
func (c *ThermostatSetCommand) Execute() error {
	c.previousTemp = c.thermostat.temperature
	c.thermostat.SetTemperature(c.temperature)
	return nil
}

// Undo restores the previous temperature
func (c *ThermostatSetCommand) Undo() error {
	c.thermostat.SetTemperature(c.previousTemp)
	return nil
}

// String returns a description of the command
func (c *ThermostatSetCommand) String() string {
	return fmt.Sprintf("Set %s temperature to %d°", c.thermostat.name, c.temperature)
}

// ThermostatModeCommand sets the mode of a thermostat
type ThermostatModeCommand struct {
	thermostat     *Thermostat
	mode           string
	previousMode   string
}

// NewThermostatModeCommand creates a new ThermostatModeCommand
func NewThermostatModeCommand(thermostat *Thermostat, mode string) *ThermostatModeCommand {
	return &ThermostatModeCommand{
		thermostat: thermostat,
		mode:       mode,
	}
}

// Execute sets the thermostat mode
func (c *ThermostatModeCommand) Execute() error {
	c.previousMode = c.thermostat.mode
	return c.thermostat.SetMode(c.mode)
}

// Undo restores the previous mode
func (c *ThermostatModeCommand) Undo() error {
	return c.thermostat.SetMode(c.previousMode)
}

// String returns a description of the command
func (c *ThermostatModeCommand) String() string {
	return fmt.Sprintf("Set %s mode to %s", c.thermostat.name, c.mode)
}

// AudioSystemOnCommand turns an audio system on
type AudioSystemOnCommand struct {
	audio *AudioSystem
}

// NewAudioSystemOnCommand creates a new AudioSystemOnCommand
func NewAudioSystemOnCommand(audio *AudioSystem) *AudioSystemOnCommand {
	return &AudioSystemOnCommand{audio: audio}
}

// Execute turns the audio system on
func (c *AudioSystemOnCommand) Execute() error {
	c.audio.TurnOn()
	return nil
}

// Undo turns the audio system off
func (c *AudioSystemOnCommand) Undo() error {
	c.audio.TurnOff()
	return nil
}

// String returns a description of the command
func (c *AudioSystemOnCommand) String() string {
	return fmt.Sprintf("Turn %s on", c.audio.name)
}

// AudioSystemOffCommand turns an audio system off
type AudioSystemOffCommand struct {
	audio       *AudioSystem
	wasPlaying  bool
	trackPlayed string
}

// NewAudioSystemOffCommand creates a new AudioSystemOffCommand
func NewAudioSystemOffCommand(audio *AudioSystem) *AudioSystemOffCommand {
	return &AudioSystemOffCommand{audio: audio}
}

// Execute turns the audio system off
func (c *AudioSystemOffCommand) Execute() error {
	c.wasPlaying = c.audio.isPlaying
	c.trackPlayed = c.audio.track
	c.audio.TurnOff()
	return nil
}

// Undo turns the audio system on and resumes playback if it was playing
func (c *AudioSystemOffCommand) Undo() error {
	c.audio.TurnOn()
	if c.wasPlaying && c.trackPlayed != "" {
		return c.audio.Play(c.trackPlayed)
	}
	return nil
}

// String returns a description of the command
func (c *AudioSystemOffCommand) String() string {
	return fmt.Sprintf("Turn %s off", c.audio.name)
}

// AudioPlayCommand plays a track on an audio system
type AudioPlayCommand struct {
	audio          *AudioSystem
	track          string
	wasPlaying     bool
	previousTrack  string
}

// NewAudioPlayCommand creates a new AudioPlayCommand
func NewAudioPlayCommand(audio *AudioSystem, track string) *AudioPlayCommand {
	return &AudioPlayCommand{
		audio: audio,
		track: track,
	}
}

// Execute plays the track
func (c *AudioPlayCommand) Execute() error {
	// Record previous state for undo
	c.wasPlaying = c.audio.isPlaying
	c.previousTrack = c.audio.track
	
	// If system is off, turn it on
	if !c.audio.isOn {
		c.audio.TurnOn()
	}
	
	return c.audio.Play(c.track)
}

// Undo stops playback or restores previous track
func (c *AudioPlayCommand) Undo() error {
	if c.wasPlaying {
		// Resume playing the previous track
		return c.audio.Play(c.previousTrack)
	} else {
		// Stop playing
		c.audio.Stop()
	}
	return nil
}

// String returns a description of the command
func (c *AudioPlayCommand) String() string {
	return fmt.Sprintf("Play '%s' on %s", c.track, c.audio.name)
}

// AudioStopCommand stops playback on an audio system
type AudioStopCommand struct {
	audio          *AudioSystem
	wasPlaying     bool
	previousTrack  string
}

// NewAudioStopCommand creates a new AudioStopCommand
func NewAudioStopCommand(audio *AudioSystem) *AudioStopCommand {
	return &AudioStopCommand{audio: audio}
}

// Execute stops playback
func (c *AudioStopCommand) Execute() error {
	// Record previous state for undo
	c.wasPlaying = c.audio.isPlaying
	c.previousTrack = c.audio.track
	
	c.audio.Stop()
	return nil
}

// Undo resumes playback if it was playing
func (c *AudioStopCommand) Undo() error {
	if c.wasPlaying && c.previousTrack != "" {
		return c.audio.Play(c.previousTrack)
	}
	return nil
}

// String returns a description of the command
func (c *AudioStopCommand) String() string {
	return fmt.Sprintf("Stop playback on %s", c.audio.name)
}

// GarageDoorOpenCommand opens a garage door
type GarageDoorOpenCommand struct {
	door *GarageDoor
}

// NewGarageDoorOpenCommand creates a new GarageDoorOpenCommand
func NewGarageDoorOpenCommand(door *GarageDoor) *GarageDoorOpenCommand {
	return &GarageDoorOpenCommand{door: door}
}

// Execute opens the garage door
func (c *GarageDoorOpenCommand) Execute() error {
	c.door.Open()
	return nil
}

// Undo closes the garage door
func (c *GarageDoorOpenCommand) Undo() error {
	c.door.Close()
	return nil
}

// String returns a description of the command
func (c *GarageDoorOpenCommand) String() string {
	return fmt.Sprintf("Open %s", c.door.name)
}

// GarageDoorCloseCommand closes a garage door
type GarageDoorCloseCommand struct {
	door *GarageDoor
}

// NewGarageDoorCloseCommand creates a new GarageDoorCloseCommand
func NewGarageDoorCloseCommand(door *GarageDoor) *GarageDoorCloseCommand {
	return &GarageDoorCloseCommand{door: door}
}

// Execute closes the garage door
func (c *GarageDoorCloseCommand) Execute() error {
	c.door.Close()
	return nil
}

// Undo opens the garage door
func (c *GarageDoorCloseCommand) Undo() error {
	c.door.Open()
	return nil
}

// String returns a description of the command
func (c *GarageDoorCloseCommand) String() string {
	return fmt.Sprintf("Close %s", c.door.name)
}

// GarageLightOnCommand turns on the garage light
type GarageLightOnCommand struct {
	door *GarageDoor
}

// NewGarageLightOnCommand creates a new GarageLightOnCommand
func NewGarageLightOnCommand(door *GarageDoor) *GarageLightOnCommand {
	return &GarageLightOnCommand{door: door}
}

// Execute turns the garage light on
func (c *GarageLightOnCommand) Execute() error {
	c.door.LightOn()
	return nil
}

// Undo turns the garage light off
func (c *GarageLightOnCommand) Undo() error {
	c.door.LightOff()
	return nil
}

// String returns a description of the command
func (c *GarageLightOnCommand) String() string {
	return fmt.Sprintf("Turn %s light on", c.door.name)
}

// GarageLightOffCommand turns off the garage light
type GarageLightOffCommand struct {
	door *GarageDoor
}

// NewGarageLightOffCommand creates a new GarageLightOffCommand
func NewGarageLightOffCommand(door *GarageDoor) *GarageLightOffCommand {
	return &GarageLightOffCommand{door: door}
}

// Execute turns the garage light off
func (c *GarageLightOffCommand) Execute() error {
	c.door.LightOff()
	return nil
}

// Undo turns the garage light on
func (c *GarageLightOffCommand) Undo() error {
	c.door.LightOn()
	return nil
}

// String returns a description of the command
func (c *GarageLightOffCommand) String() string {
	return fmt.Sprintf("Turn %s light off", c.door.name)
}

// HomeSceneCommand is a macro command that sets up a home scene
type HomeSceneCommand struct {
	MacroCommand
	sceneName string
}

// NewHomeSceneCommand creates a new HomeSceneCommand with a name
func NewHomeSceneCommand(name string) *HomeSceneCommand {
	return &HomeSceneCommand{
		MacroCommand: MacroCommand{
			commands: make([]Command, 0),
			name:     name + " Scene",
		},
		sceneName: name,
	}
}

// String returns a description of the scene command
func (s *HomeSceneCommand) String() string {
	return fmt.Sprintf("Activate '%s' Scene (%d commands)", s.sceneName, len(s.commands))
}

// CeilingFanCommand controls a ceiling fan with multiple speeds
type CeilingFanCommand struct {
	name         string
	isOn         bool
	speed        int // 0=off, 1=low, 2=medium, 3=high
	previousIsOn bool
	previousSpeed int
}

// NewCeilingFanCommand creates a new CeilingFanCommand
func NewCeilingFanCommand(name string, speed int) *CeilingFanCommand {
	if speed < 0 {
		speed = 0
	} else if speed > 3 {
		speed = 3
	}
	
	return &CeilingFanCommand{
		name:  name,
		speed: speed,
	}
}

// Execute sets the fan to the desired speed
func (c *CeilingFanCommand) Execute() error {
	// Record previous state
	c.previousIsOn = c.isOn
	c.previousSpeed = c.speed
	
	// Set new state
	if c.speed == 0 {
		c.isOn = false
		fmt.Printf("%s ceiling fan turned OFF\n", c.name)
	} else {
		c.isOn = true
		speedNames := []string{"OFF", "LOW", "MEDIUM", "HIGH"}
		fmt.Printf("%s ceiling fan set to %s\n", c.name, speedNames[c.speed])
	}
	
	return nil
}

// Undo restores the previous fan state
func (c *CeilingFanCommand) Undo() error {
	// Restore previous state
	c.isOn = c.previousIsOn
	speedTemp := c.speed
	c.speed = c.previousSpeed
	c.previousSpeed = speedTemp
	
	// Report the change
	if !c.isOn {
		fmt.Printf("%s ceiling fan turned OFF\n", c.name)
	} else {
		speedNames := []string{"OFF", "LOW", "MEDIUM", "HIGH"}
		fmt.Printf("%s ceiling fan set to %s\n", c.name, speedNames[c.speed])
	}
	
	return nil
}

// String returns a description of the command
func (c *CeilingFanCommand) String() string {
	if c.speed == 0 {
		return fmt.Sprintf("Turn %s ceiling fan OFF", c.name)
	}
	
	speedNames := []string{"OFF", "LOW", "MEDIUM", "HIGH"}
	return fmt.Sprintf("Set %s ceiling fan to %s", c.name, speedNames[c.speed])
}
