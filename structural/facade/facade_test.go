package facade

import (
	"os"
	"testing"
	"time"
)

// MockFileSystem is a mock implementation of file system operations for testing
type MockFileSystem struct {
	existingFiles map[string]bool
	directories   map[string]bool
}

// NewMockFileSystem creates a new mock file system
func NewMockFileSystem() *MockFileSystem {
	return &MockFileSystem{
		existingFiles: make(map[string]bool),
		directories:   make(map[string]bool),
	}
}

// FileExists mocks checking if a file exists
func (m *MockFileSystem) FileExists(path string) bool {
	return m.existingFiles[path]
}

// CreateDirectory mocks creating a directory
func (m *MockFileSystem) CreateDirectory(path string) error {
	m.directories[path] = true
	return nil
}

// AddExistingFile adds a file to the mock file system
func (m *MockFileSystem) AddExistingFile(path string) {
	m.existingFiles[path] = true
}

// MoveFile mocks moving a file
func (m *MockFileSystem) MoveFile(source, destination string) error {
	if !m.existingFiles[source] {
		return os.ErrNotExist
	}
	m.existingFiles[destination] = true
	return nil
}

// GetFileSize mocks getting file size
func (m *MockFileSystem) GetFileSize(path string) (int64, error) {
	if !m.existingFiles[path] {
		return 0, os.ErrNotExist
	}
	return 1024 * 1024, nil // 1MB
}

// GetFileExtension mocks getting file extension
func (m *MockFileSystem) GetFileExtension(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			return path[i:]
		}
	}
	return ""
}

// TestMediaConverterFacade tests the basic functionality of MediaConverterFacade
func TestMediaConverterFacade(t *testing.T) {
	// Create a new facade
	facade := NewMediaConverterFacade()
	
	// Test facade creation
	if facade == nil {
		t.Fatalf("Failed to create MediaConverterFacade")
	}
	
	// Check if subsystems were initialized
	if facade.videoProcessor == nil || 
	   facade.audioProcessor == nil || 
	   facade.codecManager == nil || 
	   facade.fileSystem == nil || 
	   facade.metadataHandler == nil || 
	   facade.progressReporter == nil {
		t.Fatalf("One or more subsystems not initialized")
	}
}

// TestConvertVideo tests the ConvertVideo method
func TestConvertVideo(t *testing.T) {
	// Create a facade with a mock file system
	facade := NewMediaConverterFacade()
	mockFS := NewMockFileSystem()
	facade.fileSystem = mockFS
	
	// Add some test files
	mockFS.AddExistingFile("test_video.mp4")
	
	// Test with non-existent file
	err := facade.ConvertVideo("non_existent_file.mp4", "output.mp4", FormatMP4)
	if err == nil {
		t.Errorf("ConvertVideo should fail with non-existent file")
	}
	
	// Test with existing file
	err = facade.ConvertVideo("test_video.mp4", "output.mp4", FormatMP4)
	if err != nil {
		t.Errorf("ConvertVideo failed: %v", err)
	}
	
	// Check if a job was created
	jobs := facade.GetActiveJobs()
	if len(jobs) == 0 {
		t.Errorf("No jobs created")
	}
	
	// Give the job some time to process
	time.Sleep(100 * time.Millisecond)
	
	// Check job status
	jobs = facade.GetActiveJobs()
	if len(jobs) == 0 {
		t.Errorf("No jobs found")
		return
	}
	
	job := jobs[0]
	if job.InputPath != "test_video.mp4" || job.OutputPath != "output.mp4" {
		t.Errorf("Job has incorrect paths: input=%s, output=%s", job.InputPath, job.OutputPath)
	}
}

// TestExtractAudio tests the ExtractAudio method
func TestExtractAudio(t *testing.T) {
	// Create a facade with a mock file system
	facade := NewMediaConverterFacade()
	mockFS := NewMockFileSystem()
	facade.fileSystem = mockFS
	
	// Add some test files
	mockFS.AddExistingFile("test_video.mp4")
	
	// Test with non-existent file
	err := facade.ExtractAudio("non_existent_file.mp4", "output.mp3", FormatMP3)
	if err == nil {
		t.Errorf("ExtractAudio should fail with non-existent file")
	}
	
	// Test with existing file
	err = facade.ExtractAudio("test_video.mp4", "output.mp3", FormatMP3)
	if err != nil {
		t.Errorf("ExtractAudio failed: %v", err)
	}
	
	// Allow time for processing
	time.Sleep(100 * time.Millisecond)
}

// TestOptimizeForWeb tests the OptimizeForWeb method
func TestOptimizeForWeb(t *testing.T) {
	// Create a facade with a mock file system
	facade := NewMediaConverterFacade()
	mockFS := NewMockFileSystem()
	facade.fileSystem = mockFS
	
	// Add some test files
	mockFS.AddExistingFile("test_video.mp4")
	mockFS.AddExistingFile("test_audio.mp3")
	
	// Test for video
	err := facade.OptimizeForWeb("test_video.mp4", "web_video.mp4")
	if err != nil {
		t.Errorf("OptimizeForWeb failed for video: %v", err)
	}
	
	// Test for audio
	err = facade.OptimizeForWeb("test_audio.mp3", "web_audio.mp3")
	if err != nil {
		t.Errorf("OptimizeForWeb failed for audio: %v", err)
	}
	
	// Allow time for processing
	time.Sleep(100 * time.Millisecond)
}

// TestCreateThumbnail tests the CreateThumbnail method
func TestCreateThumbnail(t *testing.T) {
	// Create a facade with a mock file system
	facade := NewMediaConverterFacade()
	mockFS := NewMockFileSystem()
	facade.fileSystem = mockFS
	
	// Add some test files
	mockFS.AddExistingFile("test_video.mp4")
	
	// Test with non-existent file
	err := facade.CreateThumbnail("non_existent_file.mp4", "thumbnail.jpg")
	if err == nil {
		t.Errorf("CreateThumbnail should fail with non-existent file")
	}
	
	// Test with existing file
	err = facade.CreateThumbnail("test_video.mp4", "thumbnail.jpg")
	if err != nil {
		t.Errorf("CreateThumbnail failed: %v", err)
	}
}

// TestBatchConvert tests the BatchConvert method
func TestBatchConvert(t *testing.T) {
	// Create a facade with a mock file system
	facade := NewMediaConverterFacade()
	mockFS := NewMockFileSystem()
	facade.fileSystem = mockFS
	
	// Add some test files
	mockFS.AddExistingFile("video1.mp4")
	mockFS.AddExistingFile("video2.avi")
	mockFS.AddExistingFile("audio1.mp3")
	mockFS.AddExistingFile("unknown.xyz")
	
	// Test batch conversion
	err := facade.BatchConvert(
		[]string{"video1.mp4", "video2.avi", "audio1.mp3", "unknown.xyz"}, 
		"output_dir", 
		FormatMP4,
	)
	
	if err != nil {
		t.Errorf("BatchConvert failed: %v", err)
	}
	
	// Check if output directory was created
	if !mockFS.directories["output_dir"] {
		t.Errorf("Output directory not created")
	}
	
	// Allow time for processing
	time.Sleep(100 * time.Millisecond)
}

// TestIsFormatSupported tests the IsFormatSupported method
func TestIsFormatSupported(t *testing.T) {
	facade := NewMediaConverterFacade()
	
	// Test valid video formats
	if !facade.IsFormatSupported(FormatMP4, true) {
		t.Errorf("MP4 should be supported for video")
	}
	
	// Test valid audio formats
	if !facade.IsFormatSupported(FormatMP3, false) {
		t.Errorf("MP3 should be supported for audio")
	}
	
	// Test invalid format combinations
	if facade.IsFormatSupported(FormatMP3, true) {
		t.Errorf("MP3 should not be supported for video")
	}
	
	if facade.IsFormatSupported(FormatMP4, false) {
		t.Errorf("MP4 should not be supported for audio")
	}
}

// TestGetSupportedFormats tests the GetSupportedFormats method
func TestGetSupportedFormats(t *testing.T) {
	facade := NewMediaConverterFacade()
	
	// Test video formats
	videoFormats := facade.GetSupportedFormats(true)
	if len(videoFormats) == 0 {
		t.Errorf("No supported video formats returned")
	}
	
	// Test audio formats
	audioFormats := facade.GetSupportedFormats(false)
	if len(audioFormats) == 0 {
		t.Errorf("No supported audio formats returned")
	}
	
	// Check for specific formats
	hasMP4 := false
	for _, format := range videoFormats {
		if format == FormatMP4 {
			hasMP4 = true
			break
		}
	}
	
	if !hasMP4 {
		t.Errorf("MP4 format not found in supported video formats")
	}
	
	hasMP3 := false
	for _, format := range audioFormats {
		if format == FormatMP3 {
			hasMP3 = true
			break
		}
	}
	
	if !hasMP3 {
		t.Errorf("MP3 format not found in supported audio formats")
	}
}

// TestCreateConversionProfile tests the CreateConversionProfile method
func TestCreateConversionProfile(t *testing.T) {
	facade := NewMediaConverterFacade()
	
	// Create a profile
	profile := facade.CreateConversionProfile(
		"HD Video",
		FormatMP4,
		Resolution1080p,
		BitrateHigh,
	)
	
	// Check profile values
	if profile.Format != FormatMP4 {
		t.Errorf("Expected format MP4, got %s", profile.Format)
	}
	
	if profile.Resolution != Resolution1080p {
		t.Errorf("Expected resolution 1080p, got %s", profile.Resolution)
	}
	
	if profile.Bitrate != BitrateHigh {
		t.Errorf("Expected bitrate high, got %s", profile.Bitrate)
	}
	
	// Check if codecs were set correctly
	if profile.VideoCodec == "" {
		t.Errorf("Video codec not set")
	}
	
	if profile.AudioCodec == "" {
		t.Errorf("Audio codec not set")
	}
}
