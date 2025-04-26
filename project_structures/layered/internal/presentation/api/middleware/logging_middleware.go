package middleware

import (
	"net/http"
	"time"

	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/infrastructure/logger"
)

// RequestLoggerMiddleware logs incoming HTTP requests
func RequestLoggerMiddleware(log *logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			// Create a response writer wrapper to capture the status code
			wrw := &responseWriterWrapper{
				ResponseWriter: w,
				statusCode:     http.StatusOK, // Default status code
			}

			// Process the request
			next.ServeHTTP(wrw, r)

			// Calculate request duration
			duration := time.Since(startTime)

			// Log request details
			log.WithFields(map[string]interface{}{
				"method":     r.Method,
				"path":       r.URL.Path,
				"query":      r.URL.RawQuery,
				"ip":         getClientIP(r),
				"user-agent": r.UserAgent(),
				"status":     wrw.statusCode,
				"duration":   duration.String(),
				"size":       wrw.bytesWritten,
			}).Info("HTTP Request")
		})
	}
}

// responseWriterWrapper is a wrapper for http.ResponseWriter to capture status code and response size
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int64
}

// WriteHeader captures the status code
func (rww *responseWriterWrapper) WriteHeader(statusCode int) {
	rww.statusCode = statusCode
	rww.ResponseWriter.WriteHeader(statusCode)
}

// Write captures the response size
func (rww *responseWriterWrapper) Write(b []byte) (int, error) {
	size, err := rww.ResponseWriter.Write(b)
	rww.bytesWritten += int64(size)
	return size, err
}

// getClientIP extracts the client IP address from the request
func getClientIP(r *http.Request) string {
	// Check for X-Forwarded-For header first (for proxies)
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return ip
	}

	// Check for X-Real-IP header
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}
