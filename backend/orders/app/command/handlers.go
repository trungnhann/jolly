package command

import (
	"context"
	"errors"
	"fmt"
	"time"

	"jolly/backend/common"
	"jolly/backend/common/module/contracts"
	"jolly/backend/common/shared"
	inventoryclient "jolly/backend/inventory/api/module/client"
	ordersclient "jolly/backend/orders/api/module/client"
	"jolly/backend/orders/app"
	"jolly/backend/orders/domain"
	paymentsclient "jolly/backend/payments/api/module/client"
)

type Handlers struct {
	modules         *contracts.Contracts
	orderRepository domain.OrderRepository
}

func NewHandlers(modules *contracts.Contracts, orderRepository domain.OrderRepository) *Handlers {
	if modules == nil {
		panic("contracts cannot be nil")
	}
	if orderRepository == nil {
		panic("order repository cannot be nil")
	}

	return &Handlers{
		modules:         modules,
		orderRepository: orderRepository,
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
	reserveItems := make([]inventoryclient.ReserveStockItem, 0, len(cmd.Items))
	for _, item := range cmd.Items {
		items = append(items, domain.LineItem{
			SKU:            item.SKU,
			Quantity:       item.Quantity,
			UnitPriceCents: item.UnitPriceCents,
		})
		reserveItems = append(reserveItems, inventoryclient.ReserveStockItem{
			SKU:      item.SKU,
			Quantity: item.Quantity,
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

	reserveErr := h.modules.Inventory.Reserve(ctx, inventoryclient.ReserveStockRequest{
		OrderID: orderUUID.String(),
		Items:   reserveItems,
	})
	if reserveErr != nil {
		return domain.Order{}, h.markOrderFailed(ctx, &order, reserveErr)
	}

	markReservedErr := order.MarkInventoryReserved()
	if markReservedErr != nil {
		return domain.Order{}, markReservedErr
	}
	updateReservedErr := h.orderRepository.UpdateOrderStatus(ctx, order.ID, order.Status, order.UpdatedAt)
	if updateReservedErr != nil {
		return domain.Order{}, updateReservedErr
	}

	_, err = h.modules.Payments.Authorize(ctx, paymentsclient.AuthorizePaymentRequest{
		OrderID:     orderUUID.String(),
		CustomerID:  order.CustomerID,
		AmountCents: order.TotalCents,
		Currency:    order.Currency.String(),
	})
	if err != nil {
		return domain.Order{}, h.markOrderFailed(ctx, &order, err)
	}

	markAuthorizedErr := order.MarkPaymentAuthorized()
	if markAuthorizedErr != nil {
		return domain.Order{}, markAuthorizedErr
	}
	updateAuthorizedErr := h.orderRepository.UpdateOrderStatus(ctx, order.ID, order.Status, order.UpdatedAt)
	if updateAuthorizedErr != nil {
		return domain.Order{}, updateAuthorizedErr
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
