package main

import (
	"fmt"
	"github.com/edgardnogueira/go-patterns/behavioral/chain"
)

func main() {
	fmt.Println("Chain of Responsibility Pattern")
	fmt.Println("==============================")
	fmt.Println("The Chain of Responsibility pattern passes a request along a chain")
	fmt.Println("of handlers. Each handler decides either to process the request or")
	fmt.Println("to pass it to the next handler in the chain.")
	fmt.Println()

	// Create a support ticket handling chain
	fmt.Println("Creating a support ticket handling chain...")
	
	// Create handlers
	level1 := chain.NewLevel1Support()
	level2 := chain.NewLevel2Support()
	level3 := chain.NewLevel3Support()
	manager := chain.NewManagerSupport()
	fallback := chain.NewFallbackHandler()
	
	// Setup the chain
	supportChain := chain.NewChain(level1)
	supportChain.AddHandler(level2)
	supportChain.AddHandler(level3)
	supportChain.AddHandler(manager)
	supportChain.AddHandler(fallback)
	
	fmt.Println("Chain setup complete!")
	fmt.Println()
	
	// Process tickets of different types and priorities
	fmt.Println("Example 1: Low priority general inquiry (should be handled by Level 1)")
	ticket1 := chain.NewSupportTicket("TKT-001", chain.General, chain.Low, 
		"Password Reset", "How do I reset my password?")
	supportChain.Process(ticket1)
	printTicketResult(ticket1)
	
	fmt.Println("Example 2: Medium priority technical issue (should be handled by Level 2)")
	ticket2 := chain.NewSupportTicket("TKT-002", chain.Technical, chain.Medium, 
		"App Crash", "The application crashes when uploading large files")
	supportChain.Process(ticket2)
	printTicketResult(ticket2)
	
	fmt.Println("Example 3: High priority bug (should be handled by Level 3)")
	ticket3 := chain.NewSupportTicket("TKT-003", chain.Bug, chain.High, 
		"Data Loss", "Customer data is being lost during transaction")
	supportChain.Process(ticket3)
	printTicketResult(ticket3)
	
	fmt.Println("Example 4: Critical security issue (should be handled by Manager)")
	ticket4 := chain.NewSupportTicket("TKT-004", chain.Security, chain.Critical, 
		"Security Breach", "Detected unauthorized access to admin accounts")
	supportChain.Process(ticket4)
	printTicketResult(ticket4)
	
	fmt.Println("Example 5: Customer complaint (should be handled by Manager)")
	ticket5 := chain.NewSupportTicket("TKT-005", chain.Complaint, chain.Medium, 
		"Poor Service", "Unhappy with the recent service quality")
	supportChain.Process(ticket5)
	printTicketResult(ticket5)
	
	// The pattern's key aspects
	fmt.Println("\nKey aspects of the Chain of Responsibility pattern:")
	fmt.Println("1. Decoupling - The sender of a request doesn't need to know which handler will process it")
	fmt.Println("2. Single Responsibility - Each handler focuses on one specific task")
	fmt.Println("3. Flexibility - Handlers can be added, removed, or reordered at runtime")
	fmt.Println("4. Open/Closed Principle - New handlers can be added without changing existing code")
	fmt.Println("5. Default behavior - Fallback handlers ensure all requests get processed in some way")
	
	fmt.Println("\nSee example/main.go for a more detailed demonstration!")
}

// printTicketResult prints the result of processing a ticket
func printTicketResult(ticket *chain.SupportTicket) {
	if ticket.IsResolved {
		fmt.Printf("Ticket #%s was resolved by: %s\n", ticket.ID, ticket.ResolvedBy)
		fmt.Printf("Resolution: %s\n", ticket.Resolution)
	} else {
		fmt.Printf("Ticket #%s could not be resolved\n", ticket.ID)
	}
	fmt.Println()
}
