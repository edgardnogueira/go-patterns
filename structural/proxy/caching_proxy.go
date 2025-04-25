package proxy

import (
	"fmt"
	"sync"
	"time"
)

// CacheEntry represents a cached image with expiration time.
type CacheEntry struct {
	image       Image
	createdTime time.Time
	lastAccess  time.Time
	accessCount int
}

// CachingProxy implements a caching mechanism to avoid reloading images.
// It stores already loaded images in memory for quick access.
type CachingProxy struct {
	BaseProxy
	cache            map[string]*CacheEntry
	mu               sync.RWMutex
	expiration       time.Duration
	imageFactory     func(string) (Image, error)
	maxCacheSize     int
	currentCacheSize int
}

// NewCachingProxy creates a new caching proxy.
func NewCachingProxy(expiration time.Duration, maxCacheSize int) *CachingProxy {
	return &CachingProxy{
		cache:        make(map[string]*CacheEntry),
		expiration:   expiration,
		maxCacheSize: maxCacheSize,
		imageFactory: func(filename string) (Image, error) {
			return NewRealImage(filename)
		},
	}
}

// Display shows the image, retrieving it from cache if available.
func (p *CachingProxy) Display(filename string) error {
	image, err := p.getFromCache(filename)
	if err != nil {
		return err
	}
	
	fmt.Println("Caching proxy delegating display call")
	return image.Display()
}

// GetFilename is not applicable for the caching proxy as it manages multiple images.
func (p *CachingProxy) GetFilename() string {
	return ""
}

// GetWidth returns the width of a specific image.
func (p *CachingProxy) GetWidth(filename string) (int, error) {
	image, err := p.getFromCache(filename)
	if err != nil {
		return 0, err
	}
	return image.GetWidth(), nil
}

// GetHeight returns the height of a specific image.
func (p *CachingProxy) GetHeight(filename string) (int, error) {
	image, err := p.getFromCache(filename)
	if err != nil {
		return 0, err
	}
	return image.GetHeight(), nil
}

// GetSize returns the size of a specific image.
func (p *CachingProxy) GetSize(filename string) (int64, error) {
	image, err := p.getFromCache(filename)
	if err != nil {
		return 0, err
	}
	return image.GetSize(), nil
}

// GetMetadata returns the metadata of a specific image.
func (p *CachingProxy) GetMetadata(filename string) (map[string]string, error) {
	image, err := p.getFromCache(filename)
	if err != nil {
		return nil, err
	}
	return image.GetMetadata(), nil
}

// getFromCache retrieves an image from cache or loads it if not available.
func (p *CachingProxy) getFromCache(filename string) (Image, error) {
	// First check if the image is in the cache (read lock)
	p.mu.RLock()
	entry, found := p.cache[filename]
	p.mu.RUnlock()

	// If found and not expired, update stats and return it
	if found && !p.isExpired(entry) {
		p.mu.Lock()
		entry.lastAccess = time.Now()
		entry.accessCount++
		p.mu.Unlock()
		fmt.Printf("Cache hit for image: %s (access count: %d)\n", filename, entry.accessCount)
		return entry.image, nil
	}

	// If not found or expired, load it (write lock)
	p.mu.Lock()
	defer p.mu.Unlock()

	// Double-check in case another goroutine loaded it while we were waiting
	entry, found = p.cache[filename]
	if found && !p.isExpired(entry) {
		entry.lastAccess = time.Now()
		entry.accessCount++
		fmt.Printf("Cache hit for image: %s (access count: %d)\n", filename, entry.accessCount)
		return entry.image, nil
	}

	// If expired, remove it from cache
	if found && p.isExpired(entry) {
		fmt.Printf("Cache entry expired for image: %s\n", filename)
		delete(p.cache, filename)
		p.currentCacheSize--
	}

	// Not in cache or expired, load the image
	fmt.Printf("Cache miss for image: %s, loading...\n", filename)
	image, err := p.imageFactory(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load image: %w", err)
	}

	// If cache is full, evict the least recently used entry
	if p.currentCacheSize >= p.maxCacheSize {
		p.evictLRU()
	}

	// Add to cache
	p.cache[filename] = &CacheEntry{
		image:       image,
		createdTime: time.Now(),
		lastAccess:  time.Now(),
		accessCount: 1,
	}
	p.currentCacheSize++

	return image, nil
}

// isExpired checks if a cache entry has expired.
func (p *CachingProxy) isExpired(entry *CacheEntry) bool {
	return p.expiration > 0 && time.Since(entry.createdTime) > p.expiration
}

// evictLRU removes the least recently used cache entry.
func (p *CachingProxy) evictLRU() {
	var oldestKey string
	var oldestTime time.Time

	// Find the least recently accessed entry
	first := true
	for filename, entry := range p.cache {
		if first || entry.lastAccess.Before(oldestTime) {
			oldestKey = filename
			oldestTime = entry.lastAccess
			first = false
		}
	}

	if oldestKey != "" {
		fmt.Printf("Evicting least recently used cache entry: %s\n", oldestKey)
		delete(p.cache, oldestKey)
		p.currentCacheSize--
	}
}

// ClearCache removes all entries from the cache.
func (p *CachingProxy) ClearCache() {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	p.cache = make(map[string]*CacheEntry)
	p.currentCacheSize = 0
	fmt.Println("Cache cleared")
}

// GetCacheStats returns statistics about the cache.
func (p *CachingProxy) GetCacheStats() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	stats := map[string]interface{}{
		"size":          p.currentCacheSize,
		"max_size":      p.maxCacheSize,
		"expiration_ms": p.expiration.Milliseconds(),
	}
	
	// Add per-item stats
	itemStats := make(map[string]map[string]interface{})
	for filename, entry := range p.cache {
		itemStats[filename] = map[string]interface{}{
			"access_count":    entry.accessCount,
			"age_seconds":     time.Since(entry.createdTime).Seconds(),
			"last_access_ago": time.Since(entry.lastAccess).Seconds(),
		}
	}
	stats["items"] = itemStats
	
	return stats
}

// SetImageFactory allows customizing how images are created when not in cache.
func (p *CachingProxy) SetImageFactory(factory func(string) (Image, error)) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.imageFactory = factory
}
