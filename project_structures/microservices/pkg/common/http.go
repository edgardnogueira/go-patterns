package common

import (
	"encoding/json"
	"errors"
	"net/http"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error       string `json:"error"`
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	RequestID   string `json:"request_id,omitempty"`
}

// JSONResponse represents a successful JSON response
type JSONResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Message   string      `json:"message,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}

// ErrorCodes for standard error responses
const (
	ErrBadRequest          = "BAD_REQUEST"
	ErrUnauthorized        = "UNAUTHORIZED"
	ErrForbidden           = "FORBIDDEN"
	ErrNotFound            = "NOT_FOUND"
	ErrConflict            = "CONFLICT"
	ErrInternalServerError = "INTERNAL_SERVER_ERROR"
)

// RespondWithError sends an error response
func RespondWithError(w http.ResponseWriter, status int, message string, code string, requestID string) {
	resp := ErrorResponse{
		Error:       message,
		Code:        code,
		RequestID:   requestID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

// RespondWithJSON sends a JSON response
func RespondWithJSON(w http.ResponseWriter, status int, data interface{}, requestID string) {
	resp := JSONResponse{
		Success:   true,
		Data:      data,
		RequestID: requestID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

// RespondWithMessage sends a message response
func RespondWithMessage(w http.ResponseWriter, status int, message string, requestID string) {
	resp := JSONResponse{
		Success:   true,
		Message:   message,
		RequestID: requestID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

// DecodeJSONBody decodes a JSON request body
func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}, requestID string) error {
	if r.Body == nil {
		return errors.New("empty request body")
	}

	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload", ErrBadRequest, requestID)
		return err
	}

	return nil
}

// MapToStatusCode maps error codes to HTTP status codes
func MapToStatusCode(errorCode string) int {
	switch errorCode {
	case ErrBadRequest:
		return http.StatusBadRequest
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrForbidden:
		return http.StatusForbidden
	case ErrNotFound:
		return http.StatusNotFound
	case ErrConflict:
		return http.StatusConflict
	case ErrInternalServerError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
