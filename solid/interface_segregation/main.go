package main

import (
	"fmt"
)

func main() {
	fmt.Println("=== Interface Segregation Principle Example ===")
	fmt.Println()
	
	fmt.Println("Before applying ISP:")
	fmt.Println("--------------------")
	demonstrateWorkerBeforeISP()
	fmt.Println()
	
	fmt.Println("After applying ISP:")
	fmt.Println("------------------")
	demonstrateWorkerAfterISP()
	
	fmt.Println()
	fmt.Println("This example demonstrates how to apply the Interface Segregation Principle:")
	fmt.Println("1. Before: We had a \"fat\" Worker interface with methods for all possible worker types")
	fmt.Println("   - Every worker had to implement every method, even if not applicable")
	fmt.Println("   - This led to empty implementations and inappropriate behaviors")
	fmt.Println("   - Adding a new method to Worker forced all implementations to change")
	fmt.Println()
	fmt.Println("2. After: We split the interface into smaller, focused interfaces")
	fmt.Println("   - BasicWorker: core methods all workers implement")
	fmt.Println("   - TeamMember: methods for attending meetings")
	fmt.Println("   - OvertimeWorker: methods for working extra hours")
	fmt.Println("   - ManagerRole: methods specific to management")
	fmt.Println("   - Each worker type only implements relevant interfaces")
	fmt.Println()
	fmt.Println("Benefits of this approach:")
	fmt.Println("- Clients depend only on the methods they actually use")
	fmt.Println("- Workers only implement methods relevant to their role")
	fmt.Println("- Type safety ensures workers are only used for appropriate tasks")
	fmt.Println("- Adding new worker types or capabilities is easier")
	fmt.Println("- Better alignment with the Single Responsibility Principle")
}
