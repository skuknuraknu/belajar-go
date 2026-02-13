package errors

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents the structure of an error response sent to the client.
// It includes an error message and optionally, more details.
type ErrorResponse struct {
	Error   string `json:"error"`             // A human-readable error message.
	Details string `json:"details,omitempty"` // Optional details about the error.
}

// APIError is a custom error type that includes an HTTP status code.
// This allows for more granular error handling in the API handlers.
type APIError struct {
	Status int    // HTTP status code.
	Err    string // Error message.
}

// Error implements the error interface for APIError.
func (e *APIError) Error() string {
	return e.Err
}

// NewAPIError creates a new APIError with the given status code and message.
func NewAPIError(status int, message string) *APIError {
	return &APIError{
		Status: status,
		Err:    message,
	}
}

// BadRequestError creates a new APIError for 400 Bad Request.
func BadRequestError(message string) *APIError {
	return NewAPIError(http.StatusBadRequest, message)
}

// NotFoundError creates a new APIError for 404 Not Found.
func NotFoundError(message string) *APIError {
	return NewAPIError(http.StatusNotFound, message)
}

// InternalServerError creates a new APIError for 500 Internal Server Error.
func InternalServerError(message string) *APIError {
	return NewAPIError(http.StatusInternalServerError, message)
}

// SendErrorResponse is a helper function to send a JSON error response.
// It sets the Content-Type header, status code, and writes the JSON error body.
func SendErrorResponse(w http.ResponseWriter, r *http.Request, status int, message string, details ...string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	errorResponse := ErrorResponse{
		Error: message,
	}
	if len(details) > 0 {
		errorResponse.Details = details[0]
	}

	err := json.NewEncoder(w).Encode(errorResponse)
	if err != nil {
		// If encoding the error response fails, log it and send a plain text error.
		// This is a fallback and should ideally never happen.
		http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
	}
}

// SendAPIErrorResponse is a convenience function to send an error response
// using an APIError type.
func SendAPIErrorResponse(w http.ResponseWriter, r *http.Request, apiErr *APIError) {
	SendErrorResponse(w, r, apiErr.Status, apiErr.Err)
}
