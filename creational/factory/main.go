package factory

import (
	"fmt"
)

// This file contains example usage of the factory pattern

// ExampleFactoryMethod demonstrates the factory method pattern in action
func ExampleFactoryMethod() {
	// Create road logistics service
	roadLogistics := CreateLogistics("road")
	fmt.Println(roadLogistics.PlanDelivery())

	// Create sea logistics service
	seaLogistics := CreateLogistics("sea")
	fmt.Println(seaLogistics.PlanDelivery())

	// The client code works with any creator or product following the interfaces
	deliverProduct(roadLogistics)
	deliverProduct(seaLogistics)
}

// deliverProduct is a client function that works with any LogisticsService
func deliverProduct(logistics LogisticsService) {
	fmt.Printf("Client: I'm not aware of the logistics class, but it still works.\n")
	fmt.Printf("Client: %s\n", logistics.PlanDelivery())
}
