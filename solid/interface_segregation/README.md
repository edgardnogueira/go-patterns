# Interface Segregation Principle (ISP)

## Definition

The Interface Segregation Principle states that clients should not be forced to depend on methods they do not use. In other words, many client-specific interfaces are better than one general-purpose interface.

## Key Concepts

- Keep interfaces small, focused, and cohesive
- Interfaces should be designed from the client's perspective
- Avoid "fat" interfaces that force clients to implement methods they don't need
- Prefer many small interfaces over a few large ones

## Benefits

- Reduces coupling between components
- Improves readability and maintainability
- Makes the system more modular and easier to refactor
- Facilitates proper separation of concerns

## Go Implementation Example

This directory contains a real-world example of the Interface Segregation Principle implemented in Go. The example demonstrates how to design and use focused interfaces that meet specific client needs.
