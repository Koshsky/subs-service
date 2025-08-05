package repositories

import (
	"gorm.io/gorm"

	"github.com/Koshsky/subs-service/core-service/internal/models"
)

type SubscriptionRepository struct{ DB *gorm.DB }

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{DB: db}
}

// GetBySubscriptionFilter gets subscriptions by subscription filter
func (sr *SubscriptionRepository) GetBySubscriptionFilter(params models.SubscriptionFilter) ([]models.Subscription, error) {
	var subs []models.Subscription

	query := sr.DB.Model(&models.Subscription{}).
		Where("start_date BETWEEN ? AND ?",
			params.StartMonth.Time(),
			params.EndMonth.Time().AddDate(0, 1, -1)) // До конца месяца

	if params.UserID != 0 {
		query = query.Where("user_id = ?", params.UserID)
	}
	if params.Service != "" {
		query = query.Where("service_name = ?", params.Service)
	}

	err := query.Find(&subs).Error
	return subs, err
}

// GetUserSubscriptions gets user subscriptions
func (sr *SubscriptionRepository) GetUserSubscriptions(userID uint) ([]models.Subscription, error) {
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
