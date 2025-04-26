package model_test

import (
	"testing"
	"time"

	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/domain/model"
	"github.com/stretchr/testify/assert"
)

func TestNewPost(t *testing.T) {
	// Test cases
	tests := []struct {
		name      string
		title     string
		content   string
		authorID  int64
		tags      []string
		wantError bool
		errorType error
	}{
		{
			name:      "Valid post",
			title:     "Test Title",
			content:   "Test Content",
			authorID:  1,
			tags:      []string{"test", "go"},
			wantError: false,
		},
		{
			name:      "Empty title",
			title:     "",
			content:   "Test Content",
			authorID:  1,
			tags:      []string{"test"},
			wantError: true,
			errorType: model.ErrEmptyTitle,
		},
		{
			name:      "Empty content",
			title:     "Test Title",
			content:   "",
			authorID:  1,
			tags:      []string{"test"},
			wantError: true,
			errorType: model.ErrEmptyContent,
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			post, err := model.NewPost(tc.title, tc.content, tc.authorID, tc.tags)

			if tc.wantError {
				assert.Error(t, err)
				if tc.errorType != nil {
					assert.Equal(t, tc.errorType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, post)
				assert.Equal(t, tc.title, post.Title)
				assert.Equal(t, tc.content, post.Content)
				assert.Equal(t, tc.authorID, post.AuthorID)
				assert.Equal(t, tc.tags, post.Tags)
				assert.False(t, post.Published)
				assert.WithinDuration(t, time.Now(), post.CreatedAt, time.Second)
				assert.WithinDuration(t, time.Now(), post.UpdatedAt, time.Second)
			}
		})
	}
}

func TestPost_Validate(t *testing.T) {
	// Test cases
	tests := []struct {
		name      string
		post      *model.Post
		wantError bool
		errorType error
	}{
		{
			name: "Valid post",
			post: &model.Post{
				ID:       1,
				Title:    "Test Title",
				Content:  "Test Content",
				AuthorID: 1,
			},
			wantError: false,
		},
		{
			name: "Empty title",
			post: &model.Post{
				ID:       1,
				Title:    "",
				Content:  "Test Content",
				AuthorID: 1,
			},
			wantError: true,
			errorType: model.ErrEmptyTitle,
		},
		{
			name: "Empty content",
			post: &model.Post{
				ID:       1,
				Title:    "Test Title",
				Content:  "",
				AuthorID: 1,
			},
			wantError: true,
			errorType: model.ErrEmptyContent,
		},
		{
			name: "Invalid ID",
			post: &model.Post{
				ID:       -1,
				Title:    "Test Title",
				Content:  "Test Content",
				AuthorID: 1,
			},
			wantError: true,
			errorType: model.ErrInvalidID,
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.post.Validate()

			if tc.wantError {
				assert.Error(t, err)
				if tc.errorType != nil {
					assert.Equal(t, tc.errorType, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPost_Update(t *testing.T) {
	// Create test post
	post := &model.Post{
		ID:        1,
		Title:     "Original Title",
		Content:   "Original Content",
		AuthorID:  1,
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now().Add(-24 * time.Hour),
		Tags:      []string{"original"},
	}

	// Test cases
	tests := []struct {
		name      string
		title     string
		content   string
		tags      []string
		wantError bool
		errorType error
	}{
		{
			name:      "Valid update",
			title:     "Updated Title",
			content:   "Updated Content",
			tags:      []string{"updated", "test"},
			wantError: false,
		},
		{
			name:      "Empty title",
			title:     "",
			content:   "Updated Content",
			tags:      []string{"updated"},
			wantError: true,
			errorType: model.ErrEmptyTitle,
		},
		{
			name:      "Empty content",
			title:     "Updated Title",
			content:   "",
			tags:      []string{"updated"},
			wantError: true,
			errorType: model.ErrEmptyContent,
		},
	}

	// Original timestamp for comparison
	originalUpdatedAt := post.UpdatedAt

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Reset post to original state for each test
			post.Title = "Original Title"
			post.Content = "Original Content"
			post.Tags = []string{"original"}
			post.UpdatedAt = originalUpdatedAt

			err := post.Update(tc.title, tc.content, tc.tags)

			if tc.wantError {
				assert.Error(t, err)
				if tc.errorType != nil {
					assert.Equal(t, tc.errorType, err)
				}
				// Ensure no changes were made
				assert.Equal(t, "Original Title", post.Title)
				assert.Equal(t, "Original Content", post.Content)
				assert.Equal(t, []string{"original"}, post.Tags)
				assert.Equal(t, originalUpdatedAt, post.UpdatedAt)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.title, post.Title)
				assert.Equal(t, tc.content, post.Content)
				assert.Equal(t, tc.tags, post.Tags)
				assert.True(t, post.UpdatedAt.After(originalUpdatedAt))
			}
		})
	}
}

func TestPost_Publish_And_Unpublish(t *testing.T) {
	// Create test post
	post := &model.Post{
		ID:        1,
		Title:     "Test Title",
		Content:   "Test Content",
		AuthorID:  1,
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now().Add(-24 * time.Hour),
		Published: false,
	}

	// Original timestamp for comparison
	originalUpdatedAt := post.UpdatedAt

	// Test Publish
	post.Publish()
	assert.True(t, post.Published)
	assert.True(t, post.UpdatedAt.After(originalUpdatedAt))

	// New timestamp for comparison
	newUpdatedAt := post.UpdatedAt

	// Test Unpublish
	post.Unpublish()
	assert.False(t, post.Published)
	assert.True(t, post.UpdatedAt.After(newUpdatedAt))
}

func TestPost_TagManagement(t *testing.T) {
	// Create test post
	post := &model.Post{
		ID:        1,
		Title:     "Test Title",
		Content:   "Test Content",
		AuthorID:  1,
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now().Add(-24 * time.Hour),
		Tags:      []string{"original", "test"},
	}

	// Original timestamp for comparison
	originalUpdatedAt := post.UpdatedAt

	// Test AddTag with new tag
	post.AddTag("new")
	assert.Equal(t, []string{"original", "test", "new"}, post.Tags)
	assert.True(t, post.UpdatedAt.After(originalUpdatedAt))

	// New timestamp for comparison
	newUpdatedAt := post.UpdatedAt

	// Test AddTag with existing tag (should not add duplicate)
	post.AddTag("test")
	assert.Equal(t, []string{"original", "test", "new"}, post.Tags)
	assert.Equal(t, newUpdatedAt, post.UpdatedAt) // Should not update timestamp

	// Test RemoveTag
	post.RemoveTag("test")
	assert.Equal(t, []string{"original", "new"}, post.Tags)
	assert.True(t, post.UpdatedAt.After(newUpdatedAt))

	// New timestamp for comparison
	finalUpdatedAt := post.UpdatedAt

	// Test RemoveTag with non-existent tag
	post.RemoveTag("nonexistent")
	assert.Equal(t, []string{"original", "new"}, post.Tags)
	assert.Equal(t, finalUpdatedAt, post.UpdatedAt) // Should not update timestamp
}

func TestPost_IsPublished(t *testing.T) {
	// Create test post
	post := &model.Post{
		ID:        1,
		Title:     "Test Title",
		Content:   "Test Content",
		AuthorID:  1,
		Published: false,
	}

	// Test initially unpublished
	assert.False(t, post.IsPublished())

	// Test after publishing
	post.Published = true
	assert.True(t, post.IsPublished())

	// Test after unpublishing
	post.Published = false
	assert.False(t, post.IsPublished())
}
