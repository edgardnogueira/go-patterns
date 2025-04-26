package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/data/entity"
	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/data/mapper"
	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/domain/service"
)

// SqliteNotificationRepository implements the NotificationRepository interface using SQLite
type SqliteNotificationRepository struct {
	db     *sql.DB
	mapper *mapper.NotificationMapper
}

// NewSqliteNotificationRepository creates a new instance of SqliteNotificationRepository
func NewSqliteNotificationRepository(db *sql.DB) service.NotificationRepository {
	return &SqliteNotificationRepository{
		db:     db,
		mapper: mapper.NewNotificationMapper(),
	}
}

// FindByID retrieves a notification by ID
func (r *SqliteNotificationRepository) FindByID(ctx context.Context, id int64) (*service.Notification, error) {
	query := `SELECT id, user_id, type, title, content, resource_id, read, created_at
			 FROM notifications WHERE id = ?`
	
	row := r.db.QueryRowContext(ctx, query, id)
	
	notificationEntity := &entity.NotificationEntity{}
	err := notificationEntity.ScanRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRows
		}
		return nil, fmt.Errorf("error scanning notification: %w", err)
	}
	
	return r.mapper.ToDomain(notificationEntity), nil
}

// FindByUserID retrieves all notifications for a user
func (r *SqliteNotificationRepository) FindByUserID(ctx context.Context, userID int64) ([]*service.Notification, error) {
	query := `SELECT id, user_id, type, title, content, resource_id, read, created_at
			 FROM notifications WHERE user_id = ? ORDER BY created_at DESC`
	
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying notifications: %w", err)
	}
	defer rows.Close()
	
	var notificationEntities []*entity.NotificationEntity
	
	for rows.Next() {
		notificationEntity := &entity.NotificationEntity{}
		err := notificationEntity.ScanRows(rows)
		if err != nil {
			return nil, fmt.Errorf("error scanning notification row: %w", err)
		}
		notificationEntities = append(notificationEntities, notificationEntity)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating notification rows: %w", err)
	}
	
	return r.mapper.ToDomainList(notificationEntities), nil
}

// FindUnreadByUserID retrieves all unread notifications for a user
func (r *SqliteNotificationRepository) FindUnreadByUserID(ctx context.Context, userID int64) ([]*service.Notification, error) {
	query := `SELECT id, user_id, type, title, content, resource_id, read, created_at
			 FROM notifications WHERE user_id = ? AND read = FALSE ORDER BY created_at DESC`
	
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying unread notifications: %w", err)
	}
	defer rows.Close()
	
	var notificationEntities []*entity.NotificationEntity
	
	for rows.Next() {
		notificationEntity := &entity.NotificationEntity{}
		err := notificationEntity.ScanRows(rows)
		if err != nil {
			return nil, fmt.Errorf("error scanning notification row: %w", err)
		}
		notificationEntities = append(notificationEntities, notificationEntity)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating notification rows: %w", err)
	}
	
	return r.mapper.ToDomainList(notificationEntities), nil
}

// Save creates a new notification
func (r *SqliteNotificationRepository) Save(ctx context.Context, notification *service.Notification) error {
	query := `INSERT INTO notifications (user_id, type, title, content, resource_id, read, created_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?)`
	
	notificationEntity := r.mapper.ToEntity(notification)
	
	result, err := r.db.ExecContext(
		ctx, 
		query, 
		notificationEntity.UserID,
		notificationEntity.Type,
		notificationEntity.Title,
		notificationEntity.Content,
		notificationEntity.ResourceID,
		notificationEntity.Read,
		notificationEntity.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("error saving notification: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %w", err)
	}
	
	notification.ID = id
	return nil
}

// Update updates an existing notification
func (r *SqliteNotificationRepository) Update(ctx context.Context, notification *service.Notification) error {
	query := `UPDATE notifications SET title = ?, content = ?, read = ? WHERE id = ?`
	
	notificationEntity := r.mapper.ToEntity(notification)
	
	_, err := r.db.ExecContext(
		ctx, 
		query, 
		notificationEntity.Title,
		notificationEntity.Content,
		notificationEntity.Read,
		notificationEntity.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating notification: %w", err)
	}
	
	return nil
}

// Delete deletes a notification by ID
func (r *SqliteNotificationRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM notifications WHERE id = ?"
	
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting notification: %w", err)
	}
	
	return nil
}

// MarkAsRead marks a notification as read
func (r *SqliteNotificationRepository) MarkAsRead(ctx context.Context, id int64) error {
	query := "UPDATE notifications SET read = TRUE WHERE id = ?"
	
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error marking notification as read: %w", err)
	}
	
	return nil
}

// MarkAllAsRead marks all notifications as read for a user
func (r *SqliteNotificationRepository) MarkAllAsRead(ctx context.Context, userID int64) error {
	query := "UPDATE notifications SET read = TRUE WHERE user_id = ?"
	
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("error marking all notifications as read: %w", err)
	}
	
	return nil
}
