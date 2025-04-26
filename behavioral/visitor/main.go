package main

import (
	"fmt"
	"github.com/edgardnogueira/go-patterns/behavioral/visitor"
)

func main() {
	// Create a document structure
	doc := createSampleDocument()

	// Export to HTML
	fmt.Println("HTML Export:")
	fmt.Println("------------")
	htmlVisitor := visitor.NewHTMLExportVisitor()
	doc.Accept(htmlVisitor)
	fmt.Println(htmlVisitor.GetHTML())
	
	// Export to Markdown
	fmt.Println("\nMarkdown Export:")
	fmt.Println("----------------")
	mdVisitor := visitor.NewMarkdownExportVisitor()
	doc.Accept(mdVisitor)
	fmt.Println(mdVisitor.GetMarkdown())
	
	// Export to Plain Text
	fmt.Println("\nPlain Text Export:")
	fmt.Println("------------------")
	txtVisitor := visitor.NewPlainTextExportVisitor()
	doc.Accept(txtVisitor)
	fmt.Println(txtVisitor.GetText())
	
	// Get Document Statistics
	fmt.Println("\nDocument Statistics:")
	fmt.Println("-------------------")
	statsVisitor := visitor.NewStatisticsVisitor()
	doc.Accept(statsVisitor)
	fmt.Printf("Text elements: %d\n", statsVisitor.TextCount)
	fmt.Printf("Image elements: %d\n", statsVisitor.ImageCount)
	fmt.Printf("Table elements: %d\n", statsVisitor.TableCount)
	fmt.Printf("Link elements: %d\n", statsVisitor.LinkCount)
	fmt.Printf("Total words: %d\n", statsVisitor.WordCount)
	fmt.Printf("Total characters: %d\n", statsVisitor.CharacterCount)
	
	// Check Spelling
	fmt.Println("\nSpell Check:")
	fmt.Println("------------")
	spellVisitor := visitor.NewSpellCheckVisitor()
	doc.Accept(spellVisitor)
	errors := spellVisitor.GetErrors()
	if len(errors) == 0 {
		fmt.Println("No spelling errors found.")
	} else {
		fmt.Printf("Found %d spelling errors:\n", len(errors))
		for _, err := range errors {
			fmt.Printf("- %s\n", err)
		}
	}
}

// createSampleDocument creates a sample document structure
func createSampleDocument() *visitor.CompositeElement {
	// Create the main document
	document := &visitor.CompositeElement{Name: "Document"}
	
	// Add a header section
	header := &visitor.CompositeElement{Name: "Header"}
	header.AddChild(&visitor.TextElement{Content: "Go Design Patterns"})
	header.AddChild(&visitor.ImageElement{
		Source: "go-logo.png",
		Alt:    "Go Language Logo",
		Width:  100,
		Height: 50,
	})
	document.AddChild(header)
	
	// Add an introduction section
	intro := &visitor.CompositeElement{Name: "Introduction"}
	intro.AddChild(&visitor.TextElement{
		Content: "This document demonstrates the Visitor design pattern in Go. " +
			"The Visitor pattern lets you separate algorithms from the objects on which they operate.",
	})
	intro.AddChild(&visitor.LinkElement{
		URL:   "https://refactoring.guru/design-patterns/visitor",
		Text:  "Learn more about Visitor pattern",
		Title: "Visitor Pattern Documentation",
	})
	document.AddChild(intro)
	
	// Add a features section
	features := &visitor.CompositeElement{Name: "Features"}
	features.AddChild(&visitor.TextElement{
		Content: "The Visitor pattern provides the following benefits:",
	})
	
	// Add a table of pattern benefits
	benefitsTable := &visitor.TableElement{
		Rows:    3,
		Columns: 2,
		Data: [][]string{
			{"Benefit", "Description"},
			{"Separation of concerns", "Algorithms are separate from object structures"},
			{"Open/Closed Principle", "Add new operations without modifying element classes"},
		},
	}
	features.AddChild(benefitsTable)
	document.AddChild(features)
	
	// Add a conclusion section
	conclusion := &visitor.CompositeElement{Name: "Conclusion"}
	conclusion.AddChild(&visitor.TextElement{
		Content: "Use the Visitor pattern when you need to perform operations on all elements " +
			"of a complex object structure, such as an abstract syntax tree or a document object model.",
	})
	conclusion.AddChild(&visitor.TextElement{
		Content: "This example showcases how different operations (HTML export, Markdown export, statistics gathering) " +
			"can be implemented without changing the element classes.",
	})
	document.AddChild(conclusion)
	
	return document
}
