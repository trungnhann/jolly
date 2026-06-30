package domain

import (
	"time"
)

const (
	TopicInventoryReserved          = "inventory.reserved"
	TopicInventoryReservationFailed = "inventory.reservation.failed"
)

type InventoryReservedEvent struct {
	OrderID   string    `json:"order_id"`
	CreatedAt time.Time `json:"occurred_at"`
}

func (e InventoryReservedEvent) EventName() string {
	return TopicInventoryReserved
}

func (e InventoryReservedEvent) OccurredAt() time.Time {
	return e.CreatedAt
}

type InventoryReservationFailedEvent struct {
	OrderID   string    `json:"order_id"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"occurred_at"`
}

func (e InventoryReservationFailedEvent) EventName() string {
	return TopicInventoryReservationFailed
}

func (e InventoryReservationFailedEvent) OccurredAt() time.Time {
	return e.CreatedAt
}
