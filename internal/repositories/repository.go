package repositories

import (
	"log"

	"github.com/Koshsky/subs-service/internal/models"
	"github.com/Koshsky/subs-service/internal/repositories/sub_repository"
	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	GetSubscriptionsByFilters(params models.SubscriptionFilter) ([]models.Subscription, error)
	GetAllSubscriptions() ([]models.Subscription, error)
	GetSubscriptionByID(id int) (models.Subscription, error)
	CreateSubscription(sub models.Subscription) (models.Subscription, error)
	UpdateSubscription(id int, updatedSub models.Subscription) (models.Subscription, error)
	DeleteSubscription(id int) error
}

func New(db interface{}) SubscriptionRepository {
	switch v := db.(type) {
	case *gorm.DB:
		return &sub_repository.GormRepository{DB: v}
	default:
		log.Fatalf("unsupported database type: %T, expected *gorm.DB", db)
		return nil
	}
}
