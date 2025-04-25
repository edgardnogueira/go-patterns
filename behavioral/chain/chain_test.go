package chain

import (
	"strings"
	"testing"
)

// TestSupportTicketCreation tests the creation of a support ticket
func TestSupportTicketCreation(t *testing.T) {
	ticket := NewSupportTicket("TKT-123", Technical, Medium, "Login Problem", "Can't log into my account")
	
	if ticket.ID != "TKT-123" {
		t.Errorf("Expected ticket ID to be TKT-123, got %s", ticket.ID)
	}
	
	if ticket.Type != Technical {
		t.Errorf("Expected ticket type to be Technical, got %s", ticket.Type)
	}
	
	if ticket.Priority != Medium {
		t.Errorf("Expected ticket priority to be Medium, got %s", ticket.Priority)
	}
	
	if ticket.Subject != "Login Problem" {
		t.Errorf("Expected ticket subject to be 'Login Problem', got %s", ticket.Subject)
	}
	
	if ticket.Description != "Can't log into my account" {
		t.Errorf("Expected ticket description to be 'Can't log into my account', got %s", ticket.Description)
	}
	
	if ticket.IsResolved {
		t.Error("New ticket should not be resolved")
	}
	
	if ticket.CreatedAt.IsZero() {
		t.Error("Ticket creation time should be set")
	}
}

// TestTicketResolution tests setting a resolution on a ticket
func TestTicketResolution(t *testing.T) {
	ticket := NewSupportTicket("TKT-123", Technical, Medium, "Login Problem", "Can't log into my account")
	
	// Resolve the ticket
	ticket.SetResolution("Level 1 Support", "Reset user password")
	
	if !ticket.IsResolved {
		t.Error("Ticket should be marked as resolved")
	}
	
	if ticket.ResolvedBy != "Level 1 Support" {
		t.Errorf("Expected ticket to be resolved by 'Level 1 Support', got %s", ticket.ResolvedBy)
	}
	
	if ticket.Resolution != "Reset user password" {
		t.Errorf("Expected resolution to be 'Reset user password', got %s", ticket.Resolution)
	}
}

// TestSingleHandler tests a chain with a single handler
func TestSingleHandler(t *testing.T) {
	// Create a single handler
	level1 := NewLevel1Support()
	
	// Create a chain with the handler
	chain := NewChain(level1)
	
	// Create a ticket that Level 1 can handle
	ticket := NewSupportTicket("TKT-123", General, Low, "Account Question", "How do I change my password?")
	
	// Process the ticket
	result := chain.Process(ticket)
	
	if !result {
		t.Error("Chain should have processed the ticket")
	}
	
	if !ticket.IsResolved {
		t.Error("Ticket should be resolved")
	}
	
	if ticket.ResolvedBy != "Level 1 Support" {
		t.Errorf("Expected ticket to be resolved by 'Level 1 Support', got %s", ticket.ResolvedBy)
	}
}

// TestSimpleChain tests a basic chain with multiple handlers
func TestSimpleChain(t *testing.T) {
	// Create handlers
	level1 := NewLevel1Support()
	level2 := NewLevel2Support()
	level3 := NewLevel3Support()
	
	// Create a chain
	chain := NewChain(level1)
	chain.AddHandler(level2)
	chain.AddHandler(level3)
	
	// Test a ticket that Level 1 should handle
	ticket1 := NewSupportTicket("TKT-101", General, Low, "Account Question", "How do I change my password?")
	chain.Process(ticket1)
	if !ticket1.IsResolved || ticket1.ResolvedBy != "Level 1 Support" {
		t.Errorf("Expected Level 1 to handle ticket, got %s", ticket1.ResolvedBy)
	}
	
	// Test a ticket that Level 2 should handle
	ticket2 := NewSupportTicket("TKT-102", Technical, Medium, "App Error", "Getting error code 404")
	chain.Process(ticket2)
	if !ticket2.IsResolved || ticket2.ResolvedBy != "Level 2 Support" {
		t.Errorf("Expected Level 2 to handle ticket, got %s", ticket2.ResolvedBy)
	}
	
	// Test a ticket that Level 3 should handle
	ticket3 := NewSupportTicket("TKT-103", Technical, High, "Database Issue", "Database connection fails intermittently")
	chain.Process(ticket3)
	if !ticket3.IsResolved || ticket3.ResolvedBy != "Level 3 Support" {
		t.Errorf("Expected Level 3 to handle ticket, got %s", ticket3.ResolvedBy)
	}
}

// TestCompleteChain tests a complete support chain with all handler types
func TestCompleteChain(t *testing.T) {
	// Create all handlers
	priorityHandler := NewPriorityUpgradeHandler()
	level1 := NewLevel1Support()
	level2 := NewLevel2Support()
	level3 := NewLevel3Support()
	security := NewSecurityHandler()
	manager := NewManagerSupport()
	fallback := NewFallbackHandler()
	
	// Create a chain with priority handler first to upgrade tickets if needed
	chain := NewChain(priorityHandler)
	chain.AddHandler(level1)
	chain.AddHandler(level2)
	chain.AddHandler(level3)
	chain.AddHandler(security)
	chain.AddHandler(manager)
	chain.AddHandler(fallback) // Fallback handler should be last
	
	// Test cases for different ticket types
	testCases := []struct {
		ticketID      string
		ticketType    TicketType
		priority      Priority
		subject       string
		description   string
		expectedLevel string
	}{
		// Level 1 cases
		{"TKT-101", General, Low, "Account Question", "How do I change my password?", "Level 1 Support"},
		{"TKT-102", Technical, Low, "Login Help", "Can't login to my account", "Level 1 Support"},
		
		// Level 2 cases
		{"TKT-201", Technical, Medium, "App Crash", "App crashes when uploading images", "Level 2 Support"},
		{"TKT-202", Bug, Medium, "UI Bug", "Button disappears when clicked", "Level 2 Support"},
		{"TKT-203", Billing, Medium, "Invoice Question", "Why was I charged twice?", "Level 2 Support"},
		
		// Level 3 cases
		{"TKT-301", Technical, High, "API Error", "API returns 500 error consistently", "Level 3 Support"},
		{"TKT-302", Bug, High, "Data Loss", "Users reporting missing data", "Level 3 Support"},
		
		// Security cases
		{"TKT-401", Security, Low, "Security Question", "How secure is my data?", "Security Team"},
		{"TKT-402", Security, Critical, "Security Breach", "Detected unauthorized access", "Security Team"},
		
		// Manager cases
		{"TKT-501", Complaint, Medium, "Unhappy Customer", "Very dissatisfied with service", "Manager Support"},
		{"TKT-502", Technical, Critical, "Critical System Down", "Main system is completely offline", "Manager Support"},
		
		// Priority upgrading case
		{"TKT-601", Technical, Low, "Urgent Login Issue", "URGENT: System shows errors when logging in", "Level 2 Support"},
		
		// Fallback case
		{"TKT-999", TicketType("Unknown"), Low, "Weird Issue", "Something strange is happening", "Fallback Handler"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.ticketID, func(t *testing.T) {
			ticket := NewSupportTicket(tc.ticketID, tc.ticketType, tc.priority, tc.subject, tc.description)
			chain.Process(ticket)
			
			if !ticket.IsResolved {
				t.Errorf("Ticket %s should be resolved", tc.ticketID)
			}
			
			if ticket.ResolvedBy != tc.expectedLevel {
				t.Errorf("Expected ticket %s to be resolved by %s, got %s", 
					tc.ticketID, tc.expectedLevel, ticket.ResolvedBy)
			}
		})
	}
}

// TestChainModification tests adding and removing handlers from the chain
func TestChainModification(t *testing.T) {
	// Create handlers
	level1 := NewLevel1Support()
	level2 := NewLevel2Support()
	level3 := NewLevel3Support()
	
	// Create a chain with level1
	chain := NewChain(level1)
	
	// Add level3, skipping level2
	chain.AddHandler(level3)
	
	// Try a ticket that level2 would normally handle
	ticket := NewSupportTicket("TKT-201", Technical, Medium, "App Crash", "App crashes when uploading images")
	chain.Process(ticket)
	
	// It should go to level3 since level2 is not in the chain
	if !ticket.IsResolved || ticket.ResolvedBy != "Level 3 Support" {
		t.Errorf("Expected Level 3 to handle ticket, got %s", ticket.ResolvedBy)
	}
	
	// Now insert level2 between level1 and level3
	inserted := chain.InsertHandler("Level 1 Support", level2)
	if !inserted {
		t.Error("Failed to insert handler")
	}
	
	// Try a new ticket that level2 should handle
	ticket2 := NewSupportTicket("TKT-202", Technical, Medium, "Another App Crash", "Different app crash")
	chain.Process(ticket2)
	
	// Now it should go to level2
	if !ticket2.IsResolved || ticket2.ResolvedBy != "Level 2 Support" {
		t.Errorf("Expected Level 2 to handle ticket, got %s", ticket2.ResolvedBy)
	}
	
	// Remove level2
	removed := chain.RemoveHandler("Level 2 Support")
	if !removed {
		t.Error("Failed to remove handler")
	}
	
	// Try another level2 ticket
	ticket3 := NewSupportTicket("TKT-203", Technical, Medium, "Yet Another App Crash", "Third app crash")
	chain.Process(ticket3)
	
	// It should go back to level3
	if !ticket3.IsResolved || ticket3.ResolvedBy != "Level 3 Support" {
		t.Errorf("Expected Level 3 to handle ticket after removing Level 2, got %s", ticket3.ResolvedBy)
	}
}

// TestLoggingHandler tests that the logging handler properly logs tickets
func TestLoggingHandler(t *testing.T) {
	// Create a log to capture messages
	var logMessages []string
	logger := func(message string) {
		logMessages = append(logMessages, message)
	}
	
	// Create handlers
	loggingHandler := NewLoggingHandler(logger)
	level1 := NewLevel1Support()
	
	// Create a chain with logging first
	chain := NewChain(loggingHandler)
	chain.AddHandler(level1)
	
	// Process a ticket
	ticket := NewSupportTicket("TKT-101", General, Low, "Simple Question", "Need help")
	chain.Process(ticket)
	
	// Check that a log message was created
	if len(logMessages) != 1 {
		t.Errorf("Expected 1 log message, got %d", len(logMessages))
	}
	
	// Check the content of the log message
	if !strings.Contains(logMessages[0], "TKT-101") {
		t.Errorf("Log message should contain the ticket ID, got: %s", logMessages[0])
	}
	
	// Check that the ticket was still processed
	if !ticket.IsResolved {
		t.Error("Ticket should still be resolved after logging")
	}
}

// TestPriorityUpgradeHandler tests that the priority upgrade handler works correctly
func TestPriorityUpgradeHandler(t *testing.T) {
	// Create handlers
	priorityHandler := NewPriorityUpgradeHandler()
	level2 := NewLevel2Support() // Should handle medium priority tech issues
	
	// Create a chain
	chain := NewChain(priorityHandler)
	chain.AddHandler(level2)
	
	// Process a ticket with urgent in the subject
	ticket1 := NewSupportTicket("TKT-101", Technical, Low, "Urgent Help Needed", "Having a problem")
	chain.Process(ticket1)
	
	// Check that priority was upgraded and ticket resolved by Level 2
	if ticket1.Priority != Critical {
		t.Errorf("Expected priority to be upgraded to Critical, got %s", ticket1.Priority)
	}
	
	if !ticket1.IsResolved || ticket1.ResolvedBy != "Level 2 Support" {
		t.Errorf("Expected Level 2 to handle upgraded ticket, got %s", ticket1.ResolvedBy)
	}
	
	// Check metadata
	if upgraded, ok := ticket1.Metadata["priority_upgraded"].(bool); !ok || !upgraded {
		t.Error("Expected priority_upgraded metadata to be set to true")
	}
	
	if originalPriority, ok := ticket1.Metadata["original_priority"].(Priority); !ok || originalPriority != Low {
		t.Error("Expected original_priority metadata to be set to Low")
	}
	
	// Add a custom keyword
	priorityHandler.AddKeyword("not working", High)
	
	// Test with custom keyword
	ticket2 := NewSupportTicket("TKT-102", Technical, Low, "Help", "System is not working properly")
	chain.Process(ticket2)
	
	// Check that priority was upgraded to High
	if ticket2.Priority != High {
		t.Errorf("Expected priority to be upgraded to High, got %s", ticket2.Priority)
	}
}

// TestFallbackHandler tests that the fallback handler processes unresolved tickets
func TestFallbackHandler(t *testing.T) {
	// Create a fallback handler
	fallback := NewFallbackHandler()
	
	// Create a chain with just the fallback
	chain := NewChain(fallback)
	
	// Create a ticket that no previous handler would resolve
	ticket := NewSupportTicket("TKT-101", TicketType("Unknown"), Low, "Strange Issue", "Something weird")
	
	// Process the ticket
	chain.Process(ticket)
	
	// Check that fallback resolved it
	if !ticket.IsResolved {
		t.Error("Fallback handler should have resolved the ticket")
	}
	
	if ticket.ResolvedBy != "Fallback Handler" {
		t.Errorf("Expected ticket to be resolved by Fallback Handler, got %s", ticket.ResolvedBy)
	}
	
	// Check metadata
	if requiresFollowup, ok := ticket.Metadata["requires_followup"].(bool); !ok || !requiresFollowup {
		t.Error("Expected requires_followup metadata to be set to true")
	}
}

// TestEmptyChain tests behavior with an empty chain
func TestEmptyChain(t *testing.T) {
	// Create an empty chain
	var chain Chain
	
	// Try to process a ticket
	ticket := NewSupportTicket("TKT-101", General, Low, "Test", "Test")
	result := chain.Process(ticket)
	
	// Should return false and leave ticket unresolved
	if result {
		t.Error("Empty chain should return false for Process")
	}
	
	if ticket.IsResolved {
		t.Error("Ticket should not be resolved by empty chain")
	}
}

// TestInvalidChainModification tests adding/removing handlers in invalid ways
func TestInvalidChainModification(t *testing.T) {
	// Create an empty chain
	var chain Chain
	
	// Try to remove from empty chain
	if chain.RemoveHandler("nonexistent") {
		t.Error("Removing from empty chain should return false")
	}
	
	// Try to insert into empty chain with non-nil handler
	handler := NewLevel1Support()
	if chain.InsertHandler("nonexistent", handler) {
		t.Error("Inserting after nonexistent handler should return false")
	}
	
	// Add a handler
	chain.AddHandler(handler)
	
	// Try to insert after nonexistent handler
	handler2 := NewLevel2Support()
	if chain.InsertHandler("nonexistent", handler2) {
		t.Error("Inserting after nonexistent handler should return false")
	}
}
