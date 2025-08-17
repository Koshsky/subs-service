package repositories

import (
	"fmt"
	"log"

	"github.com/Koshsky/subs-service/auth-service/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB DatabaseInterface
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: NewGormAdapter(db)}
}

func (ur *UserRepository) CreateUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("cannot create user with password=%s: %w", user.Password, err)
	}

	user.Password = string(hashedPassword)

	dbErr := ur.DB.Create(user).GetError()
	if dbErr != nil {
		return fmt.Errorf("cannot create user with email=%s: %w", user.Email, dbErr)
	}

	return nil
}

func (ur *UserRepository) ValidateUser(email, password string) (*models.User, error) {
	user, err := ur.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	bcryptErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if bcryptErr != nil {
		log.Printf("Error while comparing passwords: %v", bcryptErr)
		return nil, fmt.Errorf("authentication failed for user %s: %w", email, bcryptErr)
	}
	return user, nil
}

func (ur *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := ur.DB.Where("email = ?", email).First(&user).GetError()
	if err != nil {
		return nil, fmt.Errorf("cannot get user by email=%s: %w", email, err)
	}
	return &user, nil
}
