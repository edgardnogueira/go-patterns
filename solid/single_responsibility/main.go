package main

import (
	"fmt"
)

func main() {
	fmt.Println("=== Single Responsibility Principle Example ===")
	fmt.Println()
	
	fmt.Println("Before applying SRP:")
	fmt.Println("--------------------")
	demonstrateUserManagerBeforeSRP()
	fmt.Println()
	
	fmt.Println("After applying SRP:")
	fmt.Println("------------------")
	demonstrateUserManagerAfterSRP()
	
	fmt.Println()
	fmt.Println("This example demonstrates how to apply the Single Responsibility Principle:")
	fmt.Println("1. Before: UserManager class handles validation, persistence, email, and logging")
	fmt.Println("2. After: Responsibilities are separated into specialized components:")
	fmt.Println("   - UserValidator: handles only validation logic")
	fmt.Println("   - UserRepository: handles only data persistence")
	fmt.Println("   - EmailService: handles only email notifications")
	fmt.Println("   - Logger: handles only logging")
	fmt.Println("   - UserService: coordinates the above components")
	fmt.Println()
	fmt.Println("Benefits of this approach:")
	fmt.Println("- Each component is simpler and focused on one task")
	fmt.Println("- Components can be tested in isolation")
	fmt.Println("- Components can be modified without affecting others")
	fmt.Println("- New implementations of each component can be swapped in")
}
