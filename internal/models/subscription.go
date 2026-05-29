package models

import (
	"errors"
	"github.com/google/uuid"
	"onlineSubscriptions/internal/validators"
	"strings"
)

type SubscriptionRequest struct {
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	UserId    uuid.UUID `json:"user_id"`
	StartDate string    `json:"start_date"`
	EndDate   string    `json:"end_date,omitempty"`
}

type Subscription struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	UserId    uuid.UUID `json:"user_id"`
	StartDate string    `json:"start_date"`
	EndDate   string    `json:"end_date,omitempty"`
}

type SubscriptionPriceSum struct {
	Sum int `json:"sum"`
}

func (r SubscriptionRequest) Validate() error {
	var valErrors = make([]string, 0)

	if r.Price < 0 {
		valErrors = append(valErrors, "price must be greater than or equal zero")
	}

	if strings.Trim(r.Name, " ") == "" {
		valErrors = append(valErrors, "subscription name can't be empty")
	}

	if r.UserId == uuid.Nil {
		valErrors = append(valErrors, "user_id can't be empty")
	}

	if r.StartDate == "" {
		valErrors = append(valErrors, "start_date can't be empty")
	} else if err := validators.ValidateDates(r.StartDate, r.EndDate); err != nil {
		valErrors = append(valErrors, err.Error())
	}

	if len(valErrors) > 0 {
		return errors.New(strings.Join(valErrors, "; "))
	}

	return nil
}
