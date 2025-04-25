package proxy

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
)

// TestRealImage tests the basic RealImage implementation
func TestRealImage(t *testing.T) {
	// Create a new real image
	image, err := NewRealImage("test.jpg")
	if err != nil {
		t.Fatalf("Failed to create RealImage: %v", err)
	}

	// Test GetFilename
	if image.GetFilename() != "test.jpg" {
		t.Errorf("Expected filename 'test.jpg', got '%s'", image.GetFilename())
	}

	// Test Display
	err = image.Display()
	if err != nil {
		t.Errorf("Display failed: %v", err)
	}

	// Test GetWidth and GetHeight
	if image.GetWidth() != 1920 || image.GetHeight() != 1080 {
		t.Errorf("Expected dimensions 1920x1080, got %dx%d", image.GetWidth(), image.GetHeight())
	}

	// Test GetSize
	if image.GetSize() != 2048*1024 {
		t.Errorf("Expected size %d, got %d", 2048*1024, image.GetSize())
	}

	// Test GetMetadata
	metadata := image.GetMetadata()
	if metadata["format"] != "JPEG" {
		t.Errorf("Expected format 'JPEG', got '%s'", metadata["format"])
	}

	// Test error handling with non-existent image
	_, err = NewRealImage("not_found.jpg")
	if err == nil {
		t.Error("Expected error for non-existent image, got nil")
	}
}

// TestVirtualProxy tests the virtual proxy lazy loading functionality
func TestVirtualProxy(t *testing.T) {
	// Create a virtual proxy
	proxy := NewVirtualProxy("test.jpg")

	// Check initial state
	if proxy.IsLoaded() {
		t.Error("Expected image to not be loaded initially")
	}

	// Get filename without loading the image
	if proxy.GetFilename() != "test.jpg" {
		t.Errorf("Expected filename 'test.jpg', got '%s'", proxy.GetFilename())
	}

	// Still shouldn't be loaded
	if proxy.IsLoaded() {
		t.Error("Expected image to not be loaded after GetFilename")
	}

	// Now display the image, which should trigger loading
	err := proxy.Display()
	if err != nil {
		t.Errorf("Display failed: %v", err)
	}

	// Should be loaded now
	if !proxy.IsLoaded() {
		t.Error("Expected image to be loaded after Display")
	}

	// Test other methods that should use the loaded image
	if proxy.GetWidth() != 1920 || proxy.GetHeight() != 1080 {
		t.Errorf("Expected dimensions 1920x1080, got %dx%d", proxy.GetWidth(), proxy.GetHeight())
	}

	// Test with non-existent image
	badProxy := NewVirtualProxy("not_found.jpg")
	err = badProxy.Display()
	if err == nil {
		t.Error("Expected error for non-existent image, got nil")
	}
}

// TestProtectionProxy tests the access control functionality
func TestProtectionProxy(t *testing.T) {
	// Create a real image
	realImage, _ := NewRealImage("test.jpg")

	// Create an admin user
	adminUser := &User{
		Username: "admin",
		Role:     "admin",
	}

	// Create a regular user
	regularUser := &User{
		Username: "user",
		Role:     "user",
	}

	// Test with admin user
	adminProxy := NewProtectionProxy(realImage, adminUser)
	err := adminProxy.Display()
	if err != nil {
		t.Errorf("Admin user should be allowed to display image, got error: %v", err)
	}

	// Test with regular user (should be denied)
	userProxy := NewProtectionProxy(realImage, regularUser)
	err = userProxy.Display()
	if err == nil {
		t.Error("Regular user should be denied access, but no error was returned")
	}

	// Test without a user
	noUserProxy := NewProtectionProxy(realImage, nil)
	err = noUserProxy.Display()
	if err == nil {
		t.Error("Proxy without user should be denied access, but no error was returned")
	}

	// Test file type validation
	badFileProxy := NewProtectionProxy(&RealImage{filename: "bad.exe"}, adminUser)
	err = badFileProxy.Display()
	if err == nil {
		t.Error("Expected error for unsupported file type, got nil")
	}

	// Test allowing new roles
	userProxy.SetAllowedRoles([]string{"user", "admin", "guest"})
	err = userProxy.Display()
	if err != nil {
		t.Errorf("User should be allowed after role update, got error: %v", err)
	}
}

// TestLoggingProxy tests the logging functionality
func TestLoggingProxy(t *testing.T) {
	// Create a real image
	realImage, _ := NewRealImage("test.jpg")

	// Create a buffer to capture logs
	var buf bytes.Buffer

	// Create a logging proxy
	proxy := NewLoggingProxy(realImage, INFO)
	proxy.SetWriter(&buf)

	// Test display with logging
	proxy.Display()

	// Check that logs were written
	logs := buf.String()
	if !strings.Contains(logs, "Display method called") {
		t.Errorf("Expected log message about display call, got: %s", logs)
	}

	// Test with different log level
	buf.Reset()
	proxy.SetLogLevel(ERROR)
	proxy.GetWidth() // This should not log at ERROR level

	// Check that debug logs were not written
	logs = buf.String()
	if logs != "" {
		t.Errorf("Expected no logs at ERROR level for GetWidth, got: %s", logs)
	}

	// Test with custom prefix
	buf.Reset()
	proxy.SetLogLevel(INFO)
	proxy.SetPrefix("TestLogger")
	proxy.Display()

	// Check that the prefix was used
	logs = buf.String()
	if !strings.Contains(logs, "TestLogger") {
		t.Errorf("Expected log message with custom prefix, got: %s", logs)
	}
}

// TestMetricsProxy tests the metrics collection functionality
func TestMetricsProxy(t *testing.T) {
	// Create a real image
	realImage, _ := NewRealImage("test.jpg")

	// Create a metrics proxy
	proxy := NewMetricsProxy(realImage)

	// Call some methods to generate metrics
	proxy.Display()
	proxy.GetWidth()
	proxy.GetHeight()
	proxy.Display() // Call again to check counter incrementing

	// Check metrics
	metrics := proxy.GetMetrics()

	// Check Display metrics
	displayMetrics, ok := metrics["Display"]
	if !ok {
		t.Fatal("Display metrics not found")
	}
	if displayMetrics[MetricRequestCount] != 2 {
		t.Errorf("Expected 2 Display requests, got %.0f", displayMetrics[MetricRequestCount])
	}

	// Check other operation metrics exist
	if _, ok := metrics["GetWidth"]; !ok {
		t.Error("GetWidth metrics not found")
	}
	if _, ok := metrics["GetHeight"]; !ok {
		t.Error("GetHeight metrics not found")
	}

	// Test reset
	proxy.ResetMetrics()
	metrics = proxy.GetMetrics()
	if len(metrics) != 0 {
		t.Errorf("Expected empty metrics after reset, got %d entries", len(metrics))
	}

	// Test metrics hook
	var hookCalled bool
	proxy.SetMetricsHook(func(op string, metricType MetricType, val float64) {
		hookCalled = true
	})
	proxy.Display()
	if !hookCalled {
		t.Error("Metrics hook was not called")
	}
}

// TestCachingProxy tests the caching functionality
func TestCachingProxy(t *testing.T) {
	// Create a caching proxy with a 1 second expiration
	proxy := NewCachingProxy(1*time.Second, 2)

	// Access an image
	err := proxy.Display("test.jpg")
	if err != nil {
		t.Errorf("Display failed: %v", err)
	}

	// Access again - should be from cache
	err = proxy.Display("test.jpg")
	if err != nil {
		t.Errorf("Display from cache failed: %v", err)
	}

	// Access a different image
	err = proxy.Display("other.jpg")
	if err != nil {
		t.Errorf("Display of other image failed: %v", err)
	}

	// Get cache stats
	stats := proxy.GetCacheStats()
	if stats["size"].(int) != 2 {
		t.Errorf("Expected 2 items in cache, got %d", stats["size"].(int))
	}

	// Access a third image - should trigger eviction
	err = proxy.Display("third.jpg")
	if err != nil {
		t.Errorf("Display of third image failed: %v", err)
	}

	// Check cache stats again
	stats = proxy.GetCacheStats()
	if stats["size"].(int) != 2 {
		t.Errorf("Expected 2 items in cache after eviction, got %d", stats["size"].(int))
	}

	// Wait for cache to expire
	time.Sleep(1100 * time.Millisecond)

	// Access again - should reload
	err = proxy.Display("test.jpg")
	if err != nil {
		t.Errorf("Display after expiration failed: %v", err)
	}

	// Test clear cache
	proxy.ClearCache()
	stats = proxy.GetCacheStats()
	if stats["size"].(int) != 0 {
		t.Errorf("Expected 0 items after cache clear, got %d", stats["size"].(int))
	}
}

// TestRemoteProxy tests the remote resource handling
func TestRemoteProxy(t *testing.T) {
	// Create a remote proxy
	proxy := NewRemoteProxy("http://example.com/api", "test.jpg")

	// Set a shorter cache duration for testing
	proxy.SetCacheDuration(500 * time.Millisecond)

	// Initially, data should not be cached
	if proxy.IsDataCached() {
		t.Error("Data should not be cached initially")
	}

	// Access the image, which should fetch remote data
	err := proxy.Display()
	if err != nil {
		t.Errorf("Display failed: %v", err)
	}

	// Now data should be cached
	if !proxy.IsDataCached() {
		t.Error("Data should be cached after access")
	}

	// Check cache status
	status := proxy.GetCacheStatus()
	if !status["has_cached_data"].(bool) {
		t.Error("Cache status should indicate cached data")
	}
	if !status["is_fresh"].(bool) {
		t.Error("Cache status should indicate fresh data")
	}

	// Access again - should use cache
	err = proxy.Display()
	if err != nil {
		t.Errorf("Display from cache failed: %v", err)
	}

	// Clear cache
	proxy.ClearCache()
	if proxy.IsDataCached() {
		t.Error("Data should not be cached after clear")
	}

	// Wait for cache to expire
	time.Sleep(600 * time.Millisecond)
	if proxy.IsDataCached() {
		t.Error("Data should not be considered cached after expiration")
	}

	// Test with non-existent image
	badProxy := NewRemoteProxy("http://example.com/api", "not_found.jpg")
	err = badProxy.Display()
	if err == nil {
		t.Error("Expected error for non-existent remote image, got nil")
	}
}

// TestProxyChain tests the proxy chaining capabilities
func TestProxyChain(t *testing.T) {
	// Create a real image
	realImage, _ := NewRealImage("test.jpg")

	// Create a user
	adminUser := &User{
		Username: "admin",
		Role:     "admin",
	}

	// Create a chain with multiple proxies
	chain := NewProxyChain(realImage).
		AddLogging(INFO).
		AddMetrics().
		AddProtection(adminUser).
		Build()

	// Access through the chain
	err := chain.Display()
	if err != nil {
		t.Errorf("Chain display failed: %v", err)
	}

	// Test caching chain
	cachingChain := NewProxyChain(realImage).
		AddCaching(1*time.Second, 10)

	// Access through the caching chain
	err = cachingChain.Display()
	if err != nil {
		t.Errorf("Caching chain display failed: %v", err)
	}

	// Test preset chains
	presets := ProxyPresets{}
	secureChain := presets.NewSecure(realImage, adminUser)

	err = secureChain.Display()
	if err != nil {
		t.Errorf("Secure preset chain display failed: %v", err)
	}
}

// TestIntegration tests a realistic scenario with multiple proxy types
func TestIntegration(t *testing.T) {
	// Create test scenario with multiple proxies combined
	// 1. Start with a virtual proxy for lazy loading
	virtualProxy := NewVirtualProxy("important_image.jpg")

	// 2. Wrap with a logging proxy
	var logBuf bytes.Buffer
	loggingProxy := NewLoggingProxy(virtualProxy, INFO)
	loggingProxy.SetWriter(&logBuf)

	// 3. Add protection
	adminUser := &User{Username: "admin", Role: "admin"}
	protectionProxy := NewProtectionProxy(loggingProxy, adminUser)

	// 4. Add metrics collection
	metricsProxy := NewMetricsProxy(protectionProxy)

	// Execute operations through the proxy chain
	fmt.Println("=== Integration Test ===")
	
	// First display should trigger lazy loading
	err := metricsProxy.Display()
	if err != nil {
		t.Errorf("Integration display failed: %v", err)
	}

	// Check that virtual proxy is now loaded
	if !virtualProxy.IsLoaded() {
		t.Error("Virtual proxy should be loaded after access through chain")
	}

	// Check that logs were written
	logs := logBuf.String()
	if !strings.Contains(logs, "Display method called") {
		t.Errorf("Logging proxy should have logged the display call")
	}

	// Check metrics were collected
	metrics := metricsProxy.GetMetrics()
	if metrics["Display"][MetricRequestCount] != 1 {
		t.Errorf("Metrics proxy should have recorded one display call")
	}

	// Second display should use already loaded image
	err = metricsProxy.Display()
	if err != nil {
		t.Errorf("Second integration display failed: %v", err)
	}

	// Check metrics were updated
	metrics = metricsProxy.GetMetrics()
	if metrics["Display"][MetricRequestCount] != 2 {
		t.Errorf("Metrics proxy should have recorded two display calls")
	}

	fmt.Println("=== Integration Test Complete ===")
}
