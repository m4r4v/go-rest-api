package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Response structure for API responses
type APIResponse struct {
	Success    bool        `json:"success"`
	StatusCode int         `json:"status_code"`
	Status     string      `json:"status"`
	Data       interface{} `json:"data,omitempty"`
	Error      interface{} `json:"error,omitempty"`
	Timestamp  string      `json:"timestamp"`
}

// Health check response
type HealthResponse struct {
	Service   string `json:"service"`
	Version   string `json:"version"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

func main() {
	// Get port from environment variable (required by Cloud Run)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default fallback
	}

	// Validate port
	if _, err := strconv.Atoi(port); err != nil {
		log.Fatalf("Invalid PORT environment variable: %s", port)
	}

	log.Printf("Starting Go REST API Framework v2.0")
	log.Printf("Port: %s", port)

	// Setup routes
	router := setupRoutes()

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"*"},
	})

	// Create HTTP server - IMPORTANT: Listen on 0.0.0.0 as required by Cloud Run
	server := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", port),
		Handler:      c.Handler(router),
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create a deadline for shutdown (Cloud Run gives 10 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// setupRoutes configures all API routes
func setupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Health check endpoint (required for Cloud Run)
	router.HandleFunc("/health", healthHandler).Methods("GET")
	router.HandleFunc("/", rootHandler).Methods("GET")

	// API v1 routes
	v1 := router.PathPrefix("/v1").Subrouter()

	// Basic endpoints
	v1.HandleFunc("/status", statusHandler).Methods("GET")
	v1.HandleFunc("/ping", pingHandler).Methods("GET")

	// Add logging middleware
	router.Use(loggingMiddleware)
	router.Use(recoveryMiddleware)

	return router
}

// Health check handler
func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := APIResponse{
		Success:    true,
		StatusCode: http.StatusOK,
		Status:     "OK",
		Data: HealthResponse{
			Service:   "go-rest-api",
			Version:   "2.0.0",
			Status:    "healthy",
			Timestamp: time.Now().Format(time.RFC3339),
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Root handler
func rootHandler(w http.ResponseWriter, r *http.Request) {
	response := APIResponse{
		Success:    true,
		StatusCode: http.StatusOK,
		Status:     "OK",
		Data: map[string]interface{}{
			"message": "Welcome to Go REST API Framework v2.0",
			"version": "2.0.0",
			"endpoints": map[string]string{
				"health": "/health",
				"status": "/v1/status",
				"ping":   "/v1/ping",
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Status handler
func statusHandler(w http.ResponseWriter, r *http.Request) {
	response := APIResponse{
		Success:    true,
		StatusCode: http.StatusOK,
		Status:     "OK",
		Data: map[string]interface{}{
			"service":   "go-rest-api",
			"version":   "2.0.0",
			"status":    "running",
			"uptime":    "N/A",
			"timestamp": time.Now().Format(time.RFC3339),
			"environment": map[string]string{
				"port": os.Getenv("PORT"),
				"host": "0.0.0.0",
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Ping handler
func pingHandler(w http.ResponseWriter, r *http.Request) {
	response := APIResponse{
		Success:    true,
		StatusCode: http.StatusOK,
		Status:     "OK",
		Data: map[string]interface{}{
			"message":   "pong",
			"timestamp": time.Now().Format(time.RFC3339),
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	})
}

// Recovery middleware
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)

				response := APIResponse{
					Success:    false,
					StatusCode: http.StatusInternalServerError,
					Status:     "Internal Server Error",
					Error: map[string]interface{}{
						"code":    "INTERNAL_ERROR",
						"message": "An internal server error occurred",
					},
					Timestamp: time.Now().Format(time.RFC3339),
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
