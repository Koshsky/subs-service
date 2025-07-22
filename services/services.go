package services

import (
	"context"

	"github.com/Koshsky/subs-service/models"
)

// SubscriptionService provides business logic for subscriptions
// and delegates data access to the repository layer.
type SubscriptionService struct {
	repo SubscriptionRepo
}

func NewSubscriptionService(repo SubscriptionRepo) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) SumPrice(ctx context.Context, params models.SumPriceParams) (float64, error) {
	return s.repo.SumPrice(ctx, params)
}

func (s *SubscriptionService) CreateSub(ctx context.Context, sub *models.Subscription) error {
	return s.repo.Create(ctx, sub)
}

func (s *SubscriptionService) GetSub(ctx context.Context, id int) (*models.Subscription, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *SubscriptionService) GetAllSubs(ctx context.Context) ([]models.Subscription, error) {
	return s.repo.GetAll(ctx)
}

func (s *SubscriptionService) UpdateSub(ctx context.Context, id int, update models.SubscriptionUpdate) (*models.Subscription, error) {
	return s.repo.Update(ctx, id, update)
}

func (s *SubscriptionService) DeleteSub(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
