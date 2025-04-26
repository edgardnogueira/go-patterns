package dto

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Error   string      `json:"error,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(status int, message, errorType string, details interface{}) *ErrorResponse {
	return &ErrorResponse{
		Status:  status,
		Message: message,
		Error:   errorType,
		Details: details,
	}
}

// NewBadRequestError creates a bad request error response
func NewBadRequestError(message string, details interface{}) *ErrorResponse {
	return NewErrorResponse(400, message, "BadRequest", details)
}

// NewNotFoundError creates a not found error response
func NewNotFoundError(message string) *ErrorResponse {
	return NewErrorResponse(404, message, "NotFound", nil)
}

// NewInternalServerError creates an internal server error response
func NewInternalServerError(message string) *ErrorResponse {
	return NewErrorResponse(500, message, "InternalServerError", nil)
}

// NewConflictError creates a conflict error response
func NewConflictError(message string, details interface{}) *ErrorResponse {
	return NewErrorResponse(409, message, "Conflict", details)
}

// NewValidationError creates a validation error response with details
func NewValidationError(message string, validationErrors []ValidationError) *ErrorResponse {
	return NewErrorResponse(400, message, "ValidationError", validationErrors)
}

// NewUnauthorizedError creates an unauthorized error response
func NewUnauthorizedError(message string) *ErrorResponse {
	return NewErrorResponse(401, message, "Unauthorized", nil)
}

// NewForbiddenError creates a forbidden error response
func NewForbiddenError(message string) *ErrorResponse {
	return NewErrorResponse(403, message, "Forbidden", nil)
}
