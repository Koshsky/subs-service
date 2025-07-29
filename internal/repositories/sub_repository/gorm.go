package sub_repository

import (
	"github.com/Koshsky/subs-service/internal/models"
	"gorm.io/gorm"
)

type GormRepository struct {
	*gorm.DB
}

func (gr *GormRepository) GetSubscriptionsByFilters(params models.SubscriptionFilter) ([]models.Subscription, error) {
	var subs []models.Subscription

	query := gr.DB.Model(&models.Subscription{}).
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

func (gr *GormRepository) GetAllSubscriptions() ([]models.Subscription, error) {
	var subs []models.Subscription
	result := gr.DB.Find(&subs)
	return subs, result.Error
}

func (gr *GormRepository) GetSubscriptionByID(id int) (models.Subscription, error) {
	var sub models.Subscription
	result := gr.DB.First(&sub, id)
	return sub, result.Error
}

func (gr *GormRepository) CreateSubscription(sub models.Subscription) (models.Subscription, error) {
	result := gr.DB.Create(&sub)
	return sub, result.Error
}

func (gr *GormRepository) UpdateSubscription(id int, updatedSub models.Subscription) (models.Subscription, error) {
	var sub models.Subscription
	if err := gr.DB.First(&sub, id).Error; err != nil {
		return sub, err
	}

	result := gr.DB.Model(&sub).Updates(updatedSub)
	return sub, result.Error
}

func (gr *GormRepository) DeleteSubscription(id int) error {
	result := gr.DB.Delete(&models.Subscription{}, id)
	return result.Error
}
