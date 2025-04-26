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

// AuthorHandler handles HTTP requests for author endpoints
type AuthorHandler struct {
	authorService *service.AuthorAppService
	logger        *logger.Logger
	validate      *validator.Validate
}

// NewAuthorHandler creates a new author handler
func NewAuthorHandler(
	authorService *service.AuthorAppService,
	log *logger.Logger,
) *AuthorHandler {
	return &AuthorHandler{
		authorService: authorService,
		logger:        log,
		validate:      validator.New(),
	}
}

// RegisterRoutes registers routes for the author handler
func (h *AuthorHandler) RegisterRoutes(r chi.Router) {
	r.Route("/authors", func(r chi.Router) {
		r.Get("/", h.GetAllAuthors)
		r.Post("/", h.CreateAuthor)
		r.Get("/{id}", h.GetAuthorByID)
		r.Put("/{id}", h.UpdateAuthor)
		r.Delete("/{id}", h.DeleteAuthor)
		r.Get("/email/{email}", h.GetAuthorByEmail)
	})
}

// GetAllAuthors returns all authors
func (h *AuthorHandler) GetAllAuthors(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authors, err := h.authorService.GetAllAuthors(ctx)
	if err != nil {
		h.logger.WithField("error", err.Error()).Error("Failed to get all authors")
		writeJSONError(w, dto.NewInternalServerError("Failed to get authors"))
		return
	}

	response := dto.ToAuthorResponses(authors)
	writeJSON(w, http.StatusOK, response)
}

// GetAuthorByID returns an author by ID
func (h *AuthorHandler) GetAuthorByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")
	
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSONError(w, dto.NewBadRequestError("Invalid author ID", nil))
		return
	}

	author, postCount, err := h.authorService.GetAuthorWithPostCount(ctx, id)
	if err != nil {
		if err == service.ErrAuthorNotFound {
			writeJSONError(w, dto.NewNotFoundError("Author not found"))
			return
		}
		h.logger.WithFields(map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		}).Error("Failed to get author by ID")
		writeJSONError(w, dto.NewInternalServerError("Failed to get author"))
		return
	}

	response := dto.ToAuthorResponse(author)
	// Add post count as a detail
	response.PostCount = postCount

	writeJSON(w, http.StatusOK, response)
}

// CreateAuthor creates a new author
func (h *AuthorHandler) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.AuthorCreateRequest

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
	author, err := req.ToDomain()
	if err != nil {
		writeJSONError(w, dto.NewBadRequestError("Invalid author data", nil))
		return
	}

	// Create author
	if err := h.authorService.CreateAuthor(ctx, author); err != nil {
		switch err {
		case service.ErrAuthorExists:
			writeJSONError(w, dto.NewConflictError("Author with this email already exists", nil))
		case service.ErrInvalidAuthorData:
			writeJSONError(w, dto.NewBadRequestError("Invalid author data", nil))
		default:
			h.logger.WithFields(map[string]interface{}{
				"name":  req.Name,
				"email": req.Email,
				"error": err.Error(),
			}).Error("Failed to create author")
			writeJSONError(w, dto.NewInternalServerError("Failed to create author"))
		}
		return
	}

	response := dto.ToAuthorResponse(author)
	writeJSON(w, http.StatusCreated, response)
}

// UpdateAuthor updates an existing author
func (h *AuthorHandler) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")
	
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSONError(w, dto.NewBadRequestError("Invalid author ID", nil))
		return
	}

	var req dto.AuthorUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, dto.NewBadRequestError("Invalid request body", nil))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		validationErrors := translateValidationErrors(err)
		writeJSONError(w, dto.NewValidationError("Validation failed", validationErrors))
		return
	}

	// Get existing author
	existingAuthor, err := h.authorService.GetAuthorByID(ctx, id)
	if err != nil {
		if err == service.ErrAuthorNotFound {
			writeJSONError(w, dto.NewNotFoundError("Author not found"))
			return
		}
		h.logger.WithFields(map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		}).Error("Failed to get author for update")
		writeJSONError(w, dto.NewInternalServerError("Failed to update author"))
		return
	}

	// Update author fields
	if err := req.ToUpdateDomain(existingAuthor); err != nil {
		writeJSONError(w, dto.NewBadRequestError("Invalid author data", err.Error()))
		return
	}

	// Save updated author
	if err := h.authorService.UpdateAuthor(ctx, existingAuthor); err != nil {
		switch err {
		case service.ErrAuthorExists:
			writeJSONError(w, dto.NewConflictError("Email already in use by another author", nil))
		case service.ErrInvalidAuthorData:
			writeJSONError(w, dto.NewBadRequestError("Invalid author data", nil))
		default:
			h.logger.WithFields(map[string]interface{}{
				"id":    id,
				"error": err.Error(),
			}).Error("Failed to update author")
			writeJSONError(w, dto.NewInternalServerError("Failed to update author"))
		}
		return
	}

	response := dto.ToAuthorResponse(existingAuthor)
	writeJSON(w, http.StatusOK, response)
}

// DeleteAuthor deletes an author
func (h *AuthorHandler) DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")
	
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSONError(w, dto.NewBadRequestError("Invalid author ID", nil))
		return
	}

	if err := h.authorService.DeleteAuthor(ctx, id); err != nil {
		if err == service.ErrAuthorNotFound {
			writeJSONError(w, dto.NewNotFoundError("Author not found"))
			return
		}
		h.logger.WithFields(map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		}).Error("Failed to delete author")
		writeJSONError(w, dto.NewInternalServerError("Failed to delete author"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetAuthorByEmail returns an author by email
func (h *AuthorHandler) GetAuthorByEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	email := chi.URLParam(r, "email")

	author, err := h.authorService.GetAuthorByEmail(ctx, email)
	if err != nil {
		if err == service.ErrAuthorNotFound {
			writeJSONError(w, dto.NewNotFoundError("Author not found"))
			return
		}
		h.logger.WithFields(map[string]interface{}{
			"email": email,
			"error": err.Error(),
		}).Error("Failed to get author by email")
		writeJSONError(w, dto.NewInternalServerError("Failed to get author"))
		return
	}

	response := dto.ToAuthorResponse(author)
	writeJSON(w, http.StatusOK, response)
}

// AuthorWithPostCount extends AuthorResponse with post count
type AuthorWithPostCount struct {
	*dto.AuthorResponse
	PostCount int `json:"postCount"`
}
