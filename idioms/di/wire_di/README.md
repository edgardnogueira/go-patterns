# Wire-based Dependency Injection

This directory demonstrates using Google's Wire tool for compile-time dependency injection in Go.

## What is Wire?

[Wire](https://github.com/google/wire) is a code generation tool from Google that automates dependency injection in Go programs. Unlike traditional DI containers that work at runtime, Wire operates at compile time, generating code that directly instantiates and connects components.

## Key Concepts

1. **Providers**: Functions that produce values
2. **Injectors**: Generated functions that call providers in the correct order
3. **Wire Directives**: Special build tags and function calls that tell the Wire tool what to generate

## Benefits

- **Compile-time validation**: Dependency issues are detected at compile time
- **No runtime reflection**: The generated code is simply direct function calls
- **No dependency on a container**: The final binary doesn't include any DI framework
- **Clear dependency graph**: The generated code makes dependencies explicit
- **No magic**: The generated code is readable and follows the same patterns you would write by hand

## How to Use Wire

### Installation

To install the Wire tool:

```bash
go install github.com/google/wire/cmd/wire@latest
```

### Basic Usage

1. Define component structs and their constructor functions (providers)
2. Create a Wire file with provider function sets and an injector function signature
3. Run the `wire` command to generate the implementation of the injector function

### Running the Example

The example in this directory shows how to use Wire to set up a service with various dependencies:

1. View the model and service files to understand the application structure
2. Examine the `wire.go` file to see how providers are organized
3. Look at the `wire_gen.go` file to see the generated code
4. Run the example with `go run .`

### Generated Code

The Wire tool generates code based on your provider functions and injector signatures. The generated code:

1. Creates instances of your dependencies in the correct order
2. Passes them to the appropriate constructors
3. Handles error propagation

## Use Cases for Wire

Wire is especially useful for:

- Large applications with complex dependency graphs
- Applications that need compile-time validation of dependencies
- Projects with clear dependency hierarchies
- Teams that prefer explicit dependency connections over "magical" DI frameworks

## Benefits of Wire vs. Hand-written Dependency Injection

Compared to manually writing dependency injection code:

- **Less boilerplate**: Wire generates the wiring code for you
- **Easier maintenance**: Adding a new dependency only requires adding a provider, not changing existing code
- **Better code organization**: Dependencies are defined near their components

Compared to runtime DI containers:

- **Better performance**: No reflection or lookups at runtime
- **Smaller binary size**: No framework code included
- **Easier debugging**: The generated code is straightforward function calls
