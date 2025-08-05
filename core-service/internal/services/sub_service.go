package services

import (
	"github.com/Koshsky/subs-service/core-service/internal/models"
	"github.com/Koshsky/subs-service/core-service/internal/repositories"
)

type SubscriptionService struct {
	SubRepo *repositories.SubscriptionRepository
}

func NewSubscriptionService(repo *repositories.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{SubRepo: repo}
}

// Create creates a new subscription
func (s *SubscriptionService) Create(sub models.Subscription) (models.Subscription, error) {
	return s.SubRepo.Create(sub)
}

// GetByID gets a subscription by id
func (s *SubscriptionService) GetByID(id int) (models.Subscription, error) {
	return s.SubRepo.GetByID(uint(id))
}

// GetUserSubscriptions gets user subscriptions
func (s *SubscriptionService) GetUserSubscriptions(userID int) ([]models.Subscription, error) {
	return s.SubRepo.GetUserSubscriptions(uint(userID))
}

// UpdateByID updates a subscription by id
func (s *SubscriptionService) UpdateByID(id int, update models.Subscription) (models.Subscription, error) {
	return s.SubRepo.UpdateByID(uint(id), update)
}

// DeleteByID deletes a subscription by id
func (s *SubscriptionService) DeleteByID(id int) error {
	return s.SubRepo.DeleteByID(uint(id))
}

// SumPrice sums the price of subscriptions for a user
func (s *SubscriptionService) SumPrice(params models.SubscriptionFilter) (int, error) {
	subs, err := s.SubRepo.GetBySubscriptionFilter(params)
	if err != nil {
		return 0, err
	}

	var total int
	for _, sub := range subs {
		total += sub.Price
	}

	return total, nil
}
