package observability

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ContextKey is a type for context keys
type ContextKey string

// Context keys
const (
	RequestIDKey ContextKey = "request_id"
	UserIDKey    ContextKey = "user_id"
	ServiceKey   ContextKey = "service"
)

// Logger is a structured logger
type Logger struct {
	logger zerolog.Logger
}

// NewLogger creates a new structured logger
func NewLogger(serviceName string) *Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := zerolog.New(output).
		With().
		Timestamp().
		Str("service", serviceName).
		Logger()

	return &Logger{
		logger: logger,
	}
}

// WithContext returns a logger with context values
func (l *Logger) WithContext(ctx context.Context) *Logger {
	logger := l.logger.With().Logger()

	// Add request ID if present
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		logger = logger.With().Str("request_id", reqID).Logger()
	}

	// Add user ID if present
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		logger = logger.With().Str("user_id", userID).Logger()
	}

	return &Logger{
		logger: logger,
	}
}

// WithFields returns a logger with additional fields
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	ctx := l.logger.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	return &Logger{
		logger: ctx.Logger(),
	}
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		l.WithFields(fields[0]).logger.Debug().Msg(msg)
	} else {
		l.logger.Debug().Msg(msg)
	}
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		l.WithFields(fields[0]).logger.Info().Msg(msg)
	} else {
		l.logger.Info().Msg(msg)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		l.WithFields(fields[0]).logger.Warn().Msg(msg)
	} else {
		l.logger.Warn().Msg(msg)
	}
}

// Error logs an error message
func (l *Logger) Error(msg string, err error, fields ...map[string]interface{}) {
	logEvent := l.logger.Error()
	if err != nil {
		logEvent = logEvent.Err(err)
	}
	
	if len(fields) > 0 {
		for k, v := range fields[0] {
			logEvent = logEvent.Interface(k, v)
		}
	}
	
	logEvent.Msg(msg)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, err error, fields ...map[string]interface{}) {
	logEvent := l.logger.Fatal()
	if err != nil {
		logEvent = logEvent.Err(err)
	}
	
	if len(fields) > 0 {
		for k, v := range fields[0] {
			logEvent = logEvent.Interface(k, v)
		}
	}
	
	logEvent.Msg(msg)
}

// Configure sets up the global logger
func Configure(serviceName string) {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log.Logger = zerolog.New(output).
		With().
		Timestamp().
		Str("service", serviceName).
		Logger()
}
