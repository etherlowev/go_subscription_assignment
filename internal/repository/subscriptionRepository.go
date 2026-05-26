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

	CalculateSubscriptionSum(ctx context.Context, userId string, subName string, dateFrom string, dateTo string) (*models.SubscriptionPriceSum, error)
}

type PostgresSubscriptionRepository struct {
	Pool *pgxpool.Pool
}

func (repo *PostgresSubscriptionRepository) GetById(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	rows, err := repo.Pool.Query(ctx, "select id, name, price, user_id, "+
		"TO_CHAR(start_date, 'mm-YYYY') as start_date, TO_CHAR(end_date, 'mm-YYYY') as end_date "+
		"from subscriptions where id = $1", id)

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
	rows, err := repo.Pool.Query(ctx, "select id, from subscriptions where id = $1 limit 1", id)

	if err != nil {
		return false, err
	}

	return rows.Next(), nil
}

func (repo *PostgresSubscriptionRepository) Insert(ctx context.Context, uuid uuid.UUID, subscription *models.SubscriptionRequest) (bool, error) {
	tag, err := repo.Pool.Exec(
		ctx,
		"insert into subscriptions (id, name, price, user_id, start_date, end_date) "+
			"values ($1, $2, $3, $4, TO_DATE($5, 'MM-YYYY'), case when $6 = '' then null else TO_DATE($6, 'MM-YYYY') end)",
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
		"update subscriptions set name=$1, price=$2, user_id=$3, "+
			"start_date=TO_DATE($4, 'MM-YYYY'), end_date=TO_DATE($5, 'MM-YYYY') where id = $6",
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
	rows, err := repo.Pool.Query(ctx, "select id, name, price, "+
		"user_id, TO_CHAR(start_date, 'mm-YYYY') as start_date, TO_CHAR(end_date, 'mm-YYYY') as end_date "+
		"from subscriptions limit $1 offset $2",
		limit,
		offset)

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

func (repo *PostgresSubscriptionRepository) CalculateSubscriptionSum(
	ctx context.Context,
	userId string,
	subName string,
	dateFrom string,
	dateTo string) (*models.SubscriptionPriceSum, error) {

	rows, err := repo.Pool.Query(
		ctx,
		"select COALESCE(sum(price), 0) as price_sum from subscriptions "+
			"where ($1 = '' or user_id::TEXT = $1) and ($2 = '' or name = $2) and "+
			"($3 = '' or $4 = '' or start_date between TO_DATE($3, 'mm-YYYY') and TO_DATE($4, 'mm-YYYY')) ",
		userId,
		subName,
		dateFrom,
		dateTo)

	if err != nil {
		return nil, err
	}

	var priceSum models.SubscriptionPriceSum
	if rows.Next() {
		err = rows.Scan(&priceSum.Sum)
		if err != nil {
			return nil, err
		}

		return &priceSum, nil
	}
	priceSum.Sum = 0
	return &priceSum, nil
}
