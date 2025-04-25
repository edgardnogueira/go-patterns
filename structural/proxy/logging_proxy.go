package proxy

import (
	"fmt"
	"io"
	"os"
	"time"
)

// LogLevel represents the severity level of a log entry.
type LogLevel int

// Log levels
const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
)

// LoggingProxy adds logging capabilities around image operations.
// It logs method calls, parameters, and results.
type LoggingProxy struct {
	BaseProxy
	writer      io.Writer
	minLogLevel LogLevel
	prefix      string
}

// NewLoggingProxy creates a new logging proxy.
func NewLoggingProxy(realImage Image, minLogLevel LogLevel) *LoggingProxy {
	return &LoggingProxy{
		BaseProxy: BaseProxy{
			realImage: realImage,
		},
		writer:      os.Stdout,
		minLogLevel: minLogLevel,
		prefix:      "ImageLogger",
	}
}

// SetWriter sets the writer where logs will be written to.
func (p *LoggingProxy) SetWriter(writer io.Writer) {
	p.writer = writer
}

// SetLogLevel sets the minimum log level to be logged.
func (p *LoggingProxy) SetLogLevel(level LogLevel) {
	p.minLogLevel = level
}

// SetPrefix sets the prefix for log messages.
func (p *LoggingProxy) SetPrefix(prefix string) {
	p.prefix = prefix
}

// Display logs the display call and then delegates to the real image.
func (p *LoggingProxy) Display() error {
	p.log(INFO, "Display method called")
	
	startTime := time.Now()
	err := p.realImage.Display()
	elapsed := time.Since(startTime)
	
	if err != nil {
		p.log(ERROR, "Display method failed: %v", err)
	} else {
		p.log(DEBUG, "Display method completed in %v", elapsed)
	}
	
	return err
}

// GetFilename logs the call and delegates to the real image.
func (p *LoggingProxy) GetFilename() string {
	p.log(DEBUG, "GetFilename method called")
	filename := p.realImage.GetFilename()
	p.log(DEBUG, "GetFilename returned: %s", filename)
	return filename
}

// GetWidth logs the call and delegates to the real image.
func (p *LoggingProxy) GetWidth() int {
	p.log(DEBUG, "GetWidth method called")
	width := p.realImage.GetWidth()
	p.log(DEBUG, "GetWidth returned: %d", width)
	return width
}

// GetHeight logs the call and delegates to the real image.
func (p *LoggingProxy) GetHeight() int {
	p.log(DEBUG, "GetHeight method called")
	height := p.realImage.GetHeight()
	p.log(DEBUG, "GetHeight returned: %d", height)
	return height
}

// GetSize logs the call and delegates to the real image.
func (p *LoggingProxy) GetSize() int64 {
	p.log(DEBUG, "GetSize method called")
	size := p.realImage.GetSize()
	p.log(DEBUG, "GetSize returned: %d bytes", size)
	return size
}

// GetMetadata logs the call and delegates to the real image.
func (p *LoggingProxy) GetMetadata() map[string]string {
	p.log(DEBUG, "GetMetadata method called")
	metadata := p.realImage.GetMetadata()
	p.log(DEBUG, "GetMetadata returned %d entries", len(metadata))
	return metadata
}

// log writes a log message with timestamp, level, and formatted content.
func (p *LoggingProxy) log(level LogLevel, format string, args ...interface{}) {
	if level < p.minLogLevel {
		return
	}
	
	levelStr := "DEBUG"
	switch level {
	case INFO:
		levelStr = "INFO"
	case WARNING:
		levelStr = "WARNING"
	case ERROR:
		levelStr = "ERROR"
	}
	
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	message := fmt.Sprintf(format, args...)
	logEntry := fmt.Sprintf("[%s] %s [%s] %s - %s\n", 
		timestamp, p.prefix, levelStr, p.realImage.GetFilename(), message)
	
	_, _ = p.writer.Write([]byte(logEntry))
}
