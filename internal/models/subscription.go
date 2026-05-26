package models

import "github.com/google/uuid"

type SubscriptionRequest struct {
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	UserId    uuid.UUID `json:"user_id"`
	StartDate string    `json:"start_date"`
	EndDate   string    `json:"end_date"`
}

type Subscription struct {
	Id uuid.UUID `json:"id"`
	SubscriptionRequest
}
