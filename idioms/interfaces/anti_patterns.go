// Package interfaces demonstrates idiomatic Go interface implementation patterns.
package interfaces

import (
	"fmt"
	"io"
	"sync"
)

// AntiPatterns demonstrates common interface mistakes to avoid in Go.
// This file highlights patterns that should generally be avoided when
// working with interfaces in Go.

// AntiPattern1: Fat Interfaces
// --------------------------
// Problem: Creating large interfaces with many methods makes them less reusable
// and harder to implement.

// BadDatabaseClient is an example of a "fat interface" - it has too many methods
// and combines multiple responsibilities.
type BadDatabaseClient interface {
	Connect(connectionString string) error
	Disconnect() error
	Query(query string) ([]byte, error)
	Execute(command string) error
	BeginTransaction() error
	CommitTransaction() error
	RollbackTransaction() error
	CacheResults(key string, data []byte) error
	GetFromCache(key string) ([]byte, error)
	ValidateSchema() error
	MigrateDatabase() error
	LogQuery(query string) error
	Authenticate(username, password string) error
	CheckPermissions(user string, resource string) bool
	// ... many more methods
}

// Better approach: Use smaller, focused interfaces with specific responsibilities
type Queryable interface {
	Query(query string) ([]byte, error)
}

type Transactional interface {
	BeginTransaction() error
	CommitTransaction() error
	RollbackTransaction() error
}

type Connectable interface {
	Connect(connectionString string) error
	Disconnect() error
}

// ------------------------------------------------

// AntiPattern2: Returning Interfaces
// --------------------------
// Problem: Returning interfaces instead of concrete types limits the caller
// and complicates the API unnecessarily.

// BadFactory returns an interface, forcing the caller to work with a limited view
func BadFactory() io.Reader {
	return &ConcreteReader{}
}

// ConcreteReader is an implementation of io.Reader
type ConcreteReader struct {
	data []byte
	pos  int
}

// Read implements the io.Reader interface
func (r *ConcreteReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

// AdditionalMethod is a useful method that's hidden when returning just an io.Reader
func (r *ConcreteReader) AdditionalMethod() string {
	return "This method is not part of io.Reader interface"
}

// BetterFactory returns a concrete type, giving callers access to all methods
func BetterFactory() *ConcreteReader {
	return &ConcreteReader{data: []byte("Hello, World!")}
}

// ------------------------------------------------

// AntiPattern3: Interface Pollution
// --------------------------
// Problem: Creating interfaces for every type unnecessarily, 
// or creating interfaces before you need them.

// UnnecessaryRepositoryInterface is an example of an interface created unnecessarily
// before there are multiple implementations or a need for abstraction
type UnnecessaryRepositoryInterface interface {
	GetUserByID(id int) (User, error)
	CreateUser(user User) error
	UpdateUser(user User) error
	DeleteUser(id int) error
}

// UserRepository is the only implementation of this interface,
// making the interface unnecessary until there's a genuine need for abstraction
type UserRepository struct {
	// implementation details
}

// GetUserByID gets a user from storage
func (r *UserRepository) GetUserByID(id int) (User, error) {
	// Implementation...
	return User{ID: id, Name: "Example"}, nil
}

// CreateUser adds a new user to storage
func (r *UserRepository) CreateUser(user User) error {
	// Implementation...
	return nil
}

// UpdateUser updates an existing user
func (r *UserRepository) UpdateUser(user User) error {
	// Implementation...
	return nil
}

// DeleteUser removes a user from storage
func (r *UserRepository) DeleteUser(id int) error {
	// Implementation...
	return nil
}

// Better approach: Start with concrete types and extract interfaces only when needed,
// preferably from the consumer's perspective.

// ------------------------------------------------

// AntiPattern4: Pointer vs. Value Methods Confusion
// --------------------------
// Problem: Inconsistent use of pointer receivers and value receivers 
// for methods that implement an interface

// Inconsistent is an example interface
type Inconsistent interface {
	Method1()
	Method2()
}

// InconsistentImpl mixes pointer and value receivers 
type InconsistentImpl struct {
	data string
}

// Method1 uses a pointer receiver
func (i *InconsistentImpl) Method1() {
	i.data = "modified"
}

// Method2 uses a value receiver
func (i InconsistentImpl) Method2() {
	// Cannot modify i.data here
	fmt.Println(i.data)
}

// Usage demonstrates a common mistake
func InconsistentUsage() {
	// This compiles but won't work as expected
	var impl Inconsistent = InconsistentImpl{"data"}
	
	// Error: InconsistentImpl does not implement Inconsistent
	// Only *InconsistentImpl implements Method1 with a pointer receiver
	
	// Correct usage:
	var correctImpl Inconsistent = &InconsistentImpl{"data"}
	correctImpl.Method1() // Works correctly
}

// ------------------------------------------------

// AntiPattern5: Interface Type Assertions Without Checking
// --------------------------
// Problem: Performing type assertions without checking can cause panics

// UnsafeTypeAssertions demonstrates dangerous type assertions
func UnsafeTypeAssertions(value interface{}) {
	// Dangerous: Will panic if value is not a string
	str := value.(string)
	fmt.Println(str)

	// Better approach: Type assertion with check
	strValue, ok := value.(string)
	if ok {
		fmt.Println(strValue)
	} else {
		fmt.Println("Value is not a string")
	}

	// Or use a type switch for multiple possibilities
	switch v := value.(type) {
	case string:
		fmt.Println("String:", v)
	case int:
		fmt.Println("Integer:", v)
	default:
		fmt.Println("Unknown type")
	}
}

// ------------------------------------------------

// AntiPattern6: Embedding Interfaces in Structs Incorrectly
// --------------------------
// Problem: Embedding an interface in a struct without providing implementations

// BadInterfaceEmbedding embeds the io.Writer interface but doesn't implement its methods
type BadInterfaceEmbedding struct {
	io.Writer // Embedding interface without implementing Write method
	data      []byte
}

// Attempting to use this struct as a Writer will cause runtime errors

// CorrectInterfaceEmbedding properly embeds and implements the interface
type CorrectInterfaceEmbedding struct {
	buffer []byte
}

func (c *CorrectInterfaceEmbedding) Write(p []byte) (n int, err error) {
	c.buffer = append(c.buffer, p...)
	return len(p), nil
}

// ------------------------------------------------

// AntiPattern7: Not Considering Method Sets
// --------------------------
// Problem: Not understanding which methods are in the method set of values vs. pointers

// ValueReceiverOnly has methods with value receivers
type ValueReceiverOnly struct {
	counter int
}

func (v ValueReceiverOnly) Count() int {
	return v.counter
}

// PointerReceiverOnly has methods with pointer receivers
type PointerReceiverOnly struct {
	counter int
}

func (p *PointerReceiverOnly) Increment() {
	p.counter++
}

func (p *PointerReceiverOnly) Count() int {
	return p.counter
}

// Counter interface captures counting behavior
type Counter interface {
	Count() int
}

func UseCounter(c Counter) {
	fmt.Println("Count:", c.Count())
}

func MethodSetExample() {
	// This works fine - both value and pointer satisfy Counter
	value := ValueReceiverOnly{counter: 5}
	UseCounter(value)      // Works
	UseCounter(&value)     // Works too

	// For pointer receiver methods:
	pointer := PointerReceiverOnly{counter: 10}
	UseCounter(&pointer)   // Works
	// UseCounter(pointer) // Error: pointer.Count is not in the method set of PointerReceiverOnly
}

// ------------------------------------------------

// AntiPattern8: Forgetting Concurrency Considerations
// --------------------------
// Problem: Not considering thread safety when implementing interfaces

// UnsafeCounter is not safe for concurrent use
type UnsafeCounter struct {
	count int
}

func (c *UnsafeCounter) Increment() {
	c.count++ // Not atomic or protected
}

func (c *UnsafeCounter) GetCount() int {
	return c.count
}

// SafeCounter is safe for concurrent use
type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *SafeCounter) GetCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

// RunConcurrencyExample demonstrates the difference between safe and unsafe implementations
func RunConcurrencyExample() {
	unsafe := &UnsafeCounter{}
	safe := &SafeCounter{}

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

	// unsafe.count will likely be < 1000 due to race conditions
	// safe.count will be exactly 1000
	fmt.Printf("Unsafe counter: %d\n", unsafe.GetCount())
	fmt.Printf("Safe counter: %d\n", safe.GetCount())
}

// ------------------------------------------------

// AntiPattern9: Empty Interface Overuse
// --------------------------
// Problem: Using empty interfaces (interface{}) too often reduces type safety

// GenericHolder overuses empty interface
type GenericHolder struct {
	data map[string]interface{}
}

func NewGenericHolder() *GenericHolder {
	return &GenericHolder{
		data: make(map[string]interface{}),
	}
}

func (h *GenericHolder) Set(key string, value interface{}) {
	h.data[key] = value
}

func (h *GenericHolder) Get(key string) interface{} {
	return h.data[key]
}

// Usage requires type assertions and can lead to runtime errors
func UsingGenericHolder() {
	holder := NewGenericHolder()
	holder.Set("name", "John")
	holder.Set("age", 30)

	// Type assertion required, potential for panics
	name := holder.Get("name").(string)
	age := holder.Get("age").(int)
	fmt.Printf("Name: %s, Age: %d\n", name, age)

	// This will panic at runtime
	// salary := holder.Get("salary").(float64) // Key doesn't exist
}

// Better approach: Use generics (Go 1.18+) or specific types where possible

// TypedHolder uses generics for better type safety
type TypedHolder[T any] struct {
	data map[string]T
}

func NewTypedHolder[T any]() *TypedHolder[T] {
	return &TypedHolder[T]{
		data: make(map[string]T),
	}
}

func (h *TypedHolder[T]) Set(key string, value T) {
	h.data[key] = value
}

func (h *TypedHolder[T]) Get(key string) (T, bool) {
	value, ok := h.data[key]
	return value, ok
}

// Usage is type-safe and prevents runtime errors
func UsingTypedHolder() {
	stringHolder := NewTypedHolder[string]()
	stringHolder.Set("name", "John")
	
	if name, ok := stringHolder.Get("name"); ok {
		fmt.Println("Name:", name)
	}
	
	// This won't compile:
	// stringHolder.Set("age", 30) // Type error caught at compile time
}

// DemonstrateAntiPatterns shows examples of interface anti-patterns
func DemonstrateAntiPatterns() {
	fmt.Println("=== Interface Anti-Patterns ===")
	
	// Example 1: Interface pollution
	repo := &UserRepository{}
	user, _ := repo.GetUserByID(1)
	fmt.Println("User:", user.Name)
	
	// Example 2: Returning concrete types vs interfaces
	reader := BetterFactory()
	fmt.Println("Additional method:", reader.AdditionalMethod())
	
	// Example 3: Type assertions
	UnsafeTypeAssertions("hello")
	
	// Example 4: Method sets and receivers
	MethodSetExample()
	
	// Example 5: Concurrency considerations
	RunConcurrencyExample()
	
	// Example 6: Generic holder vs typed holder
	UsingGenericHolder()
	UsingTypedHolder()
}
