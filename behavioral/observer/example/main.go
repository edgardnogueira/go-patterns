package main

import (
	"fmt"
	"github.com/edgardnogueira/go-patterns/behavioral/observer"
)

func main() {
	fmt.Println("Observer Pattern Example")
	fmt.Println("=========================")
	fmt.Println("This example demonstrates a weather monitoring system where multiple displays")
	fmt.Println("observe and react to changes in weather data.\n")

	// Create the weather station (subject)
	fmt.Println("Setting up weather station...")
	weatherStation := observer.NewWeatherStation()

	// Create displays (observers)
	currentDisplay := observer.NewCurrentConditionsDisplay()
	statisticsDisplay := observer.NewStatisticsDisplay()
	forecastDisplay := observer.NewForecastDisplay()

	// Register observers with the subject
	fmt.Println("Registering observers:")
	fmt.Printf(" - %s\n", currentDisplay.GetName())
	fmt.Printf(" - %s\n", statisticsDisplay.GetName())
	fmt.Printf(" - %s\n", forecastDisplay.GetName())

	weatherStation.RegisterObserver(currentDisplay)
	weatherStation.RegisterObserver(statisticsDisplay)
	weatherStation.RegisterObserver(forecastDisplay)

	// First weather update
	fmt.Println("\n⛅ Weather Update 1: Temperature=80.0, Humidity=65.0, Pressure=30.4")
	fmt.Println("-------------------------------------------------------------------")
	weatherStation.SetMeasurements(80.0, 65.0, 30.4)

	fmt.Printf("🌡️  %s\n", currentDisplay.Display())
	fmt.Printf("📊 %s\n", statisticsDisplay.Display())
	fmt.Printf("🔮 %s\n", forecastDisplay.Display())

	// Second weather update
	fmt.Println("\n🌤️  Weather Update 2: Temperature=82.0, Humidity=70.0, Pressure=29.2")
	fmt.Println("-------------------------------------------------------------------")
	weatherStation.SetMeasurements(82.0, 70.0, 29.2)

	fmt.Printf("🌡️  %s\n", currentDisplay.Display())
	fmt.Printf("📊 %s\n", statisticsDisplay.Display())
	fmt.Printf("🔮 %s\n", forecastDisplay.Display())

	// Remove an observer
	fmt.Println("\n❌ Removing Current Conditions Display")
	weatherStation.RemoveObserver(currentDisplay)

	// Third weather update
	fmt.Println("\n🌧️  Weather Update 3: Temperature=78.0, Humidity=90.0, Pressure=29.4")
	fmt.Println("-------------------------------------------------------------------")
	weatherStation.SetMeasurements(78.0, 90.0, 29.4)

	fmt.Printf("🌡️  %s (not updated - observer removed)\n", currentDisplay.Display())
	fmt.Printf("📊 %s\n", statisticsDisplay.Display())
	fmt.Printf("🔮 %s\n", forecastDisplay.Display())

	fmt.Println("\nThe Observer pattern allows objects to be notified when state changes.")
	fmt.Println("It provides a loosely coupled design between the subject and its observers.")
}
