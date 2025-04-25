# State Pattern

## Intent

The State pattern allows an object to alter its behavior when its internal state changes, appearing to change its class. It encapsulates state-dependent behavior into separate state objects, eliminating conditional statements.

## Explanation

In this implementation, we demonstrate a package delivery system with different states (ordered, processing, shipped, delivered, returned, canceled) that change the behavior of the package object.

The pattern enables an object to change its behavior when its internal state changes. The object appears to change its class, as each state can be represented by a separate class implementing a common interface. State transitions are either triggered explicitly or can happen automatically based on internal conditions.

## Structure

- **State Interface (PackageState)**: Defines the interface for encapsulating the behavior associated with a particular state of the context.
- **Concrete States**:
  - **OrderedState**: Initial state when a package is ordered
  - **ProcessingState**: Package is being processed in the warehouse
  - **ShippedState**: Package is in transit
  - **DeliveredState**: Package has been delivered to recipient
  - **ReturnedState**: Package is being returned
  - **CanceledState**: Order has been canceled
- **Context (Package)**: Maintains an instance of a concrete state as the current state and delegates state-specific behavior to it.

## When to Use

- When an object's behavior depends on its state, and it must change its behavior at runtime depending on that state
- When operations have large, multipart conditional statements that depend on the object's state
- When state transitions follow well-defined and consistent rules
- When you need to explicitly represent and manage the states and transitions in a structured way

## Benefits

1. **Encapsulates state-dependent behavior**: Each state encapsulates its own behavior, making the code more modular and maintainable.
2. **Makes state transitions explicit**: State transitions are clearly defined and validated, making the code easier to understand.
3. **Eliminates large conditional statements**: Each state handles its own behavior, eliminating the need for complex conditional logic.
4. **Simplifies context code**: The context delegates behavior to the current state, simplifying its implementation.
5. **Facilitates adding new states**: New states can be added without modifying existing code, following the Open/Closed Principle.
6. **Improves testability**: Individual states can be tested in isolation, making the code easier to test.

## Implementation Details

In our implementation:

1. The `PackageState` interface defines the behavior for all states
2. Each concrete state class implements how the package behaves in that state
3. The `Package` context maintains a reference to the current state
4. State transitions are validated using a transition validator
5. An event system notifies about state changes
6. A history of state transitions is maintained
7. Automatic transitions can be scheduled with timeouts

## Example

The example demonstrates a package delivery system with the following states:

```
Ordered → Processing → Shipped → Delivered
    ↓           ↓          ↓        ↓ 
  Cancel      Cancel     Return    Return
```

Each state defines which operations are valid and which state transitions are allowed. For example, you can't ship a package that hasn't been processed, and you can't process a package that has been canceled.

## Usage

```go
// Create a new package
pkg := state.NewPackage("PKG123", "Smartphone")
state.InitializePackage(pkg)

// Add a logging handler
pkg.AddTransitionHandler(state.LoggingHandler(func(message string) {
    fmt.Println(message)
}))

// Process the package (Ordered → Processing)
err := pkg.HandleProcess()

// Ship the package (Processing → Shipped)
err = pkg.HandleShip()

// Deliver the package (Shipped → Delivered)
err = pkg.HandleDeliver()

// Get the package's state history
history := pkg.GetStateHistory()
```

## Related Patterns

- **State vs Strategy**: While both patterns delegate behavior to another object, State focuses on changing behavior based on internal state, while Strategy allows selecting algorithms at runtime.
- **State vs Command**: Command encapsulates a request as an object, while State encapsulates state-dependent behavior.
- **State and Flyweight**: States can be shared using the Flyweight pattern when many objects are in the same state.
