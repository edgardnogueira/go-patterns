package main

import (
	"fmt"
	"github.com/edgardnogueira/go-patterns/structural/bridge"
	"strings"
)

func main() {
	fmt.Println("Bridge Pattern Example - Drawing Application")
	fmt.Println("===========================================")
	
	// Create a new drawing app
	app := bridge.NewDrawingApp()
	
	// Create different renderers
	vectorRenderer := bridge.NewVectorRenderer()
	rasterRenderer := bridge.NewRasterRenderer()
	svgRenderer := bridge.NewSVGRenderer()
	textRenderer := bridge.NewTextRenderer()
	
	fmt.Println("\n1. Creating shapes with different renderers:")
	fmt.Println("-------------------------------------------")
	
	// Create shapes with different renderers
	circle := bridge.NewCircle(vectorRenderer, 100, 100, 50)
	square := bridge.NewSquare(rasterRenderer, 200, 200, 80)
	rectangle := bridge.NewRectangle(svgRenderer, 300, 300, 120, 60)
	triangle := bridge.NewTriangle(textRenderer, 400, 400, 450, 450, 400, 450)
	
	// Add shapes to the app
	app.AddShape(circle)
	app.AddShape(square)
	app.AddShape(rectangle)
	app.AddShape(triangle)
	
	// Draw all shapes
	fmt.Println(app.Draw())
	
	fmt.Println("\n2. Demonstrating shape operations:")
	fmt.Println("--------------------------------")
	
	// Resize circle
	circle.ResizeTo(75, 75)
	fmt.Printf("Resized circle: %s\n", circle.Draw())
	
	// Move square
	square.MoveTo(250, 250)
	fmt.Printf("Moved square: %s\n", square.Draw())
	
	fmt.Println("\n3. Creating shape collections with the same shape type but different renderers:")
	fmt.Println("----------------------------------------------------------------------------")
	
	// Clear current shapes
	app.ClearShapes()
	
	// Create collection of circles with different renderers
	app.CreateShapeCollection("circle", 150, 150)
	
	// Draw all shapes
	fmt.Println(app.Draw())
	
	fmt.Println("\n4. Showing shape descriptions:")
	fmt.Println("----------------------------")
	for _, shape := range app.GetShapes() {
		fmt.Println(shape.GetDescription())
	}
	
	fmt.Println("\n5. Changing all renderers to SVG:")
	fmt.Println("------------------------------")
	
	// Change all renderers to SVG
	app.ChangeAllRenderers("svg")
	
	// Draw all shapes
	fmt.Println(app.Draw())
	
	fmt.Println("\n6. Creating a default scene:")
	fmt.Println("-------------------------")
	
	// Clear current shapes
	app.ClearShapes()
	
	// Create default scene
	app.CreateDefaultScene()
	
	// Draw all shapes
	fmt.Println(app.Draw())
	
	fmt.Println("\n7. Filter shapes by renderer type:")
	fmt.Println("-------------------------------")
	
	// Get shapes by renderer
	vectorShapes := app.GetShapesByRenderer("Vector")
	fmt.Printf("Vector shapes (%d):\n", len(vectorShapes))
	for _, shape := range vectorShapes {
		fmt.Printf("- %s: %s\n", shape.GetName(), shape.Draw())
	}
	
	svgShapes := app.GetShapesByRenderer("SVG")
	fmt.Printf("\nSVG shapes (%d):\n", len(svgShapes))
	for _, shape := range svgShapes {
		fmt.Printf("- %s: %s\n", shape.GetName(), shape.Draw())
	}
	
	fmt.Println("\nBridge Pattern Benefits:")
	fmt.Println("----------------------")
	fmt.Println("1. Separation of abstraction (shapes) from implementation (renderers)")
	fmt.Println("2. Both hierarchies can evolve independently")
	fmt.Println("3. Shapes can change their renderer dynamically")
	fmt.Println("4. New shapes and new renderers can be added without affecting existing code")
}
