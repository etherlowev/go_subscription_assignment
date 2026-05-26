package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"onlineSubscriptions/internal/models"
	repositories "onlineSubscriptions/internal/repository"
)

type SubscriptionService interface {
	GetSubscriptionById(ctx context.Context, id uuid.UUID) (*models.Subscription, error)

	SubscriptionExistsById(ctx context.Context, id uuid.UUID) (bool, error)

	InsertSubscription(ctx context.Context, subscription *models.SubscriptionRequest) (bool, error)

	UpdateSubscription(ctx context.Context, subId uuid.UUID, subscription *models.SubscriptionRequest) (bool, error)

	RemoveSubscription(ctx context.Context, id uuid.UUID) (bool, error)

	ListSubscriptions(ctx context.Context, page int, perPage int) (*[]models.Subscription, error)
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

func (service *PostgresSubscriptionService) GetSubscriptionById(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	return service.Repository.GetById(ctx, id)
}

func (service *PostgresSubscriptionService) SubscriptionExistsById(ctx context.Context, id uuid.UUID) (bool, error) {
	return service.Repository.ExistsById(ctx, id)
}

func (service *PostgresSubscriptionService) InsertSubscription(ctx context.Context, subscription *models.SubscriptionRequest) (bool, error) {
	subId, err := uuid.NewUUID()
	if err != nil {
		return false, err
	}
	return service.Repository.Insert(ctx, subId, subscription)
}

func (service *PostgresSubscriptionService) UpdateSubscription(ctx context.Context, subId uuid.UUID, request *models.SubscriptionRequest) (bool, error) {
	if subId == uuid.Nil {
		return false, nil
	}

	exists, err := service.Repository.ExistsById(ctx, subId)
	if err != nil {
		return false, err
	}

	if exists {
		return service.Repository.Update(ctx, subId, request)
	}
	return false, nil
}

func (service *PostgresSubscriptionService) RemoveSubscription(ctx context.Context, id uuid.UUID) (bool, error) {
	return service.Repository.DeleteById(ctx, id)
}

func (service *PostgresSubscriptionService) ListSubscriptions(ctx context.Context, page int, perPage int) (*[]models.Subscription, error) {
	limit, offset := perPage, max(0, page-1)*perPage
	return service.Repository.Page(ctx, limit, offset)
}
