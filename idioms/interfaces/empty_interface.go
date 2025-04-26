// Package interfaces demonstrates idiomatic Go interface implementation patterns
package interfaces

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// Empty Interface and Type Assertions
// ----------------------------------
// The empty interface (interface{}) or "any" in Go 1.18+ can hold values of any type.
// However, using it requires type assertions or type switches to recover the underlying type.

// GenericContainer demonstrates using the empty interface to store any type
type GenericContainer struct {
	Value any // empty interface (interface{} in earlier Go versions)
}

// Print prints the container value, handling different types
func (c GenericContainer) Print() {
	// Type switch to handle different types
	switch v := c.Value.(type) {
	case int:
		fmt.Printf("Integer: %d\n", v)
	case string:
		fmt.Printf("String: %s\n", v)
	case bool:
		fmt.Printf("Boolean: %t\n", v)
	case []int:
		fmt.Printf("Integer slice: %v\n", v)
	case map[string]int:
		fmt.Printf("Map[string]int: %v\n", v)
	case fmt.Stringer:
		// Interface satisfaction check - this will match any type implementing Stringer
		fmt.Printf("Stringer: %s\n", v.String())
	default:
		// Use reflection for unknown types
		fmt.Printf("Unknown type (%T): %v\n", v, v)
	}
}

// StringerValue is an example type that implements fmt.Stringer
type StringerValue struct {
	Data string
}

// String implements fmt.Stringer
func (s StringerValue) String() string {
	return fmt.Sprintf("StringerValue{%s}", s.Data)
}

// Type assertion examples
// ---------------------

// TypeAssertionBasics demonstrates basic type assertions
func TypeAssertionBasics() {
	// Start with an empty interface value
	var i any

	// Assign different values to it
	i = 42
	
	// Type assertion with direct assignment
	// This will panic if the type doesn't match
	intVal := i.(int)
	fmt.Printf("Direct type assertion: %d\n", intVal)
	
	// Type assertion with check - safer approach
	strVal, ok := i.(string)
	if ok {
		fmt.Printf("String value: %s\n", strVal)
	} else {
		fmt.Printf("Value is not a string, it's a %T\n", i)
	}
	
	// Set value to a string and try again
	i = "hello"
	strVal, ok = i.(string)
	if ok {
		fmt.Printf("String value: %s\n", strVal)
	} else {
		fmt.Printf("Value is not a string\n")
	}
}

// TypeSwitchExample demonstrates using type switches
func TypeSwitchExample(value any) {
	// Type switch is cleaner than multiple type assertions
	switch v := value.(type) {
	case nil:
		fmt.Println("Value is nil")
	case int:
		fmt.Printf("Integer: %d\n", v)
	case float64:
		fmt.Printf("Float64: %f\n", v)
	case string:
		fmt.Printf("String: %s\n", v)
	case bool:
		fmt.Printf("Boolean: %t\n", v)
	case []byte:
		// Handle byte slices specially
		if len(v) > 10 {
			fmt.Printf("Byte slice (first 10): %v...\n", v[:10])
		} else {
			fmt.Printf("Byte slice: %v\n", v)
		}
	case fmt.Stringer:
		// Interface check - matches anything with String() method
		fmt.Printf("Stringer: %s\n", v.String())
	default:
		// Default case handles any other type
		fmt.Printf("Unknown type: %T\n", v)
	}
}

// Common use cases for empty interfaces
// -----------------------------------

// ParseValue converts a string to various types using the empty interface
func ParseValue(s string) (any, error) {
	// Try to parse as int
	if i, err := strconv.Atoi(s); err == nil {
		return i, nil
	}
	
	// Try to parse as float
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f, nil
	}
	
	// Try to parse as bool
	if b, err := strconv.ParseBool(s); err == nil {
		return b, nil
	}
	
	// Try to parse as JSON
	var jsonVal any
	if err := json.Unmarshal([]byte(s), &jsonVal); err == nil {
		return jsonVal, nil
	}
	
	// Default to string
	return s, nil
}

// Collection demonstrates using empty interface for heterogeneous collections
func Collection() {
	// Slice of any type
	collection := []any{
		42,
		"hello",
		true,
		[]int{1, 2, 3},
		map[string]int{"a": 1, "b": 2},
		StringerValue{"custom type"},
	}
	
	// Process mixed collection
	fmt.Println("Processing collection with different types:")
	for i, item := range collection {
		fmt.Printf("Item %d: ", i)
		TypeSwitchExample(item)
	}
}

// Config demonstrates a simple configuration using map[string]any
type Config map[string]any

// GetInt safely gets an int from the config
func (c Config) GetInt(key string, defaultVal int) int {
	if val, ok := c[key]; ok {
		if intVal, ok := val.(int); ok {
			return intVal
		}
		
		// Try to convert from other numeric types
		switch v := val.(type) {
		case float64:
			return int(v)
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				return i
			}
		}
	}
	return defaultVal
}

// GetString safely gets a string from the config
func (c Config) GetString(key string, defaultVal string) string {
	if val, ok := c[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
		
		// Convert other types to string
		switch v := val.(type) {
		case int:
			return strconv.Itoa(v)
		case float64:
			return strconv.FormatFloat(v, 'f', -1, 64)
		case bool:
			return strconv.FormatBool(v)
		case fmt.Stringer:
			return v.String()
		}
	}
	return defaultVal
}

// GetStringList safely gets a string slice from the config
func (c Config) GetStringList(key string, defaultVal []string) []string {
	if val, ok := c[key]; ok {
		// Check if it's already a string slice
		if strList, ok := val.([]string); ok {
			return strList
		}
		
		// Check if it's a slice of any
		if anyList, ok := val.([]any); ok {
			strList := make([]string, len(anyList))
			for i, v := range anyList {
				switch item := v.(type) {
				case string:
					strList[i] = item
				case fmt.Stringer:
					strList[i] = item.String()
				default:
					// Convert other types using fmt
					strList[i] = fmt.Sprintf("%v", item)
				}
			}
			return strList
		}
	}
	return defaultVal
}

// ConfigExample demonstrates the config use case
func ConfigExample() {
	// Create a config with various types
	cfg := Config{
		"port":      8080,
		"host":      "localhost",
		"debug":     true,
		"timeout":   30.5,
		"tags":      []any{"web", "api", 123},
		"stringer":  StringerValue{"config value"},
		"rawValues": []int{1, 2, 3},
	}
	
	// Access with type-safe getters
	port := cfg.GetInt("port", 80)
	host := cfg.GetString("host", "127.0.0.1")
	timeout := cfg.GetInt("timeout", 60)
	tags := cfg.GetStringList("tags", []string{})
	
	fmt.Printf("Config values:\n")
	fmt.Printf("  port: %d\n", port)
	fmt.Printf("  host: %s\n", host)
	fmt.Printf("  timeout: %d\n", timeout)
	fmt.Printf("  tags: %v\n", tags)
	
	// Show missing with default
	missing := cfg.GetString("missing", "default-value")
	fmt.Printf("  missing: %s\n", missing)
}

// When to avoid empty interface
// ---------------------------

// BadFunction uses empty interface unnecessarily - generally an anti-pattern
func BadFunction(data any) any {
	// This function doesn't provide type safety or clear expectations
	return data
}

// BetterFunction uses concrete types - preferred when possible
func BetterFunction(data string) int {
	return len(data)
}

// Using reflection with any
// -----------------------

// PrintReflection uses reflection to inspect a value of any type
func PrintReflection(v any) {
	// Get the reflect.Value
	val := reflect.ValueOf(v)
	
	// Print basic type information
	fmt.Printf("Type: %s, Kind: %s\n", val.Type(), val.Kind())
	
	// Different handling based on kind
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Printf("Integer value: %d\n", val.Int())
	
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fmt.Printf("Unsigned integer value: %d\n", val.Uint())
	
	case reflect.Float32, reflect.Float64:
		fmt.Printf("Float value: %f\n", val.Float())
	
	case reflect.Bool:
		fmt.Printf("Boolean value: %t\n", val.Bool())
	
	case reflect.String:
		fmt.Printf("String value: %s\n", val.String())
	
	case reflect.Struct:
		fmt.Printf("Struct with %d fields:\n", val.NumField())
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			fmt.Printf("  %s: %v\n", field.Name, val.Field(i).Interface())
		}
	
	case reflect.Map:
		fmt.Printf("Map with %d entries:\n", val.Len())
		for _, key := range val.MapKeys() {
			fmt.Printf("  %v: %v\n", key.Interface(), val.MapIndex(key).Interface())
		}
	
	case reflect.Slice, reflect.Array:
		fmt.Printf("Slice/Array with %d elements:\n", val.Len())
		for i := 0; i < val.Len(); i++ {
			fmt.Printf("  [%d]: %v\n", i, val.Index(i).Interface())
		}
	
	case reflect.Ptr:
		if val.IsNil() {
			fmt.Println("Nil pointer")
		} else {
			fmt.Printf("Pointer to %s: %v\n", val.Elem().Type(), val.Elem().Interface())
		}
	
	case reflect.Interface:
		if val.IsNil() {
			fmt.Println("Nil interface")
		} else {
			fmt.Printf("Interface containing %s: %v\n", val.Elem().Type(), val.Elem().Interface())
		}
	
	default:
		fmt.Printf("Unhandled kind: %v\n", val.Kind())
	}
}

// EmptyInterfaceDemo demonstrates the empty interface and type assertions
func EmptyInterfaceDemo() {
	fmt.Println("============================================")
	fmt.Println("Empty Interface and Type Assertions Demo")
	fmt.Println("============================================")
	
	// Basic examples of containers using empty interface
	fmt.Println("Generic container examples:")
	containers := []GenericContainer{
		{Value: 42},
		{Value: "Hello, World!"},
		{Value: true},
		{Value: []int{1, 2, 3}},
		{Value: map[string]int{"a": 1, "b": 2}},
		{Value: StringerValue{"custom data"}},
	}
	
	for _, c := range containers {
		c.Print()
	}
	
	// Type assertion examples
	fmt.Println("\nType assertion examples:")
	TypeAssertionBasics()
	
	// Parse value examples
	fmt.Println("\nParsing strings to different types:")
	examples := []string{
		"42",
		"3.14",
		"true",
		`{"name":"John","age":30}`,
		"hello",
	}
	
	for _, ex := range examples {
		val, _ := ParseValue(ex)
		fmt.Printf("Parsed '%s' to: %T (%v)\n", ex, val, val)
	}
	
	// Collection example
	fmt.Println("\nCollection example:")
	Collection()
	
	// Config example
	fmt.Println("\nConfig example:")
	ConfigExample()
	
	// Reflection example
	fmt.Println("\nReflection examples:")
	
	type Person struct {
		Name string
		Age  int
	}
	
	PrintReflection(42)
	PrintReflection("hello")
	PrintReflection(Person{Name: "Alice", Age: 30})
	PrintReflection(map[string]int{"x": 1, "y": 2})
	PrintReflection([]string{"a", "b", "c"})
	
	fmt.Println("\nBest Practices for Empty Interface:")
	fmt.Println("1. Use sparingly - prefer concrete types or specific interfaces")
	fmt.Println("2. Always use type assertions with the ok check to avoid panics")
	fmt.Println("3. Consider type switches for handling multiple types")
	fmt.Println("4. Use for config systems, plugins, or other cases needing dynamic typing")
	fmt.Println("5. Type assertions and switches are more efficient than reflection")
	fmt.Println("============================================")
}
