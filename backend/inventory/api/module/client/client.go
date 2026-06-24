package client

import "context"

type Inventory interface {
	Reserve(ctx context.Context, req ReserveStockRequest) error
}

type ReserveStockRequest struct {
	OrderID string
	Items   []ReserveStockItem
}

type ReserveStockItem struct {
	SKU      string
	Quantity int
}
