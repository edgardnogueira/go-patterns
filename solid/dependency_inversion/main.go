package main

import (
	"fmt"
)

func main() {
	fmt.Println("=== Dependency Inversion Principle Example ===")
	fmt.Println()
	
	fmt.Println("Before applying DIP:")
	fmt.Println("--------------------")
	demonstrateNotificationBeforeDIP()
	fmt.Println()
	
	fmt.Println("After applying DIP:")
	fmt.Println("------------------")
	demonstrateNotificationAfterDIP()
	
	fmt.Println()
	fmt.Println("This example demonstrates how to apply the Dependency Inversion Principle:")
	fmt.Println("1. Before: NotificationService depends directly on concrete implementations")
	fmt.Println("   - UserService directly creates and depends on NotificationService")
	fmt.Println("   - Adding a new notification method requires modifying existing code")
	fmt.Println("   - Testing is difficult because dependencies are hard-coded")
	fmt.Println()
	fmt.Println("2. After: Both high and low-level modules depend on abstractions")
	fmt.Println("   - NotificationSender interface abstracts the sending behavior")
	fmt.Println("   - NotificationService depends on this abstraction, not concrete implementations")
	fmt.Println("   - Dependencies are injected rather than created internally")
	fmt.Println("   - New notification methods can be added without changing existing code")
	fmt.Println()
	fmt.Println("Benefits of this approach:")
	fmt.Println("- Modules are loosely coupled and can evolve independently")
	fmt.Println("- New implementations can be added without modifying existing code")
	fmt.Println("- Dependencies are explicit and injected, making testing much easier")
	fmt.Println("- Code is more modular and reusable")
	fmt.Println("- Implementation details are hidden behind abstractions")
}
