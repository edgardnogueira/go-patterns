# Visitor Pattern

## Intent

The Visitor pattern lets you separate algorithms from the objects on which they operate. 
It allows you to add new operations to existing object structures without modifying them.

## Explanation

The Visitor pattern is a behavioral design pattern that represents an operation to be 
performed on the elements of an object structure. It lets you define a new operation 
without changing the classes of the elements on which it operates.

### Key Components

1. **Visitor Interface**: Declares a visit method for each type of concrete element.
2. **Concrete Visitors**: Implement the Visitor interface with operation implementations for each element type.
3. **Element Interface**: Defines an Accept method that takes a visitor object as an argument.
4. **Concrete Elements**: Implement the Element interface with their Accept method.
5. **Object Structure**: Usually a collection or a composite object that can enumerate its elements.

### Structure

```
             ┌──────────────┐
             │   Element    │
             │  Interface   │
             └──────────────┘
                    △
                    │
        ┌───────────┴───────────┐
        │                       │
┌───────────────┐        ┌──────────────┐
│  ElementA     │        │   ElementB   │
├───────────────┤        ├──────────────┤
│ Accept(v)     │        │  Accept(v)   │
└───────────────┘        └──────────────┘
        │                       │
        └───────────────────────┘
                    │
                    ▼
             ┌──────────────┐
             │   Visitor    │
             │  Interface   │
             └──────────────┘
                    △
                    │
        ┌───────────┴───────────┐
        │                       │
┌───────────────┐        ┌──────────────┐
│  VisitorA     │        │   VisitorB   │
├───────────────┤        ├──────────────┤
│ VisitElementA │        │ VisitElementA│
│ VisitElementB │        │ VisitElementB│
└───────────────┘        └──────────────┘
```

### Double Dispatch

The Visitor pattern uses a technique called "double dispatch" to determine which method to execute. When the `Accept(visitor)` method is called on an element, it calls the appropriate `Visit` method on the visitor, passing itself as an argument. This way, the concrete element type and the concrete visitor type both determine which method gets executed.

## When to Use

Use the Visitor pattern when:

1. You have a complex object structure with many distinct and unrelated operations to perform on them
2. You need to perform operations on all elements of an object structure
3. You want to add new operations to a class hierarchy without changing the classes
4. Classes defining the object structure rarely change, but you often want to define new operations on the structure
5. Related operations need to be grouped, but not necessarily in the element classes

## Benefits

1. **Separation of concerns**: Operations are kept separate from the objects they operate on
2. **Open/Closed Principle**: New operations can be added without modifying existing element classes
3. **Single Responsibility Principle**: Element classes focus on their primary behavior
4. **Accumulating state**: Visitors can maintain state as they visit elements
5. **Type-safety**: The compiler enforces that each concrete element has proper handling in each visitor

## Drawbacks

1. **Reduced encapsulation**: Elements must expose enough state for visitors to work with
2. **Rigid element hierarchy**: Adding or removing element types requires updating all visitor interfaces and implementations
3. **Complexity**: The double dispatch mechanism can be confusing to developers unfamiliar with the pattern
4. **Visitor traversal**: The pattern doesn't define how elements are traversed; this must be handled separately

## Implementation in Go

In Go, the Visitor pattern is implemented using interfaces. The Element interface defines an Accept method that takes a Visitor as an argument. Concrete Element types implement this interface and invoke the appropriate Visit method on the Visitor.

The Visitor interface declares Visit methods for each Concrete Element type. Concrete Visitors implement these methods to define operations for each element.

## Usage Examples

The Visitor pattern is commonly used in:

1. **Document object models**: Converting documents to different formats (HTML, PDF, plain text)
2. **Abstract syntax trees**: Compiler and interpreter components that process syntax trees
3. **Graphical editors**: Operations on different graphic elements without modifying their classes
4. **Static code analysis**: Analyzing different parts of the codebase with various analyzers

## Example in This Package

In this implementation, we demonstrate a document object model where different visitors can perform various operations on document elements without changing their classes.

### Element Types
- TextElement: Represents text content
- ImageElement: Represents images
- TableElement: Represents tables
- LinkElement: Represents hyperlinks
- CompositeElement: Can contain other elements (composite pattern)

### Visitor Types
- HTMLExportVisitor: Converts elements to HTML
- MarkdownExportVisitor: Converts elements to Markdown
- PlainTextExportVisitor: Extracts plain text content
- StatisticsVisitor: Collects statistics about the document
- SpellCheckVisitor: Checks spelling across document elements

## Real-World Use Cases

1. **Document Processing Systems**: Converting documents between formats
2. **Compilers and Interpreters**: Processing abstract syntax trees
3. **IDEs and Code Analysis Tools**: Performing operations on code structures
4. **CAD Systems**: Processing geometric shapes with different operations
5. **Graphics Processing**: Applying different operations to visual elements

## Related Patterns

- **Composite**: Often used with Visitor to operate on hierarchical structures
- **Iterator**: Can be used to traverse the elements in the object structure
- **Command**: Both can be used to parameterize operations, but with different focuses
- **Strategy**: Similar separation of algorithm from context, but with a simpler structure

## Sample Code Usage

```go
// Create a document structure
doc := &visitor.CompositeElement{Name: "Document"}
doc.AddChild(&visitor.TextElement{Content: "Hello, World!"})
doc.AddChild(&visitor.ImageElement{Source: "image.jpg", Alt: "Sample Image"})

// Export to HTML
htmlVisitor := visitor.NewHTMLExportVisitor()
doc.Accept(htmlVisitor)
fmt.Println(htmlVisitor.GetHTML())

// Collect statistics
statsVisitor := visitor.NewStatisticsVisitor()
doc.Accept(statsVisitor)
fmt.Printf("Text elements: %d\n", statsVisitor.TextCount)
```

## License

This implementation of the Visitor pattern is licensed under the MIT License.
