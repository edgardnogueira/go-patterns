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

// SqliteAuthorRepository implements the AuthorRepository interface using SQLite
type SqliteAuthorRepository struct {
	db     *sql.DB
	mapper *mapper.AuthorMapper
}

// NewSqliteAuthorRepository creates a new instance of SqliteAuthorRepository
func NewSqliteAuthorRepository(db *sql.DB) service.AuthorRepository {
	return &SqliteAuthorRepository{
		db:     db,
		mapper: mapper.NewAuthorMapper(),
	}
}

// FindByID retrieves an author by ID
func (r *SqliteAuthorRepository) FindByID(ctx context.Context, id int64) (*model.Author, error) {
	query := `SELECT id, name, email, bio, created_at, updated_at 
			 FROM authors WHERE id = ?`
	
	row := r.db.QueryRowContext(ctx, query, id)
	
	authorEntity := &entity.AuthorEntity{}
	err := authorEntity.ScanRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRows
		}
		return nil, fmt.Errorf("error scanning author: %w", err)
	}
	
	return r.mapper.ToDomain(authorEntity), nil
}

// FindByEmail retrieves an author by email
func (r *SqliteAuthorRepository) FindByEmail(ctx context.Context, email string) (*model.Author, error) {
	query := `SELECT id, name, email, bio, created_at, updated_at 
			 FROM authors WHERE email = ?`
	
	row := r.db.QueryRowContext(ctx, query, email)
	
	authorEntity := &entity.AuthorEntity{}
	err := authorEntity.ScanRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRows
		}
		return nil, fmt.Errorf("error scanning author: %w", err)
	}
	
	return r.mapper.ToDomain(authorEntity), nil
}

// FindAll retrieves all authors
func (r *SqliteAuthorRepository) FindAll(ctx context.Context) ([]*model.Author, error) {
	query := `SELECT id, name, email, bio, created_at, updated_at 
			 FROM authors ORDER BY name`
	
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying authors: %w", err)
	}
	defer rows.Close()
	
	var authorEntities []*entity.AuthorEntity
	
	for rows.Next() {
		authorEntity := &entity.AuthorEntity{}
		err := authorEntity.ScanRows(rows)
		if err != nil {
			return nil, fmt.Errorf("error scanning author row: %w", err)
		}
		authorEntities = append(authorEntities, authorEntity)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating author rows: %w", err)
	}
	
	return r.mapper.ToDomainList(authorEntities), nil
}

// Save creates a new author
func (r *SqliteAuthorRepository) Save(ctx context.Context, author *model.Author) error {
	query := `INSERT INTO authors (name, email, bio, created_at, updated_at) 
			 VALUES (?, ?, ?, ?, ?)`
	
	authorEntity := r.mapper.ToEntity(author)
	
	result, err := r.db.ExecContext(
		ctx, 
		query, 
		authorEntity.Name,
		authorEntity.Email,
		authorEntity.Bio,
		authorEntity.CreatedAt,
		authorEntity.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error saving author: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %w", err)
	}
	
	author.ID = id
	return nil
}

// Update updates an existing author
func (r *SqliteAuthorRepository) Update(ctx context.Context, author *model.Author) error {
	query := `UPDATE authors SET name = ?, email = ?, bio = ?, updated_at = ? 
			 WHERE id = ?`
	
	authorEntity := r.mapper.ToEntity(author)
	
	_, err := r.db.ExecContext(
		ctx, 
		query, 
		authorEntity.Name,
		authorEntity.Email,
		authorEntity.Bio,
		authorEntity.UpdatedAt,
		authorEntity.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating author: %w", err)
	}
	
	return nil
}

// Delete deletes an author by ID
func (r *SqliteAuthorRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM authors WHERE id = ?"
	
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting author: %w", err)
	}
	
	return nil
}
