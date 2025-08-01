package services

import (
	"github.com/Koshsky/subs-service/internal/models"
	"github.com/Koshsky/subs-service/internal/repositories/sub_repository"
)

type SubscriptionService struct {
	SubRepo *sub_repository.SubscriptionRepository
}

func NewSubscriptionService(repo *sub_repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{SubRepo: repo}
}

func (s *SubscriptionService) Create(sub models.Subscription) (models.Subscription, error) {
	return s.SubRepo.Create(sub)
}

func (s *SubscriptionService) GetByID(id int) (models.Subscription, error) {
	return s.SubRepo.GetByID(uint(id))
}

func (s *SubscriptionService) GetUserSubscriptions(userID int) ([]models.Subscription, error) {
	return s.SubRepo.GetUserSubscriptions(uint(userID))
}

func (s *SubscriptionService) UpdateByID(id int, update models.Subscription) (models.Subscription, error) {
	return s.SubRepo.UpdateByID(uint(id), update)
}

func (s *SubscriptionService) DeleteByID(id int) error {
	return s.SubRepo.DeleteByID(uint(id))
}

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
