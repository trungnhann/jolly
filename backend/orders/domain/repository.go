package domain

import (
	"context"
	"time"
)

type OrderRepository interface {
	SaveOrder(ctx context.Context, order Order) error
	UpdateOrderStatus(ctx context.Context, orderID OrderID, status Status, updatedAt time.Time) error
	OrderByID(ctx context.Context, orderID OrderID) (Order, error)
}
