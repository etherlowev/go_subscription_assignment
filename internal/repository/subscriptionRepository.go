package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"onlineSubscriptions/internal/models"
)

type SubscriptionRepository interface {
	GetById(ctx context.Context, id uuid.UUID) (*models.Subscription, error)

	ExistsById(ctx context.Context, id uuid.UUID) (bool, error)

	Insert(ctx context.Context, uuid uuid.UUID, subscription *models.SubscriptionRequest) (bool, error)

	Update(ctx context.Context, subId uuid.UUID, subscription *models.SubscriptionRequest) (bool, error)

	DeleteById(ctx context.Context, id uuid.UUID) (bool, error)

	Page(ctx context.Context, limit int, offset int) (*[]models.Subscription, error)
}

type PostgresSubscriptionRepository struct {
	Pool *pgxpool.Pool
}

func (repo *PostgresSubscriptionRepository) GetById(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	rows, err := repo.Pool.Query(ctx, "select * from subscriptions where id = $1", id)
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

func (repo *PostgresSubscriptionRepository) ExistsById(ctx context.Context, id uuid.UUID) (bool, error) {
	rows, err := repo.Pool.Query(ctx, "select id from subscriptions where id = $1 limit 1", id)
	if err != nil {
		return false, err
	}

	return rows.Next(), nil
}

func (repo *PostgresSubscriptionRepository) Insert(ctx context.Context, uuid uuid.UUID, subscription *models.SubscriptionRequest) (bool, error) {
	tag, err := repo.Pool.Exec(
		ctx,
		"insert into subscriptions (id, name, price, user_id, start_date, end_date) values ($1, $2, $3, $4, $5, $6)",
		uuid,
		&subscription.Name,
		&subscription.Price,
		&subscription.UserId,
		&subscription.StartDate,
		&subscription.EndDate,
	)
	if err != nil {
		return false, err
	}

	added := tag.RowsAffected() > 0
	return added, nil
}

func (repo *PostgresSubscriptionRepository) Update(ctx context.Context, subId uuid.UUID, subscription *models.SubscriptionRequest) (bool, error) {
	tag, err := repo.Pool.Exec(
		ctx,
		"update subscriptions set name=$1, price=$2, user_id=$3, start_date=$4, end_date=$5 where id = $6",
		&subscription.Name,
		&subscription.Price,
		&subscription.UserId,
		&subscription.StartDate,
		&subscription.EndDate,
		subId,
	)
	if err != nil {
		return false, err
	}

	added := tag.RowsAffected() > 0
	return added, nil
}

func (repo *PostgresSubscriptionRepository) DeleteById(ctx context.Context, id uuid.UUID) (bool, error) {
	tag, err := repo.Pool.Exec(ctx, "delete from subscription where id = $1", id)

	if err != nil {
		return false, err
	}

	added := tag.RowsAffected() > 0
	return added, nil
}

func (repo *PostgresSubscriptionRepository) Page(ctx context.Context, limit int, offset int) (*[]models.Subscription, error) {
	rows, err := repo.Pool.Query(ctx, "select * from subscriptions limit $1 offset $2", limit, offset)
	if err != nil {
		return nil, err
	}

	var results []models.Subscription

	for {
		if rows.Next() {
			var sub models.Subscription
			scanErr := rows.Scan(&sub.Id, &sub.Name, &sub.Price, &sub.UserId, &sub.StartDate, &sub.EndDate)
			if scanErr != nil {
				return nil, scanErr
			}
			results = append(results, sub)
		} else {
			break
		}
	}

	return &results, nil
}
