package db_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"jolly/backend/common"
	"jolly/backend/inventory"
	inventorydb "jolly/backend/inventory/adapters/db"
	"jolly/backend/inventory/domain"
)

func TestPostgresRepository_Integration(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = os.Getenv("POSTGRES_URL")
	}
	if dbURL == "" {
		dbURL = "postgres://jolly:jolly@localhost:5433/jolly?sslmode=disable"
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Skipf("Skipping integration test; failed to connect to Postgres: %v", err)
		return
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		t.Skipf("Skipping integration test; Postgres ping failed: %v", err)
		return
	}

	// Initialize the module to apply schema migrations
	mod := inventory.NewModule(pool, nil, nil, nil)
	if err := mod.Init(ctx); err != nil {
		t.Fatalf("failed to initialize inventory module and apply migrations: %v", err)
	}

	repo := inventorydb.NewPostgresRepository(pool)

	t.Run("Reserve SKU that does not exist in stocks table", func(t *testing.T) {
		orderID := common.NewUUIDv7()
		reservationID := common.NewUUIDv7()

		items := []struct {
			SKU      string
			Quantity int
		}{
			{SKU: "SKU-DOES-NOT-EXIST", Quantity: 1},
		}

		err := repo.Reserve(ctx, reservationID, orderID, items)
		if !errors.Is(err, domain.ErrSKUNotFound) {
			t.Errorf("expected ErrSKUNotFound error, got: %v", err)
		}
	})

	t.Run("Reserve insufficient stock quantity", func(t *testing.T) {
		orderID := common.NewUUIDv7()
		reservationID := common.NewUUIDv7()

		items := []struct {
			SKU      string
			Quantity int
		}{
			{SKU: "SKU-123", Quantity: 500}, // Seeded has 100
		}

		err := repo.Reserve(ctx, reservationID, orderID, items)
		if !errors.Is(err, domain.ErrInsufficientStock) {
			t.Errorf("expected ErrInsufficientStock error, got: %v", err)
		}
	})

	t.Run("Successful Reservation and Idempotency Check", func(t *testing.T) {
		orderID := common.NewUUIDv7()
		reservationID := common.NewUUIDv7()

		items := []struct {
			SKU      string
			Quantity int
		}{
			{SKU: "SKU-123", Quantity: 5},
		}

		// Initial Reservation
		err := repo.Reserve(ctx, reservationID, orderID, items)
		if err != nil {
			t.Fatalf("failed to reserve stock: %v", err)
		}

		// Verify Reservation exists in DB
		exists, err := repo.GetReservationForOrder(ctx, orderID)
		if err != nil {
			t.Fatalf("failed to check reservation status: %v", err)
		}
		if !exists {
			t.Error("expected reservation to exist in database")
		}

		// Idempotency: Running Reserve again with the same OrderID should succeed without error
		err = repo.Reserve(ctx, reservationID, orderID, items)
		if err != nil {
			t.Errorf("idempotency check failed; expected success, got error: %v", err)
		}
	})

	t.Run("UpsertStock creates new stock and updates existing", func(t *testing.T) {
		sku := "SKU-NEW-" + common.NewUUIDv7().String()
		stockUUID := common.NewUUIDv7()

		// Create stock
		stock, err := repo.UpsertStock(ctx, stockUUID, sku, 50)
		if err != nil {
			t.Fatalf("failed to create stock: %v", err)
		}
		if stock.SKU != sku || stock.Quantity != 50 || stock.StockUUID != stockUUID {
			t.Errorf("incorrect created stock values: %+v", stock)
		}

		// Update stock
		updatedStock, err := repo.UpsertStock(ctx, stockUUID, sku, 75)
		if err != nil {
			t.Fatalf("failed to update stock: %v", err)
		}
		if updatedStock.SKU != sku || updatedStock.Quantity != 75 {
			t.Errorf("incorrect updated stock values: %+v", updatedStock)
		}
	})
}
