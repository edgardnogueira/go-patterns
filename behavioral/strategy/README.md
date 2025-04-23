# Strategy Pattern

## Intent

The Strategy pattern defines a family of algorithms, encapsulates each one, and makes them interchangeable. It lets the algorithm vary independently from clients that use it.

## Explanation

This implementation demonstrates a payment processing system that can handle different payment methods (credit card, PayPal, cryptocurrency) as interchangeable strategies. The context (checkout process) doesn't need to know the details of how each payment method works.

## Structure

- **Strategy**: Interface that declares operations common to all supported algorithms (PaymentStrategy)
- **ConcreteStrategy**: Classes that implement the Strategy interface with specific algorithms (CreditCardStrategy, PayPalStrategy, CryptoStrategy)
- **Context**: Class that maintains a reference to a Strategy object and defines an interface that lets the strategy access its data (ShoppingCart)

## When to Use

- When you want to define a class that will have one behavior that is similar to other behaviors in a list
- When you need different variants of an algorithm
- When an algorithm uses data that clients shouldn't know about
- When a class defines many behaviors, and these appear as multiple conditional statements in its operations

## Benefits

- Isolates the implementation details of an algorithm from the code that uses it
- Helps avoid conditional logic for selecting desired behavior
- Provides an alternative to subclassing for changing behavior
- Enables runtime switching between different algorithms
