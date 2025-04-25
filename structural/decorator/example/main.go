package main

import (
	"fmt"
	"log"

	"github.com/edgardnogueira/go-patterns/structural/decorator"
)

func main() {
	fmt.Println("=== Decorator Pattern Example ===")
	fmt.Println("This example demonstrates the Decorator pattern using a text processing system.")
	fmt.Println("The pattern allows us to add behaviors to text processing dynamically.")
	fmt.Println()

	// Example 1: Basic Text Processing
	fmt.Println("=== Example 1: Basic Text Processing ===")
	basicProcessor := decorator.NewBasicTextProcessor()
	text := "Hello, this is a sample text for processing!"
	processedText, err := basicProcessor.Process(text)
	if err != nil {
		log.Fatalf("Error processing text: %v", err)
	}

	fmt.Printf("Original text: %s\n", text)
	fmt.Printf("Processed text: %s\n", processedText)
	fmt.Printf("Processor: %s - %s\n", basicProcessor.GetName(), basicProcessor.GetDescription())
	fmt.Println()

	// Example 2: Adding HTML Formatting
	fmt.Println("=== Example 2: Adding HTML Formatting ===")
	htmlProcessor := decorator.NewHTMLFormattingDecorator(basicProcessor)
	markdownText := "# Decorator Pattern\n\n" +
		"The **Decorator** pattern attaches additional responsibilities to an object dynamically.\n\n" +
		"- It provides a flexible alternative to subclassing\n" +
		"- It allows for extending functionality without modifying original code\n\n" +
		"[Learn more](https://en.wikipedia.org/wiki/Decorator_pattern)"

	processedHTML, err := htmlProcessor.Process(markdownText)
	if err != nil {
		log.Fatalf("Error processing text: %v", err)
	}

	fmt.Printf("Original markdown:\n%s\n\n", markdownText)
	fmt.Printf("Processed HTML:\n%s\n", processedHTML)
	fmt.Printf("Processor: %s - %s\n", htmlProcessor.GetName(), htmlProcessor.GetDescription())
	fmt.Println()

	// Example 3: Chaining Multiple Decorators
	fmt.Println("=== Example 3: Chaining Multiple Decorators ===")
	// Start with the basic processor
	processor := basicProcessor
	
	// Add logging
	processor = decorator.NewLoggingDecorator(processor, true, true, true, 100)
	
	// Add compression
	processor = decorator.NewCompressionDecorator(processor, "gzip", true)
	
	// Add encryption
	processor = decorator.NewEncryptionDecorator(processor, "my-secret-key", "base64")

	// Process some text through the chain
	secretMessage := "This is a confidential message that will be logged, compressed, and encrypted."
	
	fmt.Printf("Original text: %s\n", secretMessage)
	fmt.Println("Processing through a chain of decorators...")
	
	result, err := processor.Process(secretMessage)
	if err != nil {
		log.Fatalf("Error in processing chain: %v", err)
	}
	
	fmt.Printf("Result: %s\n", result)
	fmt.Println()

	// Example 4: Demonstrating unwrapping/decryption
	fmt.Println("=== Example 4: Unwrapping (Decryption) ===")
	
	// Create a chain to decrypt, decompress, and restore the text
	decryptProcessor := decorator.NewDecryptionDecorator(basicProcessor, "my-secret-key", "base64")
	decompressProcessor := decorator.NewDecompressionDecorator(decryptProcessor, "gzip", true)
	
	fmt.Println("Reversing the encryption and compression...")
	decrypted, err := decompressProcessor.Process(result)
	if err != nil {
		log.Fatalf("Error in unwrapping: %v", err)
	}
	
	fmt.Printf("Unwrapped result: %s\n", decrypted)
	fmt.Println()

	// Example 5: Formatting and Highlighting
	fmt.Println("=== Example 5: Formatting and Highlighting ===")
	
	// Start with markdown formatter
	markdownProcessor := decorator.NewMarkdownFormattingDecorator(basicProcessor)
	
	// Add highlighting for important words
	highlighter := decorator.NewHighlightingDecorator(
		markdownProcessor,
		"(?i)\\b(important|note|warning)\\b", 
		"**[", 
		"]**",
	)
	
	// Process text with important words
	infoText := "Note: This is an important message with a warning about the pattern."
	
	formatted, err := highlighter.Process(infoText)
	if err != nil {
		log.Fatalf("Error in formatting: %v", err)
	}
	
	fmt.Printf("Original: %s\n", infoText)
	fmt.Printf("Formatted with highlights: %s\n", formatted)
	fmt.Println()

	// Example 6: Using metadata and validation
	fmt.Println("=== Example 6: Metadata and Validation ===")
	
	// Create metadata
	metadata := map[string]string{
		"Author":      "Go Patterns Team",
		"Version":     "1.0",
		"Date":        "2025-04-25",
		"Description": "Decorator pattern example",
	}
	
	// Create a processor with validation and metadata
	validatingProcessor := decorator.NewValidationDecorator(
		basicProcessor,
		decorator.ValidateNotEmpty,
		decorator.ValidateMinLength(10),
		decorator.ValidateMaxLength(200),
	)
	
	metadataProcessor := decorator.NewMetadataDecorator(validatingProcessor, metadata, "suffix")
	
	// Process a valid text
	validText := "This text meets the validation criteria."
	
	result, err = metadataProcessor.Process(validText)
	if err != nil {
		log.Fatalf("Validation error: %v", err)
	}
	
	fmt.Printf("Result with metadata:\n%s\n", result)
	
	// Try with invalid text (too short)
	invalidText := "Too short"
	
	_, err = metadataProcessor.Process(invalidText)
	if err != nil {
		fmt.Printf("Expected validation error: %v\n", err)
	}
	
	fmt.Println()
	fmt.Println("=== Decorator Pattern Advantages ===")
	fmt.Println("1. Single Responsibility Principle: Each decorator has a specific functionality")
	fmt.Println("2. Open/Closed Principle: We can add new functionality without changing existing code")
	fmt.Println("3. Runtime Flexibility: Decorators can be added or removed dynamically")
	fmt.Println("4. Composition Over Inheritance: We use object composition to extend functionality")
}
