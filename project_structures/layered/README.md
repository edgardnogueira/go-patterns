# Layered Architecture in Go

This directory demonstrates a layered architecture implementation in Go, providing a structured approach to organizing code in well-defined layers.

## Overview

Layered architecture is one of the most common architectural patterns, dividing the application into horizontal layers, where each layer has a specific responsibility and communicates only with adjacent layers following the dependency rule.

### Key Principles

1. **Separation of Concerns**: Each layer has a distinct responsibility
2. **Dependency Rule**: Dependencies flow in one direction - upper layers depend on lower layers, not vice versa
3. **Abstraction**: Each layer hides its implementation details from other layers
4. **Isolation**: Changes in one layer should have minimal impact on other layers

## Architecture Diagram

```
┌─────────────────────────────────────────────────────┐
│                 Presentation Layer                  │
│           (API Controllers, Worker Handlers)        │
└─────────────────────────────────────────────────────┘
                          ▼
┌─────────────────────────────────────────────────────┐
│                 Application Layer                   │
│        (Services, Use Cases, Orchestration)         │
└─────────────────────────────────────────────────────┘
                          ▼
┌─────────────────────────────────────────────────────┐
│                   Domain Layer                      │
│            (Business Logic, Entities)               │
└─────────────────────────────────────────────────────┘
                          ▼
┌─────────────────────────────────────────────────────┐
│                   Data Layer                        │
│             (Repositories, ORM, DAOs)               │
└─────────────────────────────────────────────────────┘
                          ▼
┌─────────────────────────────────────────────────────┐
│               Infrastructure Layer                  │
│        (Configuration, External Services)           │
└─────────────────────────────────────────────────────┘
```

## Layer Responsibilities

### Presentation Layer
- Handles HTTP requests and responses
- Processes user input
- Converts DTOs to/from domain models
- Manages API validation
- Implements API endpoints and worker handlers

### Application Layer
- Orchestrates use cases and business processes
- Coordinates domain objects
- Manages transactions
- Performs input validation for business rules
- Does not contain business rules, but organizes their execution

### Domain Layer
- Contains business logic and rules
- Defines domain entities and value objects
- Implements core business operations
- Is independent of other layers
- Has no knowledge of persistence or presentation concerns

### Data Layer
- Handles data persistence
- Manages database operations
- Implements repositories
- Converts between domain models and database entities
- Abstracts the data storage mechanism

### Infrastructure Layer
- Provides technical capabilities to other layers
- Implements cross-cutting concerns (logging, configuration)
- Contains adapters for external services
- Manages third-party library integrations
- Handles technical tasks like caching, messaging, etc.

## Project Structure Example

Our project implements a simple blog system with the following components:

- API server for CRUD operations
- Worker for background processing (e.g., sending notifications)
- Layered architecture principles demonstration
- Provider pattern for external services
- Clear separation of concerns

## Design Decisions

### DTOs vs Domain Models
- Data Transfer Objects (DTOs) are used for data exchange between layers
- Domain models represent business entities with behavior
- DTOs help in decoupling layers and versioning APIs

### Dependency Injection
- Used to provide dependencies to components without tight coupling
- Promotes testability and maintainability

### Interface Segregation
- Small, focused interfaces are defined for components
- Helps in mocking dependencies for testing

### Error Handling
- Each layer handles errors appropriate to its level of abstraction
- Domain errors are translated to application/presentation level errors

## Testing Strategy

- Unit tests for each layer
- Mock dependencies for isolation
- Test business logic independently of infrastructure
- Integration tests for critical paths

## Usage

To run the API server:
```bash
cd cmd/api
go run main.go
```

To run the worker:
```bash
cd cmd/worker
go run main.go
```

To run tests:
```bash
go test ./...
```

## Sample Flow

1. Request comes in through an API endpoint
2. Presentation layer validates input
3. Application layer orchestrates the use case
4. Domain layer applies business rules
5. Data layer persists changes
6. Response flows back up through the layers
