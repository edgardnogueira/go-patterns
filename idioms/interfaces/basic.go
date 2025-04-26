// Package interfaces demonstrates idiomatic Go interface implementation patterns
package interfaces

import (
	"fmt"
	"math"
)

// Basic interface definition and implementation
// ---------------------------------------------

// Shape is a basic interface that defines a method for calculating area
type Shape interface {
	Area() float64
}

// Circle is a concrete type that implicitly implements the Shape interface
type Circle struct {
	Radius float64
}

// Area calculates the area of a circle
// This method means Circle implicitly implements the Shape interface
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Rectangle is another concrete type that implicitly implements the Shape interface
type Rectangle struct {
	Width, Height float64
}

// Area calculates the area of a rectangle
// This method means Rectangle implicitly implements the Shape interface
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// UseShape demonstrates using a function that accepts an interface
func UseShape(s Shape) {
	fmt.Printf("Area of shape: %.2f\n", s.Area())
}

// BasicInterfaceDemo demonstrates basic interface implementation in Go
func BasicInterfaceDemo() {
	// We can use Circle and Rectangle wherever a Shape is expected
	circle := Circle{Radius: 5}
	rectangle := Rectangle{Width: 4, Height: 6}

	// Both types can be assigned to a Shape variable
	var shape Shape
	shape = circle
	fmt.Printf("Circle area: %.2f\n", shape.Area())

	shape = rectangle
	fmt.Printf("Rectangle area: %.2f\n", shape.Area())

	// We can pass both types to functions accepting the Shape interface
	UseShape(circle)
	UseShape(rectangle)

	// We can create slices of the interface type containing different implementations
	shapes := []Shape{
		Circle{Radius: 3},
		Rectangle{Width: 2, Height: 5},
	}

	// We can iterate through the slice and treat all elements uniformly
	fmt.Println("Areas of different shapes:")
	for i, s := range shapes {
		fmt.Printf("  Shape %d: %.2f\n", i, s.Area())
	}
}

// Multiple interfaces implementation
// ---------------------------------

// Describer is an interface for types that can describe themselves
type Describer interface {
	Describe() string
}

// Both Circle and Rectangle can implement multiple interfaces

// Describe returns a description of the circle
func (c Circle) Describe() string {
	return fmt.Sprintf("Circle with radius %.2f", c.Radius)
}

// Describe returns a description of the rectangle
func (r Rectangle) Describe() string {
	return fmt.Sprintf("Rectangle with width %.2f and height %.2f", r.Width, r.Height)
}

// MultipleInterfacesDemo demonstrates implementing multiple interfaces
func MultipleInterfacesDemo() {
	// An object can implement multiple interfaces
	circle := Circle{Radius: 5}

	// We can use the object as either interface type
	var shape Shape = circle
	var describer Describer = circle

	fmt.Printf("As Shape: Area = %.2f\n", shape.Area())
	fmt.Printf("As Describer: %s\n", describer.Describe())

	// We can also use type assertion to get access to methods of another interface
	if desc, ok := shape.(Describer); ok {
		fmt.Printf("Shape is also a Describer: %s\n", desc.Describe())
	}
}

// Interface method sets with pointers vs values
// --------------------------------------------

// Mover is an interface for things that can move
type Mover interface {
	Move(dx, dy float64)
	GetPosition() (float64, float64)
}

// Point is a struct representing a 2D point
type Point struct {
	X, Y float64
}

// Move changes the point's position (requires a pointer receiver to modify the struct)
func (p *Point) Move(dx, dy float64) {
	p.X += dx
	p.Y += dy
}

// GetPosition returns the current position
func (p Point) GetPosition() (float64, float64) {
	return p.X, p.Y
}

// PointerReceiverDemo demonstrates interface implementation with pointer receivers
func PointerReceiverDemo() {
	p := Point{X: 1, Y: 2}

	// Important: Only *Point implements Mover, not Point
	// Uncommenting the next line would cause a compile error:
	// var m1 Mover = p // Error: Point does not implement Mover

	// Using a pointer works
	var m Mover = &p
	fmt.Printf("Position: (%.1f, %.1f)\n", m.GetPosition())
	m.Move(2, 3)
	fmt.Printf("After Move: (%.1f, %.1f)\n", m.GetPosition())

	// However, we can call pointer methods on value variables
	// Go automatically takes the address when needed
	p.Move(1, 1) // This works because Go converts it to (&p).Move(1, 1)
	fmt.Printf("Called directly: (%.1f, %.1f)\n", p.GetPosition())
}

// Interface value internals and nil interface values
// ------------------------------------------------

// NilInterfaceDemo demonstrates how nil values work with interfaces
func NilInterfaceDemo() {
	// Interface value has two components:
	// 1. Dynamic type (concrete type)
	// 2. Dynamic value (value of that type)
	
	// An interface value is nil only when both type and value are nil
	var s Shape // nil interface value (both type and value are nil)
	fmt.Printf("Nil interface: %v, Is nil: %t\n", s, s == nil)

	// A nil pointer value still has a non-nil interface type
	var ptrCircle *Circle // nil pointer
	s = ptrCircle         // non-nil interface (type is *Circle, value is nil)
	fmt.Printf("Interface with nil value: %v, Is nil: %t\n", s, s == nil)

	// This is important to understand for error handling
	// The following would panic if we tried to call s.Area()
	if s != nil && ptrCircle != nil {
		fmt.Printf("Area: %.2f\n", s.Area())
	} else {
		fmt.Println("Cannot calculate area of nil *Circle")
	}

	// Proper nil check for interface values
	if s == nil {
		fmt.Println("Interface is nil")
	} else {
		fmt.Printf("Interface is not nil (type: %T)\n", s)
	}
}
