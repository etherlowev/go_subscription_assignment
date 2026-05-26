package app

import (
	"context"
	"embed"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"onlineSubscriptions/internal/config"
	"onlineSubscriptions/internal/handler"
	services "onlineSubscriptions/internal/service"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//go:embed db/*.sql
var migrationFiles embed.FS

func migrateDb(db_url string) {
	d, err := iofs.New(migrationFiles, "db")
	if err != nil {
		log.Print("Failed to create iofs")
		log.Fatal(err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, db_url+"?sslmode=disable")
	if err != nil {
		log.Print("Failed to create migration instance")
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Print("Failed to migrate to database")
		log.Fatal(err)
	}
}

func Run() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	migrateDb(cfg.Database.Url)

	connStr := cfg.Database.Url
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	service := services.NewSubscriptionService(pool)

	router, err := handler.NewRouter(service)
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
	}

	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	shutdownSig := make(chan os.Signal, 1)

	signal.Notify(shutdownSig, os.Interrupt, syscall.SIGTERM)
	<-shutdownSig

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server shutdown gracefully")
}
