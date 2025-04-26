# Hexagonal Architecture (Ports and Adapters)

## Overview

Hexagonal Architecture, also known as Ports and Adapters, is an architectural pattern that allows an application to be equally driven by users, programs, automated tests, or batch scripts, and to be developed and tested in isolation from its runtime devices and databases.

The main idea is to create a clear separation between:
- The core business logic (domain)
- The interfaces through which the application communicates with the outside world (ports)
- The implementations of those interfaces (adapters)

## Key Concepts

### 1. Domain

The domain is the core of your application. It:
- Contains business logic and rules
- Is technology-agnostic (no dependencies on frameworks, databases, etc.)
- Has no knowledge of the outside world
- Defines entities, value objects, and services that model the business domain

### 2. Ports

Ports are interfaces that define how the domain interacts with the outside world:
- **Driving Ports (Primary/Inbound)**: Define how external actors can use the domain
- **Driven Ports (Secondary/Outbound)**: Define what the domain needs from the outside world

### 3. Adapters

Adapters implement the ports interfaces:
- **Driving Adapters (Primary/Inbound)**: Transform external requests into calls to the domain (e.g., HTTP handlers, CLI, gRPC)
- **Driven Adapters (Secondary/Outbound)**: Transform domain requests into technology-specific implementation (e.g., database repositories, external API clients)

## Architecture Diagram

```
                                         ┌─────────────────────────────────┐
                                         │                                 │
                                         │          External World         │
                                         │                                 │
                                         └───────────────┬─────────────────┘
                                                         │
                                                         │
┌ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ┐             │              ┌ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─
                                         │              │              │
│    ┌─────────────────────────────┐     │              ▼              │    ┌─────────────────────────────┐    │
     │                             │     │   ┌─────────────────────┐        │                             │
│    │      Driving Adapters       │─────┼──▶│    Driving Ports    │────┐  │    │      Driven Adapters        │    │
     │  (HTTP, CLI, gRPC, Worker)  │     │   │     (Interfaces)    │    │        │  (DB, External Services)   │
│    │                             │     │   └─────────────────────┘    │  │    │                             │    │
     └─────────────────────────────┘     │                              │        └─────────────────────────────┘
│                                         │           ┌─────┐            │  │                                   │
                                         │           │     │            │
│          Primary Adapters              │           │  D  │            │  │         Secondary Adapters         │
                                         │           │  O  │            │
│                                         │           │  M  │            │  │                                   │
     ┌─────────────────────────────┐     │           │  A  │            │        ┌─────────────────────────────┐
│    │                             │     │           │  I  │            │  │    │                             │    │
     │      Driving Adapters       │     │           │  N  │            │        │      Driven Adapters        │
│    │      (Tests, Scripts)       │─────┼──▶│       │     │            └──┼─▶│  │  (In-Memory Repositories)   │    │
     │                             │     │           └─────┘               │    │                             │
│    └─────────────────────────────┘     │                                 │    └─────────────────────────────┘    │
                                         │                                 │
└ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ┘                                 └ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─
```

## Benefits of Hexagonal Architecture

1. **Independent of external systems**: The domain doesn't depend on external technologies
2. **Highly testable**: Components can be tested in isolation
3. **Flexible**: Easy to swap implementations (e.g., switching database providers)
4. **Future-proof**: Core business logic is isolated from technological changes
5. **Maintainable**: Clear separation of concerns

## Implementation in Go

This example implements a simple order management system following Hexagonal Architecture principles, allowing you to:

1. Create new orders through HTTP API
2. Process orders asynchronously through a worker
3. Store orders in a repository (with both in-memory and simulated database implementations)
4. Notify external systems about order status changes

The project includes:
- A simple API server
- An asynchronous worker
- Domain models and services
- Clear separation between ports and adapters
- Unit tests demonstrating how components can be tested in isolation

## Running the Example

### API Server

```bash
cd cmd/api
go run main.go
```

### Worker

```bash
cd cmd/worker
go run main.go
```

### Tests

```bash
cd project_structures/hexagonal
go test ./...
```

## Key Files and Directories

- `cmd/api/main.go`: HTTP API entrypoint
- `cmd/worker/main.go`: Asynchronous worker entrypoint
- `internal/domain/model`: Domain entities and value objects
- `internal/domain/service`: Business logic services
- `internal/ports/driving`: Inbound interfaces
- `internal/ports/driven`: Outbound interfaces
- `internal/adapters/driving/http`: HTTP API implementation
- `internal/adapters/driving/worker`: Worker implementation
- `internal/adapters/driven/database`: Database adapters
- `internal/adapters/driven/external`: External service adapters

## Further Reading

- [Hexagonal Architecture by Alistair Cockburn](https://alistair.cockburn.us/hexagonal-architecture/)
- [Ports and Adapters Pattern by Juan Manuel Garrido](https://jmgarridopaz.github.io/content/hexagonalarchitecture.html)
- [Ready for changes with Hexagonal Architecture by Netflix](https://netflixtechblog.com/ready-for-changes-with-hexagonal-architecture-b315ec967749)
