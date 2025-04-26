package dto

import (
	"time"

	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/domain/model"
)

// AuthorCreateRequest represents the request body for creating a new author
type AuthorCreateRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
	Bio   string `json:"bio,omitempty"`
}

// AuthorUpdateRequest represents the request body for updating an existing author
type AuthorUpdateRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
	Bio   string `json:"bio,omitempty"`
}

// AuthorResponse represents the response body for an author
type AuthorResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Bio       string    `json:"bio,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ToResponse converts a domain Author model to an AuthorResponse DTO
func ToAuthorResponse(author *model.Author) *AuthorResponse {
	return &AuthorResponse{
		ID:        author.ID,
		Name:      author.Name,
		Email:     author.Email,
		Bio:       author.Bio,
		CreatedAt: author.CreatedAt,
		UpdatedAt: author.UpdatedAt,
	}
}

// ToAuthorResponses converts a slice of domain Author models to AuthorResponse DTOs
func ToAuthorResponses(authors []*model.Author) []*AuthorResponse {
	responses := make([]*AuthorResponse, len(authors))
	for i, author := range authors {
		responses[i] = ToAuthorResponse(author)
	}
	return responses
}

// ToDomain converts an AuthorCreateRequest to a domain Author model
func (r *AuthorCreateRequest) ToDomain() (*model.Author, error) {
	return model.NewAuthor(r.Name, r.Email, r.Bio)
}

// ToUpdateDomain converts an AuthorUpdateRequest to update a domain Author model
func (r *AuthorUpdateRequest) ToUpdateDomain(author *model.Author) error {
	return author.Update(r.Name, r.Email, r.Bio)
}
