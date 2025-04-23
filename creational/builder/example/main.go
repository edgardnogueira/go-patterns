package main

import (
	"fmt"
	"github.com/edgardnogueira/go-patterns/creational/builder"
)

func main() {
	fmt.Println("Builder Pattern Example")
	fmt.Println("=======================")

	// Create a director
	director := &builder.CarDirector{}

	// Create a sports car
	sportsCarBuilder := builder.NewSportsCarBuilder()
	director.SetBuilder(sportsCarBuilder)
	director.BuildSportsCar()
	sportsCar := sportsCarBuilder.GetCar()

	fmt.Println("\n🏎️ Sports Car Built:")
	fmt.Println("------------------")
	fmt.Println(sportsCar)

	// Create an SUV
	suvBuilder := builder.NewSUVBuilder()
	director.SetBuilder(suvBuilder)
	director.BuildSUV()
	suv := suvBuilder.GetCar()

	fmt.Println("\n🚙 SUV Built:")
	fmt.Println("------------")
	fmt.Println(suv)

	// Create a minivan
	minivanBuilder := builder.NewMinivanBuilder()
	director.SetBuilder(minivanBuilder)
	director.BuildMinivan()
	minivan := minivanBuilder.GetCar()

	fmt.Println("\n🚐 Minivan Built:")
	fmt.Println("-----------------")
	fmt.Println(minivan)

	// Create a custom car
	fmt.Println("\n🔧 Building a Custom Car...")
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

	fmt.Println("\n🏁 Custom Car Built:")
	fmt.Println("-------------------")
	fmt.Println(customCar)

	fmt.Println("\nNote: The Builder pattern allows for step-by-step construction of different")
	fmt.Println("car types using the same construction process.")
}
