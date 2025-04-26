// Package interfaces demonstrates idiomatic Go interface implementation patterns
package interfaces

import (
	"fmt"
	"io"
	"time"
)

// Interface Segregation Principle
// ------------------------------
// "Clients should not be forced to depend on methods they do not use."
// In Go, this principle is fundamental - we use small, focused interfaces
// rather than large, monolithic ones.

// Monolithic interface (anti-pattern)
// ---------------------------------

// FileManager is a monolithic interface that handles many different operations
// This violates interface segregation and is shown as an anti-pattern
type FileManager interface {
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	Seek(offset int64, whence int) (int64, error)
	Close() error
	Stat() (FileInfo, error)
	Chmod(mode uint32) error
	Chown(uid, gid int) error
	Truncate(size int64) error
	Sync() error
}

// FileInfo is a simple interface representing file information
type FileInfo interface {
	Name() string
	Size() int64
	IsDir() bool
	ModTime() time.Time
}

// Segregated interfaces (good practice)
// -----------------------------------

// Reader is an interface for reading data
type FileReader interface {
	Read(p []byte) (n int, err error)
}

// Writer is an interface for writing data
type FileWriter interface {
	Write(p []byte) (n int, err error)
}

// Seeker is an interface for seeking within a file
type FileSeeker interface {
	Seek(offset int64, whence int) (int64, error)
}

// Closer is an interface for closing resources
type FileCloser interface {
	Close() error
}

// StatProvider provides file information
type StatProvider interface {
	Stat() (FileInfo, error)
}

// FilePermissionsManager handles file permissions
type FilePermissionsManager interface {
	Chmod(mode uint32) error
	Chown(uid, gid int) error
}

// FileTruncater handles file truncation
type FileTruncater interface {
	Truncate(size int64) error
}

// FileSyncer handles syncing file to disk
type FileSyncer interface {
	Sync() error
}

// We can compose these interfaces as needed
// This allows for precise capability requirements

// ReadWriter combines read and write capabilities
type FileReadWriter interface {
	FileReader
	FileWriter
}

// ReadWriteCloser adds closing capability
type FileReadWriteCloser interface {
	FileReader
	FileWriter
	FileCloser
}

// ReadSeeker combines read and seek capabilities
type FileReadSeeker interface {
	FileReader
	FileSeeker
}

// ReadWriteSeeker combines read, write and seek capabilities
type FileReadWriteSeeker interface {
	FileReader
	FileWriter
	FileSeeker
}

// Implementation example
// --------------------

// MockFile implements multiple interfaces through composition
type MockFile struct {
	data     []byte
	position int64
	name     string
	mode     uint32
	closed   bool
}

// NewMockFile creates a new MockFile
func NewMockFile(name string, data []byte) *MockFile {
	return &MockFile{
		data:     data,
		position: 0,
		name:     name,
		mode:     0644,
	}
}

// Read implements FileReader
func (f *MockFile) Read(p []byte) (int, error) {
	if f.closed {
		return 0, fmt.Errorf("file closed")
	}
	
	if f.position >= int64(len(f.data)) {
		return 0, io.EOF
	}
	
	n := copy(p, f.data[f.position:])
	f.position += int64(n)
	return n, nil
}

// Write implements FileWriter
func (f *MockFile) Write(p []byte) (int, error) {
	if f.closed {
		return 0, fmt.Errorf("file closed")
	}
	
	pos := int(f.position)
	if pos >= len(f.data) {
		// Extend the slice if needed
		f.data = append(f.data, make([]byte, pos-len(f.data))...)
		f.data = append(f.data, p...)
	} else {
		// Overwrite existing data
		copy(f.data[pos:], p)
		if pos+len(p) > len(f.data) {
			f.data = f.data[:pos+len(p)]
		}
	}
	
	f.position += int64(len(p))
	return len(p), nil
}

// Seek implements FileSeeker
func (f *MockFile) Seek(offset int64, whence int) (int64, error) {
	if f.closed {
		return 0, fmt.Errorf("file closed")
	}
	
	var newPos int64
	switch whence {
	case io.SeekStart:
		newPos = offset
	case io.SeekCurrent:
		newPos = f.position + offset
	case io.SeekEnd:
		newPos = int64(len(f.data)) + offset
	default:
		return 0, fmt.Errorf("invalid whence: %d", whence)
	}
	
	if newPos < 0 {
		return 0, fmt.Errorf("negative position")
	}
	
	f.position = newPos
	return newPos, nil
}

// Close implements FileCloser
func (f *MockFile) Close() error {
	f.closed = true
	return nil
}

// Stat implements StatProvider
func (f *MockFile) Stat() (FileInfo, error) {
	return &mockFileInfo{
		name:    f.name,
		size:    int64(len(f.data)),
		isDir:   false,
		modTime: time.Now(),
	}, nil
}

// Chmod implements part of FilePermissionsManager
func (f *MockFile) Chmod(mode uint32) error {
	if f.closed {
		return fmt.Errorf("file closed")
	}
	f.mode = mode
	return nil
}

// Chown implements part of FilePermissionsManager
func (f *MockFile) Chown(uid, gid int) error {
	if f.closed {
		return fmt.Errorf("file closed")
	}
	// In this mock implementation, we just acknowledge the request
	return nil
}

// Truncate implements FileTruncater
func (f *MockFile) Truncate(size int64) error {
	if f.closed {
		return fmt.Errorf("file closed")
	}
	
	if size < 0 {
		return fmt.Errorf("negative size")
	}
	
	if size < int64(len(f.data)) {
		f.data = f.data[:size]
	} else if size > int64(len(f.data)) {
		f.data = append(f.data, make([]byte, size-int64(len(f.data)))...)
	}
	
	return nil
}

// Sync implements FileSyncer
func (f *MockFile) Sync() error {
	if f.closed {
		return fmt.Errorf("file closed")
	}
	// In this mock implementation, we just acknowledge the request
	return nil
}

// mockFileInfo implements FileInfo
type mockFileInfo struct {
	name    string
	size    int64
	isDir   bool
	modTime time.Time
}

func (info *mockFileInfo) Name() string {
	return info.name
}

func (info *mockFileInfo) Size() int64 {
	return info.size
}

func (info *mockFileInfo) IsDir() bool {
	return info.isDir
}

func (info *mockFileInfo) ModTime() time.Time {
	return info.modTime
}

// Usage examples
// ------------

// ReadOnly accepts only a Reader, demonstrating that we can constrain
// functions to only require the minimal interface they need
func ReadOnly(r FileReader) ([]byte, error) {
	buf := make([]byte, 1024)
	n, err := r.Read(buf)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return buf[:n], nil
}

// WriteOnly accepts only a Writer
func WriteOnly(w FileWriter, data []byte) error {
	_, err := w.Write(data)
	return err
}

// CopyData accepts minimal interfaces needed for copying
func CopyData(r FileReader, w FileWriter) (int, error) {
	buf := make([]byte, 1024)
	total := 0
	
	for {
		n, err := r.Read(buf)
		if n > 0 {
			// Write what we read
			if _, werr := w.Write(buf[:n]); werr != nil {
				return total, werr
			}
			total += n
		}
		
		if err == io.EOF {
			break
		}
		if err != nil {
			return total, err
		}
	}
	
	return total, nil
}

// AdvancedOperation requires more capabilities
func AdvancedOperation(rws FileReadWriteSeeker) error {
	// Seek to beginning
	_, err := rws.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	
	// Read some data
	buf := make([]byte, 10)
	n, err := rws.Read(buf)
	if err != nil && err != io.EOF {
		return err
	}
	
	fmt.Printf("Read %d bytes: %s\n", n, buf[:n])
	
	// Seek to end and write
	_, err = rws.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}
	
	_, err = rws.Write([]byte(" - Appended"))
	return err
}

// SegregationDemo demonstrates the interface segregation principle
func SegregationDemo() {
	fmt.Println("============================================")
	fmt.Println("Interface Segregation Principle Demo")
	fmt.Println("============================================")
	
	// Create a file that implements multiple interface capabilities
	file := NewMockFile("example.txt", []byte("Hello, Interface Segregation!"))
	
	// Use file with minimal interface requirements
	data, _ := ReadOnly(file)
	fmt.Printf("Read operation: %s\n", data)
	
	WriteOnly(file, []byte(" More text."))
	
	// Reset position
	file.Seek(0, io.SeekStart)
	
	// Read again to see combined content
	data, _ = ReadOnly(file)
	fmt.Printf("After write: %s\n", data)
	
	// Create a second file for copy operation
	destFile := NewMockFile("destination.txt", []byte{})
	
	// Reset source file position
	file.Seek(0, io.SeekStart)
	
	// Copy between files using minimal interfaces
	n, _ := CopyData(file, destFile)
	fmt.Printf("Copied %d bytes between files\n", n)
	
	// Read from destination
	destFile.Seek(0, io.SeekStart)
	destData, _ := ReadOnly(destFile)
	fmt.Printf("Destination content: %s\n", destData)
	
	// Use advanced operation requiring more capabilities
	file.Seek(0, io.SeekStart)
	AdvancedOperation(file)
	
	// Check final state
	file.Seek(0, io.SeekStart)
	finalData, _ := ReadOnly(file)
	fmt.Printf("Final content: %s\n", finalData)
	
	fmt.Println("\nBenefits of Interface Segregation:")
	fmt.Println("1. Functions only depend on methods they actually use")
	fmt.Println("2. Easier testing with minimal mock implementations")
	fmt.Println("3. More flexibility in implementation changes")
	fmt.Println("4. Better separation of concerns")
	fmt.Println("============================================")
}
