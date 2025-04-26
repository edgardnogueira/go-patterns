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

## Benefits of Dependency Injection in Go

- **Testability**: Makes unit testing easier by allowing mock dependencies
- **Modularity**: Promotes cleaner separation of concerns
- **Flexibility**: Components can be swapped without modifying consumers
- **Maintainability**: Dependencies are explicit rather than hidden
- **Reusability**: Components can be reused in different contexts

## Best Practices

- Prefer constructor injection when possible
- Use interfaces to define dependencies
- Keep dependencies explicit and visible
- Avoid global state and singletons
- Use the simplest approach that meets your needs

## Anti-Patterns

- Over-engineering with complex DI frameworks when simple approaches would suffice
- Using service locators that hide actual dependencies
- Circular dependencies
- God objects with too many injected dependencies
