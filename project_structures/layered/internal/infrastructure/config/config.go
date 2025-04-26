package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Logging  LoggingConfig  `json:"logging"`
	Worker   WorkerConfig   `json:"worker"`
}

// ServerConfig represents the server configuration
type ServerConfig struct {
	Host           string `json:"host"`
	Port           int    `json:"port"`
	ReadTimeout    int    `json:"readTimeout"`
	WriteTimeout   int    `json:"writeTimeout"`
	MaxHeaderBytes int    `json:"maxHeaderBytes"`
}

// DatabaseConfig represents the database configuration
type DatabaseConfig struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	FilePath string `json:"filePath"`
}

// LoggingConfig represents the logging configuration
type LoggingConfig struct {
	Level      string `json:"level"`
	FilePath   string `json:"filePath"`
	MaxSize    int    `json:"maxSize"`
	MaxBackups int    `json:"maxBackups"`
	MaxAge     int    `json:"maxAge"`
	Compress   bool   `json:"compress"`
}

// WorkerConfig represents the worker configuration
type WorkerConfig struct {
	Concurrency  int `json:"concurrency"`
	QueueSize    int `json:"queueSize"`
	PollInterval int `json:"pollInterval"`
}

// Load loads the configuration from a file
func Load(configPath string) (*Config, error) {
	// Set default configuration
	config := &Config{
		Server: ServerConfig{
			Host:           "localhost",
			Port:           8080,
			ReadTimeout:    30,
			WriteTimeout:   30,
			MaxHeaderBytes: 1 << 20, // 1MB
		},
		Database: DatabaseConfig{
			Type:     "sqlite",
			FilePath: "data/blog.db",
		},
		Logging: LoggingConfig{
			Level:      "info",
			FilePath:   "logs/app.log",
			MaxSize:    10,    // 10MB
			MaxBackups: 5,     // 5 files
			MaxAge:     30,    // 30 days
			Compress:   true,
		},
		Worker: WorkerConfig{
			Concurrency:  2,
			QueueSize:    100,
			PollInterval: 5, // 5 seconds
		},
	}

	// If configPath is empty, use default configuration
	if configPath == "" {
		return config, nil
	}

	// Read configuration file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse configuration
	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Override with environment variables
	config = overrideWithEnv(config)

	// Create necessary directories
	if err := ensureDirectories(config); err != nil {
		return nil, err
	}

	return config, nil
}

// overrideWithEnv overrides configuration with environment variables
func overrideWithEnv(config *Config) *Config {
	// Server configuration
	if host := os.Getenv("SERVER_HOST"); host != "" {
		config.Server.Host = host
	}
	// More environment variable overrides would be implemented here

	return config
}

// ensureDirectories ensures that necessary directories exist
func ensureDirectories(config *Config) error {
	// Ensure database directory
	if config.Database.Type == "sqlite" {
		dir := filepath.Dir(config.Database.FilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	// Ensure logs directory
	logDir := filepath.Dir(config.Logging.FilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	return nil
}

// GetConnectionString returns the database connection string
func (c *DatabaseConfig) GetConnectionString() string {
	switch strings.ToLower(c.Type) {
	case "sqlite":
		return c.FilePath
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.User, c.Password, c.Host, c.Port, c.Database)
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.Database)
	default:
		return c.FilePath
	}
}
