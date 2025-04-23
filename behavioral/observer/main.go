package observer

import (
	"fmt"
)

// This file contains example usage of the observer pattern

// ExampleObserverPattern demonstrates the observer pattern in action
func ExampleObserverPattern() {
	// Create the weather station (subject)
	weatherStation := NewWeatherStation()

	// Create displays (observers)
	currentDisplay := NewCurrentConditionsDisplay()
	statisticsDisplay := NewStatisticsDisplay()
	forecastDisplay := NewForecastDisplay()

	// Register observers with the subject
	fmt.Println("Registering observers...")
	weatherStation.RegisterObserver(currentDisplay)
	weatherStation.RegisterObserver(statisticsDisplay)
	weatherStation.RegisterObserver(forecastDisplay)

	// First weather update
	fmt.Println("\nWeather update 1:")
	fmt.Println("------------------")
	weatherStation.SetMeasurements(80.0, 65.0, 30.4)

	fmt.Println(currentDisplay.Display())
	fmt.Println(statisticsDisplay.Display())
	fmt.Println(forecastDisplay.Display())

	// Second weather update
	fmt.Println("\nWeather update 2:")
	fmt.Println("------------------")
	weatherStation.SetMeasurements(82.0, 70.0, 29.2)

	fmt.Println(currentDisplay.Display())
	fmt.Println(statisticsDisplay.Display())
	fmt.Println(forecastDisplay.Display())

	// Remove an observer
	fmt.Println("\nRemoving current conditions display...")
	weatherStation.RemoveObserver(currentDisplay)

	// Third weather update (current display should not be updated)
	fmt.Println("\nWeather update 3:")
	fmt.Println("------------------")
	weatherStation.SetMeasurements(78.0, 90.0, 29.4)

	fmt.Println(currentDisplay.Display()) // Should still show update 2 data
	fmt.Println(statisticsDisplay.Display())
	fmt.Println(forecastDisplay.Display())
}
