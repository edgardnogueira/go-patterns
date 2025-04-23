package adapter

import (
	"strings"
	"testing"
)

func TestAudioPlayerWithMP3(t *testing.T) {
	player := &AudioPlayer{}
	result := player.Play("mp3", "song.mp3")

	expected := "Playing MP3 file: song.mp3"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestAudioPlayerWithFLAC(t *testing.T) {
	player := &AudioPlayer{}
	result := player.Play("flac", "song.flac")

	expected := "Playing FLAC file: song.flac"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestAudioPlayerWithWAV(t *testing.T) {
	player := &AudioPlayer{}
	result := player.Play("wav", "song.wav")

	expected := "Playing WAV file: song.wav"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestAudioPlayerWithInvalidFormat(t *testing.T) {
	player := &AudioPlayer{}
	result := player.Play("aac", "song.aac")

	if !strings.Contains(result, "Invalid media type") {
		t.Errorf("Expected error message containing 'Invalid media type', got '%s'", result)
	}
}

func TestFLACPlayer(t *testing.T) {
	player := &FLACPlayer{}
	
	// Test playing FLAC
	result := player.PlayFLAC("song.flac")
	expected := "Playing FLAC file: song.flac"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
	
	// Test playing WAV with FLAC player
	result = player.PlayWAV("song.wav")
	if !strings.Contains(result, "cannot play WAV") {
		t.Errorf("Expected error about not being able to play WAV, got '%s'", result)
	}
}

func TestWAVPlayer(t *testing.T) {
	player := &WAVPlayer{}
	
	// Test playing WAV
	result := player.PlayWAV("song.wav")
	expected := "Playing WAV file: song.wav"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
	
	// Test playing FLAC with WAV player
	result = player.PlayFLAC("song.flac")
	if !strings.Contains(result, "cannot play FLAC") {
		t.Errorf("Expected error about not being able to play FLAC, got '%s'", result)
	}
}

func TestMediaAdapter(t *testing.T) {
	// Test FLAC adapter
	flacAdapter := NewMediaAdapter("flac")
	if flacAdapter == nil {
		t.Error("Expected FLAC adapter to be created, got nil")
	} else {
		result := flacAdapter.Play("flac", "song.flac")
		expected := "Playing FLAC file: song.flac"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	}
	
	// Test WAV adapter
	wavAdapter := NewMediaAdapter("wav")
	if wavAdapter == nil {
		t.Error("Expected WAV adapter to be created, got nil")
	} else {
		result := wavAdapter.Play("wav", "song.wav")
		expected := "Playing WAV file: song.wav"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	}
	
	// Test invalid adapter
	invalidAdapter := NewMediaAdapter("mp3")
	if invalidAdapter != nil {
		t.Errorf("Expected nil adapter for unsupported format, got %v", invalidAdapter)
	}
}
