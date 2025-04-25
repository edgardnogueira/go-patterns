// Package flyweight implements the Flyweight design pattern in Go.
//
// The Flyweight pattern minimizes memory usage by sharing as much data as possible with similar objects.
// This implementation demonstrates a text formatting system for a document editor
// where character formatting objects are shared across the document.
package flyweight

import (
	"fmt"
	"sync"
)

// TextFormat is the flyweight interface that defines methods for applying text formatting.
type TextFormat interface {
	// GetID returns a unique identifier for this text format
	GetID() string
	
	// GetFormat returns a map of formatting properties
	GetFormat() map[string]interface{}
	
	// ApplyFormatting applies this formatting to a string
	ApplyFormatting(s string) string
	
	// String returns a string representation of this format
	String() string
}

// SharedTextFormat is a concrete flyweight that stores intrinsic (shared) state.
// It represents formatting options like font, size, color, etc.
type SharedTextFormat struct {
	id          string
	fontFamily  string
	fontSize    int
	fontColor   string
	isBold      bool
	isItalic    bool
	isUnderline bool
	background  string
	alignment   string
	letterSpacing float64
	lineHeight    float64
}

// NewSharedTextFormat creates a new SharedTextFormat with the given parameters.
func NewSharedTextFormat(
	fontFamily string,
	fontSize int,
	fontColor string,
	isBold bool,
	isItalic bool,
	isUnderline bool,
	background string,
	alignment string,
	letterSpacing float64,
	lineHeight float64,
) *SharedTextFormat {
	// Generate an ID based on the formatting properties
	id := fmt.Sprintf("%s-%d-%s-%v-%v-%v-%s-%s-%.1f-%.1f",
		fontFamily, fontSize, fontColor, isBold, isItalic, isUnderline,
		background, alignment, letterSpacing, lineHeight)
	
	return &SharedTextFormat{
		id:            id,
		fontFamily:    fontFamily,
		fontSize:      fontSize,
		fontColor:     fontColor,
		isBold:        isBold,
		isItalic:      isItalic,
		isUnderline:   isUnderline,
		background:    background,
		alignment:     alignment,
		letterSpacing: letterSpacing,
		lineHeight:    lineHeight,
	}
}

// GetID returns a unique identifier for this text format.
func (f *SharedTextFormat) GetID() string {
	return f.id
}

// GetFormat returns a map of formatting properties.
func (f *SharedTextFormat) GetFormat() map[string]interface{} {
	return map[string]interface{}{
		"fontFamily":    f.fontFamily,
		"fontSize":      f.fontSize,
		"fontColor":     f.fontColor,
		"isBold":        f.isBold,
		"isItalic":      f.isItalic,
		"isUnderline":   f.isUnderline,
		"background":    f.background,
		"alignment":     f.alignment,
		"letterSpacing": f.letterSpacing,
		"lineHeight":    f.lineHeight,
	}
}

// ApplyFormatting applies this formatting to a string.
// In a real implementation, this would actually style the text.
func (f *SharedTextFormat) ApplyFormatting(s string) string {
	// This is a simplified implementation. In a real editor,
	// this would apply actual formatting to the text.
	descriptor := ""
	if f.isBold {
		descriptor += "bold "
	}
	if f.isItalic {
		descriptor += "italic "
	}
	if f.isUnderline {
		descriptor += "underlined "
	}
	
	return fmt.Sprintf("[%s%s, %dpx, %s on %s]: %s",
		descriptor, f.fontFamily, f.fontSize, f.fontColor, f.background, s)
}

// String returns a string representation of this format.
func (f *SharedTextFormat) String() string {
	return fmt.Sprintf("TextFormat[%s, %dpx, %s, bold=%v, italic=%v, underline=%v, bg=%s, align=%s]",
		f.fontFamily, f.fontSize, f.fontColor, f.isBold, f.isItalic, f.isUnderline,
		f.background, f.alignment)
}

// TextFormatFactory creates and manages flyweight objects.
// It ensures that flyweights are shared properly.
type TextFormatFactory struct {
	formats map[string]TextFormat
	mutex   sync.RWMutex
}

// NewTextFormatFactory creates a new TextFormatFactory.
func NewTextFormatFactory() *TextFormatFactory {
	return &TextFormatFactory{
		formats: make(map[string]TextFormat),
	}
}

// GetTextFormat returns a flyweight object with the specified formatting.
// If a format with the given parameters already exists, it returns the existing one.
// Otherwise, it creates a new one.
func (f *TextFormatFactory) GetTextFormat(
	fontFamily string,
	fontSize int,
	fontColor string,
	isBold bool,
	isItalic bool,
	isUnderline bool,
	background string,
	alignment string,
	letterSpacing float64,
	lineHeight float64,
) TextFormat {
	// Create a temporary format to get its ID
	tempFormat := NewSharedTextFormat(
		fontFamily, fontSize, fontColor, isBold, isItalic, isUnderline,
		background, alignment, letterSpacing, lineHeight,
	)
	id := tempFormat.GetID()
	
	// Check if we already have this format
	f.mutex.RLock()
	if format, exists := f.formats[id]; exists {
		f.mutex.RUnlock()
		return format
	}
	f.mutex.RUnlock()
	
	// If we don't have it, store the temp format we created
	f.mutex.Lock()
	defer f.mutex.Unlock()
	
	// Double-check in case another goroutine created it while we were waiting
	if format, exists := f.formats[id]; exists {
		return format
	}
	
	// Store and return the new format
	f.formats[id] = tempFormat
	return tempFormat
}

// GetFormatByID returns a format by its ID, or nil if it doesn't exist.
func (f *TextFormatFactory) GetFormatByID(id string) TextFormat {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	
	if format, exists := f.formats[id]; exists {
		return format
	}
	return nil
}

// GetCacheStats returns statistics about the format cache.
func (f *TextFormatFactory) GetCacheStats() map[string]interface{} {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	
	stats := map[string]interface{}{
		"totalFormats": len(f.formats),
	}
	
	// Count formats by properties
	boldCount := 0
	italicCount := 0
	underlineCount := 0
	fontCounts := make(map[string]int)
	sizeCounts := make(map[int]int)
	colorCounts := make(map[string]int)
	
	for _, format := range f.formats {
		props := format.GetFormat()
		
		if props["isBold"].(bool) {
			boldCount++
		}
		if props["isItalic"].(bool) {
			italicCount++
		}
		if props["isUnderline"].(bool) {
			underlineCount++
		}
		
		fontFamily := props["fontFamily"].(string)
		fontCounts[fontFamily] = fontCounts[fontFamily] + 1
		
		fontSize := props["fontSize"].(int)
		sizeCounts[fontSize] = sizeCounts[fontSize] + 1
		
		fontColor := props["fontColor"].(string)
		colorCounts[fontColor] = colorCounts[fontColor] + 1
	}
	
	stats["boldCount"] = boldCount
	stats["italicCount"] = italicCount
	stats["underlineCount"] = underlineCount
	stats["fontFamilies"] = fontCounts
	stats["fontSizes"] = sizeCounts
	stats["fontColors"] = colorCounts
	
	return stats
}

// ClearCache clears all stored formats.
func (f *TextFormatFactory) ClearCache() {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	
	f.formats = make(map[string]TextFormat)
}
