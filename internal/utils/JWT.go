package utils

import (
	"time"

	"github.com/Koshsky/subs-service/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type JWTTokenManager struct {
	jwtSecret []byte
}

func NewJWTTokenManager() *JWTTokenManager {
	return &JWTTokenManager{
		jwtSecret: []byte("your_jwt_secret_key"), // TODO: move to config (with .env) [[ or generate in-place????]]
	}
}

func (j *JWTTokenManager) GenerateJWTToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"email":   user.Email,
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.jwtSecret)
}

func (j *JWTTokenManager) ValidateToken(sessionToken string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(sessionToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return j.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	return claims, nil
}
