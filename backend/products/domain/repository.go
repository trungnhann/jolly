package domain

import "context"

type ProductRepository interface {
	SaveProduct(ctx context.Context, p Product) error
	ProductByID(ctx context.Context, id ProductUUID) (Product, error)
	ListProducts(ctx context.Context) ([]Product, error)
	VariantBySKU(ctx context.Context, sku string) (Variant, error)
	DeleteProduct(ctx context.Context, id ProductUUID) error
}
