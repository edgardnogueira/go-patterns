package main

import (
	"fmt"
)

func main() {
	fmt.Println("=== Open/Closed Principle Example ===")
	fmt.Println()
	
	fmt.Println("Before applying OCP:")
	fmt.Println("--------------------")
	demonstratePaymentProcessorBeforeOCP()
	fmt.Println()
	
	fmt.Println("After applying OCP:")
	fmt.Println("------------------")
	demonstratePaymentProcessorAfterOCP()
	
	fmt.Println()
	fmt.Println("This example demonstrates how to apply the Open/Closed Principle:")
	fmt.Println("1. Before: PaymentProcessor uses a switch statement to handle different payment types")
	fmt.Println("   - Adding a new payment method requires modifying the existing code")
	fmt.Println("   - This violates OCP because the class is not closed for modification")
	fmt.Println()
	fmt.Println("2. After: PaymentProcessor uses an interface for payment methods")
	fmt.Println("   - The processor is now closed for modification")
	fmt.Println("   - New payment methods can be added by implementing the interface")
	fmt.Println("   - We successfully added a cryptocurrency payment without changing existing code")
	fmt.Println()
	fmt.Println("Benefits of this approach:")
	fmt.Println("- Existing code remains stable and doesn't need to be retested")
	fmt.Println("- Reduced risk of introducing bugs in existing functionality")
	fmt.Println("- Better separation of concerns")
	fmt.Println("- More flexible and extensible design")
	fmt.Println("- Promotes the use of interfaces which is idiomatic in Go")
}
