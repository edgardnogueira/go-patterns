package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/edgardnogueira/go-patterns/structural/facade"
)

func main() {
	fmt.Println("=== Facade Pattern Example - Media Converter ===")
	fmt.Println("This example demonstrates how the Facade pattern simplifies interactions with complex subsystems.")
	fmt.Println("The MediaConverterFacade provides a simplified interface to various media processing components.")
	
	// Create the facade that hides the complexity of the subsystems
	converter := facade.NewMediaConverterFacade()
	
	// Example 1: Basic Video Conversion
	fmt.Println("\n=== Example 1: Basic Video Conversion ===")
	
	// In a real application, these would be actual file paths
	inputVideo := "example_files/sample_video.mp4"
	outputVideo := "example_files/converted_video.mkv"
	
	// The facade hides all the complexity of the conversion process
	fmt.Printf("Converting %s to %s...\n", inputVideo, outputVideo)
	err := converter.ConvertVideo(inputVideo, outputVideo, facade.FormatMKV)
	
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		// In a real application, we would handle this error
		// For this example, we'll simulate a successful conversion
		simulateConversion()
	} else {
		// Allow some time for the asynchronous conversion to show progress
		time.Sleep(400 * time.Millisecond)
	}
	
	// Example 2: Audio Extraction
	fmt.Println("\n=== Example 2: Audio Extraction ===")
	
	inputVideoForAudio := "example_files/sample_video.mp4"
	outputAudio := "example_files/extracted_audio.mp3"
	
	// The facade simplifies the process of extracting audio from a video
	fmt.Printf("Extracting audio from %s to %s...\n", inputVideoForAudio, outputAudio)
	err = converter.ExtractAudio(inputVideoForAudio, outputAudio, facade.FormatMP3)
	
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		// Simulate a successful extraction
		simulateAudioExtraction()
	} else {
		// Allow some time for the asynchronous extraction to show progress
		time.Sleep(200 * time.Millisecond)
	}
	
	// Example 3: Web Optimization
	fmt.Println("\n=== Example 3: Web Optimization ===")
	
	inputFile := "example_files/high_quality_video.mp4"
	outputWebFile := "example_files/web_optimized.mp4"
	
	// The facade handles all the complex settings needed for web optimization
	fmt.Printf("Optimizing %s for web streaming...\n", inputFile)
	err = converter.OptimizeForWeb(inputFile, outputWebFile)
	
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		// Simulate a successful optimization
		simulateWebOptimization()
	} else {
		// Allow some time for the asynchronous optimization to show progress
		time.Sleep(200 * time.Millisecond)
	}
	
	// Example 4: Batch Processing
	fmt.Println("\n=== Example 4: Batch Processing ===")
	
	inputFiles := []string{
		"example_files/video1.avi",
		"example_files/video2.mp4",
		"example_files/audio1.wav",
	}
	outputDir := "example_files/batch_output"
	
	// The facade simplifies batch processing of multiple files
	fmt.Printf("Converting %d files to MP4 format...\n", len(inputFiles))
	err = converter.BatchConvert(inputFiles, outputDir, facade.FormatMP4)
	
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		// Simulate a successful batch conversion
		simulateBatchConversion(inputFiles)
	} else {
		// Allow some time for the asynchronous batch conversion to show progress
		time.Sleep(300 * time.Millisecond)
	}
	
	// Example 5: Using Conversion Profiles
	fmt.Println("\n=== Example 5: Using Conversion Profiles ===")
	
	// Create a reusable profile
	hdProfile := converter.CreateConversionProfile(
		"HD Video",
		facade.FormatMP4,
		facade.Resolution1080p,
		facade.BitrateHigh,
	)
	
	// Display profile details
	fmt.Printf("Created HD Video Profile:\n")
	fmt.Printf("  Format: %s\n", hdProfile.Format)
	fmt.Printf("  Video Codec: %s\n", hdProfile.VideoCodec)
	fmt.Printf("  Audio Codec: %s\n", hdProfile.AudioCodec)
	fmt.Printf("  Resolution: %s\n", hdProfile.Resolution)
	fmt.Printf("  Bitrate: %s\n", hdProfile.Bitrate)
	
	// Example 6: Thumbnail Creation
	fmt.Println("\n=== Example 6: Thumbnail Creation ===")
	
	inputVideoForThumbnail := "example_files/sample_video.mp4"
	outputThumbnail := "example_files/thumbnail.jpg"
	
	// The facade simplifies the process of creating a thumbnail
	fmt.Printf("Creating thumbnail from %s...\n", inputVideoForThumbnail)
	err = converter.CreateThumbnail(inputVideoForThumbnail, outputThumbnail)
	
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		// Simulate a successful thumbnail creation
		simulateThumbnailCreation()
	}
	
	// Display Facade Pattern Benefits
	fmt.Println("\n=== Facade Pattern Benefits ===")
	fmt.Println("1. Simplifies complex subsystem interactions")
	fmt.Println("2. Provides a unified interface to a set of interfaces")
	fmt.Println("3. Decouples clients from subsystem components")
	fmt.Println("4. Makes subsystems easier to use")
	fmt.Println("5. Promotes loose coupling between subsystems and clients")
	
	// Without vs. With Facade Comparison
	fmt.Println("\n=== Comparison: Without vs. With Facade ===")
	
	fmt.Println("\n--- Without Facade ---")
	fmt.Println("// Initialize all subsystems separately")
	fmt.Println("videoProcessor := NewVideoProcessor()")
	fmt.Println("audioProcessor := NewAudioProcessor()")
	fmt.Println("codecManager := NewCodecManager()")
	fmt.Println("fileSystem := NewFileSystem()")
	fmt.Println("metadataHandler := NewMetadataHandler()")
	fmt.Println("progressReporter := NewProgressReporter()")
	fmt.Println("")
	fmt.Println("// Check if format and codec are supported")
	fmt.Println("if !videoProcessor.IsFormatSupported(FormatMP4) {")
	fmt.Println("    return errors.New(\"format not supported\")")
	fmt.Println("}")
	fmt.Println("if !videoProcessor.IsCodecSupported(CodecH264) {")
	fmt.Println("    return errors.New(\"codec not supported\")")
	fmt.Println("}")
	fmt.Println("")
	fmt.Println("// Extract metadata before conversion")
	fmt.Println("metadata, err := metadataHandler.ExtractMetadata(inputPath)")
	fmt.Println("if err != nil {")
	fmt.Println("    return err")
	fmt.Println("}")
	fmt.Println("")
	fmt.Println("// Set up progress tracking")
	fmt.Println("jobID := fmt.Sprintf(\"job-%d\", time.Now().Unix())")
	fmt.Println("progressCallback := progressReporter.CreateProgressCallback(jobID)")
	fmt.Println("")
	fmt.Println("// Perform the conversion")
	fmt.Println("err = videoProcessor.ConvertFormat(inputPath, outputPath, format, codec, progressCallback)")
	fmt.Println("if err != nil {")
	fmt.Println("    return err")
	fmt.Println("}")
	fmt.Println("")
	fmt.Println("// Write metadata to the output file")
	fmt.Println("err = metadataHandler.WriteMetadata(outputPath, metadata)")
	fmt.Println("if err != nil {")
	fmt.Println("    return err")
	fmt.Println("}")
	
	fmt.Println("\n--- With Facade ---")
	fmt.Println("// Initialize the facade")
	fmt.Println("converter := NewMediaConverterFacade()")
	fmt.Println("")
	fmt.Println("// Perform the conversion in one line")
	fmt.Println("err := converter.ConvertVideo(inputPath, outputPath, FormatMP4)")
	fmt.Println("if err != nil {")
	fmt.Println("    return err")
	fmt.Println("}")
}

// Simulation functions for the example

func simulateConversion() {
	for i := 0; i <= 100; i += 20 {
		fmt.Printf("Video conversion progress: %d%%\n", i)
		time.Sleep(50 * time.Millisecond)
	}
	fmt.Println("Video conversion completed successfully!")
}

func simulateAudioExtraction() {
	for i := 0; i <= 100; i += 25 {
		fmt.Printf("Audio extraction progress: %d%%\n", i)
		time.Sleep(30 * time.Millisecond)
	}
	fmt.Println("Audio extraction completed successfully!")
}

func simulateWebOptimization() {
	fmt.Println("Analyzing video...")
	time.Sleep(50 * time.Millisecond)
	fmt.Println("Optimizing video codec...")
	time.Sleep(50 * time.Millisecond)
	fmt.Println("Optimizing audio codec...")
	time.Sleep(50 * time.Millisecond)
	fmt.Println("Adjusting bitrate for streaming...")
	time.Sleep(50 * time.Millisecond)
	fmt.Println("Web optimization completed successfully!")
}

func simulateBatchConversion(files []string) {
	for i, file := range files {
		fmt.Printf("Processing file %d of %d: %s\n", i+1, len(files), file)
		for j := 0; j <= 100; j += 50 {
			fmt.Printf("  Progress: %d%%\n", j)
			time.Sleep(20 * time.Millisecond)
		}
		fmt.Printf("  Completed: %s\n", filepath.Base(file))
	}
	fmt.Println("Batch conversion completed successfully!")
}

func simulateThumbnailCreation() {
	fmt.Println("Analyzing video...")
	time.Sleep(30 * time.Millisecond)
	fmt.Println("Selecting frame at 00:00:10...")
	time.Sleep(30 * time.Millisecond)
	fmt.Println("Extracting frame...")
	time.Sleep(30 * time.Millisecond)
	fmt.Println("Optimizing image...")
	time.Sleep(30 * time.Millisecond)
	fmt.Println("Thumbnail created successfully!")
}
