package main

import (
	"fmt"
)

func main() {
	fmt.Println("=== Liskov Substitution Principle Example ===")
	fmt.Println()
	
	fmt.Println("Before applying LSP:")
	fmt.Println("--------------------")
	demonstrateFileStorageBeforeLSP()
	fmt.Println()
	
	fmt.Println("After applying LSP:")
	fmt.Println("------------------")
	demonstrateFileStorageAfterLSP()
	
	fmt.Println()
	fmt.Println("This example demonstrates how to apply the Liskov Substitution Principle:")
	fmt.Println("1. Before: ReadOnlyFileStorage implements FileStorage but throws errors for save/delete")
	fmt.Println("   - This violates LSP because you can't substitute it for a FileStorage")
	fmt.Println("   - It breaks client code expecting certain behaviors")
	fmt.Println()
	fmt.Println("2. After: We separate interfaces based on capabilities")
	fmt.Println("   - ReadableStorage for read operations")
	fmt.Println("   - WritableStorage for write operations")
	fmt.Println("   - FileStorage combines both interfaces")
	fmt.Println("   - ReadOnlyFileStorage only implements ReadableStorage")
	fmt.Println("   - Clients depend only on the interfaces they need")
	fmt.Println()
	fmt.Println("Benefits of this approach:")
	fmt.Println("- Type safety at compile-time rather than runtime errors")
	fmt.Println("- Better interface segregation (which also relates to the Interface Segregation Principle)")
	fmt.Println("- Client code only depends on the operations it actually needs")
	fmt.Println("- Clear contract for implementations")
	fmt.Println("- Proper substitutability for all interface implementations")
}
