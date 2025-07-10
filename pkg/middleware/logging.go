package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/m4r4v/go-rest-api/internal/models"
	"github.com/m4r4v/go-rest-api/pkg/auth"
	"github.com/m4r4v/go-rest-api/pkg/logger"
)

// Database interface for logging
type LogDatabase interface {
	CreateLog(logEntry *models.LogEntry) error
}

// UserInteractionLoggingMiddleware logs all user interactions with detailed information
func UserInteractionLoggingMiddleware(db LogDatabase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			requestID := uuid.New().String()

			// Create a response writer wrapper to capture status code and response
			wrapped := &loggingResponseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
				requestID:      requestID,
			}

			// Add request ID to headers for tracing
			wrapped.Header().Set("X-Request-ID", requestID)

			// Execute the request
			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)

			// Get user information from context if available
			var userID, username string
			claims := auth.GetClaimsFromContext(r.Context())
			if claims != nil {
				userID = claims.UserID
				username = claims.Username
			}

			// Determine log level based on status code and action
			level := determineLogLevel(wrapped.statusCode, r.Method, r.URL.Path)

			// Create detailed log entry
			logEntry := &models.LogEntry{
				ID:         uuid.New().String(),
				UserID:     userID,
				Username:   username,
				Level:      level,
				Message:    generateLogMessage(r.Method, r.URL.Path, wrapped.statusCode, duration),
				Action:     determineAction(r.Method, r.URL.Path),
				Resource:   r.URL.Path,
				Method:     r.Method,
				StatusCode: wrapped.statusCode,
				IPAddress:  getClientIP(r),
				UserAgent:  r.UserAgent(),
				RequestID:  requestID,
				Duration:   duration,
				Metadata:   createMetadata(r, wrapped),
			}

			// Add error information if status indicates an error
			if wrapped.statusCode >= 400 {
				logEntry.Error = http.StatusText(wrapped.statusCode)
			}

			// Log to database
			if err := db.CreateLog(logEntry); err != nil {
				logger.Errorf("Failed to save log entry: %v", err)
			}

			// Also log to standard logger for immediate visibility
			logToStandardLogger(logEntry)
		})
	}
}

// loggingResponseWriter wraps http.ResponseWriter to capture response details
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	requestID  string
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// determineLogLevel determines the appropriate log level based on status code and action
func determineLogLevel(statusCode int, method, path string) models.LogLevel {
	// Error level for 4xx and 5xx status codes
	if statusCode >= 400 {
		if statusCode >= 500 {
			return models.LogLevelError
		}
		return models.LogLevelWarning
	}

	// Warning level for potentially sensitive operations
	if method == "DELETE" || strings.Contains(path, "/admin/") {
		return models.LogLevelWarning
	}

	// Info level for normal operations
	return models.LogLevelInfo
}

// determineAction extracts the action being performed from the request
func determineAction(method, path string) string {
	// Clean up the path
	cleanPath := strings.TrimPrefix(path, "/v1")
	cleanPath = strings.TrimPrefix(cleanPath, "/")

	// Extract base resource
	parts := strings.Split(cleanPath, "/")
	if len(parts) == 0 {
		return "unknown"
	}

	resource := parts[0]
	if resource == "" {
		resource = "root"
	}

	// Determine action based on method and path structure
	switch method {
	case "GET":
		if len(parts) > 1 && parts[1] != "" {
			return "get_" + resource
		}
		return "list_" + resource
	case "POST":
		if resource == "setup" {
			return "setup_admin"
		}
		if resource == "login" {
			return "user_login"
		}
		return "create_" + resource
	case "PUT":
		return "update_" + resource
	case "DELETE":
		return "delete_" + resource
	case "OPTIONS":
		return "options_" + resource
	default:
		return method + "_" + resource
	}
}

// generateLogMessage creates a human-readable log message
func generateLogMessage(method, path string, statusCode int, duration time.Duration) string {
	action := determineAction(method, path)

	if statusCode >= 400 {
		return "Failed " + action + ": " + method + " " + path + " (" + http.StatusText(statusCode) + ") in " + duration.String()
	}

	return "Successful " + action + ": " + method + " " + path + " (" + http.StatusText(statusCode) + ") in " + duration.String()
}

// getClientIP extracts the real client IP address
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (for proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}
	return ip
}

// createMetadata creates additional metadata for the log entry
func createMetadata(r *http.Request, wrapped *loggingResponseWriter) map[string]interface{} {
	metadata := map[string]interface{}{
		"content_length": r.ContentLength,
		"protocol":       r.Proto,
		"host":           r.Host,
		"referer":        r.Referer(),
	}

	// Add query parameters if present
	if r.URL.RawQuery != "" {
		metadata["query_params"] = r.URL.RawQuery
	}

	// Add content type if present
	if ct := r.Header.Get("Content-Type"); ct != "" {
		metadata["content_type"] = ct
	}

	// Add accept header if present
	if accept := r.Header.Get("Accept"); accept != "" {
		metadata["accept"] = accept
	}

	return metadata
}

// logToStandardLogger also logs to the standard logger for immediate visibility
func logToStandardLogger(logEntry *models.LogEntry) {
	logMessage := "[" + string(logEntry.Level) + "] " + logEntry.Message + " - User: " + logEntry.Username +
		", Action: " + logEntry.Action + ", Resource: " + logEntry.Resource + ", IP: " + logEntry.IPAddress

	switch logEntry.Level {
	case models.LogLevelError:
		logger.Error(logMessage)
	case models.LogLevelWarning:
		logger.Warn(logMessage)
	default:
		logger.Info(logMessage)
	}
}
