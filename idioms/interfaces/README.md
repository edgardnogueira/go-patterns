# Interface Implementation Patterns in Go

This directory demonstrates idiomatic Go interface implementation patterns and best practices. Interfaces in Go are implicitly implemented, which creates unique patterns and approaches compared to explicitly implemented interfaces in other languages.

## Understanding Go's Interface Philosophy

Go interfaces are a collection of method signatures that objects can implement implicitly. Unlike languages with explicit interface implementation (e.g., Java, C#), Go types automatically satisfy an interface if they implement all of its methods, without any explicit declaration.

### Key Concepts

- Interfaces are satisfied implicitly
- Small, focused interfaces are preferred (Interface Segregation)
- Accept interfaces, return concrete types
- Composition over inheritance
- Interfaces foster decoupling and testability

## Patterns Implemented

1. [Basic Interface Implementation](basic.go): Simple interface definition and usage
2. [Interface Composition](composition.go): Embedding interfaces within other interfaces
3. [Accept Interfaces, Return Structs](accept_return.go): Design pattern for flexible functions
4. [Interface Segregation](segregation.go): Small, focused interfaces
5. [Empty Interface Usage](empty_interface.go): Working with the empty interface and type assertions
6. [Duck Typing in Practice](duck_typing.go): Leveraging implicit interface satisfaction
7. [Testing with Interfaces](testing.go): Using interfaces for mocks and stubs
8. [Interface Upgrade Patterns](upgrade.go): Evolving interfaces over time
9. [Anti-Patterns](anti_patterns.go): Common interface mistakes to avoid
10. [Standard Library Interface Examples](stdlib_examples.go): Real-world interface patterns from Go's standard library

## Best Practices

1. **Keep Interfaces Small**:
   ```go
   // Good: Focused interface with a single responsibility
   type Reader interface {
       Read(p []byte) (n int, err error)
   }
   ```

2. **Use Interface Composition**:
   ```go
   // Building larger interfaces through composition
   type ReadWriter interface {
       Reader
       Writer
   }
   ```

3. **Accept Interfaces, Return Concrete Types**:
   ```go
   // Accept interface (flexible input)
   func Process(r Reader) *Result {
       // Implementation...
       return &Result{} // Return concrete type
   }
   ```

4. **Use Interfaces for Seams in Testing**:
   ```go
   // Interfaces make it easy to substitute implementations for testing
   type Service interface {
       FetchData() ([]string, error)
   }
   ```

5. **Design for the Consumer, Not the Implementation**:
   - Define interfaces where they are used, not where types are defined
   - Focus on behavior needed by consumers, not all behaviors provided

## Comparison with Inheritance-Based Languages

| Go's Interface Approach | Inheritance-Based Approach |
|-------------------------|----------------------------|
| Implicit implementation | Explicit implementation |
| Composition of behaviors | Inheritance hierarchy |
| Interface at use site | Interface at definition site |
| No "is-a" relationship | Strong "is-a" relationship |
| Duck typing | Type checking |
| Runtime polymorphism | Compile-time polymorphism |

## Testing Interface Implementations

This directory includes tests that demonstrate how to properly test interfaces and their implementations, including:

- Writing tests with mock implementations
- Using interface flexibility for test fixtures
- Testing interface compliance
- Testing behavior rather than implementation
