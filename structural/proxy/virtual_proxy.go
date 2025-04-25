package proxy

import (
	"fmt"
	"sync"
)

// VirtualProxy implements lazy loading of the real image.
// It postpones the creation and loading of the expensive real object until it's actually needed.
type VirtualProxy struct {
	filename  string
	realImage Image
	once      sync.Once
	mu        sync.Mutex
}

// NewVirtualProxy creates a new virtual proxy for lazy loading.
func NewVirtualProxy(filename string) *VirtualProxy {
	return &VirtualProxy{
		filename: filename,
	}
}

// lazyInit initializes the real image on first use.
func (p *VirtualProxy) lazyInit() error {
	var err error
	p.once.Do(func() {
		fmt.Printf("Virtual proxy initializing real image %s\n", p.filename)
		var realImg *RealImage
		realImg, err = NewRealImage(p.filename)
		if err == nil {
			p.realImage = realImg
		}
	})
	return err
}

// Display loads the image lazily only when it's first displayed.
func (p *VirtualProxy) Display() error {
	err := p.lazyInit()
	if err != nil {
		return fmt.Errorf("virtual proxy initialization error: %w", err)
	}
	
	fmt.Println("Virtual proxy delegating display call to real image")
	return p.realImage.Display()
}

// GetFilename returns the image's filename.
func (p *VirtualProxy) GetFilename() string {
	return p.filename
}

// GetWidth returns the image width, initializing if needed.
func (p *VirtualProxy) GetWidth() int {
	err := p.lazyInit()
	if err != nil {
		fmt.Printf("Error getting width: %v\n", err)
		return 0
	}
	return p.realImage.GetWidth()
}

// GetHeight returns the image height, initializing if needed.
func (p *VirtualProxy) GetHeight() int {
	err := p.lazyInit()
	if err != nil {
		fmt.Printf("Error getting height: %v\n", err)
		return 0
	}
	return p.realImage.GetHeight()
}

// GetSize returns the image size, initializing if needed.
func (p *VirtualProxy) GetSize() int64 {
	err := p.lazyInit()
	if err != nil {
		fmt.Printf("Error getting size: %v\n", err)
		return 0
	}
	return p.realImage.GetSize()
}

// GetMetadata returns the image metadata, initializing if needed.
func (p *VirtualProxy) GetMetadata() map[string]string {
	err := p.lazyInit()
	if err != nil {
		fmt.Printf("Error getting metadata: %v\n", err)
		return make(map[string]string)
	}
	return p.realImage.GetMetadata()
}

// IsLoaded checks if the real image has been loaded yet.
func (p *VirtualProxy) IsLoaded() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.realImage != nil
}
