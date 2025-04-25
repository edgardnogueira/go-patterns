# Composite Pattern

## Intent
The Composite Pattern composes objects into tree structures to represent part-whole hierarchies. It lets clients treat individual objects and compositions of objects uniformly.

## Problem
When you need to represent a hierarchy where objects can be either primitive or containers that hold other objects (both primitive and containers), modeling this with inheritance alone can become complex. Clients would need to distinguish between primitives and containers to handle them differently.

## Solution
The Composite Pattern suggests creating a common interface for both containers and individual objects. This way, clients can treat them all the same, simplifying client code and enabling recursive compositions.

## Structure
- **Component**: The common interface for all concrete classes in the composition. It declares operations for both simple and complex elements.
- **Leaf**: Represents individual objects with no children. Implements the Component interface.
- **Composite**: Represents complex elements that can have children. It stores child components and implements child-related operations.
- **Client**: Works with elements through the Component interface.

## Implementation
In this implementation, we create a file system representation where:
- The `FileSystemNode` interface is the Component
- `File` is a Leaf 
- `Directory` is a Composite
- The file system can be traversed recursively with operations applied uniformly

## When to use
- When you want to represent part-whole hierarchies of objects
- When you want clients to ignore the difference between compositions of objects and individual objects
- When the structure can have an arbitrary depth and complexity
- When you want the client code to work uniformly with all objects in the composite structure

## Benefits
- Defines class hierarchies consisting of primitive and complex objects
- Makes client code simpler by treating all objects the same way
- Makes it easier to add new kinds of components
- Enables recursive operations across the entire structure

## Drawbacks
- Can make the design overly general when dealing with a restricted hierarchy
- Can sometimes be challenging to establish the constraints on the components

## Go-Specific Implementation Notes
In Go, we implement the pattern using interfaces and composition, rather than inheritance which Go doesn't support. The Component is defined as an interface, and both the Leaf and Composite types implement this interface.
