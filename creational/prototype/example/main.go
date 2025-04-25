package main

import (
	"fmt"
	"time"

	"github.com/edgardnogueira/go-patterns/creational/prototype"
)

func main() {
	fmt.Println("Prototype Pattern Example - Document Generation System")
	fmt.Println("=====================================================")

	// Create a document registry
	registry := prototype.NewDocumentRegistry()

	// Initialize with document templates (prototypes)
	initializeRegistry(registry)

	fmt.Println("\nAvailable document templates:")
	for _, name := range registry.List() {
		fmt.Printf("- %s\n", name)
	}

	// Use case 1: Clone a report document and customize it
	fmt.Println("\nUse Case 1: Creating a custom quarterly report")
	reportDoc, err := registry.DeepClone("quarterly-report")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Type assertion to access specific properties
	report, ok := reportDoc.(*prototype.ReportDocument)
	if !ok {
		fmt.Printf("Error: Expected ReportDocument, got %T\n", reportDoc)
		return
	}

	// Customize the report
	report.Document.Name = "Q1 2025 Financial Report"
	report.Title = "Q1 2025 Financial Performance"
	report.ExecutiveSummary = "Strong financial performance in Q1 2025 with revenue growth of 15%."
	report.Sections[0].Content = "The company achieved significant revenue growth in Q1 2025..."
	report.Sections[0].Charts[0].Data["Q1 2025"] = 1250000

	// Print report details
	fmt.Printf("Created report: %s (%s)\n", report.Document.Name, report.Document.ID)
	fmt.Printf("Title: %s\n", report.Title)
	fmt.Printf("Sections: %d\n", len(report.Sections))
	fmt.Printf("Charts in first section: %d\n", len(report.Sections[0].Charts))

	// Use case 2: Clone a contract document and customize it
	fmt.Println("\nUse Case 2: Creating a new service agreement")
	contractDoc, err := registry.DeepClone("service-agreement")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	contract, ok := contractDoc.(*prototype.ContractDocument)
	if !ok {
		fmt.Printf("Error: Expected ContractDocument, got %T\n", contractDoc)
		return
	}

	// Customize the contract
	contract.Document.Name = "Cloud Services Agreement - Client XYZ"
	contract.Document.ID = fmt.Sprintf("CONTRACT-%d", time.Now().Unix())
	contract.Title = "Cloud Services Agreement - Client XYZ"
	contract.EffectiveDate = time.Now().Format(time.RFC3339)
	contract.ExpirationDate = time.Now().AddDate(1, 0, 0).Format(time.RFC3339) // 1 year later
	contract.Parties[1].Name = "XYZ Corporation"
	contract.Parties[1].Contact = "ceo@xyzcorp.com"

	// Print contract details
	fmt.Printf("Created contract: %s (%s)\n", contract.Document.Name, contract.Document.ID)
	fmt.Printf("Between: %s and %s\n", contract.Parties[0].Name, contract.Parties[1].Name)
	fmt.Printf("Effective: %s\n", contract.EffectiveDate)
	fmt.Printf("Expires: %s\n", contract.ExpirationDate)
	fmt.Printf("Clauses: %d\n", len(contract.Clauses))

	// Use case 3: Clone an invoice document
	fmt.Println("\nUse Case 3: Generating a recurring invoice")
	invoiceDoc, err := registry.DeepClone("monthly-invoice")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	invoice, ok := invoiceDoc.(*prototype.InvoiceDocument)
	if !ok {
		fmt.Printf("Error: Expected InvoiceDocument, got %T\n", invoiceDoc)
		return
	}

	// Customize the invoice
	currentMonth := time.Now().Format("Jan-2006")
	invoice.Document.Name = fmt.Sprintf("Monthly Invoice - %s", currentMonth)
	invoice.InvoiceNumber = fmt.Sprintf("INV-%d", time.Now().Unix())
	invoice.CustomerInfo.Name = "Acme Corp"
	invoice.DueDate = time.Now().AddDate(0, 0, 30).Format(time.RFC3339) // Due in 30 days
	
	// Update line items
	invoice.LineItems[0].Description = fmt.Sprintf("Cloud Services - %s", currentMonth)
	invoice.LineItems[0].Amount = 1299.99
	invoice.Recalculate() // Custom method we'll imagine exists to update totals

	// Print invoice details
	fmt.Printf("Created invoice: %s (%s)\n", invoice.Document.Name, invoice.Document.ID)
	fmt.Printf("Invoice Number: %s\n", invoice.InvoiceNumber)
	fmt.Printf("Customer: %s\n", invoice.CustomerInfo.Name)
	fmt.Printf("Amount: $%.2f\n", invoice.Total)
	fmt.Printf("Due Date: %s\n", invoice.DueDate)

	// Use case 4: Demonstrate shallow vs deep cloning
	fmt.Println("\nUse Case 4: Shallow vs Deep Cloning")
	
	// Get original form
	formTemplateObj, _ := registry.Get("feedback-form")
	formTemplate, _ := formTemplateObj.(*prototype.FormDocument)
	
	// Shallow clone
	shallowCloneObj := formTemplate.Clone()
	shallowClone, _ := shallowCloneObj.(*prototype.FormDocument)
	
	// Deep clone
	deepCloneObj := formTemplate.DeepClone()
	deepClone, _ := deepCloneObj.(*prototype.FormDocument)
	
	// Modify the original
	fmt.Println("Modifying the original form's fields...")
	formTemplate.Fields[0].Label = "MODIFIED LABEL"
	
	// Compare results
	fmt.Printf("Original label: %s\n", formTemplate.Fields[0].Label)
	fmt.Printf("Shallow clone label: %s\n", shallowClone.Fields[0].Label)
	fmt.Printf("Deep clone label: %s\n", deepClone.Fields[0].Label)
	
	fmt.Println("\nDemonstration complete!")
}

// initializeRegistry initializes the registry with document templates
func initializeRegistry(registry *prototype.DocumentRegistry) {
	// Create a report document template
	reportTemplate := &prototype.ReportDocument{
		Document: prototype.Document{
			ID:      "REPORT-TEMPLATE-001",
			Name:    "Quarterly Report Template",
			Creator: "Document System",
			Created: time.Now().Format(time.RFC3339),
			Tags:    []string{"report", "financial", "template"},
			Metadata: map[string]string{
				"category": "financial",
				"version":  "1.0",
			},
		},
		Title:           "Quarterly Financial Report",
		ExecutiveSummary: "Quarterly financial performance summary",
		Introduction:    "This report presents the financial results for the quarter.",
		Sections: []prototype.ReportSection{
			{
				Title:   "Financial Results",
				Content: "Analysis of the financial performance for the quarter.",
				Charts: []prototype.Chart{
					{
						Type:  "bar",
						Title: "Quarterly Revenue",
						Data: map[string]float64{
							"Previous Quarter": 1000000,
							"Current Quarter":  0, // To be filled in when customized
						},
					},
				},
			},
			{
				Title:   "Market Analysis",
				Content: "Analysis of market trends and competitive landscape.",
				Charts:  []prototype.Chart{},
			},
		},
		Conclusion: "Overall assessment and future outlook.",
		References: []string{"Annual Financial Report", "Market Research Data"},
	}

	// Create a contract document template
	contractTemplate := &prototype.ContractDocument{
		Document: prototype.Document{
			ID:      "CONTRACT-TEMPLATE-001",
			Name:    "Service Agreement Template",
			Creator: "Legal Department",
			Created: time.Now().Format(time.RFC3339),
			Tags:    []string{"contract", "legal", "template"},
			Metadata: map[string]string{
				"category": "legal",
				"version":  "2.1",
			},
		},
		Title: "Service Agreement",
		Parties: []prototype.Party{
			{
				Name:    "Our Company Inc.",
				Type:    "provider",
				Address: "123 Business St, Business City",
				Contact: "legal@ourcompany.com",
				Details: map[string]string{
					"registration": "REG12345",
				},
			},
			{
				Name:    "CLIENT_NAME", // Placeholder to be replaced
				Type:    "client",
				Address: "CLIENT_ADDRESS", // Placeholder to be replaced
				Contact: "CLIENT_CONTACT", // Placeholder to be replaced
				Details: map[string]string{},
			},
		},
		Clauses: []prototype.Clause{
			{
				Title:     "1. Services",
				Content:   "The Provider agrees to provide the following services to the Client...",
				IsRequired: true,
			},
			{
				Title:     "2. Term",
				Content:   "This Agreement shall commence on the Effective Date and continue until...",
				IsRequired: true,
			},
			{
				Title:     "3. Payment",
				Content:   "Client agrees to pay for the services according to the following terms...",
				IsRequired: true,
			},
		},
		EffectiveDate:  "EFFECTIVE_DATE", // Placeholder to be replaced
		ExpirationDate: "EXPIRATION_DATE", // Placeholder to be replaced
		IsExecuted:     false,
	}

	// Create an invoice document template
	invoiceTemplate := &prototype.InvoiceDocument{
		Document: prototype.Document{
			ID:      "INVOICE-TEMPLATE-001",
			Name:    "Monthly Invoice Template",
			Creator: "Finance Department",
			Created: time.Now().Format(time.RFC3339),
			Tags:    []string{"invoice", "financial", "template"},
			Metadata: map[string]string{
				"category": "financial",
				"version":  "1.0",
			},
		},
		InvoiceNumber: "INV-PLACEHOLDER", // To be generated when cloned
		CustomerInfo: prototype.Customer{
			Name:    "CUSTOMER_NAME", // Placeholder to be replaced
			Address: "CUSTOMER_ADDRESS", // Placeholder to be replaced
			Email:   "CUSTOMER_EMAIL", // Placeholder to be replaced
			Phone:   "CUSTOMER_PHONE", // Placeholder to be replaced
			ID:      "CUSTOMER_ID", // Placeholder to be replaced
		},
		LineItems: []prototype.LineItem{
			{
				Description: "Service Subscription",
				Quantity:    1,
				UnitPrice:   999.99,
				Amount:      999.99,
				TaxCategory: "standard",
			},
			{
				Description: "Support Package",
				Quantity:    1,
				UnitPrice:   299.99,
				Amount:      299.99,
				TaxCategory: "standard",
			},
		},
		Subtotal:     1299.98,
		TaxRate:      0.10, // 10%
		TaxAmount:    129.998,
		Total:        1429.978,
		DueDate:      "DUE_DATE", // Placeholder to be replaced
		IsPaid:       false,
		PaymentTerms: "Net 30",
	}

	// Create a form document template
	formTemplate := &prototype.FormDocument{
		Document: prototype.Document{
			ID:      "FORM-TEMPLATE-001",
			Name:    "Customer Feedback Form Template",
			Creator: "Customer Success Team",
			Created: time.Now().Format(time.RFC3339),
			Tags:    []string{"form", "feedback", "customer", "template"},
			Metadata: map[string]string{
				"category": "customer-service",
				"version":  "1.2",
			},
		},
		Title:       "Customer Feedback Form",
		Description: "Please provide your feedback on our services",
		Fields: []prototype.FormField{
			{
				Name:        "name",
				Label:       "Your Name",
				Type:        "text",
				Placeholder: "John Doe",
				IsRequired:  true,
			},
			{
				Name:        "email",
				Label:       "Your Email",
				Type:        "email",
				Placeholder: "john.doe@example.com",
				IsRequired:  true,
				Validation:  `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
			},
			{
				Name:        "service_rating",
				Label:       "How would you rate our service?",
				Type:        "select",
				IsRequired:  true,
				Options:     []string{"Excellent", "Good", "Average", "Poor", "Very Poor"},
			},
			{
				Name:        "feedback",
				Label:       "Additional Comments",
				Type:        "textarea",
				Placeholder: "Please share your thoughts...",
				IsRequired:  false,
			},
		},
		SubmitURL:  "/api/feedback",
		IsRequired: false,
	}

	// Register all templates in the registry
	registry.Register("quarterly-report", reportTemplate)
	registry.Register("service-agreement", contractTemplate)
	registry.Register("monthly-invoice", invoiceTemplate)
	registry.Register("feedback-form", formTemplate)
}

// Method to recalculate invoice totals (just for the example)
func (i *prototype.InvoiceDocument) Recalculate() {
	i.Subtotal = 0
	for _, item := range i.LineItems {
		i.Subtotal += item.Amount
	}
	i.TaxAmount = i.Subtotal * i.TaxRate
	i.Total = i.Subtotal + i.TaxAmount
}
