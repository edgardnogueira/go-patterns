package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/data/entity"
	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/data/mapper"
	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/domain/model"
	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/domain/service"
)

// Common repository errors
var (
	ErrNoRows = errors.New("no rows found")
)

// SqlitePostRepository implements the PostRepository interface using SQLite
type SqlitePostRepository struct {
	db     *sql.DB
	mapper *mapper.PostMapper
}

// NewSqlitePostRepository creates a new instance of SqlitePostRepository
func NewSqlitePostRepository(db *sql.DB) service.PostRepository {
	return &SqlitePostRepository{
		db:     db,
		mapper: mapper.NewPostMapper(),
	}
}

// FindByID retrieves a post by its ID
func (r *SqlitePostRepository) FindByID(ctx context.Context, id int64) (*model.Post, error) {
	query := `SELECT id, title, content, author_id, created_at, updated_at, published, tags 
			 FROM posts WHERE id = ?`
	
	row := r.db.QueryRowContext(ctx, query, id)
	
	postEntity := &entity.PostEntity{}
	err := postEntity.ScanRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRows
		}
		return nil, fmt.Errorf("error scanning post: %w", err)
	}
	
	return r.mapper.ToDomain(postEntity), nil
}

// FindAll retrieves all posts
func (r *SqlitePostRepository) FindAll(ctx context.Context) ([]*model.Post, error) {
	query := `SELECT id, title, content, author_id, created_at, updated_at, published, tags 
			 FROM posts ORDER BY created_at DESC`
	
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying posts: %w", err)
	}
	defer rows.Close()
	
	var postEntities []*entity.PostEntity
	
	for rows.Next() {
		postEntity := &entity.PostEntity{}
		err := postEntity.ScanRows(rows)
		if err != nil {
			return nil, fmt.Errorf("error scanning post row: %w", err)
		}
		postEntities = append(postEntities, postEntity)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating post rows: %w", err)
	}
	
	return r.mapper.ToDomainList(postEntities), nil
}

// FindByAuthorID retrieves all posts by a specific author
func (r *SqlitePostRepository) FindByAuthorID(ctx context.Context, authorID int64) ([]*model.Post, error) {
	query := `SELECT id, title, content, author_id, created_at, updated_at, published, tags 
			 FROM posts WHERE author_id = ? ORDER BY created_at DESC`
	
	rows, err := r.db.QueryContext(ctx, query, authorID)
	if err != nil {
		return nil, fmt.Errorf("error querying posts by author: %w", err)
	}
	defer rows.Close()
	
	var postEntities []*entity.PostEntity
	
	for rows.Next() {
		postEntity := &entity.PostEntity{}
		err := postEntity.ScanRows(rows)
		if err != nil {
			return nil, fmt.Errorf("error scanning post row: %w", err)
		}
		postEntities = append(postEntities, postEntity)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating post rows: %w", err)
	}
	
	return r.mapper.ToDomainList(postEntities), nil
}

// FindPublished retrieves all published posts
func (r *SqlitePostRepository) FindPublished(ctx context.Context) ([]*model.Post, error) {
	query := `SELECT id, title, content, author_id, created_at, updated_at, published, tags 
			 FROM posts WHERE published = TRUE ORDER BY created_at DESC`
	
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying published posts: %w", err)
	}
	defer rows.Close()
	
	var postEntities []*entity.PostEntity
	
	for rows.Next() {
		postEntity := &entity.PostEntity{}
		err := postEntity.ScanRows(rows)
		if err != nil {
			return nil, fmt.Errorf("error scanning post row: %w", err)
		}
		postEntities = append(postEntities, postEntity)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating post rows: %w", err)
	}
	
	return r.mapper.ToDomainList(postEntities), nil
}

// Save creates a new post
func (r *SqlitePostRepository) Save(ctx context.Context, post *model.Post) error {
	query := `INSERT INTO posts (title, content, author_id, created_at, updated_at, published, tags) 
			 VALUES (?, ?, ?, ?, ?, ?, ?)`
	
	postEntity := r.mapper.ToEntity(post)
	
	result, err := r.db.ExecContext(
		ctx, 
		query, 
		postEntity.Title,
		postEntity.Content,
		postEntity.AuthorID,
		postEntity.CreatedAt,
		postEntity.UpdatedAt,
		postEntity.Published,
		postEntity.Tags,
	)
	if err != nil {
		return fmt.Errorf("error saving post: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %w", err)
	}
	
	post.ID = id
	return nil
}

// Update updates an existing post
func (r *SqlitePostRepository) Update(ctx context.Context, post *model.Post) error {
	query := `UPDATE posts SET title = ?, content = ?, updated_at = ?, published = ?, tags = ? 
			 WHERE id = ?`
	
	postEntity := r.mapper.ToEntity(post)
	
	_, err := r.db.ExecContext(
		ctx, 
		query, 
		postEntity.Title,
		postEntity.Content,
		postEntity.UpdatedAt,
		postEntity.Published,
		postEntity.Tags,
		postEntity.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating post: %w", err)
	}
	
	return nil
}

// Delete deletes a post by ID
func (r *SqlitePostRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM posts WHERE id = ?"
	
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting post: %w", err)
	}
	
	return nil
}
