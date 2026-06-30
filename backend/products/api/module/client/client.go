package client

import (
	"context"
	"time"

	"jolly/backend/products/domain"
)

type Products interface {
	GetProduct(ctx context.Context, req GetProductRequest) (GetProductResponse, error)
	GetVariantBySKU(ctx context.Context, req GetVariantBySKURequest) (GetVariantBySKUResponse, error)
}

type GetProductRequest struct {
	ProductUUID domain.ProductUUID
}

type VariantDTO struct {
	VariantUUID domain.VariantUUID
	SKU         string
	Name        string
	PriceCents  int64
}

type GetProductResponse struct {
	ProductUUID domain.ProductUUID
	Name        string
	Description string
	Status      domain.ProductStatus
	Variants    []VariantDTO
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type GetVariantBySKURequest struct {
	SKU string
}

type GetVariantBySKUResponse struct {
	VariantUUID domain.VariantUUID
	ProductUUID domain.ProductUUID
	SKU         string
	Name        string
	PriceCents  int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
