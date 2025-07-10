package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Claims represents the JWT claims
type Claims struct {
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// AuthService handles authentication operations
type AuthService struct {
	jwtSecret     []byte
	jwtExpiration time.Duration
	bcryptCost    int
}

// NewAuthService creates a new authentication service
func NewAuthService(jwtSecret string, jwtExpiration time.Duration, bcryptCost int) *AuthService {
	return &AuthService{
		jwtSecret:     []byte(jwtSecret),
		jwtExpiration: jwtExpiration,
		bcryptCost:    bcryptCost,
	}
}

// NewAuthServiceFromConfig creates a new authentication service from config
func NewAuthServiceFromConfig(cfg AuthConfig) *AuthService {
	return &AuthService{
		jwtSecret:     []byte(cfg.JWTSecret),
		jwtExpiration: cfg.JWTExpiration,
		bcryptCost:    cfg.BcryptCost,
	}
}

// AuthConfig represents auth configuration
type AuthConfig struct {
	JWTSecret     string
	JWTExpiration time.Duration
	BcryptCost    int
}

// GenerateToken generates a JWT token for a user
func (a *AuthService) GenerateToken(userID, username string, roles []string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.jwtExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-rest-api",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.jwtSecret)
}

// ValidateToken validates a JWT token and returns the claims
func (a *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// HashPassword hashes a password using bcrypt
func (a *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), a.bcryptCost)
	return string(bytes), err
}

// CheckPassword checks if a password matches the hash
func (a *AuthService) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ExtractBearerToken extracts the token from Authorization header
func ExtractBearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("authorization header must be Bearer token")
	}

	return parts[1], nil
}

// HasRole checks if the user has a specific role
func (c *Claims) HasRole(role string) bool {
	for _, r := range c.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// HasAnyRole checks if the user has any of the specified roles
func (c *Claims) HasAnyRole(roles ...string) bool {
	for _, role := range roles {
		if c.HasRole(role) {
			return true
		}
	}
	return false
}

// Context keys for storing claims
type contextKey string

const ClaimsContextKey contextKey = "claims"

// GetClaimsFromContext retrieves claims from request context
func GetClaimsFromContext(ctx context.Context) *Claims {
	if claims, ok := ctx.Value(ClaimsContextKey).(*Claims); ok {
		return claims
	}
	return nil
}
