package main

import (
	"fmt"
	"github.com/edgardnogueira/go-patterns/behavioral/chain"
	"time"
)

func main() {
	fmt.Println("Chain of Responsibility Pattern - Support Ticket System")
	fmt.Println("=====================================================")
	fmt.Println()

	// Create the support ticket handling chain
	supportChain := setupSupportChain()
	
	// Process different types of tickets
	fmt.Println("Processing various support tickets:")
	fmt.Println("----------------------------------")

	// Create and process tickets
	tickets := []*chain.SupportTicket{
		chain.NewSupportTicket("TKT-001", chain.General, chain.Low, 
			"Password Reset", "Need help resetting my password"),
		
		chain.NewSupportTicket("TKT-002", chain.Technical, chain.Medium, 
			"App Crash", "Application crashes when uploading files"),
		
		chain.NewSupportTicket("TKT-003", chain.Bug, chain.High, 
			"Data Loss", "Lost customer data after system update"),
		
		chain.NewSupportTicket("TKT-004", chain.Security, chain.Critical, 
			"Security Breach", "Detected unauthorized access to admin account"),
		
		chain.NewSupportTicket("TKT-005", chain.Complaint, chain.Medium, 
			"Poor Service", "Waited for 2 hours for support response"),
		
		chain.NewSupportTicket("TKT-006", chain.Technical, chain.Low, 
			"Urgent Login Problem", "URGENT: Can't access my account for important meeting"),
	}

	// Process each ticket and display results
	for _, ticket := range tickets {
		fmt.Printf("\n%s\n", ticket)
		fmt.Println("Processing ticket...")
		
		// Add a small delay to simulate processing time
		time.Sleep(500 * time.Millisecond)
		
		// Process the ticket through the chain
		supportChain.Process(ticket)
		
		// Show the result
		fmt.Printf("\nResult: %s\n", ticket)
		fmt.Printf("Resolved by: %s\n", ticket.ResolvedBy)
		fmt.Printf("Resolution: %s\n", ticket.Resolution)
		fmt.Println(strings.Repeat("-", 50))
	}

	// Demonstrate chain modification
	fmt.Println("\nDemonstrating Dynamic Chain Modification")
	fmt.Println("--------------------------------------")
	
	// Create a ticket that would normally go to Level 2 Support
	specialTicket := chain.NewSupportTicket("TKT-007", chain.Billing, chain.Medium, 
		"Billing Question", "Question about my recent invoice")
	
	// Process normally
	fmt.Println("\nProcessing billing ticket before adding special billing handler:")
	supportChain.Process(specialTicket)
	fmt.Printf("Resolved by: %s\n", specialTicket.ResolvedBy)
	
	// Create and insert a special billing handler
	fmt.Println("\nAdding specialized billing handler to the chain...")
	billingHandler := NewBillingHandler()
	supportChain.InsertHandler("Level 1 Support", billingHandler)
	
	// Create another similar ticket
	specialTicket2 := chain.NewSupportTicket("TKT-008", chain.Billing, chain.Medium, 
		"Another Billing Question", "Need help understanding my invoice")
	
	// Process with the modified chain
	fmt.Println("\nProcessing billing ticket after adding special billing handler:")
	supportChain.Process(specialTicket2)
	fmt.Printf("Resolved by: %s\n", specialTicket2.ResolvedBy)
	
	// Conclusion
	fmt.Println("\nChain of Responsibility Pattern Advantages:")
	fmt.Println("1. Decouples sender from receivers")
	fmt.Println("2. Adds flexibility in assigning responsibilities")
	fmt.Println("3. Allows dynamic modification of the chain")
	fmt.Println("4. Each handler can focus on a specific responsibility")
}

// BillingHandler is a specialized handler for billing inquiries
type BillingHandler struct {
	chain.BaseHandler
}

// NewBillingHandler creates a new billing handler
func NewBillingHandler() *BillingHandler {
	return &BillingHandler{
		chain.BaseHandler{
			name: "Billing Department",
		},
	}
}

// Handle processes billing-related tickets
func (h *BillingHandler) Handle(ticket *chain.SupportTicket) bool {
	if ticket.Type == chain.Billing {
		ticket.SetResolution(h.Name(), "Billing inquiry handled by specialized billing department")
		return true
	}
	
	// Pass to the next handler if we can't handle it
	return h.BaseHandler.Handle(ticket)
}

// setupSupportChain creates and configures the support ticket handling chain
func setupSupportChain() *chain.Chain {
	// Create handlers
	loggingHandler := chain.NewLoggingHandler(func(message string) {
		fmt.Printf("[LOG] %s\n", message)
	})
	
	priorityHandler := chain.NewPriorityUpgradeHandler()
	level1 := chain.NewLevel1Support()
	level2 := chain.NewLevel2Support()
	level3 := chain.NewLevel3Support()
	securityTeam := chain.NewSecurityHandler()
	managerSupport := chain.NewManagerSupport()
	fallbackHandler := chain.NewFallbackHandler()
	
	// Configure chain
	supportChain := chain.NewChain(loggingHandler)
	supportChain.AddHandler(priorityHandler)
	supportChain.AddHandler(level1)
	supportChain.AddHandler(level2)
	supportChain.AddHandler(level3)
	supportChain.AddHandler(securityTeam)
	supportChain.AddHandler(managerSupport)
	supportChain.AddHandler(fallbackHandler)
	
	return supportChain
}
