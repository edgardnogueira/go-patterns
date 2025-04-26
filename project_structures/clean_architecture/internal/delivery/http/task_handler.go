package http

import (
	"net/http"
	"strconv"

	"github.com/edgardnogueira/go-patterns/project_structures/clean_architecture/internal/domain/entities"
	"github.com/edgardnogueira/go-patterns/project_structures/clean_architecture/internal/usecase"
	"github.com/gin-gonic/gin"
)

// TaskHandler handles HTTP requests related to tasks
type TaskHandler struct {
	taskService *usecase.TaskService
}

// NewTaskHandler creates a new TaskHandler instance
func NewTaskHandler(taskService *usecase.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

// RegisterRoutes registers task-related routes to the given router
func (h *TaskHandler) RegisterRoutes(router *gin.Engine) {
	taskGroup := router.Group("/api/tasks")
	{
		taskGroup.GET("", h.ListTasks)
		taskGroup.GET("/:id", h.GetTask)
		taskGroup.PUT("/:id/start", h.StartTask)
		taskGroup.PUT("/:id/complete", h.CompleteTask)
		taskGroup.PUT("/:id/fail", h.FailTask)
		taskGroup.DELETE("/:id", h.DeleteTask)
	}
}

// failTaskRequest defines the structure for task failure requests
type failTaskRequest struct {
	ErrorMessage string `json:"error_message" binding:"required"`
}

// GetTask handles requests to retrieve a task by ID
func (h *TaskHandler) GetTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response{
			Success: false,
			Error:   "Task ID is required",
		})
		return
	}

	task, err := h.taskService.GetTaskByID(c.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == usecase.ErrTaskNotFound {
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
		Data:    task,
	})
}

// ListTasks handles requests to list tasks with status filtering and pagination
func (h *TaskHandler) ListTasks(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	statusStr := c.DefaultQuery("status", string(entities.TaskStatusPending))

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	status := entities.TaskStatus(statusStr)
	// Validate status
	switch status {
	case entities.TaskStatusPending, entities.TaskStatusProcessing, entities.TaskStatusCompleted, entities.TaskStatusFailed:
		// Valid status
	default:
		status = entities.TaskStatusPending
	}

	tasks, total, err := h.taskService.ListTasksByStatus(c.Request.Context(), status, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{
			Success: false,
			Error:   "Failed to fetch tasks: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response{
		Success: true,
		Data: map[string]interface{}{
			"tasks": tasks,
			"page": pageInfo{
				Total:  total,
				Limit:  limit,
				Offset: offset,
			},
		},
	})
}

// StartTask handles requests to start processing a task
func (h *TaskHandler) StartTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response{
			Success: false,
			Error:   "Task ID is required",
		})
		return
	}

	task, err := h.taskService.StartTask(c.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == usecase.ErrTaskNotFound {
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
		Message: "Task started successfully",
		Data:    task,
	})
}

// CompleteTask handles requests to mark a task as completed
func (h *TaskHandler) CompleteTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response{
			Success: false,
			Error:   "Task ID is required",
		})
		return
	}

	task, err := h.taskService.CompleteTask(c.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == usecase.ErrTaskNotFound {
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
		Message: "Task completed successfully",
		Data:    task,
	})
}

// FailTask handles requests to mark a task as failed
func (h *TaskHandler) FailTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response{
			Success: false,
			Error:   "Task ID is required",
		})
		return
	}

	var req failTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	task, err := h.taskService.FailTask(c.Request.Context(), id, req.ErrorMessage)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == usecase.ErrTaskNotFound {
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
		Message: "Task marked as failed",
		Data:    task,
	})
}

// DeleteTask handles task deletion requests
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response{
			Success: false,
			Error:   "Task ID is required",
		})
		return
	}

	err := h.taskService.DeleteTask(c.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == usecase.ErrTaskNotFound {
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
		Message: "Task deleted successfully",
	})
}
