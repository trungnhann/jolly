package domain

import (
	"errors"
	"strings"
	"time"

	"jolly/backend/common"
	"jolly/backend/common/shared"
)

type Status string

const (
	StatusPending           Status = "pending"
	StatusInventoryReserved Status = "inventory_reserved"
	StatusPaymentAuthorized Status = "payment_authorized"
	StatusFailed            Status = "failed"
)

func (s Status) String() string {
	return string(s)
}

var (
	ErrOrderIDEmpty                = errors.New("order id cannot be empty")
	ErrCustomerIDEmpty             = errors.New("customer id cannot be empty")
	ErrCurrencyEmpty               = errors.New("currency cannot be empty")
	ErrOrderItemsEmpty             = errors.New("at least one line item is required")
	ErrLineItemSKUEmpty            = errors.New("line item sku cannot be empty")
	ErrLineItemQuantityInvalid     = errors.New("line item quantity must be positive")
	ErrLineItemPriceInvalid        = errors.New("line item price cannot be negative")
	ErrInvalidOrderStateTransition = errors.New("invalid order state transition")
)

type OrderID struct {
	common.UUID
}

type LineItemUUID struct {
	common.UUID
}

type LineItem struct {
	UUID           LineItemUUID
	SKU            string
	Quantity       int
	UnitPriceCents int64
}

type Order struct {
	ID          OrderID
	CustomerID  string
	Items       []LineItem
	Currency    shared.Currency
	Status      Status
	TotalCents  int64
	PlacedAtUTC time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewOrder(id OrderID, customerID string, items []LineItem, currency shared.Currency, placedAt time.Time) (Order, error) {
	if id.IsZero() {
		return Order{}, ErrOrderIDEmpty
	}
	if strings.TrimSpace(customerID) == "" {
		return Order{}, ErrCustomerIDEmpty
	}
	if currency.IsZero() {
		return Order{}, ErrCurrencyEmpty
	}
	if len(items) == 0 {
		return Order{}, ErrOrderItemsEmpty
	}

	var total int64
	for i, item := range items {
		if item.UUID.IsZero() {
			item.UUID = LineItemUUID{UUID: common.NewUUIDv7()}
			items[i] = item
		}
		if strings.TrimSpace(item.SKU) == "" {
			return Order{}, ErrLineItemSKUEmpty
		}
		if item.Quantity <= 0 {
			return Order{}, ErrLineItemQuantityInvalid
		}
		if item.UnitPriceCents < 0 {
			return Order{}, ErrLineItemPriceInvalid
		}
		total += int64(item.Quantity) * item.UnitPriceCents
	}

	now := common.NowUTC()

	return Order{
		ID:          id,
		CustomerID:  customerID,
		Items:       items,
		Currency:    currency,
		Status:      StatusPending,
		TotalCents:  total,
		PlacedAtUTC: placedAt.UTC(),
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (o *Order) MarkInventoryReserved() error {
	if o.Status != StatusPending {
		return ErrInvalidOrderStateTransition
	}
	o.Status = StatusInventoryReserved
	o.UpdatedAt = common.NowUTC()
	return nil
}

func (o *Order) MarkPaymentAuthorized() error {
	if o.Status != StatusInventoryReserved {
		return ErrInvalidOrderStateTransition
	}
	o.Status = StatusPaymentAuthorized
	o.UpdatedAt = common.NowUTC()
	return nil
}

func (o *Order) MarkFailed() error {
	switch o.Status {
	case StatusPending, StatusInventoryReserved:
		o.Status = StatusFailed
		o.UpdatedAt = common.NowUTC()
		return nil
	default:
		return ErrInvalidOrderStateTransition
	}
}
