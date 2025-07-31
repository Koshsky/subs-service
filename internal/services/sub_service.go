package services

import (
	"github.com/Koshsky/subs-service/internal/models"
	"github.com/Koshsky/subs-service/internal/repositories/sub_repository"
)

type SubService struct{ SubRepo *sub_repository.SubRepository }

func NewSubService(repo *sub_repository.SubRepository) *SubService {
	return &SubService{SubRepo: repo}
}

func (s *SubService) Create(sub models.Subscription) (models.Subscription, error) {
	return s.SubRepo.Create(sub)
}

func (s *SubService) GetByID(id int) (models.Subscription, error) {
	return s.SubRepo.GetByID(id)
}

func (s *SubService) GetAll() ([]models.Subscription, error) {
	return s.SubRepo.GetAll()
}

func (s *SubService) UpdateByID(id int, update models.Subscription) (models.Subscription, error) {
	return s.SubRepo.UpdateByID(id, update)
}

func (s *SubService) DeleteByID(id int) error {
	return s.SubRepo.DeletByID(id)
}

func (s *SubService) SumPrice(params models.SubscriptionFilter) (int, error) {
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
