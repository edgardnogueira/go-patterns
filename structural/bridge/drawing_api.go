package bridge

// DrawingAPI is the Implementor interface in the Bridge pattern.
// It defines the interface for the implementation classes that do the actual drawing.
type DrawingAPI interface {
	// DrawCircle draws a circle with the given center (x, y) and radius.
	DrawCircle(x, y, radius float64) string

	// DrawRectangle draws a rectangle with the given coordinates.
	DrawRectangle(x1, y1, x2, y2 float64) string

	// DrawTriangle draws a triangle with the given points.
	DrawTriangle(x1, y1, x2, y2, x3, y3 float64) string

	// DrawLine draws a line from (x1, y1) to (x2, y2).
	DrawLine(x1, y1, x2, y2 float64) string

	// DrawText draws text at the specified coordinates.
	DrawText(x, y float64, text string) string

	// GetName returns the name of the renderer.
	GetName() string
}

// VectorRenderer is a ConcreteImplementor in the Bridge pattern.
// It implements the DrawingAPI interface for vector graphics.
type VectorRenderer struct{}

// NewVectorRenderer creates a new VectorRenderer instance.
func NewVectorRenderer() *VectorRenderer {
	return &VectorRenderer{}
}

// DrawCircle draws a circle using vector graphics.
func (v *VectorRenderer) DrawCircle(x, y, radius float64) string {
	return fmt.Sprintf("Vector circle at (%.1f,%.1f) with radius %.1f", x, y, radius)
}

// DrawRectangle draws a rectangle using vector graphics.
func (v *VectorRenderer) DrawRectangle(x1, y1, x2, y2 float64) string {
	return fmt.Sprintf("Vector rectangle at (%.1f,%.1f)-(%.1f,%.1f)", x1, y1, x2, y2)
}

// DrawTriangle draws a triangle using vector graphics.
func (v *VectorRenderer) DrawTriangle(x1, y1, x2, y2, x3, y3 float64) string {
	return fmt.Sprintf("Vector triangle at (%.1f,%.1f), (%.1f,%.1f), (%.1f,%.1f)", x1, y1, x2, y2, x3, y3)
}

// DrawLine draws a line using vector graphics.
func (v *VectorRenderer) DrawLine(x1, y1, x2, y2 float64) string {
	return fmt.Sprintf("Vector line from (%.1f,%.1f) to (%.1f,%.1f)", x1, y1, x2, y2)
}

// DrawText draws text using vector graphics.
func (v *VectorRenderer) DrawText(x, y float64, text string) string {
	return fmt.Sprintf("Vector text '%s' at (%.1f,%.1f)", text, x, y)
}

// GetName returns the name of the renderer.
func (v *VectorRenderer) GetName() string {
	return "Vector"
}

// RasterRenderer is a ConcreteImplementor in the Bridge pattern.
// It implements the DrawingAPI interface for raster (pixel-based) graphics.
type RasterRenderer struct{}

// NewRasterRenderer creates a new RasterRenderer instance.
func NewRasterRenderer() *RasterRenderer {
	return &RasterRenderer{}
}

// DrawCircle draws a circle using raster graphics.
func (r *RasterRenderer) DrawCircle(x, y, radius float64) string {
	return fmt.Sprintf("Raster circle at (%.1f,%.1f) with radius %.1f", x, y, radius)
}

// DrawRectangle draws a rectangle using raster graphics.
func (r *RasterRenderer) DrawRectangle(x1, y1, x2, y2 float64) string {
	return fmt.Sprintf("Raster rectangle at (%.1f,%.1f)-(%.1f,%.1f)", x1, y1, x2, y2)
}

// DrawTriangle draws a triangle using raster graphics.
func (r *RasterRenderer) DrawTriangle(x1, y1, x2, y2, x3, y3 float64) string {
	return fmt.Sprintf("Raster triangle at (%.1f,%.1f), (%.1f,%.1f), (%.1f,%.1f)", x1, y1, x2, y2, x3, y3)
}

// DrawLine draws a line using raster graphics.
func (r *RasterRenderer) DrawLine(x1, y1, x2, y2 float64) string {
	return fmt.Sprintf("Raster line from (%.1f,%.1f) to (%.1f,%.1f)", x1, y1, x2, y2)
}

// DrawText draws text using raster graphics.
func (r *RasterRenderer) DrawText(x, y float64, text string) string {
	return fmt.Sprintf("Raster text '%s' at (%.1f,%.1f)", text, x, y)
}

// GetName returns the name of the renderer.
func (r *RasterRenderer) GetName() string {
	return "Raster"
}

// SVGRenderer is a ConcreteImplementor in the Bridge pattern.
// It implements the DrawingAPI interface for SVG graphics.
type SVGRenderer struct{}

// NewSVGRenderer creates a new SVGRenderer instance.
func NewSVGRenderer() *SVGRenderer {
	return &SVGRenderer{}
}

// DrawCircle draws a circle using SVG commands.
func (s *SVGRenderer) DrawCircle(x, y, radius float64) string {
	return fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="%.1f" />`, x, y, radius)
}

// DrawRectangle draws a rectangle using SVG commands.
func (s *SVGRenderer) DrawRectangle(x1, y1, x2, y2 float64) string {
	width := x2 - x1
	height := y2 - y1
	return fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" />`, x1, y1, width, height)
}

// DrawTriangle draws a triangle using SVG commands.
func (s *SVGRenderer) DrawTriangle(x1, y1, x2, y2, x3, y3 float64) string {
	return fmt.Sprintf(`<polygon points="%.1f,%.1f %.1f,%.1f %.1f,%.1f" />`, x1, y1, x2, y2, x3, y3)
}

// DrawLine draws a line using SVG commands.
func (s *SVGRenderer) DrawLine(x1, y1, x2, y2 float64) string {
	return fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" />`, x1, y1, x2, y2)
}

// DrawText draws text using SVG commands.
func (s *SVGRenderer) DrawText(x, y float64, text string) string {
	return fmt.Sprintf(`<text x="%.1f" y="%.1f">%s</text>`, x, y, text)
}

// GetName returns the name of the renderer.
func (s *SVGRenderer) GetName() string {
	return "SVG"
}

// TextRenderer is a ConcreteImplementor in the Bridge pattern.
// It implements the DrawingAPI interface using ASCII art.
type TextRenderer struct{}

// NewTextRenderer creates a new TextRenderer instance.
func NewTextRenderer() *TextRenderer {
	return &TextRenderer{}
}

// DrawCircle draws a circle using ASCII art.
func (t *TextRenderer) DrawCircle(x, y, radius float64) string {
	return fmt.Sprintf("O (ASCII circle at %.1f,%.1f with radius %.1f)", x, y, radius)
}

// DrawRectangle draws a rectangle using ASCII art.
func (t *TextRenderer) DrawRectangle(x1, y1, x2, y2 float64) string {
	return fmt.Sprintf("[ ] (ASCII rectangle at %.1f,%.1f-%.1f,%.1f)", x1, y1, x2, y2)
}

// DrawTriangle draws a triangle using ASCII art.
func (t *TextRenderer) DrawTriangle(x1, y1, x2, y2, x3, y3 float64) string {
	return fmt.Sprintf("^ (ASCII triangle at %.1f,%.1f, %.1f,%.1f, %.1f,%.1f)", x1, y1, x2, y2, x3, y3)
}

// DrawLine draws a line using ASCII art.
func (t *TextRenderer) DrawLine(x1, y1, x2, y2 float64) string {
	return fmt.Sprintf("--- (ASCII line from %.1f,%.1f to %.1f,%.1f)", x1, y1, x2, y2)
}

// DrawText draws text using ASCII art.
func (t *TextRenderer) DrawText(x, y float64, text string) string {
	return fmt.Sprintf("[%s] (ASCII text at %.1f,%.1f)", text, x, y)
}

// GetName returns the name of the renderer.
func (t *TextRenderer) GetName() string {
	return "ASCII"
}
