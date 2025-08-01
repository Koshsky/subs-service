package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTTokenValidator interface {
	ValidateToken(string) (jwt.MapClaims, error)
}

func AuthMiddleware(jwtValidator JWTTokenValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		tokenString, err := c.Cookie("auth_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization required",
			})
			return
		}

		claims, err := jwtValidator.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid token",
				"details": err.Error(),
			})
			return
		}

		userID, ok := claims["user_id"].(float64) // JWT числа всегда float64
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid user ID in token",
			})
			return
		}

		c.Set("email", claims["email"].(string))
		c.Set("user_id", int(userID)) // Преобразуем float64 -> int
		c.Next()
	}
}
