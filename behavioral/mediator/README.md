# Mediator Pattern

## Intent

The Mediator pattern defines an object that encapsulates how a set of objects interact. This pattern promotes loose coupling by keeping objects from referring to each other explicitly, and it lets you vary their interaction independently.

## Explanation

The Mediator pattern is a behavioral design pattern that reduces chaotic dependencies between objects. Instead of communicating directly with each other, objects communicate through a mediator object. This centralizes control and reduces the connections between objects.

### Key Components

1. **Mediator Interface**: Defines the interface for communication between Colleague objects.
2. **Concrete Mediator**: Implements the Mediator interface and coordinates communication between Colleague objects. It maintains references to all the colleagues and facilitates their interaction.
3. **Colleague Interface**: Defines the interface for communication with other Colleagues through a Mediator.
4. **Concrete Colleagues**: Implement the Colleague interface and communicate with each other through the Mediator.

### Structure

```
             ┌───────────┐
             │  Mediator │
             │ Interface │
             └─────┬─────┘
                   │
                   │
             ┌─────▼────────┐
             │    Concrete  │◄────┐
             │    Mediator  │     │
             └──────────────┘     │
                   ▲               │
                   │               │
                   │               │
                   │               │
┌──────────────────┴───────────────┴──────────────────┐
│                                                      │
│                                                      │
│              ┌──────────────┐                        │
│              │   Colleague  │                        │
│              │   Interface  │                        │
│              └──────┬───────┘                        │
│                     │                                │
│                     │                                │
│        ┌────────────┴─────────────┐                 │
│        │                          │                 │
│  ┌─────▼────────┐         ┌──────▼─────────┐       │
│  │   Concrete   │         │    Concrete    │       │
│  │  Colleague A │         │   Colleague B  │       │
│  └──────────────┘         └────────────────┘       │
│                                                     │
└─────────────────────────────────────────────────────┘
```

## When to Use

Use the Mediator pattern when:

1. A set of objects communicate in well-defined but complex ways. The resulting interdependencies are unstructured and difficult to understand.
2. Reusing an object is difficult because it refers to and communicates with many other objects.
3. You want to customize a behavior that's distributed between several classes without creating too many subclasses.
4. Object protocols are very complex and/or numerous.
5. You want to create a centralized point of communication control and coordination in your application.

## Benefits

1. **Reduces coupling** between components by having them communicate indirectly through a mediator object.
2. **Simplifies object protocols** by replacing many-to-many interactions with one-to-many interactions.
3. **Centralizes control** by encapsulating collective behavior in a single mediator object.
4. **Increases component reusability** because they don't need to contain direct references to each other.
5. **Makes interactions more abstract and high-level** by separating object communication from object functionality.

## Drawbacks

1. **Centralization introduces a single point of failure**: If the mediator breaks, the entire system may fail.
2. **Mediator can become a bottleneck**: In complex systems, the mediator might handle too many responsibilities, becoming a "god object".
3. **Increased complexity of the mediator class**: As more functionality is added, the mediator can become harder to maintain.

## Implementation in Go

In Go, the Mediator pattern is implemented using interfaces. The Mediator interface declares methods to facilitate communication between colleagues. Concrete Mediator types implement this interface and maintain references to all colleague objects.

Colleagues communicate with the mediator through a Colleague interface. Each concrete colleague implements this interface and contains a reference to the mediator.

## Real-World Analogy

Air Traffic Control (ATC) is a classic example of the Mediator pattern. Planes don't communicate directly with each other to coordinate landings, takeoffs, and flight paths. Instead, they communicate with the control tower (the mediator), which coordinates all aircraft to ensure safe and efficient operations.

## Example in This Package

This implementation demonstrates an air traffic control system where a control tower (mediator) coordinates communication between various aircraft (colleagues). The system includes:

### Mediator Components
- **Mediator Interface**: Defines how aircraft can communicate with the control tower
- **AirTrafficControl**: A concrete mediator that registers aircraft and manages their communications

### Colleague Components
- **Aircraft Interface**: Defines how aircraft interact with the control tower
- **Different types of aircraft**:
  - PassengerAircraft: Commercial passenger planes
  - CargoAircraft: Cargo transport planes
  - PrivateAircraft: Private jets and small planes
  - MilitaryAircraft: Military aircraft with special protocols

### Communication Types
- Landing requests and clearances
- Takeoff requests and clearances
- Emergency notifications
- Position updates
- Weather alerts (broadcasts)

## Usage Example

```go
// Create a mediator (control tower)
controlTower := mediator.NewAirTrafficControl("JFK Tower")

// Create colleagues (aircraft)
flight1 := mediator.NewPassengerAircraft("UA123", "United Airlines", 220)
flight2 := mediator.NewCargoAircraft("FX456", "FedEx", 15000.0)

// Register aircraft with the control tower
controlTower.Register(flight1)
controlTower.Register(flight2)

// Aircraft 1 requests landing
flight1.RequestLanding()

// Aircraft 2 reports an emergency
flight2.ReportEmergency("Engine failure")

// Control tower broadcasts message to all aircraft
controlTower.Broadcast("JFK Tower", mediator.ControlMessage, 
    "Weather alert: Thunderstorm approaching", 8)
```

## Related Patterns

- **Observer**: Both handle communication between objects, but while Mediator centralizes communication, Observer broadcasts to all interested objects.
- **Facade**: Mediator focuses on coordinating interactions, while Facade provides a simpler interface to a subsystem.
- **Command**: Can be used with Mediator when commands need to be routed through a central point.
- **Singleton**: Mediators are often implemented as singletons since only one instance is typically needed.

## License

This implementation of the Mediator pattern is licensed under the MIT License.
