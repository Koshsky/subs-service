package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type TokenValidator interface {
	ValidateToken(string) (jwt.MapClaims, error)
}

func AuthMiddleware(TokenValidator TokenValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		tokenString, err := c.Cookie("auth_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization required",
			})
			return
		}

		claims, err := TokenValidator.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid token",
				"details": err.Error(),
			})
			return
		}

		c.Set("username", claims["username"])
		c.Set("user_id", claims["user_id"])
		c.Next()
	}
}
