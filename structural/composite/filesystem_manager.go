package composite

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

// FileSystemManager provides utility functions for working with a file system.
// It uses the Composite pattern for file system operations.
type FileSystemManager struct {
	root *Directory
}

// NewFileSystemManager creates a new file system manager with a root directory.
func NewFileSystemManager() *FileSystemManager {
	return &FileSystemManager{
		root: NewDirectory("/"),
	}
}

// Root returns the root directory.
func (m *FileSystemManager) Root() *Directory {
	return m.root
}

// CreateFile creates a new file at the specified path.
func (m *FileSystemManager) CreateFile(path string, content []byte) (*File, error) {
	dir, fileName := filepath.Split(path)
	
	// Get or create parent directory
	parentDir, err := m.CreateDirectoryPath(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to create parent directory: %w", err)
	}
	
	// Check if a file or directory with this name already exists
	if existingNode := parentDir.GetChild(fileName); existingNode != nil {
		return nil, fmt.Errorf("a node named '%s' already exists", fileName)
	}
	
	// Create the file
	file := NewFile(fileName, content)
	parentDir.Add(file)
	
	return file, nil
}

// CreateDirectory creates a new directory at the specified path.
func (m *FileSystemManager) CreateDirectory(path string) (*Directory, error) {
	// Handle root directory
	if path == "/" || path == "" {
		return m.root, nil
	}
	
	dir, dirName := filepath.Split(path)
	
	// Remove trailing slash if any
	dirName = strings.TrimSuffix(dirName, "/")
	
	// If the directory name is empty, use the last component of the path
	if dirName == "" {
		dirName = filepath.Base(path)
		dir = filepath.Dir(path)
		if dir == "." {
			dir = "/"
		}
	}
	
	// Get or create parent directory
	parentDir, err := m.CreateDirectoryPath(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to create parent directory: %w", err)
	}
	
	// Check if a node with this name already exists
	if existingNode := parentDir.GetChild(dirName); existingNode != nil {
		// If it's a directory, return it
		if existingDir, ok := existingNode.(*Directory); ok {
			return existingDir, nil
		}
		return nil, fmt.Errorf("a file named '%s' already exists", dirName)
	}
	
	// Create the directory
	newDir := NewDirectory(dirName)
	parentDir.Add(newDir)
	
	return newDir, nil
}

// CreateDirectoryPath creates all directories in a path and returns the last one.
func (m *FileSystemManager) CreateDirectoryPath(path string) (*Directory, error) {
	// Handle root directory
	if path == "/" || path == "" || path == "." {
		return m.root, nil
	}
	
	// Clean and split the path
	cleanPath := filepath.Clean(path)
	components := strings.Split(cleanPath, string(filepath.Separator))
	
	// If the path is absolute and starts from root
	if components[0] == "" && len(components) > 1 {
		components = components[1:]
	}
	
	// Start from the root directory
	currentDir := m.root
	
	// Create each directory in the path
	for _, component := range components {
		// Skip empty components
		if component == "" || component == "." {
			continue
		}
		
		// Check if the directory already exists
		child := currentDir.GetChild(component)
		if child != nil {
			// If it's a directory, move to it
			if dir, ok := child.(*Directory); ok {
				currentDir = dir
			} else {
				return nil, fmt.Errorf("path component '%s' exists but is a file", component)
			}
		} else {
			// Create a new directory
			newDir := NewDirectory(component)
			currentDir.Add(newDir)
			currentDir = newDir
		}
	}
	
	return currentDir, nil
}

// FindNode finds a node by its path.
func (m *FileSystemManager) FindNode(path string) (FileSystemNode, error) {
	// Handle root directory
	if path == "/" || path == "" {
		return m.root, nil
	}
	
	// Find the node starting from the root
	node := m.root.FindByPath(path)
	if node == nil {
		return nil, fmt.Errorf("node not found: %s", path)
	}
	
	return node, nil
}

// DeleteNode deletes a node by its path.
func (m *FileSystemManager) DeleteNode(path string) error {
	// Cannot delete root
	if path == "/" || path == "" {
		return errors.New("cannot delete root directory")
	}
	
	// Get the node to delete
	node, err := m.FindNode(path)
	if err != nil {
		return err
	}
	
	// Get the parent directory
	parentPath := filepath.Dir(path)
	if parentPath == "." {
		parentPath = "/"
	}
	
	parentNode, err := m.FindNode(parentPath)
	if err != nil {
		return fmt.Errorf("parent directory not found: %w", err)
	}
	
	parentDir, ok := parentNode.(*Directory)
	if !ok {
		return fmt.Errorf("parent of '%s' is not a directory", path)
	}
	
	// Remove the node from its parent
	if !parentDir.Remove(node) {
		return fmt.Errorf("failed to remove node '%s'", path)
	}
	
	return nil
}

// MoveNode moves a node from one path to another.
func (m *FileSystemManager) MoveNode(sourcePath, destPath string) error {
	// Cannot move root
	if sourcePath == "/" || sourcePath == "" {
		return errors.New("cannot move root directory")
	}
	
	// Get the source node
	sourceNode, err := m.FindNode(sourcePath)
	if err != nil {
		return err
	}
	
	// Get the destination directory
	destDir, destName := filepath.Split(destPath)
	
	// Create parent directories if needed
	destParent, err := m.CreateDirectoryPath(destDir)
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}
	
	// Check if a node with the destination name already exists
	if existingNode := destParent.GetChild(destName); existingNode != nil {
		return fmt.Errorf("a node named '%s' already exists at destination", destName)
	}
	
	// Get the source parent directory
	sourceParentPath := filepath.Dir(sourcePath)
	if sourceParentPath == "." {
		sourceParentPath = "/"
	}
	
	sourceParentNode, err := m.FindNode(sourceParentPath)
	if err != nil {
		return fmt.Errorf("source parent directory not found: %w", err)
	}
	
	sourceParent, ok := sourceParentNode.(*Directory)
	if !ok {
		return fmt.Errorf("parent of '%s' is not a directory", sourcePath)
	}
	
	// Remove the node from its current parent
	if !sourceParent.Remove(sourceNode) {
		return fmt.Errorf("failed to remove node from source '%s'", sourcePath)
	}
	
	// Rename the node if necessary
	if sourceNode.Name() != destName {
		switch node := sourceNode.(type) {
		case *File:
			node.name = destName
		case *Directory:
			node.name = destName
		}
	}
	
	// Add the node to its new parent
	destParent.Add(sourceNode)
	
	return nil
}

// CopyNode copies a node from one path to another.
func (m *FileSystemManager) CopyNode(sourcePath, destPath string) error {
	// Get the source node
	sourceNode, err := m.FindNode(sourcePath)
	if err != nil {
		return err
	}
	
	// Get the destination directory
	destDir, destName := filepath.Split(destPath)
	
	// Create parent directories if needed
	destParent, err := m.CreateDirectoryPath(destDir)
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}
	
	// Check if a node with the destination name already exists
	if existingNode := destParent.GetChild(destName); existingNode != nil {
		return fmt.Errorf("a node named '%s' already exists at destination", destName)
	}
	
	// Create a copy of the node based on its type
	var newNode FileSystemNode
	
	switch node := sourceNode.(type) {
	case *File:
		// Copy file content
		contentCopy := make([]byte, len(node.content))
		copy(contentCopy, node.content)
		
		newFile := NewFile(destName, contentCopy)
		newFile.SetPermissions(node.GetPermissions())
		newNode = newFile
		
	case *Directory:
		// Create new directory
		newDir := NewDirectory(destName)
		newDir.SetPermissions(node.GetPermissions())
		newNode = newDir
		
		// Recursively copy children
		for _, child := range node.Children() {
			childSourcePath := filepath.Join(sourcePath, child.Name())
			childDestPath := filepath.Join(destPath, child.Name())
			
			if err := m.CopyNode(childSourcePath, childDestPath); err != nil {
				return fmt.Errorf("failed to copy child '%s': %w", child.Name(), err)
			}
		}
		
	default:
		return fmt.Errorf("unknown node type")
	}
	
	// Add the new node to the destination parent
	destParent.Add(newNode)
	
	return nil
}

// ApplyVisitor applies a visitor to a node at the specified path.
func (m *FileSystemManager) ApplyVisitor(path string, visitor Visitor) error {
	// Find the node
	node, err := m.FindNode(path)
	if err != nil {
		return err
	}
	
	// Apply the visitor
	return node.Accept(visitor)
}

// PrintFileSystem returns a string representation of the file system.
func (m *FileSystemManager) PrintFileSystem() string {
	return m.root.Print("")
}

// CalculateStatistics calculates and returns statistics about the file system.
func (m *FileSystemManager) CalculateStatistics() (*StatisticsVisitor, error) {
	stats := NewStatisticsVisitor()
	err := m.ApplyVisitor("/", stats)
	return stats, err
}

// Search searches for nodes matching criteria and returns results.
func (m *FileSystemManager) Search(name, content string, bySize bool, minSize, maxSize int64, maxResults int) ([]FileSystemNode, error) {
	visitor := NewSearchVisitor(name, content, bySize, minSize, maxSize, maxResults)
	err := m.ApplyVisitor("/", visitor)
	
	// If the error is "maximum results reached", it's not a real error
	if err != nil && err.Error() == "maximum results reached" {
		err = nil
	}
	
	return visitor.Results, err
}

// UpdatePermissions updates permissions on nodes matching criteria.
func (m *FileSystemManager) UpdatePermissions(path string, add, remove Permission, affectFiles, affectDirs, recursive bool, extension string) (int, error) {
	visitor := NewPermissionUpdaterVisitor(add, remove, affectFiles, affectDirs, recursive, extension)
	
	node, err := m.FindNode(path)
	if err != nil {
		return 0, err
	}
	
	// Apply visitor only to the directory itself if it's not recursive
	if !recursive {
		switch n := node.(type) {
		case *Directory:
			if affectDirs {
				currentPerms := n.GetPermissions()
				newPerms := (currentPerms | add) &^ remove
				
				if currentPerms != newPerms {
					n.SetPermissions(newPerms)
					visitor.Modified++
				}
			}
			
			// Process immediate children only
			for _, child := range n.Children() {
				switch c := child.(type) {
				case *File:
					if affectFiles {
						if extension == "" || strings.HasSuffix(c.Name(), extension) {
							currentPerms := c.GetPermissions()
							newPerms := (currentPerms | add) &^ remove
							
							if currentPerms != newPerms {
								c.SetPermissions(newPerms)
								visitor.Modified++
							}
						}
					}
				case *Directory:
					if affectDirs {
						currentPerms := c.GetPermissions()
						newPerms := (currentPerms | add) &^ remove
						
						if currentPerms != newPerms {
							c.SetPermissions(newPerms)
							visitor.Modified++
						}
					}
				}
			}
		case *File:
			if affectFiles {
				if extension == "" || strings.HasSuffix(n.Name(), extension) {
					currentPerms := n.GetPermissions()
					newPerms := (currentPerms | add) &^ remove
					
					if currentPerms != newPerms {
						n.SetPermissions(newPerms)
						visitor.Modified++
					}
				}
			}
		}
		
		return visitor.Modified, nil
	}
	
	// Apply visitor recursively
	err = node.Accept(visitor)
	
	// If the error is "not recursive", it's not a real error because
	// we're handling non-recursive mode ourselves
	if err != nil && err.Error() == "not recursive" {
		err = nil
	}
	
	return visitor.Modified, err
}
