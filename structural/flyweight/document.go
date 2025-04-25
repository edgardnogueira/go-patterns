package flyweight

import (
	"bytes"
	"fmt"
	"strings"
)

// Character represents a character in a text document.
// It stores the extrinsic state (position, actual character) and a reference to the flyweight.
type Character struct {
	Value     rune        // The actual character value
	Row       int         // Row position in the document
	Column    int         // Column position in the document
	FormatID  string      // Reference to the format (flyweight)
}

// NewCharacter creates a new Character with the specified value, position, and format.
func NewCharacter(value rune, row, column int, formatID string) *Character {
	return &Character{
		Value:    value,
		Row:      row,
		Column:   column,
		FormatID: formatID,
	}
}

// String returns a string representation of the character.
func (c *Character) String() string {
	return fmt.Sprintf("%c@(%d,%d)[%s]", c.Value, c.Row, c.Column, c.FormatID)
}

// Document represents a text document with formatted characters.
// It's the client that uses the flyweight objects.
type Document struct {
	Name           string
	Characters     []*Character
	formatFactory  *TextFormatFactory
}

// NewDocument creates a new empty document.
func NewDocument(name string, formatFactory *TextFormatFactory) *Document {
	return &Document{
		Name:          name,
		Characters:    make([]*Character, 0),
		formatFactory: formatFactory,
	}
}

// AddText adds a string of text with the specified format to the document.
func (d *Document) AddText(text string, formatID string) {
	// Find the position where to add the text
	row, column := d.GetCurrentPosition()
	
	// Add each character with the specified format
	for i, c := range text {
		if c == '\n' {
			row++
			column = 0
		} else {
			char := NewCharacter(c, row, column+i, formatID)
			d.Characters = append(d.Characters, char)
			column++
		}
	}
}

// AddFormattedText adds a string of text with a new format to the document.
func (d *Document) AddFormattedText(text string, fontFamily string, fontSize int,
	fontColor string, isBold bool, isItalic bool, isUnderline bool,
	background string, alignment string,
	letterSpacing float64, lineHeight float64) {
	
	// Get or create the format
	format := d.formatFactory.GetTextFormat(
		fontFamily, fontSize, fontColor, isBold, isItalic, isUnderline,
		background, alignment, letterSpacing, lineHeight,
	)
	
	// Add the text with the format ID
	d.AddText(text, format.GetID())
}

// GetCurrentPosition returns the current row and column for adding new text.
func (d *Document) GetCurrentPosition() (int, int) {
	if len(d.Characters) == 0 {
		return 0, 0
	}
	
	lastChar := d.Characters[len(d.Characters)-1]
	return lastChar.Row, lastChar.Column + 1
}

// GetFormattedText returns the document text with all formatting applied.
func (d *Document) GetFormattedText() string {
	if len(d.Characters) == 0 {
		return ""
	}
	
	var result bytes.Buffer
	currentRow := 0
	currentFormatID := ""
	currentText := ""
	
	for _, char := range d.Characters {
		// Handle row changes
		if char.Row > currentRow {
			// Flush the current text with its format
			if currentText != "" {
				format := d.formatFactory.GetFormatByID(currentFormatID)
				if format != nil {
					result.WriteString(format.ApplyFormatting(currentText))
				} else {
					result.WriteString(currentText)
				}
				currentText = ""
			}
			
			// Add newlines
			for i := 0; i < char.Row-currentRow; i++ {
				result.WriteString("\n")
			}
			currentRow = char.Row
			currentFormatID = char.FormatID
		}
		
		// Handle format changes
		if char.FormatID != currentFormatID {
			// Flush the current text with its format
			if currentText != "" {
				format := d.formatFactory.GetFormatByID(currentFormatID)
				if format != nil {
					result.WriteString(format.ApplyFormatting(currentText))
				} else {
					result.WriteString(currentText)
				}
				currentText = ""
			}
			currentFormatID = char.FormatID
		}
		
		// Add the character to the current text
		currentText += string(char.Value)
	}
	
	// Flush any remaining text
	if currentText != "" {
		format := d.formatFactory.GetFormatByID(currentFormatID)
		if format != nil {
			result.WriteString(format.ApplyFormatting(currentText))
		} else {
			result.WriteString(currentText)
		}
	}
	
	return result.String()
}

// PlainText returns just the text content without formatting.
func (d *Document) PlainText() string {
	var result strings.Builder
	currentRow := 0
	
	for _, char := range d.Characters {
		// Add newlines for row changes
		if char.Row > currentRow {
			for i := 0; i < char.Row-currentRow; i++ {
				result.WriteString("\n")
			}
			currentRow = char.Row
		}
		
		// Add the character
		result.WriteRune(char.Value)
	}
	
	return result.String()
}

// GetMemoryUsage calculates the approximate memory usage of the document.
// It separates intrinsic (shared) state from extrinsic (per-character) state.
func (d *Document) GetMemoryUsage() map[string]interface{} {
	// In a real implementation, we would use more precise memory calculations
	// This is a simplified approximation
	
	// Calculate character memory (extrinsic state)
	characterCount := len(d.Characters)
	bytesPerChar := 24 // Approximate size of a Character struct in bytes (depends on architecture)
	totalCharMemory := characterCount * bytesPerChar
	
	// Get unique format count (intrinsic state)
	formatStats := d.formatFactory.GetCacheStats()
	uniqueFormatCount := formatStats["totalFormats"].(int)
	bytesPerFormat := 120 // Approximate size of a SharedTextFormat in bytes
	totalFormatMemory := uniqueFormatCount * bytesPerFormat
	
	// Without flyweight, each character would need its own format
	memoryWithoutFlyweight := characterCount * (bytesPerChar + bytesPerFormat)
	memorySaved := memoryWithoutFlyweight - (totalCharMemory + totalFormatMemory)
	savingsPercent := float64(0)
	if memoryWithoutFlyweight > 0 {
		savingsPercent = float64(memorySaved) / float64(memoryWithoutFlyweight) * 100
	}
	
	return map[string]interface{}{
		"characterCount":         characterCount,
		"uniqueFormatCount":      uniqueFormatCount,
		"characterMemory":        totalCharMemory,
		"formatMemory":           totalFormatMemory,
		"totalMemory":            totalCharMemory + totalFormatMemory,
		"memoryWithoutFlyweight": memoryWithoutFlyweight,
		"memorySaved":            memorySaved,
		"savingsPercent":         savingsPercent,
	}
}

// FormatStats returns statistics about the document formatting.
func (d *Document) FormatStats() map[string]interface{} {
	if len(d.Characters) == 0 {
		return map[string]interface{}{
			"characterCount": 0,
			"formatCount":    0,
		}
	}
	
	// Count characters by format
	formatCounts := make(map[string]int)
	for _, char := range d.Characters {
		formatCounts[char.FormatID] = formatCounts[char.FormatID] + 1
	}
	
	return map[string]interface{}{
		"characterCount": len(d.Characters),
		"formatCount":    len(formatCounts),
		"formatUsage":    formatCounts,
	}
}
