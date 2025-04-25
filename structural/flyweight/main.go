package main

import (
	"fmt"
	"github.com/edgardnogueira/go-patterns/structural/flyweight"
)

func main() {
	// Create a formatted document example
	fmt.Println("Flyweight Pattern - Document Editor Example")
	fmt.Println("==========================================")
	fmt.Println("This example demonstrates how the Flyweight pattern can be used to")
	fmt.Println("efficiently manage memory in a text document editor by sharing formatting objects.")
	fmt.Println()

	// Create factories
	textFormatFactory := flyweight.NewTextFormatFactory()
	paragraphFactory := flyweight.NewParagraphStyleFactory()

	// Create a document
	doc := flyweight.NewFormattedDocument("Sample Document", textFormatFactory, paragraphFactory)

	// Create common formats
	titleFormat := textFormatFactory.GetTextFormat(
		"Arial", 18, "blue", true, false, false,
		"white", "center", 0, 1.5,
	)

	headingFormat := textFormatFactory.GetTextFormat(
		"Arial", 14, "navy", true, false, false,
		"white", "left", 0, 1.2,
	)

	bodyFormat := textFormatFactory.GetTextFormat(
		"Times New Roman", 12, "black", false, false, false,
		"white", "left", 0, 1.0,
	)

	emphasisFormat := textFormatFactory.GetTextFormat(
		"Times New Roman", 12, "black", false, true, false,
		"white", "left", 0, 1.0,
	)

	// Create paragraph styles
	centerStyle := paragraphFactory.GetParagraphStyle(
		"center", 1.5, 15, 10, 0, 10, 10,
		"none", "black", "transparent",
	)

	bodyStyle := paragraphFactory.GetParagraphStyle(
		"left", 1.2, 5, 5, 20, 10, 10,
		"none", "black", "transparent",
	)

	// Add text with different formats to our document
	doc.AddText("The Flyweight Pattern\n", titleFormat.GetID())
	doc.AddParagraph("A structural design pattern for efficient memory usage", centerStyle.GetID())

	doc.AddText("\nIntroduction\n", headingFormat.GetID())
	doc.AddText("The Flyweight pattern minimizes memory usage by sharing as much data as possible with similar objects. ", bodyFormat.GetID())
	doc.AddText("This is particularly useful ", bodyFormat.GetID())
	doc.AddText("when dealing with a large number of objects ", emphasisFormat.GetID())
	doc.AddText("that have similar characteristics.", bodyFormat.GetID())
	doc.AddParagraph("In this document editor example, we use flyweights to efficiently store text formatting information across the document.", bodyStyle.GetID())

	doc.AddText("\nImplementation\n", headingFormat.GetID())
	doc.AddText("Our implementation demonstrates a text formatting system where character formatting objects are shared. The key components include:", bodyFormat.GetID())
	doc.AddParagraph("- TextFormat interface and SharedTextFormat implementation\n- Character and Document structures\n- TextFormatFactory to manage shared formats", bodyStyle.GetID())

	doc.AddText("\nMemory Usage Analysis\n", headingFormat.GetID())
	
	// Display memory usage statistics
	memoryStats := doc.GetMemoryUsage()
	fmt.Println("\nDocument Memory Usage Statistics:")
	fmt.Println("-----------------------------------")
	fmt.Printf("Total Characters: %d\n", memoryStats["characterCount"].(int))
	fmt.Printf("Unique Formats: %d\n", memoryStats["uniqueFormatCount"].(int))
	fmt.Printf("Character Memory: %d bytes\n", memoryStats["characterMemory"].(int))
	fmt.Printf("Format Memory: %d bytes\n", memoryStats["formatMemory"].(int))
	fmt.Printf("Total Memory: %d bytes\n", memoryStats["totalMemory"].(int))
	fmt.Printf("Memory Without Flyweight: %d bytes\n", memoryStats["memoryWithoutFlyweight"].(int))
	fmt.Printf("Memory Saved: %d bytes (%.1f%%)\n", 
		memoryStats["memorySaved"].(int), 
		memoryStats["savingsPercent"].(float64))

	// Display formatted document
	fmt.Println("\nFormatted Document:")
	fmt.Println("-----------------------------------")
	fmt.Println(doc.GetFormattedDocument())

	// Show how the cache is working
	cacheStats := textFormatFactory.GetCacheStats()
	fmt.Println("\nFormat Cache Statistics:")
	fmt.Println("-----------------------------------")
	fmt.Printf("Total formats in cache: %d\n", cacheStats["totalFormats"].(int))
	if boldCount, ok := cacheStats["boldCount"].(int); ok {
		fmt.Printf("Bold formats: %d\n", boldCount)
	}
	if italicCount, ok := cacheStats["italicCount"].(int); ok {
		fmt.Printf("Italic formats: %d\n", italicCount)
	}
	
	// Conclusion
	fmt.Println("\nConclusion:")
	fmt.Println("The Flyweight pattern has allowed us to efficiently manage")
	fmt.Println("formatting in this document by sharing format objects across")
	fmt.Printf("the document, saving around %.1f%% of memory compared to not using flyweights.\n", 
		memoryStats["savingsPercent"].(float64))
}
