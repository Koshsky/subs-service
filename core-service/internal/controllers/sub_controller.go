package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Koshsky/subs-service/core-service/internal/models"
	"github.com/Koshsky/subs-service/core-service/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SubscriptionController struct{ SubService *services.SubscriptionService }

func NewSubscriptionController(service *services.SubscriptionService) *SubscriptionController {
	return &SubscriptionController{SubService: service}
}

// getUserIDFromContext extracts and parses user UUID from gin context
func getUserIDFromContext(ctx *gin.Context) (uuid.UUID, error) {
	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		return uuid.Nil, fmt.Errorf("user_id not found in context")
	}

	userIDString, ok := userIDStr.(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("user_id is not a string")
	}

	return uuid.Parse(userIDString)
}

// Create creates a new subscription
func (c *SubscriptionController) Create(ctx *gin.Context) {
	var sub models.Subscription
	if err := ctx.ShouldBindJSON(&sub); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid user ID",
			"details": err.Error(),
		})
		return
	}
	sub.UserID = userID

	sub, err = c.SubService.Create(sub)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to create subscription",
			"details": err.Error(),
		})
		return
	}
	ctx.Set("db_affected_id", sub.ID)
	ctx.JSON(http.StatusCreated, sub)
}

// Get gets a subscription by id
func (c *SubscriptionController) Get(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid id format",
			"details": err.Error(),
		})
		return
	}

	sub, err := c.SubService.GetByID(id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "not found",
			"details": err.Error(),
		})
		return
	}
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid user ID",
			"details": err.Error(),
		})
		return
	}
	if sub.UserID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error":   "forbidden",
			"details": "you are not allowed to access this resource",
		})
		return
	}

	ctx.JSON(http.StatusOK, sub)
}

// Update updates a subscription by id
func (c *SubscriptionController) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid id format",
			"details": err.Error(),
		})
		return
	}

	var inputSub models.Subscription
	if err := ctx.ShouldBindJSON(&inputSub); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	sub, err := c.SubService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "not found",
			"details": err.Error(),
		})
		return
	}
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid user ID",
			"details": err.Error(),
		})
		return
	}
	if sub.UserID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error":   "forbidden",
			"details": "you are not allowed to access this resource",
		})
		return
	}

	updatedSub, err := c.SubService.UpdateByID(id, inputSub)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to update subscription",
			"details": err.Error(),
		})
		return
	}

	ctx.Set("db_affected_id", updatedSub.ID)
	ctx.JSON(http.StatusOK, updatedSub)
}

// List lists all subscriptions for a user
func (c *SubscriptionController) List(ctx *gin.Context) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid user ID",
			"details": err.Error(),
		})
		return
	}
	subs, err := c.SubService.GetUserSubscriptions(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to get subscriptions",
			"details": err.Error(),
		})
		return
	}
	if subs == nil {
		subs = []models.Subscription{}
	}
	ctx.JSON(http.StatusOK, subs)
}

// Delete deletes a subscription by id
func (c *SubscriptionController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid id format",
			"details": err.Error(),
		})
		return
	}

	sub, err := c.SubService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "not found",
			"details": err.Error(),
		})
		return
	}
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid user ID",
			"details": err.Error(),
		})
		return
	}
	if sub.UserID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error":   "forbidden",
			"details": "you are not allowed to access this resource",
		})
		return
	}

	if err := c.SubService.DeleteByID(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to delete subscription",
			"details": err.Error(),
		})
		return
	}
	ctx.Set("db_affected_id", id)
	ctx.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}

// SumPrice sums the price of subscriptions for a user
func (c *SubscriptionController) SumPrice(ctx *gin.Context) {
	var req models.SubscriptionFilter

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid user ID",
			"details": err.Error(),
		})
		return
	}
	req.UserID = userID
	req.Service = ctx.Query("service")

	err = req.StartMonth.UnmarshalJSON([]byte(ctx.Query("start_month")))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid start month format",
			"details": err.Error(),
		})
		return
	}
	err = req.EndMonth.UnmarshalJSON([]byte(ctx.Query("end_month")))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid end month format",
			"details": err.Error(),
		})
		return
	}

	sum, err := c.SubService.SumPrice(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to calculate total price",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"total_price": sum})
}
