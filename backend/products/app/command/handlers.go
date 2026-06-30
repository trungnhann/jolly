package command

import (
	"context"
	"errors"

	"jolly/backend/common"
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

type CreateProductVariant struct {
	SKU        string
	Name       string
	PriceCents int64
}

type CreateProduct struct {
	Name        string
	Description string
	Status      domain.ProductStatus
	Variants    []CreateProductVariant
}

func (h *Handlers) CreateProduct(ctx context.Context, cmd CreateProduct) (domain.Product, error) {
	productUUID := domain.ProductUUID{UUID: common.NewUUIDv7()}

	product, err := domain.NewProduct(productUUID, cmd.Name, cmd.Description, cmd.Status)
	if err != nil {
		return domain.Product{}, mapDomainError(err)
	}

	for _, v := range cmd.Variants {
		vUUID := domain.VariantUUID{UUID: common.NewUUIDv7()}
		if err := product.AddVariant(vUUID, v.SKU, v.Name, v.PriceCents); err != nil {
			return domain.Product{}, mapDomainError(err)
		}
	}

	if err := h.repo.SaveProduct(ctx, product); err != nil {
		return domain.Product{}, err
	}

	return product, nil
}

type UpdateProduct struct {
	ID          domain.ProductUUID
	Name        string
	Description string
	Status      domain.ProductStatus
}

func (h *Handlers) UpdateProduct(ctx context.Context, cmd UpdateProduct) (domain.Product, error) {
	product, err := h.repo.ProductByID(ctx, cmd.ID)
	if err != nil {
		return domain.Product{}, err
	}

	if err := product.UpdateDetails(cmd.Name, cmd.Description); err != nil {
		return domain.Product{}, mapDomainError(err)
	}

	if err := product.ChangeStatus(cmd.Status); err != nil {
		return domain.Product{}, mapDomainError(err)
	}

	if err := h.repo.SaveProduct(ctx, product); err != nil {
		return domain.Product{}, err
	}

	return product, nil
}

type DeleteProduct struct {
	ID domain.ProductUUID
}

func (h *Handlers) DeleteProduct(ctx context.Context, cmd DeleteProduct) error {
	// Verify product exists first
	_, err := h.repo.ProductByID(ctx, cmd.ID)
	if err != nil {
		return err
	}

	return h.repo.DeleteProduct(ctx, cmd.ID)
}

type AddVariant struct {
	ProductID  domain.ProductUUID
	SKU        string
	Name       string
	PriceCents int64
}

func (h *Handlers) AddVariant(ctx context.Context, cmd AddVariant) (domain.Product, error) {
	product, err := h.repo.ProductByID(ctx, cmd.ProductID)
	if err != nil {
		return domain.Product{}, err
	}

	vUUID := domain.VariantUUID{UUID: common.NewUUIDv7()}
	if err := product.AddVariant(vUUID, cmd.SKU, cmd.Name, cmd.PriceCents); err != nil {
		return domain.Product{}, mapDomainError(err)
	}

	if err := h.repo.SaveProduct(ctx, product); err != nil {
		return domain.Product{}, err
	}

	return product, nil
}

type UpdateVariant struct {
	ProductID   domain.ProductUUID
	VariantUUID domain.VariantUUID
	SKU         string
	Name        string
	PriceCents  int64
}

func (h *Handlers) UpdateVariant(ctx context.Context, cmd UpdateVariant) (domain.Product, error) {
	product, err := h.repo.ProductByID(ctx, cmd.ProductID)
	if err != nil {
		return domain.Product{}, err
	}

	if err := product.UpdateVariant(cmd.VariantUUID, cmd.SKU, cmd.Name, cmd.PriceCents); err != nil {
		return domain.Product{}, mapDomainError(err)
	}

	if err := h.repo.SaveProduct(ctx, product); err != nil {
		return domain.Product{}, err
	}

	return product, nil
}

type RemoveVariant struct {
	ProductID   domain.ProductUUID
	VariantUUID domain.VariantUUID
}

func (h *Handlers) RemoveVariant(ctx context.Context, cmd RemoveVariant) (domain.Product, error) {
	product, err := h.repo.ProductByID(ctx, cmd.ProductID)
	if err != nil {
		return domain.Product{}, err
	}

	if err := product.RemoveVariant(cmd.VariantUUID); err != nil {
		return domain.Product{}, mapDomainError(err)
	}

	if err := h.repo.SaveProduct(ctx, product); err != nil {
		return domain.Product{}, err
	}

	return product, nil
}

type AddVariantImage struct {
	ProductID domain.ProductUUID
	VariantID domain.VariantUUID
	URL       string
	Position  int
}

func (h *Handlers) AddVariantImage(ctx context.Context, cmd AddVariantImage) (domain.Product, error) {
	product, err := h.repo.ProductByID(ctx, cmd.ProductID)
	if err != nil {
		return domain.Product{}, err
	}

	imageID := domain.VariantImageUUID{UUID: common.NewUUIDv7()}
	if err := product.AddVariantImage(cmd.VariantID, imageID, cmd.URL, cmd.Position); err != nil {
		return domain.Product{}, mapDomainError(err)
	}

	if err := h.repo.SaveProduct(ctx, product); err != nil {
		return domain.Product{}, err
	}

	return product, nil
}

type RemoveVariantImage struct {
	ProductID domain.ProductUUID
	VariantID domain.VariantUUID
	ImageID   domain.VariantImageUUID
}

func (h *Handlers) RemoveVariantImage(ctx context.Context, cmd RemoveVariantImage) (domain.Product, error) {
	product, err := h.repo.ProductByID(ctx, cmd.ProductID)
	if err != nil {
		return domain.Product{}, err
	}

	if err := product.RemoveVariantImage(cmd.VariantID, cmd.ImageID); err != nil {
		return domain.Product{}, mapDomainError(err)
	}

	if err := h.repo.SaveProduct(ctx, product); err != nil {
		return domain.Product{}, err
	}

	return product, nil
}

func mapDomainError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, domain.ErrProductIDEmpty):
		return common.NewInvalidInputError("product_id_empty", "product id cannot be empty")
	case errors.Is(err, domain.ErrProductNameEmpty):
		return common.NewInvalidInputError("product_name_empty", "product name cannot be empty")
	case errors.Is(err, domain.ErrProductStatusEmpty):
		return common.NewInvalidInputError("product_status_empty", "product status cannot be empty")
	case errors.Is(err, domain.ErrProductStatusInvalid):
		return common.NewInvalidInputError("invalid_product_status", "invalid product status")
	case errors.Is(err, domain.ErrVariantIDEmpty):
		return common.NewInvalidInputError("variant_id_empty", "variant id cannot be empty")
	case errors.Is(err, domain.ErrVariantSKUEmpty):
		return common.NewInvalidInputError("variant_sku_empty", "variant sku cannot be empty")
	case errors.Is(err, domain.ErrVariantPriceNegative):
		return common.NewInvalidInputError("variant_price_negative", "variant price cannot be negative")
	case errors.Is(err, domain.ErrVariantNotFound):
		return common.NewNotFoundError("variant_not_found", "variant not found")
	case errors.Is(err, domain.ErrDuplicateVariantSKU):
		return common.NewConflictError("duplicate_variant_sku", "variant sku already exists in this product")
	case errors.Is(err, domain.ErrImageIDEmpty):
		return common.NewInvalidInputError("image_id_empty", "image id cannot be empty")
	case errors.Is(err, domain.ErrImageURLEmpty):
		return common.NewInvalidInputError("image_url_empty", "image url cannot be empty")
	default:
		return err
	}
}
