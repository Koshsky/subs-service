package user_repository

import (
	"log"

	"github.com/Koshsky/subs-service/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository struct{ DB *gorm.DB }

func New(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// CreateUser creates a new user with hashed password
func (ur *UserRepository) CreateUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return ur.DB.Create(user).Error
}

// ValidateUser checks credentials and returns user if valid
func (ur *UserRepository) ValidateUser(email, password string) (*models.User, error) {
	var user models.User
	err := ur.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Error while comparing passwords: %v", err)
		return nil, err
	}
	return &user, nil
}
