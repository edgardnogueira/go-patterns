# Command Pattern Documentation

## Intent

The Command pattern transforms a request into a stand-alone object containing all information about the request. This transformation allows you to:
- Parameterize objects with different requests
- Queue or log requests
- Support undoable operations
- Compose simple commands into complex ones

## Participants

### Command
The Command interface defines operations that all concrete commands must implement:

```go
type Command interface {
    Execute() error
    Undo() error
    String() string
}
```

- **Execute()**: Performs the action associated with the command
- **Undo()**: Reverses the effects of executing the command
- **String()**: Returns a human-readable description of the command

### Concrete Commands
Concrete command classes implement the Command interface. Each concrete command operates on a receiver (device) and maintains state necessary for undoing.

Examples:
- `LightOnCommand`: Turns a light on
- `ThermostatSetCommand`: Sets the thermostat to a specific temperature
- `AudioPlayCommand`: Plays media on an audio system

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

### Invoker
The invoker asks the command to carry out the request. In our implementation, `RemoteControl` serves as the invoker.

```go
type RemoteControl struct {
    onCommands  []Command
    offCommands []Command
    history     []Command
    maxHistory  int
}
```

The invoker:
- Holds command objects
- Triggers command execution
- Optionally maintains command history for undo operations

### Receiver
The receiver knows how to perform the operations associated with carrying out a request. In our implementation, devices like `Light`, `Thermostat`, and `AudioSystem` are receivers.

```go
type Light struct {
    name       string
    isOn       bool
    brightness int
}

func (l *Light) TurnOn() {
    l.isOn = true
    fmt.Printf("%s light is now ON\n", l.name)
}
```

### Client
The client creates and configures concrete Command objects. It sets a receiver for the command and potentially registers the command with an invoker.

In our example application, the `main.go` creates and configures the commands.

## Structural Relationships

1. The **Client** creates a **Concrete Command** object and specifies its **Receiver**.
2. The **Invoker** stores the **Concrete Command** object.
3. The **Invoker** invokes the command by calling its Execute method.
4. The **Concrete Command** object invokes operations on its **Receiver** to carry out the request.

## Implementation Details

### Basic Command Implementation

A simple command stores a reference to its receiver and invokes methods on that receiver:

```go
type LightOffCommand struct {
    light *Light
}

func (c *LightOffCommand) Execute() error {
    c.light.TurnOff()
    return nil
}

func (c *LightOffCommand) Undo() error {
    c.light.TurnOn()
    return nil
}
```

### Command History and Undo

The RemoteControl maintains a history of executed commands, allowing for undo operations:

```go
// RemoteControl's Undo method
func (r *RemoteControl) Undo() error {
    if len(r.history) == 0 {
        return fmt.Errorf("no commands to undo")
    }
    
    lastIndex := len(r.history) - 1
    lastCommand := r.history[lastIndex]
    r.history = r.history[:lastIndex]
    
    return lastCommand.Undo()
}
```

### Composite Commands (Macros)

The MacroCommand implements the Command interface but contains multiple commands:

```go
type MacroCommand struct {
    commands []Command
    name     string
}

func (m *MacroCommand) Execute() error {
    for _, cmd := range m.commands {
        if err := cmd.Execute(); err != nil {
            return fmt.Errorf("macro command '%s' failed: %w", m.name, err)
        }
    }
    return nil
}

func (m *MacroCommand) Undo() error {
    // Undo commands in reverse order
    for i := len(m.commands) - 1; i >= 0; i-- {
        if err := m.commands[i].Undo(); err != nil {
            return fmt.Errorf("macro command '%s' undo failed: %w", m.name, err)
        }
    }
    return nil
}
```

### Command Queuing and Scheduling

The CommandQueue allows for scheduling commands to be executed at a later time:

```go
type CommandQueue struct {
    queue []QueuedCommand
}

type QueuedCommand struct {
    Command       Command
    ExecutionTime time.Time
}
```

## Example Use Cases

Our implementation demonstrates several common use cases for the Command pattern:

1. **Basic Device Control**: Simple commands to turn devices on or off
2. **Parameterized Commands**: Commands with parameters like temperature or brightness settings
3. **Command History**: Tracking executed commands for undo functionality
4. **Macro Commands**: Combining multiple commands into a single operation (e.g., "Evening Scene")
5. **Command Scheduling**: Queueing commands for future execution

## Benefits

1. **Decoupling**: Commands decouple objects that invoke operations from objects that perform these operations
2. **Extensibility**: You can add new commands without changing existing code
3. **Composite Commands**: You can create macro commands from simple commands
4. **Undo/Redo**: Commands can support undoing and redoing operations
5. **Deferred Execution**: Commands can be scheduled to execute later
6. **Queuing**: Commands can be queued for batch processing
7. **Logging**: Command execution can be logged for audit purposes
8. **Transactions**: Commands can be used to implement transactions

## Pitfalls and Considerations

1. **Command Proliferation**: For complex systems, you may end up with many command classes
2. **Undo Limitations**: Some actions cannot be undone or require complex state tracking
3. **Memory Usage**: Storing command history for undo/redo can consume significant memory
4. **Error Handling**: Commands must handle errors appropriately, especially in queues and macros
5. **Serialization**: If commands need to be persisted or transmitted, they must be serializable

## Related Patterns

- **Composite Pattern**: Used with Command to create macro commands
- **Memento Pattern**: Can be used with Command to store state for more complex undo operations
- **Observer Pattern**: Commands can notify observers when they execute
- **Strategy Pattern**: Both encapsulate behavior, but Command represents a complete operation while Strategy usually represents an algorithm
- **Factory Method Pattern**: Can be used to create commands

## Real-World Examples

The Command pattern is widely used in:

1. GUI frameworks (buttons, menu items)
2. Transaction processing systems
3. Text editors (for undo/redo)
4. Remote controls and IoT device management
5. Task schedulers and job queues
6. Gaming input systems
7. Multi-level undo systems

## Conclusion

The Command pattern provides a powerful way to encapsulate requests as objects, enabling features like parameterization, queueing, logging, and undo functionality. Our implementation demonstrates these capabilities in a smart home automation system, showing how diverse devices can be controlled through a unified interface.
