package observer

import (
	"fmt"
)

// WeatherData represents the measurement data
type WeatherData struct {
	Temperature float64
	Humidity    float64
	Pressure    float64
}

// WeatherObserver interface for all weather observers
type WeatherObserver interface {
	Update(data WeatherData)
	GetName() string
}

// WeatherSubject interface defines methods for managing observers
type WeatherSubject interface {
	RegisterObserver(observer WeatherObserver)
	RemoveObserver(observer WeatherObserver)
	NotifyObservers()
}

// WeatherStation implements the WeatherSubject interface
type WeatherStation struct {
	observers []WeatherObserver
	data      WeatherData
}

// NewWeatherStation creates a new WeatherStation
func NewWeatherStation() *WeatherStation {
	return &WeatherStation{
		observers: make([]WeatherObserver, 0),
		data:      WeatherData{},
	}
}

// RegisterObserver adds an observer to the list
func (ws *WeatherStation) RegisterObserver(observer WeatherObserver) {
	ws.observers = append(ws.observers, observer)
}

// RemoveObserver removes an observer from the list
func (ws *WeatherStation) RemoveObserver(observer WeatherObserver) {
	for i, obs := range ws.observers {
		if obs == observer {
			ws.observers = append(ws.observers[:i], ws.observers[i+1:]...)
			break
		}
	}
}

// NotifyObservers notifies all registered observers about weather changes
func (ws *WeatherStation) NotifyObservers() {
	for _, observer := range ws.observers {
		observer.Update(ws.data)
	}
}

// SetMeasurements updates the weather data and notifies observers
func (ws *WeatherStation) SetMeasurements(temperature, humidity, pressure float64) {
	ws.data.Temperature = temperature
	ws.data.Humidity = humidity
	ws.data.Pressure = pressure
	ws.NotifyObservers()
}

// GetMeasurements returns the current weather data
func (ws *WeatherStation) GetMeasurements() WeatherData {
	return ws.data
}

// CurrentConditionsDisplay displays current weather conditions
type CurrentConditionsDisplay struct {
	temperature float64
	humidity    float64
}

// NewCurrentConditionsDisplay creates a new CurrentConditionsDisplay
func NewCurrentConditionsDisplay() *CurrentConditionsDisplay {
	return &CurrentConditionsDisplay{}
}

// Update updates the display with new weather data
func (d *CurrentConditionsDisplay) Update(data WeatherData) {
	d.temperature = data.Temperature
	d.humidity = data.Humidity
}

// GetName returns the name of the display
func (d *CurrentConditionsDisplay) GetName() string {
	return "Current Conditions Display"
}

// Display outputs the current weather conditions
func (d *CurrentConditionsDisplay) Display() string {
	return "Current conditions: " + 
		formattedTemp(d.temperature) + "F degrees and " + 
		formattedHumidity(d.humidity) + "% humidity"
}

// StatisticsDisplay displays weather statistics
type StatisticsDisplay struct {
	maxTemp     float64
	minTemp     float64
	tempSum     float64
	numReadings int
}

// NewStatisticsDisplay creates a new StatisticsDisplay
func NewStatisticsDisplay() *StatisticsDisplay {
	return &StatisticsDisplay{
		maxTemp:     -999.0,
		minTemp:     999.0,
		tempSum:     0.0,
		numReadings: 0,
	}
}

// Update updates the statistics with new weather data
func (d *StatisticsDisplay) Update(data WeatherData) {
	d.tempSum += data.Temperature
	d.numReadings++

	if data.Temperature > d.maxTemp {
		d.maxTemp = data.Temperature
	}

	if data.Temperature < d.minTemp {
		d.minTemp = data.Temperature
	}
}

// GetName returns the name of the display
func (d *StatisticsDisplay) GetName() string {
	return "Statistics Display"
}

// Display outputs the weather statistics
func (d *StatisticsDisplay) Display() string {
	avg := d.tempSum / float64(d.numReadings)
	return "Avg/Max/Min temperature: " + 
		formattedTemp(avg) + "/" + 
		formattedTemp(d.maxTemp) + "/" + 
		formattedTemp(d.minTemp)
}

// ForecastDisplay displays weather forecast
type ForecastDisplay struct {
	currentPressure float64
	lastPressure    float64
}

// NewForecastDisplay creates a new ForecastDisplay
func NewForecastDisplay() *ForecastDisplay {
	return &ForecastDisplay{
		currentPressure: 29.92, // starting pressure
		lastPressure:    0.0,
	}
}

// Update updates the forecast with new weather data
func (d *ForecastDisplay) Update(data WeatherData) {
	d.lastPressure = d.currentPressure
	d.currentPressure = data.Pressure
}

// GetName returns the name of the display
func (d *ForecastDisplay) GetName() string {
	return "Forecast Display"
}

// Display outputs the weather forecast
func (d *ForecastDisplay) Display() string {
	var forecast string

	if d.currentPressure > d.lastPressure {
		forecast = "Improving weather on the way!"
	} else if d.currentPressure == d.lastPressure {
		forecast = "More of the same"
	} else {
		forecast = "Watch out for cooler, rainy weather"
	}

	return "Forecast: " + forecast
}

// Helper function to format temperature
func formattedTemp(temp float64) string {
	return fmt.Sprintf("%.1f", temp)
}

// Helper function to format humidity
func formattedHumidity(humidity float64) string {
	return fmt.Sprintf("%.0f", humidity)
}