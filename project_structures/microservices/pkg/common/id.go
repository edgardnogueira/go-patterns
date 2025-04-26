package common

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// GenerateUUID generates a random UUID
func GenerateUUID() string {
	return uuid.New().String()
}

// GenerateOrderID generates a unique order ID
func GenerateOrderID() string {
	return fmt.Sprintf("ORD-%s-%d", randomString(8), time.Now().Unix())
}

// GenerateRequestID generates a unique request ID
func GenerateRequestID() string {
	return fmt.Sprintf("REQ-%s", randomString(12))
}

// randomString generates a random string of the given length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		// Fallback to time-based generation if crypto/rand fails
		for i := range b {
			b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		}
		return string(b)
	}

	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b)
}
