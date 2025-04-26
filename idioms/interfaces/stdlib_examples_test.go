package interfaces

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"strings"
	"testing"
	"time"
)

// TestReaderWriterExample validates the io.Reader and io.Writer example
func TestReaderWriterExample(t *testing.T) {
	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the example
	ReaderWriterExample()

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Basic verification that output contains expected elements
	output := buf.String()
	
	// Simply verify that the example ran without errors
	if !strings.Contains(output, "=== io.Reader and io.Writer Examples ===") {
		t.Error("Expected output to contain io.Reader and io.Writer examples header")
	}
	
	if !strings.Contains(output, "Decompressed content:") {
		t.Error("Expected output to contain decompressed content")
	}
	
	// Verify content was properly compressed and decompressed
	if !strings.Contains(output, "Hello, Gophers!") {
		t.Error("Expected decompressed content to contain original text")
	}
}

// TestSortInterfaceExample validates the sort.Interface example
func TestSortInterfaceExample(t *testing.T) {
	// Test sorting by age
	people := []Person{
		{"Alice", 25},
		{"Bob", 30},
		{"Charlie", 20},
		{"David", 35},
	}
	
	// Sort by age
	sort.Sort(ByAge(people))
	
	// Verify order
	if people[0].Name != "Charlie" || people[0].Age != 20 {
		t.Errorf("Expected Charlie (20) to be first, got %s (%d)", people[0].Name, people[0].Age)
	}
	
	if people[3].Name != "David" || people[3].Age != 35 {
		t.Errorf("Expected David (35) to be last, got %s (%d)", people[3].Name, people[3].Age)
	}
	
	// Reset people array and test sorting by name
	people = []Person{
		{"Alice", 25},
		{"Bob", 30},
		{"Charlie", 20},
		{"David", 35},
	}
	
	// Sort by name
	sort.Sort(ByName(people))
	
	// Verify order
	if people[0].Name != "Alice" {
		t.Errorf("Expected Alice to be first, got %s", people[0].Name)
	}
	
	if people[3].Name != "David" {
		t.Errorf("Expected David to be last, got %s", people[3].Name)
	}
	
	// Test sort.Slice
	people = []Person{
		{"Alice", 25},
		{"Bob", 30},
		{"Charlie", 20},
		{"David", 35},
	}
	
	// Sort by age descending
	sort.Slice(people, func(i, j int) bool {
		return people[i].Age > people[j].Age
	})
	
	// Verify descending order
	if people[0].Name != "David" || people[0].Age != 35 {
		t.Errorf("Expected David (35) to be first, got %s (%d)", people[0].Name, people[0].Age)
	}
	
	if people[3].Name != "Charlie" || people[3].Age != 20 {
		t.Errorf("Expected Charlie (20) to be last, got %s (%d)", people[3].Name, people[3].Age)
	}
}

// TestHttpHandlerExample validates the http.Handler example
func TestHttpHandlerExample(t *testing.T) {
	// Create a test HTTP server
	hello := &HelloHandler{Greeting: "Hello"}
	
	// Create a test request
	req, err := http.NewRequest("GET", "/hello/world", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call our handler directly
	hello.ServeHTTP(rr, req)
	
	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	
	// Check the response body
	expected := "Hello, hello/world!\n"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
	
	// Test the middleware
	mw := LoggingMiddleware(hello)
	
	// Capture stdout to test logging
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	
	// Call the middleware
	mw.ServeHTTP(rr, req)
	
	// Restore stdout
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)
	
	// Verify logging happened
	if !strings.Contains(buf.String(), "GET /hello/world") {
		t.Error("Expected logging middleware to log the request")
	}
}

// TestJSONInterfaceExample validates the json.Marshaler/Unmarshaler example
func TestJSONInterfaceExample(t *testing.T) {
	// Create an event with our custom time
	testTime := time.Date(2025, 4, 26, 14, 30, 0, 0, time.UTC)
	event := Event{
		Title: "Team Meeting",
		Time:  CustomTime{Time: testTime},
	}
	
	// Marshal to JSON
	jsonData, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Error marshaling JSON: %v", err)
	}
	
	// Verify correct format
	expected := `{"title":"Team Meeting","time":"2025-04-26 14:30:00"}`
	if string(jsonData) != expected {
		t.Errorf("JSON marshaling produced unexpected result: got %v want %v", string(jsonData), expected)
	}
	
	// Unmarshal back to a struct
	var decodedEvent Event
	err = json.Unmarshal(jsonData, &decodedEvent)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON: %v", err)
	}
	
	// Verify correct date was parsed
	if !decodedEvent.Time.Equal(testTime) {
		t.Errorf("Time unmarshaled incorrectly: got %v want %v", 
			decodedEvent.Time.Format(time.RFC3339), 
			testTime.Format(time.RFC3339))
	}
}

// TestIOCompositionExample validates the interface composition with io interfaces
func TestIOCompositionExample(t *testing.T) {
	// Create a buffer to write to
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
	buffered := NewBufferedWriteCloser(writeCloser, 16) // Small buffer for testing
	
	// Write some test data
	testData := "Testing buffered writer"
	fmt.Fprintln(buffered, testData)
	
	// Data should not be in the underlying buffer yet (because it's buffered)
	if buf.Len() > 0 {
		t.Errorf("Expected buffer to be empty before Close/Flush, but it has %d bytes", buf.Len())
	}
	
	// Close should flush and close
	err := buffered.Close()
	if err != nil {
		t.Errorf("Error closing buffered writer: %v", err)
	}
	
	// Now the data should be in the buffer
	if buf.Len() == 0 {
		t.Error("Expected buffer to contain data after Close, but it's empty")
	}
	
	// Check content
	if !strings.Contains(buf.String(), testData) {
		t.Errorf("Buffer doesn't contain expected data: got %q, want %q", buf.String(), testData)
	}
}

// TestDemonstrateStdlibInterfaces ensures the demonstration function runs without errors
func TestDemonstrateStdlibInterfaces(t *testing.T) {
	// Capture output to avoid cluttering test output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	
	// This should run without panicking
	DemonstrateStdlibInterfaces()
	
	// Restore stdout
	w.Close()
	os.Stdout = oldStdout
	
	// Read output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	
	// Basic check that something was output
	if len(buf.String()) == 0 {
		t.Error("No output from DemonstrateStdlibInterfaces")
	}
}
