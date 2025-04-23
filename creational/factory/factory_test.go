package factory

import (
	"testing"
)

func TestRoadLogistics(t *testing.T) {
	logistics := &RoadLogistics{}
	transport := logistics.CreateTransport()

	// Check if we got a Truck
	_, ok := transport.(*Truck)
	if !ok {
		t.Error("RoadLogistics should create a Truck")
	}

	expected := "Delivering by land in a truck"
	result := transport.Deliver()
	if result != expected {
		t.Errorf("Expected '%s', but got '%s'", expected, result)
	}

	expected = "Road logistics: Delivering by land in a truck"
	result = logistics.PlanDelivery()
	if result != expected {
		t.Errorf("Expected '%s', but got '%s'", expected, result)
	}
}

func TestSeaLogistics(t *testing.T) {
	logistics := &SeaLogistics{}
	transport := logistics.CreateTransport()

	// Check if we got a Ship
	_, ok := transport.(*Ship)
	if !ok {
		t.Error("SeaLogistics should create a Ship")
	}

	expected := "Delivering by sea in a ship"
	result := transport.Deliver()
	if result != expected {
		t.Errorf("Expected '%s', but got '%s'", expected, result)
	}

	expected = "Sea logistics: Delivering by sea in a ship"
	result = logistics.PlanDelivery()
	if result != expected {
		t.Errorf("Expected '%s', but got '%s'", expected, result)
	}
}

func TestCreateLogistics(t *testing.T) {
	tests := []struct {
		logisticsType string
		expectedType  string
	}{
		{"road", "*factory.RoadLogistics"},
		{"sea", "*factory.SeaLogistics"},
		{"unknown", "*factory.RoadLogistics"}, // Default case
	}

	for _, test := range tests {
		logistics := CreateLogistics(test.logisticsType)
		actualType := getTypeName(logistics)
		if actualType != test.expectedType {
			t.Errorf("For logistics type '%s', expected type '%s', but got '%s'",
				test.logisticsType, test.expectedType, actualType)
		}
	}
}

// Helper function to get the type name as a string
func getTypeName(v interface{}) string {
	switch v.(type) {
	case *RoadLogistics:
		return "*factory.RoadLogistics"
	case *SeaLogistics:
		return "*factory.SeaLogistics"
	default:
		return "unknown"
	}
}
