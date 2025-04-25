package main

import (
	"fmt"
	"github.com/edgardnogueira/go-patterns/structural/flyweight"
)

func main() {
	fmt.Println("Flyweight Pattern Example")
	fmt.Println("========================")
	fmt.Println("Demonstrating memory efficiency in a text editor using the Flyweight pattern\n")

	// Create our factories
	textFormatFactory := flyweight.NewTextFormatFactory()
	paragraphFactory := flyweight.NewParagraphStyleFactory()

	// Create our document
	doc := flyweight.NewFormattedDocument("Demo Document", textFormatFactory, paragraphFactory)

	// Create some basic formats we'll use repeatedly
	normalFormat := textFormatFactory.GetTextFormat(
		"Times New Roman", 12, "black", false, false, false,
		"white", "left", 0, 1.0,
	)

	boldFormat := textFormatFactory.GetTextFormat(
		"Times New Roman", 12, "black", true, false, false,
		"white", "left", 0, 1.0,
	)

	italicFormat := textFormatFactory.GetTextFormat(
		"Times New Roman", 12, "black", false, true, false,
		"white", "left", 0, 1.0,
	)

	// Create a paragraph style
	defaultParagraph := paragraphFactory.GetParagraphStyle(
		"left", 1.2, 10, 10, 0, 20, 20,
		"none", "black", "transparent",
	)

	// Step 1: Add text with different formats
	fmt.Println("Step 1: Adding text with different formats...")
	
	// Add some formatted text
	doc.AddText("This is ", normalFormat.GetID())
	doc.AddText("bold ", boldFormat.GetID())
	doc.AddText("and this is ", normalFormat.GetID())
	doc.AddText("italic", italicFormat.GetID())
	doc.AddText(".", normalFormat.GetID())
	
	// Add a paragraph
	doc.AddParagraph("This is a paragraph with standard formatting applied.", defaultParagraph.GetID())

	// Step 2: Add more text with the same formats
	fmt.Println("Step 2: Adding more text with the same formats...")
	
	// The key point is that we're reusing the same format objects
	doc.AddText("\nHere is more ", normalFormat.GetID())
	doc.AddText("bold ", boldFormat.GetID())
	doc.AddText("text and more ", normalFormat.GetID())
	doc.AddText("italic ", italicFormat.GetID())
	doc.AddText("text, but we're reusing the same format objects.", normalFormat.GetID())
	
	// Add another paragraph using the same style
	doc.AddParagraph("Another paragraph using the same style object as before.", defaultParagraph.GetID())

	// Step 3: Check memory usage statistics
	fmt.Println("Step 3: Analyzing memory usage...\n")
	
	// Get memory statistics
	memStats := doc.GetMemoryUsage()
	
	fmt.Println("Memory Usage Statistics:")
	fmt.Println("------------------------")
	fmt.Printf("Total characters: %d\n", memStats["characterCount"].(int))
	fmt.Printf("Unique format objects: %d\n", memStats["uniqueFormatCount"].(int))
	fmt.Printf("Memory with flyweight pattern: %d bytes\n", memStats["totalMemory"].(int))
	fmt.Printf("Memory without flyweight pattern: %d bytes\n", memStats["memoryWithoutFlyweight"].(int))
	fmt.Printf("Memory saved: %d bytes (%.1f%%)\n", 
		memStats["memorySaved"].(int), 
		memStats["savingsPercent"].(float64))

	// Step 4: Check format cache stats
	fmt.Println("\nFormat Cache Statistics:")
	fmt.Println("------------------------")
	cacheStats := textFormatFactory.GetCacheStats()
	fmt.Printf("Total formats in cache: %d\n", cacheStats["totalFormats"].(int))
	
	// Step 5: Display formatted text
	fmt.Println("\nFormatted Document:")
	fmt.Println("------------------")
	fmt.Println(doc.GetFormattedDocument())
	
	// Step 6: Demonstrate serialization/deserialization
	fmt.Println("\nStep 6: Demonstrating serialization/deserialization...")
	
	// Serialize the document
	jsonData, err := doc.Serialize()
	if err != nil {
		fmt.Printf("Error serializing document: %v\n", err)
		return
	}
	
	fmt.Printf("Document serialized to JSON (%d bytes)\n", len(jsonData))
	
	// Create new factories for deserialization
	newTextFactory := flyweight.NewTextFormatFactory()
	newParaFactory := flyweight.NewParagraphStyleFactory()
	
	// Deserialize into a new document
	newDoc, err := flyweight.DeserializeDocument(jsonData, newTextFactory, newParaFactory)
	if err != nil {
		fmt.Printf("Error deserializing document: %v\n", err)
		return
	}
	
	fmt.Println("Document successfully deserialized")
	
	// Confirm the deserialized document has the same stats
	newMemStats := newDoc.GetMemoryUsage()
	fmt.Printf("Deserialized document has %d characters and %d formats\n", 
		newMemStats["characterCount"].(int),
		newMemStats["uniqueFormatCount"].(int))
	
	// Conclusion
	fmt.Println("\nConclusion:")
	fmt.Println("-----------")
	fmt.Println("The Flyweight pattern allowed us to significantly reduce memory usage")
	fmt.Println("by sharing format objects across many characters in the document.")
	fmt.Printf("We saved %.1f%% of memory compared to storing individual formats.\n", 
		memStats["savingsPercent"].(float64))
}
