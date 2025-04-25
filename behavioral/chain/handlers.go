package chain

import (
	"fmt"
	"strings"
)

// Level1Support handles basic user inquiries and common issues
type Level1Support struct {
	BaseHandler
}

// NewLevel1Support creates a new Level1Support handler
func NewLevel1Support() *Level1Support {
	return &Level1Support{
		BaseHandler: BaseHandler{
			name: "Level 1 Support",
		},
	}
}

// Handle processes support tickets that can be handled by Level 1 Support
func (h *Level1Support) Handle(ticket *SupportTicket) bool {
	// Level 1 can handle general inquiries and low priority issues
	if ticket.Type == General && ticket.Priority <= Low {
		ticket.SetResolution(h.Name(), "Resolved basic user inquiry")
		return true
	}
	
	// Handle simple technical issues
	if ticket.Type == Technical && ticket.Priority == Low {
		// Check if it's a simple issue by keywords
		description := strings.ToLower(ticket.Description)
		if strings.Contains(description, "password reset") || 
		   strings.Contains(description, "login issue") || 
		   strings.Contains(description, "account access") {
			ticket.SetResolution(h.Name(), "Provided instructions for basic account access")
			return true
		}
	}
	
	// Handle simple feature requests by redirecting to documentation
	if ticket.Type == FeatureRequest && ticket.Priority == Low {
		// Check if it's already available
		description := strings.ToLower(ticket.Description)
		if strings.Contains(description, "how to") || 
		   strings.Contains(description, "where is") || 
		   strings.Contains(description, "can i") {
			ticket.SetResolution(h.Name(), "Provided links to relevant documentation")
			return true
		}
	}
	
	// Pass to the next handler if we can't handle it
	return h.BaseHandler.Handle(ticket)
}

// Level2Support handles technical issues that require more expertise
type Level2Support struct {
	BaseHandler
}

// NewLevel2Support creates a new Level2Support handler
func NewLevel2Support() *Level2Support {
	return &Level2Support{
		BaseHandler: BaseHandler{
			name: "Level 2 Support",
		},
	}
}

// Handle processes support tickets that can be handled by Level 2 Support
func (h *Level2Support) Handle(ticket *SupportTicket) bool {
	// Level 2 can handle technical issues up to medium priority
	if ticket.Type == Technical && ticket.Priority <= Medium {
		ticket.SetResolution(h.Name(), "Resolved technical issue after troubleshooting")
		return true
	}
	
	// Handle low and medium priority bugs
	if ticket.Type == Bug && ticket.Priority <= Medium {
		ticket.SetResolution(h.Name(), "Identified and fixed software issue")
		return true
	}
	
	// Handle billing issues
	if ticket.Type == Billing && ticket.Priority <= Medium {
		ticket.SetResolution(h.Name(), "Resolved billing inquiry and updated account")
		return true
	}
	
	// Handle feature requests of medium priority
	if ticket.Type == FeatureRequest && ticket.Priority == Medium {
		ticket.SetResolution(h.Name(), "Logged feature request for consideration")
		return true
	}
	
	// Pass to the next handler if we can't handle it
	return h.BaseHandler.Handle(ticket)
}

// Level3Support handles complex issues that require system-level access
type Level3Support struct {
	BaseHandler
}

// NewLevel3Support creates a new Level3Support handler
func NewLevel3Support() *Level3Support {
	return &Level3Support{
		BaseHandler: BaseHandler{
			name: "Level 3 Support",
		},
	}
}

// Handle processes support tickets that require Level 3 expertise
func (h *Level3Support) Handle(ticket *SupportTicket) bool {
	// Level 3 can handle high priority technical issues and bugs
	if (ticket.Type == Technical || ticket.Type == Bug) && ticket.Priority == High {
		ticket.SetResolution(h.Name(), "Resolved complex technical issue requiring system access")
		return true
	}
	
	// Handle high priority billing issues
	if ticket.Type == Billing && ticket.Priority == High {
		ticket.SetResolution(h.Name(), "Resolved complex billing issue")
		return true
	}
	
	// Handle security issues that aren't critical
	if ticket.Type == Security && ticket.Priority < Critical {
		ticket.SetResolution(h.Name(), "Addressed security concern after investigation")
		return true
	}
	
	// Handle high priority feature requests
	if ticket.Type == FeatureRequest && ticket.Priority == High {
		ticket.SetResolution(h.Name(), "Evaluated feature request and forwarded to development team")
		return true
	}
	
	// Pass to the next handler if we can't handle it
	return h.BaseHandler.Handle(ticket)
}

// ManagerSupport handles escalated issues, customer complaints, and policy exceptions
type ManagerSupport struct {
	BaseHandler
}

// NewManagerSupport creates a new ManagerSupport handler
func NewManagerSupport() *ManagerSupport {
	return &ManagerSupport{
		BaseHandler: BaseHandler{
			name: "Manager Support",
		},
	}
}

// Handle processes escalated and critical tickets
func (h *ManagerSupport) Handle(ticket *SupportTicket) bool {
	// Managers handle all customer complaints
	if ticket.Type == Complaint {
		ticket.SetResolution(h.Name(), "Customer complaint addressed by management")
		return true
	}
	
	// Managers handle all critical issues
	if ticket.Priority == Critical {
		resolution := fmt.Sprintf("Critical %s issue resolved with emergency response", ticket.Type)
		ticket.SetResolution(h.Name(), resolution)
		return true
	}
	
	// Handle critical security issues as highest priority
	if ticket.Type == Security {
		ticket.SetResolution(h.Name(), "Security issue escalated and resolved with security team")
		return true
	}
	
	// Pass to the next handler if we can't handle it
	return h.BaseHandler.Handle(ticket)
}

// SecurityHandler handles security-related issues with specialized expertise
type SecurityHandler struct {
	BaseHandler
}

// NewSecurityHandler creates a new SecurityHandler
func NewSecurityHandler() *SecurityHandler {
	return &SecurityHandler{
		BaseHandler: BaseHandler{
			name: "Security Team",
		},
	}
}

// Handle processes security-specific tickets
func (h *SecurityHandler) Handle(ticket *SupportTicket) bool {
	// Security team handles all security issues
	if ticket.Type == Security {
		// Craft resolution based on priority
		var resolution string
		switch ticket.Priority {
		case Critical:
			resolution = "Critical security vulnerability patched and verified"
		case High:
			resolution = "Security issue addressed and remediation confirmed"
		case Medium:
			resolution = "Security concern investigated and mitigated"
		default:
			resolution = "Security inquiry reviewed and addressed"
		}
		
		ticket.SetResolution(h.Name(), resolution)
		return true
	}
	
	// Pass to the next handler if it's not a security issue
	return h.BaseHandler.Handle(ticket)
}

// FallbackHandler ensures all tickets receive a response even if not fully resolved
type FallbackHandler struct {
	BaseHandler
}

// NewFallbackHandler creates a new FallbackHandler
func NewFallbackHandler() *FallbackHandler {
	return &FallbackHandler{
		BaseHandler: BaseHandler{
			name: "Fallback Handler",
		},
	}
}

// Handle provides a fallback response for any unhandled tickets
func (h *FallbackHandler) Handle(ticket *SupportTicket) bool {
	if !ticket.IsResolved {
		resolution := fmt.Sprintf("Ticket escalated for specialized review. We'll get back to you regarding this %s priority %s issue.", 
			ticket.Priority, ticket.Type)
		ticket.SetResolution(h.Name(), resolution)
		
		// Add metadata to indicate this was a fallback
		ticket.Metadata["requires_followup"] = true
		ticket.Metadata["escalated"] = true
		
		return true
	}
	
	return h.BaseHandler.Handle(ticket)
}

// LoggingHandler is a special handler that logs all tickets passing through but doesn't resolve them
type LoggingHandler struct {
	BaseHandler
	logFunc func(string)
}

// NewLoggingHandler creates a new LoggingHandler with the provided logging function
func NewLoggingHandler(logFunc func(string)) *LoggingHandler {
	return &LoggingHandler{
		BaseHandler: BaseHandler{
			name: "Logging Handler",
		},
		logFunc: logFunc,
	}
}

// Handle logs the ticket and passes it to the next handler
func (h *LoggingHandler) Handle(ticket *SupportTicket) bool {
	if h.logFunc != nil {
		logMessage := fmt.Sprintf("Processing ticket #%s: [%s, %s] - %s", 
			ticket.ID, ticket.Type, ticket.Priority, ticket.Subject)
		h.logFunc(logMessage)
	}
	
	// Always pass to the next handler
	return h.BaseHandler.Handle(ticket)
}

// PriorityUpgradeHandler upgrades ticket priority based on keywords
type PriorityUpgradeHandler struct {
	BaseHandler
	keywords map[string]Priority
}

// NewPriorityUpgradeHandler creates a new handler that can upgrade ticket priorities
func NewPriorityUpgradeHandler() *PriorityUpgradeHandler {
	handler := &PriorityUpgradeHandler{
		BaseHandler: BaseHandler{
			name: "Priority Upgrade Handler",
		},
		keywords: make(map[string]Priority),
	}
	
	// Add default critical keywords
	handler.AddKeyword("urgent", Critical)
	handler.AddKeyword("emergency", Critical)
	handler.AddKeyword("critical", Critical)
	handler.AddKeyword("broken", High)
	handler.AddKeyword("error", Medium)
	
	return handler
}

// AddKeyword adds or updates a keyword that will trigger priority upgrade
func (h *PriorityUpgradeHandler) AddKeyword(keyword string, priority Priority) {
	h.keywords[strings.ToLower(keyword)] = priority
}

// Handle checks for priority keywords and upgrades if needed
func (h *PriorityUpgradeHandler) Handle(ticket *SupportTicket) bool {
	originalPriority := ticket.Priority
	
	// Check description for priority keywords
	description := strings.ToLower(ticket.Description)
	subject := strings.ToLower(ticket.Subject)
	
	for keyword, priority := range h.keywords {
		if (strings.Contains(description, keyword) || strings.Contains(subject, keyword)) && 
		   priority > ticket.Priority {
			ticket.Priority = priority
			
			// Add metadata about the upgrade
			if ticket.Metadata == nil {
				ticket.Metadata = make(map[string]interface{})
			}
			ticket.Metadata["priority_upgraded"] = true
			ticket.Metadata["original_priority"] = originalPriority
		}
	}
	
	// Always pass to the next handler
	return h.BaseHandler.Handle(ticket)
}
