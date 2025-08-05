package utils

import (
	"log"
	"reflect"

	"github.com/Koshsky/subs-service/shared/models"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// validateUUID validates a UUID
func validateUUID(fl validator.FieldLevel) bool {
	fieldValue := fl.Field().String()
	_, err := uuid.Parse(fieldValue)
	if err != nil {
		log.Printf("UUID validation failed for field '%s': %v", fl.FieldName(), err)
		return false
	}
	return true
}

// validateMonthYear validates a MonthYear
func validateMonthYear(fl validator.FieldLevel) bool {
	field := fl.Field()
	fieldType := field.Type()

	if fieldType == reflect.TypeOf(models.MonthYear{}) {
		return true
	}

	if field.Kind() == reflect.Ptr && fieldType == reflect.TypeOf(&models.MonthYear{}) {
		return field.IsNil() || true // nil is considered valid
	}

	log.Printf("MonthYear validation failed for field '%s': invalid type %v", fl.FieldName(), fieldType)
	return false
}

// RegisterCustomValidations registers custom validations
func RegisterCustomValidations() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		log.Println("Failed to get Gin validator engine")
		return
	}

	log.Println("Registering custom validations...")

	if err := v.RegisterValidation("uuid", validateUUID); err != nil {
		log.Printf("Failed to register UUID validator: %v", err)
	} else {
		log.Println("UUID validator registered successfully")
	}

	if err := v.RegisterValidation("monthyear", validateMonthYear); err != nil {
		log.Printf("Failed to register MonthYear validator: %v", err)
	} else {
		log.Println("MonthYear validator registered successfully")
	}

	log.Println("Custom validations registration completed")
}
