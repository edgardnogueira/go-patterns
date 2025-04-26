package command

import (
	"fmt"
	"strings"
)

// Light represents a light that can be turned on and off
type Light struct {
	name      string
	isOn      bool
	brightness int // 0-100%
}

// NewLight creates a new Light with the given name
func NewLight(name string) *Light {
	return &Light{
		name:      name,
		isOn:      false,
		brightness: 100,
	}
}

// TurnOn turns the light on
func (l *Light) TurnOn() {
	l.isOn = true
	fmt.Printf("%s light is now ON\n", l.name)
}

// TurnOff turns the light off
func (l *Light) TurnOff() {
	l.isOn = false
	fmt.Printf("%s light is now OFF\n", l.name)
}

// SetBrightness sets the light brightness
func (l *Light) SetBrightness(level int) {
	if level < 0 {
		level = 0
	} else if level > 100 {
		level = 100
	}
	
	previousBrightness := l.brightness
	l.brightness = level
	
	// If light was off and we're setting brightness > 0, turn it on
	if !l.isOn && level > 0 {
		l.isOn = true
	}
	
	// If setting brightness to 0, turn light off
	if level == 0 {
		l.isOn = false
	}
	
	fmt.Printf("%s light brightness changed from %d%% to %d%%\n", 
		l.name, previousBrightness, l.brightness)
}

// GetStatus returns the current status of the light
func (l *Light) GetStatus() string {
	status := "OFF"
	if l.isOn {
		status = fmt.Sprintf("ON (Brightness: %d%%)", l.brightness)
	}
	return fmt.Sprintf("%s light: %s", l.name, status)
}

// Thermostat represents a thermostat that can be set to different temperatures
type Thermostat struct {
	name         string
	temperature  int // in degrees
	isOn         bool
	mode         string // "heat", "cool", "auto"
}

// NewThermostat creates a new Thermostat with the given name
func NewThermostat(name string) *Thermostat {
	return &Thermostat{
		name:        name,
		temperature: 72,
		isOn:        false,
		mode:        "auto",
	}
}

// TurnOn turns the thermostat on
func (t *Thermostat) TurnOn() {
	t.isOn = true
	fmt.Printf("%s thermostat is now ON (Mode: %s, Temp: %d°)\n", 
		t.name, t.mode, t.temperature)
}

// TurnOff turns the thermostat off
func (t *Thermostat) TurnOff() {
	t.isOn = false
	fmt.Printf("%s thermostat is now OFF\n", t.name)
}

// SetTemperature sets the thermostat temperature
func (t *Thermostat) SetTemperature(temp int) {
	prevTemp := t.temperature
	t.temperature = temp
	
	// If thermostat was off and we're setting a temperature, turn it on
	if !t.isOn {
		t.isOn = true
	}
	
	fmt.Printf("%s thermostat temperature changed from %d° to %d°\n", 
		t.name, prevTemp, t.temperature)
}

// SetMode sets the thermostat mode
func (t *Thermostat) SetMode(mode string) error {
	mode = strings.ToLower(mode)
	
	if mode != "heat" && mode != "cool" && mode != "auto" {
		return fmt.Errorf("invalid mode: %s (must be 'heat', 'cool', or 'auto')", mode)
	}
	
	prevMode := t.mode
	t.mode = mode
	
	fmt.Printf("%s thermostat mode changed from %s to %s\n", 
		t.name, prevMode, t.mode)
	return nil
}

// GetStatus returns the current status of the thermostat
func (t *Thermostat) GetStatus() string {
	status := "OFF"
	if t.isOn {
		status = fmt.Sprintf("ON (Mode: %s, Temp: %d°)", t.mode, t.temperature)
	}
	return fmt.Sprintf("%s thermostat: %s", t.name, status)
}

// AudioSystem represents a home audio system
type AudioSystem struct {
	name     string
	isOn     bool
	volume   int // 0-100%
	source   string
	isPlaying bool
	track    string
}

// NewAudioSystem creates a new AudioSystem with the given name
func NewAudioSystem(name string) *AudioSystem {
	return &AudioSystem{
		name:     name,
		isOn:     false,
		volume:   50,
		source:   "Bluetooth",
		isPlaying: false,
		track:    "",
	}
}

// TurnOn turns the audio system on
func (a *AudioSystem) TurnOn() {
	a.isOn = true
	fmt.Printf("%s audio system is now ON\n", a.name)
}

// TurnOff turns the audio system off
func (a *AudioSystem) TurnOff() {
	if a.isPlaying {
		a.Stop()
	}
	a.isOn = false
	fmt.Printf("%s audio system is now OFF\n", a.name)
}

// SetVolume sets the volume level
func (a *AudioSystem) SetVolume(level int) {
	if level < 0 {
		level = 0
	} else if level > 100 {
		level = 100
	}
	
	prevVolume := a.volume
	a.volume = level
	
	fmt.Printf("%s audio system volume changed from %d%% to %d%%\n", 
		a.name, prevVolume, a.volume)
}

// SetSource sets the audio source
func (a *AudioSystem) SetSource(source string) {
	prevSource := a.source
	a.source = source
	
	fmt.Printf("%s audio system source changed from %s to %s\n", 
		a.name, prevSource, a.source)
}

// Play starts playback of the given track
func (a *AudioSystem) Play(track string) error {
	if !a.isOn {
		return fmt.Errorf("cannot play: %s audio system is off", a.name)
	}
	
	a.isPlaying = true
	a.track = track
	
	fmt.Printf("%s audio system now playing: %s\n", a.name, track)
	return nil
}

// Stop stops the current playback
func (a *AudioSystem) Stop() {
	if a.isPlaying {
		a.isPlaying = false
		fmt.Printf("%s audio system stopped playing: %s\n", a.name, a.track)
	}
}

// GetStatus returns the current status of the audio system
func (a *AudioSystem) GetStatus() string {
	if !a.isOn {
		return fmt.Sprintf("%s audio system: OFF", a.name)
	}
	
	playStatus := "stopped"
	if a.isPlaying {
		playStatus = fmt.Sprintf("playing '%s'", a.track)
	}
	
	return fmt.Sprintf("%s audio system: ON (Volume: %d%%, Source: %s, Status: %s)", 
		a.name, a.volume, a.source, playStatus)
}

// GarageDoor represents a garage door that can be opened and closed
type GarageDoor struct {
	name      string
	isOpen    bool
	hasLight  bool
	lightOn   bool
}

// NewGarageDoor creates a new GarageDoor with the given name
func NewGarageDoor(name string) *GarageDoor {
	return &GarageDoor{
		name:     name,
		isOpen:   false,
		hasLight: true,
		lightOn:  false,
	}
}

// Open opens the garage door
func (g *GarageDoor) Open() {
	if !g.isOpen {
		fmt.Printf("%s garage door is opening...\n", g.name)
		g.isOpen = true
	}
}

// Close closes the garage door
func (g *GarageDoor) Close() {
	if g.isOpen {
		fmt.Printf("%s garage door is closing...\n", g.name)
		g.isOpen = false
	}
}

// LightOn turns the garage light on
func (g *GarageDoor) LightOn() {
	if g.hasLight && !g.lightOn {
		g.lightOn = true
		fmt.Printf("%s garage light is now ON\n", g.name)
	}
}

// LightOff turns the garage light off
func (g *GarageDoor) LightOff() {
	if g.hasLight && g.lightOn {
		g.lightOn = false
		fmt.Printf("%s garage light is now OFF\n", g.name)
	}
}

// GetStatus returns the current status of the garage door
func (g *GarageDoor) GetStatus() string {
	doorStatus := "CLOSED"
	if g.isOpen {
		doorStatus = "OPEN"
	}
	
	lightStatus := "N/A"
	if g.hasLight {
		if g.lightOn {
			lightStatus = "ON"
		} else {
			lightStatus = "OFF"
		}
	}
	
	return fmt.Sprintf("%s garage door: %s, Light: %s", g.name, doorStatus, lightStatus)
}
