package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

// FileStorage is a base interface for file operations
type FileStorage interface {
	Save(filename string, data []byte) error
	Read(filename string) ([]byte, error)
	Delete(filename string) error
}

// LocalFileStorage implements FileStorage for local disk
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

// Delete deletes a local file
func (fs *LocalFileStorage) Delete(filename string) error {
	fullPath := fs.BasePath + "/" + filename
	return os.Remove(fullPath)
}

// ReadOnlyFileStorage violates LSP because it can't fulfill the contract
// of the FileStorage interface - specifically, it can't save or delete files
type ReadOnlyFileStorage struct {
	BasePath string
}

// Save violates LSP by not actually saving and returning an error instead
func (fs *ReadOnlyFileStorage) Save(filename string, data []byte) error {
	// This implementation violates the Liskov Substitution Principle
	// because clients expect Save to save the file, but this throws an error
	return fmt.Errorf("cannot save: read-only file system")
}

// Read reads data from a file
func (fs *ReadOnlyFileStorage) Read(filename string) ([]byte, error) {
	fullPath := fs.BasePath + "/" + filename
	return ioutil.ReadFile(fullPath)
}

// Delete violates LSP by not actually deleting and returning an error instead
func (fs *ReadOnlyFileStorage) Delete(filename string) error {
	// This implementation violates the Liskov Substitution Principle
	// because clients expect Delete to delete the file, but this throws an error
	return fmt.Errorf("cannot delete: read-only file system")
}

// FileManager uses FileStorage to manage files
type FileManager struct {
	Storage FileStorage
}

// SaveFile saves content to a file
func (fm *FileManager) SaveFile(filename string, content string) error {
	return fm.Storage.Save(filename, []byte(content))
}

// GetFileContent gets content from a file
func (fm *FileManager) GetFileContent(filename string) (string, error) {
	data, err := fm.Storage.Read(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// DeleteFile deletes a file
func (fm *FileManager) DeleteFile(filename string) error {
	return fm.Storage.Delete(filename)
}

// This function demonstrates how the LSP violation causes problems
func demonstrateFileStorageBeforeLSP() {
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

	// Use LocalFileStorage
	fmt.Println("Using LocalFileStorage:")
	localStorage := &LocalFileStorage{BasePath: tempDir}
	fileManager := &FileManager{Storage: localStorage}

	// Save a file
	err = fileManager.SaveFile("test.txt", "Hello, World!")
	if err != nil {
		fmt.Println("  Error saving file:", err)
	} else {
		fmt.Println("  File saved successfully")
	}

	// Read a file
	content, err := fileManager.GetFileContent("sample.txt")
	if err != nil {
		fmt.Println("  Error reading file:", err)
	} else {
		fmt.Println("  Read file content:", content)
	}

	// Delete a file
	err = fileManager.DeleteFile("sample.txt")
	if err != nil {
		fmt.Println("  Error deleting file:", err)
	} else {
		fmt.Println("  File deleted successfully")
	}

	// Now use ReadOnlyFileStorage
	fmt.Println("\nUsing ReadOnlyFileStorage (violates LSP):")
	readOnlyStorage := &ReadOnlyFileStorage{BasePath: tempDir}
	fileManager = &FileManager{Storage: readOnlyStorage}

	// Try to save a file - will fail due to LSP violation
	err = fileManager.SaveFile("test2.txt", "Hello, Read-Only World!")
	if err != nil {
		fmt.Println("  Error saving file:", err)
	} else {
		fmt.Println("  File saved successfully")
	}

	// Read a file - should work
	content, err = fileManager.GetFileContent("test.txt")
	if err != nil {
		fmt.Println("  Error reading file:", err)
	} else {
		fmt.Println("  Read file content:", content)
	}

	// Try to delete a file - will fail due to LSP violation
	err = fileManager.DeleteFile("test.txt")
	if err != nil {
		fmt.Println("  Error deleting file:", err)
	} else {
		fmt.Println("  File deleted successfully")
	}

	fmt.Println("\nThe ReadOnlyFileStorage violates the Liskov Substitution Principle")
	fmt.Println("because it cannot be used in place of FileStorage without breaking")
	fmt.Println("the behavior expected by the FileManager class.")
}
