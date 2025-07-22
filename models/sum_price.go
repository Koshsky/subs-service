package models

type SumPriceParams struct {
	UserID     string    `form:"user_id" json:"user_id"`
	Service    string    `form:"service" json:"service"`
	StartMonth MonthYear `form:"start_month" json:"start_month" binding:"required"`
	EndMonth   MonthYear `form:"end_month" json:"end_month" binding:"required"`
}
