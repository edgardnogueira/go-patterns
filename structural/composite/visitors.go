package composite

import (
	"fmt"
	"strings"
	"time"
)

// SizeCalculatorVisitor calculates the total size of files
// matching specific criteria.
type SizeCalculatorVisitor struct {
	TotalSize     int64
	Extension     string
	MinSize       int64
	MaxSize       int64
	IncludeHidden bool
}

// NewSizeCalculatorVisitor creates a new SizeCalculatorVisitor
func NewSizeCalculatorVisitor(extension string, minSize, maxSize int64, includeHidden bool) *SizeCalculatorVisitor {
	return &SizeCalculatorVisitor{
		Extension:     extension,
		MinSize:       minSize,
		MaxSize:       maxSize,
		IncludeHidden: includeHidden,
	}
}

// VisitFile processes a file node
func (v *SizeCalculatorVisitor) VisitFile(file *File) error {
	// Skip hidden files if not included
	if !v.IncludeHidden && strings.HasPrefix(file.Name(), ".") {
		return nil
	}
	
	// Check extension if specified
	if v.Extension != "" && !strings.HasSuffix(file.Name(), v.Extension) {
		return nil
	}
	
	// Check file size constraints
	fileSize := file.Size()
	if (v.MinSize > 0 && fileSize < v.MinSize) || (v.MaxSize > 0 && fileSize > v.MaxSize) {
		return nil
	}
	
	// Add to total size if all criteria match
	v.TotalSize += fileSize
	return nil
}

// VisitDirectory processes a directory node
func (v *SizeCalculatorVisitor) VisitDirectory(directory *Directory) error {
	// Skip hidden directories if not included
	if !v.IncludeHidden && strings.HasPrefix(directory.Name(), ".") {
		return nil
	}
	
	// No action needed - we'll visit the contained files individually
	return nil
}

// SearchVisitor searches for nodes matching specific criteria
// and stores the results.
type SearchVisitor struct {
	Results      []FileSystemNode
	SearchName   string
	ContentMatch string
	BySize       bool
	MinSize      int64
	MaxSize      int64
	MaxResults   int
}

// NewSearchVisitor creates a new SearchVisitor
func NewSearchVisitor(searchName, contentMatch string, bySize bool, minSize, maxSize int64, maxResults int) *SearchVisitor {
	return &SearchVisitor{
		Results:      make([]FileSystemNode, 0),
		SearchName:   searchName,
		ContentMatch: contentMatch,
		BySize:       bySize,
		MinSize:      minSize,
		MaxSize:      maxSize,
		MaxResults:   maxResults,
	}
}

// VisitFile processes a file node
func (v *SearchVisitor) VisitFile(file *File) error {
	// Stop if we've reached the maximum number of results
	if v.MaxResults > 0 && len(v.Results) >= v.MaxResults {
		return fmt.Errorf("maximum results reached")
	}
	
	// Check name match
	if v.SearchName != "" && !strings.Contains(file.Name(), v.SearchName) {
		return nil
	}
	
	// Check size constraints if needed
	if v.BySize {
		fileSize := file.Size()
		if (v.MinSize > 0 && fileSize < v.MinSize) || (v.MaxSize > 0 && fileSize > v.MaxSize) {
			return nil
		}
	}
	
	// Check content match if needed
	if v.ContentMatch != "" {
		content := string(file.Content())
		if !strings.Contains(content, v.ContentMatch) {
			return nil
		}
	}
	
	// Add to results if all criteria match
	v.Results = append(v.Results, file)
	return nil
}

// VisitDirectory processes a directory node
func (v *SearchVisitor) VisitDirectory(directory *Directory) error {
	// Stop if we've reached the maximum number of results
	if v.MaxResults > 0 && len(v.Results) >= v.MaxResults {
		return fmt.Errorf("maximum results reached")
	}
	
	// Check name match
	if v.SearchName != "" && !strings.Contains(directory.Name(), v.SearchName) {
		return nil
	}
	
	// Check size constraints if needed
	if v.BySize {
		dirSize := directory.Size()
		if (v.MinSize > 0 && dirSize < v.MinSize) || (v.MaxSize > 0 && dirSize > v.MaxSize) {
			return nil
		}
	}
	
	// Add to results if criteria match
	v.Results = append(v.Results, directory)
	return nil
}

// PermissionUpdaterVisitor updates permissions on nodes matching specific criteria.
type PermissionUpdaterVisitor struct {
	PermissionsToAdd    Permission
	PermissionsToRemove Permission
	AffectFiles         bool
	AffectDirectories   bool
	Recursive           bool
	Extension           string
	Modified            int
}

// NewPermissionUpdaterVisitor creates a new PermissionUpdaterVisitor
func NewPermissionUpdaterVisitor(add, remove Permission, affectFiles, affectDirs, recursive bool, extension string) *PermissionUpdaterVisitor {
	return &PermissionUpdaterVisitor{
		PermissionsToAdd:    add,
		PermissionsToRemove: remove,
		AffectFiles:         affectFiles,
		AffectDirectories:   affectDirs,
		Recursive:           recursive,
		Extension:           extension,
	}
}

// VisitFile processes a file node
func (v *PermissionUpdaterVisitor) VisitFile(file *File) error {
	// Only proceed if we're affecting files
	if !v.AffectFiles {
		return nil
	}
	
	// Check extension if specified
	if v.Extension != "" && !strings.HasSuffix(file.Name(), v.Extension) {
		return nil
	}
	
	// Update permissions
	currentPerms := file.GetPermissions()
	newPerms := (currentPerms | v.PermissionsToAdd) &^ v.PermissionsToRemove
	
	// Only count as modified if permissions actually changed
	if currentPerms != newPerms {
		file.SetPermissions(newPerms)
		v.Modified++
	}
	
	return nil
}

// VisitDirectory processes a directory node
func (v *PermissionUpdaterVisitor) VisitDirectory(directory *Directory) error {
	// Update permissions on this directory if applicable
	if v.AffectDirectories {
		currentPerms := directory.GetPermissions()
		newPerms := (currentPerms | v.PermissionsToAdd) &^ v.PermissionsToRemove
		
		// Only count as modified if permissions actually changed
		if currentPerms != newPerms {
			directory.SetPermissions(newPerms)
			v.Modified++
		}
	}
	
	// If not recursive, skip the children by returning an error to stop traversal
	if !v.Recursive {
		return fmt.Errorf("not recursive")
	}
	
	return nil
}

// StatisticsVisitor collects statistics about the file system structure.
type StatisticsVisitor struct {
	FileCount      int
	DirectoryCount int
	TotalSize      int64
	LargestFile    *File
	LargestDir     *Directory
	OldestNode     FileSystemNode
	NewestNode     FileSystemNode
	MaxDepth       int
	CurrentDepth   int
}

// NewStatisticsVisitor creates a new StatisticsVisitor
func NewStatisticsVisitor() *StatisticsVisitor {
	return &StatisticsVisitor{
		CurrentDepth: 0,
	}
}

// VisitFile processes a file node
func (v *StatisticsVisitor) VisitFile(file *File) error {
	// Increment file count
	v.FileCount++
	
	// Add to total size
	fileSize := file.Size()
	v.TotalSize += fileSize
	
	// Check if this is the largest file
	if v.LargestFile == nil || fileSize > v.LargestFile.Size() {
		v.LargestFile = file
	}
	
	// Check if this is the oldest or newest node
	v.updateTimeStats(file)
	
	return nil
}

// VisitDirectory processes a directory node
func (v *StatisticsVisitor) VisitDirectory(directory *Directory) error {
	// Increment directory count
	v.DirectoryCount++
	
	// Check if this is the largest directory
	dirSize := directory.Size()
	if v.LargestDir == nil || dirSize > v.LargestDir.Size() {
		v.LargestDir = directory
	}
	
	// Check if this is the oldest or newest node
	v.updateTimeStats(directory)
	
	// Update depth tracking
	v.CurrentDepth++
	if v.CurrentDepth > v.MaxDepth {
		v.MaxDepth = v.CurrentDepth
	}
	
	// Defer decreasing the depth counter when we finish with this directory
	defer func() {
		v.CurrentDepth--
	}()
	
	return nil
}

// updateTimeStats updates the statistics about oldest and newest nodes
func (v *StatisticsVisitor) updateTimeStats(node FileSystemNode) {
	if v.OldestNode == nil || node.CreationTime().Before(v.OldestNode.CreationTime()) {
		v.OldestNode = node
	}
	
	if v.NewestNode == nil || node.CreationTime().After(v.NewestNode.CreationTime()) {
		v.NewestNode = node
	}
}

// PrintReport returns a formatted report of the statistics
func (v *StatisticsVisitor) PrintReport() string {
	var result strings.Builder
	
	result.WriteString("File System Statistics:\n")
	result.WriteString(fmt.Sprintf("  Total Nodes: %d (%d files, %d directories)\n", 
		v.FileCount+v.DirectoryCount, v.FileCount, v.DirectoryCount))
	result.WriteString(fmt.Sprintf("  Total Size: %d bytes\n", v.TotalSize))
	result.WriteString(fmt.Sprintf("  Maximum Depth: %d levels\n", v.MaxDepth))
	
	if v.LargestFile != nil {
		result.WriteString(fmt.Sprintf("  Largest File: %s (%d bytes)\n", 
			v.LargestFile.Path(), v.LargestFile.Size()))
	}
	
	if v.LargestDir != nil {
		result.WriteString(fmt.Sprintf("  Largest Directory: %s (%d bytes)\n", 
			v.LargestDir.Path(), v.LargestDir.Size()))
	}
	
	if v.OldestNode != nil {
		result.WriteString(fmt.Sprintf("  Oldest Node: %s (created %s)\n", 
			v.OldestNode.Path(), v.OldestNode.CreationTime().Format(time.RFC3339)))
	}
	
	if v.NewestNode != nil {
		result.WriteString(fmt.Sprintf("  Newest Node: %s (created %s)", 
			v.NewestNode.Path(), v.NewestNode.CreationTime().Format(time.RFC3339)))
	}
	
	return result.String()
}
