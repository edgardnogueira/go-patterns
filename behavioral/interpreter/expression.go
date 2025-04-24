package interpreter

// Expression defines the interface for all expression types
// in the interpreter pattern. It declares the Interpret method
// that evaluates the expression in a given context.
type Expression interface {
	// Interpret evaluates the expression with the given context
	// and returns the result as a float64
	Interpret(ctx Context) (float64, error)

	// String returns a string representation of the expression
	// useful for debugging and visualization
	String() string
}

// Context holds the state information during interpretation
// such as variable values, function definitions, etc.
type Context struct {
	// Variables stores the mapping of variable names to their values
	Variables map[string]float64

	// Parent allows for hierarchical contexts
	Parent *Context
}

// NewContext creates a new context with initialized maps
func NewContext() Context {
	return Context{
		Variables: make(map[string]float64),
	}
}

// SetVariable assigns a value to a variable in the context
func (c *Context) SetVariable(name string, value float64) {
	c.Variables[name] = value
}

// GetVariable retrieves a variable value from the context
// If the variable is not found and a parent context exists,
// it will try to retrieve it from the parent
func (c *Context) GetVariable(name string) (float64, bool) {
	if value, exists := c.Variables[name]; exists {
		return value, true
	}

	// Check parent context if available
	if c.Parent != nil {
		return c.Parent.GetVariable(name)
	}

	return 0, false
}
