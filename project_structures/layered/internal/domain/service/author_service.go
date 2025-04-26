package service

import (
	"context"
	"errors"

	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/domain/model"
)

// Common domain service errors for author
var (
	ErrAuthorNotFound = errors.New("author not found")
	ErrAuthorExists   = errors.New("author already exists")
)

// AuthorRepository defines the interface for author data access
type AuthorRepository interface {
	FindByID(ctx context.Context, id int64) (*model.Author, error)
	FindByEmail(ctx context.Context, email string) (*model.Author, error)
	FindAll(ctx context.Context) ([]*model.Author, error)
	Save(ctx context.Context, author *model.Author) error
	Update(ctx context.Context, author *model.Author) error
	Delete(ctx context.Context, id int64) error
}

// AuthorService defines the business operations for authors
type AuthorService interface {
	GetAuthorByID(ctx context.Context, id int64) (*model.Author, error)
	GetAuthorByEmail(ctx context.Context, email string) (*model.Author, error)
	GetAllAuthors(ctx context.Context) ([]*model.Author, error)
	CreateAuthor(ctx context.Context, author *model.Author) error
	UpdateAuthor(ctx context.Context, author *model.Author) error
	DeleteAuthor(ctx context.Context, id int64) error
}

// domainAuthorService implements AuthorService interface
type domainAuthorService struct {
	authorRepo AuthorRepository
}

// NewAuthorService creates a new instance of the author service
func NewAuthorService(repo AuthorRepository) AuthorService {
	return &domainAuthorService{
		authorRepo: repo,
	}
}

// GetAuthorByID retrieves an author by ID
func (s *domainAuthorService) GetAuthorByID(ctx context.Context, id int64) (*model.Author, error) {
	return s.authorRepo.FindByID(ctx, id)
}

// GetAuthorByEmail retrieves an author by email
func (s *domainAuthorService) GetAuthorByEmail(ctx context.Context, email string) (*model.Author, error) {
	return s.authorRepo.FindByEmail(ctx, email)
}

// GetAllAuthors retrieves all authors
func (s *domainAuthorService) GetAllAuthors(ctx context.Context) ([]*model.Author, error) {
	return s.authorRepo.FindAll(ctx)
}

// CreateAuthor creates a new author
func (s *domainAuthorService) CreateAuthor(ctx context.Context, author *model.Author) error {
	if err := author.Validate(); err != nil {
		return err
	}
	
	// Check if author with same email already exists
	existingAuthor, err := s.authorRepo.FindByEmail(ctx, author.Email)
	if err == nil && existingAuthor != nil {
		return ErrAuthorExists
	}
	
	return s.authorRepo.Save(ctx, author)
}

// UpdateAuthor updates an existing author
func (s *domainAuthorService) UpdateAuthor(ctx context.Context, author *model.Author) error {
	if err := author.Validate(); err != nil {
		return err
	}
	
	// Check if author exists
	_, err := s.authorRepo.FindByID(ctx, author.ID)
	if err != nil {
		return ErrAuthorNotFound
	}
	
	// Check if the updated email conflicts with another author
	existingAuthor, err := s.authorRepo.FindByEmail(ctx, author.Email)
	if err == nil && existingAuthor != nil && existingAuthor.ID != author.ID {
		return ErrAuthorExists
	}
	
	return s.authorRepo.Update(ctx, author)
}

// DeleteAuthor deletes an author by ID
func (s *domainAuthorService) DeleteAuthor(ctx context.Context, id int64) error {
	// Check if author exists
	_, err := s.authorRepo.FindByID(ctx, id)
	if err != nil {
		return ErrAuthorNotFound
	}
	
	return s.authorRepo.Delete(ctx, id)
}
