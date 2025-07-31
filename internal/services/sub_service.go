package services

import (
	"github.com/Koshsky/subs-service/internal/models"
	"github.com/Koshsky/subs-service/internal/repositories/sub_repository"
	"gorm.io/gorm"
)

func CreateSub(db *gorm.DB, sub models.Subscription) (models.Subscription, error) {
	return sub_repository.CreateSubscription(db, sub)
}

func GetSub(db *gorm.DB, id int) (models.Subscription, error) {
	return sub_repository.GetSubscriptionByID(db, id)
}

func GetAllSubs(db *gorm.DB) ([]models.Subscription, error) {
	return sub_repository.GetAllSubscriptions(db)
}

func UpdateSub(db *gorm.DB, id int, update models.Subscription) (models.Subscription, error) {
	return sub_repository.UpdateSubscription(db, id, update)
}

func DeleteSub(db *gorm.DB, id int) error {
	return sub_repository.DeleteSubscription(db, id)
}

func SumPrice(db *gorm.DB, params models.SubscriptionFilter) (int, error) {
	subs, err := sub_repository.GetSubscriptionsByFilters(db, params)
	if err != nil {
		return 0, err
	}

	var total int
	for _, sub := range subs {
		total += sub.Price
	}

	return total, nil
}
