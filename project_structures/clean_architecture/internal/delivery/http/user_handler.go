package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/edgardnogueira/go-patterns/project_structures/clean_architecture/internal/usecase"
)

// UserHandler handles HTTP requests related to users
type UserHandler struct {
	userService *usecase.UserService
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(userService *usecase.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// RegisterRoutes registers user-related routes to the given router
func (h *UserHandler) RegisterRoutes(router *gin.Engine) {
	userGroup := router.Group("/api/users")
	{
		userGroup.GET("", h.ListUsers)
		userGroup.POST("", h.CreateUser)
		userGroup.GET("/:id", h.GetUser)
		userGroup.PUT("/:id", h.UpdateUser)
		userGroup.DELETE("/:id", h.DeleteUser)
		userGroup.PATCH("/:id/password", h.UpdatePassword)
	}
}

// createUserRequest defines the structure for user creation requests
type createUserRequest struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// updateUserRequest defines the structure for user update requests
type updateUserRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email" binding:"omitempty,email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// updatePasswordRequest defines the structure for password update requests
type updatePasswordRequest struct {
	Password string `json:"password" binding:"required,min=8"`
}

// response is a generic response structure
type response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// pageInfo is used for pagination information
type pageInfo struct {
	Total  int `json:"total"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// CreateUser handles user creation requests
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	user, err := h.userService.CreateUser(
		c.Request.Context(),
		req.Username,
		req.Email,
		req.Password,
		req.FirstName,
		req.LastName,
	)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == usecase.ErrUserAlreadyExists {
			statusCode = http.StatusConflict
		}

		c.JSON(statusCode, response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response{
		Success: true,
		Message: "User created successfully",
		Data:    user,
	})
}

// GetUser handles requests to retrieve a user by ID
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response{
			Success: false,
			Error:   "User ID is required",
		})
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == usecase.ErrUserNotFound {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response{
		Success: true,
		Data:    user,
	})
}

// UpdateUser handles user update requests
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response{
			Success: false,
			Error:   "User ID is required",
		})
		return
	}

	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	user, err := h.userService.UpdateUser(
		c.Request.Context(),
		id,
		req.Username,
		req.Email,
		req.FirstName,
		req.LastName,
	)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == usecase.ErrUserNotFound {
			statusCode = http.StatusNotFound
		} else if err == usecase.ErrUserAlreadyExists {
			statusCode = http.StatusConflict
		}

		c.JSON(statusCode, response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response{
		Success: true,
		Message: "User updated successfully",
		Data:    user,
	})
}

// DeleteUser handles user deletion requests
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response{
			Success: false,
			Error:   "User ID is required",
		})
		return
	}

	err := h.userService.DeleteUser(c.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == usecase.ErrUserNotFound {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response{
		Success: true,
		Message: "User deleted successfully",
	})
}

// UpdatePassword handles password update requests
func (h *UserHandler) UpdatePassword(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response{
			Success: false,
			Error:   "User ID is required",
		})
		return
	}

	var req updatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	err := h.userService.UpdatePassword(c.Request.Context(), id, req.Password)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == usecase.ErrUserNotFound {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response{
		Success: true,
		Message: "Password updated successfully",
	})
}

// ListUsers handles requests to list users with pagination
func (h *UserHandler) ListUsers(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	users, total, err := h.userService.ListUsers(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{
			Success: false,
			Error:   "Failed to fetch users: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response{
		Success: true,
		Data: map[string]interface{}{
			"users": users,
			"page": pageInfo{
				Total:  total,
				Limit:  limit,
				Offset: offset,
			},
		},
	})
}
