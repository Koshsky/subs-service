package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Subscription represents a subscription in the core service database
type Subscription struct {
	gorm.Model
	Service   string     `json:"service_name" gorm:"column:service_name" binding:"required,min=2"`
	Price     int        `json:"price" gorm:"column:price" binding:"required,min=1"`
	UserID    uuid.UUID  `json:"user_id" gorm:"column:user_id;type:uuid;not null"` // Reference to user from auth-service
	StartDate MonthYear  `json:"start_date" gorm:"column:start_date" binding:"required"`
	EndDate   *MonthYear `json:"end_date" gorm:"column:end_date"`
}

// SubscriptionFilter represents filter parameters for subscription queries
type SubscriptionFilter struct {
	UserID     uuid.UUID `form:"user_id" json:"user_id"`
	Service    string    `form:"service" json:"service"`
	StartMonth MonthYear `form:"start_month" json:"start_month" binding:"required"`
	EndMonth   MonthYear `form:"end_month" json:"end_month" binding:"required"`
}
