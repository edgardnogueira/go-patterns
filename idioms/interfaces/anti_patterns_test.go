package interfaces

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"testing"
)

// TestFatInterfaces validates the difference between fat interfaces and smaller focused ones
func TestFatInterfaces(t *testing.T) {
	// This test demonstrates why smaller interfaces are better than large ones
	
	// Implementing the full BadDatabaseClient would be cumbersome
	// But implementing smaller, focused interfaces is much easier
	
	// Simple test implementation of Queryable
	queryClient := struct {
		Queryable
	}{
		Queryable: queryStub{},
	}
	
	result, err := queryClient.Query("SELECT * FROM users")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if string(result) != "Query executed" {
		t.Errorf("Expected 'Query executed', got %s", string(result))
	}
}

// queryStub is a test implementation of Queryable
type queryStub struct{}

func (q queryStub) Query(query string) ([]byte, error) {
	return []byte("Query executed"), nil
}

// TestReturningInterfaces validates the pattern of returning concrete types vs interfaces
func TestReturningInterfaces(t *testing.T) {
	// When we use BadFactory, we're limited to io.Reader methods
	reader := BadFactory()
	buf := make([]byte, 13)
	_, err := reader.Read(buf)
	if err != nil {
		t.Fatal("Error reading from BadFactory result")
	}
	
	// We can't access AdditionalMethod from reader
	// reader.AdditionalMethod() // Would not compile
	
	// When we use BetterFactory, we get access to all methods
	concreteReader := BetterFactory()
	additionalOutput := concreteReader.AdditionalMethod()
	if additionalOutput != "This method is not part of io.Reader interface" {
		t.Errorf("Expected correct string from AdditionalMethod, got %s", additionalOutput)
	}
}

// TestInterfacePollution validates that unnecessary interfaces add complexity
func TestInterfacePollution(t *testing.T) {
	// Direct usage of the concrete type is simpler when there's only one implementation
	repository := &UserRepository{}
	
	user, err := repository.GetUserByID(1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if user.Name != "Example" {
		t.Errorf("Expected user name 'Example', got %s", user.Name)
	}
	
	// The interface would add an unnecessary layer of abstraction
	// var repo UnnecessaryRepositoryInterface = repository
	// user, _ = repo.GetUserByID(1)
}

// TestPointerVsValueMethods validates the pitfalls of mixing receiver types
func TestPointerVsValueMethods(t *testing.T) {
	// Test that pointer receivers work as expected
	ptrImpl := &InconsistentImpl{data: "original"}
	ptrImpl.Method1()
	if ptrImpl.data != "modified" {
		t.Errorf("Expected data to be modified, but it's still %s", ptrImpl.data)
	}
	
	// Test to show value receivers don't modify the struct
	valueImpl := InconsistentImpl{data: "original"}
	valueImpl.Method1() // This modifies a copy, not the original
	if valueImpl.data == "modified" {
		t.Errorf("Expected data to remain 'original', but got %s", valueImpl.data)
	}
	
	// This would cause a compile error:
	// var impl Inconsistent = InconsistentImpl{}
	
	// This works:
	var impl Inconsistent = &InconsistentImpl{}
	impl.Method1()
}

// TestTypeAssertions validates safe vs unsafe type assertions
func TestTypeAssertions(t *testing.T) {
	// Test safe type assertion
	var value interface{} = "test string"
	
	// Safe assertion with check
	strValue, ok := value.(string)
	if !ok {
		t.Error("Type assertion should succeed for string value")
	}
	if strValue != "test string" {
		t.Errorf("Expected 'test string', got %s", strValue)
	}
	
	// Test type switch
	var result string
	switch v := value.(type) {
	case string:
		result = "string:" + v
	case int:
		result = fmt.Sprintf("int:%d", v)
	default:
		result = "unknown"
	}
	
	if result != "string:test string" {
		t.Errorf("Expected 'string:test string', got %s", result)
	}
}

// TestInterfaceEmbedding validates correct vs incorrect interface embedding
func TestInterfaceEmbedding(t *testing.T) {
	// Test the correct embedding implementation
	correct := &CorrectInterfaceEmbedding{}
	
	n, err := correct.Write([]byte("test data"))
	if err != nil {
		t.Fatalf("Write should not error: %v", err)
	}
	if n != 9 {
		t.Errorf("Expected to write 9 bytes, wrote %d", n)
	}
	
	// Bad embedding would cause a nil pointer panic at runtime
	// This is commented out because it would panic
	/*
	bad := &BadInterfaceEmbedding{}
	_, _ = bad.Write([]byte("test")) // Would panic at runtime
	*/
}

// TestMethodSets validates method set behavior with pointers and values
func TestMethodSets(t *testing.T) {
	// Value receiver methods are in the method set of both the value and pointer type
	value := ValueReceiverOnly{counter: 5}
	
	// Both of these should work
	if value.Count() != 5 {
		t.Error("Value method call failed")
	}
	
	ptr := &value
	if ptr.Count() != 5 {
		t.Error("Pointer method call failed")
	}
	
	// For pointer receiver methods
	ptrOnly := &PointerReceiverOnly{counter: 3}
	ptrOnly.Increment()
	if ptrOnly.Count() != 4 {
		t.Error("Pointer receiver increment failed")
	}
	
	// This would be a compile error because the method is not in the method set of the value type
	// valueOnly := PointerReceiverOnly{counter: 3}
	// valueOnly.Increment() // Not allowed
}

// TestConcurrencySafety validates thread-safe vs unsafe implementations
func TestConcurrencySafety(t *testing.T) {
	unsafe := &UnsafeCounter{}
	safe := &SafeCounter{}
	
	// Run a lot of concurrent increments
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			unsafe.Increment()
		}()
		go func() {
			defer wg.Done()
			safe.Increment()
		}()
	}
	wg.Wait()
	
	// Safe counter should be exactly 1000
	if safe.GetCount() != 1000 {
		t.Errorf("SafeCounter should be 1000, got %d", safe.GetCount())
	}
	
	// Unsafe counter will likely be less than 1000 due to race conditions
	// We don't assert on its exact value because it's inherently unpredictable
	t.Logf("UnsafeCounter ended up at %d (should be 1000 if there were no race conditions)", unsafe.GetCount())
}

// TestEmptyInterfaceOveruse validates issues with empty interfaces vs type-safe approaches
func TestEmptyInterfaceOveruse(t *testing.T) {
	// Test generic holder with interface{}
	genericHolder := NewGenericHolder()
	genericHolder.Set("name", "John")
	genericHolder.Set("age", 30)
	
	// Type assertions required
	name, ok := genericHolder.Get("name").(string)
	if !ok {
		t.Error("Failed to get name as string")
	}
	if name != "John" {
		t.Errorf("Expected name 'John', got %s", name)
	}
	
	// Safer with generics
	stringHolder := NewTypedHolder[string]()
	stringHolder.Set("firstName", "John")
	firstName, ok := stringHolder.Get("firstName")
	if !ok {
		t.Error("Failed to get firstName")
	}
	if firstName != "John" {
		t.Errorf("Expected firstName 'John', got %s", firstName)
	}
	
	// This would be a compile-time error (type safety benefit):
	// stringHolder.Set("age", 30)
}

// TestBadFactory ensures that type information is preserved
func TestBadFactory(t *testing.T) {
	reader := BadFactory()
	
	// We can only use reader methods
	buf := make([]byte, 100)
	_, err := reader.Read(buf)
	if err != nil && err != io.EOF {
		t.Errorf("Unexpected error: %v", err)
	}
	
	// Cannot access additional methods
	// reader.AdditionalMethod() // Would not compile
	
	// This forces a type assertion if we need the concrete type
	if concreteReader, ok := reader.(*ConcreteReader); ok {
		method := concreteReader.AdditionalMethod()
		if method != "This method is not part of io.Reader interface" {
			t.Errorf("Unexpected result: %s", method)
		}
	} else {
		t.Error("Type assertion failed")
	}
}

// TestDemonstrateAntiPatterns ensures the demonstration function doesn't panic
func TestDemonstrateAntiPatterns(t *testing.T) {
	// Capture output to avoid cluttering test output
	originalStdout := fmt.Stdout
	r, w, _ := os.Pipe()
	fmt.Stdout = w
	
	// This should run without panicking
	DemonstrateAntiPatterns()
	
	// Restore stdout
	w.Close()
	fmt.Stdout = originalStdout
	
	// Read output
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()
	
	// Basic check that something was output
	if len(output) == 0 {
		t.Error("No output from DemonstrateAntiPatterns")
	}
}
