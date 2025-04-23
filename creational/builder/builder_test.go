package builder

import (
	"strings"
	"testing"
)

func TestSportsCarBuilder(t *testing.T) {
	builder := NewSportsCarBuilder()
	director := &CarDirector{}
	director.SetBuilder(builder)
	director.BuildSportsCar()
	car := builder.GetCar()

	if car.Model != "Sports Model XZ" {
		t.Errorf("Expected model 'Sports Model XZ', got '%s'", car.Model)
	}

	if !strings.Contains(car.Engine, "V8 Turbo") {
		t.Errorf("Expected engine to contain 'V8 Turbo', got '%s'", car.Engine)
	}

	if !strings.Contains(car.Body, "Aerodynamic") {
		t.Errorf("Expected body to contain 'Aerodynamic', got '%s'", car.Body)
	}
}

func TestSUVBuilder(t *testing.T) {
	builder := NewSUVBuilder()
	director := &CarDirector{}
	director.SetBuilder(builder)
	director.BuildSUV()
	car := builder.GetCar()

	if car.Model != "Adventure SUV Pro" {
		t.Errorf("Expected model 'Adventure SUV Pro', got '%s'", car.Model)
	}

	if !strings.Contains(car.Engine, "V6 Engine with Towing") {
		t.Errorf("Expected engine to contain 'V6 Engine with Towing', got '%s'", car.Engine)
	}

	if !strings.Contains(car.Body, "High Ground Clearance") {
		t.Errorf("Expected body to contain 'High Ground Clearance', got '%s'", car.Body)
	}
}

func TestMinivanBuilder(t *testing.T) {
	builder := NewMinivanBuilder()
	director := &CarDirector{}
	director.SetBuilder(builder)
	director.BuildMinivan()
	car := builder.GetCar()

	if car.Model != "Family Comfort XL" {
		t.Errorf("Expected model 'Family Comfort XL', got '%s'", car.Model)
	}

	if !strings.Contains(car.Engine, "Efficient V6") {
		t.Errorf("Expected engine to contain 'Efficient V6', got '%s'", car.Engine)
	}

	if !strings.Contains(car.Body, "Sliding Doors") {
		t.Errorf("Expected body to contain 'Sliding Doors', got '%s'", car.Body)
	}

	if !strings.Contains(car.Interior, "7-Passenger") {
		t.Errorf("Expected interior to contain '7-Passenger', got '%s'", car.Interior)
	}
}

func TestCustomCarBuilder(t *testing.T) {
	builder := NewSportsCarBuilder()
	director := &CarDirector{}
	director.SetBuilder(builder)

	// Build a custom car
	customModel := "Custom Roadster"
	customEngine := "Electric Dual Motor"
	customTransmission := "Single-speed Automatic"
	customBody := "Convertible with Hardtop"
	customWheels := "19-inch Carbon Fiber Wheels"
	customInterior := "Premium Leather with Wood Accents"
	customElectronics := "Smart Dashboard with Voice Control"
	customSafety := "Advanced Driver Assistance Package"

	director.BuildCustomCar(
		customModel,
		customEngine,
		customTransmission,
		customBody,
		customWheels,
		customInterior,
		customElectronics,
		customSafety,
	)

	car := builder.GetCar()

	if car.Model != customModel {
		t.Errorf("Expected model '%s', got '%s'", customModel, car.Model)
	}

	if car.Engine != customEngine {
		t.Errorf("Expected engine '%s', got '%s'", customEngine, car.Engine)
	}

	if car.Transmission != customTransmission {
		t.Errorf("Expected transmission '%s', got '%s'", customTransmission, car.Transmission)
	}

	if car.Body != customBody {
		t.Errorf("Expected body '%s', got '%s'", customBody, car.Body)
	}

	if car.Wheels != customWheels {
		t.Errorf("Expected wheels '%s', got '%s'", customWheels, car.Wheels)
	}

	if car.Interior != customInterior {
		t.Errorf("Expected interior '%s', got '%s'", customInterior, car.Interior)
	}

	if car.Electronics != customElectronics {
		t.Errorf("Expected electronics '%s', got '%s'", customElectronics, car.Electronics)
	}

	if car.SafetyFeatures != customSafety {
		t.Errorf("Expected safety features '%s', got '%s'", customSafety, car.SafetyFeatures)
	}
}
