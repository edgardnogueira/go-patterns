package main

import (
	"flag"
	"fmt"
	"github.com/edgardnogueira/go-patterns/behavioral/visitor"
	"io/ioutil"
	"os"
	"strings"
)

// CustomHTMLExportVisitor extends the standard HTMLExportVisitor with custom styling
type CustomHTMLExportVisitor struct {
	visitor.HTMLExportVisitor
	CSSStyles string
}

// NewCustomHTMLExportVisitor creates a new CustomHTMLExportVisitor
func NewCustomHTMLExportVisitor() *CustomHTMLExportVisitor {
	v := &CustomHTMLExportVisitor{}
	v.CSSStyles = `
body { 
  font-family: Arial, sans-serif; 
  line-height: 1.6;
  margin: 0;
  padding: 20px;
  color: #333;
}
.document {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
  background: #fff;
  border-radius: 5px;
  box-shadow: 0 2px 5px rgba(0,0,0,0.1);
}
.header {
  border-bottom: 1px solid #eee;
  padding-bottom: 10px;
  margin-bottom: 20px;
}
h1 { color: #2c3e50; }
h2 { color: #3498db; }
a { color: #2980b9; }
table {
  width: 100%;
  border-collapse: collapse;
  margin: 15px 0;
}
th, td {
  border: 1px solid #ddd;
  padding: 8px;
  text-align: left;
}
th { background-color: #f8f9fa; }
img { max-width: 100%; }
`
	return v
}

// GetFullHTML returns a complete HTML document with styles
func (v *CustomHTMLExportVisitor) GetFullHTML() string {
	var sb strings.Builder
	sb.WriteString("<!DOCTYPE html>\n<html>\n<head>\n")
	sb.WriteString("<meta charset=\"UTF-8\">\n")
	sb.WriteString("<title>Document</title>\n")
	sb.WriteString("<style>\n")
	sb.WriteString(v.CSSStyles)
	sb.WriteString("</style>\n")
	sb.WriteString("</head>\n<body>\n")
	sb.WriteString(v.Output.String())
	sb.WriteString("</body>\n</html>")
	return sb.String()
}

// CustomMarkdownExportVisitor extends the standard MarkdownExportVisitor
// with front matter support for static site generators
type CustomMarkdownExportVisitor struct {
	visitor.MarkdownExportVisitor
	Title       string
	Author      string
	DateCreated string
	Tags        []string
}

// NewCustomMarkdownExportVisitor creates a new CustomMarkdownExportVisitor
func NewCustomMarkdownExportVisitor() *CustomMarkdownExportVisitor {
	return &CustomMarkdownExportVisitor{
		Title:       "Document",
		Author:      "Go Patterns",
		DateCreated: "2025-04-26",
		Tags:        []string{"go", "design-patterns", "visitor"},
	}
}

// GetMarkdownWithFrontMatter returns markdown with YAML front matter
func (v *CustomMarkdownExportVisitor) GetMarkdownWithFrontMatter() string {
	var sb strings.Builder
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("title: %s\n", v.Title))
	sb.WriteString(fmt.Sprintf("author: %s\n", v.Author))
	sb.WriteString(fmt.Sprintf("date: %s\n", v.DateCreated))
	sb.WriteString("tags:\n")
	for _, tag := range v.Tags {
		sb.WriteString(fmt.Sprintf("  - %s\n", tag))
	}
	sb.WriteString("---\n\n")
	sb.WriteString(v.Output.String())
	return sb.String()
}

func main() {
	// Parse command-line flags
	format := flag.String("format", "html", "Output format: html, markdown, text, stats")
	outputFile := flag.String("output", "", "Output file path (optional)")
	flag.Parse()

	// Create a sample document
	document := createSampleDocument()

	// Process the document based on the requested format
	var output string
	switch *format {
	case "html":
		visitor := NewCustomHTMLExportVisitor()
		document.Accept(visitor)
		output = visitor.GetFullHTML()
	case "markdown":
		visitor := NewCustomMarkdownExportVisitor()
		document.Accept(visitor)
		output = visitor.GetMarkdownWithFrontMatter()
	case "text":
		visitor := visitor.NewPlainTextExportVisitor()
		document.Accept(visitor)
		output = visitor.GetText()
	case "stats":
		statsVisitor := visitor.NewStatisticsVisitor()
		document.Accept(statsVisitor)
		output = fmt.Sprintf("Document Statistics:\n"+
			"-------------------\n"+
			"Text elements: %d\n"+
			"Image elements: %d\n"+
			"Table elements: %d\n"+
			"Link elements: %d\n"+
			"Total words: %d\n"+
			"Total characters: %d\n",
			statsVisitor.TextCount,
			statsVisitor.ImageCount,
			statsVisitor.TableCount,
			statsVisitor.LinkCount,
			statsVisitor.WordCount,
			statsVisitor.CharacterCount)
	default:
		fmt.Printf("Unsupported format: %s\n", *format)
		os.Exit(1)
	}

	// Write output to file or print to console
	if *outputFile != "" {
		err := ioutil.WriteFile(*outputFile, []byte(output), 0644)
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Output written to %s\n", *outputFile)
	} else {
		fmt.Println(output)
	}
}

// createSampleDocument creates a more complex document structure for demonstration
func createSampleDocument() *visitor.CompositeElement {
	// Create the main document
	document := &visitor.CompositeElement{Name: "Document"}
	
	// Add a header section
	header := &visitor.CompositeElement{Name: "Header"}
	header.AddChild(&visitor.TextElement{Content: "Visitor Pattern in Go"})
	header.AddChild(&visitor.ImageElement{
		Source: "go-visitor.png",
		Alt:    "Visitor Pattern Diagram",
		Width:  500,
		Height: 300,
	})
	document.AddChild(header)
	
	// Add introduction section
	intro := &visitor.CompositeElement{Name: "Introduction"}
	intro.AddChild(&visitor.TextElement{
		Content: "The Visitor pattern represents an operation to be performed on elements of an object structure. " +
			"It lets you define a new operation without changing the classes of the elements on which it operates.",
	})
	intro.AddChild(&visitor.TextElement{
		Content: "This pattern is particularly useful when working with complex object structures that have " +
			"multiple types of elements, and you need to perform various operations on these elements " +
			"that don't naturally belong in the element classes themselves.",
	})
	intro.AddChild(&visitor.LinkElement{
		URL:   "https://en.wikipedia.org/wiki/Visitor_pattern",
		Text:  "Wikipedia: Visitor Pattern",
		Title: "Learn more about the Visitor pattern",
	})
	document.AddChild(intro)
	
	// Add implementation section
	implementation := &visitor.CompositeElement{Name: "Implementation"}
	implementation.AddChild(&visitor.TextElement{
		Content: "In Go, the Visitor pattern can be implemented using interfaces. The key components are:",
	})
	
	// Table showing the components
	componentsTable := &visitor.TableElement{
		Rows:    5,
		Columns: 2,
		Data: [][]string{
			{"Component", "Description"},
			{"Element interface", "Defines an Accept method that takes a Visitor as an argument"},
			{"Concrete Elements", "Implement the Element interface"},
			{"Visitor interface", "Declares a Visit method for each type of Concrete Element"},
			{"Concrete Visitors", "Implement the Visitor interface with specific operations"},
		},
	}
	implementation.AddChild(componentsTable)
	
	// Code explanation
	implementation.AddChild(&visitor.TextElement{
		Content: "The key to the Visitor pattern is the 'double dispatch' mechanism, which helps determine " +
			"the correct Visit method to call based on both the Element type and Visitor type.",
	})
	document.AddChild(implementation)
	
	// Add benefits section
	benefits := &visitor.CompositeElement{Name: "Benefits"}
	benefits.AddChild(&visitor.TextElement{
		Content: "The Visitor pattern offers several advantages:",
	})
	
	benefitsList := &visitor.CompositeElement{Name: "BenefitsList"}
	benefitsList.AddChild(&visitor.TextElement{
		Content: "1. Separation of concerns: Operations are kept separate from the object structure",
	})
	benefitsList.AddChild(&visitor.TextElement{
		Content: "2. Open/Closed Principle: New operations can be added without modifying existing element classes",
	})
	benefitsList.AddChild(&visitor.TextElement{
		Content: "3. Accumulating state: Visitors can maintain state as they visit elements",
	})
	benefitsList.AddChild(&visitor.TextElement{
		Content: "4. Breaking dependencies: Elements don't need to know about operations performed on them",
	})
	benefits.AddChild(benefitsList)
	document.AddChild(benefits)
	
	// Add use cases section
	useCases := &visitor.CompositeElement{Name: "Use Cases"}
	useCases.AddChild(&visitor.TextElement{
		Content: "The Visitor pattern is useful in the following scenarios:",
	})
	
	useCasesTable := &visitor.TableElement{
		Rows:    4,
		Columns: 2,
		Data: [][]string{
			{"Use Case", "Example"},
			{"Document processing", "Converting documents to various formats (HTML, PDF, etc.)"},
			{"Abstract syntax trees", "Compilers and interpreters for processing language syntax"},
			{"Graph algorithms", "Performing different operations on nodes and edges of a graph"},
		},
	}
	useCases.AddChild(useCasesTable)
	document.AddChild(useCases)
	
	// Add conclusion
	conclusion := &visitor.CompositeElement{Name: "Conclusion"}
	conclusion.AddChild(&visitor.TextElement{
		Content: "The Visitor pattern is a powerful tool for working with complex object structures, " +
			"allowing for clean separation of algorithms from the objects they operate on. " +
			"While it requires more upfront design and can be complex to implement, " +
			"the benefits in terms of maintainability and extensibility can be significant for the right use cases.",
	})
	conclusion.AddChild(&visitor.LinkElement{
		URL:   "https://github.com/edgardnogueira/go-patterns",
		Text:  "Explore more Go design patterns",
		Title: "Go Design Patterns Repository",
	})
	document.AddChild(conclusion)
	
	return document
}
