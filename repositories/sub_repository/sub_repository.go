package sub_repository

import (
	"gorm.io/gorm"

	"github.com/Koshsky/subs-service/models"
)

func GetSubscriptionsByFilters(db *gorm.DB, params models.SubscriptionFilter) ([]models.Subscription, error) {
	var subs []models.Subscription

	query := db.Model(&models.Subscription{}).
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

func GetAllSubscriptions(db *gorm.DB) ([]models.Subscription, error) {
	var subs []models.Subscription
	result := db.Find(&subs)
	return subs, result.Error
}

func GetSubscriptionByID(db *gorm.DB, id int) (models.Subscription, error) {
	var sub models.Subscription
	result := db.First(&sub, id)
	return sub, result.Error
}

func CreateSubscription(db *gorm.DB, sub models.Subscription) (models.Subscription, error) {
	result := db.Create(&sub)
	return sub, result.Error
}

func UpdateSubscription(db *gorm.DB, id int, updatedSub models.Subscription) (models.Subscription, error) {
	var sub models.Subscription
	if err := db.First(&sub, id).Error; err != nil {
		return sub, err
	}

	result := db.Model(&sub).Updates(updatedSub)
	return sub, result.Error
}

func DeleteSubscription(db *gorm.DB, id int) error {
	result := db.Delete(&models.Subscription{}, id)
	return result.Error
}
