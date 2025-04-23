# Abstract Factory Pattern

## Intent

The Abstract Factory pattern provides an interface for creating families of related or dependent objects without specifying their concrete classes.

## Explanation

This implementation demonstrates a GUI toolkit that can create UI components for different operating systems (Modern and Vintage). The Abstract Factory provides interfaces to create different products (buttons and checkboxes) that belong to the same family.

## Structure

- **AbstractFactory**: Interface that declares creation methods for abstract products (GUIFactory)
- **ConcreteFactory**: Implements the creation methods of the AbstractFactory (ModernGUIFactory, VintageGUIFactory)
- **AbstractProduct**: Interface for a type of product object (Button, Checkbox)
- **ConcreteProduct**: Implements the AbstractProduct interface (ModernButton, VintageButton, ModernCheckbox, VintageCheckbox)
- **Client**: Works with factories and products through abstract interfaces

## When to Use

- When a system should be independent of how its products are created, composed, and represented
- When a system should be configured with one of multiple families of products
- When a family of related product objects is designed to be used together
- When you want to provide a class library of products, and you want to reveal just their interfaces, not their implementations

## Benefits

- Isolates concrete classes from the client
- Makes exchanging product families easy
- Promotes consistency among products
- Supporting new kinds of products is difficult (requires changing the AbstractFactory interface)
