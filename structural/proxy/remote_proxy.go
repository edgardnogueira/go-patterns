package proxy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// RemoteImage represents an image stored on a remote server.
type RemoteImage struct {
	Filename string            `json:"filename"`
	Width    int               `json:"width"`
	Height   int               `json:"height"`
	Size     int64             `json:"size"`
	Metadata map[string]string `json:"metadata"`
}

// RemoteProxy provides a local representation of an image on a remote server.
// It handles network communication and represents remote resources as if they were local.
type RemoteProxy struct {
	filename      string
	baseURL       string
	cachedImage   *RemoteImage
	mu            sync.RWMutex
	lastFetched   time.Time
	cacheDuration time.Duration
	httpClient    *http.Client
}

// NewRemoteProxy creates a new remote proxy for an image on a remote server.
func NewRemoteProxy(baseURL, filename string) *RemoteProxy {
	return &RemoteProxy{
		filename:      filename,
		baseURL:       baseURL,
		cacheDuration: 5 * time.Minute,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// setCacheTimeout sets how long the cached remote data remains valid.
func (p *RemoteProxy) SetCacheDuration(duration time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.cacheDuration = duration
}

// fetchRemoteData retrieves image information from the remote server.
func (p *RemoteProxy) fetchRemoteData() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Check if cache is still valid
	if p.cachedImage != nil && time.Since(p.lastFetched) < p.cacheDuration {
		return nil
	}

	fmt.Printf("Fetching remote image data for %s\n", p.filename)

	// Create the request URL
	reqURL, err := url.JoinPath(p.baseURL, "images", p.filename)
	if err != nil {
		return fmt.Errorf("failed to create request URL: %w", err)
	}

	// In a real implementation, this would make an actual HTTP request
	// Here we'll simulate it for the example
	var remoteImage *RemoteImage
	
	// Simulate network latency
	time.Sleep(300 * time.Millisecond)
	
	// Simulate the response based on the filename
	if p.filename == "not_found.jpg" {
		return fmt.Errorf("404 Not Found: image does not exist on remote server")
	}
	
	// Create a simulated remote response
	remoteImage = &RemoteImage{
		Filename: p.filename,
		Width:    1920,
		Height:   1080,
		Size:     3 * 1024 * 1024, // 3MB
		Metadata: map[string]string{
			"format": "JPEG",
			"created": time.Now().Format(time.RFC3339),
			"author": "Remote User",
			"server": "remote-image-server-01",
		},
	}

	p.cachedImage = remoteImage
	p.lastFetched = time.Now()
	
	fmt.Printf("Remote data fetched successfully for %s\n", p.filename)
	return nil
}

// Display shows the remote image (simulated).
func (p *RemoteProxy) Display() error {
	if err := p.fetchRemoteData(); err != nil {
		return fmt.Errorf("failed to display remote image: %w", err)
	}
	
	fmt.Printf("Displaying remote image: %s [%dx%d]\n", 
		p.filename, p.cachedImage.Width, p.cachedImage.Height)
	return nil
}

// GetFilename returns the image's filename.
func (p *RemoteProxy) GetFilename() string {
	return p.filename
}

// GetWidth returns the image width.
func (p *RemoteProxy) GetWidth() int {
	if err := p.fetchRemoteData(); err != nil {
		fmt.Printf("Error getting width: %v\n", err)
		return 0
	}
	return p.cachedImage.Width
}

// GetHeight returns the image height.
func (p *RemoteProxy) GetHeight() int {
	if err := p.fetchRemoteData(); err != nil {
		fmt.Printf("Error getting height: %v\n", err)
		return 0
	}
	return p.cachedImage.Height
}

// GetSize returns the image file size.
func (p *RemoteProxy) GetSize() int64 {
	if err := p.fetchRemoteData(); err != nil {
		fmt.Printf("Error getting size: %v\n", err)
		return 0
	}
	return p.cachedImage.Size
}

// GetMetadata returns the image metadata.
func (p *RemoteProxy) GetMetadata() map[string]string {
	if err := p.fetchRemoteData(); err != nil {
		fmt.Printf("Error getting metadata: %v\n", err)
		return make(map[string]string)
	}
	return p.cachedImage.Metadata
}

// IsDataCached checks if remote data is currently cached.
func (p *RemoteProxy) IsDataCached() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.cachedImage != nil && time.Since(p.lastFetched) < p.cacheDuration
}

// ClearCache invalidates the cached data, forcing a refresh on the next call.
func (p *RemoteProxy) ClearCache() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.cachedImage = nil
}

// GetCacheStatus returns information about the cache status.
func (p *RemoteProxy) GetCacheStatus() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	status := make(map[string]interface{})
	status["has_cached_data"] = p.cachedImage != nil
	status["cache_duration_seconds"] = p.cacheDuration.Seconds()
	
	if p.cachedImage != nil {
		status["last_fetched"] = p.lastFetched.Format(time.RFC3339)
		status["age_seconds"] = time.Since(p.lastFetched).Seconds()
		status["is_fresh"] = time.Since(p.lastFetched) < p.cacheDuration
		status["expires_in_seconds"] = p.cacheDuration.Seconds() - time.Since(p.lastFetched).Seconds()
	}
	
	return status
}
