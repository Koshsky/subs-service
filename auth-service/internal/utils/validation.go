package utils

import (
	"unicode"

	"github.com/go-playground/validator/v10"
)

// ValidatePassword validates password complexity requirements
func ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Check minimum length
	if len(password) < 10 {
		return false
	}

	// Check for lowercase letters
	hasLower := false
	// Check for uppercase letters
	hasUpper := false
	// Check for special characters
	hasSpecial := false

	for _, char := range password {
		if unicode.IsLower(char) {
			hasLower = true
		}
		if unicode.IsUpper(char) {
			hasUpper = true
		}
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			hasSpecial = true
		}
	}

	return hasLower && hasUpper && hasSpecial
}

// RegisterCustomValidations registers custom validations
func RegisterCustomValidations(v *validator.Validate) error {
	return v.RegisterValidation("password", ValidatePassword)
}

// NewValidator creates a new validator with custom validations
func NewValidator() *validator.Validate {
	v := validator.New()
	if err := RegisterCustomValidations(v); err != nil {
		// Log error but don't panic - validator will still work without custom validations
		// In production, you might want to handle this differently
	}
	return v
}
