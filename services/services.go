package services

import (
	"github.com/Koshsky/subs-service/models"
	"gorm.io/gorm"

	"github.com/Koshsky/subs-service/repositories/sub_repository"
)

// SubscriptionService provides business logic for subscriptions
// and delegates data access to the repository layer.
type SubscriptionService struct {
	repo *gorm.DB
}

func NewSubscriptionService(repo *gorm.DB) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) CreateSub(sub models.Subscription) (models.Subscription, error) {
	return sub_repository.CreateSubscription(s.repo, sub)
}

func (s *SubscriptionService) GetSub(id int) (models.Subscription, error) {
	return sub_repository.GetSubscriptionByID(s.repo, id)
}

func (s *SubscriptionService) GetAllSubs() ([]models.Subscription, error) {
	return sub_repository.GetAllSubscriptions(s.repo)
}

func (s *SubscriptionService) UpdateSub(id int, update models.Subscription) (models.Subscription, error) {
	return sub_repository.UpdateSubscription(s.repo, id, update)
}

func (s *SubscriptionService) DeleteSub(id int) error {
	return sub_repository.DeleteSubscription(s.repo, id)
}

func (s *SubscriptionService) SumPrice(params models.SubscriptionFilter) (int, error) {
	subs, err := sub_repository.GetSubscriptionsByFilters(s.repo, params)
	if err != nil {
		return 0, err
	}

	var total int
	for _, sub := range subs {
		total += sub.Price
	}

	return total, nil
}
