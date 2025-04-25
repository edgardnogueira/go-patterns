# Flyweight Pattern

## Intent

The Flyweight pattern minimizes memory usage by sharing as much data as possible with similar objects. It aims to use shared objects to support large numbers of fine-grained objects efficiently.

## Explanation

This implementation demonstrates a text formatting system for a document editor where character formatting objects are shared across the document to reduce memory consumption when dealing with large documents.

In a typical text editor, formatting information (font family, size, color, etc.) for each character would require significant memory if stored individually. The Flyweight pattern allows us to share these formatting objects across many characters, drastically reducing memory usage.

## Structure

- **Flyweight Interface (TextFormat)**: Defines methods for applying text formatting
- **Concrete Flyweight (SharedTextFormat)**: Stores intrinsic (shared) state like font, size, color
- **Flyweight Factory (TextFormatFactory)**: Creates and manages flyweight objects
- **Context (Character)**: Stores extrinsic state (position, actual character)
- **Client (Document)**: Uses the flyweights and maintains context

## When to Use

- When an application uses a large number of objects with similar state
- When memory usage becomes a critical concern
- When most object state can be made extrinsic (stored outside the object)
- When many groups of objects may be replaced by relatively few shared ones

## Benefits

- Reduces memory usage by sharing common state between multiple objects
- Can dramatically improve performance in applications with many similar objects
- Separates intrinsic (shared) state from extrinsic (context-specific) state
- Makes memory usage more predictable and manageable

## Implementation Details

In our implementation:

1. **TextFormat** defines the flyweight interface for text formatting
2. **SharedTextFormat** is the concrete flyweight that stores formatting properties
3. **TextFormatFactory** manages the creation and caching of flyweight objects
4. **Character** represents an individual character with a reference to its format
5. **Document** manages collections of characters with their formatting
6. **ParagraphStyle** extends the pattern to paragraph-level formatting
7. **SerializedDocument** provides serialization/deserialization capabilities

The implementation demonstrates memory savings by tracking and comparing memory usage with and without the flyweight pattern.