# Go Design Patterns

A collection of design patterns implemented in Go with practical examples to help developers learn and apply design patterns effectively in their Go projects.

## Overview

This repository contains implementations of various design patterns in Go, categorized into:

- **Creational Patterns**: Patterns dealing with object creation mechanisms
- **Structural Patterns**: Patterns concerning object composition and relationships
- **Behavioral Patterns**: Patterns focusing on communication between objects
- **Go-specific Patterns**: Idiomatic Go patterns that leverage Go's unique features

## Project Structure

```
├── creational/          # Creational design patterns
│   ├── factory/         # Factory Method pattern
│   ├── abstract/        # Abstract Factory pattern
│   ├── builder/         # Builder pattern
│   ├── prototype/       # Prototype pattern
│   └── singleton/       # Singleton pattern
│
├── structural/          # Structural design patterns
│   ├── adapter/         # Adapter pattern
│   ├── bridge/          # Bridge pattern
│   ├── composite/       # Composite pattern
│   ├── decorator/       # Decorator pattern
│   ├── facade/          # Facade pattern
│   ├── flyweight/       # Flyweight pattern
│   └── proxy/           # Proxy pattern
│
├── behavioral/          # Behavioral design patterns
│   ├── chain/           # Chain of Responsibility pattern
│   ├── command/         # Command pattern
│   ├── iterator/        # Iterator pattern
│   ├── mediator/        # Mediator pattern
│   ├── memento/         # Memento pattern
│   ├── observer/        # Observer pattern
│   ├── state/           # State pattern
│   ├── strategy/        # Strategy pattern
│   ├── template/        # Template Method pattern
│   └── visitor/         # Visitor pattern
│
└── idioms/              # Go-specific idiomatic patterns
    ├── errors/          # Error handling patterns
    ├── interfaces/      # Interface implementation patterns
    ├── context/         # Context usage patterns
    ├── concurrency/     # Concurrency patterns
    ├── options/         # Functional options pattern
    └── builder/         # Go-styled builder patterns
```

## Usage

Each pattern includes:

1. Implementation of the pattern in Go
2. Documentation explaining the pattern
3. Real-world examples
4. Tests demonstrating pattern usage

To run the examples, navigate to the specific pattern directory and run:

```bash
go run main.go
```

To run tests:

```bash
go test ./...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-pattern`)
3. Commit your changes (`git commit -m 'Add some amazing pattern'`)
4. Push to the branch (`git push origin feature/amazing-pattern`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
