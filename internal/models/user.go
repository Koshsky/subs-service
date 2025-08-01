package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email         string         `json:"email" gorm:"unique;not null" validate:"required,min=5"`
	Password      string         `json:"password" gorm:"not null" validate:"required,min=8"`
	Subscriptions []Subscription `json:"subscriptions" gorm:"foreignKey:UserID"`
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}
