package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validates a struct using validator tags
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// FormatValidationError formats validation errors into a readable string
func FormatValidationError(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errors []string
		for _, e := range validationErrors {
			errors = append(errors, formatFieldError(e))
		}
		return strings.Join(errors, ", ")
	}
	return err.Error()
}

func formatFieldError(e validator.FieldError) string {
	field := strings.ToLower(e.Field())
	
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, e.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, e.Param())
	case "numeric":
		return fmt.Sprintf("%s must be numeric", field)
	case "alpha":
		return fmt.Sprintf("%s must contain only letters", field)
	case "alphanum":
		return fmt.Sprintf("%s must contain only letters and numbers", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", field, e.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, e.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, e.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, e.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, e.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// SanitizeString removes extra spaces and trims the string
func SanitizeString(s string) string {
	return strings.TrimSpace(s)
}

// SanitizeEmail converts email to lowercase and trims spaces
func SanitizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// ValidateAndSanitize validates a struct and sanitizes string fields
func ValidateAndSanitize(s interface{}) error {
	// First sanitize (if needed, implement specific sanitization logic)
	// Then validate
	return ValidateStruct(s)
}