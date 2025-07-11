package models

import (
	"encoding/json"
	"net/http"
	"time"
)

// StandardResponse represents the standardized API response format
type StandardResponse struct {
	Success    bool        `json:"success"`
	StatusCode int         `json:"status_code"`
	Status     string      `json:"status"`
	Timestamp  string      `json:"timestamp"`
	Endpoint   string      `json:"endpoint"`
	Method     string      `json:"method"`
	User       *string     `json:"user"`    // Pointer to allow null for unauthenticated requests
	UserID     *string     `json:"user_id"` // Pointer to allow null for unauthenticated requests
	Response   interface{} `json:"response"`
}

// ResponseWriter handles writing standardized responses
type ResponseWriter struct{}

// NewResponseWriter creates a new response writer
func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{}
}

// WriteSuccess writes a successful response
func (rw *ResponseWriter) WriteSuccess(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}, user *string, userID *string) {
	rw.writeResponse(w, r, true, statusCode, data, user, userID)
}

// WriteError writes an error response
func (rw *ResponseWriter) WriteError(w http.ResponseWriter, r *http.Request, statusCode int, errorData interface{}, user *string, userID *string) {
	errorResponse := map[string]interface{}{
		"error": errorData,
	}
	rw.writeResponse(w, r, false, statusCode, errorResponse, user, userID)
}

// writeResponse writes the standardized response format
func (rw *ResponseWriter) writeResponse(w http.ResponseWriter, r *http.Request, success bool, statusCode int, data interface{}, user *string, userID *string) {
	// Set security headers (CORS is handled by middleware)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("X-API-Framework", "Go-REST-API-v2.0")

	// Set status code
	w.WriteHeader(statusCode)

	// Create response
	response := StandardResponse{
		Success:    success,
		StatusCode: statusCode,
		Status:     http.StatusText(statusCode),
		Timestamp:  time.Now().Format(time.RFC3339),
		Endpoint:   r.URL.Path,
		Method:     r.Method,
		User:       user,
		UserID:     userID,
		Response:   data,
	}

	// Write JSON response
	json.NewEncoder(w).Encode(response)
}

// GetStatusCodeForMethod returns the appropriate success status code for HTTP methods
func GetStatusCodeForMethod(method string) int {
	switch method {
	case http.MethodGet:
		return http.StatusOK // 200
	case http.MethodPost:
		return http.StatusCreated // 201
	case http.MethodPut:
		return http.StatusCreated // 201 (Updated)
	case http.MethodDelete:
		return http.StatusCreated // 201 (Deleted)
	default:
		return http.StatusOK // 200
	}
}

// ErrorResponse represents a standardized error structure
type ErrorResponse struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(code, message string, details map[string]interface{}) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Common error responses
var (
	ErrInvalidCredentials   = NewErrorResponse("INVALID_CREDENTIALS", "Invalid username or password", nil)
	ErrUnauthorized         = NewErrorResponse("UNAUTHORIZED", "Authentication required", nil)
	ErrForbidden            = NewErrorResponse("FORBIDDEN", "Insufficient permissions", nil)
	ErrUserNotFoundResp     = NewErrorResponse("USER_NOT_FOUND", "User not found or you don't have permission to access it", nil)
	ErrResourceNotFoundResp = NewErrorResponse("RESOURCE_NOT_FOUND", "Resource not found or you don't have permission to access it", nil)
	ErrValidationError      = NewErrorResponse("VALIDATION_ERROR", "Invalid input data", nil)
	ErrSetupComplete        = NewErrorResponse("SETUP_ALREADY_COMPLETE", "Initial setup has already been completed", nil)
	ErrEndpointConflict     = NewErrorResponse("ENDPOINT_CONFLICT", "Dynamic endpoint conflicts with existing routes", nil)
	ErrInternalServer       = NewErrorResponse("INTERNAL_SERVER_ERROR", "An internal server error occurred", nil)
)
