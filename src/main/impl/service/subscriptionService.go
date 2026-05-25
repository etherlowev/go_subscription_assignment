package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	models "onlineSubscriptions/src/main/impl/model"
	repositories "onlineSubscriptions/src/main/impl/repository"
)

type SubscriptionService interface {
	Find(ctx context.Context, id uuid.UUID) (*models.Subscription, error)

	Insert(ctx context.Context, subscription *models.Subscription) error

	Update(ctx context.Context, subscription *models.Subscription) error

	Delete(ctx context.Context, id uuid.UUID) error

	List(ctx context.Context, page int, perPage int) (*[]models.Subscription, error)
}

type PostgresSubscriptionService struct {
	Repository repositories.SubscriptionRepository
}

func NewSubscriptionService(Pool *pgxpool.Pool) *PostgresSubscriptionService {
	repo := &repositories.PostgresSubscriptionRepository{
		Pool: Pool,
	}

	service := &PostgresSubscriptionService{
		Repository: repo,
	}

	return service
}

func (service *PostgresSubscriptionService) Find(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	return service.Repository.Find(ctx, id)
}

func (service *PostgresSubscriptionService) Insert(ctx context.Context, subscription *models.Subscription) error {
	return service.Repository.Insert(ctx, subscription)
}

func (service *PostgresSubscriptionService) Update(ctx context.Context, subscription *models.Subscription) error {
	return service.Repository.Update(ctx, subscription)
}

func (service *PostgresSubscriptionService) Delete(ctx context.Context, id uuid.UUID) error {
	return service.Repository.Delete(ctx, id)
}

func (service *PostgresSubscriptionService) List(ctx context.Context, page int, perPage int) (*[]models.Subscription, error) {
	return service.Repository.List(ctx, page, perPage)
}
