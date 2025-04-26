// Package interfaces demonstrates idiomatic Go interface implementation patterns
package interfaces

import (
	"fmt"
	"io"
	"strings"
)

// Interface composition through embedding
// --------------------------------------

// Reader is an interface for types that can read data
type Reader interface {
	Read(p []byte) (n int, err error)
}

// Writer is an interface for types that can write data
type Writer interface {
	Write(p []byte) (n int, err error)
}

// ReadWriter composes the Reader and Writer interfaces
// This is interface embedding - a fundamental form of composition in Go
type ReadWriter interface {
	Reader
	Writer
}

// StringReadWriter is a concrete type that implements the ReadWriter interface
type StringReadWriter struct {
	data string
	pos  int
}

// Read implements the Reader interface
func (s *StringReadWriter) Read(p []byte) (n int, err error) {
	if s.pos >= len(s.data) {
		return 0, io.EOF
	}
	
	n = copy(p, []byte(s.data[s.pos:]))
	s.pos += n
	return n, nil
}

// Write implements the Writer interface
func (s *StringReadWriter) Write(p []byte) (n int, err error) {
	s.data += string(p)
	return len(p), nil
}

// CompositionDemo demonstrates interface composition
func CompositionDemo() {
	// Create an instance of StringReadWriter
	rw := &StringReadWriter{data: "Hello, World!"}
	
	// Use it as a Writer
	var w Writer = rw
	w.Write([]byte(" More text."))
	
	// Use it as a Reader
	var r Reader = rw
	buf := make([]byte, 10)
	n, _ := r.Read(buf)
	fmt.Printf("Read %d bytes: %s\n", n, buf[:n])
	
	// Use it as a ReadWriter (composed interface)
	var rwr ReadWriter = rw
	n, _ = rwr.Read(buf)
	fmt.Printf("Read %d more bytes: %s\n", n, buf[:n])
	rwr.Write([]byte(" Even more text."))
	
	// Display the final contents
	fmt.Printf("Final data: %s\n", rw.data)
}

// Multi-level interface composition
// --------------------------------

// Closer is an interface for types that can be closed
type Closer interface {
	Close() error
}

// ReadWriteCloser composes three interfaces
type ReadWriteCloser interface {
	Reader
	Writer
	Closer
}

// StringReadWriteCloser implements the ReadWriteCloser interface
type StringReadWriteCloser struct {
	StringReadWriter
	closed bool
}

// Close implements the Closer interface
func (s *StringReadWriteCloser) Close() error {
	s.closed = true
	return nil
}

// Read overrides the embedded Read method to check if closed
func (s *StringReadWriteCloser) Read(p []byte) (n int, err error) {
	if s.closed {
		return 0, fmt.Errorf("read from closed StringReadWriteCloser")
	}
	return s.StringReadWriter.Read(p)
}

// Write overrides the embedded Write method to check if closed
func (s *StringReadWriteCloser) Write(p []byte) (n int, err error) {
	if s.closed {
		return 0, fmt.Errorf("write to closed StringReadWriteCloser")
	}
	return s.StringReadWriter.Write(p)
}

// MultiLevelCompositionDemo demonstrates multi-level interface composition
func MultiLevelCompositionDemo() {
	// Create a StringReadWriteCloser
	rwc := &StringReadWriteCloser{
		StringReadWriter: StringReadWriter{data: "Initial text."},
	}
	
	// Use it through different interface views
	var r Reader = rwc
	var w Writer = rwc
	var c Closer = rwc
	var rw ReadWriter = rwc
	var rwc2 ReadWriteCloser = rwc
	
	// Demonstrate using all interfaces
	w.Write([]byte(" More content."))
	
	buf := make([]byte, 12)
	n, _ := r.Read(buf)
	fmt.Printf("Read through Reader: %s\n", buf[:n])
	
	rwc2.Write([]byte(" Final content."))
	
	// Reset position to read from the beginning
	rwc.StringReadWriter.pos = 0
	
	// Read full content through ReadWriter
	var sb strings.Builder
	buffer := make([]byte, 8) // Small buffer to demonstrate multiple reads
	
	for {
		n, err := rw.Read(buffer)
		if err == io.EOF {
			break
		}
		sb.Write(buffer[:n])
	}
	
	fmt.Printf("Full content: %s\n", sb.String())
	
	// Close and try to use
	c.Close()
	
	// This will return an error
	_, err := rwc2.Write([]byte("This will fail"))
	fmt.Printf("Write after close: %v\n", err)
}

// Extending standard library interfaces
// -----------------------------------

// LoggingReader adds logging to any Reader
type LoggingReader struct {
	r        Reader
	loggedOp string
}

// NewLoggingReader creates a new LoggingReader
func NewLoggingReader(r Reader, operation string) *LoggingReader {
	return &LoggingReader{r: r, loggedOp: operation}
}

// Read implements the Reader interface and adds logging
func (lr *LoggingReader) Read(p []byte) (n int, err error) {
	n, err = lr.r.Read(p)
	fmt.Printf("[%s] Read %d bytes, err: %v\n", lr.loggedOp, n, err)
	return n, err
}

// ExtendingInterfacesDemo demonstrates extending standard library interfaces
func ExtendingInterfacesDemo() {
	// Create a string reader
	sr := strings.NewReader("Hello, interface composition!")
	
	// Wrap it with our logging reader
	loggingReader := NewLoggingReader(sr, "STRING-READ")
	
	// Use it as a Reader
	buf := make([]byte, 10)
	for {
		n, err := loggingReader.Read(buf)
		if err == io.EOF {
			break
		}
		fmt.Printf("Got: %s\n", buf[:n])
	}
	
	// We can chain and compose wrappers
	sr2 := strings.NewReader("Another example")
	doubleLogger := NewLoggingReader(
		NewLoggingReader(sr2, "INNER"), 
		"OUTER",
	)
	
	fmt.Println("\nReading with double logger:")
	io.ReadAll(doubleLogger) // Use io.ReadAll to read everything
}

// Interface embedding in structs
// ----------------------------

// DataStore defines behaviors for storing and retrieving data
type DataStore interface {
	Get(key string) (string, bool)
	Set(key, value string)
}

// SimpleStore is a basic implementation of DataStore
type SimpleStore struct {
	data map[string]string
}

// NewSimpleStore creates a new SimpleStore
func NewSimpleStore() *SimpleStore {
	return &SimpleStore{
		data: make(map[string]string),
	}
}

// Get retrieves a value from the store
func (s *SimpleStore) Get(key string) (string, bool) {
	val, ok := s.data[key]
	return val, ok
}

// Set stores a value in the store
func (s *SimpleStore) Set(key, value string) {
	s.data[key] = value
}

// CachingService embeds a DataStore interface
// This demonstrates embedding an interface inside a struct
type CachingService struct {
	DataStore // Interface embedding in struct
	cache     map[string]string
	misses    int
}

// NewCachingService creates a new CachingService
func NewCachingService(store DataStore) *CachingService {
	return &CachingService{
		DataStore: store,
		cache:     make(map[string]string),
	}
}

// Get overrides the embedded Get method to add caching
func (c *CachingService) Get(key string) (string, bool) {
	// Check cache first
	if val, ok := c.cache[key]; ok {
		return val, true
	}
	
	// Call embedded implementation if not in cache
	val, ok := c.DataStore.Get(key)
	if ok {
		// Add to cache
		c.cache[key] = val
	} else {
		c.misses++
	}
	
	return val, ok
}

// CacheMisses returns the number of cache misses
func (c *CachingService) CacheMisses() int {
	return c.misses
}

// InterfaceEmbeddingDemo demonstrates embedding interfaces in structs
func InterfaceEmbeddingDemo() {
	// Create the base store
	store := NewSimpleStore()
	store.Set("key1", "value1")
	store.Set("key2", "value2")
	
	// Create the caching service that embeds the store
	service := NewCachingService(store)
	
	// First access (miss)
	val1, _ := service.Get("key1")
	fmt.Printf("First access for key1: %s, Misses: %d\n", val1, service.CacheMisses())
	
	// Second access (hit - should use cache)
	val1Again, _ := service.Get("key1")
	fmt.Printf("Second access for key1: %s, Misses: %d\n", val1Again, service.CacheMisses())
	
	// Access key2 (miss)
	val2, _ := service.Get("key2")
	fmt.Printf("First access for key2: %s, Misses: %d\n", val2, service.CacheMisses())
	
	// Access nonexistent key (miss)
	_, found := service.Get("key3")
	fmt.Printf("key3 found: %t, Misses: %d\n", found, service.CacheMisses())
	
	// We can still use Set from the embedded interface
	service.Set("key3", "value3")
	val3, _ := service.Get("key3")
	fmt.Printf("After setting key3: %s\n", val3)
}
