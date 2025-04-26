// Package interfaces demonstrates idiomatic Go interface implementation patterns
package interfaces

import (
	"context"
	"fmt"
	"io"
	"time"
)

// Interface Upgrade Patterns
// -------------------------
// This file demonstrates patterns for evolving interfaces over time while
// maintaining backward compatibility.

// V1 of an interface - simple operations
// ------------------------------------

// StorageV1 is the first version of our storage interface
type StorageV1 interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Delete(key string) error
}

// SimpleStorage implements the V1 interface
type SimpleStorage struct {
	data map[string][]byte
}

// NewSimpleStorage creates a new SimpleStorage
func NewSimpleStorage() *SimpleStorage {
	return &SimpleStorage{
		data: make(map[string][]byte),
	}
}

// Get retrieves a value from storage
func (s *SimpleStorage) Get(key string) ([]byte, error) {
	value, ok := s.data[key]
	if !ok {
		return nil, fmt.Errorf("key not found: %s", key)
	}
	return value, nil
}

// Set stores a value in storage
func (s *SimpleStorage) Set(key string, value []byte) error {
	s.data[key] = value
	return nil
}

// Delete removes a value from storage
func (s *SimpleStorage) Delete(key string) error {
	if _, ok := s.data[key]; !ok {
		return fmt.Errorf("key not found: %s", key)
	}
	delete(s.data, key)
	return nil
}

// ClientV1 uses the V1 interface
type ClientV1 struct {
	storage StorageV1
}

// NewClientV1 creates a new ClientV1
func NewClientV1(storage StorageV1) *ClientV1 {
	return &ClientV1{storage: storage}
}

// StoreData stores data using V1 interface
func (c *ClientV1) StoreData(key string, value []byte) error {
	return c.storage.Set(key, value)
}

// GetData retrieves data using V1 interface
func (c *ClientV1) GetData(key string) ([]byte, error) {
	return c.storage.Get(key)
}

// V2 of the interface - adding new methods
// --------------------------------------

// StorageV2 extends StorageV1 with new capabilities
type StorageV2 interface {
	StorageV1                             // Embed V1 interface for backward compatibility
	GetWithExpiration(key string) ([]byte, time.Time, error)
	SetWithExpiration(key string, value []byte, expiration time.Duration) error
	List(prefix string) ([]string, error)
}

// EnhancedStorage implements the V2 interface
type EnhancedStorage struct {
	SimpleStorage                    // Embed V1 implementation
	expiration   map[string]time.Time // Track expiration times
}

// NewEnhancedStorage creates a new EnhancedStorage
func NewEnhancedStorage() *EnhancedStorage {
	return &EnhancedStorage{
		SimpleStorage: *NewSimpleStorage(),
		expiration:    make(map[string]time.Time),
	}
}

// GetWithExpiration retrieves a value and its expiration time
func (s *EnhancedStorage) GetWithExpiration(key string) ([]byte, time.Time, error) {
	value, err := s.Get(key)
	if err != nil {
		return nil, time.Time{}, err
	}
	
	exp, ok := s.expiration[key]
	if !ok {
		exp = time.Time{} // No expiration
	} else if !exp.IsZero() && time.Now().After(exp) {
		// Expired key
		s.Delete(key)
		return nil, time.Time{}, fmt.Errorf("key expired: %s", key)
	}
	
	return value, exp, nil
}

// SetWithExpiration stores a value with an expiration time
func (s *EnhancedStorage) SetWithExpiration(key string, value []byte, expiration time.Duration) error {
	err := s.Set(key, value)
	if err != nil {
		return err
	}
	
	if expiration > 0 {
		s.expiration[key] = time.Now().Add(expiration)
	} else {
		s.expiration[key] = time.Time{} // No expiration
	}
	
	return nil
}

// List returns all keys with the given prefix
func (s *EnhancedStorage) List(prefix string) ([]string, error) {
	var keys []string
	for key := range s.data {
		// Check if key has the prefix
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			keys = append(keys, key)
		}
	}
	return keys, nil
}

// ClientV2 uses the V2 interface
type ClientV2 struct {
	storage StorageV2
}

// NewClientV2 creates a new ClientV2
func NewClientV2(storage StorageV2) *ClientV2 {
	return &ClientV2{storage: storage}
}

// StoreData stores data with expiration
func (c *ClientV2) StoreData(key string, value []byte, expiration time.Duration) error {
	return c.storage.SetWithExpiration(key, value, expiration)
}

// GetData retrieves data checking expiration
func (c *ClientV2) GetData(key string) ([]byte, time.Time, error) {
	return c.storage.GetWithExpiration(key)
}

// ListData lists keys with a prefix
func (c *ClientV2) ListData(prefix string) ([]string, error) {
	return c.storage.List(prefix)
}

// V3 of the interface - adding context support
// -----------------------------------------

// StorageV3 adds context support to our storage interface
type StorageV3 interface {
	// Keep original methods
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Delete(key string) error
	GetWithExpiration(key string) ([]byte, time.Time, error)
	SetWithExpiration(key string, value []byte, expiration time.Duration) error
	List(prefix string) ([]string, error)
	
	// Context-aware methods
	GetWithContext(ctx context.Context, key string) ([]byte, error)
	SetWithContext(ctx context.Context, key string, value []byte) error
	DeleteWithContext(ctx context.Context, key string) error
}

// AdvancedStorage implements the V3 interface
type AdvancedStorage struct {
	*EnhancedStorage // Embed V2 implementation
}

// NewAdvancedStorage creates a new AdvancedStorage
func NewAdvancedStorage() *AdvancedStorage {
	return &AdvancedStorage{
		EnhancedStorage: NewEnhancedStorage(),
	}
}

// GetWithContext gets a value with context support
func (s *AdvancedStorage) GetWithContext(ctx context.Context, key string) ([]byte, error) {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		// Proceed with the operation
		return s.Get(key)
	}
}

// SetWithContext sets a value with context support
func (s *AdvancedStorage) SetWithContext(ctx context.Context, key string, value []byte) error {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Proceed with the operation
		return s.Set(key, value)
	}
}

// DeleteWithContext deletes a value with context support
func (s *AdvancedStorage) DeleteWithContext(ctx context.Context, key string) error {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Proceed with the operation
		return s.Delete(key)
	}
}

// ClientV3 uses the V3 interface
type ClientV3 struct {
	storage StorageV3
}

// NewClientV3 creates a new ClientV3
func NewClientV3(storage StorageV3) *ClientV3 {
	return &ClientV3{storage: storage}
}

// StoreData stores data with context
func (c *ClientV3) StoreData(ctx context.Context, key string, value []byte) error {
	return c.storage.SetWithContext(ctx, key, value)
}

// GetData retrieves data with context
func (c *ClientV3) GetData(ctx context.Context, key string) ([]byte, error) {
	return c.storage.GetWithContext(ctx, key)
}

// Alternative V3 - using adapter pattern
// ------------------------------------

// StorageV3Alt is a simpler interface that only defines new context methods
type StorageV3Alt interface {
	GetWithContext(ctx context.Context, key string) ([]byte, error)
	SetWithContext(ctx context.Context, key string, value []byte) error
	DeleteWithContext(ctx context.Context, key string) error
}

// StorageAdapter adapts V2 to V3Alt
type StorageAdapter struct {
	v2 StorageV2
}

// NewStorageAdapter creates a new adapter
func NewStorageAdapter(v2 StorageV2) *StorageAdapter {
	return &StorageAdapter{v2: v2}
}

// GetWithContext implements StorageV3Alt
func (a *StorageAdapter) GetWithContext(ctx context.Context, key string) ([]byte, error) {
	// Check context first
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		// Delegate to wrapped implementation
		return a.v2.Get(key)
	}
}

// SetWithContext implements StorageV3Alt
func (a *StorageAdapter) SetWithContext(ctx context.Context, key string, value []byte) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return a.v2.Set(key, value)
	}
}

// DeleteWithContext implements StorageV3Alt
func (a *StorageAdapter) DeleteWithContext(ctx context.Context, key string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return a.v2.Delete(key)
	}
}

// Wrapper to combine both V2 and V3Alt capabilities
type CombinedStorage struct {
	StorageV2
	StorageV3Alt
}

// Upgrading while using deprecated functions
// ---------------------------------------

// Original interface
type LoggerV1 interface {
	Log(message string)
	LogError(err error)
}

// Updated interface
type LoggerV2 interface {
	LogWithLevel(level string, message string)
	LogErrorWithLevel(level string, err error)
}

// Implementation of both interfaces
type CombinedLogger struct {
	output io.Writer
}

// NewCombinedLogger creates a new CombinedLogger
func NewCombinedLogger(output io.Writer) *CombinedLogger {
	return &CombinedLogger{output: output}
}

// Log implements LoggerV1.Log (deprecated)
func (l *CombinedLogger) Log(message string) {
	// Delegate to the new method
	l.LogWithLevel("INFO", message)
}

// LogError implements LoggerV1.LogError (deprecated)
func (l *CombinedLogger) LogError(err error) {
	// Delegate to the new method
	l.LogErrorWithLevel("ERROR", err)
}

// LogWithLevel implements LoggerV2.LogWithLevel
func (l *CombinedLogger) LogWithLevel(level string, message string) {
	fmt.Fprintf(l.output, "[%s] %s: %s\n", time.Now().Format(time.RFC3339), level, message)
}

// LogErrorWithLevel implements LoggerV2.LogErrorWithLevel
func (l *CombinedLogger) LogErrorWithLevel(level string, err error) {
	fmt.Fprintf(l.output, "[%s] %s: %s\n", time.Now().Format(time.RFC3339), level, err.Error())
}

// Demo functions for showing interface upgrades
// ------------------------------------------

// ShowV1 demonstrates the V1 interface
func ShowV1() {
	fmt.Println("=== V1 Interface Demo ===")
	storage := NewSimpleStorage()
	client := NewClientV1(storage)
	
	// Store some data
	err := client.StoreData("key1", []byte("value1"))
	if err != nil {
		fmt.Println("Error storing data:", err)
		return
	}
	
	// Retrieve data
	data, err := client.GetData("key1")
	if err != nil {
		fmt.Println("Error getting data:", err)
		return
	}
	
	fmt.Printf("Retrieved data: %s\n", data)
}

// ShowV2 demonstrates the V2 interface
func ShowV2() {
	fmt.Println("\n=== V2 Interface Demo ===")
	storage := NewEnhancedStorage()
	client := NewClientV2(storage)
	
	// Store some data with expiration
	err := client.StoreData("key1", []byte("value1"), 2*time.Second)
	if err != nil {
		fmt.Println("Error storing data:", err)
		return
	}
	
	// Store another key with no expiration
	err = client.StoreData("key2", []byte("value2"), 0)
	if err != nil {
		fmt.Println("Error storing data:", err)
		return
	}
	
	// Retrieve data with expiration info
	data, exp, err := client.GetData("key1")
	if err != nil {
		fmt.Println("Error getting data:", err)
		return
	}
	
	fmt.Printf("Retrieved data: %s, Expires: %v\n", data, exp)
	
	// List keys
	keys, err := client.ListData("key")
	if err != nil {
		fmt.Println("Error listing data:", err)
		return
	}
	
	fmt.Printf("Keys: %v\n", keys)
	
	// Wait for expiration
	fmt.Println("Waiting for key1 to expire...")
	time.Sleep(3 * time.Second)
	
	// Try to get expired key
	data, exp, err = client.GetData("key1")
	if err != nil {
		fmt.Printf("After waiting: %v\n", err)
	} else {
		fmt.Printf("After waiting - Retrieved data: %s, Expires: %v\n", data, exp)
	}
	
	// key2 should still be available
	data, exp, err = client.GetData("key2")
	if err != nil {
		fmt.Printf("Error getting key2: %v\n", err)
	} else {
		fmt.Printf("Key2 data: %s, Expires: %v\n", data, exp)
	}
}

// ShowV3 demonstrates the V3 interface with context
func ShowV3() {
	fmt.Println("\n=== V3 Interface with Context Demo ===")
	storage := NewAdvancedStorage()
	client := NewClientV3(storage)
	
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	
	// Store and retrieve data with context
	err := client.StoreData(ctx, "key3", []byte("value3"))
	if err != nil {
		fmt.Println("Error storing data with context:", err)
		return
	}
	
	data, err := client.GetData(ctx, "key3")
	if err != nil {
		fmt.Println("Error getting data with context:", err)
		return
	}
	
	fmt.Printf("Retrieved data with context: %s\n", data)
	
	// Demonstrate context timeout
	fmt.Println("Demonstrating context timeout...")
	
	// Create expired context
	expiredCtx, expiredCancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer expiredCancel()
	
	// Sleep to ensure context expires
	time.Sleep(5 * time.Millisecond)
	
	// Try to use expired context
	_, err = client.GetData(expiredCtx, "key3")
	fmt.Printf("Using expired context: %v\n", err)
}

// ShowAdapterPattern demonstrates the adapter pattern for interface upgrades
func ShowAdapterPattern() {
	fmt.Println("\n=== Adapter Pattern Demo ===")
	
	// Create a V2 implementation
	v2Storage := NewEnhancedStorage()
	
	// Create an adapter to make it compatible with V3Alt
	adapter := NewStorageAdapter(v2Storage)
	
	// Now we can use the adapter with V3Alt context methods
	ctx := context.Background()
	err := adapter.SetWithContext(ctx, "adapted-key", []byte("adapted-value"))
	if err != nil {
		fmt.Println("Error setting data through adapter:", err)
		return
	}
	
	data, err := adapter.GetWithContext(ctx, "adapted-key")
	if err != nil {
		fmt.Println("Error getting data through adapter:", err)
		return
	}
	
	fmt.Printf("Retrieved data through adapter: %s\n", data)
	
	// We can also still use the V2 methods directly
	v2Storage.Set("direct-key", []byte("direct-value"))
	directData, _ := v2Storage.Get("direct-key")
	fmt.Printf("Retrieved direct data: %s\n", directData)
}

// ShowDeprecatedFunction demonstrates upgrading with deprecated functions
func ShowDeprecatedFunction() {
	fmt.Println("\n=== Deprecated Function Demo ===")
	
	// Create a logger that implements both V1 and V2 interfaces
	logger := NewCombinedLogger(io.Discard)
	
	// Old code still works
	logger.Log("This uses the old interface")
	logger.LogError(fmt.Errorf("old interface error"))
	
	// New code uses the new interface
	logger.LogWithLevel("DEBUG", "This uses the new interface")
	logger.LogErrorWithLevel("WARN", fmt.Errorf("new interface error"))
	
	fmt.Println("Both old and new interface methods work")
	fmt.Println("Log messages would be written to the output writer")
}

// InterfaceUpgradeDemo demonstrates interface upgrade patterns
func InterfaceUpgradeDemo() {
	fmt.Println("============================================")
	fmt.Println("Interface Upgrade Patterns Demo")
	fmt.Println("============================================")
	
	// Show the evolution of our interfaces
	ShowV1()
	ShowV2()
	ShowV3()
	ShowAdapterPattern()
	ShowDeprecatedFunction()
	
	fmt.Println("\nInterface Upgrade Best Practices:")
	fmt.Println("1. Keep backward compatibility when possible")
	fmt.Println("2. Use embedding to extend interfaces")
	fmt.Println("3. Consider providing adapters for new capabilities")
	fmt.Println("4. Delegate deprecated methods to newer implementations")
	fmt.Println("5. Use versioning in interface names when breaking changes are necessary")
	fmt.Println("6. Document upgrade paths and migration strategies")
	fmt.Println("============================================")
}
