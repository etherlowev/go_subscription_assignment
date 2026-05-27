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
	"onlineSubscriptions/internal/services"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//go:embed db/*.sql
var migrationFiles embed.FS

var logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

func migrateDb(dbUrl string) {
	logger.Printf("Attempting to migrate database")
	d, err := iofs.New(migrationFiles, "db")
	if err != nil {
		logger.Print("Failed to create iofs")
		logger.Fatal(err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dbUrl+"?sslmode=disable")
	if err != nil {
		logger.Print("Failed to create migration instance")
		logger.Fatal(err)
	}

	migrationAppError := m.Up()

	if migrationAppError != nil && !errors.Is(migrationAppError, migrate.ErrNoChange) {
		logger.Print("Failed to migrate to database")
		logger.Fatal(err)
	} else if errors.Is(migrationAppError, migrate.ErrNoChange) {
		logger.Print("No changes applied to the database")
	}

	logger.Printf("Database migration is complete")
}

func Run() {
	logger.Print("Starting server")
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	migrateDb(cfg.Database.Url)

	connStr := cfg.Database.Url
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		logger.Fatal(err)
	}
	defer pool.Close()

	service := services.NewSubscriptionService(pool)

	router, err := handler.NewRouter(service)
	if err != nil {
		logger.Fatal(err)
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
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Println("Server shutdown gracefully")
}
