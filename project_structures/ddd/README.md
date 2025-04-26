# Domain-Driven Design (DDD) in Go

This directory provides an example implementation of a project structure using Domain-Driven Design (DDD) principles in Go.

## Overview

Domain-Driven Design (DDD) is an approach to software development that focuses on understanding the problem domain and modeling it effectively in code. DDD emphasizes:

- Collaboration between domain experts and developers
- Creating a ubiquitous language shared by all team members
- Modeling complex domains through strategic and tactical design patterns
- Separating the domain from implementation details

This example demonstrates how to implement DDD principles in a Go application, showing how to organize code, define boundaries, and implement domain logic.

## Structure

The example implements an e-commerce system with the following bounded contexts:
- **Order Management**: Handles the creation and processing of customer orders
- **Inventory Management**: Manages product inventory and stock levels
- **User Management**: Handles user accounts and authentication

### Project Layout

```
project_structures/ddd/
├── cmd/                              # Application entry points
│   ├── api/                          # API server
│   │   └── main.go
│   └── worker/                       # Background worker for processing events
│       └── main.go
├── internal/                         # Private application code
│   ├── order/                        # Order Management bounded context
│   │   ├── domain/                   # Domain layer
│   │   │   ├── aggregate/            # Aggregate roots
│   │   │   ├── entity/               # Entities
│   │   │   ├── valueobject/          # Value objects
│   │   │   ├── event/                # Domain events
│   │   │   ├── repository/           # Repository interfaces
│   │   │   └── service/              # Domain services
│   │   ├── application/              # Application layer
│   │   │   ├── service/              # Application services
│   │   │   ├── dto/                  # Data Transfer Objects
│   │   │   └── mapper/               # DTO-to-Domain mappers
│   │   └── infrastructure/           # Infrastructure layer
│   │       ├── repository/           # Repository implementations
│   │       ├── event/                # Event handling
│   │       └── provider/             # External service providers
│   ├── inventory/                    # Inventory Management bounded context
│   │   └── ... (similar structure)
│   ├── user/                         # User Management bounded context
│   │   └── ... (similar structure)
│   └── shared/                       # Shared code between contexts
│       ├── domain/                   # Shared domain types and interfaces
│       └── infrastructure/           # Shared infrastructure
├── pkg/                              # Public libraries that can be used by external applications
│   ├── eventbus/                     # Event bus implementation
│   └── common/                       # Common utilities
└── Makefile                          # Build automation
```

## Key DDD Concepts Implemented

### Strategic Design

- **Bounded Contexts**: Clear separation between different business domains
- **Context Mapping**: Defined relationships between contexts using events
- **Ubiquitous Language**: Domain terminology consistently used in code

### Tactical Design

- **Entities**: Objects with identity and lifecycle (e.g., Order, User)
- **Value Objects**: Immutable objects without identity (e.g., Address, Money)
- **Aggregates**: Clusters of domain objects treated as a unit (e.g., Order with OrderItems)
- **Domain Events**: Objects representing something that happened in the domain
- **Repositories**: Persistence abstractions for retrieving and storing aggregates
- **Domain Services**: Operations that don't belong to any specific entity
- **Application Services**: Orchestration of domain objects for use cases

## Example Flow

The example demonstrates:

1. Creating a new order (Order context)
2. Checking inventory availability (Inventory context)
3. Processing payment (Order context)
4. Updating inventory upon successful order (Inventory context, via domain events)
5. Notifying the customer (User context, via domain events)

## DDD Principles Applied

- **Rich Domain Model**: Business logic encapsulated in domain objects
- **Immutability**: Value objects are immutable
- **Encapsulation**: Internal state protected through proper encapsulation
- **Invariants**: Business rules enforced within aggregates
- **Persistence Ignorance**: Domain model unaware of storage mechanisms
- **Dependency Inversion**: Domain doesn't depend on infrastructure
- **CQRS Pattern**: Separation of read and write operations where appropriate

## Domain Model Diagram

```
┌─────────────────────────────┐      ┌──────────────────────────┐      ┌─────────────────────────┐
│      Order Context          │      │    Inventory Context     │      │      User Context       │
│                             │      │                          │      │                         │
│  ┌─────────┐  ┌─────────┐   │      │   ┌────────┐             │      │  ┌─────┐  ┌─────────┐  │
│  │  Order  │──│OrderItem│   │      │   │ Product│             │      │  │User │──│ Address │  │
│  └─────────┘  └─────────┘   │      │   └────────┘             │      │  └─────┘  └─────────┘  │
│       │                     │      │       │                  │      │     │                  │
│  ┌─────────┐                │      │  ┌────────┐              │      │  ┌─────────┐           │
│  │ Payment │                │      │  │ Stock  │              │      │  │ Profile │           │
│  └─────────┘                │      │  └────────┘              │      │  └─────────┘           │
│                             │      │                          │      │                         │
└──────────┬──────────────────┘      └──────────┬───────────────┘      └─────────┬───────────────┘
           │                                    │                                │
           │                                    │                                │
           │                                    │                                │
           │          ┌───────────────────────────────────────┐                  │
           └──────────┤           Event Bus                   ├──────────────────┘
                      └───────────────────────────────────────┘
```

## Running the Example

To run the API server:
```bash
cd cmd/api
go run main.go
```

To run the worker that processes domain events:
```bash
cd cmd/worker
go run main.go
```

To run tests:
```bash
go test ./...
```

## Best Practices Demonstrated

- Clear separation between domain and infrastructure
- Use of interfaces for external dependencies
- Immutable value objects
- Domain events for communication between bounded contexts
- Validation at domain object creation
- Aggregate roots controlling access to their components
- Repository pattern for persistence abstraction
- Application services for orchestrating use cases
