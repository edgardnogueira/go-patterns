// Package visitor implements the Visitor design pattern.
//
// The Visitor pattern lets you separate algorithms from the objects on which they operate.
// It allows adding new operations to existing object structures without modifying them.
package visitor

import (
	"fmt"
	"strings"
	"unicode"
)

// Element defines the Accept method for visitors
type Element interface {
	Accept(visitor Visitor) error
}

// TextElement represents text content in a document
type TextElement struct {
	Content string
}

// Accept implements the Element interface for TextElement
func (t *TextElement) Accept(visitor Visitor) error {
	return visitor.VisitText(t)
}

// ImageElement represents an image in a document
type ImageElement struct {
	Source string
	Alt    string
	Width  int
	Height int
}

// Accept implements the Element interface for ImageElement
func (i *ImageElement) Accept(visitor Visitor) error {
	return visitor.VisitImage(i)
}

// TableElement represents a table in a document
type TableElement struct {
	Rows    int
	Columns int
	Data    [][]string
}

// Accept implements the Element interface for TableElement
func (t *TableElement) Accept(visitor Visitor) error {
	return visitor.VisitTable(t)
}

// LinkElement represents a hyperlink in a document
type LinkElement struct {
	URL   string
	Text  string
	Title string
}

// Accept implements the Element interface for LinkElement
func (l *LinkElement) Accept(visitor Visitor) error {
	return visitor.VisitLink(l)
}

// CompositeElement can contain multiple elements
type CompositeElement struct {
	Name     string
	Children []Element
}

// Accept implements the Element interface for CompositeElement
// It visits the composite element itself and then all child elements
func (c *CompositeElement) Accept(visitor Visitor) error {
	// First visit the composite itself
	if err := visitor.VisitComposite(c); err != nil {
		return err
	}
	
	// Then visit all children
	for _, child := range c.Children {
		if err := child.Accept(visitor); err != nil {
			return err
		}
	}
	
	return nil
}

// AddChild adds a child element to the composite
func (c *CompositeElement) AddChild(element Element) {
	c.Children = append(c.Children, element)
}

// Visitor interface declares visit methods for each element type
type Visitor interface {
	VisitText(text *TextElement) error
	VisitImage(image *ImageElement) error
	VisitTable(table *TableElement) error
	VisitLink(link *LinkElement) error
	VisitComposite(composite *CompositeElement) error
}

// HTMLExportVisitor converts elements to HTML
type HTMLExportVisitor struct {
	Output      strings.Builder
	indentLevel int
}

// NewHTMLExportVisitor creates a new HTMLExportVisitor
func NewHTMLExportVisitor() *HTMLExportVisitor {
	return &HTMLExportVisitor{}
}

// indent returns the current indentation string
func (v *HTMLExportVisitor) indent() string {
	return strings.Repeat("  ", v.indentLevel)
}

// VisitText implements the Visitor interface for TextElement
func (v *HTMLExportVisitor) VisitText(text *TextElement) error {
	if text == nil {
		return fmt.Errorf("nil text element")
	}
	
	v.Output.WriteString(v.indent())
	v.Output.WriteString("<p>")
	v.Output.WriteString(text.Content)
	v.Output.WriteString("</p>\n")
	return nil
}

// VisitImage implements the Visitor interface for ImageElement
func (v *HTMLExportVisitor) VisitImage(image *ImageElement) error {
	if image == nil {
		return fmt.Errorf("nil image element")
	}
	
	v.Output.WriteString(v.indent())
	v.Output.WriteString(fmt.Sprintf("<img src=\"%s\" alt=\"%s\" width=\"%d\" height=\"%d\">\n", 
		image.Source, image.Alt, image.Width, image.Height))
	return nil
}

// VisitTable implements the Visitor interface for TableElement
func (v *HTMLExportVisitor) VisitTable(table *TableElement) error {
	if table == nil {
		return fmt.Errorf("nil table element")
	}
	
	v.Output.WriteString(v.indent())
	v.Output.WriteString("<table>\n")
	v.indentLevel++
	
	for i := 0; i < table.Rows; i++ {
		v.Output.WriteString(v.indent())
		v.Output.WriteString("<tr>\n")
		v.indentLevel++
		
		for j := 0; j < table.Columns; j++ {
			v.Output.WriteString(v.indent())
			cellData := ""
			if i < len(table.Data) && j < len(table.Data[i]) {
				cellData = table.Data[i][j]
			}
			v.Output.WriteString(fmt.Sprintf("<td>%s</td>\n", cellData))
		}
		
		v.indentLevel--
		v.Output.WriteString(v.indent())
		v.Output.WriteString("</tr>\n")
	}
	
	v.indentLevel--
	v.Output.WriteString(v.indent())
	v.Output.WriteString("</table>\n")
	return nil
}

// VisitLink implements the Visitor interface for LinkElement
func (v *HTMLExportVisitor) VisitLink(link *LinkElement) error {
	if link == nil {
		return fmt.Errorf("nil link element")
	}
	
	v.Output.WriteString(v.indent())
	if link.Title != "" {
		v.Output.WriteString(fmt.Sprintf("<a href=\"%s\" title=\"%s\">%s</a>\n", 
			link.URL, link.Title, link.Text))
	} else {
		v.Output.WriteString(fmt.Sprintf("<a href=\"%s\">%s</a>\n", 
			link.URL, link.Text))
	}
	return nil
}

// VisitComposite implements the Visitor interface for CompositeElement
func (v *HTMLExportVisitor) VisitComposite(composite *CompositeElement) error {
	if composite == nil {
		return fmt.Errorf("nil composite element")
	}
	
	v.Output.WriteString(v.indent())
	v.Output.WriteString(fmt.Sprintf("<div class=\"%s\">\n", strings.ToLower(composite.Name)))
	v.indentLevel++
	
	// Note: We don't process children here because the composite's Accept method does that
	
	return nil
}

// GetHTML returns the generated HTML as a string
func (v *HTMLExportVisitor) GetHTML() string {
	return v.Output.String()
}

// MarkdownExportVisitor converts elements to Markdown
type MarkdownExportVisitor struct {
	Output strings.Builder
	nesting int
}

// NewMarkdownExportVisitor creates a new MarkdownExportVisitor
func NewMarkdownExportVisitor() *MarkdownExportVisitor {
	return &MarkdownExportVisitor{}
}

// VisitText implements the Visitor interface for TextElement
func (v *MarkdownExportVisitor) VisitText(text *TextElement) error {
	if text == nil {
		return fmt.Errorf("nil text element")
	}
	
	v.Output.WriteString(text.Content)
	v.Output.WriteString("\n\n")
	return nil
}

// VisitImage implements the Visitor interface for ImageElement
func (v *MarkdownExportVisitor) VisitImage(image *ImageElement) error {
	if image == nil {
		return fmt.Errorf("nil image element")
	}
	
	v.Output.WriteString(fmt.Sprintf("![%s](%s)\n\n", image.Alt, image.Source))
	return nil
}

// VisitTable implements the Visitor interface for TableElement
func (v *MarkdownExportVisitor) VisitTable(table *TableElement) error {
	if table == nil {
		return fmt.Errorf("nil table element")
	}
	
	// Header row
	if table.Rows > 0 {
		for j := 0; j < table.Columns; j++ {
			v.Output.WriteString("| ")
			if j < len(table.Data[0]) {
				v.Output.WriteString(table.Data[0][j])
			}
			v.Output.WriteString(" ")
		}
		v.Output.WriteString("|\n")
		
		// Separator row
		for j := 0; j < table.Columns; j++ {
			v.Output.WriteString("| --- ")
		}
		v.Output.WriteString("|\n")
		
		// Data rows
		for i := 1; i < table.Rows; i++ {
			for j := 0; j < table.Columns; j++ {
				v.Output.WriteString("| ")
				if i < len(table.Data) && j < len(table.Data[i]) {
					v.Output.WriteString(table.Data[i][j])
				}
				v.Output.WriteString(" ")
			}
			v.Output.WriteString("|\n")
		}
	}
	
	v.Output.WriteString("\n")
	return nil
}

// VisitLink implements the Visitor interface for LinkElement
func (v *MarkdownExportVisitor) VisitLink(link *LinkElement) error {
	if link == nil {
		return fmt.Errorf("nil link element")
	}
	
	v.Output.WriteString(fmt.Sprintf("[%s](%s)", link.Text, link.URL))
	return nil
}

// VisitComposite implements the Visitor interface for CompositeElement
func (v *MarkdownExportVisitor) VisitComposite(composite *CompositeElement) error {
	if composite == nil {
		return fmt.Errorf("nil composite element")
	}
	
	v.nesting++
	v.Output.WriteString(fmt.Sprintf("%s %s\n\n", strings.Repeat("#", v.nesting), composite.Name))
	
	// Note: We don't process children here because the composite's Accept method does that
	
	return nil
}

// GetMarkdown returns the generated markdown as a string
func (v *MarkdownExportVisitor) GetMarkdown() string {
	return v.Output.String()
}

// PlainTextExportVisitor extracts plain text content
type PlainTextExportVisitor struct {
	Output strings.Builder
}

// NewPlainTextExportVisitor creates a new PlainTextExportVisitor
func NewPlainTextExportVisitor() *PlainTextExportVisitor {
	return &PlainTextExportVisitor{}
}

// VisitText implements the Visitor interface for TextElement
func (v *PlainTextExportVisitor) VisitText(text *TextElement) error {
	if text == nil {
		return fmt.Errorf("nil text element")
	}
	
	v.Output.WriteString(text.Content)
	v.Output.WriteString("\n\n")
	return nil
}

// VisitImage implements the Visitor interface for ImageElement
func (v *PlainTextExportVisitor) VisitImage(image *ImageElement) error {
	if image == nil {
		return fmt.Errorf("nil image element")
	}
	
	v.Output.WriteString(fmt.Sprintf("[Image: %s]\n", image.Alt))
	return nil
}

// VisitTable implements the Visitor interface for TableElement
func (v *PlainTextExportVisitor) VisitTable(table *TableElement) error {
	if table == nil {
		return fmt.Errorf("nil table element")
	}
	
	for i := 0; i < table.Rows; i++ {
		for j := 0; j < table.Columns; j++ {
			if j > 0 {
				v.Output.WriteString("\t")
			}
			
			if i < len(table.Data) && j < len(table.Data[i]) {
				v.Output.WriteString(table.Data[i][j])
			}
		}
		v.Output.WriteString("\n")
	}
	v.Output.WriteString("\n")
	return nil
}

// VisitLink implements the Visitor interface for LinkElement
func (v *PlainTextExportVisitor) VisitLink(link *LinkElement) error {
	if link == nil {
		return fmt.Errorf("nil link element")
	}
	
	v.Output.WriteString(link.Text)
	v.Output.WriteString(" (")
	v.Output.WriteString(link.URL)
	v.Output.WriteString(")")
	return nil
}

// VisitComposite implements the Visitor interface for CompositeElement
func (v *PlainTextExportVisitor) VisitComposite(composite *CompositeElement) error {
	if composite == nil {
		return fmt.Errorf("nil composite element")
	}
	
	// Note: We don't process children here because the composite's Accept method does that
	return nil
}

// GetText returns the generated plain text
func (v *PlainTextExportVisitor) GetText() string {
	return v.Output.String()
}

// StatisticsVisitor collects document statistics
type StatisticsVisitor struct {
	TextCount     int
	ImageCount    int
	TableCount    int
	LinkCount     int
	WordCount     int
	CharacterCount int
}

// NewStatisticsVisitor creates a new StatisticsVisitor
func NewStatisticsVisitor() *StatisticsVisitor {
	return &StatisticsVisitor{}
}

// VisitText implements the Visitor interface for TextElement
func (v *StatisticsVisitor) VisitText(text *TextElement) error {
	if text == nil {
		return fmt.Errorf("nil text element")
	}
	
	v.TextCount++
	v.WordCount += len(strings.Fields(text.Content))
	v.CharacterCount += len(text.Content)
	return nil
}

// VisitImage implements the Visitor interface for ImageElement
func (v *StatisticsVisitor) VisitImage(image *ImageElement) error {
	if image == nil {
		return fmt.Errorf("nil image element")
	}
	
	v.ImageCount++
	return nil
}

// VisitTable implements the Visitor interface for TableElement
func (v *StatisticsVisitor) VisitTable(table *TableElement) error {
	if table == nil {
		return fmt.Errorf("nil table element")
	}
	
	v.TableCount++
	for i := 0; i < len(table.Data); i++ {
		for j := 0; j < len(table.Data[i]); j++ {
			content := table.Data[i][j]
			v.WordCount += len(strings.Fields(content))
			v.CharacterCount += len(content)
		}
	}
	return nil
}

// VisitLink implements the Visitor interface for LinkElement
func (v *StatisticsVisitor) VisitLink(link *LinkElement) error {
	if link == nil {
		return fmt.Errorf("nil link element")
	}
	
	v.LinkCount++
	v.WordCount += len(strings.Fields(link.Text))
	v.CharacterCount += len(link.Text)
	return nil
}

// VisitComposite implements the Visitor interface for CompositeElement
func (v *StatisticsVisitor) VisitComposite(composite *CompositeElement) error {
	if composite == nil {
		return fmt.Errorf("nil composite element")
	}
	
	// We don't count the composite itself as a statistic
	// Just process all its children (but this is done in the Accept method)
	return nil
}

// SpellCheckVisitor checks spelling across document elements
type SpellCheckVisitor struct {
	Errors []string
	dictionary map[string]bool // Simple dictionary implementation
}

// NewSpellCheckVisitor creates a new SpellCheckVisitor with a basic dictionary
func NewSpellCheckVisitor() *SpellCheckVisitor {
	// This is a very simple dictionary for demonstration purposes
	// In a real implementation, you would use a proper spell checking library
	dictionary := map[string]bool{
		"the": true, "and": true, "a": true, "to": true, "in": true,
		"that": true, "it": true, "with": true, "as": true, "for": true,
		"was": true, "on": true, "are": true, "by": true, "this": true,
		"be": true, "is": true, "from": true, "at": true, "an": true,
		"but": true, "not": true, "or": true, "what": true, "all": true,
		"were": true, "when": true, "we": true, "there": true, "can": true,
		"an": true, "your": true, "which": true, "their": true, "said": true,
		"if": true, "will": true, "each": true, "about": true, "how": true,
		"up": true, "out": true, "them": true, "then": true, "she": true,
		"many": true, "some": true, "so": true, "these": true, "would": true,
		"other": true, "into": true, "has": true, "more": true, "two": true,
		"like": true, "him": true, "see": true, "time": true, "could": true,
		"no": true, "make": true, "than": true, "first": true, "been": true,
		"its": true, "who": true, "now": true, "people": true, "my": true,
		"made": true, "over": true, "did": true, "down": true, "only": true,
		"way": true, "find": true, "use": true, "may": true, "water": true,
		"long": true, "little": true, "very": true, "after": true, "words": true,
		"called": true, "just": true, "where": true, "most": true, "know": true,
		// Add more common words as needed
	}
	
	return &SpellCheckVisitor{
		dictionary: dictionary,
	}
}

// checkWord validates if a word is in the dictionary
func (v *SpellCheckVisitor) checkWord(word string) bool {
	// Convert to lowercase and strip punctuation for checking
	cleaned := strings.TrimFunc(word, func(r rune) bool {
		return !unicode.IsLetter(r)
	})
	
	cleaned = strings.ToLower(cleaned)
	if cleaned == "" {
		return true // Skip empty strings
	}
	
	return v.dictionary[cleaned]
}

// VisitText implements the Visitor interface for TextElement
func (v *SpellCheckVisitor) VisitText(text *TextElement) error {
	if text == nil {
		return fmt.Errorf("nil text element")
	}
	
	words := strings.Fields(text.Content)
	for _, word := range words {
		if !v.checkWord(word) {
			v.Errors = append(v.Errors, fmt.Sprintf("Possible spelling error: %s", word))
		}
	}
	return nil
}

// VisitImage implements the Visitor interface for ImageElement
func (v *SpellCheckVisitor) VisitImage(image *ImageElement) error {
	if image == nil {
		return fmt.Errorf("nil image element")
	}
	
	// Check alt text
	words := strings.Fields(image.Alt)
	for _, word := range words {
		if !v.checkWord(word) {
			v.Errors = append(v.Errors, fmt.Sprintf("Possible spelling error in image alt text: %s", word))
		}
	}
	return nil
}

// VisitTable implements the Visitor interface for TableElement
func (v *SpellCheckVisitor) VisitTable(table *TableElement) error {
	if table == nil {
		return fmt.Errorf("nil table element")
	}
	
	for i := 0; i < len(table.Data); i++ {
		for j := 0; j < len(table.Data[i]); j++ {
			content := table.Data[i][j]
			words := strings.Fields(content)
			for _, word := range words {
				if !v.checkWord(word) {
					v.Errors = append(v.Errors, 
						fmt.Sprintf("Possible spelling error in table at [%d,%d]: %s", i, j, word))
				}
			}
		}
	}
	return nil
}

// VisitLink implements the Visitor interface for LinkElement
func (v *SpellCheckVisitor) VisitLink(link *LinkElement) error {
	if link == nil {
		return fmt.Errorf("nil link element")
	}
	
	// Check link text
	words := strings.Fields(link.Text)
	for _, word := range words {
		if !v.checkWord(word) {
			v.Errors = append(v.Errors, fmt.Sprintf("Possible spelling error in link text: %s", word))
		}
	}
	
	// Check title if available
	if link.Title != "" {
		words = strings.Fields(link.Title)
		for _, word := range words {
			if !v.checkWord(word) {
				v.Errors = append(v.Errors, fmt.Sprintf("Possible spelling error in link title: %s", word))
			}
		}
	}
	return nil
}

// VisitComposite implements the Visitor interface for CompositeElement
func (v *SpellCheckVisitor) VisitComposite(composite *CompositeElement) error {
	if composite == nil {
		return fmt.Errorf("nil composite element")
	}
	
	// Check name
	words := strings.Fields(composite.Name)
	for _, word := range words {
		if !v.checkWord(word) {
			v.Errors = append(v.Errors, fmt.Sprintf("Possible spelling error in section name: %s", word))
		}
	}
	return nil
}

// GetErrors returns all the spelling errors found
func (v *SpellCheckVisitor) GetErrors() []string {
	return v.Errors
}
