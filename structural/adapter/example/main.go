package main

import (
	"fmt"
	"github.com/edgardnogueira/go-patterns/structural/adapter"
)

func main() {
	fmt.Println("Adapter Pattern Example")
	fmt.Println("=======================")
	fmt.Println("This example demonstrates adapting advanced media players to work with a"
		+ "\nsimple media player interface that only natively supports MP3.")
	fmt.Println()

	// Create audio player
	audioPlayer := &adapter.AudioPlayer{}

	// Play different file formats
	fmt.Println("\nTesting different audio formats:")
	fmt.Println("--------------------------")

	fmt.Println("1. Play MP3 file (native support):")
	result := audioPlayer.Play("mp3", "beyond_the_horizon.mp3")
	fmt.Printf("   » %s\n", result)

	fmt.Println("\n2. Play FLAC file (using adapter):")
	result = audioPlayer.Play("flac", "alone.flac")
	fmt.Printf("   » %s\n", result)

	fmt.Println("\n3. Play WAV file (using adapter):")
	result = audioPlayer.Play("wav", "far_away.wav")
	fmt.Printf("   » %s\n", result)

	fmt.Println("\n4. Try to play unsupported format (AAC):")
	result = audioPlayer.Play("aac", "mind_me.aac")
	fmt.Printf("   » %s\n", result)

	fmt.Println("\nThe Adapter pattern allows our AudioPlayer to work with advanced")
	fmt.Println("media formats without changing its interface. The client code")
	fmt.Println("doesn't need to know about the adapters or advanced players.")
}
