package facade

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// MediaConverterFacade is the facade that simplifies interactions with the complex media subsystems
type MediaConverterFacade struct {
	videoProcessor  *VideoProcessor
	audioProcessor  *AudioProcessor
	codecManager    *CodecManager
	fileSystem      *FileSystem
	metadataHandler *MetadataHandler
	progressReporter *ProgressReporter
	activeJobs      map[string]*ConversionJob
}

// NewMediaConverterFacade creates a new MediaConverterFacade
func NewMediaConverterFacade() *MediaConverterFacade {
	return &MediaConverterFacade{
		videoProcessor:   NewVideoProcessor(),
		audioProcessor:   NewAudioProcessor(),
		codecManager:     NewCodecManager(),
		fileSystem:       NewFileSystem(),
		metadataHandler:  NewMetadataHandler(),
		progressReporter: NewProgressReporter(),
		activeJobs:       make(map[string]*ConversionJob),
	}
}

// --------------------------------
// Simplified interface methods
// --------------------------------

// ConvertVideo converts a video file to a different format with sensible defaults
func (m *MediaConverterFacade) ConvertVideo(inputPath, outputPath string, format MediaFormat) error {
	// Validate input
	if !m.fileSystem.FileExists(inputPath) {
		return fmt.Errorf("input file does not exist: %s", inputPath)
	}
	
	// Create a job ID
	jobID := fmt.Sprintf("video-convert-%d", time.Now().Unix())
	
	// Get optimal codecs for the format
	videoCodec, audioCodec := m.codecManager.GetOptimalCodecForFormat(format)
	
	// Create a job
	job := &ConversionJob{
		ID:         jobID,
		InputPath:  inputPath,
		OutputPath: outputPath,
		OutputFormat: OutputFormat{
			Format:       format,
			VideoCodec:   videoCodec,
			AudioCodec:   audioCodec,
			Resolution:   Resolution1080p,
			Bitrate:      BitrateMedium,
			KeepMetadata: true,
		},
		StartTime: time.Now(),
		Status:    "started",
	}
	
	// Store the job
	m.activeJobs[jobID] = job
	
	// Register a progress listener
	m.progressReporter.RegisterListener(jobID, func(progress float64) {
		job.Progress = progress
		fmt.Printf("Video conversion progress: %.1f%%\n", progress*100)
	})
	
	// Create a progress callback
	progressCallback := m.progressReporter.CreateProgressCallback(jobID)
	
	// Start the conversion
	go func() {
		// Extract metadata before conversion
		metadata, err := m.metadataHandler.ExtractMetadata(inputPath)
		if err != nil {
			job.Status = "failed"
			job.Error = fmt.Errorf("failed to extract metadata: %w", err)
			return
		}
		
		// Convert the video
		err = m.videoProcessor.ConvertFormat(
			inputPath,
			outputPath,
			format,
			videoCodec,
			progressCallback,
		)
		
		if err != nil {
			job.Status = "failed"
			job.Error = fmt.Errorf("video conversion failed: %w", err)
			return
		}
		
		// Write metadata if requested
		if job.OutputFormat.KeepMetadata {
			err = m.metadataHandler.WriteMetadata(outputPath, metadata)
			if err != nil {
				fmt.Printf("Warning: Failed to write metadata: %v\n", err)
				// Continue anyway, this is not critical
			}
		}
		
		// Mark job as completed
		job.Status = "completed"
		job.EndTime = time.Now()
		job.Progress = 1.0
		
		// Unregister the progress listener
		m.progressReporter.UnregisterListener(jobID)
	}()
	
	return nil
}

// ExtractAudio extracts audio from a video file
func (m *MediaConverterFacade) ExtractAudio(videoPath, outputPath string, format MediaFormat) error {
	// Validate input
	if !m.fileSystem.FileExists(videoPath) {
		return fmt.Errorf("input file does not exist: %s", videoPath)
	}
	
	// Create a job ID
	jobID := fmt.Sprintf("audio-extract-%d", time.Now().Unix())
	
	// Create a progress callback
	progressCallback := m.progressReporter.CreateProgressCallback(jobID)
	
	// Start the extraction
	go func() {
		err := m.audioProcessor.ExtractAudio(
			videoPath,
			outputPath,
			format,
			progressCallback,
		)
		
		if err != nil {
			fmt.Printf("Audio extraction failed: %v\n", err)
		} else {
			fmt.Println("Audio extraction completed successfully")
		}
		
		// Unregister the progress listener
		m.progressReporter.UnregisterListener(jobID)
	}()
	
	return nil
}

// OptimizeForWeb prepares media for web streaming with optimal settings
func (m *MediaConverterFacade) OptimizeForWeb(inputPath, outputPath string) error {
	// Determine if it's video or audio based on extension
	inputExt := m.fileSystem.GetFileExtension(inputPath)
	
	// Determine output format
	var format MediaFormat
	isVideo := false
	
	switch inputExt {
	case ".mp4", ".avi", ".mkv", ".mov", ".wmv":
		format = FormatMP4
		isVideo = true
	case ".mp3", ".wav", ".flac", ".ogg":
		format = FormatMP3
		isVideo = false
	default:
		// Default to MP4 for unknown formats
		format = FormatMP4
		isVideo = true
	}
	
	// Create a job ID
	jobID := fmt.Sprintf("web-optimize-%d", time.Now().Unix())
	
	// Get optimal codecs for web streaming
	videoCodec, audioCodec := CodecH264, CodecAAC
	
	// Create a progress callback
	progressCallback := m.progressReporter.CreateProgressCallback(jobID)
	
	// Process based on content type
	go func() {
		if isVideo {
			err := m.videoProcessor.ConvertFormat(
				inputPath,
				outputPath,
				format,
				videoCodec,
				progressCallback,
			)
			
			if err != nil {
				fmt.Printf("Web optimization failed: %v\n", err)
			} else {
				fmt.Println("Video optimized for web successfully")
			}
		} else {
			err := m.audioProcessor.ConvertFormat(
				inputPath,
				outputPath,
				format,
				audioCodec,
				progressCallback,
			)
			
			if err != nil {
				fmt.Printf("Audio optimization failed: %v\n", err)
			} else {
				fmt.Println("Audio optimized for web successfully")
			}
		}
		
		// Unregister the progress listener
		m.progressReporter.UnregisterListener(jobID)
	}()
	
	return nil
}

// CreateThumbnail generates a thumbnail from a video file
func (m *MediaConverterFacade) CreateThumbnail(videoPath, outputPath string) error {
	// Validate input
	if !m.fileSystem.FileExists(videoPath) {
		return fmt.Errorf("input file does not exist: %s", videoPath)
	}
	
	// Create the thumbnail at a reasonable position (10% into the video)
	// In a real implementation, we would determine the video duration
	timeOffset := 10 * time.Second
	
	// Generate the thumbnail
	err := m.videoProcessor.CreateThumbnail(videoPath, outputPath, timeOffset)
	if err != nil {
		return fmt.Errorf("failed to create thumbnail: %w", err)
	}
	
	return nil
}

// BatchConvert converts multiple files with the same settings
func (m *MediaConverterFacade) BatchConvert(inputPaths []string, outputDir string, format MediaFormat) error {
	// Ensure output directory exists
	err := m.fileSystem.CreateDirectory(outputDir)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	
	// Process each file
	for _, inputPath := range inputPaths {
		// Generate output path based on input filename
		fileName := filepath.Base(inputPath)
		fileExt := filepath.Ext(fileName)
		baseName := strings.TrimSuffix(fileName, fileExt)
		outputPath := filepath.Join(outputDir, baseName+"."+string(format))
		
		// Convert the file
		inputExt := strings.ToLower(fileExt)
		
		// Determine if it's video or audio based on extension
		isVideo := false
		switch inputExt {
		case ".mp4", ".avi", ".mkv", ".mov", ".wmv":
			isVideo = true
		case ".mp3", ".wav", ".flac", ".ogg":
			isVideo = false
		default:
			// Skip unknown formats
			fmt.Printf("Skipping unknown format: %s\n", inputPath)
			continue
		}
		
		if isVideo {
			err := m.ConvertVideo(inputPath, outputPath, format)
			if err != nil {
				fmt.Printf("Failed to convert video %s: %v\n", inputPath, err)
				// Continue with the next file
				continue
			}
		} else {
			_, audioCodec := m.codecManager.GetOptimalCodecForFormat(format)
			
			// Create a job ID for tracking
			jobID := fmt.Sprintf("batch-audio-%d-%s", time.Now().Unix(), baseName)
			
			// Create a progress callback
			progressCallback := m.progressReporter.CreateProgressCallback(jobID)
			
			err := m.audioProcessor.ConvertFormat(
				inputPath,
				outputPath,
				format,
				audioCodec,
				progressCallback,
			)
			
			if err != nil {
				fmt.Printf("Failed to convert audio %s: %v\n", inputPath, err)
				// Continue with the next file
				continue
			}
		}
	}
	
	return nil
}

// --------------------------------
// Additional convenience methods
// --------------------------------

// GetConversionJob returns details about a specific conversion job
func (m *MediaConverterFacade) GetConversionJob(jobID string) (*ConversionJob, bool) {
	job, exists := m.activeJobs[jobID]
	return job, exists
}

// GetActiveJobs returns all active conversion jobs
func (m *MediaConverterFacade) GetActiveJobs() []*ConversionJob {
	jobs := make([]*ConversionJob, 0, len(m.activeJobs))
	for _, job := range m.activeJobs {
		jobs = append(jobs, job)
	}
	return jobs
}

// CancelJob cancels an active conversion job
func (m *MediaConverterFacade) CancelJob(jobID string) error {
	job, exists := m.activeJobs[jobID]
	if !exists {
		return fmt.Errorf("job not found: %s", jobID)
	}
	
	// Mark the job as cancelled
	job.Status = "cancelled"
	job.EndTime = time.Now()
	
	// Unregister the progress listener
	m.progressReporter.UnregisterListener(jobID)
	
	fmt.Printf("Job %s cancelled\n", jobID)
	return nil
}

// IsFormatSupported checks if a specific format is supported
func (m *MediaConverterFacade) IsFormatSupported(format MediaFormat, isVideo bool) bool {
	if isVideo {
		return m.videoProcessor.IsFormatSupported(format)
	}
	return m.audioProcessor.IsFormatSupported(format)
}

// GetSupportedFormats returns a list of supported formats
func (m *MediaConverterFacade) GetSupportedFormats(isVideo bool) []MediaFormat {
	if isVideo {
		return []MediaFormat{FormatMP4, FormatAVI, FormatMKV, FormatWEBM}
	}
	return []MediaFormat{FormatMP3, FormatAAC, FormatWAV}
}

// CreateConversionProfile creates a reusable conversion profile
func (m *MediaConverterFacade) CreateConversionProfile(name string, format MediaFormat, resolution ResolutionPreset, bitrate BitratePreset) OutputFormat {
	videoCodec, audioCodec := m.codecManager.GetOptimalCodecForFormat(format)
	
	return OutputFormat{
		Format:       format,
		VideoCodec:   videoCodec,
		AudioCodec:   audioCodec,
		Resolution:   resolution,
		Bitrate:      bitrate,
		KeepMetadata: true,
	}
}
