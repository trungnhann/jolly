package module

import (
	"context"

	"jolly/backend/orders/api/module/client"
	"jolly/backend/orders/app"
	"jolly/backend/orders/app/command"
	"jolly/backend/orders/app/query"
	"jolly/backend/orders/domain"
)

type Orders struct {
	commandHandlers *command.Handlers
	queryHandlers   *query.Handlers
}

func New(commandHandlers *command.Handlers, queryHandlers *query.Handlers) *Orders {
	if commandHandlers == nil {
		panic("orders command handlers cannot be nil")
	}
	if queryHandlers == nil {
		panic("orders query handlers cannot be nil")
	}

	return &Orders{
		commandHandlers: commandHandlers,
		queryHandlers:   queryHandlers,
	}
}

func (o *Orders) PlaceOrder(ctx context.Context, req client.PlaceOrderRequest) (client.PlaceOrderResponse, error) {
	order, err := o.commandHandlers.PlaceOrder(ctx, command.PlaceOrder{
		CustomerID: req.CustomerID,
		Currency:   req.Currency,
		Items:      req.Items,
	})
	if err != nil {
		return client.PlaceOrderResponse{}, err
	}

	return client.PlaceOrderResponse{
		OrderUUID: app.OrderUUID{UUID: order.ID.UUID},
		Status:    order.Status,
	}, nil
}

func (o *Orders) GetOrder(ctx context.Context, req client.GetOrderRequest) (client.GetOrderResponse, error) {
	order, err := o.queryHandlers.GetOrder(ctx, query.GetOrder{
		OrderID: domain.OrderID{UUID: req.OrderUUID.UUID},
	})
	if err != nil {
		return client.GetOrderResponse{}, err
	}

	items := make([]client.GetOrderLineItem, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, client.GetOrderLineItem{
			LineItemUUID:   item.UUID,
			SKU:            item.SKU,
			Quantity:       item.Quantity,
			UnitPriceCents: item.UnitPriceCents,
		})
	}

	return client.GetOrderResponse{
		OrderUUID:   req.OrderUUID,
		CustomerID:  order.CustomerID,
		Currency:    order.Currency,
		Status:      order.Status,
		TotalCents:  order.TotalCents,
		PlacedAtUTC: order.PlacedAtUTC,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
		Items:       items,
	}, nil
}
