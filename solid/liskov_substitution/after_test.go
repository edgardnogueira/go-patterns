package main

import (
	"errors"
	"testing"
)

// MockReadableStorage implements ReadableStorage for testing
type MockReadableStorage struct {
	Data       map[string][]byte
	ExistsFunc func(filename string) (bool, error)
	ReadFunc   func(filename string) ([]byte, error)
}

func (m *MockReadableStorage) Read(filename string) ([]byte, error) {
	if m.ReadFunc != nil {
		return m.ReadFunc(filename)
	}
	data, ok := m.Data[filename]
	if !ok {
		return nil, errors.New("file not found")
	}
	return data, nil
}

func (m *MockReadableStorage) Exists(filename string) (bool, error) {
	if m.ExistsFunc != nil {
		return m.ExistsFunc(filename)
	}
	_, ok := m.Data[filename]
	return ok, nil
}

// MockWritableStorage implements WritableStorage for testing
type MockWritableStorage struct {
	SavedData   map[string][]byte
	DeletedKeys []string
	SaveFunc    func(filename string, data []byte) error
	DeleteFunc  func(filename string) error
}

func (m *MockWritableStorage) Save(filename string, data []byte) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(filename, data)
	}
	if m.SavedData == nil {
		m.SavedData = make(map[string][]byte)
	}
	m.SavedData[filename] = data
	return nil
}

func (m *MockWritableStorage) Delete(filename string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(filename)
	}
	if m.SavedData != nil {
		delete(m.SavedData, filename)
	}
	m.DeletedKeys = append(m.DeletedKeys, filename)
	return nil
}

// MockFullStorage implements both ReadableStorage and WritableStorage
type MockFullStorage struct {
	MockReadableStorage
	MockWritableStorage
}

func TestLSPWithFileReader(t *testing.T) {
	// Test with a read-only storage
	t.Run("file reader with read-only storage", func(t *testing.T) {
		// Create a mock readable storage
		mockStorage := &MockReadableStorage{
			Data: map[string][]byte{
				"test.txt": []byte("test content"),
			},
		}
		
		// Create a file reader with the mock storage
		reader := &FileReader{Storage: mockStorage}
		
		// Test GetFileContent
		content, err := reader.GetFileContent("test.txt")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if content != "test content" {
			t.Errorf("Expected 'test content', got '%s'", content)
		}
		
		// Test FileExists
		exists, err := reader.FileExists("test.txt")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !exists {
			t.Errorf("Expected file to exist, got false")
		}
		
		// Test FileExists for non-existent file
		exists, err = reader.FileExists("nonexistent.txt")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if exists {
			t.Errorf("Expected file to not exist, got true")
		}
	})
	
	// Test with a full storage
	t.Run("file reader with full storage", func(t *testing.T) {
		// Create a mock full storage
		mockStorage := &MockFullStorage{
			MockReadableStorage: MockReadableStorage{
				Data: map[string][]byte{
					"test.txt": []byte("test content"),
				},
			},
		}
		
		// Create a file reader with the mock storage
		reader := &FileReader{Storage: mockStorage}
		
		// The reader should work exactly the same way as with a read-only storage
		// This demonstrates LSP: the full storage can be used wherever a readable storage is expected
		
		content, err := reader.GetFileContent("test.txt")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if content != "test content" {
			t.Errorf("Expected 'test content', got '%s'", content)
		}
		
		exists, err := reader.FileExists("test.txt")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !exists {
			t.Errorf("Expected file to exist, got false")
		}
	})
	
	// Test error conditions
	t.Run("file reader with errors", func(t *testing.T) {
		// Create a mock readable storage that returns errors
		mockStorage := &MockReadableStorage{
			ReadFunc: func(filename string) ([]byte, error) {
				return nil, errors.New("read error")
			},
			ExistsFunc: func(filename string) (bool, error) {
				return false, errors.New("exists error")
			},
		}
		
		// Create a file reader with the mock storage
		reader := &FileReader{Storage: mockStorage}
		
		// Test GetFileContent with error
		_, err := reader.GetFileContent("test.txt")
		if err == nil {
			t.Error("Expected read error, got nil")
		}
		
		// Test FileExists with error
		_, err = reader.FileExists("test.txt")
		if err == nil {
			t.Error("Expected exists error, got nil")
		}
	})
}

func TestLSPWithFileWriter(t *testing.T) {
	// Test with a writable storage
	t.Run("file writer with writable storage", func(t *testing.T) {
		// Create a mock writable storage
		mockStorage := &MockWritableStorage{
			SavedData: make(map[string][]byte),
		}
		
		// Create a file writer with the mock storage
		writer := &FileWriter{Storage: mockStorage}
		
		// Test SaveFile
		err := writer.SaveFile("test.txt", "new content")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		// Verify the file was saved
		savedContent, exists := mockStorage.SavedData["test.txt"]
		if !exists {
			t.Error("Expected file to be saved, but it wasn't")
		}
		if string(savedContent) != "new content" {
			t.Errorf("Expected 'new content', got '%s'", string(savedContent))
		}
		
		// Test DeleteFile
		err = writer.DeleteFile("test.txt")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		// Verify the file was deleted
		if len(mockStorage.DeletedKeys) != 1 || mockStorage.DeletedKeys[0] != "test.txt" {
			t.Errorf("Expected 'test.txt' to be deleted, got %v", mockStorage.DeletedKeys)
		}
	})
	
	// Test with a full storage
	t.Run("file writer with full storage", func(t *testing.T) {
		// Create a mock full storage
		mockStorage := &MockFullStorage{
			MockWritableStorage: MockWritableStorage{
				SavedData: make(map[string][]byte),
			},
		}
		
		// Create a file writer with the mock storage
		writer := &FileWriter{Storage: mockStorage}
		
		// The writer should work exactly the same way as with a writable-only storage
		// This demonstrates LSP: the full storage can be used wherever a writable storage is expected
		
		err := writer.SaveFile("test.txt", "new content")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		savedContent, exists := mockStorage.SavedData["test.txt"]
		if !exists {
			t.Error("Expected file to be saved, but it wasn't")
		}
		if string(savedContent) != "new content" {
			t.Errorf("Expected 'new content', got '%s'", string(savedContent))
		}
	})
	
	// Test error conditions
	t.Run("file writer with errors", func(t *testing.T) {
		// Create a mock writable storage that returns errors
		mockStorage := &MockWritableStorage{
			SaveFunc: func(filename string, data []byte) error {
				return errors.New("save error")
			},
			DeleteFunc: func(filename string) error {
				return errors.New("delete error")
			},
		}
		
		// Create a file writer with the mock storage
		writer := &FileWriter{Storage: mockStorage}
		
		// Test SaveFile with error
		err := writer.SaveFile("test.txt", "new content")
		if err == nil {
			t.Error("Expected save error, got nil")
		}
		
		// Test DeleteFile with error
		err = writer.DeleteFile("test.txt")
		if err == nil {
			t.Error("Expected delete error, got nil")
		}
	})
}

func TestLSPDemonstration(t *testing.T) {
	t.Run("demonstrate LSP benefits", func(t *testing.T) {
		// This test demonstrates how LSP allows for proper substitution
		
		// Create a full storage implementation
		fullStorage := &MockFullStorage{
			MockReadableStorage: MockReadableStorage{
				Data: map[string][]byte{
					"existing.txt": []byte("existing content"),
				},
			},
			MockWritableStorage: MockWritableStorage{
				SavedData: make(map[string][]byte),
			},
		}
		
		// Use it as a ReadableStorage
		reader := &FileReader{Storage: fullStorage}
		
		// Use it as a WritableStorage
		writer := &FileWriter{Storage: fullStorage}
		
		// Both can work with the same storage instance without conflicts
		content, err := reader.GetFileContent("existing.txt")
		if err != nil {
			t.Errorf("Expected no error when reading, got %v", err)
		}
		if content != "existing content" {
			t.Errorf("Expected 'existing content', got '%s'", content)
		}
		
		err = writer.SaveFile("new.txt", "new content")
		if err != nil {
			t.Errorf("Expected no error when writing, got %v", err)
		}
		
		// After writing, we should be able to read the new file
		newContent, err := reader.GetFileContent("new.txt")
		if err != nil {
			t.Errorf("Expected no error when reading new file, got %v", err)
		}
		if newContent != "new content" {
			t.Errorf("Expected 'new content', got '%s'", newContent)
		}
		
		// This demonstrates how the LSP allows different clients to use
		// different interfaces of the same implementation seamlessly
	})
}
