package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/domain/model"
	domainservice "github.com/edgardnogueira/go-patterns/project_structures/layered/internal/domain/service"
)

// Application-level errors
var (
	ErrPostNotFound   = errors.New("post not found")
	ErrAuthorNotFound = errors.New("author not found")
	ErrInvalidInput   = errors.New("invalid input data")
	ErrPostExists     = errors.New("post already exists")
)

// PostAppService orchestrates post-related use cases
type PostAppService struct {
	postService   domainservice.PostService
	authorService domainservice.AuthorService
}

// NewPostAppService creates a new instance of PostAppService
func NewPostAppService(
	postService domainservice.PostService,
	authorService domainservice.AuthorService,
) *PostAppService {
	return &PostAppService{
		postService:   postService,
		authorService: authorService,
	}
}

// GetPostByID retrieves a post and its author by post ID
func (s *PostAppService) GetPostByID(ctx context.Context, id int64) (*model.Post, *model.Author, error) {
	post, err := s.postService.GetPostByID(ctx, id)
	if err != nil {
		if errors.Is(err, domainservice.ErrPostNotFound) {
			return nil, nil, ErrPostNotFound
		}
		return nil, nil, fmt.Errorf("failed to get post: %w", err)
	}

	author, err := s.authorService.GetAuthorByID(ctx, post.AuthorID)
	if err != nil {
		// We still return the post even if we can't get the author
		return post, nil, fmt.Errorf("failed to get author: %w", err)
	}

	return post, author, nil
}

// GetAllPosts retrieves all posts with author information
func (s *PostAppService) GetAllPosts(ctx context.Context) ([]*model.Post, map[int64]*model.Author, error) {
	posts, err := s.postService.GetAllPosts(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get posts: %w", err)
	}

	// Get all authors for the posts
	authorMap, err := s.getAuthorsForPosts(ctx, posts)
	if err != nil {
		// Still return posts even if we can't get authors
		return posts, nil, fmt.Errorf("failed to get authors: %w", err)
	}

	return posts, authorMap, nil
}

// GetPublishedPosts retrieves all published posts with author information
func (s *PostAppService) GetPublishedPosts(ctx context.Context) ([]*model.Post, map[int64]*model.Author, error) {
	posts, err := s.postService.GetPublishedPosts(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get published posts: %w", err)
	}

	// Get all authors for the posts
	authorMap, err := s.getAuthorsForPosts(ctx, posts)
	if err != nil {
		// Still return posts even if we can't get authors
		return posts, nil, fmt.Errorf("failed to get authors: %w", err)
	}

	return posts, authorMap, nil
}

// GetPostsByAuthor retrieves all posts by a specific author
func (s *PostAppService) GetPostsByAuthor(ctx context.Context, authorID int64) ([]*model.Post, *model.Author, error) {
	// First verify author exists
	author, err := s.authorService.GetAuthorByID(ctx, authorID)
	if err != nil {
		if errors.Is(err, domainservice.ErrAuthorNotFound) {
			return nil, nil, ErrAuthorNotFound
		}
		return nil, nil, fmt.Errorf("failed to get author: %w", err)
	}

	posts, err := s.postService.GetPostsByAuthor(ctx, authorID)
	if err != nil {
		return nil, author, fmt.Errorf("failed to get posts by author: %w", err)
	}

	return posts, author, nil
}

// CreatePost creates a new post after validating the author exists
func (s *PostAppService) CreatePost(ctx context.Context, post *model.Post) error {
	// Validate author exists
	_, err := s.authorService.GetAuthorByID(ctx, post.AuthorID)
	if err != nil {
		if errors.Is(err, domainservice.ErrAuthorNotFound) {
			return ErrAuthorNotFound
		}
		return fmt.Errorf("failed to validate author: %w", err)
	}

	// Create post
	if err := s.postService.CreatePost(ctx, post); err != nil {
		if errors.Is(err, domainservice.ErrPostExists) {
			return ErrPostExists
		}
		return fmt.Errorf("failed to create post: %w", err)
	}

	return nil
}

// UpdatePost updates an existing post
func (s *PostAppService) UpdatePost(ctx context.Context, post *model.Post) error {
	// Check if post exists
	existingPost, err := s.postService.GetPostByID(ctx, post.ID)
	if err != nil {
		if errors.Is(err, domainservice.ErrPostNotFound) {
			return ErrPostNotFound
		}
		return fmt.Errorf("failed to get post: %w", err)
	}

	// Preserve author ID and creation date
	post.AuthorID = existingPost.AuthorID
	post.CreatedAt = existingPost.CreatedAt

	// Update post
	if err := s.postService.UpdatePost(ctx, post); err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}

	return nil
}

// DeletePost deletes a post
func (s *PostAppService) DeletePost(ctx context.Context, id int64) error {
	if err := s.postService.DeletePost(ctx, id); err != nil {
		if errors.Is(err, domainservice.ErrPostNotFound) {
			return ErrPostNotFound
		}
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}

// PublishPost publishes a post
func (s *PostAppService) PublishPost(ctx context.Context, id int64) error {
	if err := s.postService.PublishPost(ctx, id); err != nil {
		if errors.Is(err, domainservice.ErrPostNotFound) {
			return ErrPostNotFound
		}
		return fmt.Errorf("failed to publish post: %w", err)
	}

	return nil
}

// UnpublishPost unpublishes a post
func (s *PostAppService) UnpublishPost(ctx context.Context, id int64) error {
	if err := s.postService.UnpublishPost(ctx, id); err != nil {
		if errors.Is(err, domainservice.ErrPostNotFound) {
			return ErrPostNotFound
		}
		return fmt.Errorf("failed to unpublish post: %w", err)
	}

	return nil
}

// Helper function to get authors for a list of posts
func (s *PostAppService) getAuthorsForPosts(ctx context.Context, posts []*model.Post) (map[int64]*model.Author, error) {
	// Create a map to avoid fetching the same author multiple times
	authorIds := make(map[int64]struct{})
	for _, post := range posts {
		authorIds[post.AuthorID] = struct{}{}
	}

	// Fetch all authors
	authorMap := make(map[int64]*model.Author)
	for authorID := range authorIds {
		author, err := s.authorService.GetAuthorByID(ctx, authorID)
		if err != nil {
			// Skip authors we can't find, but log error
			continue
		}
		authorMap[authorID] = author
	}

	return authorMap, nil
}
