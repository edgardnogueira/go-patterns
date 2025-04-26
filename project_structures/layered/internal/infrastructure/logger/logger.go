package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/infrastructure/config"
)

// LogLevel represents the logging level
type LogLevel int

const (
	// DebugLevel represents debug log level
	DebugLevel LogLevel = iota
	// InfoLevel represents info log level
	InfoLevel
	// WarnLevel represents warning log level
	WarnLevel
	// ErrorLevel represents error log level
	ErrorLevel
	// FatalLevel represents fatal log level
	FatalLevel
)

// Logger represents a logger
type Logger struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
	level       LogLevel
	output      io.Writer
}

// NewLogger creates a new logger
func NewLogger(cfg config.LoggingConfig) (*Logger, error) {
	// Ensure log directory exists
	err := os.MkdirAll(strings.TrimSuffix(cfg.FilePath, "/app.log"), 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open log file
	file, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Create multi-writer for console and file
	multiWriter := io.MultiWriter(os.Stdout, file)

	// Parse log level
	level := parseLevel(cfg.Level)

	// Create logger instance
	logger := &Logger{
		debugLogger: log.New(multiWriter, "DEBUG: ", log.Ldate|log.Ltime),
		infoLogger:  log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime),
		warnLogger:  log.New(multiWriter, "WARN: ", log.Ldate|log.Ltime),
		errorLogger: log.New(multiWriter, "ERROR: ", log.Ldate|log.Ltime),
		fatalLogger: log.New(multiWriter, "FATAL: ", log.Ldate|log.Ltime),
		level:       level,
		output:      multiWriter,
	}

	logger.Info("Logger initialized with level: " + cfg.Level)
	return logger, nil
}

// parseLevel parses a log level string
func parseLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	default:
		return InfoLevel
	}
}

// Debug logs a debug message
func (l *Logger) Debug(message string) {
	if l.level <= DebugLevel {
		l.debugLogger.Println(message)
	}
}

// Info logs an info message
func (l *Logger) Info(message string) {
	if l.level <= InfoLevel {
		l.infoLogger.Println(message)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(message string) {
	if l.level <= WarnLevel {
		l.warnLogger.Println(message)
	}
}

// Error logs an error message
func (l *Logger) Error(message string) {
	if l.level <= ErrorLevel {
		l.errorLogger.Println(message)
	}
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(message string) {
	if l.level <= FatalLevel {
		l.fatalLogger.Println(message)
		os.Exit(1)
	}
}

// Debugf logs a formatted debug message
func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.level <= DebugLevel {
		l.debugLogger.Printf(format, args...)
	}
}

// Infof logs a formatted info message
func (l *Logger) Infof(format string, args ...interface{}) {
	if l.level <= InfoLevel {
		l.infoLogger.Printf(format, args...)
	}
}

// Warnf logs a formatted warning message
func (l *Logger) Warnf(format string, args ...interface{}) {
	if l.level <= WarnLevel {
		l.warnLogger.Printf(format, args...)
	}
}

// Errorf logs a formatted error message
func (l *Logger) Errorf(format string, args ...interface{}) {
	if l.level <= ErrorLevel {
		l.errorLogger.Printf(format, args...)
	}
}

// Fatalf logs a formatted fatal message and exits
func (l *Logger) Fatalf(format string, args ...interface{}) {
	if l.level <= FatalLevel {
		l.fatalLogger.Printf(format, args...)
		os.Exit(1)
	}
}

// WithField adds a field to the log message
func (l *Logger) WithField(key string, value interface{}) *Entry {
	return &Entry{
		logger: l,
		fields: map[string]interface{}{key: value},
	}
}

// WithFields adds multiple fields to the log message
func (l *Logger) WithFields(fields map[string]interface{}) *Entry {
	return &Entry{
		logger: l,
		fields: fields,
	}
}

// Entry represents a log entry with fields
type Entry struct {
	logger *Logger
	fields map[string]interface{}
}

// formatFields formats fields as a string
func (e *Entry) formatFields() string {
	if len(e.fields) == 0 {
		return ""
	}

	var fields string
	for key, value := range e.fields {
		fields += fmt.Sprintf("%s=%v ", key, value)
	}

	return "[" + strings.TrimSpace(fields) + "] "
}

// Debug logs a debug message with fields
func (e *Entry) Debug(message string) {
	if e.logger.level <= DebugLevel {
		e.logger.debugLogger.Println(e.formatFields() + message)
	}
}

// Info logs an info message with fields
func (e *Entry) Info(message string) {
	if e.logger.level <= InfoLevel {
		e.logger.infoLogger.Println(e.formatFields() + message)
	}
}

// Warn logs a warning message with fields
func (e *Entry) Warn(message string) {
	if e.logger.level <= WarnLevel {
		e.logger.warnLogger.Println(e.formatFields() + message)
	}
}

// Error logs an error message with fields
func (e *Entry) Error(message string) {
	if e.logger.level <= ErrorLevel {
		e.logger.errorLogger.Println(e.formatFields() + message)
	}
}

// Fatal logs a fatal message with fields and exits
func (e *Entry) Fatal(message string) {
	if e.logger.level <= FatalLevel {
		e.logger.fatalLogger.Println(e.formatFields() + message)
		os.Exit(1)
	}
}

// Now returns the current timestamp
func Now() string {
	return time.Now().Format(time.RFC3339)
}
