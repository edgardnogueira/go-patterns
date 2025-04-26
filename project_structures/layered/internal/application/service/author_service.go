package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/domain/model"
	domainservice "github.com/edgardnogueira/go-patterns/project_structures/layered/internal/domain/service"
)

// Application-level errors for author service
var (
	ErrAuthorNotFound   = errors.New("author not found")
	ErrAuthorExists     = errors.New("author already exists")
	ErrInvalidAuthorData = errors.New("invalid author data")
)

// AuthorAppService orchestrates author-related use cases
type AuthorAppService struct {
	authorService domainservice.AuthorService
	postService   domainservice.PostService
}

// NewAuthorAppService creates a new instance of AuthorAppService
func NewAuthorAppService(
	authorService domainservice.AuthorService,
	postService domainservice.PostService,
) *AuthorAppService {
	return &AuthorAppService{
		authorService: authorService,
		postService:   postService,
	}
}

// GetAuthorByID retrieves an author by ID
func (s *AuthorAppService) GetAuthorByID(ctx context.Context, id int64) (*model.Author, error) {
	author, err := s.authorService.GetAuthorByID(ctx, id)
	if err != nil {
		if errors.Is(err, domainservice.ErrAuthorNotFound) {
			return nil, ErrAuthorNotFound
		}
		return nil, fmt.Errorf("failed to get author: %w", err)
	}
	
	return author, nil
}

// GetAuthorByEmail retrieves an author by email
func (s *AuthorAppService) GetAuthorByEmail(ctx context.Context, email string) (*model.Author, error) {
	author, err := s.authorService.GetAuthorByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domainservice.ErrAuthorNotFound) {
			return nil, ErrAuthorNotFound
		}
		return nil, fmt.Errorf("failed to get author by email: %w", err)
	}
	
	return author, nil
}

// GetAllAuthors retrieves all authors
func (s *AuthorAppService) GetAllAuthors(ctx context.Context) ([]*model.Author, error) {
	authors, err := s.authorService.GetAllAuthors(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get authors: %w", err)
	}
	
	return authors, nil
}

// CreateAuthor creates a new author
func (s *AuthorAppService) CreateAuthor(ctx context.Context, author *model.Author) error {
	// Check if email already exists
	existingAuthor, err := s.authorService.GetAuthorByEmail(ctx, author.Email)
	if err == nil && existingAuthor != nil {
		return ErrAuthorExists
	}
	
	// Create author
	if err := s.authorService.CreateAuthor(ctx, author); err != nil {
		if errors.Is(err, domainservice.ErrAuthorExists) {
			return ErrAuthorExists
		}
		if errors.Is(err, model.ErrInvalidEmail) || errors.Is(err, model.ErrEmptyName) {
			return ErrInvalidAuthorData
		}
		return fmt.Errorf("failed to create author: %w", err)
	}
	
	return nil
}

// UpdateAuthor updates an existing author
func (s *AuthorAppService) UpdateAuthor(ctx context.Context, author *model.Author) error {
	// Check if author exists
	existingAuthor, err := s.authorService.GetAuthorByID(ctx, author.ID)
	if err != nil {
		if errors.Is(err, domainservice.ErrAuthorNotFound) {
			return ErrAuthorNotFound
		}
		return fmt.Errorf("failed to get author: %w", err)
	}
	
	// Check if the new email conflicts with another author
	if author.Email != existingAuthor.Email {
		emailAuthor, err := s.authorService.GetAuthorByEmail(ctx, author.Email)
		if err == nil && emailAuthor != nil && emailAuthor.ID != author.ID {
			return ErrAuthorExists
		}
	}
	
	// Preserve creation date
	author.CreatedAt = existingAuthor.CreatedAt
	
	// Update author
	if err := s.authorService.UpdateAuthor(ctx, author); err != nil {
		if errors.Is(err, model.ErrInvalidEmail) || errors.Is(err, model.ErrEmptyName) {
			return ErrInvalidAuthorData
		}
		return fmt.Errorf("failed to update author: %w", err)
	}
	
	return nil
}

// DeleteAuthor deletes an author and all their posts
func (s *AuthorAppService) DeleteAuthor(ctx context.Context, id int64) error {
	// Check if author exists
	_, err := s.authorService.GetAuthorByID(ctx, id)
	if err != nil {
		if errors.Is(err, domainservice.ErrAuthorNotFound) {
			return ErrAuthorNotFound
		}
		return fmt.Errorf("failed to get author: %w", err)
	}
	
	// Get all posts by this author
	posts, err := s.postService.GetPostsByAuthor(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get author's posts: %w", err)
	}
	
	// Delete all author's posts first
	for _, post := range posts {
		if err := s.postService.DeletePost(ctx, post.ID); err != nil {
			return fmt.Errorf("failed to delete author's post: %w", err)
		}
	}
	
	// Delete author
	if err := s.authorService.DeleteAuthor(ctx, id); err != nil {
		return fmt.Errorf("failed to delete author: %w", err)
	}
	
	return nil
}

// GetAuthorWithPostCount retrieves an author with their post count
func (s *AuthorAppService) GetAuthorWithPostCount(ctx context.Context, id int64) (*model.Author, int, error) {
	author, err := s.authorService.GetAuthorByID(ctx, id)
	if err != nil {
		if errors.Is(err, domainservice.ErrAuthorNotFound) {
			return nil, 0, ErrAuthorNotFound
		}
		return nil, 0, fmt.Errorf("failed to get author: %w", err)
	}
	
	posts, err := s.postService.GetPostsByAuthor(ctx, id)
	if err != nil {
		// Return author even if we can't get post count
		return author, 0, fmt.Errorf("failed to get author's posts: %w", err)
	}
	
	return author, len(posts), nil
}
