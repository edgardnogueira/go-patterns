package decorator

import (
	"strings"
	"testing"
)

// TestBasicTextProcessor tests the basic processor functionality
func TestBasicTextProcessor(t *testing.T) {
	processor := NewBasicTextProcessor()
	
	testText := "Sample text"
	result, err := processor.Process(testText)
	
	if err != nil {
		t.Fatalf("Basic processor returned error: %v", err)
	}
	
	if result != testText {
		t.Errorf("Basic processor changed the text. Expected '%s', got '%s'", testText, result)
	}
	
	if processor.GetName() != "Basic Text Processor" {
		t.Errorf("Incorrect processor name. Expected 'Basic Text Processor', got '%s'", processor.GetName())
	}
}

// TestFormattingDecorators tests the formatting decorators
func TestFormattingDecorators(t *testing.T) {
	basic := NewBasicTextProcessor()
	
	// Test HTML Formatter
	t.Run("HTML Formatter", func(t *testing.T) {
		htmlFormatter := NewHTMLFormattingDecorator(basic)
		markdown := "# Title\n\nThis is **bold**."
		result, err := htmlFormatter.Process(markdown)
		
		if err != nil {
			t.Fatalf("HTML formatter returned error: %v", err)
		}
		
		if !strings.Contains(result, "<h1>Title</h1>") {
			t.Errorf("HTML formatting didn't convert heading properly: %s", result)
		}
		
		if !strings.Contains(result, "<strong>bold</strong>") {
			t.Errorf("HTML formatting didn't convert bold properly: %s", result)
		}
	})
	
	// Test Markdown Formatter
	t.Run("Markdown Formatter", func(t *testing.T) {
		markdownFormatter := NewMarkdownFormattingDecorator(basic)
		text := "Important title"
		result, err := markdownFormatter.Process(text)
		
		if err != nil {
			t.Fatalf("Markdown formatter returned error: %v", err)
		}
		
		if !strings.HasPrefix(result, "# Important title") {
			t.Errorf("Markdown formatting didn't add heading properly: %s", result)
		}
	})
	
	// Test Plain Text Formatter (removes formatting)
	t.Run("Plain Text Formatter", func(t *testing.T) {
		plainTextFormatter := NewPlainTextFormattingDecorator(basic)
		formatted := "<h1>Title</h1>\n\nThis is **bold** and [a link](https://example.com)."
		result, err := plainTextFormatter.Process(formatted)
		
		if err != nil {
			t.Fatalf("Plain text formatter returned error: %v", err)
		}
		
		if strings.Contains(result, "<h1>") || strings.Contains(result, "**") {
			t.Errorf("Plain text formatting didn't remove tags properly: %s", result)
		}
		
		if !strings.Contains(result, "Title") || !strings.Contains(result, "bold") || !strings.Contains(result, "a link") {
			t.Errorf("Plain text formatting removed content: %s", result)
		}
	})
}

// TestSecurityDecorators tests the security-related decorators
func TestSecurityDecorators(t *testing.T) {
	basic := NewBasicTextProcessor()
	
	// Test Encryption/Decryption
	t.Run("Encryption and Decryption", func(t *testing.T) {
		// Use simpler encryption for predictable test results
		encryptor := NewEncryptionDecorator(basic, "test-key", "base64")
		decryptor := NewDecryptionDecorator(basic, "test-key", "base64")
		
		original := "Secret message for testing"
		encrypted, err := encryptor.Process(original)
		
		if err != nil {
			t.Fatalf("Encryption returned error: %v", err)
		}
		
		if encrypted == original {
			t.Errorf("Text wasn't encrypted: %s", encrypted)
		}
		
		decrypted, err := decryptor.Process(encrypted)
		
		if err != nil {
			t.Fatalf("Decryption returned error: %v", err)
		}
		
		if decrypted != original {
			t.Errorf("Decryption failed. Expected '%s', got '%s'", original, decrypted)
		}
	})
	
	// Test Validation
	t.Run("Validation", func(t *testing.T) {
		validator := NewValidationDecorator(
			basic,
			ValidateNotEmpty,
			ValidateMinLength(5),
			ValidateMaxLength(20),
		)
		
		// Valid text
		_, err := validator.Process("Valid text")
		if err != nil {
			t.Errorf("Validation failed for valid text: %v", err)
		}
		
		// Empty text
		_, err = validator.Process("")
		if err == nil {
			t.Errorf("Validation should have failed for empty text")
		}
		
		// Too short
		_, err = validator.Process("Hi")
		if err == nil {
			t.Errorf("Validation should have failed for too short text")
		}
		
		// Too long
		_, err = validator.Process("This text is definitely too long for validation to pass")
		if err == nil {
			t.Errorf("Validation should have failed for too long text")
		}
	})
	
	// Test Hashing
	t.Run("Hashing", func(t *testing.T) {
		hasher := NewHashingDecorator(basic, "md5", false)
		text := "Test hashing"
		hash, err := hasher.Process(text)
		
		if err != nil {
			t.Fatalf("Hashing returned error: %v", err)
		}
		
		if len(hash) != 32 { // MD5 hash is 32 hex chars
			t.Errorf("Invalid MD5 hash length. Expected 32, got %d: %s", len(hash), hash)
		}
		
		// Test with hash appending
		hasherAppend := NewHashingDecorator(basic, "md5", true)
		result, err := hasherAppend.Process(text)
		
		if err != nil {
			t.Fatalf("Hashing with append returned error: %v", err)
		}
		
		if !strings.Contains(result, "Hash:") || !strings.Contains(result, text) {
			t.Errorf("Hash wasn't properly appended: %s", result)
		}
	})
}

// TestUtilityDecorators tests the utility decorators
func TestUtilityDecorators(t *testing.T) {
	basic := NewBasicTextProcessor()
	
	// Test Compression/Decompression
	t.Run("Compression and Decompression", func(t *testing.T) {
		compressor := NewCompressionDecorator(basic, "gzip", true)
		decompressor := NewDecompressionDecorator(basic, "gzip", true)
		
		original := "This is a test message that should be compressible with lots of repeating text. " +
			"This is a test message that should be compressible with lots of repeating text."
		
		compressed, err := compressor.Process(original)
		
		if err != nil {
			t.Fatalf("Compression returned error: %v", err)
		}
		
		if compressed == original {
			t.Errorf("Text wasn't compressed")
		}
		
		decompressed, err := decompressor.Process(compressed)
		
		if err != nil {
			t.Fatalf("Decompression returned error: %v", err)
		}
		
		if decompressed != original {
			t.Errorf("Decompression failed. Expected '%s', got '%s'", original, decompressed)
		}
	})
	
	// Test Metadata
	t.Run("Metadata", func(t *testing.T) {
		metadata := map[string]string{
			"Author": "Test",
			"Date":   "2025-04-25",
		}
		
		prefixMeta := NewMetadataDecorator(basic, metadata, "prefix")
		suffixMeta := NewMetadataDecorator(basic, metadata, "suffix")
		
		text := "Sample text"
		
		// Test prefix metadata
		prefixResult, err := prefixMeta.Process(text)
		if err != nil {
			t.Fatalf("Prefix metadata returned error: %v", err)
		}
		
		if !strings.HasPrefix(prefixResult, "--- Metadata ---") || !strings.Contains(prefixResult, text) {
			t.Errorf("Prefix metadata not applied correctly: %s", prefixResult)
		}
		
		// Test suffix metadata
		suffixResult, err := suffixMeta.Process(text)
		if err != nil {
			t.Fatalf("Suffix metadata returned error: %v", err)
		}
		
		if !strings.HasSuffix(suffixResult, "-------------") || !strings.Contains(suffixResult, text) {
			t.Errorf("Suffix metadata not applied correctly: %s", suffixResult)
		}
	})
	
	// Test Highlighting
	t.Run("Highlighting", func(t *testing.T) {
		highlighter := NewHighlightingDecorator(basic, "important|critical", "**", "**")
		
		text := "This is an important message with critical information."
		result, err := highlighter.Process(text)
		
		if err != nil {
			t.Fatalf("Highlighting returned error: %v", err)
		}
		
		if !strings.Contains(result, "**important**") || !strings.Contains(result, "**critical**") {
			t.Errorf("Highlighting not applied correctly: %s", result)
		}
	})
}

// TestChainOfDecorators tests chaining multiple decorators together
func TestChainOfDecorators(t *testing.T) {
	// Create a chain of decorators
	processor := NewBasicTextProcessor()
	processor = NewHighlightingDecorator(processor, "important", "_", "_")
	processor = NewMarkdownFormattingDecorator(processor)
	processor = NewLoggingDecorator(processor, false, false, false, 0)
	
	// Process text through the chain
	text := "This is an important test of decorator chaining."
	result, err := processor.Process(text)
	
	if err != nil {
		t.Fatalf("Decorator chain returned error: %v", err)
	}
	
	// The text should be processed by all decorators in order
	if !strings.Contains(result, "_important_") {
		t.Errorf("Highlighting wasn't applied: %s", result)
	}
	
	if !strings.HasPrefix(result, "# This") {
		t.Errorf("Markdown formatting wasn't applied: %s", result)
	}
	
	// Test that the chain order is correctly maintained
	if textDecorator, ok := processor.(*LoggingDecorator); ok {
		if formattingDecorator, ok := textDecorator.wrapped.(*HighlightingDecorator); ok {
			if basicProcessor, ok := formattingDecorator.wrapped.(*MarkdownFormattingDecorator); ok {
				// Chain is in wrong order
				t.Errorf("Decorator chain is in wrong order: %T -> %T -> %T", 
					textDecorator, formattingDecorator, basicProcessor)
			}
		}
	}
}

// TestGetProcessingChain tests the GetProcessingChain method
func TestGetProcessingChain(t *testing.T) {
	basic := NewBasicTextProcessor()
	decorated := NewHighlightingDecorator(
		NewMarkdownFormattingDecorator(
			NewLoggingDecorator(basic, false, false, false, 0),
		),
		"test", "*", "*",
	)
	
	chain := decorated.(*HighlightingDecorator).GetProcessingChain()
	
	// Check that all decorators are listed in the chain
	if !strings.Contains(chain, "Text Highlighter") ||
	   !strings.Contains(chain, "Markdown Formatter") ||
	   !strings.Contains(chain, "Logging Processor") ||
	   !strings.Contains(chain, "Basic Text Processor") {
		t.Errorf("Processing chain doesn't list all decorators: %s", chain)
	}
	
	// Check the order of the chain
	if !strings.HasPrefix(chain, "Text Highlighter") {
		t.Errorf("Processing chain has incorrect order: %s", chain)
	}
	
	if !strings.HasSuffix(chain, "Basic Text Processor") {
		t.Errorf("Processing chain doesn't end with the base processor: %s", chain)
	}
}
