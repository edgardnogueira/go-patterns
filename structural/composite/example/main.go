package main

import (
	"fmt"
	"github.com/edgardnogueira/go-patterns/structural/composite"
)

func main() {
	fmt.Println("Composite Pattern Example - File System")
	fmt.Println("======================================")
	
	// Create a file system manager
	fsm := composite.NewFileSystemManager()
	
	// Initialize a sample file system structure
	createSampleFileSystem(fsm)
	
	// Show the file system structure
	fmt.Println("\n1. File System Structure:")
	fmt.Println("------------------------")
	fmt.Println(fsm.PrintFileSystem())
	
	// Demonstrate finding nodes by path
	fmt.Println("\n2. Finding Nodes by Path:")
	fmt.Println("------------------------")
	findAndPrintNode(fsm, "/documents/projects/project1/README.md")
	findAndPrintNode(fsm, "/documents/notes")
	findAndPrintNode(fsm, "/images")
	
	// Demonstrate operations on nodes
	fmt.Println("\n3. Performing Operations:")
	fmt.Println("------------------------")
	
	// Create a new file
	fmt.Println("Creating a new file:")
	_, err := fsm.CreateFile("/documents/notes/todo.txt", []byte("1. Finish the Composite Pattern\n2. Implement more patterns"))
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
	} else {
		findAndPrintNode(fsm, "/documents/notes/todo.txt")
	}
	
	// Move a file
	fmt.Println("\nMoving a file:")
	err = fsm.MoveNode("/documents/notes/todo.txt", "/todo.txt")
	if err != nil {
		fmt.Printf("Error moving file: %v\n", err)
	} else {
		findAndPrintNode(fsm, "/todo.txt")
	}
	
	// Copy a directory
	fmt.Println("\nCopying a directory:")
	err = fsm.CopyNode("/documents/projects/project1", "/documents/projects/project1_backup")
	if err != nil {
		fmt.Printf("Error copying directory: %v\n", err)
	} else {
		findAndPrintNode(fsm, "/documents/projects/project1_backup")
	}
	
	// Delete a node
	fmt.Println("\nDeleting a node:")
	err = fsm.DeleteNode("/images/logo.png")
	if err != nil {
		fmt.Printf("Error deleting node: %v\n", err)
	} else {
		fmt.Println("Deleted /images/logo.png successfully")
	}
	
	// Demonstrate using visitors
	fmt.Println("\n4. Using Visitors:")
	fmt.Println("----------------")
	
	// Calculate statistics
	fmt.Println("File System Statistics:")
	stats, err := fsm.CalculateStatistics()
	if err != nil {
		fmt.Printf("Error calculating statistics: %v\n", err)
	} else {
		fmt.Println(stats.PrintReport())
	}
	
	// Search for files
	fmt.Println("\nSearching for files containing 'README':")
	results, err := fsm.Search("README", "", false, 0, 0, 0)
	if err != nil {
		fmt.Printf("Error searching: %v\n", err)
	} else {
		fmt.Printf("Found %d results:\n", len(results))
		for i, node := range results {
			fmt.Printf("  %d. %s (%s)\n", i+1, node.Path(), nodeType(node))
		}
	}
	
	// Calculate size of specific file types
	fmt.Println("\nCalculating size of all .md files:")
	visitor := composite.NewSizeCalculatorVisitor(".md", 0, 0, false)
	err = fsm.ApplyVisitor("/", visitor)
	if err != nil {
		fmt.Printf("Error applying visitor: %v\n", err)
	} else {
		fmt.Printf("Total size of .md files: %d bytes\n", visitor.TotalSize)
	}
	
	// Update permissions
	fmt.Println("\nUpdating permissions for all .txt files:")
	modified, err := fsm.UpdatePermissions("/", composite.Execute, 0, true, false, true, ".txt")
	if err != nil {
		fmt.Printf("Error updating permissions: %v\n", err)
	} else {
		fmt.Printf("Modified permissions for %d files\n", modified)
	}
	
	// Show final file system structure
	fmt.Println("\n5. Final File System Structure:")
	fmt.Println("------------------------------")
	fmt.Println(fsm.PrintFileSystem())
	
	fmt.Println("\nComposite Pattern Benefits:")
	fmt.Println("-------------------------")
	fmt.Println("1. Treats individual objects and compositions uniformly")
	fmt.Println("2. Makes it easy to add new types of components")
	fmt.Println("3. Simplifies client code by avoiding type checks and special case handling")
	fmt.Println("4. Enables recursive operations across the entire structure")
}

// createSampleFileSystem initializes a sample file system structure for the example
func createSampleFileSystem(fsm *composite.FileSystemManager) {
	// Create a sample file system structure
	
	// Create directories
	fsm.CreateDirectoryPath("/documents/projects/project1")
	fsm.CreateDirectoryPath("/documents/projects/project2")
	fsm.CreateDirectory("/documents/notes")
	fsm.CreateDirectory("/images")
	fsm.CreateDirectory("/config")
	
	// Create some files with content
	fsm.CreateFile("/documents/projects/project1/README.md", []byte("# Project 1\nThis is a sample project."))
	fsm.CreateFile("/documents/projects/project1/main.go", []byte("package main\n\nfunc main() {\n\tfmt.Println(\"Hello, Project 1\")\n}"))
	fsm.CreateFile("/documents/projects/project2/README.md", []byte("# Project 2\nAnother sample project."))
	fsm.CreateFile("/documents/notes/meeting.txt", []byte("Meeting notes from 2025-04-10:\n- Discuss design patterns\n- Plan implementation"))
	fsm.CreateFile("/images/logo.png", []byte("Binary image content"))
	fsm.CreateFile("/images/background.jpg", []byte("Binary image content"))
	fsm.CreateFile("/config/settings.json", []byte("{\n  \"debug\": true,\n  \"theme\": \"dark\"\n}"))
	
	// Create a hidden file
	fsm.CreateFile("/.hidden", []byte("This is a hidden file"))
}

// findAndPrintNode finds a node by path and prints its information
func findAndPrintNode(fsm *composite.FileSystemManager, path string) {
	node, err := fsm.FindNode(path)
	if err != nil {
		fmt.Printf("Error finding %s: %v\n", path, err)
		return
	}
	
	fmt.Printf("Found %s (%s):\n", path, nodeType(node))
	fmt.Printf("  Name: %s\n", node.Name())
	fmt.Printf("  Size: %d bytes\n", node.Size())
	
	// If it's a file, show some content
	if !node.IsDirectory() {
		if file, ok := node.(*composite.File); ok {
			content := string(file.Content())
			if len(content) > 50 {
				content = content[:50] + "..."
			}
			fmt.Printf("  Content: %s\n", content)
		}
	}
}

// nodeType returns a string description of the node type
func nodeType(node composite.FileSystemNode) string {
	if node.IsDirectory() {
		return "directory"
	}
	return "file"
}
