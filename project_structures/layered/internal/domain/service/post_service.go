package service

import (
	"context"
	"errors"

	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/domain/model"
)

// Common domain service errors
var (
	ErrPostNotFound = errors.New("post not found")
	ErrPostExists   = errors.New("post already exists")
)

// PostRepository defines the interface for post data access
type PostRepository interface {
	FindByID(ctx context.Context, id int64) (*model.Post, error)
	FindAll(ctx context.Context) ([]*model.Post, error)
	FindByAuthorID(ctx context.Context, authorID int64) ([]*model.Post, error)
	FindPublished(ctx context.Context) ([]*model.Post, error)
	Save(ctx context.Context, post *model.Post) error
	Update(ctx context.Context, post *model.Post) error
	Delete(ctx context.Context, id int64) error
}

// PostService defines the business operations for blog posts
type PostService interface {
	GetPostByID(ctx context.Context, id int64) (*model.Post, error)
	GetAllPosts(ctx context.Context) ([]*model.Post, error)
	GetPostsByAuthor(ctx context.Context, authorID int64) ([]*model.Post, error)
	GetPublishedPosts(ctx context.Context) ([]*model.Post, error)
	CreatePost(ctx context.Context, post *model.Post) error
	UpdatePost(ctx context.Context, post *model.Post) error
	DeletePost(ctx context.Context, id int64) error
	PublishPost(ctx context.Context, id int64) error
	UnpublishPost(ctx context.Context, id int64) error
}

// domainPostService implements PostService interface
type domainPostService struct {
	postRepo PostRepository
}

// NewPostService creates a new instance of the post service
func NewPostService(repo PostRepository) PostService {
	return &domainPostService{
		postRepo: repo,
	}
}

// GetPostByID retrieves a post by its ID
func (s *domainPostService) GetPostByID(ctx context.Context, id int64) (*model.Post, error) {
	return s.postRepo.FindByID(ctx, id)
}

// GetAllPosts retrieves all posts
func (s *domainPostService) GetAllPosts(ctx context.Context) ([]*model.Post, error) {
	return s.postRepo.FindAll(ctx)
}

// GetPostsByAuthor retrieves all posts by a specific author
func (s *domainPostService) GetPostsByAuthor(ctx context.Context, authorID int64) ([]*model.Post, error) {
	return s.postRepo.FindByAuthorID(ctx, authorID)
}

// GetPublishedPosts retrieves all published posts
func (s *domainPostService) GetPublishedPosts(ctx context.Context) ([]*model.Post, error) {
	return s.postRepo.FindPublished(ctx)
}

// CreatePost creates a new post
func (s *domainPostService) CreatePost(ctx context.Context, post *model.Post) error {
	if err := post.Validate(); err != nil {
		return err
	}
	return s.postRepo.Save(ctx, post)
}

// UpdatePost updates an existing post
func (s *domainPostService) UpdatePost(ctx context.Context, post *model.Post) error {
	if err := post.Validate(); err != nil {
		return err
	}
	
	// Check if post exists
	_, err := s.postRepo.FindByID(ctx, post.ID)
	if err != nil {
		return ErrPostNotFound
	}
	
	return s.postRepo.Update(ctx, post)
}

// DeletePost deletes a post by ID
func (s *domainPostService) DeletePost(ctx context.Context, id int64) error {
	// Check if post exists
	_, err := s.postRepo.FindByID(ctx, id)
	if err != nil {
		return ErrPostNotFound
	}
	
	return s.postRepo.Delete(ctx, id)
}

// PublishPost marks a post as published
func (s *domainPostService) PublishPost(ctx context.Context, id int64) error {
	post, err := s.postRepo.FindByID(ctx, id)
	if err != nil {
		return ErrPostNotFound
	}
	
	post.Publish()
	return s.postRepo.Update(ctx, post)
}

// UnpublishPost marks a post as unpublished
func (s *domainPostService) UnpublishPost(ctx context.Context, id int64) error {
	post, err := s.postRepo.FindByID(ctx, id)
	if err != nil {
		return ErrPostNotFound
	}
	
	post.Unpublish()
	return s.postRepo.Update(ctx, post)
}
