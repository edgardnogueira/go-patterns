package bridge

import (
	"strings"
	"testing"
)

func TestVectorRenderer(t *testing.T) {
	renderer := NewVectorRenderer()
	
	tests := []struct {
		name     string
		method   func() string
		expected string
	}{
		{
			name: "DrawCircle",
			method: func() string {
				return renderer.DrawCircle(10, 20, 30)
			},
			expected: "Vector circle at (10.0,20.0) with radius 30.0",
		},
		{
			name: "DrawRectangle",
			method: func() string {
				return renderer.DrawRectangle(10, 20, 30, 40)
			},
			expected: "Vector rectangle at (10.0,20.0)-(30.0,40.0)",
		},
		{
			name: "DrawTriangle",
			method: func() string {
				return renderer.DrawTriangle(10, 20, 30, 40, 50, 60)
			},
			expected: "Vector triangle at (10.0,20.0), (30.0,40.0), (50.0,60.0)",
		},
		{
			name: "DrawLine",
			method: func() string {
				return renderer.DrawLine(10, 20, 30, 40)
			},
			expected: "Vector line from (10.0,20.0) to (30.0,40.0)",
		},
		{
			name: "DrawText",
			method: func() string {
				return renderer.DrawText(10, 20, "Test")
			},
			expected: "Vector text 'Test' at (10.0,20.0)",
		},
		{
			name: "GetName",
			method: func() string {
				return renderer.GetName()
			},
			expected: "Vector",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.method()
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSVGRenderer(t *testing.T) {
	renderer := NewSVGRenderer()
	
	tests := []struct {
		name     string
		method   func() string
		expected string
	}{
		{
			name: "DrawCircle",
			method: func() string {
				return renderer.DrawCircle(10, 20, 30)
			},
			expected: `<circle cx="10.0" cy="20.0" r="30.0" />`,
		},
		{
			name: "DrawRectangle",
			method: func() string {
				return renderer.DrawRectangle(10, 20, 40, 50)
			},
			expected: `<rect x="10.0" y="20.0" width="30.0" height="30.0" />`,
		},
		{
			name: "DrawTriangle",
			method: func() string {
				return renderer.DrawTriangle(10, 20, 30, 40, 50, 60)
			},
			expected: `<polygon points="10.0,20.0 30.0,40.0 50.0,60.0" />`,
		},
		{
			name: "DrawLine",
			method: func() string {
				return renderer.DrawLine(10, 20, 30, 40)
			},
			expected: `<line x1="10.0" y1="20.0" x2="30.0" y2="40.0" />`,
		},
		{
			name: "DrawText",
			method: func() string {
				return renderer.DrawText(10, 20, "Test")
			},
			expected: `<text x="10.0" y="20.0">Test</text>`,
		},
		{
			name: "GetName",
			method: func() string {
				return renderer.GetName()
			},
			expected: "SVG",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.method()
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestCircle(t *testing.T) {
	renderer := NewVectorRenderer()
	circle := NewCircle(renderer, 10, 20, 30)
	
	// Test Draw method
	expectedDraw := "Vector circle at (10.0,20.0) with radius 30.0"
	if result := circle.Draw(); result != expectedDraw {
		t.Errorf("Circle.Draw() = %q, want %q", result, expectedDraw)
	}
	
	// Test ResizeTo method
	circle.ResizeTo(100, 80)
	expectedRadius := 45.0 // (100+80)/2/2
	if circle.radius != expectedRadius {
		t.Errorf("Circle.radius = %f, want %f", circle.radius, expectedRadius)
	}
	
	// Test MoveTo method
	circle.MoveTo(50, 60)
	if circle.x != 50 || circle.y != 60 {
		t.Errorf("Circle position = (%f,%f), want (50,60)", circle.x, circle.y)
	}
	
	// Test GetName and GetDescription
	if name := circle.GetName(); name != "Circle" {
		t.Errorf("Circle.GetName() = %q, want %q", name, "Circle")
	}
	
	expectedDesc := "Circle rendered using Vector renderer"
	if desc := circle.GetDescription(); desc != expectedDesc {
		t.Errorf("Circle.GetDescription() = %q, want %q", desc, expectedDesc)
	}
}

func TestRectangle(t *testing.T) {
	renderer := NewRasterRenderer()
	rectangle := NewRectangle(renderer, 10, 20, 30, 40)
	
	// Test Draw method
	expectedDraw := "Raster rectangle at (10.0,20.0)-(40.0,60.0)"
	if result := rectangle.Draw(); result != expectedDraw {
		t.Errorf("Rectangle.Draw() = %q, want %q", result, expectedDraw)
	}
	
	// Test ResizeTo method
	rectangle.ResizeTo(50, 60)
	if rectangle.width != 50 || rectangle.height != 60 {
		t.Errorf("Rectangle dimensions = (%f,%f), want (50,60)", rectangle.width, rectangle.height)
	}
	
	// Test MoveTo method
	rectangle.MoveTo(50, 60)
	if rectangle.x != 50 || rectangle.y != 60 {
		t.Errorf("Rectangle position = (%f,%f), want (50,60)", rectangle.x, rectangle.y)
	}
}

func TestDrawingApp(t *testing.T) {
	app := NewDrawingApp()
	
	// Test empty drawing
	if count := app.GetShapeCount(); count != 0 {
		t.Errorf("Initial shape count = %d, want 0", count)
	}
	
	// Add shapes and test count
	vectorRenderer := NewVectorRenderer()
	svgRenderer := NewSVGRenderer()
	
	app.AddShape(NewCircle(vectorRenderer, 10, 20, 30))
	app.AddShape(NewRectangle(svgRenderer, 40, 50, 60, 70))
	
	if count := app.GetShapeCount(); count != 2 {
		t.Errorf("Shape count after adding = %d, want 2", count)
	}
	
	// Test GetShapesByRenderer
	vectorShapes := app.GetShapesByRenderer("Vector")
	if len(vectorShapes) != 1 {
		t.Errorf("Vector shapes count = %d, want 1", len(vectorShapes))
	}
	
	svgShapes := app.GetShapesByRenderer("SVG")
	if len(svgShapes) != 1 {
		t.Errorf("SVG shapes count = %d, want 1", len(svgShapes))
	}
	
	// Test Draw method
	drawing := app.Draw()
	if !strings.Contains(drawing, "Vector circle") || !strings.Contains(drawing, "<rect") {
		t.Errorf("Drawing output doesn't contain expected content: %s", drawing)
	}
	
	// Test changing renderers
	app.ChangeAllRenderers("text")
	
	for _, shape := range app.GetShapes() {
		if renderer := shape.GetDrawingAPI().GetName(); renderer != "ASCII" {
			t.Errorf("Renderer after change = %q, want 'ASCII'", renderer)
		}
	}
	
	// Test ClearShapes
	app.ClearShapes()
	if count := app.GetShapeCount(); count != 0 {
		t.Errorf("Shape count after clearing = %d, want 0", count)
	}
	
	// Test CreateShapeCollection
	app.CreateShapeCollection("circle", 100, 100)
	if count := app.GetShapeCount(); count != 4 {
		t.Errorf("Shape count after CreateShapeCollection = %d, want 4", count)
	}
}

func TestShapeTransformations(t *testing.T) {
	renderer := NewVectorRenderer()
	
	// Test Square resize
	square := NewSquare(renderer, 10, 20, 30)
	square.ResizeTo(40, 50) // For a square, it should take the smaller dimension
	expectedSide := 40.0    // min(40, 50)
	if square.sideLength != expectedSide {
		t.Errorf("Square.sideLength after resize = %f, want %f", square.sideLength, expectedSide)
	}
	
	// Test Triangle move
	triangle := NewTriangle(renderer, 10, 20, 30, 40, 50, 60)
	triangle.MoveTo(100, 200)
	
	// The relative positions should be maintained
	expectedX2 := 120.0 // 100 + (30 - 10)
	expectedY2 := 220.0 // 200 + (40 - 20)
	expectedX3 := 140.0 // 100 + (50 - 10)
	expectedY3 := 240.0 // 200 + (60 - 20)
	
	if triangle.x2 != expectedX2 || triangle.y2 != expectedY2 ||
	   triangle.x3 != expectedX3 || triangle.y3 != expectedY3 {
		t.Errorf("Triangle points after move = (%f,%f), (%f,%f), want (%f,%f), (%f,%f)",
			triangle.x2, triangle.y2, triangle.x3, triangle.y3,
			expectedX2, expectedY2, expectedX3, expectedY3)
	}
	
	// Test Text content
	text := NewText(renderer, 10, 20, "Hello")
	if content := text.GetContent(); content != "Hello" {
		t.Errorf("Text.GetContent() = %q, want %q", content, "Hello")
	}
	
	text.SetContent("Updated Text")
	if content := text.GetContent(); content != "Updated Text" {
		t.Errorf("Text.GetContent() after update = %q, want %q", content, "Updated Text")
	}
	
	// Width should be updated based on content length
	expectedWidth := float64(len("Updated Text") * 10)
	if text.width != expectedWidth {
		t.Errorf("Text.width after content update = %f, want %f", text.width, expectedWidth)
	}
}
