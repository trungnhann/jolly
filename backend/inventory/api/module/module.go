package module

import (
	"context"

	"jolly/backend/inventory/api/module/client"
	"jolly/backend/inventory/app/command"
)

type Inventory struct {
	commands *command.Handlers
}

func New(commands *command.Handlers) *Inventory {
	if commands == nil {
		panic("inventory command handlers cannot be nil")
	}

	return &Inventory{commands: commands}
}

func (i *Inventory) Reserve(ctx context.Context, req client.ReserveStockRequest) error {
	cmd := command.Reserve{
		OrderID: req.OrderID,
		Items: make([]struct {
			SKU      string
			Quantity int
		}, 0, len(req.Items)),
	}
	for _, item := range req.Items {
		cmd.Items = append(cmd.Items, struct {
			SKU      string
			Quantity int
		}{
			SKU:      item.SKU,
			Quantity: item.Quantity,
		})
	}
	return i.commands.Reserve(ctx, cmd)
}
