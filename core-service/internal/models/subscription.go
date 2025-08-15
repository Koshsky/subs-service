package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subscription struct {
	gorm.Model
	Service   string     `json:"service_name" gorm:"column:service_name" binding:"required,min=2"`
	Price     int        `json:"price" gorm:"column:price" binding:"required,min=1"`
	UserID    uuid.UUID  `json:"user_id" gorm:"column:user_id;type:uuid;not null"`
	StartDate MonthYear  `json:"start_date" gorm:"column:start_date" binding:"required"`
	EndDate   *MonthYear `json:"end_date" gorm:"column:end_date"`
}
