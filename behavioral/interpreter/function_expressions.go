package interpreter

import (
	"fmt"
	"math"
)

// FunctionExpression represents a mathematical function
type FunctionExpression struct {
	name      string
	argument  Expression
	operation func(float64) (float64, error)
}

// NewFunctionExpression creates a new function expression
func NewFunctionExpression(name string, argument Expression) (*FunctionExpression, error) {
	var operation func(float64) (float64, error)

	switch name {
	case "sin":
		operation = func(x float64) (float64, error) {
			return math.Sin(x), nil
		}
	case "cos":
		operation = func(x float64) (float64, error) {
			return math.Cos(x), nil
		}
	case "tan":
		operation = func(x float64) (float64, error) {
			return math.Tan(x), nil
		}
	case "sqrt":
		operation = func(x float64) (float64, error) {
			if x < 0 {
				return 0, fmt.Errorf("cannot calculate square root of negative number: %f", x)
			}
			return math.Sqrt(x), nil
		}
	case "log":
		operation = func(x float64) (float64, error) {
			if x <= 0 {
				return 0, fmt.Errorf("cannot calculate logarithm of non-positive number: %f", x)
			}
			return math.Log(x), nil
		}
	case "abs":
		operation = func(x float64) (float64, error) {
			return math.Abs(x), nil
		}
	default:
		return nil, fmt.Errorf("unknown function: %s", name)
	}

	return &FunctionExpression{
		name:      name,
		argument:  argument,
		operation: operation,
	}, nil
}

// Interpret evaluates the function with the argument
func (f *FunctionExpression) Interpret(ctx Context) (float64, error) {
	// First, interpret the argument
	argValue, err := f.argument.Interpret(ctx)
	if err != nil {
		return 0, err
	}

	// Then apply the function
	return f.operation(argValue)
}

// String returns a string representation of the function expression
func (f *FunctionExpression) String() string {
	return fmt.Sprintf("%s(%s)", f.name, f.argument.String())
}
