package builder

import (
	"fmt"
)

// This file contains example usage of the builder pattern

// ExampleBuilderPattern demonstrates the builder pattern in action
func ExampleBuilderPattern() {
	// Create a director
	director := &CarDirector{}

	// Create a sports car
	sportsCarBuilder := NewSportsCarBuilder()
	director.SetBuilder(sportsCarBuilder)
	director.BuildSportsCar()
	sportsCar := sportsCarBuilder.GetCar()

	fmt.Println("Sports Car Built:")
	fmt.Println(sportsCar)

	fmt.Println()

	// Create an SUV
	suvBuilder := NewSUVBuilder()
	director.SetBuilder(suvBuilder)
	director.BuildSUV()
	suv := suvBuilder.GetCar()

	fmt.Println("SUV Built:")
	fmt.Println(suv)

	fmt.Println()

	// Create a minivan
	minivanBuilder := NewMinivanBuilder()
	director.SetBuilder(minivanBuilder)
	director.BuildMinivan()
	minivan := minivanBuilder.GetCar()

	fmt.Println("Minivan Built:")
	fmt.Println(minivan)

	fmt.Println()

	// Create a custom car
	director.SetBuilder(sportsCarBuilder)
	director.BuildCustomCar(
		"Custom Roadster",
		"Electric Dual Motor",
		"Single-speed Automatic",
		"Convertible with Hardtop",
		"19-inch Carbon Fiber Wheels",
		"Premium Leather with Wood Accents",
		"Smart Dashboard with Voice Control",
		"Advanced Driver Assistance Package",
	)
	customCar := sportsCarBuilder.GetCar()

	fmt.Println("Custom Car Built:")
	fmt.Println(customCar)
}
