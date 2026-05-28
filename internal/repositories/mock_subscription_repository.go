package repositories

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"onlineSubscriptions/internal/models"
	"sync"
	"time"
)

type MockSubscriptionRepository struct {
	mu     sync.Mutex
	data   map[uuid.UUID]*models.Subscription
	nextId int
}

func NewMockSubscriptionRepository() *MockSubscriptionRepository {
	return &MockSubscriptionRepository{
		data: make(map[uuid.UUID]*models.Subscription),
	}
}

func (repo *MockSubscriptionRepository) GetById(_ context.Context, id uuid.UUID) (*models.Subscription, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	sub, ok := repo.data[id]
	if !ok {
		return nil, errors.New("subscription not found")
	}
	return sub, nil
}

func (repo *MockSubscriptionRepository) ExistsById(_ context.Context, id uuid.UUID) (bool, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	_, ok := repo.data[id]
	return ok, nil
}

func (repo *MockSubscriptionRepository) Insert(_ context.Context, uuid uuid.UUID, subscription *models.SubscriptionRequest) (bool, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, ok := repo.data[uuid]; ok {
		return false, nil
	}

	sub := &models.Subscription{
		Id:        uuid,
		Name:      subscription.Name,
		Price:     subscription.Price,
		UserId:    subscription.UserId,
		StartDate: subscription.StartDate,
		EndDate:   subscription.EndDate,
	}
	repo.data[uuid] = sub

	return true, nil
}

func (repo *MockSubscriptionRepository) Update(_ context.Context, subId uuid.UUID, subscription *models.SubscriptionRequest) (bool, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	sub, ok := repo.data[subId]
	if !ok {
		return false, errors.New("subscription not found")
	}

	sub.Id = subId
	sub.Name = subscription.Name
	sub.Price = subscription.Price
	sub.UserId = subscription.UserId
	sub.StartDate = subscription.StartDate
	sub.EndDate = subscription.EndDate

	repo.data[subId] = sub

	return true, nil
}

func (repo *MockSubscriptionRepository) DeleteById(_ context.Context, id uuid.UUID) (bool, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	_, ok := repo.data[id]
	if !ok {
		return false, errors.New("subscription not found")
	}

	delete(repo.data, id)
	return true, nil
}

func (repo *MockSubscriptionRepository) Page(_ context.Context, limit int, offset int) (*[]models.Subscription, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var subs []models.Subscription

	i := 0
	for _, sub := range repo.data {
		if i >= offset {
			subs = append(subs, *sub)
		}

		if i > limit-1 {
			break
		}

		i++
	}
	return &subs, nil
}

func (repo *MockSubscriptionRepository) CalculateSubscriptionSum(_ context.Context, userId string, subName string, dateFrom string, dateTo string) (*models.SubscriptionPriceSum, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	dateLayout := "01-2006"

	filterDateFrom, _ := time.Parse(dateLayout, dateFrom)
	filterDateTo, _ := time.Parse(dateLayout, dateTo)

	sum := 0

	for _, sub := range repo.data {
		if subName != "" && sub.Name != subName {
			continue
		}

		if userId != "" {
			uid, err := uuid.Parse(userId)
			if err != nil {
				return nil, err
			}

			if sub.UserId != uid {
				continue
			}
		}

		if dateFrom != "" && dateTo != "" {
			subStartDate, err := time.Parse(dateLayout, sub.StartDate)
			if err != nil {
				return nil, err
			}

			if filterDateFrom.After(subStartDate) || filterDateTo.Before(subStartDate) {
				continue
			}
		}

		sum += sub.Price
	}

	return &models.SubscriptionPriceSum{Sum: sum}, nil
}
