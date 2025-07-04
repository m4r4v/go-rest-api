package middleware

import (
	"encoding/json"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/m4r4v/go-rest-api/pkg/auth"
	"github.com/m4r4v/go-rest-api/pkg/config"
	"github.com/m4r4v/go-rest-api/pkg/errors"
	"github.com/m4r4v/go-rest-api/pkg/logger"
	"github.com/rs/cors"
	"golang.org/x/time/rate"
)

// Middleware represents a middleware function
type Middleware func(http.Handler) http.Handler

// Chain applies middlewares to a handler
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// Logger middleware logs HTTP requests
func Logger() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create a response writer wrapper to capture status code
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			logger.WithFields(logger.GetLogger().WithFields(map[string]interface{}{
				"method":      r.Method,
				"path":        r.URL.Path,
				"status":      wrapped.statusCode,
				"duration":    time.Since(start),
				"remote_addr": r.RemoteAddr,
				"user_agent":  r.UserAgent(),
			}).Data).Info("HTTP Request")
		})
	}
}

// Recovery middleware recovers from panics
func Recovery() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.WithFields(logger.GetLogger().WithFields(map[string]interface{}{
						"error": err,
						"stack": string(debug.Stack()),
						"path":  r.URL.Path,
					}).Data).Error("Panic recovered")

					appErr := errors.InternalServerError("Internal server error")
					writeErrorResponse(w, appErr)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// CORS middleware handles Cross-Origin Resource Sharing
func CORS(cfg *config.Config) Middleware {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Configure based on your needs
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	return func(next http.Handler) http.Handler {
		return c.Handler(next)
	}
}

// RateLimit middleware implements rate limiting
func RateLimit(requestsPerSecond int, burst int) Middleware {
	limiter := rate.NewLimiter(rate.Limit(requestsPerSecond), burst)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				appErr := errors.NewAppError("RATE_LIMIT_EXCEEDED", "Rate limit exceeded", http.StatusTooManyRequests)
				writeErrorResponse(w, appErr)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// ContentType middleware enforces JSON content type for specific methods
func ContentType() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only check content type for methods that typically have a body
			if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
				contentType := r.Header.Get("Content-Type")
				if contentType != "application/json" {
					appErr := errors.BadRequest("Content-Type must be application/json")
					writeErrorResponse(w, appErr)
					return
				}
			}

			// Set response content type
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	}
}

// Auth middleware validates JWT tokens
func Auth(authService *auth.AuthService) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			token, err := auth.ExtractBearerToken(authHeader)
			if err != nil {
				writeErrorResponse(w, err.(*errors.AppError))
				return
			}

			claims, err := authService.ValidateToken(token)
			if err != nil {
				appErr := errors.Unauthorized("Invalid or expired token")
				writeErrorResponse(w, appErr)
				return
			}

			// Add claims to request context
			ctx := r.Context()
			ctx = auth.SetClaimsInContext(ctx, claims)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// RequireRole middleware checks if user has required role
func RequireRole(roles ...string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := auth.GetClaimsFromContext(r.Context())
			if claims == nil {
				appErr := errors.Unauthorized("Authentication required")
				writeErrorResponse(w, appErr)
				return
			}

			if !claims.HasAnyRole(roles...) {
				appErr := errors.Forbidden("Insufficient permissions")
				writeErrorResponse(w, appErr)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Security middleware adds security headers
func Security() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

			next.ServeHTTP(w, r)
		})
	}
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

// writeErrorResponse writes an error response
func writeErrorResponse(w http.ResponseWriter, appErr *errors.AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Status)

	response := errors.ErrorResponse(appErr)
	json.NewEncoder(w).Encode(response)
}
