package model

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Subscription struct {
	Id        pgtype.UUID `json:"id"`
	Name      string      `json:"name"`
	Price     int         `json:"price"`
	UserId    pgtype.UUID `json:"user_id"`
	StartDate string      `json:"start_date"`
	EndDate   string      `json:"end_date"`
}
