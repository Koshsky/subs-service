package services

import (
	"github.com/Koshsky/subs-service/internal/models"

	"github.com/Koshsky/subs-service/internal/repositories"
)

// SubscriptionService provides business logic for subscriptions
// and delegates data access to the repository layer.
type SubscriptionService struct {
	repo repositories.SubscriptionRepository
}

func New(repo repositories.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) CreateSub(sub models.Subscription) (models.Subscription, error) {
	return s.repo.CreateSubscription(sub)
}

func (s *SubscriptionService) GetSub(id int) (models.Subscription, error) {
	return s.repo.GetSubscriptionByID(id)
}

func (s *SubscriptionService) GetAllSubs() ([]models.Subscription, error) {
	return s.repo.GetAllSubscriptions()
}

func (s *SubscriptionService) UpdateSub(id int, update models.Subscription) (models.Subscription, error) {
	return s.repo.UpdateSubscription(id, update)
}

func (s *SubscriptionService) DeleteSub(id int) error {
	return s.repo.DeleteSubscription(id)
}

func (s *SubscriptionService) SumPrice(params models.SubscriptionFilter) (int, error) {
	subs, err := s.repo.GetSubscriptionsByFilters(params)
	if err != nil {
		return 0, err
	}

	var total int
	for _, sub := range subs {
		total += sub.Price
	}

	return total, nil
}
