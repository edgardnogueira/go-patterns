# Factory Method Pattern

## Intent

The Factory Method pattern defines an interface for creating an object, but lets subclasses decide which class to instantiate. Factory Method lets a class defer instantiation to subclasses.

## Explanation

In this example, we have a logistics management system where we need to create different types of transport vehicles (trucks and ships). Using the Factory Method pattern, we can let the specific logistics services decide which type of transport to create.

## Structure

- **Product**: Defines the interface for objects the factory method creates (Transport)
- **ConcreteProduct**: Implements the Product interface (Truck, Ship)
- **Creator**: Declares the factory method that returns a Product object (LogisticsService)
- **ConcreteCreator**: Overrides the factory method to return a ConcreteProduct (RoadLogistics, SeaLogistics)

## When to Use

- When you don't know ahead of time what concrete classes you will need
- When you want to provide users of your library or framework a way to extend its internal components
- When you want to save system resources by reusing existing objects instead of rebuilding them each time

## Benefits

- Eliminates the need to bind application-specific classes into your code
- Provides hooks for subclasses to extend a core component
- Connects parallel class hierarchies
