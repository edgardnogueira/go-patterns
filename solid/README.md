# SOLID Principles in Go

This directory contains examples and explanations of SOLID principles implemented in Go.

## Overview

SOLID is an acronym for five design principles intended to make software designs more understandable, flexible, and maintainable:

- **S**ingle Responsibility Principle
- **O**pen/Closed Principle
- **L**iskov Substitution Principle
- **I**nterface Segregation Principle
- **D**ependency Inversion Principle

## Directory Structure

```
└── solid/
    ├── single_responsibility/   # Single Responsibility Principle
    ├── open_closed/            # Open/Closed Principle
    ├── liskov_substitution/    # Liskov Substitution Principle
    ├── interface_segregation/  # Interface Segregation Principle
    └── dependency_inversion/   # Dependency Inversion Principle
```

## Usage

Each subdirectory contains:

1. Implementation of the principle in Go
2. Documentation explaining the principle
3. Real-world examples
4. Tests demonstrating principle application

To run the examples, navigate to the specific principle directory and run:

```bash
go run main.go
```

To run tests:

```bash
go test ./...
```
