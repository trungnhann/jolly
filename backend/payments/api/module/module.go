package module

import (
	"context"

	"jolly/backend/payments/api/module/client"
	"jolly/backend/payments/app"
)

type Payments struct {
	service *app.Service
}

func New(service *app.Service) *Payments {
	if service == nil {
		panic("payments service cannot be nil")
	}

	return &Payments{service: service}
}

func (p *Payments) Authorize(ctx context.Context, req client.AuthorizePaymentRequest) (client.AuthorizePaymentResponse, error) {
	return p.service.Authorize(ctx, req)
}
