package models

import (
	"gorm.io/gorm"
)

// User represents a user in the auth service database
type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"unique;not null" validate:"required,min=5"`
	Password string `json:"password" gorm:"not null" validate:"required,min=8"`
}
