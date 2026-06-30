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
		TotalCents: order.TotalCents,
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
		publishErr := fmt.Errorf("failed to publish order created event: %w", err)
		if markErr := h.markOrderFailed(ctx, &order); markErr != nil {
			return domain.Order{}, errors.Join(publishErr, markErr)
		}
		return domain.Order{}, publishErr
	}

	return order, nil
}

func (h *Handlers) markOrderFailed(ctx context.Context, order *domain.Order) error {
	if err := order.MarkFailed(); err != nil {
		return err
	}
	if err := h.orderRepository.UpdateOrderStatus(ctx, order.ID, order.Status, order.UpdatedAt); err != nil {
		return err
	}
	return nil
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
	return h.markOrderFailed(ctx, &order)
}

func (h *Handlers) MarkOrderConfirmed(ctx context.Context, orderID string) error {
	var id common.UUID
	if err := id.UnmarshalText([]byte(orderID)); err != nil {
		return err
	}
	order, err := h.orderRepository.OrderByID(ctx, domain.OrderID{UUID: id})
	if err != nil {
		return err
	}
	if err := order.MarkConfirmed(); err != nil {
		return err
	}
	return h.orderRepository.UpdateOrderStatus(ctx, order.ID, order.Status, order.UpdatedAt)
}

func (h *Handlers) MarkOrderPaid(ctx context.Context, orderID string) error {
	var id common.UUID
	if err := id.UnmarshalText([]byte(orderID)); err != nil {
		return err
	}
	order, err := h.orderRepository.OrderByID(ctx, domain.OrderID{UUID: id})
	if err != nil {
		return err
	}
	if err := order.MarkPaid(); err != nil {
		return err
	}
	if err := h.orderRepository.UpdateOrderStatus(ctx, order.ID, order.Status, order.UpdatedAt); err != nil {
		return err
	}

	eventPayload := domain.OrderPaidEvent{
		OrderID:   orderID,
		Items:     make([]domain.OrderPaidItem, 0, len(order.Items)),
		CreatedAt: common.NowUTC(),
	}
	for _, item := range order.Items {
		eventPayload.Items = append(eventPayload.Items, domain.OrderPaidItem{
			SKU:      item.SKU,
			Quantity: item.Quantity,
		})
	}

	return h.publisher.Publish(ctx, eventPayload.EventName(), eventPayload)
}
