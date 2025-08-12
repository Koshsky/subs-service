package utils

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// validateUUID validates a UUID
func validateUUID(fl validator.FieldLevel) bool {
	fieldValue := fl.Field().String()
	_, err := uuid.Parse(fieldValue)
	if err != nil {
		return false
	}
	return true
}

// RegisterCustomValidations registers custom validations
func RegisterCustomValidations() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return
	}

	if err := v.RegisterValidation("uuid", validateUUID); err != nil {
		panic(err)
	}
}
