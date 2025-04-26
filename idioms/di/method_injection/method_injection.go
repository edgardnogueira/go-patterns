// Package method_injection demonstrates dependency injection via method parameters.
package method_injection

import (
	"context"
	"fmt"
	"time"
)

// Formatter formats data for display
type Formatter interface {
	Format(data interface{}) string
}

// JSONFormatter formats data as JSON
type JSONFormatter struct{}

func (f *JSONFormatter) Format(data interface{}) string {
	// This is a simplified implementation
	return fmt.Sprintf(`{"value":"%v"}`, data)
}

// TextFormatter formats data as plain text
type TextFormatter struct{}

func (f *TextFormatter) Format(data interface{}) string {
	return fmt.Sprintf("%v", data)
}

// ReportGenerator creates reports with injected formatters
type ReportGenerator struct {
	// No formatter dependency as a field
	// Formatters will be injected via method parameters
}

// GenerateReport creates a report using the provided formatter
// The formatter is injected via the method parameter
func (g *ReportGenerator) GenerateReport(data []string, formatter Formatter) string {
	report := "Report:\n"
	
	for i, item := range data {
		formattedItem := formatter.Format(item)
		report += fmt.Sprintf("Item %d: %s\n", i+1, formattedItem)
	}
	
	return report
}

// Context-based method injection with request scoped dependencies
type RequestHandler struct {
	// Core dependencies could be here
}

// DatabaseConnection represents a database connection
type DatabaseConnection interface {
	Query(query string) ([]string, error)
}

// AuthService provides authentication-related functionality
type AuthService interface {
	CheckPermission(userID, resource string) bool
}

// Cache provides caching functionality
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, expiration time.Duration)
}

// RequestContext holds per-request dependencies
type RequestContext struct {
	context.Context
	DB            DatabaseConnection
	Auth          AuthService
	Cache         Cache
	User          string
	RequestID     string
	TraceEnabled  bool
}

// HandleUserProfile handles a user profile request
// Dependencies are injected via the RequestContext parameter
func (h *RequestHandler) HandleUserProfile(ctx *RequestContext, userID string) (string, error) {
	// Use the auth service from context
	if !ctx.Auth.CheckPermission(ctx.User, "user_profiles") {
		return "", fmt.Errorf("permission denied")
	}
	
	// Check cache first
	cacheKey := fmt.Sprintf("user_profile:%s", userID)
	if cachedProfile, found := ctx.Cache.Get(cacheKey); found {
		if ctx.TraceEnabled {
			fmt.Printf("[%s] Cache hit for user profile: %s\n", ctx.RequestID, userID)
		}
		return cachedProfile.(string), nil
	}
	
	// Query the database
	query := fmt.Sprintf("SELECT * FROM users WHERE id = '%s'", userID)
	results, err := ctx.DB.Query(query)
	if err != nil {
		return "", fmt.Errorf("database error: %w", err)
	}
	
	// Process results
	profile := fmt.Sprintf("Profile data for user %s: %v", userID, results)
	
	// Cache the result
	ctx.Cache.Set(cacheKey, profile, 5*time.Minute)
	
	return profile, nil
}

// ProcessWorkflow demonstrates chained method injection
func ProcessWorkflow(items []string) []string {
	// Initial data
	results := items
	
	// Apply multiple operations with different injected dependencies
	results = filterItems(results, func(item string) bool {
		return len(item) > 3
	})
	
	results = transformItems(results, func(item string) string {
		return fmt.Sprintf("[%s]", item)
	})
	
	results = sortItems(results, func(a, b string) bool {
		return len(a) < len(b)
	})
	
	return results
}

// Helper functions that use method injection for behavior customization

// filterItems filters a slice using the provided predicate
func filterItems(items []string, predicate func(string) bool) []string {
	var result []string
	for _, item := range items {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// transformItems transforms each item using the provided transformer
func transformItems(items []string, transformer func(string) string) []string {
	result := make([]string, len(items))
	for i, item := range items {
		result[i] = transformer(item)
	}
	return result
}

// sortItems sorts items using the provided comparator
func sortItems(items []string, lessFn func(string, string) bool) []string {
	result := make([]string, len(items))
	copy(result, items)
	
	// Simple bubble sort for demonstration
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if !lessFn(result[i], result[j]) {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	
	return result
}
