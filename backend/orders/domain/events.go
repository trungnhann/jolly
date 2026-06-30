package domain

import (
	"time"

	"jolly/backend/common/shared"
)

const TopicOrderCreated = "order.created"

type OrderCreatedItem struct {
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
}

type OrderCreatedEvent struct {
	OrderID    string             `json:"order_id"`
	CustomerID string             `json:"customer_id"`
	Currency   shared.Currency    `json:"currency"`
	TotalCents int64              `json:"total_cents"`
	Items      []OrderCreatedItem `json:"items"`
	CreatedAt  time.Time          `json:"occurred_at"`
}

func (e OrderCreatedEvent) EventName() string {
	return TopicOrderCreated
}

func (e OrderCreatedEvent) OccurredAt() time.Time {
	return e.CreatedAt
}

const TopicOrderPaid = "order.paid"

type OrderPaidItem struct {
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
}

type OrderPaidEvent struct {
	OrderID   string          `json:"order_id"`
	Items     []OrderPaidItem `json:"items"`
	CreatedAt time.Time       `json:"occurred_at"`
}

func (e OrderPaidEvent) EventName() string {
	return TopicOrderPaid
}

func (e OrderPaidEvent) OccurredAt() time.Time {
	return e.CreatedAt
}
