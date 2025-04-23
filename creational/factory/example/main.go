package main

import (
	"fmt"
	"github.com/edgardnogueira/go-patterns/creational/factory"
)

func main() {
	fmt.Println("Factory Method Pattern Example")
	fmt.Println("================================")

	// Create road logistics
	roadLogistics := factory.CreateLogistics("road")
	fmt.Println(roadLogistics.PlanDelivery())

	// Create sea logistics
	seaLogistics := factory.CreateLogistics("sea")
	fmt.Println(seaLogistics.PlanDelivery())

	fmt.Println("\nUsing client code that works with any logistics service:")
	fmt.Println("---------------------------------------------------")
	deliverProduct(roadLogistics)
	deliverProduct(seaLogistics)
}

// deliverProduct is a client function that works with any LogisticsService
func deliverProduct(logistics factory.LogisticsService) {
	fmt.Printf("Client: I'm not aware of the logistics class, but it still works.\n")
	fmt.Printf("Client: %s\n", logistics.PlanDelivery())
}
