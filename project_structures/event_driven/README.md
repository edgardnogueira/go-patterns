# Event-Driven Architecture in Go

This directory demonstrates an implementation of Event-Driven Architecture (EDA) in Go, showcasing how to structure a Go project using event-driven principles and patterns.

## What is Event-Driven Architecture?

Event-Driven Architecture is an architectural pattern that promotes the production, detection, consumption of, and reaction to events. An event represents a significant change in state or an action that has happened within a system.

Key characteristics of Event-Driven Architecture include:
- **Loosely Coupled Components**: Components communicate through events without direct knowledge of each other
- **Asynchronous Processing**: Events are typically processed asynchronously
- **Scalability**: Easy to scale by adding more consumers to process events
- **Resilience**: Failure in one component doesn't necessarily affect others
- **Event History**: Events can be stored as an immutable log (Event Sourcing)

## Architecture Diagram

```
┌─────────────────────────────────────────────────┐
│                 External World                   │
└────────────────────┬────────────────────────────┘
                     │
┌────────────────────┼────────────────────────────┐
│                API Layer                         │
│  ┌─────────────────┴──────────────────┐         │
│  │            API Gateway              │         │
│  └─────────────────┬──────────────────┘         │
│                    │                             │
│  ┌─────────────────┴──────────────────┐         │
│  │          Command Handlers           │         │
│  └─────────────────┬──────────────────┘         │
└────────────────────┼────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────┐
│                                                  │
│              Command Bus / Dispatcher            │
│                                                  │
└────────────────────┬────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────┐
│                Domain Layer                      │
│                                                  │
│  ┌─────────────────────────────────────┐        │
│  │   Command Processing & Validation    │        │
│  └───────────────────┬─────────────────┘        │
│                      │                           │
│  ┌───────────────────▼─────────────────┐        │
│  │         Aggregate Operations         │        │
│  └───────────────────┬─────────────────┘        │
│                      │                           │
│  ┌───────────────────▼─────────────────┐        │
│  │           Event Creation             │        │
│  └───────────────────┬─────────────────┘        │
└──────────────────────┼───────────────────────────┘
                       │
┌─────────────────────▼───────────────────────────┐
│             Event Publication                    │
│                                                  │
│  ┌────────────────────────────────────┐         │
│  │           Event Store              │         │
│  └───────────────┬────────────────────┘         │
│                  │                              │
│  ┌───────────────▼────────────────────┐         │
│  │         Message Broker              │         │
│  └───────────────┬────────────────────┘         │
└──────────────────┼──────────────────────────────┘
                   │
┌──────────────────┼──────────────────────────────┐
│              Event Consumption                   │
│  ┌───────────────▼────────────────────┐         │
│  │         Event Consumers            │         │
│  └───────────────┬────────────────────┘         │
│                  │                              │
│  ┌───────────────▼────────────────────┐         │
│  │    Projections / Read Models        │         │
│  └────────────────────────────────────┘         │
│                                                  │
│  ┌────────────────────────────────────┐         │
│  │     Asynchronous Processes         │         │
│  └────────────────────────────────────┘         │
│                                                  │
│  ┌────────────────────────────────────┐         │
│  │      Integration / Notifications    │         │
│  └────────────────────────────────────┘         │
└──────────────────────────────────────────────────┘
```

## Key Components

### 1. Command Processing
- **Commands**: Represent intent to change the system state
- **Command Handlers**: Process commands and validate business rules
- **Command Bus**: Routes commands to appropriate handlers

### 2. Event Production & Storage
- **Events**: Immutable records of state changes that have occurred
- **Event Store**: Persistent storage for events
- **Event Publisher**: Publishes events to message brokers

### 3. Event Consumption
- **Event Consumers**: Listen for and react to events
- **Projections**: Build read-optimized views from events
- **Workers**: Process events asynchronously

### 4. Infrastructure
- **Message Brokers**: Distribute events to consumers (Kafka, RabbitMQ)
- **Databases**: Store events and projections (Event Store, SQL, NoSQL)

## Implemented Patterns

### Command Query Responsibility Segregation (CQRS)
Separates the write operations (commands) from read operations (queries):
- Commands modify state and produce events
- Queries read from optimized projections/views

### Event Sourcing (ES)
Stores all changes to application state as a sequence of events:
- The event store is the system of record
- Application state can be reconstructed by replaying events
- Provides a complete audit trail and time-travel capabilities

### Saga Pattern
Manages distributed transactions across multiple services:
- Each step in a transaction produces events and listens for others
- If a step fails, compensating transactions are triggered

### Outbox Pattern
Ensures consistency between database operations and message publishing:
- Events are stored in an "outbox" table within the same transaction as state changes
- A separate process reliably processes the outbox and publishes events

### Idempotent Consumers
Ensures events can be processed multiple times without side effects:
- Consumers track already processed events
- Duplicate events are detected and ignored

## Example Flow

This example implements a simplified inventory management system that demonstrates:

1. How a user initiates a product stock reduction (command)
2. How the command is validated and processed
3. How events are generated and stored
4. How events are published to a message broker
5. How consumers process events to update projections
6. How to handle failures and implement retries
7. How to ensure consistency across distributed components

### Example Use Case: Reducing Product Stock

1. API receives a request to reduce product stock
2. Command is created and sent to the command bus
3. Command handler validates business rules (sufficient stock available)
4. Domain logic processes the command, reducing stock in the aggregate
5. Events are generated (StockReducedEvent)
6. Events are stored in the event store
7. Events are published to the message broker
8. Multiple consumers process the event:
   - Projection updater maintains the current stock level view
   - Notification service informs about low stock
   - Analytics service updates metrics
9. If any consumer fails, retry mechanisms ensure eventual processing

## Project Structure

```
project_structures/
└── event_driven/
    ├── cmd/
    │   ├── api/
    │   │   └── main.go
    │   ├── publisher/
    │   │   └── main.go
    │   └── worker/
    │       └── main.go
    ├── internal/
    │   ├── commands/
    │   │   ├── handler/
    │   │   └── model/
    │   ├── events/
    │   │   ├── publisher/
    │   │   ├── consumer/
    │   │   ├── store/
    │   │   └── types/
    │   ├── queries/
    │   │   ├── handler/
    │   │   └── model/
    │   ├── projections/
    │   │   ├── builder/
    │   │   └── repository/
    │   └── infrastructure/
    │       ├── messaging/
    │       │   ├── kafka/
    │       │   └── rabbitmq/
    │       ├── eventstore/
    │       │   ├── postgres/
    │       │   └── mongo/
    │       └── config/
    ├── pkg/
    │   ├── cqrs/
    │   ├── eventsourcing/
    │   └── common/
    ├── README.md
    └── Makefile
```

## Running the Example

To run the API server:

```bash
make run-api
```

To run the event publisher:

```bash
make run-publisher
```

To run the worker:

```bash
make run-worker
```

To run tests:

```bash
make test
```

## Benefits of Event-Driven Architecture

- **Scalability**: Components can be scaled independently based on load
- **Loose Coupling**: Services can evolve independently
- **Resilience**: Failure in one component doesn't bring down the entire system
- **Real-time Processing**: Events can be processed as they occur
- **Audit Trail**: Complete history of all state changes
- **Time Travel**: Ability to reconstruct system state at any point in time
- **Extensibility**: New consumers can be added without modifying existing components

## Challenges and Considerations

- **Eventual Consistency**: Data consistency takes time to propagate
- **Complexity**: More moving parts to manage and debug
- **Monitoring**: Requires comprehensive monitoring across distributed components
- **Error Handling**: Must handle failures in event processing
- **Ordering**: May need to preserve event order in some scenarios
- **Event Evolution**: Versioning and backward compatibility of events