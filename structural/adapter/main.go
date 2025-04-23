package adapter

import (
	"fmt"
)

// This file contains example usage of the adapter pattern

// ExampleAdapterPattern demonstrates the adapter pattern in action
func ExampleAdapterPattern() {
	audioPlayer := &AudioPlayer{}

	// Play MP3 file - native support
	result := audioPlayer.Play("mp3", "beyond_the_horizon.mp3")
	fmt.Println(result)

	// Play FLAC file - using adapter
	result = audioPlayer.Play("flac", "alone.flac")
	fmt.Println(result)

	// Play WAV file - using adapter
	result = audioPlayer.Play("wav", "far_away.wav")
	fmt.Println(result)

	// Try to play unsupported format
	result = audioPlayer.Play("aac", "mind_me.aac")
	fmt.Println(result)
}
