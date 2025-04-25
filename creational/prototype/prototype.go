// Package prototype implements the Prototype design pattern in Go.
//
// The Prototype pattern is a creational design pattern that allows creating new objects by copying 
// existing ones without coupling to their specific classes. This implementation demonstrates a document
// generation system where documents can be created by cloning prototypes.
package prototype

// Prototype defines the interface for objects that can be cloned.
type Prototype interface {
	// Clone creates a copy of the object.
	Clone() Prototype
	
	// DeepClone creates a deep copy of the object.
	DeepClone() Prototype
}

// Document is an abstract base struct that represents a document.
// It implements the Prototype interface and provides common functionality for all documents.
type Document struct {
	ID       string
	Name     string
	Creator  string
	Created  string // ISO 8601 date format
	Modified string // ISO 8601 date format
	Tags     []string
	Metadata map[string]string
}

// Clone creates a shallow copy of the Document.
func (d *Document) Clone() Prototype {
	// Create a simple copy of the document
	cloned := &Document{
		ID:       d.ID,
		Name:     d.Name + " (Copy)",
		Creator:  d.Creator,
		Created:  d.Created,
		Modified: d.Modified,
		// For shallow copy, we just copy the slice reference
		Tags:     d.Tags,
		// For shallow copy, we just copy the map reference
		Metadata: d.Metadata,
	}
	
	return cloned
}

// DeepClone creates a deep copy of the Document.
func (d *Document) DeepClone() Prototype {
	// Create a deep copy of the document
	cloned := &Document{
		ID:       d.ID,
		Name:     d.Name + " (Copy)",
		Creator:  d.Creator,
		Created:  d.Created,
		Modified: d.Modified,
		// For deep copy, we create a new slice with the same content
		Tags:     make([]string, len(d.Tags)),
		// For deep copy, we create a new map with the same content
		Metadata: make(map[string]string, len(d.Metadata)),
	}
	
	// Copy the tags
	copy(cloned.Tags, d.Tags)
	
	// Copy the metadata
	for k, v := range d.Metadata {
		cloned.Metadata[k] = v
	}
	
	return cloned
}
