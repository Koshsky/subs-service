package controllers

import (
	"net/http"
	"strconv"

	"github.com/Koshsky/subs-service/models"
	"github.com/Koshsky/subs-service/services"
	"github.com/gin-gonic/gin"
)

type SubscriptionController struct {
	service *services.SubscriptionService
}

func NewSubscriptionController(service *services.SubscriptionService) *SubscriptionController {
	return &SubscriptionController{service: service}
}

func (c *SubscriptionController) Create(ctx *gin.Context) {
	var sub models.Subscription
	if err := ctx.ShouldBindJSON(&sub); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Error{
			Error:   "invalid request body",
			Code:    models.ErrCodeInvalidRequest,
			Details: err.Error(),
		})
		return
	}
	if err := c.service.CreateSub(ctx, &sub); err != nil {
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

func (c *SubscriptionController) Get(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Error{
			Error:   "invalid id format",
			Code:    models.ErrCodeInvalidID,
			Details: err.Error(),
		})
		return
	}
	sub, err := c.service.GetSub(ctx, id)
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

func (c *SubscriptionController) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Error{
			Error:   "invalid id format",
			Code:    models.ErrCodeInvalidID,
			Details: err.Error(),
		})
		return
	}

	var inputSub models.SubscriptionUpdate
	if err := ctx.ShouldBindJSON(&inputSub); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Error{
			Error:   "invalid request body",
			Code:    models.ErrCodeInvalidRequest,
			Details: err.Error(),
		})
		return
	}

	updatedSub, err := c.service.UpdateSub(ctx, id, inputSub)
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

func (c *SubscriptionController) List(ctx *gin.Context) {
	subs, err := c.service.GetAllSubs(ctx)
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

func (c *SubscriptionController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Error{
			Error:   "invalid id format",
			Code:    models.ErrCodeInvalidID,
			Details: err.Error(),
		})
		return
	}
	if err := c.service.DeleteSub(ctx, id); err != nil {
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

func (c *SubscriptionController) SumPrice(ctx *gin.Context) {
	var req models.SumPriceParams

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

	sum, err := c.service.SumPrice(ctx, req)
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
