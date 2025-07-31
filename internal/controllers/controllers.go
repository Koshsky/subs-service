package controllers

import (
	"net/http"
	"strconv"

	"github.com/Koshsky/subs-service/internal/models"
	"github.com/Koshsky/subs-service/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Create(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var sub models.Subscription
		if err := ctx.ShouldBindJSON(&sub); err != nil {
			ctx.JSON(http.StatusBadRequest, models.Error{
				Error:   "invalid request body",
				Code:    models.ErrCodeInvalidRequest,
				Details: err.Error(),
			})
			return
		}
		sub, err := services.CreateSub(db, sub)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, models.Error{
				Error:   "failed to create subscription",
				Code:    models.ErrCodeDatabaseOperation,
				Details: err.Error(),
			})
			return
		}
		ctx.Set("db_affected_id", sub.ID)
		ctx.JSON(http.StatusCreated, sub)
	}
}

func Get(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, models.Error{
				Error:   "invalid id format",
				Code:    models.ErrCodeInvalidID,
				Details: err.Error(),
			})
			return
		}
		sub, err := services.GetSub(db, id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, models.Error{
				Error:   "not found",
				Code:    models.ErrCodeNotFound,
				Details: err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, sub)
	}
}

func Update(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, models.Error{
				Error:   "invalid id format",
				Code:    models.ErrCodeInvalidID,
				Details: err.Error(),
			})
			return
		}

		var inputSub models.Subscription
		if err := ctx.ShouldBindJSON(&inputSub); err != nil {
			ctx.JSON(http.StatusBadRequest, models.Error{
				Error:   "invalid request body",
				Code:    models.ErrCodeInvalidRequest,
				Details: err.Error(),
			})
			return
		}

		updatedSub, err := services.UpdateSub(db, id, inputSub)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, models.Error{
				Error:   "failed to update subscription",
				Code:    models.ErrCodeDatabaseOperation,
				Details: err.Error(),
			})
			return
		}
		ctx.Set("db_affected_id", updatedSub.ID)
		ctx.JSON(http.StatusOK, updatedSub)
	}
}

func List(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		subs, err := services.GetAllSubs(db)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, models.Error{
				Error:   "failed to get subscriptions",
				Code:    models.ErrCodeDatabaseOperation,
				Details: err.Error(),
			})
			return
		}
		if subs == nil {
			subs = []models.Subscription{}
		}
		ctx.JSON(http.StatusOK, subs)
	}
}

func Delete(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, models.Error{
				Error:   "invalid id format",
				Code:    models.ErrCodeInvalidID,
				Details: err.Error(),
			})
			return
		}
		if err := services.DeleteSub(db, id); err != nil {
			ctx.JSON(http.StatusInternalServerError, models.Error{
				Error:   "failed to delete subscription",
				Code:    models.ErrCodeDatabaseOperation,
				Details: err.Error(),
			})
			return
		}
		ctx.Set("db_affected_id", id)
		ctx.Status(http.StatusNoContent)
	}
}

func SumPrice(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req models.SubscriptionFilter

		req.UserID = ctx.Query("user_id")
		req.Service = ctx.Query("service")

		var err error
		err = req.StartMonth.UnmarshalJSON([]byte(ctx.Query("start_month")))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, models.Error{
				Error:   "invalid start month format",
				Code:    models.ErrCodeInvalidDate,
				Details: err.Error(),
			})
			return
		}
		err = req.EndMonth.UnmarshalJSON([]byte(ctx.Query("end_month")))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, models.Error{
				Error:   "invalid end month format",
				Code:    models.ErrCodeInvalidDate,
				Details: err.Error(),
			})
			return
		}

		sum, err := services.SumPrice(db, req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, models.Error{
				Error:   "failed to calculate total price",
				Code:    models.ErrCodeDatabaseOperation,
				Details: err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"total_price": sum})
	}
}
