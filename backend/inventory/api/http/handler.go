package http

import (
	"context"

	"jolly/backend/common"
	"jolly/backend/inventory/app/command"
	"jolly/backend/inventory/app/query"
)

type Handler struct {
	commands *command.Handlers
	queries  *query.Handlers
}

func NewHandler(commands *command.Handlers, queries *query.Handlers) *Handler {
	return &Handler{
		commands: commands,
		queries:  queries,
	}
}

func Register(ctx context.Context, e common.EchoRouter, commands *command.Handlers, queries *query.Handlers) error {
	handler := NewHandler(commands, queries)
	RegisterHandlers(e, NewStrictHandler(handler, nil))
	return nil
}

func (h *Handler) UpsertStock(ctx context.Context, req UpsertStockRequestObject) (UpsertStockResponseObject, error) {
	if req.Body == nil {
		return nil, common.NewInvalidInputError("empty_body", "body is required")
	}

	stock, err := h.commands.UpsertStock(ctx, command.UpsertStock{
		SKU:      req.Body.Sku,
		Quantity: req.Body.Quantity,
	})
	if err != nil {
		return nil, err
	}

	return UpsertStock200JSONResponse(StockResponse{
		StockUuid: stock.StockUUID,
		Sku:       stock.SKU,
		Quantity:  stock.Quantity,
	}), nil
}
