package models

type SubscriptionUpdate struct {
	Service   *string    `json:"service_name,omitempty" binding:"omitempty,min=2"`
	Price     *int       `json:"price,omitempty" binding:"omitempty,min=1"`
	UserID    *string    `json:"user_id,omitempty" binding:"omitempty,uuid4"`
	StartDate *MonthYear `json:"start_date,omitempty" binding:"omitempty"`
	EndDate   *MonthYear `json:"end_date,omitempty" binding:"omitempty"`
}
