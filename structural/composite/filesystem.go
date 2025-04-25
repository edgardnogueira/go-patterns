package composite

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// FileSystemNode is the Component interface in the Composite pattern.
// It defines operations common to both leaf nodes (files) and composite nodes (directories).
type FileSystemNode interface {
	// Name returns the name of the node
	Name() string
	
	// Path returns the full path of the node
	Path() string
	
	// Size returns the size of the node in bytes
	Size() int64
	
	// IsDirectory returns true if the node is a directory
	IsDirectory() bool
	
	// CreationTime returns when the node was created
	CreationTime() time.Time
	
	// Print displays the node information
	Print(prefix string) string
	
	// Accept allows a visitor to visit the node
	Accept(visitor Visitor) error
}

// Permission represents file system permissions
type Permission int

const (
	// Read permission
	Read Permission = 1 << iota
	// Write permission
	Write
	// Execute permission
	Execute
)

// Visitor defines an interface for visiting nodes in the file system.
// This allows operations to be performed on the structure without changing the classes.
type Visitor interface {
	// VisitFile is called when a file is visited
	VisitFile(file *File) error
	
	// VisitDirectory is called when a directory is visited
	VisitDirectory(directory *Directory) error
}

// baseNode contains common attributes for both file and directory nodes
type baseNode struct {
	name         string
	parent       *Directory
	permissions  Permission
	creationTime time.Time
}

// Name returns the name of the node
func (b *baseNode) Name() string {
	return b.name
}

// Path returns the full path of the node
func (b *baseNode) Path() string {
	if b.parent == nil {
		return b.name
	}
	
	return filepath.Join(b.parent.Path(), b.name)
}

// CreationTime returns when the node was created
func (b *baseNode) CreationTime() time.Time {
	return b.creationTime
}

// SetPermissions sets the permissions for the node
func (b *baseNode) SetPermissions(perm Permission) {
	b.permissions = perm
}

// GetPermissions returns the permissions for the node
func (b *baseNode) GetPermissions() Permission {
	return b.permissions
}

// HasPermission checks if the node has the specified permission
func (b *baseNode) HasPermission(perm Permission) bool {
	return b.permissions&perm != 0
}

// PermissionsString returns a string representation of the permissions
func (b *baseNode) PermissionsString() string {
	perms := []string{}
	
	if b.HasPermission(Read) {
		perms = append(perms, "read")
	}
	
	if b.HasPermission(Write) {
		perms = append(perms, "write")
	}
	
	if b.HasPermission(Execute) {
		perms = append(perms, "execute")
	}
	
	if len(perms) == 0 {
		return "none"
	}
	
	return strings.Join(perms, ", ")
}

// File is a Leaf in the Composite pattern.
// It represents a file in the file system with no children.
type File struct {
	baseNode
	content []byte
}

// NewFile creates a new file with the given name and content
func NewFile(name string, content []byte) *File {
	return &File{
		baseNode: baseNode{
			name:         name,
			permissions:  Read | Write,
			creationTime: time.Now(),
		},
		content: content,
	}
}

// Size returns the size of the file in bytes
func (f *File) Size() int64 {
	return int64(len(f.content))
}

// IsDirectory returns false for files
func (f *File) IsDirectory() bool {
	return false
}

// Content returns the content of the file
func (f *File) Content() []byte {
	return f.content
}

// SetContent updates the content of the file
func (f *File) SetContent(content []byte) {
	f.content = content
}

// Print returns a string representation of the file
func (f *File) Print(prefix string) string {
	return fmt.Sprintf("%s- %s (file, size: %d bytes, permissions: %s)",
		prefix, f.name, f.Size(), f.PermissionsString())
}

// Accept allows a visitor to visit the file
func (f *File) Accept(visitor Visitor) error {
	return visitor.VisitFile(f)
}

// Directory is a Composite in the Composite pattern.
// It represents a directory in the file system which can contain other files and directories.
type Directory struct {
	baseNode
	children []FileSystemNode
}

// NewDirectory creates a new directory with the given name
func NewDirectory(name string) *Directory {
	return &Directory{
		baseNode: baseNode{
			name:         name,
			permissions:  Read | Write | Execute,
			creationTime: time.Now(),
		},
		children: []FileSystemNode{},
	}
}

// Size returns the total size of all children in the directory
func (d *Directory) Size() int64 {
	var size int64
	for _, child := range d.children {
		size += child.Size()
	}
	return size
}

// IsDirectory returns true for directories
func (d *Directory) IsDirectory() bool {
	return true
}

// Add adds a child node to the directory
func (d *Directory) Add(node FileSystemNode) {
	// Set parent for the node if it's a baseNode
	switch n := node.(type) {
	case *File:
		n.parent = d
	case *Directory:
		n.parent = d
	}
	
	d.children = append(d.children, node)
}

// Remove removes a child node from the directory
func (d *Directory) Remove(node FileSystemNode) bool {
	for i, child := range d.children {
		if child == node {
			d.children = append(d.children[:i], d.children[i+1:]...)
			return true
		}
	}
	return false
}

// GetChild returns a child node by name
func (d *Directory) GetChild(name string) FileSystemNode {
	for _, child := range d.children {
		if child.Name() == name {
			return child
		}
	}
	return nil
}

// Children returns all child nodes
func (d *Directory) Children() []FileSystemNode {
	return d.children
}

// Print returns a string representation of the directory and its children
func (d *Directory) Print(prefix string) string {
	var result strings.Builder
	
	result.WriteString(fmt.Sprintf("%s+ %s (directory, size: %d bytes, permissions: %s)\n",
		prefix, d.name, d.Size(), d.PermissionsString()))
	
	childPrefix := prefix + "  "
	for _, child := range d.children {
		result.WriteString(child.Print(childPrefix) + "\n")
	}
	
	// Remove the last newline
	resultStr := result.String()
	if len(resultStr) > 0 {
		resultStr = resultStr[:len(resultStr)-1]
	}
	
	return resultStr
}

// Accept allows a visitor to visit the directory and its children
func (d *Directory) Accept(visitor Visitor) error {
	err := visitor.VisitDirectory(d)
	if err != nil {
		return err
	}
	
	for _, child := range d.children {
		err = child.Accept(visitor)
		if err != nil {
			return err
		}
	}
	
	return nil
}

// FindByPath finds a node by its path
func (d *Directory) FindByPath(path string) FileSystemNode {
	if path == "" || path == "." {
		return d
	}
	
	// Split the path into components
	components := strings.Split(filepath.Clean(path), string(filepath.Separator))
	
	// If the path is absolute and starts from root
	if components[0] == "" && len(components) > 1 {
		components = components[1:]
	}
	
	// If the first component is the current directory
	if components[0] == d.name {
		if len(components) == 1 {
			return d
		}
		components = components[1:]
	}
	
	// Find the first component
	child := d.GetChild(components[0])
	if child == nil {
		return nil
	}
	
	// If this is the last component, return the child
	if len(components) == 1 {
		return child
	}
	
	// If the child is a directory, continue searching in it
	if dir, ok := child.(*Directory); ok {
		return dir.FindByPath(filepath.Join(components[1:]...))
	}
	
	// If the child is a file but we have more path components, the path doesn't exist
	return nil
}
