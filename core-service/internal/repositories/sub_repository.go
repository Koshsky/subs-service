package repositories

import (
	"github.com/Koshsky/subs-service/core-service/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionRepository struct{ DB *gorm.DB }

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{DB: db}
}

// GetUserSubscriptions gets user subscriptions
func (sr *SubscriptionRepository) GetUserSubscriptions(userID uuid.UUID) ([]models.Subscription, error) {
	var subs []models.Subscription
	result := sr.DB.Where("user_id = ?", userID).Find(&subs)
	return subs, result.Error
}

// GetByID gets a subscription by id
func (sr *SubscriptionRepository) GetByID(id uint) (models.Subscription, error) {
	var sub models.Subscription
	result := sr.DB.First(&sub, id)
	return sub, result.Error
}

// Create creates a new subscription
func (sr *SubscriptionRepository) Create(sub models.Subscription) (models.Subscription, error) {
	result := sr.DB.Create(&sub)
	return sub, result.Error
}

// UpdateByID updates a subscription by id
func (sr *SubscriptionRepository) UpdateByID(id uint, updatedSub models.Subscription) (models.Subscription, error) {
	var sub models.Subscription
	if err := sr.DB.First(&sub, id).Error; err != nil {
		return sub, err
	}

	result := sr.DB.Model(&sub).Updates(updatedSub)
	return sub, result.Error
}

// DeleteByID deletes a subscription by id
func (sr *SubscriptionRepository) DeleteByID(id uint) error {
	result := sr.DB.Delete(&models.Subscription{}, id)
	return result.Error
}
