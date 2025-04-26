package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

// ReadableStorage defines a contract for read-only storage operations
// By separating read-only operations from write operations, we allow for better LSP adherence
type ReadableStorage interface {
	Read(filename string) ([]byte, error)
	Exists(filename string) (bool, error)
}

// WritableStorage defines a contract for write operations
type WritableStorage interface {
	Save(filename string, data []byte) error
	Delete(filename string) error
}

// FileStorage combines both readable and writable operations
// This forms a hierarchy of interfaces that allows for proper substitution
type FileStorage interface {
	ReadableStorage
	WritableStorage
}

// LocalFileStorage implements the full FileStorage interface
type LocalFileStorage struct {
	BasePath string
}

// Save saves data to a local file
func (fs *LocalFileStorage) Save(filename string, data []byte) error {
	fullPath := fs.BasePath + "/" + filename
	return ioutil.WriteFile(fullPath, data, 0644)
}

// Read reads data from a local file
func (fs *LocalFileStorage) Read(filename string) ([]byte, error) {
	fullPath := fs.BasePath + "/" + filename
	return ioutil.ReadFile(fullPath)
}

// Exists checks if a file exists
func (fs *LocalFileStorage) Exists(filename string) (bool, error) {
	fullPath := fs.BasePath + "/" + filename
	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err // Some other error occurred
}

// Delete deletes a local file
func (fs *LocalFileStorage) Delete(filename string) error {
	fullPath := fs.BasePath + "/" + filename
	return os.Remove(fullPath)
}

// ReadOnlyFileStorage implements only the ReadableStorage interface
// This follows LSP because it doesn't promise write capabilities
type ReadOnlyFileStorage struct {
	BasePath string
}

// Read reads data from a file
func (fs *ReadOnlyFileStorage) Read(filename string) ([]byte, error) {
	fullPath := fs.BasePath + "/" + filename
	return ioutil.ReadFile(fullPath)
}

// Exists checks if a file exists
func (fs *ReadOnlyFileStorage) Exists(filename string) (bool, error) {
	fullPath := fs.BasePath + "/" + filename
	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err // Some other error occurred
}

// FileReader manages read operations on files
// It expects only a ReadableStorage, so it works with both storage types
type FileReader struct {
	Storage ReadableStorage
}

// GetFileContent gets content from a file
func (fr *FileReader) GetFileContent(filename string) (string, error) {
	data, err := fr.Storage.Read(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FileExists checks if a file exists
func (fr *FileReader) FileExists(filename string) (bool, error) {
	return fr.Storage.Exists(filename)
}

// FileWriter manages write operations on files
// It requires a WritableStorage, so it only works with full implementations
type FileWriter struct {
	Storage WritableStorage
}

// SaveFile saves content to a file
func (fw *FileWriter) SaveFile(filename string, content string) error {
	return fw.Storage.Save(filename, []byte(content))
}

// DeleteFile deletes a file
func (fw *FileWriter) DeleteFile(filename string) error {
	return fw.Storage.Delete(filename)
}

// This function demonstrates how to apply LSP correctly
func demonstrateFileStorageAfterLSP() {
	// Create a temporary directory for testing
	tempDir, err := ioutil.TempDir("", "file-storage-test")
	if err != nil {
		fmt.Println("Error creating temp directory:", err)
		return
	}
	defer os.RemoveAll(tempDir)

	// Create sample file
	sampleFile := tempDir + "/sample.txt"
	ioutil.WriteFile(sampleFile, []byte("This is a test file"), 0644)

	// Use LocalFileStorage with both FileReader and FileWriter
	fmt.Println("Using LocalFileStorage (can read and write):")
	localStorage := &LocalFileStorage{BasePath: tempDir}
	
	fileReader := &FileReader{Storage: localStorage}
	fileWriter := &FileWriter{Storage: localStorage}

	// Use the writer to save a file
	err = fileWriter.SaveFile("test.txt", "Hello, World!")
	if err != nil {
		fmt.Println("  Error saving file:", err)
	} else {
		fmt.Println("  File saved successfully")
	}

	// Use the reader to read a file
	content, err := fileReader.GetFileContent("sample.txt")
	if err != nil {
		fmt.Println("  Error reading file:", err)
	} else {
		fmt.Println("  Read file content:", content)
	}

	// Use the writer to delete a file
	err = fileWriter.DeleteFile("sample.txt")
	if err != nil {
		fmt.Println("  Error deleting file:", err)
	} else {
		fmt.Println("  File deleted successfully")
	}

	// Now use ReadOnlyFileStorage with just FileReader
	fmt.Println("\nUsing ReadOnlyFileStorage (follows LSP by only implementing ReadableStorage):")
	readOnlyStorage := &ReadOnlyFileStorage{BasePath: tempDir}
	
	// We can only use a FileReader with ReadOnlyFileStorage
	fileReader = &FileReader{Storage: readOnlyStorage}

	// Read a file - should work fine
	content, err = fileReader.GetFileContent("test.txt")
	if err != nil {
		fmt.Println("  Error reading file:", err)
	} else {
		fmt.Println("  Read file content:", content)
	}

	// Check if a file exists
	exists, err := fileReader.FileExists("test.txt")
	if err != nil {
		fmt.Println("  Error checking if file exists:", err)
	} else if exists {
		fmt.Println("  File exists")
	} else {
		fmt.Println("  File does not exist")
	}

	// We can't use a FileWriter with ReadOnlyFileStorage because it doesn't implement WritableStorage
	// This is a compile-time check that prevents LSP violations
	// The following would cause a compile error:
	// fileWriter = &FileWriter{Storage: readOnlyStorage}

	fmt.Println("\nThis approach follows the Liskov Substitution Principle because:")
	fmt.Println("- We have separated the interfaces based on responsibilities")
	fmt.Println("- ReadableStorage can be substituted with any implementation that can read files")
	fmt.Println("- WritableStorage can be substituted with any implementation that can write files")
	fmt.Println("- Each implementation clearly states its capabilities through the interfaces it implements")
	fmt.Println("- Clients depend only on the interfaces they actually need")
}
