package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email         string         `json:"email" gorm:"unique;not null" validate:"required,min=5"`
	Password      string         `json:"password" gorm:"not null" validate:"required,min=8"`
	Subscriptions []Subscription `json:"subscriptions" gorm:"foreignKey:UserID"`
}
