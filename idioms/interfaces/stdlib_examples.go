// Package interfaces demonstrates idiomatic Go interface implementation patterns.
package interfaces

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

// StdlibExamples demonstrates real-world interface usage from Go's standard library.
// The Go standard library makes extensive use of interfaces for flexibility and
// composition. This file showcases key interfaces from the standard library and
// how they are used effectively.

// Example 1: io.Reader and io.Writer
// ----------------------------------
// These interfaces are among the most widely used in Go's standard library.
// They provide a common contract for reading from and writing to different sources.

// ReaderWriterExample shows the power of the io.Reader and io.Writer interfaces
func ReaderWriterExample() {
	fmt.Println("=== io.Reader and io.Writer Examples ===")

	// Create a simple string source
	source := strings.NewReader("Hello, Gophers! This is an example of io interfaces.")

	// Chain readers together for processing
	// gzip -> buffered -> original source
	var compressed bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressed)
	
	// Copy from source to gzip writer
	_, err := io.Copy(gzipWriter, source)
	if err != nil {
		fmt.Println("Error compressing:", err)
		return
	}
	gzipWriter.Close()
	
	fmt.Printf("Original size: %d bytes\n", source.Size())
	fmt.Printf("Compressed size: %d bytes\n", compressed.Len())
	
	// Decompress by chaining readers in the opposite direction
	// original source -> gzip reader -> bufio scanner
	gzipReader, err := gzip.NewReader(&compressed)
	if err != nil {
		fmt.Println("Error creating gzip reader:", err)
		return
	}
	defer gzipReader.Close()
	
	scanner := bufio.NewScanner(gzipReader)
	fmt.Println("Decompressed content:")
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

// Example 2: sort.Interface
// -------------------------
// The sort package defines the Interface interface to allow sorting of
// any collection without knowing its underlying type.

// Person represents a person with a name and age
type Person struct {
	Name string
	Age  int
}

// ByAge implements sort.Interface for []Person based on Age field
type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }

// ByName implements sort.Interface for []Person based on Name field
type ByName []Person

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

// SortInterfaceExample demonstrates the power of sort.Interface
func SortInterfaceExample() {
	fmt.Println("\n=== sort.Interface Example ===")

	people := []Person{
		{"Alice", 25},
		{"Bob", 30},
		{"Charlie", 20},
		{"David", 35},
	}

	// Create a copy for each sort
	byAge := make([]Person, len(people))
	byName := make([]Person, len(people))
	copy(byAge, people)
	copy(byName, people)

	// Sort by age
	sort.Sort(ByAge(byAge))
	fmt.Println("Sorted by age:")
	for _, p := range byAge {
		fmt.Printf("  %s: %d\n", p.Name, p.Age)
	}

	// Sort by name
	sort.Sort(ByName(byName))
	fmt.Println("Sorted by name:")
	for _, p := range byName {
		fmt.Printf("  %s: %d\n", p.Name, p.Age)
	}

	// Using sort.Slice as a more modern alternative (Go 1.8+)
	// This is a convenience function that uses functional approaches
	// rather than requiring interface implementation
	sort.Slice(people, func(i, j int) bool {
		return people[i].Age > people[j].Age
	})
	fmt.Println("Sorted by age (descending) using sort.Slice:")
	for _, p := range people {
		fmt.Printf("  %s: %d\n", p.Name, p.Age)
	}
}

// Example 3: http.Handler
// ----------------------
// The http package uses the Handler interface to process HTTP requests.
// This is a cornerstone of Go's web development ecosystem.

// HelloHandler is a simple implementation of http.Handler
type HelloHandler struct {
	Greeting string
}

func (h HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s, %s!\n", h.Greeting, r.URL.Path[1:])
}

// LoggingMiddleware wraps an http.Handler and logs each request
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request
		fmt.Printf("[%s] %s %s\n", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)
		
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// HttpHandlerExample demonstrates the http.Handler interface
func HttpHandlerExample() {
	fmt.Println("\n=== http.Handler Example ===")
	
	// Create a new ServeMux to register our handlers
	mux := http.NewServeMux()
	
	// Register a HelloHandler
	hello := &HelloHandler{Greeting: "Hello"}
	
	// Wrap the handler with our logging middleware
	loggingHello := LoggingMiddleware(hello)
	mux.Handle("/hello/", loggingHello)
	
	// Add a handler using the HandlerFunc adapter
	mux.HandleFunc("/goodbye/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Goodbye, %s!", r.URL.Path[9:])
	})
	
	fmt.Println("If this were a real server, it would be listening on port 8080 now...")
	fmt.Println("GET /hello/world would trigger LoggingMiddleware, then HelloHandler")
	fmt.Println("GET /goodbye/world would trigger the goodbye HandlerFunc")
	
	// In a real application, you would start the server:
	// http.ListenAndServe(":8080", mux)
}

// Example 4: json.Marshaler/Unmarshaler
// --------------------------------------
// The encoding/json package uses interfaces to allow custom 
// serialization/deserialization behavior.

// CustomTime demonstrates implementing json.Marshaler and json.Unmarshaler
type CustomTime struct {
	time.Time
}

// MarshalJSON implements json.Marshaler for custom time formatting
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	// Use a custom format for the timestamp
	stamp := fmt.Sprintf("\"%s\"", ct.Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

// UnmarshalJSON implements json.Unmarshaler for custom time parsing
func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	// Remove the quotes
	str := string(data)
	str = strings.Trim(str, "\"")
	
	// Parse the date in our custom format
	t, err := time.Parse("2006-01-02 15:04:05", str)
	if err != nil {
		return err
	}
	
	ct.Time = t
	return nil
}

// Event represents an event with a title and timestamp
type Event struct {
	Title   string     `json:"title"`
	Time    CustomTime `json:"time"`
}

// JSONInterfaceExample demonstrates json.Marshaler and json.Unmarshaler
func JSONInterfaceExample() {
	fmt.Println("\n=== json.Marshaler/Unmarshaler Example ===")
	
	// Create an event with our custom time
	event := Event{
		Title: "Team Meeting",
		Time:  CustomTime{Time: time.Date(2025, 4, 26, 14, 30, 0, 0, time.UTC)},
	}
	
	// Marshal to JSON with our custom formatting
	jsonData, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	
	fmt.Println("Custom JSON format:")
	fmt.Println(string(jsonData))
	
	// Unmarshal back to a struct
	var decodedEvent Event
	err = json.Unmarshal(jsonData, &decodedEvent)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}
	
	fmt.Println("\nDecoded event:")
	fmt.Printf("  Title: %s\n", decodedEvent.Title)
	fmt.Printf("  Time: %s\n", decodedEvent.Time.Format(time.RFC3339))
}

// Example 5: io.ReadWriteCloser with composition
// ----------------------------------------------
// The io package demonstrates effective interface composition

// BufferedWriteCloser combines a buffered writer with a closer
type BufferedWriteCloser struct {
	*bufio.Writer
	io.Closer
}

// NewBufferedWriteCloser creates a new BufferedWriteCloser
func NewBufferedWriteCloser(w io.WriteCloser, size int) *BufferedWriteCloser {
	return &BufferedWriteCloser{
		Writer: bufio.NewWriterSize(w, size),
		Closer: w,
	}
}

// Close flushes the buffer and closes the underlying writer
func (bwc *BufferedWriteCloser) Close() error {
	if err := bwc.Writer.Flush(); err != nil {
		return err
	}
	return bwc.Closer.Close()
}

// IOCompositionExample demonstrates interface composition
func IOCompositionExample() {
	fmt.Println("\n=== io Interface Composition Example ===")
	
	// Create a file-like buffer
	var buf bytes.Buffer
	
	// Wrap it in a struct that implements io.WriteCloser
	writeCloser := struct {
		io.Writer
		io.Closer
	}{
		Writer: &buf,
		Closer: io.NopCloser(nil), // Dummy closer that does nothing
	}
	
	// Create our buffered write closer
	buffered := NewBufferedWriteCloser(writeCloser, 4096)
	
	// Write some data
	fmt.Fprintln(buffered, "This is buffered until Close() or Flush() is called")
	
	// Data should not be in the underlying buffer yet (because it's buffered)
	fmt.Printf("Before Close: Buffer has %d bytes\n", buf.Len())
	
	// Close will flush and close
	buffered.Close()
	
	// Now the data should be in the buffer
	fmt.Printf("After Close: Buffer has %d bytes\n", buf.Len())
	fmt.Printf("Buffer contents: %q\n", buf.String())
}

// DemonstrateStdlibInterfaces shows all the standard library interface examples
func DemonstrateStdlibInterfaces() {
	fmt.Println("=== Standard Library Interface Patterns ===")
	
	ReaderWriterExample()
	SortInterfaceExample()
	HttpHandlerExample()
	JSONInterfaceExample()
	IOCompositionExample()
}
