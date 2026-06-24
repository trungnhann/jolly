package client

import (
	"context"
	"time"

	"jolly/backend/common/shared"
	"jolly/backend/orders/app"
	"jolly/backend/orders/domain"
)

type Orders interface {
	PlaceOrder(ctx context.Context, req PlaceOrderRequest) (PlaceOrderResponse, error)
	GetOrder(ctx context.Context, req GetOrderRequest) (GetOrderResponse, error)
}

type PlaceOrderRequest struct {
	CustomerID string          `json:"customer_id"`
	Currency   shared.Currency `json:"currency"`
	Items      []OrderItem     `json:"items"`
}

type OrderItem struct {
	SKU            string `json:"sku"`
	Quantity       int    `json:"quantity"`
	UnitPriceCents int64  `json:"unit_price_cents"`
}

type PlaceOrderResponse struct {
	OrderUUID app.OrderUUID `json:"order_uuid"`
	Status    domain.Status `json:"status"`
}

type GetOrderRequest struct {
	OrderUUID app.OrderUUID `json:"order_uuid"`
}

type GetOrderResponse struct {
	OrderUUID   app.OrderUUID      `json:"order_uuid"`
	CustomerID  string             `json:"customer_id"`
	Currency    shared.Currency    `json:"currency"`
	Status      domain.Status      `json:"status"`
	TotalCents  int64              `json:"total_cents"`
	PlacedAtUTC time.Time          `json:"placed_at_utc"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	Items       []GetOrderLineItem `json:"items"`
}

type GetOrderLineItem struct {
	LineItemUUID   domain.LineItemUUID `json:"line_item_uuid"`
	SKU            string              `json:"sku"`
	Quantity       int                 `json:"quantity"`
	UnitPriceCents int64               `json:"unit_price_cents"`
}
