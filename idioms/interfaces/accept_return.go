// Package interfaces demonstrates idiomatic Go interface implementation patterns
package interfaces

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Accept interfaces, return structs pattern
// ----------------------------------------
// This pattern makes your functions more flexible by accepting interfaces
// (which allows for many implementations) while returning concrete types
// that give callers the full range of functionality they need.

// DataProcessor represents something that can process data
type DataProcessor interface {
	Process(data string) string
}

// SimpleProcessor is a basic implementation of DataProcessor
type SimpleProcessor struct {
	prefix string
}

// Process implements the DataProcessor interface
func (p SimpleProcessor) Process(data string) string {
	return p.prefix + ": " + data
}

// Result is a concrete type returned by our functions
type Result struct {
	Data       string
	LineCount  int
	WordCount  int
	ByteCount  int
	IsModified bool
}

// ProcessData accepts an interface and returns a struct
// This demonstrates the "accept interfaces, return structs" pattern
func ProcessData(processor DataProcessor, data string) *Result {
	// Use the interface to process the data
	processed := processor.Process(data)
	
	// Return a concrete type with rich functionality
	return &Result{
		Data:       processed,
		LineCount:  len(strings.Split(processed, "\n")),
		WordCount:  len(strings.Fields(processed)),
		ByteCount:  len(processed),
		IsModified: processed != data,
	}
}

// InterfacePromise represents a promise to provide data
type InterfacePromise interface {
	Fulfill() (string, error)
}

// StringSource is a simple implementation of InterfacePromise
type StringSource struct {
	data string
	err  error
}

// Fulfill implements the InterfacePromise interface
func (s StringSource) Fulfill() (string, error) {
	return s.data, s.err
}

// DataLoader loads data from a promise and transforms it
func DataLoader(promise InterfacePromise) (*Result, error) {
	// Accept an interface (flexible)
	data, err := promise.Fulfill()
	if err != nil {
		return nil, fmt.Errorf("data loading failed: %w", err)
	}
	
	// Return a struct (concrete)
	return &Result{
		Data:       data,
		LineCount:  len(strings.Split(data, "\n")),
		WordCount:  len(strings.Fields(data)),
		ByteCount:  len(data),
		IsModified: false,
	}, nil
}

// AcceptReturnDemo demonstrates the "accept interfaces, return structs" pattern
func AcceptReturnDemo() {
	fmt.Println("=== Accept Interfaces, Return Structs Demo ===")
	
	// Create a processor
	processor := SimpleProcessor{prefix: "PROCESSED"}
	
	// Use our function that accepts an interface and returns a struct
	result := ProcessData(processor, "Hello, world!\nThis is a test.")
	
	// Now we have a rich struct to work with
	fmt.Printf("Processed data: %s\n", result.Data)
	fmt.Printf("Stats: %d lines, %d words, %d bytes\n", 
		result.LineCount, result.WordCount, result.ByteCount)
	
	// Example with data loading
	fmt.Println("\nLoading data from a promise:")
	
	source := StringSource{data: "Promised data\nwith multiple lines"}
	loadResult, _ := DataLoader(source)
	
	fmt.Printf("Loaded data: %s\n", loadResult.Data)
	fmt.Printf("Stats: %d lines, %d words, %d bytes\n",
		loadResult.LineCount, loadResult.WordCount, loadResult.ByteCount)
}

// Additional examples of the pattern
// --------------------------------

// Reader and Writer are interfaces from the io package
// io.Copy is a perfect example of this pattern - it accepts interfaces (Reader, Writer)
// and returns a concrete int64 and error.

// StringSorter creates a function that sorts lines in a string
func StringSorter(r io.Reader) (string, error) {
	// Accept an interface (io.Reader)
	
	// Read all data from the reader
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	
	// Split into lines and sort
	lines := strings.Split(string(data), "\n")
	sort.Strings(lines)
	
	// Return a concrete type (string)
	return strings.Join(lines, "\n"), nil
}

// Example of making a utility function flexible with interfaces
// ----------------------------------------------------------

// Calculator defines calculation behavior
type Calculator interface {
	Calculate(a, b int) int
}

// Adder implements Calculator by adding numbers
type Adder struct{}

// Calculate implements the Calculator interface
func (a Adder) Calculate(x, y int) int {
	return x + y
}

// Multiplier implements Calculator by multiplying numbers
type Multiplier struct{}

// Calculate implements the Calculator interface
func (m Multiplier) Calculate(x, y int) int {
	return x * y
}

// CalculationResult is a concrete return type
type CalculationResult struct {
	Input1     int
	Input2     int
	Result     int
	Operation  string
	IsPositive bool
}

// PerformCalculation demonstrates the pattern with calculations
func PerformCalculation(calc Calculator, a, b int) *CalculationResult {
	// Accept an interface (Calculator)
	result := calc.Calculate(a, b)
	
	// Figure out what operation was performed
	var operation string
	switch calc.(type) {
	case Adder:
		operation = "addition"
	case Multiplier:
		operation = "multiplication"
	default:
		operation = "unknown"
	}
	
	// Return a concrete struct with rich information
	return &CalculationResult{
		Input1:     a,
		Input2:     b,
		Result:     result,
		Operation:  operation,
		IsPositive: result >= 0,
	}
}

// CalculatorDemo demonstrates the accept interfaces pattern with calculations
func CalculatorDemo() {
	fmt.Println("\n=== Calculator Demo ===")
	
	add := Adder{}
	mult := Multiplier{}
	
	// Using the same function with different implementations
	addResult := PerformCalculation(add, 5, 3)
	multResult := PerformCalculation(mult, 5, 3)
	
	fmt.Printf("Addition: %d %s %d = %d\n", 
		addResult.Input1, addResult.Operation, addResult.Input2, addResult.Result)
	fmt.Printf("Multiplication: %d %s %d = %d\n", 
		multResult.Input1, multResult.Operation, multResult.Input2, multResult.Result)
}

// Real-world example: standard library
// ----------------------------------

// SortData demonstrates another example of accepting interfaces
func SortData() {
	fmt.Println("\n=== Standard Library Pattern Example ===")
	
	// strings.NewReader accepts a concrete type (string)
	// and returns something that satisfies many interfaces (io.Reader, etc.)
	r := strings.NewReader("banana\napple\ncherry\ndate")
	
	// Our function accepts an interface (io.Reader)
	sorted, _ := StringSorter(r)
	fmt.Printf("Sorted data:\n%s\n", sorted)
	
	// io.Copy is another great example from the standard library:
	// func Copy(dst Writer, src Reader) (written int64, err error)
	// It accepts interfaces but returns concrete types (int64, error)
}

// Functions Returning Interfaces
// ----------------------------

// NewProcessor creates a new processor with the given prefix
// This function returns an interface, which is acceptable when:
// 1. You want to hide implementation details
// 2. You need polymorphic behavior determined at runtime
func NewProcessor(prefix string) DataProcessor {
	return SimpleProcessor{prefix: prefix}
}

// ReturningInterfacesDemo shows when it makes sense to return an interface
func ReturningInterfacesDemo() {
	fmt.Println("\n=== Returning Interfaces Demo ===")
	
	// Factory function returns an interface
	processor := NewProcessor("PREFIX")
	
	// We can only use methods defined by the interface
	result := processor.Process("Some data")
	fmt.Printf("Result: %s\n", result)
	
	// We don't know the concrete type (and that's the point)
	fmt.Printf("Type: %T\n", processor)
}

// Demonstration of the full pattern
// -------------------------------

// RunAcceptReturnDemo demonstrates all aspects of the "accept interfaces, return structs" pattern
func RunAcceptReturnDemo() {
	fmt.Println("============================================")
	fmt.Println("Accept Interfaces, Return Structs Pattern")
	fmt.Println("============================================")
	
	AcceptReturnDemo()
	CalculatorDemo()
	SortData()
	ReturningInterfacesDemo()
	
	fmt.Println("\nKey benefits of this pattern:")
	fmt.Println("1. Flexibility - functions can accept any type implementing the interface")
	fmt.Println("2. Concrete returns - callers get all functionality they need")
	fmt.Println("3. Testability - easy to provide mock implementations for testing")
	fmt.Println("4. Decoupling - minimal dependencies between components")
	fmt.Println("============================================")
}
