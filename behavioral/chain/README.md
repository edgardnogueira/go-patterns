# Chain of Responsibility Pattern

## Intent

The Chain of Responsibility pattern passes a request along a chain of handlers. Upon receiving a request, each handler decides either to process the request or to pass it to the next handler in the chain.

## Explanation

This implementation demonstrates a support ticket system where tickets are processed based on their priority and type through different support levels. The pattern decouples the sender of a request from its receivers, allowing multiple objects to handle the request without the sender needing to know which one will ultimately process it.

## Structure

- **Handler Interface**: Defines the interface for handling requests and passing them along the chain
- **BaseHandler**: Provides common functionality for all handlers
- **Concrete Handlers**:
  - **Level1Support**: Handles basic user inquiries and common issues
  - **Level2Support**: Handles technical issues that require more expertise
  - **Level3Support**: Handles complex issues that require system-level access
  - **ManagerSupport**: Handles escalated issues, customer complaints, and policy exceptions
  - **SecurityHandler**: Handles security-related issues with specialized expertise
  - **FallbackHandler**: Ensures all tickets receive a response even if not fully resolved
  - **LoggingHandler**: Logs all tickets passing through but doesn't resolve them
  - **PriorityUpgradeHandler**: Upgrades ticket priority based on keywords
- **Chain**: Manages the chain of responsibility by maintaining references to the handlers
- **Client**: Creates the chain and sends requests to it

## When to Use

- When more than one object may handle a request, and the handler isn't known a priori
- When you want to pass a request to one of several objects without specifying the receiver explicitly
- When the set of objects that can handle a request should be specified dynamically
- When you want to avoid coupling the sender of a request to its receiver

## Benefits

1. **Reduced coupling**: The sender doesn't need to know which object will handle the request
2. **Flexibility in assigning responsibilities**: You can add or remove handlers dynamically
3. **Simplified object connections**: Each handler only knows about its successor, not the entire chain
4. **Single Responsibility Principle**: Each handler focuses on a specific responsibility
5. **Open/Closed Principle**: You can add new handlers without changing existing code

## Implementation Details

In our implementation:

1. The `Handler` interface defines the methods for handling requests and passing them along the chain
2. The `BaseHandler` provides common functionality for all handlers
3. Concrete handlers implement specific handling logic for different types of requests
4. The `Chain` class manages the chain of responsibility, allowing dynamic modification
5. Each handler decides whether to process the request or pass it to the next handler
6. Special handlers like `LoggingHandler` and `PriorityUpgradeHandler` provide cross-cutting functionality
7. A `FallbackHandler` ensures all requests receive a response

## Example

```go
// Create handlers
level1 := chain.NewLevel1Support()
level2 := chain.NewLevel2Support()
level3 := chain.NewLevel3Support()
manager := chain.NewManagerSupport()

// Setup the chain
supportChain := chain.NewChain(level1)
supportChain.AddHandler(level2)
supportChain.AddHandler(level3)
supportChain.AddHandler(manager)

// Create a ticket
ticket := chain.NewSupportTicket("TKT-001", chain.Technical, chain.Medium, 
    "App Crash", "Application crashes when uploading files")

// Process the ticket through the chain
supportChain.Process(ticket)

// Check the result
if ticket.IsResolved {
    fmt.Printf("Ticket resolved by: %s\n", ticket.ResolvedBy)
    fmt.Printf("Resolution: %s\n", ticket.Resolution)
}
```

## Dynamic Chain Modification

Our implementation supports adding, removing, and inserting handlers at runtime:

```go
// Add a handler to the end of the chain
supportChain.AddHandler(newHandler)

// Insert a handler after a specific one
supportChain.InsertHandler("Level 1 Support", specializedHandler)

// Remove a handler
supportChain.RemoveHandler("Level 2 Support")
```

## Request Prioritization

The implementation supports prioritizing requests using the `PriorityUpgradeHandler`:

```go
// Create a priority handler with default keywords
priorityHandler := chain.NewPriorityUpgradeHandler()

// Add custom keywords that will upgrade priority
priorityHandler.AddKeyword("urgent", chain.Critical)
priorityHandler.AddKeyword("important", chain.High)

// Add it at the beginning of the chain
supportChain := chain.NewChain(priorityHandler)
supportChain.AddHandler(regularHandlers...)
```

## Related Patterns

- **Command**: Chain of Responsibility can be used with Command to implement a chain of command objects
- **Composite**: Chain of Responsibility often uses Composite to represent the chain
- **Decorator**: Chain of Responsibility and Decorator are both used to add behavior to an object dynamically
