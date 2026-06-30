package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"jolly/backend/common"
	"jolly/backend/inventory/adapters/db/dbmodels"
	"jolly/backend/inventory/domain"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	if db == nil {
		panic("db pool cannot be nil")
	}
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetReservationForOrder(ctx context.Context, orderID common.UUID) (bool, error) {
	queries := dbmodels.New(r.db)
	_, err := queries.GetReservationForOrder(ctx, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *PostgresRepository) Reserve(ctx context.Context, reservationUUID, orderID common.UUID, items []struct {
	SKU      string
	Quantity int
}) error {
	return common.UpdateInTx(ctx, r.db, func(ctx context.Context, tx pgx.Tx) error {
		queries := dbmodels.New(tx)

		// 1. Check if reservation already exists (idempotency check inside transaction)
		_, err := queries.GetReservationForOrder(ctx, orderID)
		if err == nil {
			// Reservation already exists, idempotent success!
			return nil
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("failed to check existing reservation: %w", err)
		}

		// 2. Create the reservation record
		err = queries.CreateReservation(ctx, dbmodels.CreateReservationParams{
			ReservationUuid: reservationUUID,
			OrderUuid:       orderID,
			CreatedAt:       time.Now().UTC(),
		})
		if err != nil {
			return fmt.Errorf("failed to create reservation: %w", err)
		}

		// 3. For each item: lock stock row, check quantity, update stock, and insert reservation item
		for _, item := range items {
			// Lock stock row using FOR UPDATE
			stock, err := queries.GetStockForUpdate(ctx, item.SKU)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return fmt.Errorf("%w: %s", domain.ErrSKUNotFound, item.SKU)
				}
				return fmt.Errorf("failed to get stock for sku %s: %w", item.SKU, err)
			}

			// Validate inventory level
			if int(stock.Quantity) < item.Quantity {
				return fmt.Errorf("%w: for sku %s (requested %d, available %d)", domain.ErrInsufficientStock, item.SKU, item.Quantity, stock.Quantity)
			}

			// Calculate new quantity
			newQuantity := int(stock.Quantity) - item.Quantity

			// Update stock level
			err = queries.UpdateStockQuantity(ctx, dbmodels.UpdateStockQuantityParams{
				StockUuid: stock.StockUuid,
				Quantity:  int32(newQuantity),
			})
			if err != nil {
				return fmt.Errorf("failed to update stock for sku %s: %w", item.SKU, err)
			}

			// Insert reservation item
			err = queries.CreateReservationItem(ctx, dbmodels.CreateReservationItemParams{
				ReservationUuid: reservationUUID,
				StockUuid:       stock.StockUuid,
				Quantity:        int32(item.Quantity),
			})
			if err != nil {
				return fmt.Errorf("failed to create reservation item for sku %s: %w", item.SKU, err)
			}
		}

		return nil
	})
}

func (r *PostgresRepository) UpsertStock(ctx context.Context, stockUUID common.UUID, sku string, quantity int) (domain.Stock, error) {
	queries := dbmodels.New(r.db)
	row, err := queries.UpsertStock(ctx, dbmodels.UpsertStockParams{
		StockUuid: stockUUID,
		Sku:       sku,
		Quantity:  int32(quantity),
	})
	if err != nil {
		return domain.Stock{}, err
	}
	return domain.Stock{
		StockUUID: row.StockUuid,
		SKU:       row.Sku,
		Quantity:  int(row.Quantity),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}
