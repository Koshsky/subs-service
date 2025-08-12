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

func (s *AuthService) Register(ctx context.Context, email, password string) (*models.User, error) {
	startTime := time.Now()
	log.Printf("[AUTH_SERVICE] [%s] Starting Register for email: %s", startTime.Format("15:04:05.000"), email)

	user := &models.User{
		Email:    email,
		Password: password,
	}

	if err := s.Validator.Struct(user); err != nil {
		totalDuration := time.Since(startTime)
		log.Printf("[AUTH_SERVICE] [%s] Validation FAILED after %v: %v", time.Now().Format("15:04:05.000"), totalDuration, err)
		return nil, err
	}

	err := s.UserRepo.CreateUser(user)

	if err != nil {
		totalDuration := time.Since(startTime)
		log.Printf("[AUTH_SERVICE] [%s] Register FAILED after %v (database error: %v)", time.Now().Format("15:04:05.000"), totalDuration, err)
		return nil, err
	}

	user.Password = ""

	totalDuration := time.Since(startTime)
	log.Printf("[AUTH_SERVICE] [%s] Register SUCCESS in %v", time.Now().Format("15:04:05.000"), totalDuration)

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
