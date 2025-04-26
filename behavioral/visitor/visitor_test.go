package visitor

import (
	"reflect"
	"strings"
	"testing"
)

// TestTextElementVisitorAcceptance tests the TextElement's Accept method
func TestTextElementVisitorAcceptance(t *testing.T) {
	element := &TextElement{Content: "Sample text"}
	visitor := NewHTMLExportVisitor()
	
	err := element.Accept(visitor)
	if err != nil {
		t.Errorf("TextElement.Accept returned error: %v", err)
	}
	
	expected := "<p>Sample text</p>\n"
	if visitor.Output.String() != expected {
		t.Errorf("TextElement.Accept produced incorrect HTML.\nExpected: %q\nGot: %q", expected, visitor.Output.String())
	}
}

// TestImageElementVisitorAcceptance tests the ImageElement's Accept method
func TestImageElementVisitorAcceptance(t *testing.T) {
	element := &ImageElement{
		Source: "image.jpg",
		Alt:    "Sample Image",
		Width:  100,
		Height: 80,
	}
	visitor := NewHTMLExportVisitor()
	
	err := element.Accept(visitor)
	if err != nil {
		t.Errorf("ImageElement.Accept returned error: %v", err)
	}
	
	expected := "<img src=\"image.jpg\" alt=\"Sample Image\" width=\"100\" height=\"80\">\n"
	if visitor.Output.String() != expected {
		t.Errorf("ImageElement.Accept produced incorrect HTML.\nExpected: %q\nGot: %q", expected, visitor.Output.String())
	}
}

// TestLinkElementVisitorAcceptance tests the LinkElement's Accept method
func TestLinkElementVisitorAcceptance(t *testing.T) {
	element := &LinkElement{
		URL:   "https://example.com",
		Text:  "Example",
		Title: "Visit Example",
	}
	visitor := NewHTMLExportVisitor()
	
	err := element.Accept(visitor)
	if err != nil {
		t.Errorf("LinkElement.Accept returned error: %v", err)
	}
	
	expected := "<a href=\"https://example.com\" title=\"Visit Example\">Example</a>\n"
	if visitor.Output.String() != expected {
		t.Errorf("LinkElement.Accept produced incorrect HTML.\nExpected: %q\nGot: %q", expected, visitor.Output.String())
	}
}

// TestTableElementVisitorAcceptance tests the TableElement's Accept method
func TestTableElementVisitorAcceptance(t *testing.T) {
	element := &TableElement{
		Rows:    2,
		Columns: 2,
		Data: [][]string{
			{"Header 1", "Header 2"},
			{"Cell 1", "Cell 2"},
		},
	}
	visitor := NewHTMLExportVisitor()
	
	err := element.Accept(visitor)
	if err != nil {
		t.Errorf("TableElement.Accept returned error: %v", err)
	}
	
	expected := "<table>\n  <tr>\n    <td>Header 1</td>\n    <td>Header 2</td>\n  </tr>\n  <tr>\n    <td>Cell 1</td>\n    <td>Cell 2</td>\n  </tr>\n</table>\n"
	if visitor.Output.String() != expected {
		t.Errorf("TableElement.Accept produced incorrect HTML.\nExpected: %q\nGot: %q", expected, visitor.Output.String())
	}
}

// TestCompositeElementVisitorAcceptance tests the CompositeElement's Accept method
func TestCompositeElementVisitorAcceptance(t *testing.T) {
	composite := &CompositeElement{
		Name: "Section",
		Children: []Element{
			&TextElement{Content: "Hello"},
		},
	}
	visitor := NewHTMLExportVisitor()
	
	err := composite.Accept(visitor)
	if err != nil {
		t.Errorf("CompositeElement.Accept returned error: %v", err)
	}
	
	expected := "<div class=\"section\">\n<p>Hello</p>\n"
	if visitor.Output.String() != expected {
		t.Errorf("CompositeElement.Accept produced incorrect HTML.\nExpected: %q\nGot: %q", expected, visitor.Output.String())
	}
}

// TestComplexDocument tests a complex document structure with nested elements
func TestComplexDocument(t *testing.T) {
	// Create a complex document
	document := &CompositeElement{
		Name: "Document",
		Children: []Element{
			&CompositeElement{
				Name: "Header",
				Children: []Element{
					&TextElement{Content: "Title"},
					&ImageElement{Source: "logo.png", Alt: "Logo", Width: 50, Height: 30},
				},
			},
			&CompositeElement{
				Name: "Body",
				Children: []Element{
					&TextElement{Content: "Lorem ipsum dolor sit amet"},
					&LinkElement{URL: "https://example.com", Text: "Example", Title: ""},
					&TableElement{
						Rows:    2,
						Columns: 2,
						Data: [][]string{
							{"Name", "Value"},
							{"Item", "10"},
						},
					},
				},
			},
		},
	}
	
	// Test HTML export
	htmlVisitor := NewHTMLExportVisitor()
	err := document.Accept(htmlVisitor)
	if err != nil {
		t.Errorf("Document.Accept with HTMLExportVisitor returned error: %v", err)
	}
	
	html := htmlVisitor.GetHTML()
	if !strings.Contains(html, "<div class=\"document\">") ||
	   !strings.Contains(html, "<div class=\"header\">") ||
	   !strings.Contains(html, "<div class=\"body\">") {
		t.Error("HTML output doesn't contain expected structure")
	}
	
	// Test Markdown export
	mdVisitor := NewMarkdownExportVisitor()
	err = document.Accept(mdVisitor)
	if err != nil {
		t.Errorf("Document.Accept with MarkdownExportVisitor returned error: %v", err)
	}
	
	markdown := mdVisitor.GetMarkdown()
	if !strings.Contains(markdown, "# Document") ||
	   !strings.Contains(markdown, "## Header") ||
	   !strings.Contains(markdown, "## Body") {
		t.Error("Markdown output doesn't contain expected structure")
	}
	
	// Test Statistics visitor
	statsVisitor := NewStatisticsVisitor()
	err = document.Accept(statsVisitor)
	if err != nil {
		t.Errorf("Document.Accept with StatisticsVisitor returned error: %v", err)
	}
	
	if statsVisitor.TextCount != 2 ||
	   statsVisitor.ImageCount != 1 ||
	   statsVisitor.LinkCount != 1 ||
	   statsVisitor.TableCount != 1 {
		t.Errorf("StatisticsVisitor count is incorrect. Got: %+v", statsVisitor)
	}
}

// TestHTMLExportVisitor tests the HTML export visitor on all element types
func TestHTMLExportVisitor(t *testing.T) {
	visitor := NewHTMLExportVisitor()
	
	// Test visit text element
	textElement := &TextElement{Content: "Sample text"}
	err := visitor.VisitText(textElement)
	if err != nil {
		t.Errorf("HTMLExportVisitor.VisitText returned error: %v", err)
	}
	
	// Test visit image element
	imageElement := &ImageElement{Source: "image.jpg", Alt: "Image", Width: 100, Height: 80}
	err = visitor.VisitImage(imageElement)
	if err != nil {
		t.Errorf("HTMLExportVisitor.VisitImage returned error: %v", err)
	}
	
	// Test visit link element
	linkElement := &LinkElement{URL: "https://example.com", Text: "Link", Title: ""}
	err = visitor.VisitLink(linkElement)
	if err != nil {
		t.Errorf("HTMLExportVisitor.VisitLink returned error: %v", err)
	}
	
	// Test visit table element
	tableElement := &TableElement{
		Rows:    2,
		Columns: 2,
		Data: [][]string{
			{"H1", "H2"},
			{"D1", "D2"},
		},
	}
	err = visitor.VisitTable(tableElement)
	if err != nil {
		t.Errorf("HTMLExportVisitor.VisitTable returned error: %v", err)
	}
	
	// Test visit composite element
	compositeElement := &CompositeElement{Name: "Test"}
	err = visitor.VisitComposite(compositeElement)
	if err != nil {
		t.Errorf("HTMLExportVisitor.VisitComposite returned error: %v", err)
	}
	
	// Verify output contains all expected elements
	output := visitor.GetHTML()
	if !strings.Contains(output, "<p>Sample text</p>") ||
	   !strings.Contains(output, "<img src=\"image.jpg\"") ||
	   !strings.Contains(output, "<a href=\"https://example.com\">Link</a>") ||
	   !strings.Contains(output, "<table>") ||
	   !strings.Contains(output, "<div class=\"test\">") {
		t.Errorf("HTMLExportVisitor output missing expected elements: %s", output)
	}
}

// TestMarkdownExportVisitor tests the Markdown export visitor on all element types
func TestMarkdownExportVisitor(t *testing.T) {
	visitor := NewMarkdownExportVisitor()
	
	// Test visit text element
	textElement := &TextElement{Content: "Sample text"}
	err := visitor.VisitText(textElement)
	if err != nil {
		t.Errorf("MarkdownExportVisitor.VisitText returned error: %v", err)
	}
	
	// Test visit image element
	imageElement := &ImageElement{Source: "image.jpg", Alt: "Image"}
	err = visitor.VisitImage(imageElement)
	if err != nil {
		t.Errorf("MarkdownExportVisitor.VisitImage returned error: %v", err)
	}
	
	// Test visit link element
	linkElement := &LinkElement{URL: "https://example.com", Text: "Link"}
	err = visitor.VisitLink(linkElement)
	if err != nil {
		t.Errorf("MarkdownExportVisitor.VisitLink returned error: %v", err)
	}
	
	// Test visit table element
	tableElement := &TableElement{
		Rows:    2,
		Columns: 2,
		Data: [][]string{
			{"H1", "H2"},
			{"D1", "D2"},
		},
	}
	err = visitor.VisitTable(tableElement)
	if err != nil {
		t.Errorf("MarkdownExportVisitor.VisitTable returned error: %v", err)
	}
	
	// Test visit composite element
	compositeElement := &CompositeElement{Name: "Test"}
	err = visitor.VisitComposite(compositeElement)
	if err != nil {
		t.Errorf("MarkdownExportVisitor.VisitComposite returned error: %v", err)
	}
	
	// Verify output contains all expected elements
	output := visitor.GetMarkdown()
	if !strings.Contains(output, "Sample text") ||
	   !strings.Contains(output, "![Image](image.jpg)") ||
	   !strings.Contains(output, "[Link](https://example.com)") ||
	   !strings.Contains(output, "| H1 | H2 |") ||
	   !strings.Contains(output, "# Test") {
		t.Errorf("MarkdownExportVisitor output missing expected elements: %s", output)
	}
}

// TestPlainTextExportVisitor tests the Plain Text export visitor on all element types
func TestPlainTextExportVisitor(t *testing.T) {
	visitor := NewPlainTextExportVisitor()
	
	// Test visit text element
	textElement := &TextElement{Content: "Sample text"}
	err := visitor.VisitText(textElement)
	if err != nil {
		t.Errorf("PlainTextExportVisitor.VisitText returned error: %v", err)
	}
	
	// Test visit image element
	imageElement := &ImageElement{Source: "image.jpg", Alt: "Image"}
	err = visitor.VisitImage(imageElement)
	if err != nil {
		t.Errorf("PlainTextExportVisitor.VisitImage returned error: %v", err)
	}
	
	// Test visit link element
	linkElement := &LinkElement{URL: "https://example.com", Text: "Link"}
	err = visitor.VisitLink(linkElement)
	if err != nil {
		t.Errorf("PlainTextExportVisitor.VisitLink returned error: %v", err)
	}
	
	// Verify output contains all expected elements
	output := visitor.GetText()
	if !strings.Contains(output, "Sample text") ||
	   !strings.Contains(output, "[Image: Image]") ||
	   !strings.Contains(output, "Link (https://example.com)") {
		t.Errorf("PlainTextExportVisitor output missing expected elements: %s", output)
	}
}

// TestStatisticsVisitor tests the Statistics visitor on all element types
func TestStatisticsVisitor(t *testing.T) {
	visitor := NewStatisticsVisitor()
	
	// Test visit text element
	textElement := &TextElement{Content: "One two three"}
	err := visitor.VisitText(textElement)
	if err != nil {
		t.Errorf("StatisticsVisitor.VisitText returned error: %v", err)
	}
	
	// Test visit image element
	imageElement := &ImageElement{Source: "image.jpg", Alt: "Image"}
	err = visitor.VisitImage(imageElement)
	if err != nil {
		t.Errorf("StatisticsVisitor.VisitImage returned error: %v", err)
	}
	
	// Test visit link element
	linkElement := &LinkElement{URL: "https://example.com", Text: "Example Link"}
	err = visitor.VisitLink(linkElement)
	if err != nil {
		t.Errorf("StatisticsVisitor.VisitLink returned error: %v", err)
	}
	
	// Test visit table element
	tableElement := &TableElement{
		Rows:    2,
		Columns: 2,
		Data: [][]string{
			{"Header1", "Header2"},
			{"Data1", "Data2"},
		},
	}
	err = visitor.VisitTable(tableElement)
	if err != nil {
		t.Errorf("StatisticsVisitor.VisitTable returned error: %v", err)
	}
	
	// Verify statistics are correct
	if visitor.TextCount != 1 {
		t.Errorf("StatisticsVisitor.TextCount is incorrect. Expected 1, got %d", visitor.TextCount)
	}
	if visitor.ImageCount != 1 {
		t.Errorf("StatisticsVisitor.ImageCount is incorrect. Expected 1, got %d", visitor.ImageCount)
	}
	if visitor.LinkCount != 1 {
		t.Errorf("StatisticsVisitor.LinkCount is incorrect. Expected 1, got %d", visitor.LinkCount)
	}
	if visitor.TableCount != 1 {
		t.Errorf("StatisticsVisitor.TableCount is incorrect. Expected 1, got %d", visitor.TableCount)
	}
	if visitor.WordCount != 7 { // "One two three" + "Example Link" + 4 words in table
		t.Errorf("StatisticsVisitor.WordCount is incorrect. Expected 7, got %d", visitor.WordCount)
	}
}

// TestSpellCheckVisitor tests the SpellCheck visitor on all element types
func TestSpellCheckVisitor(t *testing.T) {
	visitor := NewSpellCheckVisitor()
	
	// Test visit text element with valid words
	textElement := &TextElement{Content: "the and to with"}
	err := visitor.VisitText(textElement)
	if err != nil {
		t.Errorf("SpellCheckVisitor.VisitText returned error: %v", err)
	}
	if len(visitor.Errors) > 0 {
		t.Errorf("SpellCheckVisitor found errors in valid text: %v", visitor.Errors)
	}
	
	// Reset errors
	visitor.Errors = []string{}
	
	// Test visit text element with invalid words
	textElement = &TextElement{Content: "the and xyzzyx flubbergasted"}
	err = visitor.VisitText(textElement)
	if err != nil {
		t.Errorf("SpellCheckVisitor.VisitText returned error: %v", err)
	}
	if len(visitor.Errors) != 2 {
		t.Errorf("SpellCheckVisitor should have found 2 errors, found %d: %v", len(visitor.Errors), visitor.Errors)
	}
}

// TestNilElementHandling tests that visitors properly handle nil elements
func TestNilElementHandling(t *testing.T) {
	visitor := NewHTMLExportVisitor()
	
	// Test visit nil text element
	err := visitor.VisitText(nil)
	if err == nil {
		t.Error("HTMLExportVisitor.VisitText did not return error for nil element")
	}
	
	// Test visit nil image element
	err = visitor.VisitImage(nil)
	if err == nil {
		t.Error("HTMLExportVisitor.VisitImage did not return error for nil element")
	}
	
	// Test visit nil link element
	err = visitor.VisitLink(nil)
	if err == nil {
		t.Error("HTMLExportVisitor.VisitLink did not return error for nil element")
	}
	
	// Test visit nil table element
	err = visitor.VisitTable(nil)
	if err == nil {
		t.Error("HTMLExportVisitor.VisitTable did not return error for nil element")
	}
	
	// Test visit nil composite element
	err = visitor.VisitComposite(nil)
	if err == nil {
		t.Error("HTMLExportVisitor.VisitComposite did not return error for nil element")
	}
}

// TestEmptyCompositeElement tests handling of a composite with no children
func TestEmptyCompositeElement(t *testing.T) {
	composite := &CompositeElement{
		Name:     "Empty",
		Children: []Element{},
	}
	
	visitor := NewHTMLExportVisitor()
	err := composite.Accept(visitor)
	if err != nil {
		t.Errorf("Empty CompositeElement.Accept returned error: %v", err)
	}
	
	expected := "<div class=\"empty\">\n"
	if visitor.Output.String() != expected {
		t.Errorf("Empty CompositeElement.Accept produced incorrect HTML.\nExpected: %q\nGot: %q", expected, visitor.Output.String())
	}
}
