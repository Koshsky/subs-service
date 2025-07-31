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

	if params.UserID != "" {
		query = query.Where("user_id = ?", params.UserID)
	}
	if params.Service != "" {
		query = query.Where("service_name = ?", params.Service)
	}

	err := query.Find(&subs).Error
	return subs, err
}

func (sr *SubscriptionRepository) GetAll() ([]models.Subscription, error) {
	var subs []models.Subscription
	result := sr.DB.Find(&subs)
	return subs, result.Error
}

func (sr *SubscriptionRepository) GetByID(id int) (models.Subscription, error) {
	var sub models.Subscription
	result := sr.DB.First(&sub, id)
	return sub, result.Error
}

func (sr *SubscriptionRepository) Create(sub models.Subscription) (models.Subscription, error) {
	result := sr.DB.Create(&sub)
	return sub, result.Error
}

func (sr *SubscriptionRepository) UpdateByID(id int, updatedSub models.Subscription) (models.Subscription, error) {
	var sub models.Subscription
	if err := sr.DB.First(&sub, id).Error; err != nil {
		return sub, err
	}

	result := sr.DB.Model(&sub).Updates(updatedSub)
	return sub, result.Error
}

func (sr *SubscriptionRepository) DeleteByID(id int) error {
	result := sr.DB.Delete(&models.Subscription{}, id)
	return result.Error
}
