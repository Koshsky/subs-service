package services

import (
	"context"
	"time"

	"github.com/Koshsky/subs-service/auth-service/internal/repositories"
	"github.com/Koshsky/subs-service/shared/models"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	UserRepo  *repositories.UserRepository
	JWTSecret []byte
	Validator *validator.Validate
}

func NewAuthService(repo *repositories.UserRepository, jwtSecret []byte) *AuthService {
	return &AuthService{
		UserRepo:  repo,
		JWTSecret: jwtSecret,
		Validator: validator.New(),
	}
}

// Register creates a new user and returns the user and error if it's not
func (s *AuthService) Register(ctx context.Context, email, password string) (*models.User, error) {
	user := &models.User{
		Email:    email,
		Password: password,
	}

	if err := s.Validator.Struct(user); err != nil {
		return nil, err
	}

	err := s.UserRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	// Убираем пароль из ответа
	user.Password = ""
	return user, nil
}

// Login logs in a user and returns the token and user and error if it's not
func (s *AuthService) Login(ctx context.Context, email, password string) (string, *models.User, error) {
	user, err := s.UserRepo.ValidateUser(email, password)
	if err != nil {
		return "", nil, err
	}

	token, err := s.GenerateJWTToken(user)
	if err != nil {
		return "", nil, err
	}

	// Убираем пароль из ответа
	user.Password = ""
	return token, user, nil
}

// ValidateToken checks if the token is valid and returns the claims and error if it's not
func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return s.JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	return claims, nil
}

// GenerateJWTToken generates a JWT token for a user
func (s *AuthService) GenerateJWTToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"email":   user.Email,
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.JWTSecret)
}
