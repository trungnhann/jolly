package db

import (
	"context"
	"errors"
	"fmt"

	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"jolly/backend/common"
	"jolly/backend/products/adapters/db/dbmodels"
	"jolly/backend/products/domain"
)

const variantsSKUUniqueConstraint = "variants_sku_unique"

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	if db == nil {
		panic("db connection pool cannot be nil")
	}

	return &PostgresRepository{db: db}
}

func toPgUUID(u *common.UUID) pgtype.UUID {
	if u == nil || u.IsZero() {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{
		Bytes: [16]byte(*u),
		Valid: true,
	}
}

func fromPgUUID(pgUUID pgtype.UUID) *common.UUID {
	if !pgUUID.Valid {
		return nil
	}
	u := common.UUID(pgUUID.Bytes)
	return &u
}

func (r *PostgresRepository) SaveProduct(ctx context.Context, product domain.Product) error {
	return common.UpdateInTx(ctx, r.db, func(ctx context.Context, tx pgx.Tx) error {
		queries := dbmodels.New(tx)

		_, err := queries.GetProduct(ctx, product.ID())
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				// Insert new product
				err = queries.CreateProduct(ctx, dbmodels.CreateProductParams{
					ProductUuid:  product.ID(),
					Name:         product.Name(),
					Description:  product.Description(),
					Status:       product.Status(),
					CategoryUuid: product.CategoryUUID(),
					BrandUuid:    product.BrandUUID(),
					CreatedAt:    product.CreatedAt(),
					UpdatedAt:    product.UpdatedAt(),
				})
				if err != nil {
					return fmt.Errorf("failed to create product %s: %w", product.ID(), err)
				}
			} else {
				return fmt.Errorf("failed to get product %s: %w", product.ID(), err)
			}
		} else {
			// Update existing product
			err = queries.UpdateProduct(ctx, dbmodels.UpdateProductParams{
				ProductUuid:  product.ID(),
				Name:         product.Name(),
				Description:  product.Description(),
				Status:       product.Status(),
				CategoryUuid: product.CategoryUUID(),
				BrandUuid:    product.BrandUUID(),
				UpdatedAt:    product.UpdatedAt(),
			})
			if err != nil {
				return fmt.Errorf("failed to update product %s: %w", product.ID(), err)
			}
		}

		// Delete existing variants for product
		err = queries.DeleteVariantsForProduct(ctx, product.ID())
		if err != nil {
			return fmt.Errorf("failed to delete variants for product %s: %w", product.ID(), err)
		}

		// Insert variants
		for _, v := range product.Variants() {
			err = queries.CreateVariant(ctx, dbmodels.CreateVariantParams{
				VariantUuid: v.ID(),
				ProductUuid: product.ID(),
				Sku:         v.SKU(),
				Name:        v.Name(),
				PriceCents:  v.PriceCents(),
				CreatedAt:   v.CreatedAt(),
				UpdatedAt:   v.UpdatedAt(),
			})
			if err != nil {
				if common.IsUniqueViolationError(err, variantsSKUUniqueConstraint) {
					return common.NewConflictError("sku_already_exists", "sku '%s' already exists", v.SKU())
				}
				return fmt.Errorf("failed to create variant %s for product %s: %w", v.ID(), product.ID(), err)
			}

			// Insert variant images
			for _, img := range v.Images() {
				err = queries.CreateVariantImage(ctx, dbmodels.CreateVariantImageParams{
					ImageUuid:   img.ID(),
					VariantUuid: v.ID(),
					Url:         img.URL(),
					Position:    int32(img.Position()),
					CreatedAt:   img.CreatedAt(),
				})
				if err != nil {
					return fmt.Errorf("failed to create variant image %s for variant %s: %w", img.ID(), v.ID(), err)
				}
			}
		}

		return nil
	})
}

func (r *PostgresRepository) ProductByID(ctx context.Context, id domain.ProductUUID) (domain.Product, error) {
	queries := dbmodels.New(r.db)

	dbProduct, err := queries.GetProduct(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Product{}, common.NewNotFoundError("product_not_found", "product not found")
		}
		return domain.Product{}, fmt.Errorf("failed to get product %s: %w", id, err)
	}

	dbVariants, err := queries.GetVariantsForProduct(ctx, id)
	if err != nil {
		return domain.Product{}, fmt.Errorf("failed to get variants for product %s: %w", id, err)
	}

	variants := make([]domain.Variant, 0, len(dbVariants))
	for _, row := range dbVariants {
		dbImages, err := queries.GetImagesForVariant(ctx, row.VariantUuid)
		if err != nil {
			return domain.Product{}, fmt.Errorf("failed to get images for variant %s: %w", row.VariantUuid, err)
		}

		images := make([]domain.VariantImage, 0, len(dbImages))
		for _, imgRow := range dbImages {
			images = append(images, domain.UnmarshalVariantImage(
				imgRow.ImageUuid,
				imgRow.VariantUuid,
				imgRow.Url,
				int(imgRow.Position),
				imgRow.CreatedAt,
			))
		}

		variants = append(variants, domain.UnmarshalVariant(
			row.VariantUuid,
			row.ProductUuid,
			row.Sku,
			row.Name,
			row.PriceCents,
			images,
			row.CreatedAt,
			row.UpdatedAt,
		))
	}

	return domain.UnmarshalProduct(
		dbProduct.ProductUuid,
		dbProduct.Name,
		dbProduct.Description,
		dbProduct.Status,
		dbProduct.CategoryUuid,
		dbProduct.BrandUuid,
		variants,
		dbProduct.CreatedAt,
		dbProduct.UpdatedAt,
	), nil
}

func (r *PostgresRepository) ListProducts(ctx context.Context) ([]domain.Product, error) {
	queries := dbmodels.New(r.db)

	dbProducts, err := queries.ListProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	products := make([]domain.Product, 0, len(dbProducts))
	for _, p := range dbProducts {
		dbVariants, err := queries.GetVariantsForProduct(ctx, p.ProductUuid)
		if err != nil {
			return nil, fmt.Errorf("failed to get variants for product %s: %w", p.ProductUuid, err)
		}

		variants := make([]domain.Variant, 0, len(dbVariants))
		for _, row := range dbVariants {
			dbImages, err := queries.GetImagesForVariant(ctx, row.VariantUuid)
			if err != nil {
				return nil, fmt.Errorf("failed to get images for variant %s: %w", row.VariantUuid, err)
			}

			images := make([]domain.VariantImage, 0, len(dbImages))
			for _, imgRow := range dbImages {
				images = append(images, domain.UnmarshalVariantImage(
					imgRow.ImageUuid,
					imgRow.VariantUuid,
					imgRow.Url,
					int(imgRow.Position),
					imgRow.CreatedAt,
				))
			}

			variants = append(variants, domain.UnmarshalVariant(
				row.VariantUuid,
				row.ProductUuid,
				row.Sku,
				row.Name,
				row.PriceCents,
				images,
				row.CreatedAt,
				row.UpdatedAt,
			))
		}

		products = append(products, domain.UnmarshalProduct(
			p.ProductUuid,
			p.Name,
			p.Description,
			p.Status,
			p.CategoryUuid,
			p.BrandUuid,
			variants,
			p.CreatedAt,
			p.UpdatedAt,
		))
	}

	return products, nil
}

func (r *PostgresRepository) VariantBySKU(ctx context.Context, sku string) (domain.Variant, error) {
	queries := dbmodels.New(r.db)

	dbVariant, err := queries.GetVariantBySKU(ctx, sku)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Variant{}, common.NewNotFoundError("variant_not_found", "variant with sku '%s' not found", sku)
		}
		return domain.Variant{}, fmt.Errorf("failed to get variant by sku %s: %w", sku, err)
	}

	dbImages, err := queries.GetImagesForVariant(ctx, dbVariant.VariantUuid)
	if err != nil {
		return domain.Variant{}, fmt.Errorf("failed to get images for variant %s: %w", dbVariant.VariantUuid, err)
	}

	images := make([]domain.VariantImage, 0, len(dbImages))
	for _, row := range dbImages {
		images = append(images, domain.UnmarshalVariantImage(
			row.ImageUuid,
			row.VariantUuid,
			row.Url,
			int(row.Position),
			row.CreatedAt,
		))
	}

	return domain.UnmarshalVariant(
		dbVariant.VariantUuid,
		dbVariant.ProductUuid,
		dbVariant.Sku,
		dbVariant.Name,
		dbVariant.PriceCents,
		images,
		dbVariant.CreatedAt,
		dbVariant.UpdatedAt,
	), nil
}

func (r *PostgresRepository) DeleteProduct(ctx context.Context, id domain.ProductUUID) error {
	queries := dbmodels.New(r.db)

	err := queries.DeleteProduct(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete product %s: %w", id, err)
	}

	return nil
}

// Categories

func (r *PostgresRepository) SaveCategory(ctx context.Context, category domain.Category) error {
	queries := dbmodels.New(r.db)
	err := queries.CreateCategory(ctx, dbmodels.CreateCategoryParams{
		CategoryUuid:       category.ID(),
		ParentCategoryUuid: category.ParentCategoryUUID(),
		Name:               category.Name(),
		Slug:               category.Slug(),
		CreatedAt:          category.CreatedAt(),
		UpdatedAt:          category.UpdatedAt(),
	})
	if err != nil {
		if common.IsUniqueViolationError(err, "categories_slug_unique") {
			return common.NewConflictError("category_slug_already_exists", "category slug '%s' already exists", category.Slug())
		}
		return fmt.Errorf("failed to save category: %w", err)
	}
	return nil
}

func (r *PostgresRepository) CategoryByID(ctx context.Context, id domain.CategoryUUID) (domain.Category, error) {
	queries := dbmodels.New(r.db)
	row, err := queries.GetCategory(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Category{}, common.NewNotFoundError("category_not_found", "category not found")
		}
		return domain.Category{}, fmt.Errorf("failed to get category: %w", err)
	}
	return domain.UnmarshalCategory(
		row.CategoryUuid,
		row.ParentCategoryUuid,
		row.Name,
		row.Slug,
		row.CreatedAt,
		row.UpdatedAt,
	), nil
}

func (r *PostgresRepository) ListCategories(ctx context.Context) ([]domain.Category, error) {
	queries := dbmodels.New(r.db)
	rows, err := queries.ListCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	categories := make([]domain.Category, 0, len(rows))
	for _, row := range rows {
		categories = append(categories, domain.UnmarshalCategory(
			row.CategoryUuid,
			row.ParentCategoryUuid,
			row.Name,
			row.Slug,
			row.CreatedAt,
			row.UpdatedAt,
		))
	}
	return categories, nil
}

// Brands

func (r *PostgresRepository) SaveBrand(ctx context.Context, brand domain.Brand) error {
	queries := dbmodels.New(r.db)
	err := queries.CreateBrand(ctx, dbmodels.CreateBrandParams{
		BrandUuid: brand.ID(),
		Name:      brand.Name(),
		Slug:      brand.Slug(),
		CreatedAt: brand.CreatedAt(),
		UpdatedAt: brand.UpdatedAt(),
	})
	if err != nil {
		if common.IsUniqueViolationError(err, "brands_slug_unique") {
			return common.NewConflictError("brand_slug_already_exists", "brand slug '%s' already exists", brand.Slug())
		}
		return fmt.Errorf("failed to save brand: %w", err)
	}
	return nil
}

func (r *PostgresRepository) BrandByID(ctx context.Context, id domain.BrandUUID) (domain.Brand, error) {
	queries := dbmodels.New(r.db)
	row, err := queries.GetBrand(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Brand{}, common.NewNotFoundError("brand_not_found", "brand not found")
		}
		return domain.Brand{}, fmt.Errorf("failed to get brand: %w", err)
	}
	return domain.UnmarshalBrand(
		row.BrandUuid,
		row.Name,
		row.Slug,
		row.CreatedAt,
		row.UpdatedAt,
	), nil
}

func (r *PostgresRepository) ListBrands(ctx context.Context) ([]domain.Brand, error) {
	queries := dbmodels.New(r.db)
	rows, err := queries.ListBrands(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list brands: %w", err)
	}
	brands := make([]domain.Brand, 0, len(rows))
	for _, row := range rows {
		brands = append(brands, domain.UnmarshalBrand(
			row.BrandUuid,
			row.Name,
			row.Slug,
			row.CreatedAt,
			row.UpdatedAt,
		))
	}
	return brands, nil
}
