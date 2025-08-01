package services

import (
	"github.com/Koshsky/subs-service/internal/models"
	"github.com/Koshsky/subs-service/internal/repositories/user_repository"
)

type UserService struct {
	UserRepo *user_repository.UserRepository
}

func NewUserService(repo *user_repository.UserRepository) *UserService {
	return &UserService{
		UserRepo: repo,
	}
}

// RegisterUser handles user registration (password hashing done in repository)
func (us *UserService) RegisterUser(user *models.User) error {
	return us.UserRepo.CreateUser(user)
}

// ValidateCredentials validates user credentials
func (us *UserService) ValidateCredentials(email, password string) (*models.User, error) {
	return us.UserRepo.ValidateUser(email, password)
}
