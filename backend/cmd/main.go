package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	backoff "github.com/cenkalti/backoff/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"jolly/backend"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	dsn := os.Getenv("POSTGRES_URL")
	if dsn == "" {
		dsn = os.Getenv("DATABASE_URL")
	}
	if dsn == "" {
		log.Fatal("POSTGRES_URL or DATABASE_URL environment variable is required")
	}

	dbPgx, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}

	b := backoff.NewExponentialBackOff()
	b.InitialInterval = 500 * time.Millisecond
	b.MaxInterval = 3 * time.Second

	_, err = backoff.Retry(
		ctx,
		func() (struct{}, error) {
			if pingErr := dbPgx.Ping(ctx); pingErr != nil {
				log.Printf("waiting for postgres: %v", pingErr)
				return struct{}{}, pingErr
			}
			return struct{}{}, nil
		},
		backoff.WithBackOff(b),
		backoff.WithMaxTries(20),
	)
	if err != nil {
		log.Fatal(err)
	}

	app, err := backend.New(ctx, dbPgx, backend.Config{
		HTTPAddr: ":8080",
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("jolly listening on :8080")
	if err := app.Start(ctx); err != nil {
		log.Fatal(err)
	}
}
