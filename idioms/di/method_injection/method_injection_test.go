package method_injection

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestReportGenerator(t *testing.T) {
	// Create report generator
	generator := &ReportGenerator{}
	
	// Test data
	data := []string{"apple", "banana", "cherry"}
	
	// Test with JSON formatter
	jsonFormatter := &JSONFormatter{}
	jsonReport := generator.GenerateReport(data, jsonFormatter)
	
	// Verify JSON formatting
	if !strings.Contains(jsonReport, `{"value":"apple"}`) {
		t.Errorf("Expected JSON formatting, got: %s", jsonReport)
	}
	
	// Test with text formatter
	textFormatter := &TextFormatter{}
	textReport := generator.GenerateReport(data, textFormatter)
	
	// Verify text formatting
	if !strings.Contains(textReport, "apple") {
		t.Errorf("Expected plain text formatting, got: %s", textReport)
	}
	
	// Different formatter returns different output
	if jsonReport == textReport {
		t.Error("Expected different outputs for different formatters")
	}
}

// Mock implementations for testing RequestHandler
type MockDB struct {
	results []string
	err     error
}

func (m *MockDB) Query(query string) ([]string, error) {
	return m.results, m.err
}

type MockAuth struct {
	allowed bool
}

func (m *MockAuth) CheckPermission(userID, resource string) bool {
	return m.allowed
}

type MockCache struct {
	data map[string]interface{}
}

func NewMockCache() *MockCache {
	return &MockCache{
		data: make(map[string]interface{}),
	}
}

func (m *MockCache) Get(key string) (interface{}, bool) {
	val, ok := m.data[key]
	return val, ok
}

func (m *MockCache) Set(key string, value interface{}, expiration time.Duration) {
	m.data[key] = value
}

func TestRequestHandler(t *testing.T) {
	// Create handler
	handler := &RequestHandler{}
	
	// Create mock dependencies
	mockDB := &MockDB{
		results: []string{"user data"},
	}
	
	mockAuth := &MockAuth{
		allowed: true,
	}
	
	mockCache := NewMockCache()
	
	// Create request context with injected dependencies
	ctx := &RequestContext{
		Context:      context.Background(),
		DB:           mockDB,
		Auth:         mockAuth,
		Cache:        mockCache,
		User:         "admin",
		RequestID:    "req-123",
		TraceEnabled: true,
	}
	
	// Call method with injected context
	profile, err := handler.HandleUserProfile(ctx, "user-42")
	
	// Verify method used injected dependencies correctly
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if !strings.Contains(profile, "user data") {
		t.Errorf("Expected profile to contain user data, got: %s", profile)
	}
	
	// Verify cache was populated
	cacheKey := "user_profile:user-42"
	cachedProfile, found := mockCache.Get(cacheKey)
	if !found {
		t.Error("Expected profile to be cached")
	}
	
	if cachedProfile != profile {
		t.Error("Cached profile doesn't match returned profile")
	}
	
	// Test cache hit path
	profile2, err := handler.HandleUserProfile(ctx, "user-42")
	if err != nil {
		t.Errorf("Expected no error on cache hit, got %v", err)
	}
	
	if profile2 != profile {
		t.Error("Profile from cache hit doesn't match expected profile")
	}
	
	// Test permission denied
	mockAuth.allowed = false
	_, err = handler.HandleUserProfile(ctx, "user-42")
	if err == nil || !strings.Contains(err.Error(), "permission denied") {
		t.Errorf("Expected permission denied error, got: %v", err)
	}
}

func TestProcessWorkflow(t *testing.T) {
	// Test data
	items := []string{"a", "foo", "barbaz", "testing", "x"}
	
	// Process items with method-injected behaviors
	results := ProcessWorkflow(items)
	
	// Verify processing - should have filtered, transformed, and sorted
	expectedResults := []string{"[foo]", "[barbaz]", "[testing]"}
	
	if len(results) != len(expectedResults) {
		t.Errorf("Expected %d results, got %d", len(expectedResults), len(results))
	}
	
	// Verify sorting is by length
	for i := 0; i < len(results)-1; i++ {
		if len(results[i]) > len(results[i+1]) {
			t.Errorf("Expected sorting by length, but got: %v", results)
		}
	}
	
	// Verify all items are transformed (have brackets)
	for _, item := range results {
		if !strings.HasPrefix(item, "[") || !strings.HasSuffix(item, "]") {
			t.Errorf("Item not properly transformed: %s", item)
		}
	}
	
	// Verify all short items were filtered
	for _, item := range results {
		// Remove brackets to check original length
		original := item[1 : len(item)-1]
		if len(original) <= 3 {
			t.Errorf("Expected items with length > 3, got: %s", original)
		}
	}
}

// Custom formatters for testing flexibility
type UppercaseFormatter struct{}

func (f *UppercaseFormatter) Format(data interface{}) string {
	return strings.ToUpper(fmt.Sprintf("%v", data))
}

func TestCustomFormatters(t *testing.T) {
	// Show how method injection allows for easy extension
	generator := &ReportGenerator{}
	data := []string{"apple", "banana"}
	
	// Use a custom formatter not defined in the original code
	uppercaseFormatter := &UppercaseFormatter{}
	report := generator.GenerateReport(data, uppercaseFormatter)
	
	if !strings.Contains(report, "APPLE") {
		t.Errorf("Expected uppercase formatting, got: %s", report)
	}
}
