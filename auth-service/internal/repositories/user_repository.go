package repositories

import (
	"log"
	"time"

	"github.com/Koshsky/subs-service/auth-service/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository struct{ DB *gorm.DB }

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (ur *UserRepository) CreateUser(user *models.User) error {
	startTime := time.Now()
	log.Printf("[USER_REPO] [%s] Starting CreateUser for email: %s", startTime.Format("15:04:05.000"), user.Email)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		totalDuration := time.Since(startTime)
		log.Printf("[USER_REPO] [%s] Password hashing FAILED after %v: %v", time.Now().Format("15:04:05.000"), totalDuration, err)
		return err
	}

	user.Password = string(hashedPassword)

	dbErr := ur.DB.Create(user).Error
	if dbErr != nil {
		totalDuration := time.Since(startTime)
		log.Printf("[USER_REPO] [%s] CreateUser FAILED after %v (database error: %v)", time.Now().Format("15:04:05.000"), totalDuration, dbErr)
		return dbErr
	}

	totalDuration := time.Since(startTime)
	log.Printf("[USER_REPO] [%s] CreateUser SUCCESS in %v", time.Now().Format("15:04:05.000"), totalDuration)

	return nil
}

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
