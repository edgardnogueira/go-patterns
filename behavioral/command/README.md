# Command Pattern

## Overview

The Command pattern is a behavioral design pattern that turns a request into a stand-alone object that contains all information about the request. This transformation allows you to:

- Parameterize methods with different requests
- Queue or log requests
- Support undoable operations

The pattern decouples the object that invokes the operation from the one that has the knowledge to perform it.

## Problem

Imagine you're developing a smart home application where users can control various home devices through a remote control, mobile app, or voice commands. Each device has different operations and interfaces. How do you:

1. Create a unified interface for controlling diverse devices?
2. Allow for operations to be undone?
3. Support scheduling operations for later execution?
4. Create composite operations (macros) that execute multiple commands at once?

## Solution

The Command pattern addresses these issues by encapsulating a request as an object, thereby allowing you to:

1. Define a common interface for all commands
2. Parameterize objects with commands
3. Store command history for undoing operations
4. Compose commands into larger commands
5. Queue or schedule command execution

## Structure

- **Command**: The interface that declares methods for executing and undoing operations.
- **Concrete Command**: Classes implementing the Command interface. They encapsulate a specific action and its parameters.
- **Invoker**: The class that requests the command execution (like a remote control).
- **Receiver**: The class that performs the actual work (like a light, thermostat, etc.).
- **Client**: Creates and configures command objects.

## Implementation in Go

### Command Interface

```go
type Command interface {
    Execute() error
    Undo() error
    String() string
}
```

### Concrete Commands

Various command implementations that control different devices, for example:

```go
type LightOnCommand struct {
    light *Light
}

func (c *LightOnCommand) Execute() error {
    c.light.TurnOn()
    return nil
}

func (c *LightOnCommand) Undo() error {
    c.light.TurnOff()
    return nil
}
```

### Invoker (RemoteControl)

```go
type RemoteControl struct {
    onCommands  []Command
    offCommands []Command
    history     []Command
}
```

### Receivers (Devices)

Various device classes:

```go
type Light struct {
    name       string
    isOn       bool
    brightness int
}

type Thermostat struct {
    name        string
    temperature int
    isOn        bool
    mode        string
}
```

## When to Use

- When you want to parameterize objects with operations
- When you need operations to be queued, executed at different times, or undone
- When you want to implement callbacks, request logging, or transaction systems
- When you need to structure a system around high-level operations built on primitives

## Benefits

1. **Single Responsibility Principle**: Command objects encapsulate specific operations
2. **Open/Closed Principle**: You can add new commands without changing existing code
3. **Undo/Redo Support**: Command history enables operation reversal
4. **Macro Commands**: Commands can be combined for complex operations
5. **Deferred Execution**: Commands can be scheduled for later execution

## Example Use Cases

1. **Remote Controls**: For TVs, home automation systems
2. **GUI Buttons and Menu Items**: Each button press executes a command
3. **Multi-level Undo/Redo**: Text editors, graphics applications
4. **Transaction Processing**: Financial systems where operations must be committed or rolled back
5. **Task Scheduling**: Running tasks at specific times
6. **Wizards and Multi-step Processes**: Breaking complex operations into steps

## Related Patterns

- **Composite Pattern**: Used with Command to create macro commands
- **Memento Pattern**: Can store command history with states for undo/redo
- **Observer Pattern**: Commands can be subscribers to events
- **Strategy Pattern**: Both encapsulate behavior, but Command focuses on complete operations with undo support

## Implementation Details

In this implementation, we've created a smart home automation system where commands control various devices:

1. **Devices**: Light, Thermostat, AudioSystem, GarageDoor
2. **Commands**: LightOnCommand, LightOffCommand, ThermostatSetCommand, etc.
3. **MacroCommand**: Executes multiple commands in sequence
4. **RemoteControl**: Invokes commands and maintains history for undo
5. **CommandQueue**: Supports scheduling commands for future execution

See the example directory for a demonstration of how to use this pattern in a complete application.
