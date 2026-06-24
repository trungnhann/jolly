package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"jolly/backend/common"
	"jolly/backend/orders/adapters/db/dbmodels"
	"jolly/backend/orders/domain"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

const ordersCustomerIDForeignKeyConstraint = "orders_orders_customer_id_fkey"

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	if db == nil {
		panic("db connection pool cannot be nil")
	}

	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) SaveOrder(ctx context.Context, order domain.Order) error {
	return common.UpdateInTx(ctx, r.db, func(ctx context.Context, tx pgx.Tx) error {
		queries := dbmodels.New(tx)
		customerUUID := common.MustUUIDFromString(order.CustomerID)

		if err := queries.SaveOrder(ctx, dbmodels.SaveOrderParams{
			OrderUuid:   order.ID,
			CustomerID:  customerUUID,
			Currency:    order.Currency,
			Status:      order.Status,
			TotalCents:  order.TotalCents,
			PlacedAtUtc: order.PlacedAtUTC,
			CreatedAt:   order.CreatedAt,
			UpdatedAt:   order.UpdatedAt,
		}); err != nil {
			if common.IsForeignKeyViolationError(err, ordersCustomerIDForeignKeyConstraint) {
				return common.NewInvalidInputError("customer_not_found", "customer not found")
			}
			return fmt.Errorf("failed to save order %s: %w", order.ID, err)
		}

		for _, item := range order.Items {
			if err := queries.SaveOrderLineItem(ctx, dbmodels.SaveOrderLineItemParams{
				LineItemUuid:   item.UUID,
				OrderUuid:      order.ID,
				Sku:            item.SKU,
				Quantity:       int32(item.Quantity),
				UnitPriceCents: item.UnitPriceCents,
			}); err != nil {
				return fmt.Errorf("failed to save line item %s for order %s: %w", item.UUID, order.ID, err)
			}
		}

		return nil
	})
}

func (r *PostgresRepository) UpdateOrderStatus(ctx context.Context, orderID domain.OrderID, status domain.Status, updatedAt time.Time) error {
	queries := dbmodels.New(r.db)

	if err := queries.UpdateOrderStatus(ctx, dbmodels.UpdateOrderStatusParams{
		OrderUuid: orderID,
		Status:    status,
		UpdatedAt: updatedAt,
	}); err != nil {
		return fmt.Errorf("failed to update status for order %s: %w", orderID, err)
	}

	return nil
}

func (r *PostgresRepository) OrderByID(ctx context.Context, orderID domain.OrderID) (domain.Order, error) {
	queries := dbmodels.New(r.db)

	dbOrder, err := queries.GetOrder(ctx, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Order{}, common.NewNotFoundError("order_not_found", "order not found")
		}
		return domain.Order{}, fmt.Errorf("failed to get order %s: %w", orderID, err)
	}

	dbLineItems, err := queries.GetOrderLineItems(ctx, orderID)
	if err != nil {
		return domain.Order{}, fmt.Errorf("failed to get order line items %s: %w", orderID, err)
	}

	lineItems := make([]domain.LineItem, 0, len(dbLineItems))
	for _, row := range dbLineItems {
		lineItems = append(lineItems, domain.LineItem{
			UUID:           row.LineItemUuid,
			SKU:            row.Sku,
			Quantity:       int(row.Quantity),
			UnitPriceCents: row.UnitPriceCents,
		})
	}

	return domain.UnmarshalOrder(
		dbOrder.OrderUuid,
		dbOrder.CustomerID.String(),
		lineItems,
		dbOrder.Currency,
		dbOrder.Status,
		dbOrder.TotalCents,
		dbOrder.PlacedAtUtc,
		dbOrder.CreatedAt,
		dbOrder.UpdatedAt,
	), nil
}
