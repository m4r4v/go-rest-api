package validation

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/m4r4v/go-rest-api/pkg/errors"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// ValidateJSON validates JSON request body
func ValidateJSON(r *http.Request, v interface{}) *errors.AppError {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return errors.BadRequest("Invalid JSON format")
	}

	if err := validate.Struct(v); err != nil {
		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, formatValidationError(err))
		}
		return errors.ValidationError(strings.Join(errorMessages, "; "))
	}

	return nil
}

// formatValidationError formats validation errors into human-readable messages
func formatValidationError(err validator.FieldError) string {
	field := err.Field()
	tag := err.Tag()

	switch tag {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email address"
	case "min":
		return field + " must be at least " + err.Param() + " characters long"
	case "max":
		return field + " must be at most " + err.Param() + " characters long"
	case "oneof":
		return field + " must be one of: " + err.Param()
	default:
		return field + " is invalid"
	}
}

// ValidatePassword checks if password meets strength requirements
func ValidatePassword(password string) *errors.AppError {
	if len(password) < 6 {
		return errors.ValidationError("Password must be at least 6 characters long")
	}
	return nil
}

// ValidateEmail checks if email format is valid
func ValidateEmail(email string) *errors.AppError {
	if err := validate.Var(email, "email"); err != nil {
		return errors.ValidationError("Invalid email format")
	}
	return nil
}
