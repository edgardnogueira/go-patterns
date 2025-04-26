package entity

import (
	"database/sql"
	"time"
)

// PostEntity represents the database entity for a blog post
type PostEntity struct {
	ID        int64
	Title     string
	Content   string
	AuthorID  int64
	CreatedAt time.Time
	UpdatedAt time.Time
	Published bool
	Tags      string // Stored as comma-separated values in the database
}

// Schema returns the SQL schema for the post table
func (p *PostEntity) Schema() string {
	return `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		author_id INTEGER NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL,
		published BOOLEAN DEFAULT FALSE,
		tags TEXT DEFAULT '',
		FOREIGN KEY (author_id) REFERENCES authors(id)
	);
	`
}

// ScanRow scans a database row into a PostEntity
func (p *PostEntity) ScanRow(row *sql.Row) error {
	return row.Scan(
		&p.ID,
		&p.Title,
		&p.Content,
		&p.AuthorID,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.Published,
		&p.Tags,
	)
}

// ScanRows scans a database rows into a PostEntity
func (p *PostEntity) ScanRows(rows *sql.Rows) error {
	return rows.Scan(
		&p.ID,
		&p.Title,
		&p.Content,
		&p.AuthorID,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.Published,
		&p.Tags,
	)
}
