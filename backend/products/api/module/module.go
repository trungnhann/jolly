package module

import (
	"context"

	"jolly/backend/products/api/module/client"
	"jolly/backend/products/app/query"
)

type Products struct {
	queryHandlers *query.Handlers
}

func New(queryHandlers *query.Handlers) *Products {
	if queryHandlers == nil {
		panic("products query handlers cannot be nil")
	}

	return &Products{
		queryHandlers: queryHandlers,
	}
}

func (p *Products) GetProduct(ctx context.Context, req client.GetProductRequest) (client.GetProductResponse, error) {
	product, err := p.queryHandlers.GetProduct(ctx, query.GetProduct{ID: req.ProductUUID})
	if err != nil {
		return client.GetProductResponse{}, err
	}

	variants := make([]client.VariantDTO, 0, len(product.Variants()))
	for _, v := range product.Variants() {
		variants = append(variants, client.VariantDTO{
			VariantUUID: v.ID(),
			SKU:         v.SKU(),
			Name:        v.Name(),
			PriceCents:  v.PriceCents(),
		})
	}

	return client.GetProductResponse{
		ProductUUID: product.ID(),
		Name:        product.Name(),
		Description: product.Description(),
		Status:      product.Status(),
		Variants:    variants,
		CreatedAt:   product.CreatedAt(),
		UpdatedAt:   product.UpdatedAt(),
	}, nil
}

func (p *Products) GetVariantBySKU(ctx context.Context, req client.GetVariantBySKURequest) (client.GetVariantBySKUResponse, error) {
	v, err := p.queryHandlers.GetVariantBySKU(ctx, query.GetVariantBySKU{SKU: req.SKU})
	if err != nil {
		return client.GetVariantBySKUResponse{}, err
	}

	return client.GetVariantBySKUResponse{
		VariantUUID: v.ID(),
		ProductUUID: v.ProductID(),
		SKU:         v.SKU(),
		Name:        v.Name(),
		PriceCents:  v.PriceCents(),
		CreatedAt:   v.CreatedAt(),
		UpdatedAt:   v.UpdatedAt(),
	}, nil
}
