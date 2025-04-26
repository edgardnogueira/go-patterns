package mapper

import (
	"strings"

	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/data/entity"
	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/domain/model"
)

// PostMapper converts between domain Post models and PostEntity data objects
type PostMapper struct{}

// NewPostMapper creates a new instance of PostMapper
func NewPostMapper() *PostMapper {
	return &PostMapper{}
}

// ToEntity converts a domain Post model to a PostEntity
func (m *PostMapper) ToEntity(post *model.Post) *entity.PostEntity {
	// Join tags array into comma-separated string for storage
	tags := strings.Join(post.Tags, ",")

	return &entity.PostEntity{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		AuthorID:  post.AuthorID,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		Published: post.Published,
		Tags:      tags,
	}
}

// ToDomain converts a PostEntity to a domain Post model
func (m *PostMapper) ToDomain(entity *entity.PostEntity) *model.Post {
	// Split comma-separated tags string into array
	var tags []string
	if entity.Tags != "" {
		tags = strings.Split(entity.Tags, ",")
	} else {
		tags = []string{}
	}

	return &model.Post{
		ID:        entity.ID,
		Title:     entity.Title,
		Content:   entity.Content,
		AuthorID:  entity.AuthorID,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		Published: entity.Published,
		Tags:      tags,
	}
}

// ToDomainList converts a slice of PostEntity to a slice of domain Post models
func (m *PostMapper) ToDomainList(entities []*entity.PostEntity) []*model.Post {
	posts := make([]*model.Post, len(entities))
	for i, entity := range entities {
		posts[i] = m.ToDomain(entity)
	}
	return posts
}
