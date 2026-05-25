package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	models "onlineSubscriptions/src/main/impl/model"
)

type SubscriptionRepository interface {
	Find(ctx context.Context, id uuid.UUID) (*models.Subscription, error)

	Insert(ctx context.Context, subscription *models.Subscription) error

	Update(ctx context.Context, subscription *models.Subscription) error

	Delete(ctx context.Context, id uuid.UUID) error

	List(ctx context.Context, page int, perPage int) (*[]models.Subscription, error)
}

type PostgresSubscriptionRepository struct {
	Pool *pgxpool.Pool
}

func (repo *PostgresSubscriptionRepository) Find(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	rows, err := repo.Pool.Query(ctx, "select * from subscriptions where id = $1", id.String())
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		var result models.Subscription
		scanErr := rows.Scan(&result.Id, &result.Name, &result.Price, &result.UserId, &result.StartDate, &result.EndDate)
		if scanErr != nil {
			return nil, scanErr
		}
		return &result, nil
	}

	return nil, nil
}

func (repo *PostgresSubscriptionRepository) Insert(ctx context.Context, subscription *models.Subscription) error {
	return nil
}

func (repo *PostgresSubscriptionRepository) Update(ctx context.Context, subscription *models.Subscription) error {
	return nil
}

func (repo *PostgresSubscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (repo *PostgresSubscriptionRepository) List(ctx context.Context, page int, perPage int) (*[]models.Subscription, error) {
	return nil, nil
}
