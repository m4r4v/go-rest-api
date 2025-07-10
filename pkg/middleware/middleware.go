package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/m4r4v/go-rest-api/pkg/auth"
	"github.com/m4r4v/go-rest-api/pkg/errors"
	"github.com/m4r4v/go-rest-api/pkg/logger"
)

// AuthMiddleware validates JWT tokens and sets user context
func AuthMiddleware(authService *auth.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeErrorResponse(w, errors.Unauthorized("Authorization header is required"))
				return
			}

			token, err := auth.ExtractBearerToken(authHeader)
			if err != nil {
				writeErrorResponse(w, errors.Unauthorized("Invalid authorization header format"))
				return
			}

			claims, err := authService.ValidateToken(token)
			if err != nil {
				writeErrorResponse(w, errors.Unauthorized("Invalid or expired token"))
				return
			}

			// Add claims to request context
			ctx := context.WithValue(r.Context(), auth.ClaimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole middleware checks if user has required role
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := auth.GetClaimsFromContext(r.Context())
			if claims == nil {
				writeErrorResponse(w, errors.Unauthorized("Authentication required"))
				return
			}

			if !claims.HasAnyRole(roles...) {
				writeErrorResponse(w, errors.Forbidden("Insufficient permissions"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		logger.Infof("HTTP %s %s %d %v %s",
			r.Method,
			r.RequestURI,
			wrapped.statusCode,
			duration,
			r.RemoteAddr,
		)
	})
}

// RecoveryMiddleware recovers from panics
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("Panic recovered: %v", err)
				writeErrorResponse(w, errors.InternalServerError("An internal server error occurred"))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// CORSMiddleware handles CORS headers
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// writeErrorResponse writes an error response in the standard format
func writeErrorResponse(w http.ResponseWriter, appErr *errors.AppError) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(appErr.Status)

	response := map[string]interface{}{
		"http_status_code":    appErr.Status,
		"http_status_message": http.StatusText(appErr.Status),
		"resource":            "",
		"app":                 "Go REST API Framework",
		"timestamp":           time.Now().Format(time.RFC3339),
		"response": map[string]interface{}{
			"error": map[string]interface{}{
				"code":    appErr.Code,
				"message": appErr.Message,
			},
		},
	}

	json.NewEncoder(w).Encode(response)
}
