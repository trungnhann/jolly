package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"jolly/backend/payments/api/module/client"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Authorize(ctx context.Context, req client.AuthorizePaymentRequest) (client.AuthorizePaymentResponse, error) {
	if strings.TrimSpace(req.OrderID) == "" {
		return client.AuthorizePaymentResponse{}, fmt.Errorf("order id cannot be empty")
	}
	if req.AmountCents <= 0 {
		return client.AuthorizePaymentResponse{}, fmt.Errorf("amount must be positive")
	}
	if strings.TrimSpace(req.Currency) == "" {
		return client.AuthorizePaymentResponse{}, fmt.Errorf("currency cannot be empty")
	}

	return client.AuthorizePaymentResponse{
		PaymentID: fmt.Sprintf("pay_%d", time.Now().UnixNano()),
		Status:    "authorized",
	}, nil
}
