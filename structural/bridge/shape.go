package bridge

import (
	"fmt"
	"strings"
)

// Shape is the Abstraction in the Bridge pattern.
// It defines the interface for all shapes.
type Shape interface {
	// Draw renders the shape using the associated drawing API.
	Draw() string
	
	// ResizeTo resizes the shape to the specified dimensions.
	ResizeTo(width, height float64) Shape
	
	// MoveTo moves the shape to the specified coordinates.
	MoveTo(x, y float64) Shape
	
	// GetName returns the name of the shape.
	GetName() string
	
	// GetDrawingAPI returns the drawing API used by the shape.
	GetDrawingAPI() DrawingAPI
	
	// GetDescription returns a description of the shape.
	GetDescription() string
}

// BaseShape is a base struct that implements common functionality for shapes.
// It's a refined Abstraction in the Bridge pattern.
type BaseShape struct {
	name       string
	drawingAPI DrawingAPI
	x          float64
	y          float64
	width      float64
	height     float64
}

// NewBaseShape creates a new BaseShape with the given parameters.
func NewBaseShape(name string, drawingAPI DrawingAPI, x, y, width, height float64) BaseShape {
	return BaseShape{
		name:       name,
		drawingAPI: drawingAPI,
		x:          x,
		y:          y,
		width:      width,
		height:     height,
	}
}

// GetName returns the name of the shape.
func (b *BaseShape) GetName() string {
	return b.name
}

// GetDrawingAPI returns the drawing API used by the shape.
func (b *BaseShape) GetDrawingAPI() DrawingAPI {
	return b.drawingAPI
}

// GetDescription returns a description of the shape.
func (b *BaseShape) GetDescription() string {
	return fmt.Sprintf("%s rendered using %s renderer", b.GetName(), b.drawingAPI.GetName())
}

// Circle is a concrete Shape implementation that represents a circle.
type Circle struct {
	BaseShape
	radius float64
}

// NewCircle creates a new Circle with the given drawing API, center coordinates, and radius.
func NewCircle(drawingAPI DrawingAPI, x, y, radius float64) *Circle {
	return &Circle{
		BaseShape: NewBaseShape("Circle", drawingAPI, x, y, 2*radius, 2*radius),
		radius:    radius,
	}
}

// Draw renders the circle using the associated drawing API.
func (c *Circle) Draw() string {
	return c.drawingAPI.DrawCircle(c.x, c.y, c.radius)
}

// ResizeTo resizes the circle to match the specified width and height (uses average).
func (c *Circle) ResizeTo(width, height float64) Shape {
	// For a circle, we take the average of width and height as the diameter
	avgSize := (width + height) / 2
	c.radius = avgSize / 2
	c.width = avgSize
	c.height = avgSize
	return c
}

// MoveTo moves the circle to the specified center coordinates.
func (c *Circle) MoveTo(x, y float64) Shape {
	c.x = x
	c.y = y
	return c
}

// Square is a concrete Shape implementation that represents a square.
type Square struct {
	BaseShape
	sideLength float64
}

// NewSquare creates a new Square with the given drawing API, top-left coordinates, and side length.
func NewSquare(drawingAPI DrawingAPI, x, y, sideLength float64) *Square {
	return &Square{
		BaseShape:  NewBaseShape("Square", drawingAPI, x, y, sideLength, sideLength),
		sideLength: sideLength,
	}
}

// Draw renders the square using the associated drawing API.
func (s *Square) Draw() string {
	return s.drawingAPI.DrawRectangle(s.x, s.y, s.x+s.sideLength, s.y+s.sideLength)
}

// ResizeTo resizes the square (uses the smallest of width and height to maintain square).
func (s *Square) ResizeTo(width, height float64) Shape {
	// Use the smaller dimension to maintain a square
	s.sideLength = min(width, height)
	s.width = s.sideLength
	s.height = s.sideLength
	return s
}

// MoveTo moves the square to the specified top-left coordinates.
func (s *Square) MoveTo(x, y float64) Shape {
	s.x = x
	s.y = y
	return s
}

// Rectangle is a concrete Shape implementation that represents a rectangle.
type Rectangle struct {
	BaseShape
}

// NewRectangle creates a new Rectangle with the given drawing API, top-left coordinates, width, and height.
func NewRectangle(drawingAPI DrawingAPI, x, y, width, height float64) *Rectangle {
	return &Rectangle{
		BaseShape: NewBaseShape("Rectangle", drawingAPI, x, y, width, height),
	}
}

// Draw renders the rectangle using the associated drawing API.
func (r *Rectangle) Draw() string {
	return r.drawingAPI.DrawRectangle(r.x, r.y, r.x+r.width, r.y+r.height)
}

// ResizeTo resizes the rectangle to the specified width and height.
func (r *Rectangle) ResizeTo(width, height float64) Shape {
	r.width = width
	r.height = height
	return r
}

// MoveTo moves the rectangle to the specified top-left coordinates.
func (r *Rectangle) MoveTo(x, y float64) Shape {
	r.x = x
	r.y = y
	return r
}

// Triangle is a concrete Shape implementation that represents a triangle.
type Triangle struct {
	BaseShape
	x2, y2, x3, y3 float64
}

// NewTriangle creates a new Triangle with the given drawing API and three points.
func NewTriangle(drawingAPI DrawingAPI, x1, y1, x2, y2, x3, y3 float64) *Triangle {
	// Calculate width and height based on the bounding box
	minX := min(min(x1, x2), x3)
	maxX := max(max(x1, x2), x3)
	minY := min(min(y1, y2), y3)
	maxY := max(max(y1, y2), y3)
	
	return &Triangle{
		BaseShape: NewBaseShape("Triangle", drawingAPI, x1, y1, maxX-minX, maxY-minY),
		x2:        x2,
		y2:        y2,
		x3:        x3,
		y3:        y3,
	}
}

// Draw renders the triangle using the associated drawing API.
func (t *Triangle) Draw() string {
	return t.drawingAPI.DrawTriangle(t.x, t.y, t.x2, t.y2, t.x3, t.y3)
}

// ResizeTo resizes the triangle while maintaining its proportions.
func (t *Triangle) ResizeTo(width, height float64) Shape {
	// Calculate current bounding box
	minX := min(min(t.x, t.x2), t.x3)
	maxX := max(max(t.x, t.x2), t.x3)
	minY := min(min(t.y, t.y2), t.y3)
	maxY := max(max(t.y, t.y2), t.y3)
	
	// Calculate scale factors
	scaleX := width / (maxX - minX)
	scaleY := height / (maxY - minY)
	
	// Save the original coordinates
	x1 := t.x
	y1 := t.y
	
	// Scale the coordinates relative to the first point
	t.x2 = x1 + (t.x2 - x1) * scaleX
	t.y2 = y1 + (t.y2 - y1) * scaleY
	t.x3 = x1 + (t.x3 - x1) * scaleX
	t.y3 = y1 + (t.y3 - y1) * scaleY
	
	// Update the width and height
	t.width = width
	t.height = height
	
	return t
}

// MoveTo moves the triangle to the specified coordinates, preserving its shape.
func (t *Triangle) MoveTo(x, y float64) Shape {
	// Calculate the offset
	dx := x - t.x
	dy := y - t.y
	
	// Apply the offset to all points
	t.x = x
	t.y = y
	t.x2 += dx
	t.y2 += dy
	t.x3 += dx
	t.y3 += dy
	
	return t
}

// Text is a concrete Shape implementation that represents text.
type Text struct {
	BaseShape
	content string
}

// NewText creates a new Text with the given drawing API, position, and content.
func NewText(drawingAPI DrawingAPI, x, y float64, content string) *Text {
	// Approximate width and height based on character count
	width := float64(len(content) * 10)
	height := 20.0
	
	return &Text{
		BaseShape: NewBaseShape("Text", drawingAPI, x, y, width, height),
		content:   content,
	}
}

// Draw renders the text using the associated drawing API.
func (t *Text) Draw() string {
	return t.drawingAPI.DrawText(t.x, t.y, t.content)
}

// ResizeTo does not actually resize the text, but updates internal dimensions.
func (t *Text) ResizeTo(width, height float64) Shape {
	t.width = width
	t.height = height
	return t
}

// MoveTo moves the text to the specified coordinates.
func (t *Text) MoveTo(x, y float64) Shape {
	t.x = x
	t.y = y
	return t
}

// GetContent returns the text content.
func (t *Text) GetContent() string {
	return t.content
}

// SetContent sets the text content.
func (t *Text) SetContent(content string) {
	t.content = content
	// Update width based on the new content length
	t.width = float64(len(content) * 10)
}

// Helper functions for min and max
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
