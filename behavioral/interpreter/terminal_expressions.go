package interpreter

import (
	"fmt"
	"strconv"
)

// NumberExpression represents a terminal expression for numeric literals
type NumberExpression struct {
	value float64
}

// NewNumberExpression creates a new number expression with the given value
func NewNumberExpression(value float64) *NumberExpression {
	return &NumberExpression{value: value}
}

// Interpret returns the value of the number expression
func (n *NumberExpression) Interpret(_ Context) (float64, error) {
	return n.value, nil
}

// String returns the string representation of the number
func (n *NumberExpression) String() string {
	return strconv.FormatFloat(n.value, 'f', -1, 64)
}

// VariableExpression represents a terminal expression for variables
type VariableExpression struct {
	name string
}

// NewVariableExpression creates a new variable expression with the given name
func NewVariableExpression(name string) *VariableExpression {
	return &VariableExpression{name: name}
}

// Interpret looks up the variable in the context and returns its value
func (v *VariableExpression) Interpret(ctx Context) (float64, error) {
	value, exists := ctx.GetVariable(v.name)
	if !exists {
		return 0, fmt.Errorf("variable '%s' not defined", v.name)
	}
	return value, nil
}

// String returns the name of the variable
func (v *VariableExpression) String() string {
	return v.name
}
