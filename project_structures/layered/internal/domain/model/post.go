package model

import (
	"errors"
	"time"
)

// Common domain errors
var (
	ErrEmptyTitle   = errors.New("post title cannot be empty")
	ErrEmptyContent = errors.New("post content cannot be empty")
	ErrInvalidID    = errors.New("invalid post ID")
)

// Post represents a blog post entity in the domain layer
type Post struct {
	ID        int64
	Title     string
	Content   string
	AuthorID  int64
	CreatedAt time.Time
	UpdatedAt time.Time
	Published bool
	Tags      []string
}

// NewPost creates a new blog post with validation
func NewPost(title, content string, authorID int64, tags []string) (*Post, error) {
	if title == "" {
		return nil, ErrEmptyTitle
	}
	if content == "" {
		return nil, ErrEmptyContent
	}

	now := time.Now()
	
	return &Post{
		Title:     title,
		Content:   content,
		AuthorID:  authorID,
		CreatedAt: now,
		UpdatedAt: now,
		Published: false,
		Tags:      tags,
	}, nil
}

// Validate ensures the post is in a valid state
func (p *Post) Validate() error {
	if p.Title == "" {
		return ErrEmptyTitle
	}
	if p.Content == "" {
		return ErrEmptyContent
	}
	if p.ID < 0 {
		return ErrInvalidID
	}
	return nil
}

// Publish marks the post as published
func (p *Post) Publish() {
	p.Published = true
	p.UpdatedAt = time.Now()
}

// Unpublish marks the post as unpublished/draft
func (p *Post) Unpublish() {
	p.Published = false
	p.UpdatedAt = time.Now()
}

// Update updates the post content with validation
func (p *Post) Update(title, content string, tags []string) error {
	if title == "" {
		return ErrEmptyTitle
	}
	if content == "" {
		return ErrEmptyContent
	}
	
	p.Title = title
	p.Content = content
	p.Tags = tags
	p.UpdatedAt = time.Now()
	
	return nil
}

// IsPublished checks if the post is published
func (p *Post) IsPublished() bool {
	return p.Published
}

// AddTag adds a new tag to the post if it doesn't already exist
func (p *Post) AddTag(tag string) {
	for _, existingTag := range p.Tags {
		if existingTag == tag {
			return
		}
	}
	p.Tags = append(p.Tags, tag)
	p.UpdatedAt = time.Now()
}

// RemoveTag removes a tag from the post
func (p *Post) RemoveTag(tag string) {
	for i, existingTag := range p.Tags {
		if existingTag == tag {
			p.Tags = append(p.Tags[:i], p.Tags[i+1:]...)
			p.UpdatedAt = time.Now()
			return
		}
	}
}
