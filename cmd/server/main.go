package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/m4r4v/go-rest-api/internal/handlers"
	"github.com/m4r4v/go-rest-api/pkg/auth"
	"github.com/m4r4v/go-rest-api/pkg/config"
	"github.com/m4r4v/go-rest-api/pkg/logger"
	"github.com/m4r4v/go-rest-api/pkg/middleware"
)

// StandardResponse represents the standard API response format
type StandardResponse struct {
	HTTPStatusCode    string      `json:"http_status_code"`
	HTTPStatusMessage string      `json:"http_status_message"`
	Resource          string      `json:"resource"`
	App               string      `json:"app"`
	Timestamp         string      `json:"timestamp"`
	Response          interface{} `json:"response"`
}

func main() {
	// Load environment variables from .env file (for local development)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger.Init(cfg.Log.Level, cfg.Log.Format)

	logger.Infof("Starting Go REST API Framework v2.0")
	logger.Infof("Port: %s", cfg.Server.Port)

	// Initialize auth service
	authService := auth.NewAuthService(cfg.Auth.JWTSecret, cfg.Auth.JWTExpiration, cfg.Auth.BcryptCost)

	// Initialize handlers
	apiHandlers := handlers.NewAPIHandlers(authService)

	// Setup routes
	router := setupRoutes(apiHandlers, authService)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		logger.Infof("Server starting on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Create a deadline for shutdown (Cloud Run gives 10 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited")
}

// setupRoutes configures all API routes according to the specification
func setupRoutes(apiHandlers *handlers.APIHandlers, authService *auth.AuthService) *mux.Router {
	router := mux.NewRouter()

	// Apply global middleware
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.RecoveryMiddleware)
	router.Use(middleware.CORSMiddleware)

	// Public endpoints (no authentication required)

	// /setup - Setup SuperAdmin account (POST only)
	router.HandleFunc("/setup", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeStandardError(w, http.StatusMethodNotAllowed, "/setup", "Method not allowed")
			return
		}
		apiHandlers.Setup(w, r)
	}).Methods("POST", "OPTIONS")

	// /login - User authentication (POST only)
	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeStandardError(w, http.StatusMethodNotAllowed, "/login", "Method not allowed")
			return
		}
		apiHandlers.Login(w, r)
	}).Methods("POST", "OPTIONS")

	// /status - Server status check (GET only)
	router.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeStandardError(w, http.StatusMethodNotAllowed, "/status", "Method not allowed")
			return
		}
		statusHandler(w, r)
	}).Methods("GET", "OPTIONS")

	// /v1/ping - Sample resource (GET only)
	router.HandleFunc("/v1/ping", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeStandardError(w, http.StatusMethodNotAllowed, "/v1/ping", "Method not allowed")
			return
		}
		pingHandler(w, r)
	}).Methods("GET", "OPTIONS")

	// Protected routes (require authentication)
	protected := router.PathPrefix("/v1").Subrouter()
	protected.Use(middleware.AuthMiddleware(authService))

	// User management endpoints
	protected.HandleFunc("/users/me", apiHandlers.GetMe).Methods("GET")
	protected.HandleFunc("/users/me", apiHandlers.UpdateMe).Methods("PUT")

	// Resource management endpoints
	protected.HandleFunc("/resources", apiHandlers.ListResources).Methods("GET")
	protected.HandleFunc("/resources", apiHandlers.CreateResource).Methods("POST")
	protected.HandleFunc("/resources/{id}", apiHandlers.GetResource).Methods("GET")
	protected.HandleFunc("/resources/{id}", apiHandlers.UpdateResource).Methods("PUT")
	protected.HandleFunc("/resources/{id}", apiHandlers.DeleteResource).Methods("DELETE")

	// Admin-only routes
	admin := protected.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.RequireRole("admin"))

	admin.HandleFunc("/users", apiHandlers.ListUsers).Methods("GET")
	admin.HandleFunc("/users", apiHandlers.CreateUser).Methods("POST")
	admin.HandleFunc("/users/{id}", apiHandlers.UpdateUserByAdmin).Methods("PUT")
	admin.HandleFunc("/users/{id}", apiHandlers.DeleteUser).Methods("DELETE")

	// Health check endpoint (for Cloud Run)
	router.HandleFunc("/health", healthHandler).Methods("GET")

	return router
}

// statusHandler handles the /status endpoint
func statusHandler(w http.ResponseWriter, r *http.Request) {
	response := StandardResponse{
		HTTPStatusCode:    "200",
		HTTPStatusMessage: "OK",
		Resource:          "/status",
		App:               "Go REST API Framework",
		Timestamp:         time.Now().Format(time.RFC3339),
		Response: map[string]interface{}{
			"status":  "healthy",
			"version": "2.0.0",
			"uptime":  "N/A",
		},
	}

	writeStandardResponse(w, response)
}

// pingHandler handles the /v1/ping endpoint
func pingHandler(w http.ResponseWriter, r *http.Request) {
	response := StandardResponse{
		HTTPStatusCode:    "200",
		HTTPStatusMessage: "OK",
		Resource:          "/v1/ping",
		App:               "Go REST API Framework",
		Timestamp:         time.Now().Format(time.RFC3339),
		Response: map[string]interface{}{
			"message": "pong",
		},
	}

	writeStandardResponse(w, response)
}

// healthHandler handles the /health endpoint for Cloud Run
func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := StandardResponse{
		HTTPStatusCode:    "200",
		HTTPStatusMessage: "OK",
		Resource:          "/health",
		App:               "Go REST API Framework",
		Timestamp:         time.Now().Format(time.RFC3339),
		Response: map[string]interface{}{
			"status":    "healthy",
			"service":   "go-rest-api",
			"version":   "2.0.0",
			"timestamp": time.Now().Format(time.RFC3339),
		},
	}

	writeStandardResponse(w, response)
}

// writeStandardResponse writes a response in the standard format
func writeStandardResponse(w http.ResponseWriter, response StandardResponse) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// writeStandardError writes an error response in the standard format
func writeStandardError(w http.ResponseWriter, statusCode int, resource, message string) {
	response := StandardResponse{
		HTTPStatusCode:    fmt.Sprintf("%d", statusCode),
		HTTPStatusMessage: http.StatusText(statusCode),
		Resource:          resource,
		App:               "Go REST API Framework",
		Timestamp:         time.Now().Format(time.RFC3339),
		Response: map[string]interface{}{
			"error": map[string]interface{}{
				"message": message,
			},
		},
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
