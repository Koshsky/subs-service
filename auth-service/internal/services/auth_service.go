package services

import (
	"context"
	"log"
	"time"

	"github.com/Koshsky/subs-service/auth-service/internal/models"
	"github.com/Koshsky/subs-service/auth-service/internal/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	UserRepo        *repositories.UserRepository
	RabbitMQService *RabbitMQService
	JWTSecret       []byte
	Validator       *validator.Validate
}

func NewAuthService(repo *repositories.UserRepository, rabbitmqService *RabbitMQService, jwtSecret []byte) *AuthService {
	return &AuthService{
		UserRepo:        repo,
		RabbitMQService: rabbitmqService,
		JWTSecret:       jwtSecret,
		Validator:       validator.New(),
	}
}

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

	// Publish user created event
	if s.RabbitMQService != nil {
		if err := s.RabbitMQService.PublishUserCreated(user); err != nil {
			log.Printf("Failed to publish user.created event: %v", err)
			// Don't return error as main operation was successful
		}
	}

	user.Password = ""
	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, *models.User, error) {
	user, err := s.UserRepo.ValidateUser(email, password)
	if err != nil {
		return "", nil, err
	}

	token, err := s.GenerateJWTToken(user)
	if err != nil {
		return "", nil, err
	}

	// Remove password from response
	user.Password = ""
	return token, user, nil
}

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

func (s *AuthService) GenerateJWTToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"email":   user.Email,
		"user_id": user.ID.String(),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.JWTSecret)
}
