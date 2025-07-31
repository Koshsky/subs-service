package models

import "gorm.io/gorm"

type Subscription struct {
	gorm.Model
	Service   string     `json:"service_name" gorm:"column:service_name" binding:"required,min=2"`
	Price     int        `json:"price" gorm:"column:price" binding:"required,min=1"`
	UserID    string     `json:"user_id" gorm:"column:user_id" binding:"required,uuid"`
	StartDate MonthYear  `json:"start_date" gorm:"column:start_date" binding:"required"`
	EndDate   *MonthYear `json:"end_date" gorm:"column:end_date"`
}

type SubscriptionFilter struct {
	UserID     string    `form:"user_id" json:"user_id"`
	Service    string    `form:"service" json:"service"`
	StartMonth MonthYear `form:"start_month" json:"start_month" binding:"required"`
	EndMonth   MonthYear `form:"end_month" json:"end_month" binding:"required"`
}
