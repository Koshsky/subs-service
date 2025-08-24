package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Koshsky/subs-service/core-service/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SubscriptionService defines the operations controller requires
// Placed here to depend on behavior, not concrete implementation
// This allows services to implement this contract in their own package
// without controllers needing to import service implementations
// The interface remains small and focused on controller needs
// and can evolve with controller requirements
// Consumers provide an implementation at wiring time
// Note: keep types in shared models package
// to avoid circular deps
type SubscriptionService interface {
	Create(sub models.Subscription) (models.Subscription, error)
	GetByID(id int) (models.Subscription, error)
	GetUserSubscriptions(userID uuid.UUID) ([]models.Subscription, error)
	UpdateByID(id int, update models.Subscription) (models.Subscription, error)
	DeleteByID(id int) error
}

type SubscriptionController struct{ SubService SubscriptionService }

func NewSubscriptionController(service SubscriptionService) *SubscriptionController {
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
			"GetError": "invalid request body",
			"details":  err.Error(),
		})
		return
	}

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"GetError": "invalid user ID",
			"details":  err.Error(),
		})
		return
	}
	sub.UserID = userID

	sub, err = c.SubService.Create(sub)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"GetError": "failed to create subscription",
			"details":  err.Error(),
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
			"GetError": "invalid id format",
			"details":  err.Error(),
		})
		return
	}

	sub, err := c.SubService.GetByID(id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"GetError": "not found",
			"details":  err.Error(),
		})
		return
	}
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"GetError": "invalid user ID",
			"details":  err.Error(),
		})
		return
	}
	if sub.UserID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"GetError": "forbidden",
			"details":  "you are not allowed to access this resource",
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
			"GetError": "invalid id format",
			"details":  err.Error(),
		})
		return
	}

	var inputSub models.Subscription
	if err := ctx.ShouldBindJSON(&inputSub); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"GetError": "invalid request body",
			"details":  err.Error(),
		})
		return
	}

	sub, err := c.SubService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"GetError": "not found",
			"details":  err.Error(),
		})
		return
	}
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"GetError": "invalid user ID",
			"details":  err.Error(),
		})
		return
	}
	if sub.UserID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"GetError": "forbidden",
			"details":  "you are not allowed to access this resource",
		})
		return
	}

	updatedSub, err := c.SubService.UpdateByID(id, inputSub)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"GetError": "failed to update subscription",
			"details":  err.Error(),
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
			"GetError": "invalid user ID",
			"details":  err.Error(),
		})
		return
	}
	subs, err := c.SubService.GetUserSubscriptions(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"GetError": "failed to get subscriptions",
			"details":  err.Error(),
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
			"GetError": "invalid id format",
			"details":  err.Error(),
		})
		return
	}

	sub, err := c.SubService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"GetError": "not found",
			"details":  err.Error(),
		})
		return
	}
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"GetError": "invalid user ID",
			"details":  err.Error(),
		})
		return
	}
	if sub.UserID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"GetError": "forbidden",
			"details":  "you are not allowed to access this resource",
		})
		return
	}

	if err := c.SubService.DeleteByID(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"GetError": "failed to delete subscription",
			"details":  err.Error(),
		})
		return
	}
	ctx.Set("db_affected_id", id)
	ctx.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}
