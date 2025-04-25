package prototype

// ReportDocument represents a business report with standardized sections.
type ReportDocument struct {
	Document
	Title           string
	ExecutiveSummary string
	Introduction    string
	Sections        []ReportSection
	Conclusion      string
	References      []string
}

// ReportSection represents a section in a report.
type ReportSection struct {
	Title   string
	Content string
	Charts  []Chart
}

// Chart represents a chart in a report section.
type Chart struct {
	Type  string // "bar", "line", "pie", etc.
	Title string
	Data  map[string]float64
}

// Clone creates a shallow copy of the ReportDocument.
func (r *ReportDocument) Clone() Prototype {
	clonedDoc := r.Document.Clone().(*Document)
	
	cloned := &ReportDocument{
		Document:        *clonedDoc,
		Title:           r.Title,
		ExecutiveSummary: r.ExecutiveSummary,
		Introduction:    r.Introduction,
		Sections:        r.Sections,  // Shallow copy
		Conclusion:      r.Conclusion,
		References:      r.References, // Shallow copy
	}
	
	return cloned
}

// DeepClone creates a deep copy of the ReportDocument.
func (r *ReportDocument) DeepClone() Prototype {
	clonedDoc := r.Document.DeepClone().(*Document)
	
	cloned := &ReportDocument{
		Document:        *clonedDoc,
		Title:           r.Title,
		ExecutiveSummary: r.ExecutiveSummary,
		Introduction:    r.Introduction,
		Conclusion:      r.Conclusion,
		// Deep copy of sections
		Sections:        make([]ReportSection, len(r.Sections)),
		// Deep copy of references
		References:      make([]string, len(r.References)),
	}
	
	// Copy sections
	for i, section := range r.Sections {
		cloned.Sections[i] = ReportSection{
			Title:   section.Title,
			Content: section.Content,
			Charts:  make([]Chart, len(section.Charts)),
		}
		
		// Copy charts
		for j, chart := range section.Charts {
			cloned.Sections[i].Charts[j] = Chart{
				Type:  chart.Type,
				Title: chart.Title,
				Data:  make(map[string]float64, len(chart.Data)),
			}
			
			// Copy chart data
			for k, v := range chart.Data {
				cloned.Sections[i].Charts[j].Data[k] = v
			}
		}
	}
	
	// Copy references
	copy(cloned.References, r.References)
	
	return cloned
}

// FormDocument represents a document with fillable form fields.
type FormDocument struct {
	Document
	Title       string
	Description string
	Fields      []FormField
	SubmitURL   string
	IsRequired  bool
}

// FormField represents a field in a form.
type FormField struct {
	Name        string
	Label       string
	Type        string  // "text", "select", "checkbox", etc.
	Value       string
	Placeholder string
	IsRequired  bool
	Validation  string  // Validation rules (regex, etc.)
	Options     []string // For select, radio, etc.
}

// Clone creates a shallow copy of the FormDocument.
func (f *FormDocument) Clone() Prototype {
	clonedDoc := f.Document.Clone().(*Document)
	
	cloned := &FormDocument{
		Document:    *clonedDoc,
		Title:       f.Title,
		Description: f.Description,
		Fields:      f.Fields,  // Shallow copy
		SubmitURL:   f.SubmitURL,
		IsRequired:  f.IsRequired,
	}
	
	return cloned
}

// DeepClone creates a deep copy of the FormDocument.
func (f *FormDocument) DeepClone() Prototype {
	clonedDoc := f.Document.DeepClone().(*Document)
	
	cloned := &FormDocument{
		Document:    *clonedDoc,
		Title:       f.Title,
		Description: f.Description,
		// Deep copy of fields
		Fields:      make([]FormField, len(f.Fields)),
		SubmitURL:   f.SubmitURL,
		IsRequired:  f.IsRequired,
	}
	
	// Copy fields
	for i, field := range f.Fields {
		cloned.Fields[i] = FormField{
			Name:        field.Name,
			Label:       field.Label,
			Type:        field.Type,
			Value:       field.Value,
			Placeholder: field.Placeholder,
			IsRequired:  field.IsRequired,
			Validation:  field.Validation,
			Options:     make([]string, len(field.Options)),
		}
		
		// Copy options
		copy(cloned.Fields[i].Options, field.Options)
	}
	
	return cloned
}

// ContractDocument represents a legal contract with clauses and terms.
type ContractDocument struct {
	Document
	Title          string
	Parties        []Party
	Clauses        []Clause
	EffectiveDate  string // ISO 8601 date format
	ExpirationDate string // ISO 8601 date format
	IsExecuted     bool
	Signatures     []Signature
}

// Party represents a party in a contract.
type Party struct {
	Name      string
	Type      string // "individual", "company", etc.
	Address   string
	Contact   string
	Details   map[string]string
}

// Clause represents a clause in a contract.
type Clause struct {
	Title   string
	Content string
	IsRequired bool
}

// Signature represents a signature on a contract.
type Signature struct {
	PartyName  string
	SignedBy   string
	SignedDate string // ISO 8601 date format
	IPAddress  string
}

// Clone creates a shallow copy of the ContractDocument.
func (c *ContractDocument) Clone() Prototype {
	clonedDoc := c.Document.Clone().(*Document)
	
	cloned := &ContractDocument{
		Document:       *clonedDoc,
		Title:          c.Title,
		Parties:        c.Parties,  // Shallow copy
		Clauses:        c.Clauses,  // Shallow copy
		EffectiveDate:  c.EffectiveDate,
		ExpirationDate: c.ExpirationDate,
		IsExecuted:     false, // Reset execution status
		Signatures:     nil,   // Reset signatures
	}
	
	return cloned
}

// DeepClone creates a deep copy of the ContractDocument.
func (c *ContractDocument) DeepClone() Prototype {
	clonedDoc := c.Document.DeepClone().(*Document)
	
	cloned := &ContractDocument{
		Document:       *clonedDoc,
		Title:          c.Title,
		// Deep copy of parties
		Parties:        make([]Party, len(c.Parties)),
		// Deep copy of clauses
		Clauses:        make([]Clause, len(c.Clauses)),
		EffectiveDate:  c.EffectiveDate,
		ExpirationDate: c.ExpirationDate,
		IsExecuted:     false, // Reset execution status
		Signatures:     nil,   // Reset signatures
	}
	
	// Copy parties
	for i, party := range c.Parties {
		cloned.Parties[i] = Party{
			Name:    party.Name,
			Type:    party.Type,
			Address: party.Address,
			Contact: party.Contact,
			Details: make(map[string]string, len(party.Details)),
		}
		
		// Copy details
		for k, v := range party.Details {
			cloned.Parties[i].Details[k] = v
		}
	}
	
	// Copy clauses
	for i, clause := range c.Clauses {
		cloned.Clauses[i] = Clause{
			Title:     clause.Title,
			Content:   clause.Content,
			IsRequired: clause.IsRequired,
		}
	}
	
	return cloned
}

// InvoiceDocument represents an invoice with line items and totals.
type InvoiceDocument struct {
	Document
	InvoiceNumber string
	CustomerInfo  Customer
	LineItems     []LineItem
	Subtotal      float64
	TaxRate       float64
	TaxAmount     float64
	Total         float64
	DueDate       string // ISO 8601 date format
	IsPaid        bool
	PaymentTerms  string
}

// Customer represents a customer in an invoice.
type Customer struct {
	Name    string
	Address string
	Email   string
	Phone   string
	ID      string
}

// LineItem represents a line item in an invoice.
type LineItem struct {
	Description string
	Quantity    int
	UnitPrice   float64
	Amount      float64
	TaxCategory string
}

// Clone creates a shallow copy of the InvoiceDocument.
func (i *InvoiceDocument) Clone() Prototype {
	clonedDoc := i.Document.Clone().(*Document)
	
	cloned := &InvoiceDocument{
		Document:      *clonedDoc,
		InvoiceNumber: "NEW-" + i.InvoiceNumber, // Generate a new invoice number
		CustomerInfo:  i.CustomerInfo,  // Shallow copy
		LineItems:     i.LineItems,     // Shallow copy
		Subtotal:      i.Subtotal,
		TaxRate:       i.TaxRate,
		TaxAmount:     i.TaxAmount,
		Total:         i.Total,
		DueDate:       i.DueDate,
		IsPaid:        false, // Reset payment status
		PaymentTerms:  i.PaymentTerms,
	}
	
	return cloned
}

// DeepClone creates a deep copy of the InvoiceDocument.
func (i *InvoiceDocument) DeepClone() Prototype {
	clonedDoc := i.Document.DeepClone().(*Document)
	
	cloned := &InvoiceDocument{
		Document:      *clonedDoc,
		InvoiceNumber: "NEW-" + i.InvoiceNumber, // Generate a new invoice number
		// Deep copy of customer info
		CustomerInfo: Customer{
			Name:    i.CustomerInfo.Name,
			Address: i.CustomerInfo.Address,
			Email:   i.CustomerInfo.Email,
			Phone:   i.CustomerInfo.Phone,
			ID:      i.CustomerInfo.ID,
		},
		// Deep copy of line items
		LineItems:    make([]LineItem, len(i.LineItems)),
		Subtotal:     i.Subtotal,
		TaxRate:      i.TaxRate,
		TaxAmount:    i.TaxAmount,
		Total:        i.Total,
		DueDate:      i.DueDate,
		IsPaid:       false, // Reset payment status
		PaymentTerms: i.PaymentTerms,
	}
	
	// Copy line items
	for j, item := range i.LineItems {
		cloned.LineItems[j] = LineItem{
			Description: item.Description,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Amount:      item.Amount,
			TaxCategory: item.TaxCategory,
		}
	}
	
	return cloned
}
