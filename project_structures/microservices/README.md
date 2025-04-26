# Microservices Architecture in Go

This directory contains an example project structure demonstrating how to organize code for a microservices architecture in Go.

## Overview

Microservices architecture is an approach to application development where a large application is built as a suite of small, independently deployable services. Each service runs in its own process and communicates with other services through well-defined APIs.

### Key Characteristics

- **Decentralized**: Each service can be developed, deployed, and scaled independently
- **Domain-Driven**: Services are organized around business capabilities
- **Resilient**: Failure in one service doesn't bring down the entire system
- **Scalable**: Individual services can be scaled based on demand
- **Technology Diversity**: Different services can use different technologies when appropriate

## Project Structure

This example demonstrates a simple e-commerce system with the following components:

- **API Gateway**: Entry point for external clients, handling routing and authentication
- **Order Service**: Manages order creation and processing
- **Inventory Service**: Handles product inventory and availability
- **Notification Service**: Sends notifications to users
- **Worker**: Processes background tasks and scheduled jobs

```
microservices/
├── api-gateway/            # API Gateway service
│   ├── main.go             # Entry point
│   ├── routers/            # Route definitions
│   └── middleware/         # Authentication, logging middleware
├── order-service/          # Order management service
│   ├── cmd/
│   │   ├── api/            # HTTP API
│   │   │   └── main.go
│   │   └── worker/         # Background worker
│   │       └── main.go
│   └── internal/
│       ├── domain/         # Domain models and business logic
│       ├── handlers/       # Request handlers
│       ├── repositories/   # Data access layer
│       └── providers/      # External service clients
├── inventory-service/      # Inventory management service
│   ├── cmd/
│   │   ├── api/
│   │   │   └── main.go
│   │   └── worker/
│   │       └── main.go
│   └── internal/
│       ├── domain/
│       ├── handlers/
│       ├── repositories/
│       └── providers/
├── pkg/                    # Shared packages
│   ├── messaging/          # Message broker integration
│   ├── observability/      # Logging, metrics, tracing
│   └── common/             # Shared utilities
├── Makefile                # Build and deployment commands
├── docker-compose.yml      # Local development environment
└── README.md               # Project documentation
```

## Communication Patterns

### Synchronous Communication

Services communicate synchronously via HTTP/REST or gRPC for operations that require immediate responses:

- **HTTP/REST**: Simple to implement and understand, good for public APIs
- **gRPC**: Efficient binary protocol with strict typing, ideal for internal service communication

### Asynchronous Communication

Services communicate asynchronously via message queues for operations that can be processed later:

- **Event-Driven**: Services publish events when state changes
- **Command Queue**: Services send commands for other services to process
- **Message Broker**: RabbitMQ, Kafka, or NATS used as the messaging infrastructure

## Flow Diagram

The following diagram illustrates the communication flow between services:

```
┌─────────────┐      ┌───────────────┐      ┌─────────────────┐
│             │ HTTP │               │ HTTP │                 │
│  Client     ├─────►│  API Gateway  ├─────►│  Order Service  │
│             │      │               │      │                 │
└─────────────┘      └───────────────┘      └────────┬────────┘
                                                     │
                                                     │ HTTP/gRPC
                                                     ▼
┌─────────────┐      ┌───────────────┐      ┌─────────────────┐
│             │      │               │◄─────┤                 │
│  Worker     │◄─────┤  Message Bus  │      │ Inventory Serv. │
│             │      │               │      │                 │
└─────────────┘      └───────────────┘      └─────────────────┘
      △                      △                      │
      │                      │                      │
      │                      │                      │
      │                      │                      │
      └──────────────────────┼──────────────────────┘
                             │
                      ┌──────┴────────┐
                      │               │
                      │ Notification  │
                      │   Service     │
                      │               │
                      └───────────────┘
```

## Example Flow: Order Processing

1. Client sends an order request to the API Gateway
2. API Gateway authenticates the request and forwards it to the Order Service
3. Order Service validates the order and calls the Inventory Service to check stock
4. If stock is available, Order Service creates the order and publishes an "OrderCreated" event
5. Inventory Service consumes the "OrderCreated" event and reserves inventory
6. Notification Service consumes the "OrderCreated" event and sends a confirmation email
7. Worker processes the order for fulfillment

## Observability

Distributed systems require robust observability to understand system behavior:

- **Structured Logging**: Consistent log format across services with correlation IDs
- **Metrics**: Monitor service health, performance, and business metrics
- **Distributed Tracing**: Track requests as they flow through multiple services

## Running the Example

```bash
# Start all services locally
make up

# Run tests
make test

# Stop all services
make down
```

## Additional Features

- **Service Discovery**: Services register themselves and discover other services
- **Circuit Breakers**: Prevent cascading failures when services are unavailable
- **Distributed Configuration**: Centralized configuration management
- **API Documentation**: OpenAPI/Swagger specifications for each service
