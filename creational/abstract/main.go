package abstract

import (
	"fmt"
)

// This file contains example usage of the abstract factory pattern

// Application represents a client using the GUI factories and products
type Application struct {
	button   Button
	checkbox Checkbox
}

// NewApplication creates a new application with UI elements from the specified factory
func NewApplication(factory GUIFactory) *Application {
	return &Application{
		button:   factory.CreateButton(),
		checkbox: factory.CreateCheckbox(),
	}
}

// RenderUI renders all UI components of the application
func (a *Application) RenderUI() []string {
	result := []string{
		a.button.Render(),
		a.checkbox.Render(),
	}
	return result
}

// ExecuteActions simulates user interactions with UI components
func (a *Application) ExecuteActions() []string {
	result := []string{
		a.button.OnClick(),
		a.checkbox.Toggle(),
	}
	return result
}

// ExampleAbstractFactory demonstrates the abstract factory pattern in action
func ExampleAbstractFactory() {
	// Create a modern style application
	modernFactory := CreateGUIFactory("modern")
	modernApp := NewApplication(modernFactory)

	fmt.Println("Modern UI:")
	for _, result := range modernApp.RenderUI() {
		fmt.Println("- " + result)
	}
	for _, result := range modernApp.ExecuteActions() {
		fmt.Println("- " + result)
	}

	fmt.Println()

	// Create a vintage style application
	vintageFactory := CreateGUIFactory("vintage")
	vintageApp := NewApplication(vintageFactory)

	fmt.Println("Vintage UI:")
	for _, result := range vintageApp.RenderUI() {
		fmt.Println("- " + result)
	}
	for _, result := range vintageApp.ExecuteActions() {
		fmt.Println("- " + result)
	}
}
