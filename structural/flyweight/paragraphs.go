package flyweight

import (
	"fmt"
	"sync"
)

// ParagraphStyle represents formatting options for an entire paragraph.
// These are also flyweights, shared between many paragraphs.
type ParagraphStyle interface {
	// GetID returns a unique identifier for this paragraph style
	GetID() string
	
	// GetStyle returns a map of style properties
	GetStyle() map[string]interface{}
	
	// FormatParagraph applies the style to a paragraph
	FormatParagraph(content string) string
	
	// String returns a string representation of this style
	String() string
}

// SharedParagraphStyle is a concrete flyweight that stores intrinsic (shared) state
// for paragraph formatting.
type SharedParagraphStyle struct {
	id             string
	alignment      string // "left", "center", "right", "justify"
	lineSpacing    float64
	beforeSpacing  float64
	afterSpacing   float64
	firstLineIndent float64
	leftMargin     float64
	rightMargin    float64
	borderStyle    string // "none", "single", "double", etc.
	borderColor    string
	backgroundColor string
}

// NewSharedParagraphStyle creates a new SharedParagraphStyle.
func NewSharedParagraphStyle(
	alignment string,
	lineSpacing float64,
	beforeSpacing float64,
	afterSpacing float64,
	firstLineIndent float64,
	leftMargin float64,
	rightMargin float64,
	borderStyle string,
	borderColor string,
	backgroundColor string,
) *SharedParagraphStyle {
	// Generate an ID based on the properties
	id := fmt.Sprintf("%s-%.1f-%.1f-%.1f-%.1f-%.1f-%.1f-%s-%s-%s",
		alignment, lineSpacing, beforeSpacing, afterSpacing,
		firstLineIndent, leftMargin, rightMargin,
		borderStyle, borderColor, backgroundColor)
	
	return &SharedParagraphStyle{
		id:             id,
		alignment:      alignment,
		lineSpacing:    lineSpacing,
		beforeSpacing:  beforeSpacing,
		afterSpacing:   afterSpacing,
		firstLineIndent: firstLineIndent,
		leftMargin:     leftMargin,
		rightMargin:    rightMargin,
		borderStyle:    borderStyle,
		borderColor:    borderColor,
		backgroundColor: backgroundColor,
	}
}

// GetID returns a unique identifier for this paragraph style.
func (p *SharedParagraphStyle) GetID() string {
	return p.id
}

// GetStyle returns a map of style properties.
func (p *SharedParagraphStyle) GetStyle() map[string]interface{} {
	return map[string]interface{}{
		"alignment":      p.alignment,
		"lineSpacing":    p.lineSpacing,
		"beforeSpacing":  p.beforeSpacing,
		"afterSpacing":   p.afterSpacing,
		"firstLineIndent": p.firstLineIndent,
		"leftMargin":     p.leftMargin,
		"rightMargin":    p.rightMargin,
		"borderStyle":    p.borderStyle,
		"borderColor":    p.borderColor,
		"backgroundColor": p.backgroundColor,
	}
}

// FormatParagraph applies the style to a paragraph.
// In a real implementation, this would apply actual formatting.
func (p *SharedParagraphStyle) FormatParagraph(content string) string {
	indentStr := ""
	if p.firstLineIndent > 0 {
		for i := 0; i < int(p.firstLineIndent); i++ {
			indentStr += " "
		}
	}
	
	// Apply margins in a basic way
	leftMarginStr := ""
	for i := 0; i < int(p.leftMargin); i++ {
		leftMarginStr += " "
	}
	
	// Simple representation of paragraph formatting
	return fmt.Sprintf("[PARAGRAPH: %s align, %.1f spacing, %s%s%s]\n%s%s%s",
		p.alignment, p.lineSpacing,
		p.borderStyle != "none" ? "[" + p.borderStyle + " border] " : "",
		p.backgroundColor != "transparent" ? "[" + p.backgroundColor + " bg] " : "",
		p.firstLineIndent > 0 ? fmt.Sprintf("[%.1fpx indent] ", p.firstLineIndent) : "",
		leftMarginStr, indentStr, content)
}

// String returns a string representation of this style.
func (p *SharedParagraphStyle) String() string {
	return fmt.Sprintf("ParagraphStyle[%s, spacing=%.1f, indent=%.1f, margins=%.1f/%.1f]",
		p.alignment, p.lineSpacing, p.firstLineIndent, p.leftMargin, p.rightMargin)
}

// ParagraphStyleFactory creates and manages paragraph style flyweights.
type ParagraphStyleFactory struct {
	styles map[string]ParagraphStyle
	mutex  sync.RWMutex
}

// NewParagraphStyleFactory creates a new paragraph style factory.
func NewParagraphStyleFactory() *ParagraphStyleFactory {
	return &ParagraphStyleFactory{
		styles: make(map[string]ParagraphStyle),
	}
}

// GetParagraphStyle returns a flyweight paragraph style with the specified properties.
func (f *ParagraphStyleFactory) GetParagraphStyle(
	alignment string,
	lineSpacing float64,
	beforeSpacing float64,
	afterSpacing float64,
	firstLineIndent float64,
	leftMargin float64,
	rightMargin float64,
	borderStyle string,
	borderColor string,
	backgroundColor string,
) ParagraphStyle {
	// Create a temporary style to get its ID
	tempStyle := NewSharedParagraphStyle(
		alignment, lineSpacing, beforeSpacing, afterSpacing,
		firstLineIndent, leftMargin, rightMargin,
		borderStyle, borderColor, backgroundColor,
	)
	id := tempStyle.GetID()
	
	// Check if we already have this style
	f.mutex.RLock()
	if style, exists := f.styles[id]; exists {
		f.mutex.RUnlock()
		return style
	}
	f.mutex.RUnlock()
	
	// If we don't have it, store the temp style we created
	f.mutex.Lock()
	defer f.mutex.Unlock()
	
	// Double-check in case another goroutine created it while we were waiting
	if style, exists := f.styles[id]; exists {
		return style
	}
	
	// Store and return the new style
	f.styles[id] = tempStyle
	return tempStyle
}

// GetStyleByID returns a style by its ID, or nil if it doesn't exist.
func (f *ParagraphStyleFactory) GetStyleByID(id string) ParagraphStyle {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	
	if style, exists := f.styles[id]; exists {
		return style
	}
	return nil
}

// GetCacheStats returns statistics about the style cache.
func (f *ParagraphStyleFactory) GetCacheStats() map[string]interface{} {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	
	stats := map[string]interface{}{
		"totalStyles": len(f.styles),
	}
	
	// Count styles by properties
	alignCounts := make(map[string]int)
	borderCounts := make(map[string]int)
	bgCounts := make(map[string]int)
	
	for _, style := range f.styles {
		props := style.GetStyle()
		
		alignment := props["alignment"].(string)
		alignCounts[alignment] = alignCounts[alignment] + 1
		
		borderStyle := props["borderStyle"].(string)
		borderCounts[borderStyle] = borderCounts[borderStyle] + 1
		
		bgColor := props["backgroundColor"].(string)
		bgCounts[bgColor] = bgCounts[bgColor] + 1
	}
	
	stats["alignments"] = alignCounts
	stats["borderStyles"] = borderCounts
	stats["backgroundColors"] = bgCounts
	
	return stats
}

// ClearCache clears all stored styles.
func (f *ParagraphStyleFactory) ClearCache() {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	
	f.styles = make(map[string]ParagraphStyle)
}

// Paragraph represents a paragraph in a document,
// combining a paragraph style with its content.
type Paragraph struct {
	Content     string
	StyleID     string
	Row         int
}

// NewParagraph creates a new paragraph with the specified content and style.
func NewParagraph(content string, styleID string, row int) *Paragraph {
	return &Paragraph{
		Content: content,
		StyleID: styleID,
		Row:     row,
	}
}

// FormattedDocument extends the basic Document with paragraph styling.
type FormattedDocument struct {
	*Document
	Paragraphs       []*Paragraph
	paragraphFactory *ParagraphStyleFactory
}

// NewFormattedDocument creates a new document with paragraph formatting.
func NewFormattedDocument(name string, textFactory *TextFormatFactory, paragraphFactory *ParagraphStyleFactory) *FormattedDocument {
	return &FormattedDocument{
		Document:         NewDocument(name, textFactory),
		Paragraphs:       make([]*Paragraph, 0),
		paragraphFactory: paragraphFactory,
	}
}

// AddParagraph adds a paragraph with formatting to the document.
func (d *FormattedDocument) AddParagraph(content string, styleID string) {
	row := 0
	if len(d.Paragraphs) > 0 {
		lastPara := d.Paragraphs[len(d.Paragraphs)-1]
		// Calculate row based on previous paragraph position and lines
		row = lastPara.Row + 1 + len(content)/80 // Rough estimate
	}
	
	para := NewParagraph(content, styleID, row)
	d.Paragraphs = append(d.Paragraphs, para)
}

// GetFormattedDocument returns the document with all text and paragraph formatting applied.
func (d *FormattedDocument) GetFormattedDocument() string {
	result := ""
	
	for _, para := range d.Paragraphs {
		style := d.paragraphFactory.GetStyleByID(para.StyleID)
		if style != nil {
			result += style.FormatParagraph(para.Content) + "\n\n"
		} else {
			result += para.Content + "\n\n"
		}
	}
	
	return result
}
