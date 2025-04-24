# Bridge Pattern

## Intent
The Bridge Pattern separates an abstraction from its implementation so that the two can vary independently. It's particularly useful when both the abstraction and its implementation need to be extended using subclasses.

## Problem
When we have a class hierarchy that has different dimensions of variation, inheritance alone isn't enough. For instance, if you have different shapes (Circle, Square) and different renderers (Vector, Raster), creating a class for each combination leads to an explosion of subclasses.

## Solution
The Bridge Pattern suggests splitting the hierarchy into two: one for the abstraction (e.g., Shape) and another for the implementation (e.g., DrawingAPI). This way, both can evolve independently without affecting each other.

## Structure
- **Abstraction**: Defines the abstraction's interface and maintains a reference to an object of type Implementor.
- **RefinedAbstraction**: Extends the interface defined by Abstraction.
- **Implementor**: Defines the interface for implementation classes.
- **ConcreteImplementor**: Implements the Implementor interface.

## Implementation
In this implementation, we create a drawing application that demonstrates the Bridge Pattern by separating shapes (abstraction) from rendering methods (implementation).

The hierarchy is as follows:
- Shape (Abstraction)
  - Circle
  - Square
  - Triangle
  - Rectangle
- DrawingAPI (Implementor)
  - VectorRenderer
  - RasterRenderer
  - SVGRenderer
  - TextRenderer

## When to use
- When you want to avoid a permanent binding between an abstraction and its implementation.
- When both the abstraction and its implementation should be extensible by subclassing.
- When changes in the implementation should not impact the client code.
- When you have a proliferation of classes due to a combinatorial explosion of possibilities.

## Benefits
- Decouples interface from implementation.
- Improves extensibility.
- Hides implementation details from clients.
- Allows for dynamic switching of implementations at runtime.

## Drawbacks
- Increases complexity due to the introduction of additional interfaces and classes.
- Can be overkill for simple class hierarchies.

## Go-Specific Implementation Notes
In Go, we implement the pattern using interfaces and composition, rather than inheritance which Go doesn't support. The Shape interface contains a reference to a DrawingAPI interface, which provides the actual drawing capabilities.
