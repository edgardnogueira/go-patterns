package abstract

import (
	"testing"
)

func TestModernGUIFactory(t *testing.T) {
	factory := &ModernGUIFactory{}

	// Test button creation
	button := factory.CreateButton()
	_, ok := button.(*ModernButton)
	if !ok {
		t.Error("ModernGUIFactory should create a ModernButton")
	}

	expectedRender := "Rendered a modern style button"
	if result := button.Render(); result != expectedRender {
		t.Errorf("Expected render '%s', but got '%s'", expectedRender, result)
	}

	expectedClick := "Modern button clicked with smooth animation"
	if result := button.OnClick(); result != expectedClick {
		t.Errorf("Expected onClick '%s', but got '%s'", expectedClick, result)
	}

	// Test checkbox creation
	checkbox := factory.CreateCheckbox()
	_, ok = checkbox.(*ModernCheckbox)
	if !ok {
		t.Error("ModernGUIFactory should create a ModernCheckbox")
	}

	expectedRender = "Rendered a modern style checkbox"
	if result := checkbox.Render(); result != expectedRender {
		t.Errorf("Expected render '%s', but got '%s'", expectedRender, result)
	}

	expectedToggle := "Modern checkbox toggled with sliding animation"
	if result := checkbox.Toggle(); result != expectedToggle {
		t.Errorf("Expected toggle '%s', but got '%s'", expectedToggle, result)
	}
}

func TestVintageGUIFactory(t *testing.T) {
	factory := &VintageGUIFactory{}

	// Test button creation
	button := factory.CreateButton()
	_, ok := button.(*VintageButton)
	if !ok {
		t.Error("VintageGUIFactory should create a VintageButton")
	}

	expectedRender := "Rendered a vintage style button"
	if result := button.Render(); result != expectedRender {
		t.Errorf("Expected render '%s', but got '%s'", expectedRender, result)
	}

	expectedClick := "Vintage button clicked with click sound"
	if result := button.OnClick(); result != expectedClick {
		t.Errorf("Expected onClick '%s', but got '%s'", expectedClick, result)
	}

	// Test checkbox creation
	checkbox := factory.CreateCheckbox()
	_, ok = checkbox.(*VintageCheckbox)
	if !ok {
		t.Error("VintageGUIFactory should create a VintageCheckbox")
	}

	expectedRender = "Rendered a vintage style checkbox"
	if result := checkbox.Render(); result != expectedRender {
		t.Errorf("Expected render '%s', but got '%s'", expectedRender, result)
	}

	expectedToggle := "Vintage checkbox toggled with mechanical sound"
	if result := checkbox.Toggle(); result != expectedToggle {
		t.Errorf("Expected toggle '%s', but got '%s'", expectedToggle, result)
	}
}

func TestCreateGUIFactory(t *testing.T) {
	tests := []struct {
		style          string
		expectedType   string
		expectedButton string
	}{
		{"modern", "*abstract.ModernGUIFactory", "*abstract.ModernButton"},
		{"vintage", "*abstract.VintageGUIFactory", "*abstract.VintageButton"},
		{"unknown", "*abstract.ModernGUIFactory", "*abstract.ModernButton"}, // Default case
	}

	for _, test := range tests {
		factory := CreateGUIFactory(test.style)
		actualType := getTypeName(factory)
		if actualType != test.expectedType {
			t.Errorf("For style '%s', expected factory type '%s', but got '%s'",
				test.style, test.expectedType, actualType)
		}

		button := factory.CreateButton()
		actualButtonType := getButtonTypeName(button)
		if actualButtonType != test.expectedButton {
			t.Errorf("For style '%s', expected button type '%s', but got '%s'",
				test.style, test.expectedButton, actualButtonType)
		}
	}
}

// Helper function to get the factory type name as a string
func getTypeName(v interface{}) string {
	switch v.(type) {
	case *ModernGUIFactory:
		return "*abstract.ModernGUIFactory"
	case *VintageGUIFactory:
		return "*abstract.VintageGUIFactory"
	default:
		return "unknown"
	}
}

// Helper function to get the button type name as a string
func getButtonTypeName(v interface{}) string {
	switch v.(type) {
	case *ModernButton:
		return "*abstract.ModernButton"
	case *VintageButton:
		return "*abstract.VintageButton"
	default:
		return "unknown"
	}
}
