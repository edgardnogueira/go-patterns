package cqrs

import (
	"errors"
	"fmt"
	"sync"
)

// CommandHandler is a function that handles a command
type CommandHandler func(interface{}) (interface{}, error)

// CommandBus routes commands to their appropriate handlers
type CommandBus struct {
	handlers map[string]CommandHandler
	mutex    sync.RWMutex
}

// NewCommandBus creates a new CommandBus
func NewCommandBus() *CommandBus {
	return &CommandBus{
		handlers: make(map[string]CommandHandler),
	}
}

// Register registers a command handler for a specific command type
func (b *CommandBus) Register(commandType string, handler CommandHandler) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.handlers[commandType] = handler
}

// Dispatch sends a command to its registered handler
func (b *CommandBus) Dispatch(commandType string, command interface{}) (interface{}, error) {
	b.mutex.RLock()
	handler, exists := b.handlers[commandType]
	b.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no handler registered for command type: %s", commandType)
	}

	return handler(command)
}
