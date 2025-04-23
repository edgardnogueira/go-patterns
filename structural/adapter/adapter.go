package adapter

// MediaPlayer is the target interface that the client expects
type MediaPlayer interface {
	Play(audioType, fileName string)
}

// AdvancedMediaPlayer is the adaptee interface
type AdvancedMediaPlayer interface {
	PlayFLAC(fileName string)
	PlayWAV(fileName string)
}

// FLACPlayer is a concrete implementation of AdvancedMediaPlayer for FLAC files
type FLACPlayer struct{}

// PlayFLAC plays FLAC files
func (p *FLACPlayer) PlayFLAC(fileName string) string {
	return "Playing FLAC file: " + fileName
}

// PlayWAV is not implemented for FLAC player
func (p *FLACPlayer) PlayWAV(fileName string) string {
	return "FLAC player cannot play WAV files"
}

// WAVPlayer is a concrete implementation of AdvancedMediaPlayer for WAV files
type WAVPlayer struct{}

// PlayFLAC is not implemented for WAV player
func (p *WAVPlayer) PlayFLAC(fileName string) string {
	return "WAV player cannot play FLAC files"
}

// PlayWAV plays WAV files
func (p *WAVPlayer) PlayWAV(fileName string) string {
	return "Playing WAV file: " + fileName
}

// MediaAdapter is the adapter that adapts AdvancedMediaPlayer to MediaPlayer
type MediaAdapter struct {
	AdvancedMediaPlayer AdvancedMediaPlayer
}

// NewMediaAdapter creates a new MediaAdapter for the given audio type
func NewMediaAdapter(audioType string) *MediaAdapter {
	switch audioType {
	case "flac":
		return &MediaAdapter{AdvancedMediaPlayer: &FLACPlayer{}}
	case "wav":
		return &MediaAdapter{AdvancedMediaPlayer: &WAVPlayer{}}
	default:
		return nil
	}
}

// Play adapts the AdvancedMediaPlayer to MediaPlayer interface
func (a *MediaAdapter) Play(audioType, fileName string) string {
	switch audioType {
	case "flac":
		return a.AdvancedMediaPlayer.PlayFLAC(fileName)
	case "wav":
		return a.AdvancedMediaPlayer.PlayWAV(fileName)
	default:
		return "Invalid media type"
	}
}

// AudioPlayer is the client that uses the MediaPlayer interface
type AudioPlayer struct{}

// Play plays the audio file
func (p *AudioPlayer) Play(audioType, fileName string) string {
	// Native support for mp3 format
	if audioType == "mp3" {
		return "Playing MP3 file: " + fileName
	}

	// For other formats, use adapter
	if audioType == "flac" || audioType == "wav" {
		adapter := NewMediaAdapter(audioType)
		if adapter != nil {
			return adapter.Play(audioType, fileName)
		}
	}

	return "Invalid media type: " + audioType + ", file not supported"
}
