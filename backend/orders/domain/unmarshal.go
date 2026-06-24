package domain

import (
	"time"

	"jolly/backend/common/shared"
)

func UnmarshalOrder(
	id OrderID,
	customerID string,
	items []LineItem,
	currency shared.Currency,
	status Status,
	totalCents int64,
	placedAtUTC time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) Order {
	return Order{
		ID:          id,
		CustomerID:  customerID,
		Items:       items,
		Currency:    currency,
		Status:      status,
		TotalCents:  totalCents,
		PlacedAtUTC: placedAtUTC,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}
