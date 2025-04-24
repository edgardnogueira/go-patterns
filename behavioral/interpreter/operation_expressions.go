package interpreter

import (
	"fmt"
)

// BinaryOperation is the base type for all binary operations
type BinaryOperation struct {
	left  Expression
	right Expression
}

// AddExpression represents the addition operation
type AddExpression struct {
	BinaryOperation
}

// NewAddExpression creates a new addition expression
func NewAddExpression(left, right Expression) *AddExpression {
	return &AddExpression{
		BinaryOperation: BinaryOperation{
			left:  left,
			right: right,
		},
	}
}

// Interpret implements the addition operation
func (a *AddExpression) Interpret(ctx Context) (float64, error) {
	leftVal, err := a.left.Interpret(ctx)
	if err != nil {
		return 0, err
	}

	rightVal, err := a.right.Interpret(ctx)
	if err != nil {
		return 0, err
	}

	return leftVal + rightVal, nil
}

// String returns a string representation of the addition expression
func (a *AddExpression) String() string {
	return fmt.Sprintf("(%s + %s)", a.left.String(), a.right.String())
}

// SubtractExpression represents the subtraction operation
type SubtractExpression struct {
	BinaryOperation
}

// NewSubtractExpression creates a new subtraction expression
func NewSubtractExpression(left, right Expression) *SubtractExpression {
	return &SubtractExpression{
		BinaryOperation: BinaryOperation{
			left:  left,
			right: right,
		},
	}
}

// Interpret implements the subtraction operation
func (s *SubtractExpression) Interpret(ctx Context) (float64, error) {
	leftVal, err := s.left.Interpret(ctx)
	if err != nil {
		return 0, err
	}

	rightVal, err := s.right.Interpret(ctx)
	if err != nil {
		return 0, err
	}

	return leftVal - rightVal, nil
}

// String returns a string representation of the subtraction expression
func (s *SubtractExpression) String() string {
	return fmt.Sprintf("(%s - %s)", s.left.String(), s.right.String())
}

// MultiplyExpression represents the multiplication operation
type MultiplyExpression struct {
	BinaryOperation
}

// NewMultiplyExpression creates a new multiplication expression
func NewMultiplyExpression(left, right Expression) *MultiplyExpression {
	return &MultiplyExpression{
		BinaryOperation: BinaryOperation{
			left:  left,
			right: right,
		},
	}
}

// Interpret implements the multiplication operation
func (m *MultiplyExpression) Interpret(ctx Context) (float64, error) {
	leftVal, err := m.left.Interpret(ctx)
	if err != nil {
		return 0, err
	}

	rightVal, err := m.right.Interpret(ctx)
	if err != nil {
		return 0, err
	}

	return leftVal * rightVal, nil
}

// String returns a string representation of the multiplication expression
func (m *MultiplyExpression) String() string {
	return fmt.Sprintf("(%s * %s)", m.left.String(), m.right.String())
}

// DivideExpression represents the division operation
type DivideExpression struct {
	BinaryOperation
}

// NewDivideExpression creates a new division expression
func NewDivideExpression(left, right Expression) *DivideExpression {
	return &DivideExpression{
		BinaryOperation: BinaryOperation{
			left:  left,
			right: right,
		},
	}
}

// Interpret implements the division operation
func (d *DivideExpression) Interpret(ctx Context) (float64, error) {
	leftVal, err := d.left.Interpret(ctx)
	if err != nil {
		return 0, err
	}

	rightVal, err := d.right.Interpret(ctx)
	if err != nil {
		return 0, err
	}

	if rightVal == 0 {
		return 0, fmt.Errorf("division by zero")
	}

	return leftVal / rightVal, nil
}

// String returns a string representation of the division expression
func (d *DivideExpression) String() string {
	return fmt.Sprintf("(%s / %s)", d.left.String(), d.right.String())
}
