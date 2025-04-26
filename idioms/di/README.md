# Dependency Injection in Go

This directory contains examples of various dependency injection patterns implemented in Go.

## Overview

Dependency injection is a design pattern that allows for loose coupling between components by passing dependencies to an object rather than having the object create them. Go's simplicity and explicitness make dependency injection straightforward without requiring heavy frameworks.

## Patterns Included

1. **Constructor Injection**: Passing dependencies via constructors/factory functions
2. **Field Injection**: Setting dependencies after construction
3. **Method Injection**: Providing dependencies via method parameters
4. **Interface-based Injection**: Using interfaces for loose coupling
5. **Functional Dependency Injection**: Using higher-order functions
6. **Service Locator Pattern**: Centralized dependency registry (with caution)
7. **Wire-based DI**: Using Google's Wire code generation tool
8. **Dependency Injection Containers**: Lightweight DI containers in Go

## Directory Structure

- `/constructor_injection`: Examples of injecting dependencies through constructors
- `/field_injection`: Examples of injecting dependencies through struct fields after construction
- `/method_injection`: Examples of injecting dependencies through method parameters
- `/interface_injection`: Examples of using interfaces for flexible dependency injection
- `/functional_injection`: Examples of using higher-order functions for dependency injection
- `/service_locator`: Examples of using a service locator for centralized dependency management
- `/di_container`: Examples of lightweight dependency injection containers
- `/wire_di`: Examples of using Google's Wire tool for compile-time dependency injection

## Benefits of Dependency Injection in Go

- **Testability**: Makes unit testing easier by allowing mock dependencies
- **Modularity**: Promotes cleaner separation of concerns
- **Flexibility**: Components can be swapped without modifying consumers
- **Maintainability**: Dependencies are explicit rather than hidden
- **Reusability**: Components can be reused in different contexts

## Comparison of Different DI Approaches

| Approach | Pros | Cons | Best Use Cases |
|----------|------|------|---------------|
| Constructor Injection | Simple, explicit, compile-time safety | Can lead to large constructors | General purpose, most cases |
| Field Injection | Flexible, allows optional dependencies | Less explicit, runtime errors | Configuration objects, optional dependencies |
| Method Injection | Dependencies scoped to methods, reduced coupling | Can complicate method signatures | Per-operation dependencies |
| Interface Injection | Maximum decoupling, great for testing | Requires more code | Complex systems, testable code |
| Functional Injection | Composable, stateless, great for middleware | Can be hard to understand | HTTP middleware, decorators |
| Service Locator | Centralized dependency management | Hides dependencies, service name strings | Legacy systems, framework integration |
| DI Containers | Automates dependency creation, manages lifecycle | Runtime errors, reflection overhead | Larger applications |
| Wire DI | Compile-time safety, no runtime overhead | Requires code generation step | Production applications |

## Testing with Dependency Injection

Dependency injection greatly simplifies testing:

1. Define interfaces for dependencies
2. Create mock implementations for testing
3. Inject mocks instead of real implementations during tests
4. Verify interactions with dependencies

Examples of testing with DI can be found in each pattern's test files.

## Best Practices

- Prefer constructor injection when possible
- Use interfaces to define dependencies
- Keep dependencies explicit and visible
- Avoid global state and singletons
- Use the simplest approach that meets your needs
- Consider compile-time DI (like Wire) for larger applications
- Inject only what you need, when you need it

## Anti-Patterns

- Over-engineering with complex DI frameworks when simple approaches would suffice
- Using service locators that hide actual dependencies
- Circular dependencies
- God objects with too many injected dependencies
- Excessive layering of decorators in functional injection
- Relying on string-based lookup in service locators
- Runtime DI with insufficient error handling

## Real-World Go Projects Using DI

Many popular Go projects demonstrate good DI practices:

- The Go standard library's `http` package uses functional middleware patterns
- Kubernetes uses constructor injection extensively
- The `database/sql` package uses interface-based injection
- Many web frameworks use method injection for request handlers

## Further Reading

- [Dependency Injection in Go](https://github.com/google/wire/blob/main/docs/guide.md)
- [Wire User Guide](https://github.com/google/wire/blob/main/docs/guide.md)
- [Go Interfaces and DI](https://www.alexedwards.net/blog/interfaces-explained)
- [Dave Cheney's Functional Options Pattern](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)