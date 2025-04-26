package entity

import (
	"database/sql"
	"time"
)

// NotificationEntity represents the database entity for a notification
type NotificationEntity struct {
	ID         int64
	UserID     int64
	Type       string
	Title      string
	Content    string
	ResourceID int64
	Read       bool
	CreatedAt  time.Time
}

// Schema returns the SQL schema for the notification table
func (n *NotificationEntity) Schema() string {
	return `
	CREATE TABLE IF NOT EXISTS notifications (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		type TEXT NOT NULL,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		resource_id INTEGER NOT NULL,
		read BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP NOT NULL
	);
	`
}

// ScanRow scans a database row into a NotificationEntity
func (n *NotificationEntity) ScanRow(row *sql.Row) error {
	return row.Scan(
		&n.ID,
		&n.UserID,
		&n.Type,
		&n.Title,
		&n.Content,
		&n.ResourceID,
		&n.Read,
		&n.CreatedAt,
	)
}

// ScanRows scans database rows into a NotificationEntity
func (n *NotificationEntity) ScanRows(rows *sql.Rows) error {
	return rows.Scan(
		&n.ID,
		&n.UserID,
		&n.Type,
		&n.Title,
		&n.Content,
		&n.ResourceID,
		&n.Read,
		&n.CreatedAt,
	)
}
