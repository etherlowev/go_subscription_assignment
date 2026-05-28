package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"onlineSubscriptions/internal/models"
	repos "onlineSubscriptions/internal/repositories"
	"testing"
)

var (
	user1UID = uuid.MustParse("e1036ff5-2876-45b9-bc45-37309f20629c")
	user2UID = uuid.MustParse("f8b63790-092a-461b-8e7a-d2d1e45aac5f")

	subUUID1 = uuid.MustParse("c19bf741-1191-46bb-9e4a-446779836371")
	subUUID2 = uuid.MustParse("fa1b0097-197c-4f12-89db-afd9a2f4bab1")
	subUUID3 = uuid.MustParse("5b1e7cfb-b4b7-48e6-8894-b43b28483518")
)

// Setup
func setup() (repos.SubscriptionRepository, SubscriptionService) {
	repo := repos.NewMockSubscriptionRepository()
	service := NewSubscriptionService(repo)

	return repo, service
}

func setupSubscriptions(ctx context.Context, repo repos.SubscriptionRepository) ([]uuid.UUID, error) {

	var subscriptionUids = []uuid.UUID{subUUID1, subUUID2, subUUID3}

	var subscriptionRequests = []models.SubscriptionRequest{
		{
			Name:      "netflix",
			Price:     100,
			UserId:    user1UID,
			StartDate: "05-2025",
		},
		{
			Name:      "paket",
			Price:     100,
			UserId:    user1UID,
			StartDate: "03-2025",
			EndDate:   "09-2025",
		},
		{
			Name:      "ivi",
			Price:     100,
			UserId:    user2UID,
			StartDate: "01-2025",
			EndDate:   "05-2025",
		},
	}

	for i := 0; i < len(subscriptionRequests); i++ {
		request := subscriptionRequests[i]
		subUUID := subscriptionUids[i]

		inserted, err := repo.Insert(ctx, subUUID, &request)
		if err != nil {
			return nil, err
		}
		if !inserted {
			return nil, errors.New(fmt.Sprintf("failed to insert subscription: %v", &request))
		}
	}
	return subscriptionUids, nil
}

func tearDownSubscriptions(ctx context.Context, repo repos.SubscriptionRepository) {
	subUUIDS := []uuid.UUID{subUUID1, subUUID2, subUUID3}

	for _, subUUID := range subUUIDS {
		_, _ = repo.DeleteById(ctx, subUUID)
	}
}

// End setup

// Tests

func TestGetById(t *testing.T) {
	repo, service := setup()
	ctx := context.Background()
	uids, err := setupSubscriptions(ctx, repo)
	defer tearDownSubscriptions(ctx, repo)

	if err != nil {
		t.Fatal(err)
	}

	for _, subId := range uids {
		_, err = service.GetSubscriptionById(ctx, subId)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestInsert(t *testing.T) {
	repo, service := setup()
	ctx := context.Background()
	defer tearDownSubscriptions(ctx, repo)

	inserted, err := service.InsertSubscription(ctx, &models.SubscriptionRequest{
		Name:      "netflix",
		Price:     100,
		UserId:    user1UID,
		StartDate: "05-2025",
		EndDate:   "08-2025",
	})

	if err != nil {
		t.Fatal(err)
	}

	if !inserted {
		t.Error("failed to insert subscription")
	}
}

func TestUpdate(t *testing.T) {
	repo, service := setup()
	ctx := context.Background()
	_, err := setupSubscriptions(ctx, repo)
	defer tearDownSubscriptions(ctx, repo)
	if err != nil {
		t.Fatal(err)
	}

	subscription, err := service.GetSubscriptionById(ctx, subUUID1)
	if err != nil {
		t.Fatal(err)
	}
	if subscription == nil {
		t.Error("TestUpdate() >> subscription is nil")
	}

	subUpdateRequest := models.SubscriptionRequest{
		Name:      "netflix2",
		Price:     200,
		UserId:    user2UID,
		StartDate: "05-2026",
	}

	updated, err := service.UpdateSubscription(ctx, subUUID1, &subUpdateRequest)
	if err != nil {
		t.Fatal(err)
	}
	if !updated {
		t.Error("TestUpdate() >> failed to update subscription")
	}
}

func TestSubscriptionSum(t *testing.T) {
	repo, service := setup()
	ctx := context.Background()
	_, err := setupSubscriptions(ctx, repo)
	defer tearDownSubscriptions(ctx, repo)

	if err != nil {
		t.Fatal(err)
	}

	noFilterSum, err := service.SubscriptionSum(ctx, "", "", "", "")
	if err != nil {
		t.Fatal(err)
	}

	if noFilterSum.Sum != 300 {
		t.Fatalf("TestSubscriptionSum() >> invalid sum, expected 300, got %v", noFilterSum.Sum)
	}

	user1Sum, err := service.SubscriptionSum(ctx, user1UID.String(), "", "", "")
	if err != nil {
		t.Fatal(err)
	}

	if user1Sum.Sum != 200 {
		t.Fatalf("TestSubscriptionSum() >> invalid sum, expected 200, got %v", user1Sum.Sum)
	}

	iviSum, err := service.SubscriptionSum(ctx, "", "ivi", "", "")

	if err != nil {
		t.Fatal(err)
	}

	if iviSum.Sum != 100 {
		t.Fatalf("TestSubscriptionSum() >> invalid sum, expected 100, got %v", iviSum.Sum)
	}

	dateSum, err := service.SubscriptionSum(ctx, "", "", "02-2025", "04-2025")
	if err != nil {
		t.Fatal(err)
	}

	if dateSum.Sum != 100 {
		t.Fatalf("TestSubscriptionSum() >> invalid sum, expected 100, got %v", dateSum.Sum)
	}
}

// End tests
