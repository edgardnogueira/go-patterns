package flyweight

import (
	"fmt"
	"strings"
	"testing"
)

// TestSharedTextFormat tests the basic functionality of SharedTextFormat
func TestSharedTextFormat(t *testing.T) {
	// Create a text format
	format := NewSharedTextFormat(
		"Arial", 12, "black", true, false, false,
		"white", "left", 0, 1.2,
	)
	
	// Test GetID
	if format.GetID() == "" {
		t.Error("Format ID should not be empty")
	}
	
	// Test GetFormat
	properties := format.GetFormat()
	if properties["fontFamily"] != "Arial" {
		t.Errorf("Expected fontFamily to be 'Arial', got '%s'", properties["fontFamily"])
	}
	if properties["fontSize"] != 12 {
		t.Errorf("Expected fontSize to be 12, got %d", properties["fontSize"])
	}
	if properties["isBold"] != true {
		t.Errorf("Expected isBold to be true, got %v", properties["isBold"])
	}
	
	// Test ApplyFormatting
	formatted := format.ApplyFormatting("Hello World")
	if !strings.Contains(formatted, "Hello World") {
		t.Errorf("Formatted text should contain 'Hello World', got '%s'", formatted)
	}
	if !strings.Contains(formatted, "bold") {
		t.Errorf("Formatted text should indicate bold formatting, got '%s'", formatted)
	}
	
	// Test String
	str := format.String()
	if !strings.Contains(str, "Arial") || !strings.Contains(str, "12px") {
		t.Errorf("String representation should include format details, got '%s'", str)
	}
}

// TestTextFormatFactory tests the flyweight factory functionality
func TestTextFormatFactory(t *testing.T) {
	factory := NewTextFormatFactory()
	
	// Get a format
	format1 := factory.GetTextFormat(
		"Arial", 12, "black", true, false, false,
		"white", "left", 0, 1.2,
	)
	
	// Get the same format again
	format2 := factory.GetTextFormat(
		"Arial", 12, "black", true, false, false,
		"white", "left", 0, 1.2,
	)
	
	// They should be the same object
	if format1 != format2 {
		t.Error("Factory should return the same object for identical formats")
	}
	
	// Get a different format
	format3 := factory.GetTextFormat(
		"Times New Roman", 14, "blue", false, true, false,
		"white", "center", 0, 1.5,
	)
	
	// They should be different objects
	if format1 == format3 {
		t.Error("Factory should return different objects for different formats")
	}
	
	// Check cache stats
	stats := factory.GetCacheStats()
	if stats["totalFormats"].(int) != 2 {
		t.Errorf("Expected 2 formats in cache, got %d", stats["totalFormats"].(int))
	}
	
	// Clear cache and check stats
	factory.ClearCache()
	stats = factory.GetCacheStats()
	if stats["totalFormats"].(int) != 0 {
		t.Errorf("Expected 0 formats after clearing cache, got %d", stats["totalFormats"].(int))
	}
}

// TestDocument tests the Document functionality
func TestDocument(t *testing.T) {
	factory := NewTextFormatFactory()
	doc := NewDocument("Test Document", factory)
	
	// Create some formats
	boldFormat := factory.GetTextFormat(
		"Arial", 12, "black", true, false, false,
		"white", "left", 0, 1.2,
	)
	
	italicFormat := factory.GetTextFormat(
		"Arial", 12, "black", false, true, false,
		"white", "left", 0, 1.2,
	)
	
	// Add text with different formats
	doc.AddText("Hello ", boldFormat.GetID())
	doc.AddText("World!", italicFormat.GetID())
	
	// Check plain text
	if doc.PlainText() != "Hello World!" {
		t.Errorf("Expected plain text 'Hello World!', got '%s'", doc.PlainText())
	}
	
	// Check formatted text
	formatted := doc.GetFormattedText()
	if !strings.Contains(formatted, "bold") || !strings.Contains(formatted, "italic") {
		t.Errorf("Formatted text should contain both bold and italic sections, got '%s'", formatted)
	}
	
	// Check memory usage
	memory := doc.GetMemoryUsage()
	if memory["characterCount"].(int) != 12 {
		t.Errorf("Expected 12 characters, got %d", memory["characterCount"].(int))
	}
	if memory["uniqueFormatCount"].(int) != 2 {
		t.Errorf("Expected 2 unique formats, got %d", memory["uniqueFormatCount"].(int))
	}
	if memory["savingsPercent"].(float64) <= 0 {
		t.Errorf("Expected memory savings > 0, got %f", memory["savingsPercent"].(float64))
	}
}

// TestParagraphStyle tests the paragraph styling functionality
func TestParagraphStyle(t *testing.T) {
	// Create a paragraph style
	style := NewSharedParagraphStyle(
		"center", 1.5, 10, 10, 20, 15, 15,
		"single", "black", "transparent",
	)
	
	// Test GetID
	if style.GetID() == "" {
		t.Error("Style ID should not be empty")
	}
	
	// Test GetStyle
	properties := style.GetStyle()
	if properties["alignment"] != "center" {
		t.Errorf("Expected alignment to be 'center', got '%s'", properties["alignment"])
	}
	if properties["lineSpacing"] != 1.5 {
		t.Errorf("Expected lineSpacing to be 1.5, got %f", properties["lineSpacing"].(float64))
	}
	
	// Test FormatParagraph
	formatted := style.FormatParagraph("Test paragraph")
	if !strings.Contains(formatted, "Test paragraph") {
		t.Errorf("Formatted paragraph should contain 'Test paragraph', got '%s'", formatted)
	}
	if !strings.Contains(formatted, "center align") {
		t.Errorf("Formatted paragraph should indicate center alignment, got '%s'", formatted)
	}
	
	// Test String
	str := style.String()
	if !strings.Contains(str, "center") || !strings.Contains(str, "1.5") {
		t.Errorf("String representation should include style details, got '%s'", str)
	}
}

// TestParagraphStyleFactory tests the paragraph style factory functionality
func TestParagraphStyleFactory(t *testing.T) {
	factory := NewParagraphStyleFactory()
	
	// Get a style
	style1 := factory.GetParagraphStyle(
		"center", 1.5, 10, 10, 20, 15, 15,
		"single", "black", "transparent",
	)
	
	// Get the same style again
	style2 := factory.GetParagraphStyle(
		"center", 1.5, 10, 10, 20, 15, 15,
		"single", "black", "transparent",
	)
	
	// They should be the same object
	if style1 != style2 {
		t.Error("Factory should return the same object for identical styles")
	}
	
	// Get a different style
	style3 := factory.GetParagraphStyle(
		"left", 1.2, 5, 5, 0, 10, 10,
		"none", "black", "transparent",
	)
	
	// They should be different objects
	if style1 == style3 {
		t.Error("Factory should return different objects for different styles")
	}
	
	// Check cache stats
	stats := factory.GetCacheStats()
	if stats["totalStyles"].(int) != 2 {
		t.Errorf("Expected 2 styles in cache, got %d", stats["totalStyles"].(int))
	}
}

// TestFormattedDocument tests the formatted document functionality
func TestFormattedDocument(t *testing.T) {
	textFactory := NewTextFormatFactory()
	paraFactory := NewParagraphStyleFactory()
	doc := NewFormattedDocument("Test Document", textFactory, paraFactory)
	
	// Create some formats
	boldFormat := textFactory.GetTextFormat(
		"Arial", 12, "black", true, false, false,
		"white", "left", 0, 1.2,
	)
	
	// Create a paragraph style
	centerStyle := paraFactory.GetParagraphStyle(
		"center", 1.5, 10, 10, 20, 15, 15,
		"single", "black", "transparent",
	)
	
	// Add text with formatting
	doc.AddText("Hello World!", boldFormat.GetID())
	
	// Add a paragraph
	doc.AddParagraph("This is a test paragraph.", centerStyle.GetID())
	
	// Get the formatted document
	formatted := doc.GetFormattedDocument()
	if !strings.Contains(formatted, "This is a test paragraph") {
		t.Errorf("Formatted document should contain the paragraph text, got '%s'", formatted)
	}
	if !strings.Contains(formatted, "center align") {
		t.Errorf("Formatted document should indicate paragraph alignment, got '%s'", formatted)
	}
}

// TestSerialization tests the serialization and deserialization functionality
func TestSerialization(t *testing.T) {
	// Create a document with some content
	textFactory := NewTextFormatFactory()
	paraFactory := NewParagraphStyleFactory()
	doc := NewFormattedDocument("Test Document", textFactory, paraFactory)
	
	// Create some formats
	boldFormat := textFactory.GetTextFormat(
		"Arial", 12, "black", true, false, false,
		"white", "left", 0, 1.2,
	)
	
	italicFormat := textFactory.GetTextFormat(
		"Times New Roman", 14, "blue", false, true, false,
		"white", "left", 0, 1.5,
	)
	
	// Create a paragraph style
	centerStyle := paraFactory.GetParagraphStyle(
		"center", 1.5, 10, 10, 20, 15, 15,
		"single", "black", "transparent",
	)
	
	// Add text with different formats
	doc.AddText("Hello ", boldFormat.GetID())
	doc.AddText("World!", italicFormat.GetID())
	
	// Add a paragraph
	doc.AddParagraph("This is a test paragraph.", centerStyle.GetID())
	
	// Serialize the document
	serialized, err := doc.Serialize()
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}
	
	// Check that the serialized data contains expected content
	if !strings.Contains(serialized, "Test Document") {
		t.Errorf("Serialized data should contain document name, got '%s'", serialized)
	}
	
	// Create new factories for deserialization
	newTextFactory := NewTextFormatFactory()
	newParaFactory := NewParagraphStyleFactory()
	
	// Deserialize the document
	newDoc, err := DeserializeDocument(serialized, newTextFactory, newParaFactory)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}
	
	// Check that the deserialized document contains the same content
	if newDoc.Name != "Test Document" {
		t.Errorf("Expected document name 'Test Document', got '%s'", newDoc.Name)
	}
	
	if len(newDoc.Characters) != len(doc.Characters) {
		t.Errorf("Expected %d characters, got %d", len(doc.Characters), len(newDoc.Characters))
	}
	
	if len(newDoc.Paragraphs) != len(doc.Paragraphs) {
		t.Errorf("Expected %d paragraphs, got %d", len(doc.Paragraphs), len(newDoc.Paragraphs))
	}
	
	// Check formatted output
	origFormatted := doc.GetFormattedDocument()
	newFormatted := newDoc.GetFormattedDocument()
	if origFormatted != newFormatted {
		t.Errorf("Formatted output should match after serialization/deserialization")
		t.Errorf("Original: %s", origFormatted)
		t.Errorf("New: %s", newFormatted)
	}
}

// BenchmarkWithFlyweight benchmarks document creation with the flyweight pattern
func BenchmarkWithFlyweight(b *testing.B) {
	// Setup
	textFactory := NewTextFormatFactory()
	
	// Create some formats
	format1 := textFactory.GetTextFormat(
		"Arial", 12, "black", true, false, false,
		"white", "left", 0, 1.2,
	)
	
	format2 := textFactory.GetTextFormat(
		"Times New Roman", 14, "blue", false, true, false,
		"white", "left", 0, 1.5,
	)
	
	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		doc := NewDocument(fmt.Sprintf("Document %d", i), textFactory)
		
		// Add lots of text with just a few formats
		for j := 0; j < 1000; j++ {
			if j % 2 == 0 {
				doc.AddText("A", format1.GetID())
			} else {
				doc.AddText("B", format2.GetID())
			}
		}
	}
}

// BenchmarkWithoutFlyweight simulates document creation without flyweight pattern
func BenchmarkWithoutFlyweight(b *testing.B) {
	// Setup
	textFactory := NewTextFormatFactory()
	
	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		doc := NewDocument(fmt.Sprintf("Document %d", i), textFactory)
		
		// Simulate creating a new format for each character
		for j := 0; j < 1000; j++ {
			if j % 2 == 0 {
				format := textFactory.GetTextFormat(
					"Arial", 12, "black", true, false, false,
					"white", "left", 0, 1.2,
				)
				doc.AddText("A", format.GetID())
			} else {
				format := textFactory.GetTextFormat(
					"Times New Roman", 14, "blue", false, true, false,
					"white", "left", 0, 1.5,
				)
				doc.AddText("B", format.GetID())
			}
			
			// Clear cache to simulate not using flyweight
			textFactory.ClearCache()
		}
	}
}
