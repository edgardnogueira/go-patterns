# Running the Decorator Pattern Example

This directory contains an implementation of the Decorator pattern in Go. The implementation demonstrates a text processing system where various behaviors (formatting, encryption, compression, etc.) can be added to text dynamically through decorators.

## Core Components

- **TextProcessor** interface: Defines the basic operations that can be performed on text
- **BasicTextProcessor**: The concrete component that implements the basic text processing functionality
- **TextProcessorDecorator**: The base decorator that wraps a TextProcessor and delegates operations to it
- **Concrete Decorators**: Various implementations that add specific behaviors:
  - Formatting decorators (HTML, Markdown, plain text)
  - Security decorators (encryption, decryption, hashing, validation)
  - Utility decorators (compression, logging, metadata, translation)

## Running the Example

To run the example, navigate to the example directory and run the main.go file:

```bash
cd structural/decorator/example
go run main.go
```

The example demonstrates:

1. Basic text processing
2. Adding HTML formatting to Markdown text
3. Chaining multiple decorators (logging, compression, encryption)
4. Unwrapping (decrypting and decompressing)
5. Formatting with highlighting
6. Adding metadata and validation

## Running the Tests

To run the tests, navigate to the decorator directory and run the go test command:

```bash
cd structural/decorator
go test -v
```

The tests cover:

- Basic text processor functionality
- Formatting decorators (HTML, Markdown, plain text)
- Security decorators (encryption, validation, hashing)
- Utility decorators (compression, metadata, highlighting)
- Chaining multiple decorators
- Processing chain visualization

## Decorator Pattern Benefits

1. **Single Responsibility Principle**: Each decorator has a specific functionality
2. **Open/Closed Principle**: You can add new functionality without changing existing code
3. **Runtime Flexibility**: Decorators can be added or removed dynamically
4. **Composition Over Inheritance**: Object composition is used to extend functionality

## Go-Specific Implementation Notes

In Go, the Decorator pattern is implemented through composition and interfaces rather than inheritance. The decorator and the component implement the same interface, and the decorator wraps the component, delegating requests to it while potentially adding behavior before or after.

### Example Usage in Real Code

```go
// Create a basic text processor
processor := decorator.NewBasicTextProcessor()

// Add encryption
processor = decorator.NewEncryptionDecorator(processor, "secret-key", "aes")

// Add compression
processor = decorator.NewCompressionDecorator(processor, "gzip", true)

// Add logging
processor = decorator.NewLoggingDecorator(processor, true, true, true, 100)

// Process text through the chain
result, err := processor.Process("Some text to process")
if err != nil {
    log.Fatalf("Processing error: %v", err)
}
```

## When to Use the Decorator Pattern

- When you need to add responsibilities to objects dynamically without affecting other objects
- When extension by subclassing is impractical or impossible
- When you need to add/remove responsibilities at runtime
- When you want to add functionalities in various combinations
