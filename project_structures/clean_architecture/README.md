# Clean Architecture in Go

This directory demonstrates a Clean Architecture implementation in Go, showcasing how to structure a Go project according to Clean Architecture principles established by Robert C. Martin (Uncle Bob).

## What is Clean Architecture?

Clean Architecture is a software design philosophy that separates the concerns of a software system into layers. The main goal is to make the system:

- Independent of frameworks
- Testable
- Independent of the UI
- Independent of the database
- Independent of any external agency

The fundamental rule is that dependencies can only point inward. Inner layers should not know anything about outer layers.

## Architecture Diagram

```
┌──────────────────────────────────────────────────────────┐
│                     External World                       │
└───────────────────────────┬──────────────────────────────┘
                            │
┌───────────────────────────▼──────────────────────────────┐
│              Frameworks & Drivers Layer                  │
│                                                          │
│   ┌─────────────────┐         ┌───────────────────────┐  │
│   │   HTTP Server   │         │  Database Drivers     │  │
│   └─────────────────┘         └───────────────────────┘  │
└───────────────────────────┬──────────────────────────────┘
                            │
┌───────────────────────────▼──────────────────────────────┐
│                Interface Adapters Layer                  │
│                                                          │
│   ┌─────────────────┐         ┌───────────────────────┐  │
│   │   Controllers   │         │    Repositories       │  │
│   └─────────────────┘         └───────────────────────┘  │
└───────────────────────────┬──────────────────────────────┘
                            │
┌───────────────────────────▼──────────────────────────────┐
│             Use Cases / Application Layer                │
│                                                          │
│   ┌─────────────────────────────────────────────────┐    │
│   │         Business Rules Orchestration            │    │
│   └─────────────────────────────────────────────────┘    │
└───────────────────────────┬──────────────────────────────┘
                            │
┌───────────────────────────▼──────────────────────────────┐
│                 Entities / Domain Layer                  │
│                                                          │
│   ┌─────────────────────────────────────────────────┐    │
│   │     Enterprise Business Rules and Entities      │    │
│   └─────────────────────────────────────────────────┘    │
└──────────────────────────────────────────────────────────┘
```

## Layers Explained

### 1. Entities / Domain Layer
- Core business entities and rules
- Independent of other layers
- Contains business models and core validation rules
- Has no dependencies on outer layers or external frameworks

### 2. Use Cases / Application Layer
- Contains application-specific business rules
- Implements and orchestrates the use cases of the system
- Depends only on the domain layer
- Uses domain entities but doesn't alter them

### 3. Interface Adapters Layer
- Converts data between the use cases and external layers
- Includes controllers, presenters, gateways, repositories
- Adapts external technologies to internal needs
- Depends on use cases and domain layers

### 4. Frameworks & Drivers Layer
- Outermost layer with frameworks and tools
- Database, web frameworks, devices, UI, etc.
- Glues the system to external tools
- Contains very little code as it mainly delegates to inner layers

## Flow of Control

In a Clean Architecture system, the flow of control typically follows this path:

1. External input comes through the Frameworks & Drivers layer (e.g., HTTP request)
2. The Interface Adapters layer converts and routes this input to the appropriate Use Case
3. The Use Case orchestrates the fulfillment of the user's request using Domain Entities
4. Results flow back up through the layers, with each layer transforming the data as needed

## Benefits of Clean Architecture

- **Independence from external frameworks**: Frameworks become tools rather than having your system constrained by their limitations.
- **Testability**: Business rules can be tested without UI, database, web server, or any external element.
- **Independence from UI**: The UI can change without affecting the rest of the system.
- **Independence from Database**: Business rules don't depend on a specific database, making it easy to switch databases.
- **Independence from any external agency**: Business rules don't know anything about the outside world.

## Project Structure

Our implementation follows this structure:

```
clean_architecture/
├── cmd/
│   ├── api/        # HTTP API entry point
│   │   └── main.go
│   └── worker/     # Background worker entry point
│       └── main.go
├── internal/
│   ├── domain/     # Entities Layer
│   │   ├── entities/
│   │   └── repositories/
│   ├── usecase/    # Use Cases Layer
│   │   └── service.go
│   ├── delivery/   # Interface Layer
│   │   ├── http/
│   │   └── worker/
│   └── infrastructure/ # Frameworks & Drivers Layer
│       ├── db/
│       └── providers/
├── README.md
└── Makefile
```

## Example Flow

This example implements a simple user management system with CRUD operations. The typical flow for creating a user works as follows:

1. An HTTP request comes in through the API server in `cmd/api`
2. The request is handled by a controller in `delivery/http` 
3. The controller transforms the HTTP request into a use case input
4. The use case in `usecase` processes the request, applying business rules
5. Domain entities in `domain/entities` are manipulated according to business rules
6. The repository interface in `domain/repositories` defines how entities are persisted
7. The concrete repository implementation in `infrastructure/db` handles actual storage
8. Results flow back up through the layers and are returned to the client

Additionally, some operations may be handled asynchronously by the worker:

1. A use case might dispatch a task to be processed later
2. The worker in `cmd/worker` picks up the task 
3. The worker handler in `delivery/worker` processes the task, using the same use cases and domain logic

## Running the Example

To run the API server:

```bash
make run-api
```

To run the worker:

```bash
make run-worker
```

To run tests:

```bash
make test
```
