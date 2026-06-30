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
	Items      []OrderCreatedItem `json:"items"`
	CreatedAt  time.Time          `json:"occurred_at"`
}

func (e OrderCreatedEvent) EventName() string {
	return TopicOrderCreated
}

func (e OrderCreatedEvent) OccurredAt() time.Time {
	return e.CreatedAt
}
