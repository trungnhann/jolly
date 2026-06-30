package query

import (
	"context"

	"jolly/backend/products/domain"
)

type Handlers struct {
	repo domain.ProductRepository
}

func NewHandlers(repo domain.ProductRepository) *Handlers {
	if repo == nil {
		panic("product repository cannot be nil")
	}

	return &Handlers{repo: repo}
}

type GetProduct struct {
	ID domain.ProductUUID
}

func (h *Handlers) GetProduct(ctx context.Context, q GetProduct) (domain.Product, error) {
	return h.repo.ProductByID(ctx, q.ID)
}

type ListProducts struct{}

func (h *Handlers) ListProducts(ctx context.Context, q ListProducts) ([]domain.Product, error) {
	return h.repo.ListProducts(ctx)
}

type GetVariantBySKU struct {
	SKU string
}

func (h *Handlers) GetVariantBySKU(ctx context.Context, q GetVariantBySKU) (domain.Variant, error) {
	return h.repo.VariantBySKU(ctx, q.SKU)
}
