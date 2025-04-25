package flyweight

import (
	"encoding/json"
	"fmt"
)

// SerializedCharacter represents a Character for serialization.
type SerializedCharacter struct {
	Value    rune   `json:"value"`
	Row      int    `json:"row"`
	Column   int    `json:"column"`
	FormatID string `json:"formatId"`
}

// SerializedParagraph represents a Paragraph for serialization.
type SerializedParagraph struct {
	Content string `json:"content"`
	StyleID string `json:"styleId"`
	Row     int    `json:"row"`
}

// SerializedFormat represents the TextFormat for serialization.
type SerializedFormat struct {
	ID            string  `json:"id"`
	FontFamily    string  `json:"fontFamily"`
	FontSize      int     `json:"fontSize"`
	FontColor     string  `json:"fontColor"`
	IsBold        bool    `json:"isBold"`
	IsItalic      bool    `json:"isItalic"`
	IsUnderline   bool    `json:"isUnderline"`
	Background    string  `json:"background"`
	Alignment     string  `json:"alignment"`
	LetterSpacing float64 `json:"letterSpacing"`
	LineHeight    float64 `json:"lineHeight"`
}

// SerializedParagraphStyle represents the ParagraphStyle for serialization.
type SerializedParagraphStyle struct {
	ID              string  `json:"id"`
	Alignment       string  `json:"alignment"`
	LineSpacing     float64 `json:"lineSpacing"`
	BeforeSpacing   float64 `json:"beforeSpacing"`
	AfterSpacing    float64 `json:"afterSpacing"`
	FirstLineIndent float64 `json:"firstLineIndent"`
	LeftMargin      float64 `json:"leftMargin"`
	RightMargin     float64 `json:"rightMargin"`
	BorderStyle     string  `json:"borderStyle"`
	BorderColor     string  `json:"borderColor"`
	BackgroundColor string  `json:"backgroundColor"`
}

// SerializedDocument represents a Document for serialization.
type SerializedDocument struct {
	Name            string                   `json:"name"`
	Characters      []SerializedCharacter    `json:"characters"`
	Paragraphs      []SerializedParagraph    `json:"paragraphs"`
	Formats         []SerializedFormat       `json:"formats"`
	ParagraphStyles []SerializedParagraphStyle `json:"paragraphStyles"`
}

// Serialize converts a FormattedDocument to a JSON string.
func (d *FormattedDocument) Serialize() (string, error) {
	serialized := SerializedDocument{
		Name:            d.Name,
		Characters:      make([]SerializedCharacter, len(d.Characters)),
		Paragraphs:      make([]SerializedParagraph, len(d.Paragraphs)),
		Formats:         []SerializedFormat{},
		ParagraphStyles: []SerializedParagraphStyle{},
	}

	// Serialize characters
	for i, char := range d.Characters {
		serialized.Characters[i] = SerializedCharacter{
			Value:    char.Value,
			Row:      char.Row,
			Column:   char.Column,
			FormatID: char.FormatID,
		}
	}

	// Serialize paragraphs
	for i, para := range d.Paragraphs {
		serialized.Paragraphs[i] = SerializedParagraph{
			Content: para.Content,
			StyleID: para.StyleID,
			Row:     para.Row,
		}
	}

	// Serialize formats
	formatIDs := make(map[string]bool)
	for _, char := range d.Characters {
		formatIDs[char.FormatID] = true
	}

	for formatID := range formatIDs {
		format := d.formatFactory.GetFormatByID(formatID)
		if format == nil {
			continue
		}

		formatProps := format.GetFormat()
		serialized.Formats = append(serialized.Formats, SerializedFormat{
			ID:            formatID,
			FontFamily:    formatProps["fontFamily"].(string),
			FontSize:      formatProps["fontSize"].(int),
			FontColor:     formatProps["fontColor"].(string),
			IsBold:        formatProps["isBold"].(bool),
			IsItalic:      formatProps["isItalic"].(bool),
			IsUnderline:   formatProps["isUnderline"].(bool),
			Background:    formatProps["background"].(string),
			Alignment:     formatProps["alignment"].(string),
			LetterSpacing: formatProps["letterSpacing"].(float64),
			LineHeight:    formatProps["lineHeight"].(float64),
		})
	}

	// Serialize paragraph styles
	styleIDs := make(map[string]bool)
	for _, para := range d.Paragraphs {
		styleIDs[para.StyleID] = true
	}

	for styleID := range styleIDs {
		style := d.paragraphFactory.GetStyleByID(styleID)
		if style == nil {
			continue
		}

		styleProps := style.GetStyle()
		serialized.ParagraphStyles = append(serialized.ParagraphStyles, SerializedParagraphStyle{
			ID:              styleID,
			Alignment:       styleProps["alignment"].(string),
			LineSpacing:     styleProps["lineSpacing"].(float64),
			BeforeSpacing:   styleProps["beforeSpacing"].(float64),
			AfterSpacing:    styleProps["afterSpacing"].(float64),
			FirstLineIndent: styleProps["firstLineIndent"].(float64),
			LeftMargin:      styleProps["leftMargin"].(float64),
			RightMargin:     styleProps["rightMargin"].(float64),
			BorderStyle:     styleProps["borderStyle"].(string),
			BorderColor:     styleProps["borderColor"].(string),
			BackgroundColor: styleProps["backgroundColor"].(string),
		})
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(serialized, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to serialize document: %w", err)
	}

	return string(jsonData), nil
}

// DeserializeDocument creates a FormattedDocument from a JSON string.
func DeserializeDocument(jsonData string, formatFactory *TextFormatFactory, paragraphFactory *ParagraphStyleFactory) (*FormattedDocument, error) {
	var serialized SerializedDocument
	err := json.Unmarshal([]byte(jsonData), &serialized)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize document: %w", err)
	}

	// Create a new document
	doc := NewFormattedDocument(serialized.Name, formatFactory, paragraphFactory)

	// Recreate all formats first
	for _, format := range serialized.Formats {
		formatFactory.GetTextFormat(
			format.FontFamily,
			format.FontSize,
			format.FontColor,
			format.IsBold,
			format.IsItalic,
			format.IsUnderline,
			format.Background,
			format.Alignment,
			format.LetterSpacing,
			format.LineHeight,
		)
	}

	// Recreate all paragraph styles
	for _, style := range serialized.ParagraphStyles {
		paragraphFactory.GetParagraphStyle(
			style.Alignment,
			style.LineSpacing,
			style.BeforeSpacing,
			style.AfterSpacing,
			style.FirstLineIndent,
			style.LeftMargin,
			style.RightMargin,
			style.BorderStyle,
			style.BorderColor,
			style.BackgroundColor,
		)
	}

	// Recreate characters
	for _, char := range serialized.Characters {
		doc.Characters = append(doc.Characters, &Character{
			Value:    char.Value,
			Row:      char.Row,
			Column:   char.Column,
			FormatID: char.FormatID,
		})
	}

	// Recreate paragraphs
	for _, para := range serialized.Paragraphs {
		doc.Paragraphs = append(doc.Paragraphs, &Paragraph{
			Content: para.Content,
			StyleID: para.StyleID,
			Row:     para.Row,
		})
	}

	return doc, nil
}
