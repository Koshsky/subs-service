package middleware

import (
	"context"
	"net/http"

	"github.com/Koshsky/subs-service/core-service/internal/corepb"
	"github.com/gin-gonic/gin"
)

// ValidateTokenFunc is a function that validates a token and returns user info
type ValidateTokenFunc func(ctx context.Context, token string) (*corepb.UserResponse, error)

// AuthMiddleware is a middleware that validates the token
func AuthMiddleware(validateToken ValidateTokenFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("auth_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization required",
			})
			return
		}

		resp, err := validateToken(c.Request.Context(), tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid token",
				"details": err.Error(),
			})
			return
		}

		if !resp.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid token",
				"details": resp.Error,
			})
			return
		}

		c.Set("email", resp.Email)
		c.Set("user_id", resp.UserId)
		c.Next()
	}
}
