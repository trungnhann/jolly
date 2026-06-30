package domain

import (
	"context"
	"errors"
	"time"

	"jolly/backend/common"
)

var (
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrSKUNotFound       = errors.New("sku not found")
)

type Stock struct {
	StockUUID common.UUID
	SKU       string
	Quantity  int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ReservationItem struct {
	StockUUID common.UUID
	Quantity  int
}

type Reservation struct {
	ReservationUUID common.UUID
	OrderID         common.UUID
	CreatedAt       time.Time
	Items           []ReservationItem
}

type Repository interface {
	// Reserve atomic reservation: locks stock rows, checks availability, decrements stock, and saves reservation.
	Reserve(ctx context.Context, reservationUUID, orderID common.UUID, items []struct {
		SKU      string
		Quantity int
	}) error
	// GetReservationForOrder check if a reservation already exists for the given order ID
	GetReservationForOrder(ctx context.Context, orderID common.UUID) (bool, error)
	// UpsertStock creates or updates the stock level for a SKU
	UpsertStock(ctx context.Context, stockUUID common.UUID, sku string, quantity int) (Stock, error)
}
