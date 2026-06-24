package http

import (
	"context"

	"jolly/backend/common"
	moduleclient "jolly/backend/orders/api/module/client"
	"jolly/backend/orders/app/command"
	"jolly/backend/orders/app/query"
	"jolly/backend/orders/domain"
)

type Handler struct {
	commands *command.Handlers
	queries  *query.Handlers
}

func NewHandler(commands *command.Handlers, queries *query.Handlers) *Handler {
	if commands == nil {
		panic("orders command handlers cannot be nil")
	}
	if queries == nil {
		panic("orders query handlers cannot be nil")
	}

	return &Handler{
		commands: commands,
		queries:  queries,
	}
}

func (h Handler) CreateOrder(ctx context.Context, request CreateOrderRequestObject) (CreateOrderResponseObject, error) {
	if request.Body == nil {
		return nil, common.NewInvalidInputError("empty-body", "request body is required")
	}

	items := make([]moduleclient.OrderItem, 0, len(request.Body.Items))
	for _, item := range request.Body.Items {
		items = append(items, moduleclient.OrderItem{
			SKU:            item.Sku,
			Quantity:       item.Quantity,
			UnitPriceCents: item.UnitPriceCents,
		})
	}

	resp, err := h.commands.PlaceOrder(ctx, command.PlaceOrder{
		CustomerID: request.Body.CustomerId,
		Currency:   request.Body.Currency,
		Items:      items,
	})
	if err != nil {
		return nil, err
	}

	return CreateOrder201JSONResponse{
		OrderUuid: OrderUUID{UUID: resp.ID.UUID},
		Status:    OrderStatus(resp.Status),
	}, nil
}

func (h Handler) GetOrder(ctx context.Context, request GetOrderRequestObject) (GetOrderResponseObject, error) {
	order, err := h.queries.GetOrder(ctx, query.GetOrder{
		OrderID: domain.OrderID{UUID: request.OrderUuid.UUID},
	})
	if err != nil {
		return nil, err
	}

	items := make([]OrderLineItem, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, OrderLineItem{
			LineItemUuid:   item.UUID,
			Sku:            item.SKU,
			Quantity:       int32(item.Quantity),
			UnitPriceCents: item.UnitPriceCents,
		})
	}

	return GetOrder200JSONResponse{
		Currency:    order.Currency,
		CreatedAt:   order.CreatedAt,
		CustomerId:  order.CustomerID,
		Items:       items,
		OrderUuid:   OrderUUID{UUID: order.ID.UUID},
		PlacedAtUtc: order.PlacedAtUTC,
		Status:      OrderStatus(order.Status),
		TotalCents:  order.TotalCents,
		UpdatedAt:   order.UpdatedAt,
	}, nil
}

func Register(ctx context.Context, e common.EchoRouter, commands *command.Handlers, queries *query.Handlers) error {
	_ = ctx

	handler := Handler{
		commands: commands,
		queries:  queries,
	}

	RegisterHandlers(e, NewStrictHandler(handler, nil))
	return nil
}
