package sub_repository

import (
	"github.com/Koshsky/subs-service/internal/models"
	"gorm.io/gorm"
)

type SubscriptionRepository struct{ DB *gorm.DB }

func New(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{DB: db}
}

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

func (sr *SubscriptionRepository) GetUserSubscriptions(userID uint) ([]models.Subscription, error) {
	var subs []models.Subscription
	result := sr.DB.Where("user_id = ?", userID).Find(&subs)
	return subs, result.Error
}

func (sr *SubscriptionRepository) GetByID(id uint) (models.Subscription, error) {
	var sub models.Subscription
	result := sr.DB.First(&sub, id)
	return sub, result.Error
}

func (sr *SubscriptionRepository) Create(sub models.Subscription) (models.Subscription, error) {
	result := sr.DB.Create(&sub)
	return sub, result.Error
}

func (sr *SubscriptionRepository) UpdateByID(id uint, updatedSub models.Subscription) (models.Subscription, error) {
	var sub models.Subscription
	if err := sr.DB.First(&sub, id).Error; err != nil {
		return sub, err
	}

	result := sr.DB.Model(&sub).Updates(updatedSub)
	return sub, result.Error
}

func (sr *SubscriptionRepository) DeleteByID(id uint) error {
	result := sr.DB.Delete(&models.Subscription{}, id)
	return result.Error
}
