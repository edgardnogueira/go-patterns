package composite

import (
	"strings"
	"testing"
	"time"
)

func TestFileOperations(t *testing.T) {
	// Create a file
	content := []byte("Hello, World!")
	file := NewFile("test.txt", content)
	
	if file.Name() != "test.txt" {
		t.Errorf("Expected name to be 'test.txt', got '%s'", file.Name())
	}
	
	if file.Size() != int64(len(content)) {
		t.Errorf("Expected size to be %d, got %d", len(content), file.Size())
	}
	
	if file.IsDirectory() {
		t.Error("Expected IsDirectory() to be false")
	}
	
	if string(file.Content()) != "Hello, World!" {
		t.Errorf("Expected content to be 'Hello, World!', got '%s'", string(file.Content()))
	}
	
	// Test permissions
	file.SetPermissions(Read | Write)
	if !file.HasPermission(Read) {
		t.Error("Expected file to have read permission")
	}
	if !file.HasPermission(Write) {
		t.Error("Expected file to have write permission")
	}
	if file.HasPermission(Execute) {
		t.Error("Expected file not to have execute permission")
	}
	
	// Test permission string
	permsStr := file.PermissionsString()
	if !strings.Contains(permsStr, "read") || !strings.Contains(permsStr, "write") {
		t.Errorf("Incorrect permissions string: %s", permsStr)
	}
	
	// Test creation time
	if time.Since(file.CreationTime()) > 1*time.Minute {
		t.Error("Creation time is too far in the past")
	}
	
	// Test print
	output := file.Print("")
	if !strings.Contains(output, "test.txt") || !strings.Contains(output, "13 bytes") {
		t.Errorf("Unexpected print output: %s", output)
	}
}

func TestDirectoryOperations(t *testing.T) {
	// Create a directory
	dir := NewDirectory("docs")
	
	if dir.Name() != "docs" {
		t.Errorf("Expected name to be 'docs', got '%s'", dir.Name())
	}
	
	if !dir.IsDirectory() {
		t.Error("Expected IsDirectory() to be true")
	}
	
	if dir.Size() != 0 {
		t.Errorf("Expected empty directory size to be 0, got %d", dir.Size())
	}
	
	// Add a file
	file := NewFile("hello.txt", []byte("Hello!"))
	dir.Add(file)
	
	// Check size update
	if dir.Size() != 6 {
		t.Errorf("Expected directory size to be 6, got %d", dir.Size())
	}
	
	// Test GetChild
	child := dir.GetChild("hello.txt")
	if child == nil {
		t.Error("Expected to find child 'hello.txt'")
	}
	
	// Test Children
	children := dir.Children()
	if len(children) != 1 {
		t.Errorf("Expected 1 child, got %d", len(children))
	}
	
	// Test Remove
	removed := dir.Remove(file)
	if !removed {
		t.Error("Expected file to be removed")
	}
	
	if len(dir.Children()) != 0 {
		t.Error("Expected directory to be empty after removal")
	}
	
	// Test print
	file2 := NewFile("test.txt", []byte("Test"))
	dir.Add(file2)
	
	output := dir.Print("")
	if !strings.Contains(output, "docs") || !strings.Contains(output, "test.txt") {
		t.Errorf("Unexpected print output: %s", output)
	}
}

func TestHierarchy(t *testing.T) {
	// Create a hierarchy
	root := NewDirectory("root")
	
	docsDir := NewDirectory("docs")
	root.Add(docsDir)
	
	readme := NewFile("readme.md", []byte("# Documentation"))
	docsDir.Add(readme)
	
	secretsDir := NewDirectory("secrets")
	docsDir.Add(secretsDir)
	
	secretFile := NewFile("key.txt", []byte("12345"))
	secretsDir.Add(secretFile)
	
	dataDir := NewDirectory("data")
	root.Add(dataDir)
	
	dataFile := NewFile("data.csv", []byte("id,name\n1,John"))
	dataDir.Add(dataFile)
	
	// Test sizes
	if root.Size() != 36 {
		t.Errorf("Expected root size to be 36, got %d", root.Size())
	}
	
	if docsDir.Size() != 21 {
		t.Errorf("Expected docs size to be 21, got %d", docsDir.Size())
	}
	
	// Test paths
	if readme.Path() != "root/docs/readme.md" {
		t.Errorf("Expected path to be 'root/docs/readme.md', got '%s'", readme.Path())
	}
	
	if secretFile.Path() != "root/docs/secrets/key.txt" {
		t.Errorf("Expected path to be 'root/docs/secrets/key.txt', got '%s'", secretFile.Path())
	}
	
	// Test FindByPath
	foundNode := root.FindByPath("docs/secrets/key.txt")
	if foundNode == nil {
		t.Error("Expected to find node by path")
	}
	
	if foundNode.Name() != "key.txt" {
		t.Errorf("Expected to find 'key.txt', got '%s'", foundNode.Name())
	}
	
	// Test non-existent path
	notFound := root.FindByPath("docs/not-exists")
	if notFound != nil {
		t.Error("Expected nil for non-existent path")
	}
}

func TestFileSystemManager(t *testing.T) {
	fsm := NewFileSystemManager()
	
	// Test root
	root := fsm.Root()
	if root.Name() != "/" {
		t.Errorf("Expected root name to be '/', got '%s'", root.Name())
	}
	
	// Test creating files
	file1, err := fsm.CreateFile("/test.txt", []byte("Test file"))
	if err != nil {
		t.Errorf("Error creating file: %v", err)
	}
	
	if file1.Name() != "test.txt" {
		t.Errorf("Expected name to be 'test.txt', got '%s'", file1.Name())
	}
	
	// Test creating files in subdirectories
	file2, err := fsm.CreateFile("/docs/readme.md", []byte("# README"))
	if err != nil {
		t.Errorf("Error creating file in subdirectory: %v", err)
	}
	
	if file2.Path() != "/docs/readme.md" {
		t.Errorf("Expected path to be '/docs/readme.md', got '%s'", file2.Path())
	}
	
	// Test creating directories
	dir1, err := fsm.CreateDirectory("/config")
	if err != nil {
		t.Errorf("Error creating directory: %v", err)
	}
	
	if dir1.Name() != "config" {
		t.Errorf("Expected name to be 'config', got '%s'", dir1.Name())
	}
	
	// Test creating nested directories
	dir2, err := fsm.CreateDirectoryPath("/data/processed/2023")
	if err != nil {
		t.Errorf("Error creating nested directories: %v", err)
	}
	
	if dir2.Path() != "/data/processed/2023" {
		t.Errorf("Expected path to be '/data/processed/2023', got '%s'", dir2.Path())
	}
	
	// Test finding nodes
	found, err := fsm.FindNode("/docs/readme.md")
	if err != nil {
		t.Errorf("Error finding node: %v", err)
	}
	
	if found.Name() != "readme.md" {
		t.Errorf("Expected to find 'readme.md', got '%s'", found.Name())
	}
	
	// Test deleting nodes
	err = fsm.DeleteNode("/docs/readme.md")
	if err != nil {
		t.Errorf("Error deleting node: %v", err)
	}
	
	_, err = fsm.FindNode("/docs/readme.md")
	if err == nil {
		t.Error("Expected error finding deleted node")
	}
	
	// Test moving nodes
	file3, err := fsm.CreateFile("/temp.txt", []byte("Temporary"))
	if err != nil {
		t.Errorf("Error creating file: %v", err)
	}
	
	err = fsm.MoveNode("/temp.txt", "/data/temp.txt")
	if err != nil {
		t.Errorf("Error moving node: %v", err)
	}
	
	_, err = fsm.FindNode("/temp.txt")
	if err == nil {
		t.Error("Expected error finding moved node at old path")
	}
	
	moved, err := fsm.FindNode("/data/temp.txt")
	if err != nil {
		t.Errorf("Error finding moved node: %v", err)
	}
	
	if moved.Name() != "temp.txt" {
		t.Errorf("Expected moved node name to be 'temp.txt', got '%s'", moved.Name())
	}
	
	// Test copying nodes
	err = fsm.CopyNode("/data/temp.txt", "/backup/temp.txt")
	if err != nil {
		t.Errorf("Error copying node: %v", err)
	}
	
	original, err := fsm.FindNode("/data/temp.txt")
	if err != nil {
		t.Errorf("Error finding original node: %v", err)
	}
	
	copied, err := fsm.FindNode("/backup/temp.txt")
	if err != nil {
		t.Errorf("Error finding copied node: %v", err)
	}
	
	if original == copied {
		t.Error("Expected copied node to be a different instance")
	}
	
	originalFile := original.(*File)
	copiedFile := copied.(*File)
	
	if string(originalFile.Content()) != string(copiedFile.Content()) {
		t.Error("Expected copied file to have the same content")
	}
}

func TestVisitors(t *testing.T) {
	fsm := NewFileSystemManager()
	
	// Create a test file system
	fsm.CreateFile("/docs/readme.md", []byte("# README"))
	fsm.CreateFile("/docs/guide.md", []byte("# User Guide\nThis is a guide."))
	fsm.CreateDirectory("/docs/examples")
	fsm.CreateFile("/docs/examples/example1.md", []byte("# Example 1"))
	fsm.CreateFile("/docs/examples/example2.md", []byte("# Example 2"))
	fsm.CreateFile("/data/numbers.csv", []byte("1,2,3\n4,5,6"))
	fsm.CreateFile("/data/names.csv", []byte("John,Jane,Bob"))
	fsm.CreateFile("/config.json", []byte("{\"debug\": true}"))
	fsm.CreateFile("/.hidden", []byte("secret"))
	
	// Test SizeCalculatorVisitor
	sizeVisitor := NewSizeCalculatorVisitor(".md", 0, 0, false)
	err := fsm.ApplyVisitor("/docs", sizeVisitor)
	if err != nil {
		t.Errorf("Error applying size visitor: %v", err)
	}
	
	// We have 4 .md files with total size 50 bytes (excluding hidden files)
	if sizeVisitor.TotalSize != 50 {
		t.Errorf("Expected .md files size to be 50, got %d", sizeVisitor.TotalSize)
	}
	
	// Test SearchVisitor
	searchVisitor := NewSearchVisitor("example", "", false, 0, 0, 0)
	err = fsm.ApplyVisitor("/", searchVisitor)
	if err != nil {
		t.Errorf("Error applying search visitor: %v", err)
	}
	
	// Should find directory /docs/examples and files example1.md and example2.md
	if len(searchVisitor.Results) != 3 {
		t.Errorf("Expected to find 3 nodes, got %d", len(searchVisitor.Results))
	}
	
	// Test content search
	contentSearchVisitor := NewSearchVisitor("", "guide", false, 0, 0, 0)
	err = fsm.ApplyVisitor("/", contentSearchVisitor)
	if err != nil {
		t.Errorf("Error applying content search visitor: %v", err)
	}
	
	// Should find only guide.md
	if len(contentSearchVisitor.Results) != 1 {
		t.Errorf("Expected to find 1 node, got %d", len(contentSearchVisitor.Results))
	}
	
	// Test PermissionUpdaterVisitor
	// First check current permissions
	file, err := fsm.FindNode("/docs/readme.md")
	if err != nil {
		t.Errorf("Error finding file: %v", err)
	}
	
	fileNode := file.(*File)
	if !fileNode.HasPermission(Write) {
		t.Error("Expected file to have write permission by default")
	}
	
	// Remove write permission from .md files
	modified, err := fsm.UpdatePermissions("/docs", 0, Write, true, false, true, ".md")
	if err != nil {
		t.Errorf("Error updating permissions: %v", err)
	}
	
	// Should modify 4 files
	if modified != 4 {
		t.Errorf("Expected to modify 4 files, got %d", modified)
	}
	
	// Check that permissions were updated
	file, err = fsm.FindNode("/docs/readme.md")
	if err != nil {
		t.Errorf("Error finding file: %v", err)
	}
	
	fileNode = file.(*File)
	if fileNode.HasPermission(Write) {
		t.Error("Expected file to have write permission removed")
	}
	
	// Test StatisticsVisitor
	stats, err := fsm.CalculateStatistics()
	if err != nil {
		t.Errorf("Error calculating statistics: %v", err)
	}
	
	if stats.FileCount != 9 {
		t.Errorf("Expected 9 files, got %d", stats.FileCount)
	}
	
	if stats.DirectoryCount != 4 {
		t.Errorf("Expected 4 directories (including root), got %d", stats.DirectoryCount)
	}
	
	// Print statistics
	report := stats.PrintReport()
	t.Logf("File System Statistics:\n%s", report)
}
