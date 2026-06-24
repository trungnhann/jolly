package common

import (
	"context"
	"errors"
	"fmt"
	"io/fs"

	migrate "github.com/golang-migrate/migrate/v4"
	pgxMigrate "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	"jolly/backend/common/log"
)

func MigrateDatabaseUp(
	ctx context.Context,
	moduleName string,
	pool *pgxpool.Pool,
	files fs.FS,
	migrationsDir string,
) error {
	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	d, err := iofs.New(files, migrationsDir)
	if err != nil {
		return fmt.Errorf("could not create iofs driver: %w", err)
	}

	if _, err := db.ExecContext(ctx, "CREATE SCHEMA IF NOT EXISTS "+moduleName); err != nil {
		return fmt.Errorf("could not create schema %s: %w", moduleName, err)
	}

	migDB, err := pgxMigrate.WithInstance(db, &pgxMigrate.Config{
		SchemaName:      moduleName,
		MigrationsTable: "schema_migrations",
	})
	if err != nil {
		return fmt.Errorf("could not connect to pgx migrations database: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", d, "pgx", migDB)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			log.FromContext(ctx).With("error", srcErr).Error("closing migration source failed")
		}
		if dbErr != nil {
			log.FromContext(ctx).With("error", dbErr).Error("closing migration database failed")
		}
	}()

	finished := make(chan struct{})
	defer close(finished)

	go func() {
		select {
		case <-finished:
			return
		case <-ctx.Done():
			log.FromContext(ctx).Info("Interrupt received, stopping migrations...")
			m.GracefulStop <- true
		}
	}()

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migration up failed: %w", err)
	}

	return nil
}
