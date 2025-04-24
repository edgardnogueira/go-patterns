package bridge

import (
	"strings"
)

// DrawingApp demonstrates how the Bridge Pattern allows shapes and drawing
// implementations to vary independently.
type DrawingApp struct {
	shapes []Shape
}

// NewDrawingApp creates a new DrawingApp instance.
func NewDrawingApp() *DrawingApp {
	return &DrawingApp{
		shapes: make([]Shape, 0),
	}
}

// AddShape adds a shape to the drawing app.
func (d *DrawingApp) AddShape(shape Shape) {
	d.shapes = append(d.shapes, shape)
}

// ClearShapes removes all shapes from the drawing app.
func (d *DrawingApp) ClearShapes() {
	d.shapes = make([]Shape, 0)
}

// GetShapes returns all shapes in the drawing app.
func (d *DrawingApp) GetShapes() []Shape {
	return d.shapes
}

// Draw renders all shapes in the drawing app.
func (d *DrawingApp) Draw() string {
	var result strings.Builder
	
	for _, shape := range d.shapes {
		result.WriteString(shape.Draw())
		result.WriteString("\n")
	}
	
	return result.String()
}

// CreateDefaultScene creates a default scene with various shapes and renderers.
func (d *DrawingApp) CreateDefaultScene() {
	// Create different types of renderers
	vectorRenderer := NewVectorRenderer()
	rasterRenderer := NewRasterRenderer()
	svgRenderer := NewSVGRenderer()
	textRenderer := NewTextRenderer()
	
	// Add a circle with vector renderer
	d.AddShape(NewCircle(vectorRenderer, 100, 100, 50))
	
	// Add a square with raster renderer
	d.AddShape(NewSquare(rasterRenderer, 200, 200, 80))
	
	// Add a rectangle with SVG renderer
	d.AddShape(NewRectangle(svgRenderer, 300, 300, 120, 60))
	
	// Add a triangle with text renderer
	d.AddShape(NewTriangle(textRenderer, 400, 400, 450, 450, 400, 450))
	
	// Add text with various renderers
	d.AddShape(NewText(vectorRenderer, 100, 500, "Vector Text Example"))
	d.AddShape(NewText(svgRenderer, 300, 500, "SVG Text Example"))
}

// CreateShapeCollection creates a collection of the same shape with different renderers
// to demonstrate the benefit of the Bridge Pattern.
func (d *DrawingApp) CreateShapeCollection(shapeType string, x, y float64) {
	// Create different types of renderers
	vectorRenderer := NewVectorRenderer()
	rasterRenderer := NewRasterRenderer()
	svgRenderer := NewSVGRenderer()
	textRenderer := NewTextRenderer()
	
	renderers := []DrawingAPI{vectorRenderer, rasterRenderer, svgRenderer, textRenderer}
	yOffset := 0.0
	
	for _, renderer := range renderers {
		var shape Shape
		
		// Create the appropriate shape type with the current renderer
		switch strings.ToLower(shapeType) {
		case "circle":
			shape = NewCircle(renderer, x, y+yOffset, 30)
		case "square":
			shape = NewSquare(renderer, x, y+yOffset, 40)
		case "rectangle":
			shape = NewRectangle(renderer, x, y+yOffset, 80, 40)
		case "triangle":
			shape = NewTriangle(renderer, x, y+yOffset, x+40, y+yOffset-40, x+80, y+yOffset)
		case "text":
			shape = NewText(renderer, x, y+yOffset, "Hello, Bridge Pattern!")
		default:
			// Default to a circle if the shape type is unknown
			shape = NewCircle(renderer, x, y+yOffset, 30)
		}
		
		d.AddShape(shape)
		yOffset += 100 // Offset for the next shape
	}
}

// ChangeAllRenderers changes the renderer for all shapes to the specified type.
// This demonstrates the flexibility of the Bridge Pattern by changing the
// implementation (renderer) without changing the abstraction (shapes).
func (d *DrawingApp) ChangeAllRenderers(rendererType string) {
	var renderer DrawingAPI
	
	// Create the appropriate renderer
	switch strings.ToLower(rendererType) {
	case "vector":
		renderer = NewVectorRenderer()
	case "raster":
		renderer = NewRasterRenderer()
	case "svg":
		renderer = NewSVGRenderer()
	case "text", "ascii":
		renderer = NewTextRenderer()
	default:
		// Default to vector if the renderer type is unknown
		renderer = NewVectorRenderer()
	}
	
	// Create new shapes with the same properties but different renderer
	newShapes := make([]Shape, 0, len(d.shapes))
	
	for _, shape := range d.shapes {
		var newShape Shape
		
		switch s := shape.(type) {
		case *Circle:
			newShape = NewCircle(renderer, s.x, s.y, s.radius)
		case *Square:
			newShape = NewSquare(renderer, s.x, s.y, s.sideLength)
		case *Rectangle:
			newShape = NewRectangle(renderer, s.x, s.y, s.width, s.height)
		case *Triangle:
			newShape = NewTriangle(renderer, s.x, s.y, s.x2, s.y2, s.x3, s.y3)
		case *Text:
			newShape = NewText(renderer, s.x, s.y, s.content)
		}
		
		if newShape != nil {
			newShapes = append(newShapes, newShape)
		}
	}
	
	// Replace the shapes
	d.shapes = newShapes
}

// GetShapesByRenderer returns a list of shapes that use the specified renderer type.
func (d *DrawingApp) GetShapesByRenderer(rendererType string) []Shape {
	filteredShapes := make([]Shape, 0)
	
	for _, shape := range d.shapes {
		if strings.EqualFold(shape.GetDrawingAPI().GetName(), rendererType) {
			filteredShapes = append(filteredShapes, shape)
		}
	}
	
	return filteredShapes
}

// GetShapeCount returns the total number of shapes.
func (d *DrawingApp) GetShapeCount() int {
	return len(d.shapes)
}
