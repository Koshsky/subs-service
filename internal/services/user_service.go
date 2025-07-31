package services

import (
	"time"

	"github.com/Koshsky/subs-service/internal/models"
	"github.com/Koshsky/subs-service/internal/repositories/user_repository"
	"github.com/golang-jwt/jwt/v5"
)

type UserService struct {
	UserRepo  *user_repository.UserRepository
	jwtSecret []byte
}

func NewUserService(repo *user_repository.UserRepository) *UserService {
	return &UserService{
		UserRepo:  repo,
		jwtSecret: []byte("your_jwt_secret_key"), // TODO: move to config (with .env) [[ or generate in-place????]]
	}
}

// RegisterUser handles user registration (password hashing done in repository)
func (us *UserService) RegisterUser(user *models.User) error {
	return us.UserRepo.CreateUser(user)
}

// ValidateCredentials validates user credentials
func (us *UserService) ValidateCredentials(username, password string) (*models.User, error) {
	return us.UserRepo.ValidateUser(username, password)
}

func (us *UserService) GenerateJWTToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"username": user.Username,
		"user_id":  user.ID,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(us.jwtSecret)
}

func (us *UserService) ValidateToken(sessionToken string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(sessionToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return us.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	return claims, nil
}
