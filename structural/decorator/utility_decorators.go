package decorator

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"strings"
	"time"
)

// CompressionDecorator is a concrete decorator that compresses and decompresses text.
type CompressionDecorator struct {
	TextProcessorDecorator
	algorithm  string
	compress   bool
	encodeBase64 bool
}

// NewCompressionDecorator creates a decorator that compresses text.
func NewCompressionDecorator(processor TextProcessor, algorithm string, encodeBase64 bool) *CompressionDecorator {
	return &CompressionDecorator{
		TextProcessorDecorator: TextProcessorDecorator{
			wrapped:     processor,
			name:        "Compression Processor",
			description: fmt.Sprintf("Compresses text using %s algorithm", algorithm),
		},
		algorithm:  algorithm,
		compress:   true,
		encodeBase64: encodeBase64,
	}
}

// NewDecompressionDecorator creates a decorator that decompresses text.
func NewDecompressionDecorator(processor TextProcessor, algorithm string, decodeBase64 bool) *CompressionDecorator {
	return &CompressionDecorator{
		TextProcessorDecorator: TextProcessorDecorator{
			wrapped:     processor,
			name:        "Decompression Processor",
			description: fmt.Sprintf("Decompresses text using %s algorithm", algorithm),
		},
		algorithm:  algorithm,
		compress:   false,
		encodeBase64: decodeBase64,
	}
}

// Process first processes the text using the wrapped processor,
// then compresses or decompresses the text using the specified algorithm.
func (c *CompressionDecorator) Process(text string) (string, error) {
	// First, let the wrapped processor do its work
	processedText, err := c.wrapped.Process(text)
	if err != nil {
		return "", fmt.Errorf("error in wrapped processor: %w", err)
	}

	// Then apply compression or decompression
	if c.compress {
		return c.compressText(processedText)
	} else {
		return c.decompressText(processedText)
	}
}

// compressText compresses the text using the specified algorithm.
func (c *CompressionDecorator) compressText(text string) (string, error) {
	var buf bytes.Buffer
	var compressor io.WriteCloser
	var err error

	// Create the appropriate compressor
	switch c.algorithm {
	case "gzip":
		compressor = gzip.NewWriter(&buf)
	case "zlib":
		compressor = zlib.NewWriter(&buf)
	default:
		return text, fmt.Errorf("unsupported compression algorithm: %s", c.algorithm)
	}

	// Write the text to the compressor
	_, err = compressor.Write([]byte(text))
	if err != nil {
		return "", fmt.Errorf("compression write error: %w", err)
	}

	// Close the compressor to flush any pending data
	err = compressor.Close()
	if err != nil {
		return "", fmt.Errorf("compression close error: %w", err)
	}

	// If configured, encode the compressed data as base64
	if c.encodeBase64 {
		return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
	}
	
	// Otherwise return as raw bytes (which might contain non-printable characters)
	return string(buf.Bytes()), nil
}

// decompressText decompresses the text using the specified algorithm.
func (c *CompressionDecorator) decompressText(text string) (string, error) {
	var data []byte
	var err error

	// If the data is base64 encoded, decode it first
	if c.encodeBase64 {
		data, err = base64.StdEncoding.DecodeString(text)
		if err != nil {
			return "", fmt.Errorf("base64 decode error: %w", err)
		}
	} else {
		data = []byte(text)
	}

	// Create a bytes reader
	buf := bytes.NewReader(data)
	var decompressor io.ReadCloser

	// Create the appropriate decompressor
	switch c.algorithm {
	case "gzip":
		decompressor, err = gzip.NewReader(buf)
		if err != nil {
			return "", fmt.Errorf("gzip reader creation error: %w", err)
		}
	case "zlib":
		decompressor, err = zlib.NewReader(buf)
		if err != nil {
			return "", fmt.Errorf("zlib reader creation error: %w", err)
		}
	default:
		return text, fmt.Errorf("unsupported decompression algorithm: %s", c.algorithm)
	}

	// Read the decompressed data
	decompressed, err := io.ReadAll(decompressor)
	if err != nil {
		return "", fmt.Errorf("decompression read error: %w", err)
	}

	// Close the decompressor
	err = decompressor.Close()
	if err != nil {
		return "", fmt.Errorf("decompression close error: %w", err)
	}

	return string(decompressed), nil
}

// LoggingDecorator is a concrete decorator that logs text processing.
type LoggingDecorator struct {
	TextProcessorDecorator
	logInput       bool
	logOutput      bool
	logTiming      bool
	maxContentLength int
}

// NewLoggingDecorator creates a decorator that logs the processing details.
func NewLoggingDecorator(processor TextProcessor, logInput, logOutput, logTiming bool, maxContentLength int) *LoggingDecorator {
	return &LoggingDecorator{
		TextProcessorDecorator: TextProcessorDecorator{
			wrapped:     processor,
			name:        "Logging Processor",
			description: "Logs details about text processing",
		},
		logInput:       logInput,
		logOutput:      logOutput,
		logTiming:      logTiming,
		maxContentLength: maxContentLength,
	}
}

// Process logs details about the text processing and then delegates
// to the wrapped processor.
func (l *LoggingDecorator) Process(text string) (string, error) {
	// Truncate the text for logging if necessary
	truncatedInput := text
	if l.maxContentLength > 0 && len(text) > l.maxContentLength {
		truncatedInput = text[:l.maxContentLength] + "... (truncated)"
	}

	// Log the input if configured to do so
	if l.logInput {
		log.Printf("Processing text with %s\nInput: %s", l.GetProcessingChain(), truncatedInput)
	} else {
		log.Printf("Processing text with %s", l.GetProcessingChain())
	}

	// Start timing if configured to do so
	var startTime time.Time
	if l.logTiming {
		startTime = time.Now()
	}

	// Process the text
	result, err := l.wrapped.Process(text)

	// Log any error
	if err != nil {
		log.Printf("Error processing text: %v", err)
		return "", err
	}

	// Log the timing if configured to do so
	if l.logTiming {
		elapsed := time.Since(startTime)
		log.Printf("Processing took %v", elapsed)
	}

	// Log the output if configured to do so
	if l.logOutput {
		truncatedOutput := result
		if l.maxContentLength > 0 && len(result) > l.maxContentLength {
			truncatedOutput = result[:l.maxContentLength] + "... (truncated)"
		}
		log.Printf("Output: %s", truncatedOutput)
	}

	return result, nil
}

// MetadataDecorator is a concrete decorator that adds metadata to the text.
type MetadataDecorator struct {
	TextProcessorDecorator
	metadata map[string]string
	position string // "prefix" or "suffix"
}

// NewMetadataDecorator creates a decorator that adds metadata to the text.
func NewMetadataDecorator(processor TextProcessor, metadata map[string]string, position string) *MetadataDecorator {
	return &MetadataDecorator{
		TextProcessorDecorator: TextProcessorDecorator{
			wrapped:     processor,
			name:        "Metadata Processor",
			description: "Adds metadata to the text",
		},
		metadata: metadata,
		position: strings.ToLower(position),
	}
}

// Process first processes the text using the wrapped processor,
// then adds metadata to the text.
func (m *MetadataDecorator) Process(text string) (string, error) {
	// First, let the wrapped processor do its work
	processedText, err := m.wrapped.Process(text)
	if err != nil {
		return "", fmt.Errorf("error in wrapped processor: %w", err)
	}

	// Generate the metadata section
	var metadataSection strings.Builder
	metadataSection.WriteString("--- Metadata ---\n")
	for key, value := range m.metadata {
		metadataSection.WriteString(fmt.Sprintf("%s: %s\n", key, value))
	}
	metadataSection.WriteString("---------------\n")

	// Add the metadata based on the position
	if m.position == "prefix" {
		return metadataSection.String() + processedText, nil
	} else {
		return processedText + "\n" + metadataSection.String(), nil
	}
}

// LanguageTranslationDecorator is a concrete decorator that simulates language translation.
type LanguageTranslationDecorator struct {
	TextProcessorDecorator
	sourceLanguage string
	targetLanguage string
	translations  map[string]map[string]string
}

// NewLanguageTranslationDecorator creates a decorator that simulates translating text.
// Note: This is a simple simulation for demonstration purposes.
func NewLanguageTranslationDecorator(processor TextProcessor, sourceLanguage, targetLanguage string) *LanguageTranslationDecorator {
	// Initialize with some common words and phrases for demonstration
	translations := map[string]map[string]string{
		"en": {
			"fr": map[string]string{
				"hello": "bonjour",
				"world": "monde",
				"welcome": "bienvenue",
				"thank you": "merci",
				"please": "s'il vous plaît",
				"goodbye": "au revoir",
			},
			"es": map[string]string{
				"hello": "hola",
				"world": "mundo",
				"welcome": "bienvenido",
				"thank you": "gracias",
				"please": "por favor",
				"goodbye": "adiós",
			},
		},
	}

	return &LanguageTranslationDecorator{
		TextProcessorDecorator: TextProcessorDecorator{
			wrapped:     processor,
			name:        "Translation Processor",
			description: fmt.Sprintf("Translates text from %s to %s", sourceLanguage, targetLanguage),
		},
		sourceLanguage: sourceLanguage,
		targetLanguage: targetLanguage,
		translations:  translations,
	}
}

// Process first processes the text using the wrapped processor,
// then translates the text from the source language to the target language.
func (t *LanguageTranslationDecorator) Process(text string) (string, error) {
	// First, let the wrapped processor do its work
	processedText, err := t.wrapped.Process(text)
	if err != nil {
		return "", fmt.Errorf("error in wrapped processor: %w", err)
	}

	// Translate the text
	// Note: This is a simple demonstration; in a real app you would use a translation API
	if t.sourceLanguage == t.targetLanguage {
		return processedText, nil
	}

	// Check if we have translations for the source language
	sourceLang, ok := t.translations[t.sourceLanguage]
	if !ok {
		return processedText, fmt.Errorf("translation from %s not supported", t.sourceLanguage)
	}

	// Check if we have translations for the target language
	targetDict, ok := sourceLang[t.targetLanguage]
	if !ok {
		return processedText, fmt.Errorf("translation to %s not supported", t.targetLanguage)
	}

	// Simple word-by-word translation for demonstration
	words := strings.Fields(processedText)
	for i, word := range words {
		// Convert to lowercase for dictionary lookup
		lowerWord := strings.ToLower(word)
		
		// Strip any trailing punctuation
		trimmedWord := lowerWord
		var punctuation string
		for j := len(trimmedWord) - 1; j >= 0; j-- {
			if !strings.ContainsRune(",.!?:;", rune(trimmedWord[j])) {
				break
			}
			punctuation = string(trimmedWord[j]) + punctuation
			trimmedWord = trimmedWord[:j]
		}
		
		// Check if we have a translation for this word
		if translation, ok := targetDict[trimmedWord]; ok {
			// Preserve capitalization
			if word[0] >= 'A' && word[0] <= 'Z' {
				translation = strings.ToUpper(translation[:1]) + translation[1:]
			}
			// Reattach punctuation
			words[i] = translation + punctuation
		}
	}

	return strings.Join(words, " "), nil
}
