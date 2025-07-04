package validation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/m4r4v/go-rest-api/pkg/errors"
)

var validate *validator.Validate

// Init initializes the validator
func Init() {
	validate = validator.New()

	// Register custom tag name function to use json tags
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// GetValidator returns the validator instance
func GetValidator() *validator.Validate {
	if validate == nil {
		Init()
	}
	return validate
}

// ValidateStruct validates a struct and returns formatted errors
func ValidateStruct(s interface{}) *errors.AppError {
	err := GetValidator().Struct(s)
	if err == nil {
		return nil
	}

	var validationErrors []string
	for _, err := range err.(validator.ValidationErrors) {
		validationErrors = append(validationErrors, formatValidationError(err))
	}

	return errors.ValidationError(strings.Join(validationErrors, "; "))
}

// ValidateJSON validates JSON request body and binds to struct
func ValidateJSON(r *http.Request, dst interface{}) *errors.AppError {
	if r.Header.Get("Content-Type") != "application/json" {
		return errors.BadRequest("Content-Type must be application/json")
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return errors.BadRequest(fmt.Sprintf("Invalid JSON: %s", err.Error()))
	}

	return ValidateStruct(dst)
}

// formatValidationError formats a validation error into a human-readable message
func formatValidationError(err validator.FieldError) string {
	field := err.Field()
	tag := err.Tag()
	param := err.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, param)
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", field, param)
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", field, param)
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, param)
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, param)
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, param)
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, param)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, param)
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "alpha":
		return fmt.Sprintf("%s must contain only alphabetic characters", field)
	case "alphanum":
		return fmt.Sprintf("%s must contain only alphanumeric characters", field)
	case "numeric":
		return fmt.Sprintf("%s must be numeric", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// Common validation structs
type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type UpdateUserRequest struct {
	Username string `json:"username,omitempty" validate:"omitempty,min=3,max=50,alphanum"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
}
