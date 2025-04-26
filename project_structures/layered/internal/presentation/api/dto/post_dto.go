package dto

import (
	"time"

	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/domain/model"
)

// PostCreateRequest represents the request body for creating a new post
type PostCreateRequest struct {
	Title    string   `json:"title" validate:"required,min=3,max=100"`
	Content  string   `json:"content" validate:"required,min=10"`
	AuthorID int64    `json:"authorId" validate:"required,gt=0"`
	Tags     []string `json:"tags,omitempty"`
}

// PostUpdateRequest represents the request body for updating an existing post
type PostUpdateRequest struct {
	Title    string   `json:"title" validate:"required,min=3,max=100"`
	Content  string   `json:"content" validate:"required,min=10"`
	Tags     []string `json:"tags,omitempty"`
}

// PostResponse represents the response body for a post
type PostResponse struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorID  int64     `json:"authorId"`
	Author    string    `json:"author,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Published bool      `json:"published"`
	Tags      []string  `json:"tags,omitempty"`
}

// ToResponse converts a domain Post model to a PostResponse DTO
func ToPostResponse(post *model.Post, authorName string) *PostResponse {
	return &PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		AuthorID:  post.AuthorID,
		Author:    authorName,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		Published: post.Published,
		Tags:      post.Tags,
	}
}

// ToPostResponses converts a slice of domain Post models to PostResponse DTOs
func ToPostResponses(posts []*model.Post, authorMap map[int64]string) []*PostResponse {
	responses := make([]*PostResponse, len(posts))
	for i, post := range posts {
		authorName := ""
		if name, ok := authorMap[post.AuthorID]; ok {
			authorName = name
		}
		responses[i] = ToPostResponse(post, authorName)
	}
	return responses
}

// ToBasicPostResponses converts a slice of domain Post models to basic PostResponse DTOs (without author names)
func ToBasicPostResponses(posts []*model.Post) []*PostResponse {
	responses := make([]*PostResponse, len(posts))
	for i, post := range posts {
		responses[i] = ToPostResponse(post, "")
	}
	return responses
}

// ToDomain converts a PostCreateRequest to a domain Post model
func (r *PostCreateRequest) ToDomain() (*model.Post, error) {
	return model.NewPost(r.Title, r.Content, r.AuthorID, r.Tags)
}
