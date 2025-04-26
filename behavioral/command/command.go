// Package command implements the Command design pattern in Go.
//
// The Command pattern encapsulates a request as an object, thereby letting
// you parameterize clients with different requests, queue or log requests,
// and support undoable operations. This implementation demonstrates a
// smart home automation system where various commands can be issued to
// control different devices.
package command

import (
	"fmt"
	"time"
)

// Command is the interface that wraps the basic Execute and Undo methods.
// Any struct that implements these two methods can be used as a command
// in our system.
type Command interface {
	// Execute runs the command
	Execute() error
	
	// Undo reverses the effects of executing the command
	Undo() error
	
	// String returns a description of the command
	String() string
}

// NoOpCommand is a command that does nothing
// Useful as a null object pattern implementation
type NoOpCommand struct{}

// Execute implementation for NoOpCommand
func (c *NoOpCommand) Execute() error {
	return nil
}

// Undo implementation for NoOpCommand
func (c *NoOpCommand) Undo() error {
	return nil
}

// String returns a description of the NoOpCommand
func (c *NoOpCommand) String() string {
	return "No operation"
}

// MacroCommand is a command that executes multiple commands in sequence
type MacroCommand struct {
	commands []Command
	name     string
}

// NewMacroCommand creates a new MacroCommand with the given name and commands
func NewMacroCommand(name string, commands ...Command) *MacroCommand {
	return &MacroCommand{
		commands: commands,
		name:     name,
	}
}

// Execute runs all commands in the macro in sequence
func (m *MacroCommand) Execute() error {
	for _, cmd := range m.commands {
		if err := cmd.Execute(); err != nil {
			return fmt.Errorf("macro command '%s' failed: %w", m.name, err)
		}
	}
	return nil
}

// Undo reverses all commands in the macro in reverse order
func (m *MacroCommand) Undo() error {
	// Undo commands in reverse order
	for i := len(m.commands) - 1; i >= 0; i-- {
		if err := m.commands[i].Undo(); err != nil {
			return fmt.Errorf("macro command '%s' undo failed: %w", m.name, err)
		}
	}
	return nil
}

// String returns a description of the MacroCommand
func (m *MacroCommand) String() string {
	return fmt.Sprintf("Macro: %s (%d commands)", m.name, len(m.commands))
}

// AddCommand adds a command to the macro
func (m *MacroCommand) AddCommand(cmd Command) {
	m.commands = append(m.commands, cmd)
}

// RemoteControl is the invoker that executes commands
type RemoteControl struct {
	onCommands  []Command
	offCommands []Command
	history     []Command
	maxHistory  int
}

// NewRemoteControl creates a new RemoteControl with the specified number of slots
func NewRemoteControl(slots int) *RemoteControl {
	onCommands := make([]Command, slots)
	offCommands := make([]Command, slots)
	
	// Initialize with NoOpCommand
	noOp := &NoOpCommand{}
	for i := 0; i < slots; i++ {
		onCommands[i] = noOp
		offCommands[i] = noOp
	}
	
	return &RemoteControl{
		onCommands:  onCommands,
		offCommands: offCommands,
		history:     make([]Command, 0),
		maxHistory:  20, // Default history size
	}
}

// SetCommand assigns commands to a specific slot
func (r *RemoteControl) SetCommand(slot int, onCommand, offCommand Command) error {
	if slot < 0 || slot >= len(r.onCommands) {
		return fmt.Errorf("invalid slot: %d", slot)
	}
	
	r.onCommands[slot] = onCommand
	r.offCommands[slot] = offCommand
	return nil
}

// PressOn executes the on command for the specified slot
func (r *RemoteControl) PressOn(slot int) error {
	if slot < 0 || slot >= len(r.onCommands) {
		return fmt.Errorf("invalid slot: %d", slot)
	}
	
	cmd := r.onCommands[slot]
	err := cmd.Execute()
	if err == nil {
		r.addToHistory(cmd)
	}
	return err
}

// PressOff executes the off command for the specified slot
func (r *RemoteControl) PressOff(slot int) error {
	if slot < 0 || slot >= len(r.offCommands) {
		return fmt.Errorf("invalid slot: %d", slot)
	}
	
	cmd := r.offCommands[slot]
	err := cmd.Execute()
	if err == nil {
		r.addToHistory(cmd)
	}
	return err
}

// Undo reverts the last command executed
func (r *RemoteControl) Undo() error {
	if len(r.history) == 0 {
		return fmt.Errorf("no commands to undo")
	}
	
	lastIndex := len(r.history) - 1
	lastCommand := r.history[lastIndex]
	r.history = r.history[:lastIndex]
	
	return lastCommand.Undo()
}

// GetHistory returns the command history
func (r *RemoteControl) GetHistory() []Command {
	return r.history
}

// ClearHistory clears the command history
func (r *RemoteControl) ClearHistory() {
	r.history = make([]Command, 0)
}

// SetMaxHistory sets the maximum number of commands to keep in history
func (r *RemoteControl) SetMaxHistory(max int) {
	r.maxHistory = max
	// Trim history if needed
	if len(r.history) > r.maxHistory {
		r.history = r.history[len(r.history)-r.maxHistory:]
	}
}

// addToHistory adds a command to the history
func (r *RemoteControl) addToHistory(cmd Command) {
	r.history = append(r.history, cmd)
	// Trim history if needed
	if len(r.history) > r.maxHistory {
		r.history = r.history[1:]
	}
}

// CommandQueue represents a queue of commands to be executed
type CommandQueue struct {
	queue []QueuedCommand
}

// QueuedCommand is a command with execution time
type QueuedCommand struct {
	Command       Command
	ExecutionTime time.Time
}

// NewCommandQueue creates a new CommandQueue
func NewCommandQueue() *CommandQueue {
	return &CommandQueue{
		queue: make([]QueuedCommand, 0),
	}
}

// AddCommand adds a command to the queue to be executed immediately
func (q *CommandQueue) AddCommand(cmd Command) {
	q.AddScheduledCommand(cmd, time.Now())
}

// AddScheduledCommand adds a command to the queue to be executed at a specific time
func (q *CommandQueue) AddScheduledCommand(cmd Command, executionTime time.Time) {
	qc := QueuedCommand{
		Command:       cmd,
		ExecutionTime: executionTime,
	}
	
	// Find the right position to insert the command based on execution time
	index := 0
	for index < len(q.queue) && q.queue[index].ExecutionTime.Before(executionTime) {
		index++
	}
	
	// Insert the command at the correct position
	if index == len(q.queue) {
		// Append to the end if it's the latest command
		q.queue = append(q.queue, qc)
	} else {
		// Insert in the middle by creating a new slice with room for the new command
		newQueue := make([]QueuedCommand, 0, len(q.queue)+1)
		newQueue = append(newQueue, q.queue[:index]...)
		newQueue = append(newQueue, qc)
		newQueue = append(newQueue, q.queue[index:]...)
		q.queue = newQueue
	}
}

// ExecuteDue executes all commands that are due
func (q *CommandQueue) ExecuteDue() (int, error) {
	now := time.Now()
	executedCount := 0
	
	// Find all commands that are due (execution time <= now)
	dueIndex := 0
	for dueIndex < len(q.queue) && !q.queue[dueIndex].ExecutionTime.After(now) {
		dueIndex++
	}
	
	// Execute due commands
	for i := 0; i < dueIndex; i++ {
		if err := q.queue[i].Command.Execute(); err != nil {
			// If a command fails, return the error but keep track of how many were executed
			return executedCount, fmt.Errorf("command execution failed: %w", err)
		}
		executedCount++
	}
	
	// Remove executed commands from the queue
	if dueIndex > 0 {
		q.queue = q.queue[dueIndex:]
	}
	
	return executedCount, nil
}

// Clear removes all commands from the queue
func (q *CommandQueue) Clear() {
	q.queue = make([]QueuedCommand, 0)
}

// Size returns the number of commands in the queue
func (q *CommandQueue) Size() int {
	return len(q.queue)
}

// Peek returns the next command to be executed without removing it
func (q *CommandQueue) Peek() (Command, time.Time, error) {
	if len(q.queue) == 0 {
		return nil, time.Time{}, fmt.Errorf("queue is empty")
	}
	return q.queue[0].Command, q.queue[0].ExecutionTime, nil
}
