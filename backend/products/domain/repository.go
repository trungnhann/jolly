package domain

import "context"

type ProductRepository interface {
	SaveProduct(ctx context.Context, p Product) error
	ProductByID(ctx context.Context, id ProductUUID) (Product, error)
	ListProducts(ctx context.Context) ([]Product, error)
	VariantBySKU(ctx context.Context, sku string) (Variant, error)
	DeleteProduct(ctx context.Context, id ProductUUID) error

	// Categories
	SaveCategory(ctx context.Context, c Category) error
	CategoryByID(ctx context.Context, id CategoryUUID) (Category, error)
	ListCategories(ctx context.Context) ([]Category, error)

	// Brands
	SaveBrand(ctx context.Context, b Brand) error
	BrandByID(ctx context.Context, id BrandUUID) (Brand, error)
	ListBrands(ctx context.Context) ([]Brand, error)
}
