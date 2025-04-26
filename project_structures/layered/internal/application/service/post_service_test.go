package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/application/service"
	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/domain/model"
	domainservice "github.com/edgardnogueira/go-patterns/project_structures/layered/internal/domain/service"
)

// MockPostRepository is a mock implementation of domain post repository
type MockPostRepository struct {
	mock.Mock
}

// FindByID mocks the FindByID method
func (m *MockPostRepository) FindByID(ctx context.Context, id int64) (*model.Post, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Post), args.Error(1)
}

// FindAll mocks the FindAll method
func (m *MockPostRepository) FindAll(ctx context.Context) ([]*model.Post, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Post), args.Error(1)
}

// FindByAuthorID mocks the FindByAuthorID method
func (m *MockPostRepository) FindByAuthorID(ctx context.Context, authorID int64) ([]*model.Post, error) {
	args := m.Called(ctx, authorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Post), args.Error(1)
}

// FindPublished mocks the FindPublished method
func (m *MockPostRepository) FindPublished(ctx context.Context) ([]*model.Post, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Post), args.Error(1)
}

// Save mocks the Save method
func (m *MockPostRepository) Save(ctx context.Context, post *model.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

// Update mocks the Update method
func (m *MockPostRepository) Update(ctx context.Context, post *model.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockPostRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockAuthorRepository is a mock implementation of domain author repository
type MockAuthorRepository struct {
	mock.Mock
}

// FindByID mocks the FindByID method
func (m *MockAuthorRepository) FindByID(ctx context.Context, id int64) (*model.Author, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Author), args.Error(1)
}

// FindByEmail mocks the FindByEmail method
func (m *MockAuthorRepository) FindByEmail(ctx context.Context, email string) (*model.Author, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Author), args.Error(1)
}

// FindAll mocks the FindAll method
func (m *MockAuthorRepository) FindAll(ctx context.Context) ([]*model.Author, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Author), args.Error(1)
}

// Save mocks the Save method
func (m *MockAuthorRepository) Save(ctx context.Context, author *model.Author) error {
	args := m.Called(ctx, author)
	return args.Error(0)
}

// Update mocks the Update method
func (m *MockAuthorRepository) Update(ctx context.Context, author *model.Author) error {
	args := m.Called(ctx, author)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockAuthorRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestPostAppService_GetPostByID(t *testing.T) {
	// Create mocks
	mockPostRepo := new(MockPostRepository)
	mockAuthorRepo := new(MockAuthorRepository)

	// Create test data
	testPost := &model.Post{
		ID:        1,
		Title:     "Test Post",
		Content:   "Test Content",
		AuthorID:  2,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testAuthor := &model.Author{
		ID:    2,
		Name:  "Test Author",
		Email: "test@example.com",
	}

	// Test cases
	tests := []struct {
		name            string
		id              int64
		mockPostSetup   func()
		mockAuthorSetup func()
		wantPost        *model.Post
		wantAuthor      *model.Author
		wantErr         bool
		expectedErr     error
	}{
		{
			name: "Success - post and author found",
			id:   1,
			mockPostSetup: func() {
				mockPostRepo.On("FindByID", mock.Anything, int64(1)).Return(testPost, nil)
			},
			mockAuthorSetup: func() {
				mockAuthorRepo.On("FindByID", mock.Anything, int64(2)).Return(testAuthor, nil)
			},
			wantPost:   testPost,
			wantAuthor: testAuthor,
			wantErr:    false,
		},
		{
			name: "Post not found",
			id:   999,
			mockPostSetup: func() {
				mockPostRepo.On("FindByID", mock.Anything, int64(999)).Return(nil, domainservice.ErrPostNotFound)
			},
			mockAuthorSetup: func() {
				// No setup needed as author lookup won't happen
			},
			wantPost:    nil,
			wantAuthor:  nil,
			wantErr:     true,
			expectedErr: service.ErrPostNotFound,
		},
		{
			name: "Post found but author not found",
			id:   2,
			mockPostSetup: func() {
				mockPostRepo.On("FindByID", mock.Anything, int64(2)).Return(testPost, nil)
			},
			mockAuthorSetup: func() {
				mockAuthorRepo.On("FindByID", mock.Anything, int64(2)).Return(nil, domainservice.ErrAuthorNotFound)
			},
			wantPost:   testPost,
			wantAuthor: nil,
			wantErr:    false, // Not an error - we still return the post even if author lookup fails
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Reset mocks
			mockPostRepo = new(MockPostRepository)
			mockAuthorRepo = new(MockAuthorRepository)

			// Setup mocks
			tc.mockPostSetup()
			tc.mockAuthorSetup()

			// Create domain services
			postService := domainservice.NewPostService(mockPostRepo)
			authorService := domainservice.NewAuthorService(mockAuthorRepo)

			// Create application service
			postAppService := service.NewPostAppService(postService, authorService)

			// Call the method
			post, author, err := postAppService.GetPostByID(context.Background(), tc.id)

			// Assert expectations
			if tc.wantErr {
				assert.Error(t, err)
				if tc.expectedErr != nil {
					assert.Equal(t, tc.expectedErr, err)
				}
				assert.Nil(t, post)
				assert.Nil(t, author)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantPost, post)
				assert.Equal(t, tc.wantAuthor, author)
			}

			// Verify mocks
			mockPostRepo.AssertExpectations(t)
			mockAuthorRepo.AssertExpectations(t)
		})
	}
}

func TestPostAppService_CreatePost(t *testing.T) {
	// Create mocks
	mockPostRepo := new(MockPostRepository)
	mockAuthorRepo := new(MockAuthorRepository)

	// Create test data
	testPost, _ := model.NewPost("Test Post", "Test Content", 1, []string{"test"})
	testAuthor := &model.Author{
		ID:   1,
		Name: "Test Author",
	}

	// Test cases
	tests := []struct {
		name            string
		post            *model.Post
		mockPostSetup   func()
		mockAuthorSetup func()
		wantErr         bool
		expectedErr     error
	}{
		{
			name: "Success - author exists",
			post: testPost,
			mockAuthorSetup: func() {
				mockAuthorRepo.On("FindByID", mock.Anything, int64(1)).Return(testAuthor, nil)
			},
			mockPostSetup: func() {
				mockPostRepo.On("Save", mock.Anything, testPost).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Error - author not found",
			post: testPost,
			mockAuthorSetup: func() {
				mockAuthorRepo.On("FindByID", mock.Anything, int64(1)).Return(nil, domainservice.ErrAuthorNotFound)
			},
			mockPostSetup: func() {
				// No setup needed as save won't happen
			},
			wantErr:     true,
			expectedErr: service.ErrAuthorNotFound,
		},
		{
			name: "Error - post save fails",
			post: testPost,
			mockAuthorSetup: func() {
				mockAuthorRepo.On("FindByID", mock.Anything, int64(1)).Return(testAuthor, nil)
			},
			mockPostSetup: func() {
				mockPostRepo.On("Save", mock.Anything, testPost).Return(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Reset mocks
			mockPostRepo = new(MockPostRepository)
			mockAuthorRepo = new(MockAuthorRepository)

			// Setup mocks
			tc.mockAuthorSetup()
			tc.mockPostSetup()

			// Create domain services
			postService := domainservice.NewPostService(mockPostRepo)
			authorService := domainservice.NewAuthorService(mockAuthorRepo)

			// Create application service
			postAppService := service.NewPostAppService(postService, authorService)

			// Call the method
			err := postAppService.CreatePost(context.Background(), tc.post)

			// Assert expectations
			if tc.wantErr {
				assert.Error(t, err)
				if tc.expectedErr != nil {
					assert.Equal(t, tc.expectedErr, err)
				}
			} else {
				assert.NoError(t, err)
			}

			// Verify mocks
			mockPostRepo.AssertExpectations(t)
			mockAuthorRepo.AssertExpectations(t)
		})
	}
}
