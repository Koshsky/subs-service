package middleware

import (
	"net/http"

	"github.com/Koshsky/subs-service/core-service/internal/services"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a middleware that validates the token
func AuthMiddleware(authClient *services.AuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("auth_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization required",
			})
			return
		}

		resp, err := authClient.ValidateToken(c.Request.Context(), tokenString)
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
		c.Set("user_id", int(resp.UserId))
		c.Next()
	}
}
