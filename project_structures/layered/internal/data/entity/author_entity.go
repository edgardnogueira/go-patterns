package entity

import (
	"database/sql"
	"time"
)

// AuthorEntity represents the database entity for an author
type AuthorEntity struct {
	ID        int64
	Name      string
	Email     string
	Bio       string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Schema returns the SQL schema for the author table
func (a *AuthorEntity) Schema() string {
	return `
	CREATE TABLE IF NOT EXISTS authors (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		bio TEXT,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);
	`
}

// ScanRow scans a database row into an AuthorEntity
func (a *AuthorEntity) ScanRow(row *sql.Row) error {
	return row.Scan(
		&a.ID,
		&a.Name,
		&a.Email,
		&a.Bio,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
}

// ScanRows scans database rows into an AuthorEntity
func (a *AuthorEntity) ScanRows(rows *sql.Rows) error {
	return rows.Scan(
		&a.ID,
		&a.Name,
		&a.Email,
		&a.Bio,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
}
