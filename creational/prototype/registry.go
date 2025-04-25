package prototype

import (
	"errors"
	"sync"
)

// DocumentRegistry is a registry that stores and manages document prototypes.
// It allows registering, retrieving, and cloning prototypes by name.
type DocumentRegistry struct {
	prototypes map[string]Prototype
	mutex      sync.RWMutex
}

// NewDocumentRegistry creates a new document registry.
func NewDocumentRegistry() *DocumentRegistry {
	return &DocumentRegistry{
		prototypes: make(map[string]Prototype),
	}
}

// Register adds a prototype to the registry.
func (r *DocumentRegistry) Register(name string, prototype Prototype) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.prototypes[name] = prototype
}

// Unregister removes a prototype from the registry.
func (r *DocumentRegistry) Unregister(name string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	delete(r.prototypes, name)
}

// Get retrieves a prototype from the registry without cloning it.
func (r *DocumentRegistry) Get(name string) (Prototype, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	prototype, exists := r.prototypes[name]
	if !exists {
		return nil, errors.New("prototype not found: " + name)
	}
	
	return prototype, nil
}

// Clone retrieves a shallow clone of a prototype from the registry.
func (r *DocumentRegistry) Clone(name string) (Prototype, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	prototype, exists := r.prototypes[name]
	if !exists {
		return nil, errors.New("prototype not found: " + name)
	}
	
	return prototype.Clone(), nil
}

// DeepClone retrieves a deep clone of a prototype from the registry.
func (r *DocumentRegistry) DeepClone(name string) (Prototype, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	prototype, exists := r.prototypes[name]
	if !exists {
		return nil, errors.New("prototype not found: " + name)
	}
	
	return prototype.DeepClone(), nil
}

// List returns a list of all prototype names in the registry.
func (r *DocumentRegistry) List() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	names := make([]string, 0, len(r.prototypes))
	for name := range r.prototypes {
		names = append(names, name)
	}
	
	return names
}

// Count returns the number of prototypes in the registry.
func (r *DocumentRegistry) Count() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	return len(r.prototypes)
}
