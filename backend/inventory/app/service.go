package app

import (
	"context"
	"fmt"
	"strings"

	"jolly/backend/inventory/api/module/client"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Reserve(ctx context.Context, req client.ReserveStockRequest) error {
	if strings.TrimSpace(req.OrderID) == "" {
		return fmt.Errorf("order id cannot be empty")
	}
	if len(req.Items) == 0 {
		return fmt.Errorf("at least one stock item is required")
	}

	for _, item := range req.Items {
		if strings.TrimSpace(item.SKU) == "" {
			return fmt.Errorf("stock item sku cannot be empty")
		}
		if item.Quantity <= 0 {
			return fmt.Errorf("stock item quantity must be positive")
		}
	}

	return nil
}
