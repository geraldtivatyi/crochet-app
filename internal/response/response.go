package response

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse defines the standard JSON structure for all API errors
type ErrorResponse struct {
	Status  int               `json:"status"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"` // Field-level validation errors
}

// JSON renders any Go struct as JSON with the given HTTP status code
func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

// Error renders a standardized JSON error response
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, ErrorResponse{
		Status:  status,
		Message: message,
	})
}

// ValidationError renders a 400 Bad Request response with detailed field-level validation errors
func ValidationError(w http.ResponseWriter, details map[string]string) {
	JSON(w, http.StatusBadRequest, ErrorResponse{
		Status:  http.StatusBadRequest,
		Message: "Validation failed for request payload",
		Details: details,
	})
}
