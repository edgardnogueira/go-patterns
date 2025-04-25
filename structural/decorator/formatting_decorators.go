package decorator

import (
	"fmt"
	"regexp"
	"strings"
)

// FormattingDecorator is a concrete decorator that formats text.
// It decorates a TextProcessor by adding formatting capabilities.
type FormattingDecorator struct {
	TextProcessorDecorator
	formatType string
}

// NewHTMLFormattingDecorator creates a decorator that converts text to HTML.
func NewHTMLFormattingDecorator(processor TextProcessor) *FormattingDecorator {
	return &FormattingDecorator{
		TextProcessorDecorator: TextProcessorDecorator{
			wrapped:     processor,
			name:        "HTML Formatter",
			description: "Formats text as HTML with proper tags",
		},
		formatType: "html",
	}
}

// NewMarkdownFormattingDecorator creates a decorator that converts text to Markdown.
func NewMarkdownFormattingDecorator(processor TextProcessor) *FormattingDecorator {
	return &FormattingDecorator{
		TextProcessorDecorator: TextProcessorDecorator{
			wrapped:     processor,
			name:        "Markdown Formatter",
			description: "Formats text as Markdown with proper syntax",
		},
		formatType: "markdown",
	}
}

// NewPlainTextFormattingDecorator creates a decorator that strips formatting.
func NewPlainTextFormattingDecorator(processor TextProcessor) *FormattingDecorator {
	return &FormattingDecorator{
		TextProcessorDecorator: TextProcessorDecorator{
			wrapped:     processor,
			name:        "Plain Text Formatter",
			description: "Strips formatting and converts to plain text",
		},
		formatType: "plain",
	}
}

// Process first processes the text using the wrapped processor,
// then applies the appropriate formatting based on the formatType.
func (f *FormattingDecorator) Process(text string) (string, error) {
	// First, let the wrapped processor do its work
	processedText, err := f.wrapped.Process(text)
	if err != nil {
		return "", fmt.Errorf("error in wrapped processor: %w", err)
	}

	// Then apply formatting based on the type
	switch f.formatType {
	case "html":
		return f.formatAsHTML(processedText), nil
	case "markdown":
		return f.formatAsMarkdown(processedText), nil
	case "plain":
		return f.formatAsPlainText(processedText), nil
	default:
		return processedText, nil
	}
}

// formatAsHTML converts text to HTML format.
func (f *FormattingDecorator) formatAsHTML(text string) string {
	// Convert paragraphs
	paragraphs := strings.Split(text, "\n\n")
	for i, p := range paragraphs {
		if p != "" {
			paragraphs[i] = "<p>" + p + "</p>"
		}
	}
	
	// Convert headings
	re := regexp.MustCompile(`(?m)^(#{1,6})\s+(.+)$`)
	text = re.ReplaceAllStringFunc(strings.Join(paragraphs, "\n"), func(match string) string {
		matches := re.FindStringSubmatch(match)
		level := len(matches[1])
		return fmt.Sprintf("<h%d>%s</h%d>", level, matches[2], level)
	})
	
	// Convert bold
	text = regexp.MustCompile(`\*\*(.+?)\*\*`).ReplaceAllString(text, "<strong>$1</strong>")
	
	// Convert italic
	text = regexp.MustCompile(`\*(.+?)\*`).ReplaceAllString(text, "<em>$1</em>")
	
	// Convert links
	text = regexp.MustCompile(`\[(.+?)\]\((.+?)\)`).ReplaceAllString(text, "<a href=\"$2\">$1</a>")
	
	// Convert lists
	listItems := regexp.MustCompile(`(?m)^-\s+(.+)$`)
	text = listItems.ReplaceAllStringFunc(text, func(match string) string {
		item := listItems.FindStringSubmatch(match)
		return "<li>" + item[1] + "</li>"
	})
	
	// Wrap list items in ul tags
	if listItems.MatchString(text) {
		text = "<ul>\n" + text + "\n</ul>"
	}
	
	return fmt.Sprintf("<!DOCTYPE html>\n<html>\n<body>\n%s\n</body>\n</html>", text)
}

// formatAsMarkdown converts text to Markdown format or enhances existing Markdown.
func (f *FormattingDecorator) formatAsMarkdown(text string) string {
	// If text already has some markdown, we'll enhance it
	// Otherwise, add some basic markdown formatting
	
	// Add heading to the first line if it doesn't already have one
	if !strings.HasPrefix(text, "#") && len(text) > 0 {
		lines := strings.Split(text, "\n")
		lines[0] = "# " + lines[0]
		text = strings.Join(lines, "\n")
	}
	
	// Add emphasis to important words
	importantWords := []string{"important", "note", "warning", "caution", "danger"}
	for _, word := range importantWords {
		re := regexp.MustCompile(`(?i)\b` + word + `\b`)
		text = re.ReplaceAllString(text, "**"+word+"**")
	}
	
	// Convert simple bullet points
	text = regexp.MustCompile(`(?m)^(\*)\s+(.+)$`).ReplaceAllString(text, "- $2")
	
	// Add horizontal rules for better section separation
	paragraphs := strings.Split(text, "\n\n")
	result := []string{}
	
	for _, p := range paragraphs {
		result = append(result, p)
		if strings.HasPrefix(p, "##") {
			// Add a horizontal rule after section headings
			result = append(result, "---")
		}
	}
	
	return strings.Join(result, "\n\n")
}

// formatAsPlainText strips formatting and converts to plain text.
func (f *FormattingDecorator) formatAsPlainText(text string) string {
	// Remove HTML tags
	text = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(text, "")
	
	// Remove Markdown headings
	text = regexp.MustCompile(`(?m)^#{1,6}\s+(.+)$`).ReplaceAllString(text, "$1")
	
	// Remove Markdown bold and italic
	text = regexp.MustCompile(`\*\*(.+?)\*\*`).ReplaceAllString(text, "$1")
	text = regexp.MustCompile(`\*(.+?)\*`).ReplaceAllString(text, "$1")
	
	// Extract link text from Markdown links
	text = regexp.MustCompile(`\[(.+?)\]\(.+?\)`).ReplaceAllString(text, "$1")
	
	// Remove bullet points
	text = regexp.MustCompile(`(?m)^-\s+`).ReplaceAllString(text, "")
	
	// Remove horizontal rules
	text = regexp.MustCompile(`(?m)^---+$`).ReplaceAllString(text, "")
	
	// Remove multiple newlines
	text = regexp.MustCompile(`\n{3,}`).ReplaceAllString(text, "\n\n")
	
	return text
}

// HighlightingDecorator is a concrete decorator that highlights specific text.
type HighlightingDecorator struct {
	TextProcessorDecorator
	pattern        string
	highlightStart string
	highlightEnd   string
}

// NewHighlightingDecorator creates a decorator that highlights matching text.
func NewHighlightingDecorator(processor TextProcessor, pattern, startTag, endTag string) *HighlightingDecorator {
	return &HighlightingDecorator{
		TextProcessorDecorator: TextProcessorDecorator{
			wrapped:     processor,
			name:        "Text Highlighter",
			description: fmt.Sprintf("Highlights text matching pattern '%s'", pattern),
		},
		pattern:        pattern,
		highlightStart: startTag,
		highlightEnd:   endTag,
	}
}

// Process first processes the text using the wrapped processor,
// then highlights text matching the specified pattern.
func (h *HighlightingDecorator) Process(text string) (string, error) {
	// First, let the wrapped processor do its work
	processedText, err := h.wrapped.Process(text)
	if err != nil {
		return "", fmt.Errorf("error in wrapped processor: %w", err)
	}

	// Then apply highlighting
	re, err := regexp.Compile(h.pattern)
	if err != nil {
		return processedText, fmt.Errorf("invalid highlight pattern: %w", err)
	}

	return re.ReplaceAllString(processedText, h.highlightStart+"$0"+h.highlightEnd), nil
}

// IndentationDecorator is a concrete decorator that indents text.
type IndentationDecorator struct {
	TextProcessorDecorator
	indentation    string
	indentFirstLine bool
}

// NewIndentationDecorator creates a decorator that indents text.
func NewIndentationDecorator(processor TextProcessor, indentation string, indentFirstLine bool) *IndentationDecorator {
	return &IndentationDecorator{
		TextProcessorDecorator: TextProcessorDecorator{
			wrapped:     processor,
			name:        "Indentation Processor",
			description: fmt.Sprintf("Indents each line with '%s'", indentation),
		},
		indentation:    indentation,
		indentFirstLine: indentFirstLine,
	}
}

// Process first processes the text using the wrapped processor,
// then indents each line of the text.
func (i *IndentationDecorator) Process(text string) (string, error) {
	// First, let the wrapped processor do its work
	processedText, err := i.wrapped.Process(text)
	if err != nil {
		return "", fmt.Errorf("error in wrapped processor: %w", err)
	}

	// Split the text into lines
	lines := strings.Split(processedText, "\n")
	
	// Apply indentation to each line
	for j, line := range lines {
		if j == 0 && !i.indentFirstLine {
			// Skip first line if not indenting it
			continue
		}
		if line != "" {
			lines[j] = i.indentation + line
		}
	}
	
	return strings.Join(lines, "\n"), nil
}
