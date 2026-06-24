package query

import (
	"context"

	"jolly/backend/orders/domain"
)

type Handlers struct {
	orderRepository domain.OrderRepository
}

func NewHandlers(orderRepository domain.OrderRepository) *Handlers {
	if orderRepository == nil {
		panic("order repository cannot be nil")
	}

	return &Handlers{orderRepository: orderRepository}
}

type GetOrder struct {
	OrderID domain.OrderID
}

func (h *Handlers) GetOrder(ctx context.Context, q GetOrder) (domain.Order, error) {
	return h.orderRepository.OrderByID(ctx, q.OrderID)
}
