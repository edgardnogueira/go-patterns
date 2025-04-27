package cqrs_test

import (
	"errors"
	"testing"

	"github.com/edgardnogueira/go-patterns/project_structures/event_driven/pkg/cqrs"
)

// TestCommandBusRegisterAndDispatch tests the command bus registration and dispatching
func TestCommandBusRegisterAndDispatch(t *testing.T) {
	// Create a new command bus
	bus := cqrs.NewCommandBus()

	// Define test command and result
	type TestCommand struct {
		Message string
	}

	expectedResult := "Processed: Hello, World!"

	// Register a handler
	bus.Register("TestCommand", func(cmd interface{}) (interface{}, error) {
		testCmd, ok := cmd.(TestCommand)
		if !ok {
			return nil, errors.New("invalid command type")
		}
		return "Processed: " + testCmd.Message, nil
	})

	// Dispatch a command
	result, err := bus.Dispatch("TestCommand", TestCommand{Message: "Hello, World!"})

	// Check results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != expectedResult {
		t.Errorf("Expected result %q, got %q", expectedResult, result)
	}
}

// TestCommandBusUnregisteredCommand tests dispatching an unregistered command
func TestCommandBusUnregisteredCommand(t *testing.T) {
	// Create a new command bus
	bus := cqrs.NewCommandBus()

	// Dispatch a command with no registered handler
	_, err := bus.Dispatch("UnregisteredCommand", struct{}{})

	// Check error
	if err == nil {
		t.Errorf("Expected error for unregistered command, got nil")
	}
}

// TestCommandBusInvalidCommand tests handling an invalid command
func TestCommandBusInvalidCommand(t *testing.T) {
	// Create a new command bus
	bus := cqrs.NewCommandBus()

	// Register a handler that expects a specific type
	bus.Register("TypedCommand", func(cmd interface{}) (interface{}, error) {
		_, ok := cmd.(struct{ ID int })
		if !ok {
			return nil, errors.New("invalid command type")
		}
		return "success", nil
	})

	// Dispatch a command with wrong type
	_, err := bus.Dispatch("TypedCommand", "wrong type")

	// Check error
	if err == nil {
		t.Errorf("Expected error for invalid command type, got nil")
	}
}
