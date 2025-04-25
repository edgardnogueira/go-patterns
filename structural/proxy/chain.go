package proxy

import (
	"fmt"
	"time"
)

// ProxyChain helps create and manage chains of proxies.
// It allows building complex proxy chains with a fluent interface.
type ProxyChain struct {
	// The current chain of proxies
	current Image
}

// NewProxyChain creates a new proxy chain starting with a specific image.
func NewProxyChain(image Image) *ProxyChain {
	return &ProxyChain{
		current: image,
	}
}

// AddVirtual adds a virtual proxy to the chain for lazy loading.
func (c *ProxyChain) AddVirtual() *ProxyChain {
	// For VirtualProxy, we need to use the filename
	// Since we're in a chain, we can't directly create a VirtualProxy
	// Instead, we'll create a proxy that delegates to VirtualProxy
	filename := c.current.GetFilename()
	virtualProxy := NewVirtualProxy(filename)
	
	// Create a proxy that first delegates to the virtual proxy, but replaces
	// the real subject with our current chain when the virtual proxy initializes
	proxy := &chainedVirtualProxy{
		VirtualProxy: virtualProxy,
		replacement:  c.current,
	}
	
	c.current = proxy
	return c
}

// AddLogging adds a logging proxy to the chain.
func (c *ProxyChain) AddLogging(level LogLevel) *ProxyChain {
	c.current = NewLoggingProxy(c.current, level)
	return c
}

// AddMetrics adds a metrics proxy to the chain.
func (c *ProxyChain) AddMetrics() *ProxyChain {
	c.current = NewMetricsProxy(c.current)
	return c
}

// AddProtection adds a protection proxy to the chain.
func (c *ProxyChain) AddProtection(user *User) *ProxyChain {
	c.current = NewProtectionProxy(c.current, user)
	return c
}

// AddCaching adds a caching proxy to the chain.
// This is a special case since CachingProxy works differently
func (c *ProxyChain) AddCaching(expiration time.Duration, maxSize int) *ProxyChainWithCaching {
	cachingProxy := NewCachingProxy(expiration, maxSize)
	
	// Set a custom image factory that returns our chain
	filename := c.current.GetFilename()
	cachingProxy.SetImageFactory(func(f string) (Image, error) {
		if f == filename {
			return c.current, nil
		}
		return nil, fmt.Errorf("image not found: %s", f)
	})
	
	return &ProxyChainWithCaching{
		cachingProxy: cachingProxy,
		filename:     filename,
	}
}

// Build returns the final proxy chain as an Image.
func (c *ProxyChain) Build() Image {
	return c.current
}

// ProxyChainWithCaching is a special case for chains that end with a caching proxy.
type ProxyChainWithCaching struct {
	cachingProxy *CachingProxy
	filename     string
}

// Display delegates to the caching proxy with the filename.
func (c *ProxyChainWithCaching) Display() error {
	return c.cachingProxy.Display(c.filename)
}

// GetWidth delegates to the caching proxy with the filename.
func (c *ProxyChainWithCaching) GetWidth() int {
	width, err := c.cachingProxy.GetWidth(c.filename)
	if err != nil {
		return 0
	}
	return width
}

// GetHeight delegates to the caching proxy with the filename.
func (c *ProxyChainWithCaching) GetHeight() int {
	height, err := c.cachingProxy.GetHeight(c.filename)
	if err != nil {
		return 0
	}
	return height
}

// GetSize delegates to the caching proxy with the filename.
func (c *ProxyChainWithCaching) GetSize() int64 {
	size, err := c.cachingProxy.GetSize(c.filename)
	if err != nil {
		return 0
	}
	return size
}

// GetMetadata delegates to the caching proxy with the filename.
func (c *ProxyChainWithCaching) GetMetadata() map[string]string {
	metadata, err := c.cachingProxy.GetMetadata(c.filename)
	if err != nil {
		return make(map[string]string)
	}
	return metadata
}

// GetFilename returns the filename.
func (c *ProxyChainWithCaching) GetFilename() string {
	return c.filename
}

// Helper for chaining a VirtualProxy with existing proxies
type chainedVirtualProxy struct {
	*VirtualProxy
	replacement Image
}

// Override lazyInit to replace the realImage with our chain
func (p *chainedVirtualProxy) lazyInit() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if p.realImage == nil {
		fmt.Printf("Virtual proxy in chain initializing for %s\n", p.filename)
		p.realImage = p.replacement
	}
	return nil
}

// ProxyPresets provides common preset proxy chains for convenience.
type ProxyPresets struct{}

// NewPerformanceMonitoring creates a preset chain for performance monitoring.
// It adds logging and metrics proxies to track performance.
func (ProxyPresets) NewPerformanceMonitoring(realImage Image) Image {
	return NewProxyChain(realImage).
		AddLogging(INFO).
		AddMetrics().
		Build()
}

// NewSecure creates a preset chain for secure access.
// It adds protection and logging for security monitoring.
func (ProxyPresets) NewSecure(realImage Image, user *User) Image {
	return NewProxyChain(realImage).
		AddProtection(user).
		AddLogging(WARNING). // Log at warning level for security events
		Build()
}

// NewOptimized creates a preset chain for optimized access.
// It adds virtual loading and caching for better performance.
func (ProxyPresets) NewOptimized(filename string) *ProxyChainWithCaching {
	virtualProxy := NewVirtualProxy(filename)
	
	return NewProxyChain(virtualProxy).
		AddCaching(10*time.Minute, 100)
}
