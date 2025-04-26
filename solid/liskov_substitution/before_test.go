package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestFileStorageBeforeLSP(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := ioutil.TempDir("", "file-storage-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a sample file for testing
	sampleFilePath := filepath.Join(tempDir, "sample.txt")
	sampleContent := []byte("Sample file content")
	err = ioutil.WriteFile(sampleFilePath, sampleContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create sample file: %v", err)
	}

	// Test the LocalFileStorage implementation
	t.Run("local file storage", func(t *testing.T) {
		localStorage := &LocalFileStorage{BasePath: tempDir}

		// Test Read
		data, err := localStorage.Read("sample.txt")
		if err != nil {
			t.Errorf("Expected no error when reading, got %v", err)
		}
		if string(data) != "Sample file content" {
			t.Errorf("Expected 'Sample file content', got '%s'", string(data))
		}

		// Test Save
		err = localStorage.Save("new.txt", []byte("New file content"))
		if err != nil {
			t.Errorf("Expected no error when saving, got %v", err)
		}

		// Verify the file was saved
		newFilePath := filepath.Join(tempDir, "new.txt")
		if _, err := os.Stat(newFilePath); os.IsNotExist(err) {
			t.Error("Expected new file to exist, but it doesn't")
		}

		// Test Delete
		err = localStorage.Delete("new.txt")
		if err != nil {
			t.Errorf("Expected no error when deleting, got %v", err)
		}

		// Verify the file was deleted
		if _, err := os.Stat(newFilePath); !os.IsNotExist(err) {
			t.Error("Expected new file to be deleted, but it still exists")
		}
	})

	// Test the ReadOnlyFileStorage implementation - This demonstrates LSP violations
	t.Run("read only file storage - LSP violation", func(t *testing.T) {
		readOnlyStorage := &ReadOnlyFileStorage{BasePath: tempDir}

		// Test Read - Should work
		data, err := readOnlyStorage.Read("sample.txt")
		if err != nil {
			t.Errorf("Expected no error when reading, got %v", err)
		}
		if string(data) != "Sample file content" {
			t.Errorf("Expected 'Sample file content', got '%s'", string(data))
		}

		// Test Save - Should FAIL, violating LSP
		err = readOnlyStorage.Save("new.txt", []byte("New file content"))
		if err == nil {
			t.Error("Expected error when trying to save to read-only storage, got nil")
		} else {
			t.Logf("Got expected error for Save operation: %v", err)
		}

		// Test Delete - Should FAIL, violating LSP
		err = readOnlyStorage.Delete("sample.txt")
		if err == nil {
			t.Error("Expected error when trying to delete from read-only storage, got nil")
		} else {
			t.Logf("Got expected error for Delete operation: %v", err)
		}
	})

	// Test the FileManager with different storage implementations
	t.Run("file manager with different storage implementations", func(t *testing.T) {
		// Test with local storage
		localStorage := &LocalFileStorage{BasePath: tempDir}
		fileManager := &FileManager{Storage: localStorage}

		// Save file - Should work
		err := fileManager.SaveFile("test.txt", "Test content")
		if err != nil {
			t.Errorf("Expected no error when saving with local storage, got %v", err)
		}

		// Read file - Should work
		content, err := fileManager.GetFileContent("test.txt")
		if err != nil {
			t.Errorf("Expected no error when reading with local storage, got %v", err)
		}
		if content != "Test content" {
			t.Errorf("Expected 'Test content', got '%s'", content)
		}

		// Delete file - Should work
		err = fileManager.DeleteFile("test.txt")
		if err != nil {
			t.Errorf("Expected no error when deleting with local storage, got %v", err)
		}

		// Now test with read-only storage
		readOnlyStorage := &ReadOnlyFileStorage{BasePath: tempDir}
		fileManager = &FileManager{Storage: readOnlyStorage}

		// Save file - Should FAIL, breaking client expectations
		err = fileManager.SaveFile("readonly.txt", "This won't work")
		if err == nil {
			t.Error("Expected error when saving with read-only storage, got nil")
		} else {
			t.Logf("Got expected error: %v", err)
		}

		// This demonstrates the LSP violation: FileManager expects to work with any
		// FileStorage implementation, but it breaks when using ReadOnlyFileStorage
    })
}

func TestLSPViolationIssues(t *testing.T) {
	t.Run("highlight issues with LSP violations", func(t *testing.T) {
		// This "test" highlights the issues with LSP violations
		
		// ISSUE 1: Client code breaks when using subtypes
		t.Log("ISSUE 1: Client code breaks when using subtypes")
		t.Log("- FileManager expects to use any FileStorage implementation")
		t.Log("- But ReadOnlyFileStorage fails for Save and Delete operations")
		t.Log("- This violates client expectations and breaks the contract")
		
		// ISSUE 2: Requires checking implementation type at runtime
		t.Log("ISSUE 2: Requires checking implementation type or handling errors")
		t.Log("- Client code must either check implementation type before operations")
		t.Log("- Or handle unexpected errors that shouldn't occur with proper substitution")
		t.Log("- This leads to defensive programming and complex error handling")
		
		// ISSUE 3: Makes testing more difficult
		t.Log("ISSUE 3: Makes testing more difficult")
		t.Log("- Can't reliably mock the interface because implementations don't respect the contract")
		t.Log("- Tests must account for method failures that shouldn't happen")
		
		// ISSUE 4: Violates interface segregation as well
		t.Log("ISSUE 4: Often also violates Interface Segregation Principle")
		t.Log("- Fat interfaces lead to implementations that can't fulfill all methods")
		t.Log("- A proper solution would separate interfaces by capability")
		
		// This is not a real test, just documentation
		t.Skip("This is not a real test, just documentation of LSP violation issues")
	})
}
