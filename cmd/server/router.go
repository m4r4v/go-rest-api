package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/m4r4v/go-rest-api/internal/handlers"
	"github.com/m4r4v/go-rest-api/pkg/auth"
	"github.com/m4r4v/go-rest-api/pkg/middleware"
)

// DynamicRouter wraps mux.Router with dynamic endpoint management
type DynamicRouter struct {
	*mux.Router
	protectedRouter *mux.Router
	apiHandlers     *handlers.APIHandlers
}

// NewDynamicRouter creates a new dynamic router
func NewDynamicRouter(apiHandlers *handlers.APIHandlers) *DynamicRouter {
	return &DynamicRouter{
		Router:      mux.NewRouter(),
		apiHandlers: apiHandlers,
	}
}

// SetProtectedRouter sets the protected subrouter for dynamic endpoints
func (dr *DynamicRouter) SetProtectedRouter(protectedRouter *mux.Router) {
	dr.protectedRouter = protectedRouter
}

// AddDynamicEndpoint adds a new dynamic endpoint based on resource data to the protected router
func (dr *DynamicRouter) AddDynamicEndpoint(endpoint, method string, response interface{}) {
	if dr.protectedRouter != nil {
		dr.apiHandlers.AddDynamicEndpoint(dr.protectedRouter, endpoint, method, response)
	} else {
		// Fallback to main router if protected router not set
		dr.apiHandlers.AddDynamicEndpoint(dr.Router, endpoint, method, response)
	}
}

// RemoveDynamicEndpoint removes a dynamic endpoint
func (dr *DynamicRouter) RemoveDynamicEndpoint(endpoint, method string) {
	// Note: Gorilla mux doesn't support removing routes dynamically
	// In a production system, you'd use a different router or implement route versioning
	// For now, we'll log this limitation
}

// setupRoutes configures all API routes
func setupRoutes(apiHandlers *handlers.APIHandlers, authService *auth.AuthService) *DynamicRouter {
	dynamicRouter := NewDynamicRouter(apiHandlers)
	router := dynamicRouter.Router

	// API v1 routes
	v1 := router.PathPrefix("/v1").Subrouter()

	// Setup endpoint (no authentication required)
	v1.HandleFunc("/setup", apiHandlers.Setup).Methods("POST")

	// Authentication endpoints (no authentication required)
	auth := v1.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", apiHandlers.Login).Methods("POST")

	// Protected routes (require authentication)
	protected := v1.PathPrefix("").Subrouter()
	protected.Use(func(next http.Handler) http.Handler {
		return middleware.Auth(authService)(next)
	})

	// User self-management endpoints
	protected.HandleFunc("/auth/me", apiHandlers.GetMe).Methods("GET")
	protected.HandleFunc("/users/me", apiHandlers.UpdateMe).Methods("PUT")

	// Admin-only user management endpoints
	adminUsers := protected.PathPrefix("/admin/users").Subrouter()
	adminUsers.Use(func(next http.Handler) http.Handler {
		return middleware.RequireRole("admin")(next)
	})
	adminUsers.HandleFunc("", apiHandlers.CreateUser).Methods("POST")
	adminUsers.HandleFunc("", apiHandlers.ListUsers).Methods("GET")
	adminUsers.HandleFunc("/{id}", apiHandlers.UpdateUserByAdmin).Methods("PUT")
	adminUsers.HandleFunc("/{id}", apiHandlers.DeleteUser).Methods("DELETE")

	// Resource management endpoints (all authenticated users)
	resources := protected.PathPrefix("/resources").Subrouter()
	resources.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		apiHandlers.CreateResourceWithDynamicEndpoint(w, r, dynamicRouter)
	}).Methods("POST")
	resources.HandleFunc("", apiHandlers.ListResources).Methods("GET")
	resources.HandleFunc("/{id}", apiHandlers.GetResource).Methods("GET")
	resources.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		apiHandlers.UpdateResourceWithDynamicEndpoint(w, r, dynamicRouter)
	}).Methods("PUT")
	resources.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		apiHandlers.DeleteResourceWithDynamicEndpoint(w, r, dynamicRouter)
	}).Methods("DELETE")

	// System endpoints (all authenticated users)
	protected.HandleFunc("/status", apiHandlers.GetStatus).Methods("GET")
	protected.HandleFunc("/health", apiHandlers.GetHealth).Methods("GET")

	// Set the protected router for dynamic endpoints
	dynamicRouter.SetProtectedRouter(protected)

	// Load existing dynamic endpoints on startup
	apiHandlers.LoadExistingEndpoints(dynamicRouter)

	// Add 404 handler with proper response format
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set development-friendly headers
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		w.WriteHeader(http.StatusNotFound)

		response := map[string]interface{}{
			"success":     false,
			"status_code": http.StatusNotFound,
			"status":      "Not Found",
			"error": map[string]interface{}{
				"code":    "ENDPOINT_NOT_FOUND",
				"message": "The requested endpoint was not found",
				"path":    r.URL.Path,
				"method":  r.Method,
			},
			"timestamp": "2025-07-03T20:44:00-04:00",
		}

		json.NewEncoder(w).Encode(response)
	})

	// Add middleware
	router.Use(func(next http.Handler) http.Handler {
		return middleware.Logger()(next)
	})
	router.Use(func(next http.Handler) http.Handler {
		return middleware.Recovery()(next)
	})
	router.Use(func(next http.Handler) http.Handler {
		return middleware.Security()(next)
	})

	return dynamicRouter
}
