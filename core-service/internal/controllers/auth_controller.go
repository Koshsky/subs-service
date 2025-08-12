package controllers

import (
	"context"
	"net/http"

	"github.com/Koshsky/subs-service/core-service/internal/corepb"
	"github.com/gin-gonic/gin"
)

// AuthClient defines methods controller needs from auth client
// Matching the concrete client signatures for simple wiring
type AuthClient interface {
	Register(ctx context.Context, email, password string) (*corepb.RegisterResponse, error)
	Login(ctx context.Context, email, password string) (*corepb.LoginResponse, error)
}

type AuthController struct {
	AuthClient AuthClient
}

func NewAuthController(authClient AuthClient) *AuthController {
	return &AuthController{
		AuthClient: authClient,
	}
}

// Register handles user registration via gRPC
func (ac *AuthController) Register(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	resp, err := ac.AuthClient.Register(c.Request.Context(), credentials.Email, credentials.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to register user",
			"details": err.Error(),
		})
		return
	}

	if !resp.Success {
		c.JSON(http.StatusConflict, gin.H{
			"error":   resp.Error,
			"details": resp.Error,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": resp.Message,
		"user": gin.H{
			"id":    resp.UserId,
			"email": resp.Email,
		},
	})
}

// Login handles user authentication via gRPC
func (ac *AuthController) Login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid credentials format",
			"details": err.Error(),
		})
		return
	}

	resp, err := ac.AuthClient.Login(c.Request.Context(), credentials.Email, credentials.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to authenticate",
			"details": err.Error(),
		})
		return
	}

	if !resp.Success {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   resp.Error,
			"details": resp.Error,
		})
		return
	}

	c.SetCookie("auth_token", resp.Token, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": resp.Message,
	})
}
