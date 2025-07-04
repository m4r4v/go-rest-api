package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/m4r4v/go-rest-api/internal/handlers"
	"github.com/m4r4v/go-rest-api/pkg/auth"
	"github.com/m4r4v/go-rest-api/pkg/config"
	"github.com/m4r4v/go-rest-api/pkg/logger"
	"github.com/m4r4v/go-rest-api/pkg/validation"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger.Init(&cfg.Logger)
	logger.Info("Starting Go REST API Framework v2.0")

	// Initialize validation
	validation.Init()

	// Initialize auth service
	authService := auth.NewAuthService(&cfg.Auth)

	// Initialize API handlers with new architecture
	apiHandlers := handlers.NewAPIHandlers(authService)

	// Setup routes
	dynamicRouter := setupRoutes(apiHandlers, authService)

	// Create HTTP server
	server := &http.Server{
		Addr:         cfg.Server.Host + ":" + cfg.Server.Port,
		Handler:      dynamicRouter.Router,
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

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited")
}
