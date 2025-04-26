package mapper

import (
	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/data/entity"
	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/domain/model"
)

// AuthorMapper converts between domain Author models and AuthorEntity data objects
type AuthorMapper struct{}

// NewAuthorMapper creates a new instance of AuthorMapper
func NewAuthorMapper() *AuthorMapper {
	return &AuthorMapper{}
}

// ToEntity converts a domain Author model to an AuthorEntity
func (m *AuthorMapper) ToEntity(author *model.Author) *entity.AuthorEntity {
	return &entity.AuthorEntity{
		ID:        author.ID,
		Name:      author.Name,
		Email:     author.Email,
		Bio:       author.Bio,
		CreatedAt: author.CreatedAt,
		UpdatedAt: author.UpdatedAt,
	}
}

// ToDomain converts an AuthorEntity to a domain Author model
func (m *AuthorMapper) ToDomain(entity *entity.AuthorEntity) *model.Author {
	return &model.Author{
		ID:        entity.ID,
		Name:      entity.Name,
		Email:     entity.Email,
		Bio:       entity.Bio,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

// ToDomainList converts a slice of AuthorEntity to a slice of domain Author models
func (m *AuthorMapper) ToDomainList(entities []*entity.AuthorEntity) []*model.Author {
	authors := make([]*model.Author, len(entities))
	for i, entity := range entities {
		authors[i] = m.ToDomain(entity)
	}
	return authors
}
