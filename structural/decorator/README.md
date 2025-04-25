# Decorator Pattern

## Intent
The Decorator Pattern attaches additional responsibilities to an object dynamically. Decorators provide a flexible alternative to subclassing for extending functionality.

## Problem
You want to add responsibilities to individual objects, not to an entire class, and you want these additions to be dynamic (at runtime) and transparent to clients.

## Solution
The Decorator Pattern suggests creating a set of decorator classes that are used to wrap concrete components. Decorators mirror the type of the components they decorate (they implement the same interface) but add or override behavior. Multiple decorators can wrap a component, each adding different functionalities.

## Structure
- **Component**: Defines the interface for objects that can have responsibilities added to them dynamically.
- **ConcreteComponent**: The basic object that can have responsibilities added to it.
- **Decorator**: Maintains a reference to a Component object and defines an interface that conforms to Component's interface.
- **ConcreteDecorator**: Adds responsibilities to the component.

## Implementation
In this implementation, we create a text processing system where various decorators can add formatting, encryption, compression, validation, and other transformations to text data.

The key elements are:
- The TextProcessor interface (Component)
- The BaseTextProcessor concrete component
- The TextProcessorDecorator abstract decorator
- Multiple concrete decorators that add different functionalities

## When to use
- When you need to add responsibilities to objects dynamically without affecting other objects
- When extension by subclassing is impractical or impossible
- When you need to add/remove responsibilities at runtime
- When you want to add functionalities in various combinations

## Benefits
- Greater flexibility than static inheritance
- Avoids feature-laden classes high up in the hierarchy
- Allows for mixing and matching of responsibilities
- Enhances the Single Responsibility Principle by dividing functionality into separate classes
- Enables runtime composition of behavior

## Drawbacks
- Can result in many small objects with similar structure
- Can be complex for clients to understand the full behavior
- Can increase complexity when trying to debug the system

## Go-Specific Implementation Notes
In Go, the Decorator Pattern is implemented through composition and interfaces rather than inheritance. The decorator and the component implement the same interface, and the decorator wraps the component, delegating requests to it while potentially adding behavior before or after.
