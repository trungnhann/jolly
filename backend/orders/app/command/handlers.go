package command

import (
	"context"
	"errors"
	"fmt"
	"time"

	"jolly/backend/common"
	"jolly/backend/common/event"
	"jolly/backend/common/module/contracts"
	"jolly/backend/common/shared"
	ordersclient "jolly/backend/orders/api/module/client"
	"jolly/backend/orders/app"
	"jolly/backend/orders/domain"
)

type Handlers struct {
	modules         *contracts.Contracts
	orderRepository domain.OrderRepository
	publisher       event.Publisher
}

func NewHandlers(modules *contracts.Contracts, orderRepository domain.OrderRepository, publisher event.Publisher) *Handlers {
	if modules == nil {
		panic("contracts cannot be nil")
	}
	if orderRepository == nil {
		panic("order repository cannot be nil")
	}

	return &Handlers{
		modules:         modules,
		orderRepository: orderRepository,
		publisher:       publisher,
	}
}

type PlaceOrder struct {
	CustomerID string
	Currency   shared.Currency
	Items      []ordersclient.OrderItem
}

func (h *Handlers) PlaceOrder(ctx context.Context, cmd PlaceOrder) (domain.Order, error) {
	if h.modules.Inventory == nil {
		return domain.Order{}, fmt.Errorf("inventory contract is not registered")
	}
	if h.modules.Payments == nil {
		return domain.Order{}, fmt.Errorf("payments contract is not registered")
	}
	var customerUUID common.UUID
	if err := customerUUID.UnmarshalText([]byte(cmd.CustomerID)); err != nil {
		return domain.Order{}, common.NewInvalidInputError("invalid_customer_id", "customer_id must be a valid UUID")
	}

	orderUUID := app.OrderUUID{UUID: common.NewUUIDv7()}

	items := make([]domain.LineItem, 0, len(cmd.Items))
	for _, item := range cmd.Items {
		items = append(items, domain.LineItem{
			SKU:            item.SKU,
			Quantity:       item.Quantity,
			UnitPriceCents: item.UnitPriceCents,
		})
	}

	order, err := domain.NewOrder(domain.OrderID{UUID: orderUUID.UUID}, cmd.CustomerID, items, cmd.Currency, time.Now())
	if err != nil {
		return domain.Order{}, err
	}

	saveErr := h.orderRepository.SaveOrder(ctx, order)
	if saveErr != nil {
		return domain.Order{}, saveErr
	}

	eventPayload := domain.OrderCreatedEvent{
		OrderID:    orderUUID.String(),
		CustomerID: order.CustomerID,
		Currency:   order.Currency,
		Items:      make([]domain.OrderCreatedItem, 0, len(cmd.Items)),
		CreatedAt:  common.NowUTC(),
	}
	for _, item := range cmd.Items {
		eventPayload.Items = append(eventPayload.Items, domain.OrderCreatedItem{
			SKU:      item.SKU,
			Quantity: item.Quantity,
		})
	}

	if err := h.publisher.Publish(ctx, eventPayload.EventName(), eventPayload); err != nil {
		return domain.Order{}, h.markOrderFailed(ctx, &order, fmt.Errorf("failed to publish order created event: %w", err))
	}

	return order, nil
}

func (h *Handlers) markOrderFailed(ctx context.Context, order *domain.Order, cause error) error {
	if err := order.MarkFailed(); err != nil {
		return errors.Join(cause, err)
	}
	if err := h.orderRepository.UpdateOrderStatus(ctx, order.ID, order.Status, order.UpdatedAt); err != nil {
		return errors.Join(cause, err)
	}
	return cause
}

func (h *Handlers) MarkOrderFailed(ctx context.Context, orderID string) error {
	var id common.UUID
	if err := id.UnmarshalText([]byte(orderID)); err != nil {
		return err
	}
	order, err := h.orderRepository.OrderByID(ctx, domain.OrderID{UUID: id})
	if err != nil {
		return err
	}
	return h.markOrderFailed(ctx, &order, errors.New("inventory reservation failed"))
}

func (h *Handlers) MarkOrderInventoryReserved(ctx context.Context, orderID string) error {
	var id common.UUID
	if err := id.UnmarshalText([]byte(orderID)); err != nil {
		return err
	}
	order, err := h.orderRepository.OrderByID(ctx, domain.OrderID{UUID: id})
	if err != nil {
		return err
	}
	if err := order.MarkInventoryReserved(); err != nil {
		return err
	}
	return h.orderRepository.UpdateOrderStatus(ctx, order.ID, order.Status, order.UpdatedAt)
}

func (h *Handlers) MarkOrderPaymentAuthorized(ctx context.Context, orderID string) error {
	var id common.UUID
	if err := id.UnmarshalText([]byte(orderID)); err != nil {
		return err
	}
	order, err := h.orderRepository.OrderByID(ctx, domain.OrderID{UUID: id})
	if err != nil {
		return err
	}
	if err := order.MarkPaymentAuthorized(); err != nil {
		return err
	}
	return h.orderRepository.UpdateOrderStatus(ctx, order.ID, order.Status, order.UpdatedAt)
}
