package models

type Subscription struct {
	ID        int        `json:"id"`
	Service   string     `json:"service_name" db:"service_name" binding:"required,min=2"`
	Price     int        `json:"price" db:"price" binding:"required,min=1"`
	UserID    string     `json:"user_id" db:"user_id" binding:"required,uuid"`
	StartDate MonthYear  `json:"start_date" db:"start_date" binding:"required"`
	EndDate   *MonthYear `json:"end_date" db:"end_date"`
}
