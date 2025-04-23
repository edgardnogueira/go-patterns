package main

import (
	"fmt"
	"github.com/edgardnogueira/go-patterns/creational/abstract"
)

func main() {
	fmt.Println("Abstract Factory Pattern Example")
	fmt.Println("==================================")

	// Create a client application with Modern GUI
	modernFactory := abstract.CreateGUIFactory("modern")
	modernApp := abstract.NewApplication(modernFactory)

	fmt.Println("\nModern UI:")
	fmt.Println("----------")

	fmt.Println("Rendering UI components:")
	for _, result := range modernApp.RenderUI() {
		fmt.Println("  » " + result)
	}

	fmt.Println("\nUser interactions:")
	for _, result := range modernApp.ExecuteActions() {
		fmt.Println("  » " + result)
	}

	// Create a client application with Vintage GUI
	vintageFactory := abstract.CreateGUIFactory("vintage")
	vintageApp := abstract.NewApplication(vintageFactory)

	fmt.Println("\nVintage UI:")
	fmt.Println("------------")

	fmt.Println("Rendering UI components:")
	for _, result := range vintageApp.RenderUI() {
		fmt.Println("  » " + result)
	}

	fmt.Println("\nUser interactions:")
	for _, result := range vintageApp.ExecuteActions() {
		fmt.Println("  » " + result)
	}

	fmt.Println("\nThe client code works with factories and products through abstract interfaces,")
	fmt.Println("so it doesn't matter which factory or product variant is used.")
}
