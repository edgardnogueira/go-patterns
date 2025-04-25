package decorator

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

// EncryptionDecorator is a concrete decorator that encrypts and decrypts text.
type EncryptionDecorator struct {
	TextProcessorDecorator
	key         []byte
	encrypt     bool
	encryptMode string
}

// NewEncryptionDecorator creates a decorator that encrypts text.
func NewEncryptionDecorator(processor TextProcessor, key string, mode string) *EncryptionDecorator {
	// Create a fixed size key using MD5 (for simplicity - not for production use)
	hasher := md5.New()
	hasher.Write([]byte(key))
	keyBytes := hasher.Sum(nil)

	return &EncryptionDecorator{
		TextProcessorDecorator: TextProcessorDecorator{
			wrapped:     processor,
			name:        "Encryption Processor",
			description: fmt.Sprintf("Encrypts text using %s mode", mode),
		},
		key:         keyBytes,
		encrypt:     true,
		encryptMode: mode,
	}
}

// NewDecryptionDecorator creates a decorator that decrypts text.
func NewDecryptionDecorator(processor TextProcessor, key string, mode string) *EncryptionDecorator {
	// Create a fixed size key using MD5 (for simplicity - not for production use)
	hasher := md5.New()
	hasher.Write([]byte(key))
	keyBytes := hasher.Sum(nil)

	return &EncryptionDecorator{
		TextProcessorDecorator: TextProcessorDecorator{
			wrapped:     processor,
			name:        "Decryption Processor",
			description: fmt.Sprintf("Decrypts text using %s mode", mode),
		},
		key:         keyBytes,
		encrypt:     false,
		encryptMode: mode,
	}
}

// Process first processes the text using the wrapped processor,
// then encrypts or decrypts the text using the specified mode.
func (e *EncryptionDecorator) Process(text string) (string, error) {
	// First, let the wrapped processor do its work
	processedText, err := e.wrapped.Process(text)
	if err != nil {
		return "", fmt.Errorf("error in wrapped processor: %w", err)
	}

	// Then apply encryption or decryption
	if e.encrypt {
		switch e.encryptMode {
		case "aes":
			return e.encryptAES(processedText)
		case "base64":
			return e.encryptBase64(processedText), nil
		case "rot13":
			return e.encryptRot13(processedText), nil
		default:
			return processedText, fmt.Errorf("unsupported encryption mode: %s", e.encryptMode)
		}
	} else {
		switch e.encryptMode {
		case "aes":
			return e.decryptAES(processedText)
		case "base64":
			return e.decryptBase64(processedText)
		case "rot13":
			return e.decryptRot13(processedText), nil
		default:
			return processedText, fmt.Errorf("unsupported decryption mode: %s", e.encryptMode)
		}
	}
}

// encryptAES encrypts text using AES-256 in CFB mode.
func (e *EncryptionDecorator) encryptAES(text string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	// Create the initialization vector
	ciphertext := make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// Encrypt
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(text))

	// Return as base64 encoded string
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decryptAES decrypts text using AES-256 in CFB mode.
func (e *EncryptionDecorator) decryptAES(text string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	// Check the ciphertext length
	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	// Get the initialization vector
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// Decrypt
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

// encryptBase64 encodes text using base64.
func (e *EncryptionDecorator) encryptBase64(text string) string {
	return base64.StdEncoding.EncodeToString([]byte(text))
}

// decryptBase64 decodes text using base64.
func (e *EncryptionDecorator) decryptBase64(text string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// encryptRot13 encodes text using ROT13 (Caesar cipher with shift 13).
func (e *EncryptionDecorator) encryptRot13(text string) string {
	return strings.Map(rot13, text)
}

// decryptRot13 decodes text using ROT13 (applying ROT13 again).
func (e *EncryptionDecorator) decryptRot13(text string) string {
	return strings.Map(rot13, text)
}

// rot13 applies ROT13 transformation to a character.
func rot13(r rune) rune {
	switch {
	case r >= 'a' && r <= 'z':
		return 'a' + (r-'a'+13)%26
	case r >= 'A' && r <= 'Z':
		return 'A' + (r-'A'+13)%26
	default:
		return r
	}
}

// HashingDecorator is a concrete decorator that adds a hash to text.
type HashingDecorator struct {
	TextProcessorDecorator
	algorithm string
	appendHash bool
}

// NewHashingDecorator creates a decorator that hashes text.
func NewHashingDecorator(processor TextProcessor, algorithm string, appendHash bool) *HashingDecorator {
	return &HashingDecorator{
		TextProcessorDecorator: TextProcessorDecorator{
			wrapped:     processor,
			name:        "Hashing Processor",
			description: fmt.Sprintf("Hashes text using %s algorithm", algorithm),
		},
		algorithm: algorithm,
		appendHash: appendHash,
	}
}

// Process first processes the text using the wrapped processor,
// then adds a hash of the text using the specified algorithm.
func (h *HashingDecorator) Process(text string) (string, error) {
	// First, let the wrapped processor do its work
	processedText, err := h.wrapped.Process(text)
	if err != nil {
		return "", fmt.Errorf("error in wrapped processor: %w", err)
	}

	// Calculate the hash
	var hash string
	switch h.algorithm {
	case "md5":
		hasher := md5.New()
		hasher.Write([]byte(processedText))
		hash = hex.EncodeToString(hasher.Sum(nil))
	default:
		return processedText, fmt.Errorf("unsupported hash algorithm: %s", h.algorithm)
	}

	// Either append the hash or return it
	if h.appendHash {
		return fmt.Sprintf("%s [Hash: %s]", processedText, hash), nil
	}
	return hash, nil
}

// ValidationDecorator is a concrete decorator that validates text.
type ValidationDecorator struct {
	TextProcessorDecorator
	validators []func(string) error
}

// NewValidationDecorator creates a decorator that validates text.
func NewValidationDecorator(processor TextProcessor, validators ...func(string) error) *ValidationDecorator {
	return &ValidationDecorator{
		TextProcessorDecorator: TextProcessorDecorator{
			wrapped:     processor,
			name:        "Validation Processor",
			description: "Validates text against specified rules",
		},
		validators: validators,
	}
}

// Process first processes the text using the wrapped processor,
// then validates the text against the specified rules.
func (v *ValidationDecorator) Process(text string) (string, error) {
	// First, let the wrapped processor do its work
	processedText, err := v.wrapped.Process(text)
	if err != nil {
		return "", fmt.Errorf("error in wrapped processor: %w", err)
	}

	// Then validate the text
	for _, validator := range v.validators {
		if err := validator(processedText); err != nil {
			return "", fmt.Errorf("validation error: %w", err)
		}
	}

	return processedText, nil
}

// Common validation functions that can be used with ValidationDecorator

// ValidateNotEmpty checks if the text is not empty.
func ValidateNotEmpty(text string) error {
	if strings.TrimSpace(text) == "" {
		return fmt.Errorf("text cannot be empty")
	}
	return nil
}

// ValidateMaxLength checks if the text does not exceed the specified length.
func ValidateMaxLength(maxLength int) func(string) error {
	return func(text string) error {
		if len(text) > maxLength {
			return fmt.Errorf("text exceeds maximum length of %d characters", maxLength)
		}
		return nil
	}
}

// ValidateMinLength checks if the text meets the minimum length.
func ValidateMinLength(minLength int) func(string) error {
	return func(text string) error {
		if len(text) < minLength {
			return fmt.Errorf("text does not meet minimum length of %d characters", minLength)
		}
		return nil
	}
}

// ValidateContains checks if the text contains the specified substring.
func ValidateContains(substring string) func(string) error {
	return func(text string) error {
		if !strings.Contains(text, substring) {
			return fmt.Errorf("text does not contain required substring: %s", substring)
		}
		return nil
	}
}
