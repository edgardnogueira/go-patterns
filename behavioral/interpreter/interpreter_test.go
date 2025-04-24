package interpreter

import (
	"testing"
)

func TestNumberExpression(t *testing.T) {
	ctx := NewContext()
	expr := NewNumberExpression(42.5)

	result, err := expr.Interpret(ctx)
	if err != nil {
		t.Errorf("Failed to interpret number expression: %v", err)
	}

	if result != 42.5 {
		t.Errorf("Expected 42.5, got %f", result)
	}
}

func TestVariableExpression(t *testing.T) {
	ctx := NewContext()
	ctx.SetVariable("x", 10)
	expr := NewVariableExpression("x")

	result, err := expr.Interpret(ctx)
	if err != nil {
		t.Errorf("Failed to interpret variable expression: %v", err)
	}

	if result != 10 {
		t.Errorf("Expected 10, got %f", result)
	}

	// Test undefined variable
	expr = NewVariableExpression("y")
	_, err = expr.Interpret(ctx)
	if err == nil {
		t.Error("Expected error for undefined variable, got nil")
	}
}

func TestAddExpression(t *testing.T) {
	ctx := NewContext()
	left := NewNumberExpression(5)
	right := NewNumberExpression(3)
	expr := NewAddExpression(left, right)

	result, err := expr.Interpret(ctx)
	if err != nil {
		t.Errorf("Failed to interpret add expression: %v", err)
	}

	if result != 8 {
		t.Errorf("Expected 8, got %f", result)
	}
}

func TestSubtractExpression(t *testing.T) {
	ctx := NewContext()
	left := NewNumberExpression(5)
	right := NewNumberExpression(3)
	expr := NewSubtractExpression(left, right)

	result, err := expr.Interpret(ctx)
	if err != nil {
		t.Errorf("Failed to interpret subtract expression: %v", err)
	}

	if result != 2 {
		t.Errorf("Expected 2, got %f", result)
	}
}

func TestMultiplyExpression(t *testing.T) {
	ctx := NewContext()
	left := NewNumberExpression(5)
	right := NewNumberExpression(3)
	expr := NewMultiplyExpression(left, right)

	result, err := expr.Interpret(ctx)
	if err != nil {
		t.Errorf("Failed to interpret multiply expression: %v", err)
	}

	if result != 15 {
		t.Errorf("Expected 15, got %f", result)
	}
}

func TestDivideExpression(t *testing.T) {
	ctx := NewContext()
	left := NewNumberExpression(6)
	right := NewNumberExpression(3)
	expr := NewDivideExpression(left, right)

	result, err := expr.Interpret(ctx)
	if err != nil {
		t.Errorf("Failed to interpret divide expression: %v", err)
	}

	if result != 2 {
		t.Errorf("Expected 2, got %f", result)
	}

	// Test division by zero
	right = NewNumberExpression(0)
	expr = NewDivideExpression(left, right)
	_, err = expr.Interpret(ctx)
	if err == nil {
		t.Error("Expected error for division by zero, got nil")
	}
}

func TestFunctionExpression(t *testing.T) {
	ctx := NewContext()
	arg := NewNumberExpression(0)
	
	// Test sin function
	expr, err := NewFunctionExpression("sin", arg)
	if err != nil {
		t.Errorf("Failed to create sin function: %v", err)
	}

	result, err := expr.Interpret(ctx)
	if err != nil {
		t.Errorf("Failed to interpret sin function: %v", err)
	}

	if result != 0 {
		t.Errorf("Expected sin(0) = 0, got %f", result)
	}

	// Test sqrt function with negative argument
	arg = NewNumberExpression(-1)
	expr, err = NewFunctionExpression("sqrt", arg)
	if err != nil {
		t.Errorf("Failed to create sqrt function: %v", err)
	}

	_, err = expr.Interpret(ctx)
	if err == nil {
		t.Error("Expected error for sqrt of negative number, got nil")
	}

	// Test unknown function
	_, err = NewFunctionExpression("unknown", arg)
	if err == nil {
		t.Error("Expected error for unknown function, got nil")
	}
}

func TestParser(t *testing.T) {
	tests := []struct {
		input    string
		vars     map[string]float64
		expected float64
	}{
		{"5 + 3", nil, 8},
		{"5 - 3", nil, 2},
		{"5 * 3", nil, 15},
		{"6 / 3", nil, 2},
		{"(5 + 3) * 2", nil, 16},
		{"5 + 3 * 2", nil, 11},
		{"x + y", map[string]float64{"x": 5, "y": 3}, 8},
		{"sin(0)", nil, 0},
		{"cos(0)", nil, 1},
		{"sqrt(9)", nil, 3},
		{"abs(-5)", nil, 5},
	}

	for _, test := range tests {
		parser := NewParser(test.input)
		expr, err := parser.Parse()
		if err != nil {
			t.Errorf("Failed to parse '%s': %v", test.input, err)
			continue
		}

		ctx := NewContext()
		for k, v := range test.vars {
			ctx.SetVariable(k, v)
		}

		result, err := expr.Interpret(ctx)
		if err != nil {
			t.Errorf("Failed to interpret '%s': %v", test.input, err)
			continue
		}

		// Use a small epsilon for floating point comparison
		epsilon := 0.0001
		if result < test.expected-epsilon || result > test.expected+epsilon {
			t.Errorf("For '%s', expected %f, got %f", test.input, test.expected, result)
		}
	}
}

func TestParserErrors(t *testing.T) {
	tests := []struct {
		input string
		error bool
	}{
		{"", true},                // Empty expression
		{"5 +", true},             // Incomplete expression
		{"5 + )", true},           // Unmatched parenthesis
		{"(5 + 3", true},          // Missing closing parenthesis
		{"log(-1)", true},         // Runtime error: log of negative number
		{"sqrt(-4)", true},        // Runtime error: sqrt of negative number
		{"5 / 0", true},           // Runtime error: division by zero
		{"unknown(5)", true},      // Unknown function
	}

	for _, test := range tests {
		parser := NewParser(test.input)
		expr, err := parser.Parse()
		
		if err != nil {
			// If parsing failed, the test passes if we expected an error
			if !test.error {
				t.Errorf("Unexpected parse error for '%s': %v", test.input, err)
			}
			continue
		}

		ctx := NewContext()
		_, err = expr.Interpret(ctx)
		
		if (err != nil) != test.error {
			if test.error {
				t.Errorf("Expected error for '%s', got none", test.input)
			} else {
				t.Errorf("Unexpected error for '%s': %v", test.input, err)
			}
		}
	}
}
