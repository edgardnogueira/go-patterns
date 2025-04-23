package observer

import (
	"fmt"
	"strings"
	"testing"
)

func TestWeatherStation(t *testing.T) {
	// Create the weather station (subject)
	weatherStation := NewWeatherStation()

	// Create displays (observers)
	currentDisplay := NewCurrentConditionsDisplay()
	statisticsDisplay := NewStatisticsDisplay()
	forecastDisplay := NewForecastDisplay()

	// Register observers with the subject
	weatherStation.RegisterObserver(currentDisplay)
	weatherStation.RegisterObserver(statisticsDisplay)
	weatherStation.RegisterObserver(forecastDisplay)

	// Test initial measurements
	weatherStation.SetMeasurements(80.0, 65.0, 30.4)

	// Test current conditions display
	currentResult := currentDisplay.Display()
	expectedCurrent := "Current conditions: 80.0F degrees and 65% humidity"
	if currentResult != expectedCurrent {
		t.Errorf("Expected '%s', got '%s'", expectedCurrent, currentResult)
	}

	// Test statistics display after first reading
	statsResult := statisticsDisplay.Display()
	if !strings.Contains(statsResult, "Avg/Max/Min temperature: 80.0/80.0/80.0") {
		t.Errorf("Expected stats to contain '80.0/80.0/80.0', got '%s'", statsResult)
	}

	// Test forecast display after first reading
	// Initial pressure change can't be tested meaningfully since we start with a default

	// Change measurements and test again
	weatherStation.SetMeasurements(82.0, 70.0, 29.2)

	// Test current conditions updated
	currentResult = currentDisplay.Display()
	expectedCurrent = "Current conditions: 82.0F degrees and 70% humidity"
	if currentResult != expectedCurrent {
		t.Errorf("Expected '%s', got '%s'", expectedCurrent, currentResult)
	}

	// Test statistics display after second reading
	statsResult = statisticsDisplay.Display()
	if !strings.Contains(statsResult, "Avg/Max/Min temperature: 81.0/82.0/80.0") {
		t.Errorf("Expected stats to contain '81.0/82.0/80.0', got '%s'", statsResult)
	}

	// Test forecast display after pressure drop
	forecastResult := forecastDisplay.Display()
	expectedForecast := "Forecast: Watch out for cooler, rainy weather"
	if forecastResult != expectedForecast {
		t.Errorf("Expected '%s', got '%s'", expectedForecast, forecastResult)
	}

	// Test removing an observer
	weatherStation.RemoveObserver(currentDisplay)

	// Change measurements again
	weatherStation.SetMeasurements(78.0, 90.0, 29.2)

	// Current display should not be updated (still showing previous values)
	currentResult = currentDisplay.Display()
	if currentResult != expectedCurrent {
		t.Errorf("Expected '%s' (not updated after removal), got '%s'", expectedCurrent, currentResult)
	}

	// Statistics display should still be updated
	statsResult = statisticsDisplay.Display()
	if !strings.Contains(statsResult, "Avg/Max/Min temperature: 80.0/82.0/78.0") {
		t.Errorf("Expected stats to contain '80.0/82.0/78.0', got '%s'", statsResult)
	}
}
