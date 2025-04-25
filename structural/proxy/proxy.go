// Package proxy implements the Proxy design pattern in Go.
//
// The Proxy pattern provides a surrogate or placeholder for another object to control access to it.
// This implementation demonstrates various types of proxies for an image loading system
// that can provide protection, caching, lazy loading, and remote access capabilities.
package proxy

import (
	"fmt"
	"sync"
	"time"
)

// Image defines the interface for both the RealImage and its proxies.
// This is the Subject interface in the Proxy pattern.
type Image interface {
	// Display renders the image on screen
	Display() error

	// GetFilename returns the image's filename
	GetFilename() string

	// GetWidth returns the image width in pixels
	GetWidth() int

	// GetHeight returns the image height in pixels
	GetHeight() int

	// GetSize returns the image file size in bytes
	GetSize() int64

	// GetMetadata returns a map of image metadata
	GetMetadata() map[string]string
}

// RealImage is the concrete implementation of the Image interface.
// This is the RealSubject in the Proxy pattern that does the actual work.
type RealImage struct {
	filename string
	width    int
	height   int
	size     int64
	metadata map[string]string
	data     []byte // Simulates the actual image data in memory
	loaded   bool
}

// NewRealImage creates a new RealImage instance and loads the image from disk.
func NewRealImage(filename string) (*RealImage, error) {
	image := &RealImage{
		filename: filename,
		metadata: make(map[string]string),
		loaded:   false,
	}

	// Simulate loading the image from disk
	err := image.loadFromDisk()
	if err != nil {
		return nil, err
	}

	return image, nil
}

// loadFromDisk simulates loading the image data from a file.
func (r *RealImage) loadFromDisk() error {
	fmt.Printf("Loading image %s from disk\n", r.filename)

	// Simulate a delay for disk I/O
	time.Sleep(200 * time.Millisecond)

	// Simulate image properties based on the filename
	// In a real implementation, this would be determined from the actual file
	if r.filename == "not_found.jpg" {
		return fmt.Errorf("image not found: %s", r.filename)
	}

	// Set some example image properties
	r.width = 1920
	r.height = 1080
	r.size = 2048 * 1024 // 2MB
	
	// Set some example metadata
	r.metadata["format"] = "JPEG"
	r.metadata["created"] = time.Now().Format(time.RFC3339)
	r.metadata["colorspace"] = "RGB"

	// Simulate image data (in a real app, this would be the actual pixels)
	r.data = make([]byte, r.size)
	
	r.loaded = true
	fmt.Printf("Image %s loaded successfully\n", r.filename)

	return nil
}

// Display shows the image on screen (simulated).
func (r *RealImage) Display() error {
	if !r.loaded {
		err := r.loadFromDisk()
		if err != nil {
			return err
		}
	}

	fmt.Printf("Displaying image: %s [%dx%d]\n", r.filename, r.width, r.height)
	return nil
}

// GetFilename returns the image's filename.
func (r *RealImage) GetFilename() string {
	return r.filename
}

// GetWidth returns the image width in pixels.
func (r *RealImage) GetWidth() int {
	return r.width
}

// GetHeight returns the image height in pixels.
func (r *RealImage) GetHeight() int {
	return r.height
}

// GetSize returns the image file size in bytes.
func (r *RealImage) GetSize() int64 {
	return r.size
}

// GetMetadata returns a map of image metadata.
func (r *RealImage) GetMetadata() map[string]string {
	return r.metadata
}

// BaseProxy is a base struct that can be embedded by concrete proxy implementations.
// It implements forwarding methods to the real subject.
type BaseProxy struct {
	realImage Image
}

// GetFilename forwards the call to the real subject.
func (p *BaseProxy) GetFilename() string {
	return p.realImage.GetFilename()
}

// GetWidth forwards the call to the real subject.
func (p *BaseProxy) GetWidth() int {
	return p.realImage.GetWidth()
}

// GetHeight forwards the call to the real subject.
func (p *BaseProxy) GetHeight() int {
	return p.realImage.GetHeight()
}

// GetSize forwards the call to the real subject.
func (p *BaseProxy) GetSize() int64 {
	return p.realImage.GetSize()
}

// GetMetadata forwards the call to the real subject.
func (p *BaseProxy) GetMetadata() map[string]string {
	return p.realImage.GetMetadata()
}

// Display is implemented by concrete proxy types to add behavior.
func (p *BaseProxy) Display() error {
	// This should be overridden by concrete proxy types
	return p.realImage.Display()
}
