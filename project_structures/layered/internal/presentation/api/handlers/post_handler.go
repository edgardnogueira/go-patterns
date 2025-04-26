package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/application/service"
	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/domain/model"
	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/infrastructure/logger"
	"github.com/edgardnogueira/go-patterns/project_structures/layered/internal/presentation/api/dto"
)

// PostHandler handles HTTP requests for post endpoints
type PostHandler struct {
	postService   *service.PostAppService
	authorService *service.AuthorAppService
	logger        *logger.Logger
	validate      *validator.Validate
}

// NewPostHandler creates a new post handler
func NewPostHandler(
	postService *service.PostAppService,
	authorService *service.AuthorAppService,
	log *logger.Logger,
) *PostHandler {
	return &PostHandler{
		postService:   postService,
		authorService: authorService,
		logger:        log,
		validate:      validator.New(),
	}
}

// RegisterRoutes registers routes for the post handler
func (h *PostHandler) RegisterRoutes(r chi.Router) {
	r.Route("/posts", func(r chi.Router) {
		r.Get("/", h.GetAllPosts)
		r.Post("/", h.CreatePost)
		r.Get("/published", h.GetPublishedPosts)
		r.Get("/{id}", h.GetPostByID)
		r.Put("/{id}", h.UpdatePost)
		r.Delete("/{id}", h.DeletePost)
		r.Put("/{id}/publish", h.PublishPost)
		r.Put("/{id}/unpublish", h.UnpublishPost)
		r.Get("/author/{authorId}", h.GetPostsByAuthor)
	})
}

// GetAllPosts returns all posts
func (h *PostHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	posts, authorMap, err := h.postService.GetAllPosts(ctx)
	if err != nil {
		h.logger.WithField("error", err.Error()).Error("Failed to get all posts")
		writeJSONError(w, dto.NewInternalServerError("Failed to get posts"))
		return
	}

	// Convert authorMap to a map of author names by ID for simpler usage
	authorNames := make(map[int64]string)
	for id, author := range authorMap {
		authorNames[id] = author.Name
	}

	response := dto.ToPostResponses(posts, authorNames)
	writeJSON(w, http.StatusOK, response)
}

// GetPublishedPosts returns only published posts
func (h *PostHandler) GetPublishedPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	posts, authorMap, err := h.postService.GetPublishedPosts(ctx)
	if err != nil {
		h.logger.WithField("error", err.Error()).Error("Failed to get published posts")
		writeJSONError(w, dto.NewInternalServerError("Failed to get published posts"))
		return
	}

	// Convert authorMap to a map of author names by ID for simpler usage
	authorNames := make(map[int64]string)
	for id, author := range authorMap {
		authorNames[id] = author.Name
	}

	response := dto.ToPostResponses(posts, authorNames)
	writeJSON(w, http.StatusOK, response)
}

// GetPostByID returns a post by ID
func (h *PostHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")
	
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSONError(w, dto.NewBadRequestError("Invalid post ID", nil))
		return
	}

	post, author, err := h.postService.GetPostByID(ctx, id)
	if err != nil {
		if err == service.ErrPostNotFound {
			writeJSONError(w, dto.NewNotFoundError("Post not found"))
			return
		}
		h.logger.WithFields(map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		}).Error("Failed to get post by ID")
		writeJSONError(w, dto.NewInternalServerError("Failed to get post"))
		return
	}

	authorName := ""
	if author != nil {
		authorName = author.Name
	}

	response := dto.ToPostResponse(post, authorName)
	writeJSON(w, http.StatusOK, response)
}

// CreatePost creates a new post
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.PostCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, dto.NewBadRequestError("Invalid request body", nil))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		validationErrors := translateValidationErrors(err)
		writeJSONError(w, dto.NewValidationError("Validation failed", validationErrors))
		return
	}

	// Convert DTO to domain model
	post, err := req.ToDomain()
	if err != nil {
		writeJSONError(w, dto.NewBadRequestError("Invalid post data", nil))
		return
	}

	// Create post
	if err := h.postService.CreatePost(ctx, post); err != nil {
		switch err {
		case service.ErrAuthorNotFound:
			writeJSONError(w, dto.NewBadRequestError("Author not found", nil))
		case service.ErrPostExists:
			writeJSONError(w, dto.NewConflictError("Post already exists", nil))
		default:
			h.logger.WithFields(map[string]interface{}{
				"authorId": req.AuthorID,
				"title":    req.Title,
				"error":    err.Error(),
			}).Error("Failed to create post")
			writeJSONError(w, dto.NewInternalServerError("Failed to create post"))
		}
		return
	}

	// Get author name for response
	author, err := h.authorService.GetAuthorByID(ctx, post.AuthorID)
	authorName := ""
	if err == nil && author != nil {
		authorName = author.Name
	}

	response := dto.ToPostResponse(post, authorName)
	writeJSON(w, http.StatusCreated, response)
}

// UpdatePost updates an existing post
func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")
	
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSONError(w, dto.NewBadRequestError("Invalid post ID", nil))
		return
	}

	var req dto.PostUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, dto.NewBadRequestError("Invalid request body", nil))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		validationErrors := translateValidationErrors(err)
		writeJSONError(w, dto.NewValidationError("Validation failed", validationErrors))
		return
	}

	// Get existing post
	existingPost, _, err := h.postService.GetPostByID(ctx, id)
	if err != nil {
		if err == service.ErrPostNotFound {
			writeJSONError(w, dto.NewNotFoundError("Post not found"))
			return
		}
		h.logger.WithFields(map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		}).Error("Failed to get post for update")
		writeJSONError(w, dto.NewInternalServerError("Failed to update post"))
		return
	}

	// Update post fields
	if err := existingPost.Update(req.Title, req.Content, req.Tags); err != nil {
		writeJSONError(w, dto.NewBadRequestError("Invalid post data", err.Error()))
		return
	}

	// Save updated post
	if err := h.postService.UpdatePost(ctx, existingPost); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		}).Error("Failed to update post")
		writeJSONError(w, dto.NewInternalServerError("Failed to update post"))
		return
	}

	// Get author name for response
	author, err := h.authorService.GetAuthorByID(ctx, existingPost.AuthorID)
	authorName := ""
	if err == nil && author != nil {
		authorName = author.Name
	}

	response := dto.ToPostResponse(existingPost, authorName)
	writeJSON(w, http.StatusOK, response)
}

// DeletePost deletes a post
func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")
	
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSONError(w, dto.NewBadRequestError("Invalid post ID", nil))
		return
	}

	if err := h.postService.DeletePost(ctx, id); err != nil {
		if err == service.ErrPostNotFound {
			writeJSONError(w, dto.NewNotFoundError("Post not found"))
			return
		}
		h.logger.WithFields(map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		}).Error("Failed to delete post")
		writeJSONError(w, dto.NewInternalServerError("Failed to delete post"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// PublishPost publishes a post
func (h *PostHandler) PublishPost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")
	
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSONError(w, dto.NewBadRequestError("Invalid post ID", nil))
		return
	}

	if err := h.postService.PublishPost(ctx, id); err != nil {
		if err == service.ErrPostNotFound {
			writeJSONError(w, dto.NewNotFoundError("Post not found"))
			return
		}
		h.logger.WithFields(map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		}).Error("Failed to publish post")
		writeJSONError(w, dto.NewInternalServerError("Failed to publish post"))
		return
	}

	// Get updated post for response
	post, author, err := h.postService.GetPostByID(ctx, id)
	if err != nil {
		// Just return success without post details
		w.WriteHeader(http.StatusOK)
		return
	}

	authorName := ""
	if author != nil {
		authorName = author.Name
	}

	response := dto.ToPostResponse(post, authorName)
	writeJSON(w, http.StatusOK, response)
}

// UnpublishPost unpublishes a post
func (h *PostHandler) UnpublishPost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")
	
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSONError(w, dto.NewBadRequestError("Invalid post ID", nil))
		return
	}

	if err := h.postService.UnpublishPost(ctx, id); err != nil {
		if err == service.ErrPostNotFound {
			writeJSONError(w, dto.NewNotFoundError("Post not found"))
			return
		}
		h.logger.WithFields(map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		}).Error("Failed to unpublish post")
		writeJSONError(w, dto.NewInternalServerError("Failed to unpublish post"))
		return
	}

	// Get updated post for response
	post, author, err := h.postService.GetPostByID(ctx, id)
	if err != nil {
		// Just return success without post details
		w.WriteHeader(http.StatusOK)
		return
	}

	authorName := ""
	if author != nil {
		authorName = author.Name
	}

	response := dto.ToPostResponse(post, authorName)
	writeJSON(w, http.StatusOK, response)
}

// GetPostsByAuthor returns all posts by a specific author
func (h *PostHandler) GetPostsByAuthor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorIdStr := chi.URLParam(r, "authorId")
	
	authorId, err := strconv.ParseInt(authorIdStr, 10, 64)
	if err != nil {
		writeJSONError(w, dto.NewBadRequestError("Invalid author ID", nil))
		return
	}

	posts, author, err := h.postService.GetPostsByAuthor(ctx, authorId)
	if err != nil {
		if err == service.ErrAuthorNotFound {
			writeJSONError(w, dto.NewNotFoundError("Author not found"))
			return
		}
		h.logger.WithFields(map[string]interface{}{
			"authorId": authorId,
			"error":    err.Error(),
		}).Error("Failed to get posts by author")
		writeJSONError(w, dto.NewInternalServerError("Failed to get posts by author"))
		return
	}

	authorName := ""
	if author != nil {
		authorName = author.Name
	}

	// Create a map with just this author
	authorNames := map[int64]string{authorId: authorName}
	
	response := dto.ToPostResponses(posts, authorNames)
	writeJSON(w, http.StatusOK, response)
}

// Helper function to write JSON response
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.NewInternalServerError("Failed to encode response"))
		}
	}
}

// Helper function to write JSON error response
func writeJSONError(w http.ResponseWriter, err *dto.ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(err)
}

// Helper function to translate validation errors
func translateValidationErrors(err error) []dto.ValidationError {
	var validationErrors []dto.ValidationError

	validationErrs, ok := err.(validator.ValidationErrors)
	if ok {
		for _, e := range validationErrs {
			validationErrors = append(validationErrors, dto.ValidationError{
				Field:   e.Field(),
				Message: getValidationErrorMsg(e),
			})
		}
	}

	return validationErrors
}

// Helper function to get validation error message
func getValidationErrorMsg(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return "This field must be at least " + e.Param() + " characters long"
	case "max":
		return "This field must be no more than " + e.Param() + " characters long"
	case "email":
		return "This field must be a valid email"
	default:
		return "This field is invalid"
	}
}
