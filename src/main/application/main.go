package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"onlineSubscriptions/src/main/impl/api"
	services "onlineSubscriptions/src/main/impl/service"
)

// @title Subscription api
// @version 1.0
// @description Api for handling subscriptions
// @host localhost:8080
// @BasePath /api/subscriptions
func main() {
	connStr := "postgres://postgres:postgres@localhost:5432/subscriptions"
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	service := services.NewSubscriptionService(pool)

	router, err := api.NewRouter(service)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":8080", router))
}
