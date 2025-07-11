package tests

import (
	"testing"
	"time"

	"github.com/m4r4v/go-rest-api/pkg/auth"
	"github.com/m4r4v/go-rest-api/pkg/config"
)

func TestAuthService(t *testing.T) {
	// Create test config
	cfg := &config.AuthConfig{
		JWTSecret:     "test-secret-key",
		JWTExpiration: time.Hour,
		BcryptCost:    4, // Lower cost for faster tests
	}

	authService := auth.NewAuthService(cfg.JWTSecret, cfg.JWTExpiration, cfg.BcryptCost)

	t.Run("HashPassword", func(t *testing.T) {
		password := "testpassword123"
		hash, err := authService.HashPassword(password)
		if err != nil {
			t.Fatalf("Failed to hash password: %v", err)
		}

		if hash == "" {
			t.Fatal("Hash should not be empty")
		}

		if hash == password {
			t.Fatal("Hash should not equal original password")
		}
	})

	t.Run("CheckPassword", func(t *testing.T) {
		password := "testpassword123"
		hash, err := authService.HashPassword(password)
		if err != nil {
			t.Fatalf("Failed to hash password: %v", err)
		}

		// Test correct password
		if !authService.CheckPassword(password, hash) {
			t.Fatal("CheckPassword should return true for correct password")
		}

		// Test incorrect password
		if authService.CheckPassword("wrongpassword", hash) {
			t.Fatal("CheckPassword should return false for incorrect password")
		}
	})

	t.Run("GenerateToken", func(t *testing.T) {
		userID := "123"
		username := "testuser"
		roles := []string{"user", "admin"}

		token, err := authService.GenerateToken(userID, username, roles)
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		if token == "" {
			t.Fatal("Token should not be empty")
		}
	})

	t.Run("ValidateToken", func(t *testing.T) {
		userID := "123"
		username := "testuser"
		roles := []string{"user", "admin"}

		// Generate token
		token, err := authService.GenerateToken(userID, username, roles)
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		// Validate token
		claims, err := authService.ValidateToken(token)
		if err != nil {
			t.Fatalf("Failed to validate token: %v", err)
		}

		if claims.UserID != userID {
			t.Fatalf("Expected UserID %s, got %s", userID, claims.UserID)
		}

		if claims.Username != username {
			t.Fatalf("Expected Username %s, got %s", username, claims.Username)
		}

		if len(claims.Roles) != len(roles) {
			t.Fatalf("Expected %d roles, got %d", len(roles), len(claims.Roles))
		}

		for i, role := range roles {
			if claims.Roles[i] != role {
				t.Fatalf("Expected role %s, got %s", role, claims.Roles[i])
			}
		}
	})

	t.Run("ValidateInvalidToken", func(t *testing.T) {
		invalidToken := "invalid.token.here"

		_, err := authService.ValidateToken(invalidToken)
		if err == nil {
			t.Fatal("Should return error for invalid token")
		}
	})

	t.Run("ExtractBearerToken", func(t *testing.T) {
		// Test valid bearer token
		authHeader := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
		token, err := auth.ExtractBearerToken(authHeader)
		if err != nil {
			t.Fatalf("Failed to extract bearer token: %v", err)
		}

		expected := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
		if token != expected {
			t.Fatalf("Expected token %s, got %s", expected, token)
		}

		// Test invalid format
		_, err = auth.ExtractBearerToken("InvalidFormat")
		if err == nil {
			t.Fatal("Should return error for invalid format")
		}

		// Test empty header
		_, err = auth.ExtractBearerToken("")
		if err == nil {
			t.Fatal("Should return error for empty header")
		}
	})

	t.Run("HasRole", func(t *testing.T) {
		claims := &auth.Claims{
			UserID:   "123",
			Username: "testuser",
			Roles:    []string{"user", "admin"},
		}

		if !claims.HasRole("user") {
			t.Fatal("Should have user role")
		}

		if !claims.HasRole("admin") {
			t.Fatal("Should have admin role")
		}

		if claims.HasRole("superadmin") {
			t.Fatal("Should not have superadmin role")
		}
	})

	t.Run("HasAnyRole", func(t *testing.T) {
		claims := &auth.Claims{
			UserID:   "123",
			Username: "testuser",
			Roles:    []string{"user"},
		}

		if !claims.HasAnyRole("user", "admin") {
			t.Fatal("Should have at least one of the specified roles")
		}

		if !claims.HasAnyRole("admin", "user") {
			t.Fatal("Should have at least one of the specified roles")
		}

		if claims.HasAnyRole("admin", "superadmin") {
			t.Fatal("Should not have any of the specified roles")
		}
	})
}
