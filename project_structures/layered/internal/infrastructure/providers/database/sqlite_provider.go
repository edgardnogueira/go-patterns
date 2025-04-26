package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/data/entity"
	_ "github.com/mattn/go-sqlite3"
)

// SqliteProvider represents a SQLite database provider
type SqliteProvider struct {
	db         *sql.DB
	dbPath     string
	initialized bool
}

// NewSqliteProvider creates a new SQLite database provider
func NewSqliteProvider(dbPath string) (*SqliteProvider, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &SqliteProvider{
		db:     db,
		dbPath: dbPath,
	}, nil
}

// Initialize initializes the database schema
func (p *SqliteProvider) Initialize() error {
	if p.initialized {
		return nil
	}

	// Create schemas
	schemas := []string{
		(&entity.AuthorEntity{}).Schema(),
		(&entity.PostEntity{}).Schema(),
		(&entity.NotificationEntity{}).Schema(),
	}

	// Execute each schema
	for _, schema := range schemas {
		_, err := p.db.Exec(schema)
		if err != nil {
			return fmt.Errorf("failed to create schema: %w", err)
		}
	}

	p.initialized = true
	return nil
}

// GetDB returns the database connection
func (p *SqliteProvider) GetDB() *sql.DB {
	return p.db
}

// Close closes the database connection
func (p *SqliteProvider) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

// RunMigrations runs database migrations
func (p *SqliteProvider) RunMigrations() error {
	// In a real application, this would run migration scripts
	// For this example, we're just initializing the database
	return p.Initialize()
}

// Backup creates a backup of the database
func (p *SqliteProvider) Backup(backupPath string) error {
	// Ensure the database is closed before backup
	p.db.Close()

	// Copy the database file
	source, err := os.Open(p.dbPath)
	if err != nil {
		return fmt.Errorf("failed to open source database: %w", err)
	}
	defer source.Close()

	// Create backup file
	destination, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer destination.Close()

	// Copy data
	if _, err := destination.ReadFrom(source); err != nil {
		return fmt.Errorf("failed to copy database: %w", err)
	}

	// Reopen the database
	db, err := sql.Open("sqlite3", p.dbPath)
	if err != nil {
		return fmt.Errorf("failed to reopen database: %w", err)
	}
	p.db = db

	return nil
}
