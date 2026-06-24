package module

import (
	"context"

	"jolly/backend/inventory/api/module/client"
	"jolly/backend/inventory/app"
)

type Inventory struct {
	service *app.Service
}

func New(service *app.Service) *Inventory {
	if service == nil {
		panic("inventory service cannot be nil")
	}

	return &Inventory{service: service}
}

func (i *Inventory) Reserve(ctx context.Context, req client.ReserveStockRequest) error {
	return i.service.Reserve(ctx, req)
}
