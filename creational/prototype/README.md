# Prototype Pattern

## Intent
The Prototype pattern is a creational design pattern that allows cloning objects, even complex ones, without coupling to their specific classes. It creates objects by copying existing ones, known as prototypes.

## Problem Solved
The Prototype pattern solves the following problems:
- Creating new objects by copying an existing one is sometimes more efficient than creating new ones from scratch
- When classes to be instantiated are specified at runtime
- When creating objects that are similar to existing ones but with slight modifications
- When the creation of objects should be independent of the system they're part of
- When object creation is expensive or complex

## Structure
![Prototype Pattern Structure](https://github.com/edgardnogueira/go-patterns/raw/main/assets/images/prototype.png)

### Components
1. **Prototype**: An interface that declares the cloning methods (Clone/DeepClone)
2. **Concrete Prototype**: Classes that implement the Prototype interface
3. **Client**: Creates new objects by asking a prototype to clone itself
4. **Registry (optional)**: A catalog of available prototypes that can be cloned

## Implementation in Go

### Prototype Interface
```go
type Prototype interface {
    Clone() Prototype      // For shallow copying
    DeepClone() Prototype  // For deep copying
}
```

### Document Abstract Base (Concrete Prototype)
```go
type Document struct {
    ID       string
    Name     string
    Creator  string
    Created  string
    Modified string
    Tags     []string
    Metadata map[string]string
}

// Clone creates a shallow copy
func (d *Document) Clone() Prototype {
    // Implementation...
}

// DeepClone creates a deep copy
func (d *Document) DeepClone() Prototype {
    // Implementation...
}
```

### Concrete Document Types
We've implemented several document types as concrete prototypes:
- `ReportDocument`: For business reports with standardized sections
- `FormDocument`: For form templates with user-fillable fields
- `ContractDocument`: For legal contracts with clauses and terms
- `InvoiceDocument`: For invoice generation with line items and totals

### Document Registry
We've also implemented a registry to store and manage document prototypes:
```go
type DocumentRegistry struct {
    prototypes map[string]Prototype
    mutex      sync.RWMutex
}
```

## When to Use
Use the Prototype pattern when:
1. Your code shouldn't depend on the concrete classes of objects that you need to copy
2. You want to reduce the number of subclasses that only differ in initialization
3. You need to instantiate classes at runtime that are specified by configuration or user input
4. You want to avoid building a parallel class hierarchy of factories
5. Creating an object is more expensive than copying an existing one

## Advantages
- You can clone objects without coupling to their concrete classes
- You can get rid of repeated initialization code in favor of cloning pre-built prototypes
- You can produce complex objects more conveniently
- You get an alternative to inheritance when dealing with configuration presets

## Disadvantages
- Cloning complex objects with circular references might be challenging
- Deep copying might be complex for objects with many fields and nested objects

## Example Use Case
Our implementation showcases a document generation system where different types of documents (reports, forms, contracts, invoices) can be created by cloning pre-configured templates.

### Sample Code
```go
// Create a registry
registry := NewDocumentRegistry()

// Register a prototype
reportTemplate := &ReportDocument{
    Document: Document{
        ID:   "REPORT-TEMPLATE",
        Name: "Quarterly Report Template",
    },
    Title: "Quarterly Financial Report",
    // ... more fields
}
registry.Register("quarterly-report", reportTemplate)

// Clone a prototype to create a new instance
reportDoc, err := registry.DeepClone("quarterly-report")
if err != nil {
    // Handle error
}

// Customize the cloned object
report := reportDoc.(*ReportDocument)
report.Title = "Q1 2025 Financial Performance"
```

### Running the Example
A complete example application is provided in the `example` directory. To run it:

```bash
cd creational/prototype/example
go run main.go
```

## Related Patterns
- **Abstract Factory**: Prototype can be an alternative to Abstract Factory when dealing with many possible classes
- **Composite**: Prototypes can be used with Composite to abstract the cloning process for complex structures
- **Decorator**: Prototype can help when you need to save the state of a decorated object
- **Command**: Prototype can be used to store command history or for saving command snapshots
- **Memento**: Prototype can sometimes be used as an alternative to Memento for saving object state

## Go-Specific Considerations
- Go doesn't have built-in cloning mechanisms, so manual implementation is required
- In Go, implementing deep copying requires careful handling of reference types (slices, maps, pointers)
- Go's lack of generics (before Go 1.18) makes type assertions necessary when working with the Prototype interface
- Thread safety should be considered in concurrent applications (our DocumentRegistry uses mutex for this)

## References
- "Design Patterns: Elements of Reusable Object-Oriented Software" by Gamma, Helm, Johnson, and Vlissides
- "Head First Design Patterns" by Eric Freeman and Elisabeth Robson
- "Design Patterns in Go" by Mario Castro Contreras
