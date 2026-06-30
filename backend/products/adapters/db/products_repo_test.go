package db_test

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"jolly/backend/common"
	"jolly/backend/common/file"
	"jolly/backend/products"
	productsdb "jolly/backend/products/adapters/db"
	"jolly/backend/products/domain"
)

func TestPostgresRepository_Integration(t *testing.T) {
	// Skip integration test if not running locally with access to DB
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = os.Getenv("POSTGRES_URL")
	}
	if dbURL == "" {
		// Fallback for local machine dev
		dbURL = "postgres://jolly:jolly@localhost:5433/jolly?sslmode=disable"
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Skipf("Skipping integration test; failed to connect to Postgres: %v", err)
		return
	}
	defer pool.Close()

	// Ping the DB to make sure connection is alive
	if err := pool.Ping(ctx); err != nil {
		t.Skipf("Skipping integration test; Postgres ping failed: %v", err)
		return
	}

	fileStorage := file.NewPublicStorage("/tmp/jolly_test_storage", "http://localhost:8080/public")

	// Initialize the module to apply schema migrations
	mod := products.NewModule(pool, nil, fileStorage)
	if err := mod.Init(ctx); err != nil {
		t.Fatalf("failed to initialize products module and apply migrations: %v", err)
	}

	repo := productsdb.NewPostgresRepository(pool)

	t.Run("Create and Save Product with Variants", func(t *testing.T) {

		pID := domain.ProductUUID{UUID: common.NewUUIDv7()}
		p, err := domain.NewProduct(pID, "Integration Test Product "+pID.String(), "Integration Description", domain.ProductStatusDraft(), nil, nil)
		if err != nil {
			t.Fatalf("failed to construct product: %v", err)
		}

		vID1 := domain.VariantUUID{UUID: common.NewUUIDv7()}
		sku1 := "SKU-INT-1-" + pID.String()
		if err := p.AddVariant(vID1, sku1, "Size L", 2500); err != nil {
			t.Fatalf("failed to add variant 1: %v", err)
		}

		vID2 := domain.VariantUUID{UUID: common.NewUUIDv7()}
		sku2 := "SKU-INT-2-" + pID.String()
		if err := p.AddVariant(vID2, sku2, "Size XL", 3000); err != nil {
			t.Fatalf("failed to add variant 2: %v", err)
		}

		// Add variant images
		imgID1 := domain.VariantImageUUID{UUID: common.NewUUIDv7()}
		if err := p.AddVariantImage(vID1, imgID1, "http://test.com/img1.png", 1); err != nil {
			t.Fatalf("failed to add variant image 1: %v", err)
		}
		imgID2 := domain.VariantImageUUID{UUID: common.NewUUIDv7()}
		if err := p.AddVariantImage(vID2, imgID2, "http://test.com/img2.png", 2); err != nil {
			t.Fatalf("failed to add variant image 2: %v", err)
		}

		// Save the product
		if err := repo.SaveProduct(ctx, p); err != nil {
			t.Fatalf("failed to save product: %v", err)
		}

		// Retrieve the product
		retrieved, err := repo.ProductByID(ctx, pID)
		if err != nil {
			t.Fatalf("failed to load product: %v", err)
		}

		if retrieved.Name() != p.Name() {
			t.Errorf("expected product name %s, got %s", p.Name(), retrieved.Name())
		}
		if len(retrieved.Variants()) != 2 {
			t.Errorf("expected 2 variants, got %d", len(retrieved.Variants()))
		}

		// Verify images are reloaded
		v1 := retrieved.Variants()[0]
		if len(v1.Images()) != 1 || v1.Images()[0].URL() != "http://test.com/img1.png" {
			t.Errorf("failed to reload variant 1 images properly")
		}

		// Get Variant by SKU
		variant, err := repo.VariantBySKU(ctx, sku1)
		if err != nil {
			t.Fatalf("failed to get variant by sku: %v", err)
		}
		if variant.ID() != vID1 || variant.PriceCents() != 2500 || len(variant.Images()) != 1 {
			t.Errorf("retrieved variant properties or images incorrect")
		}

		// Update product and variants
		if err := p.UpdateDetails("Updated Integration Product", "New desc"); err != nil {
			t.Fatalf("failed to update details: %v", err)
		}
		// Remove variant 1 (should cascade delete variant 1 image), update variant 2
		if err := p.RemoveVariant(vID1); err != nil {
			t.Fatalf("failed to remove variant 1: %v", err)
		}
		if err := p.UpdateVariant(vID2, sku2, "Size XL - Updated", 3500); err != nil {
			t.Fatalf("failed to update variant 2: %v", err)
		}

		if err := repo.SaveProduct(ctx, p); err != nil {
			t.Fatalf("failed to save updated product: %v", err)
		}

		// Retrieve again
		retrieved2, err := repo.ProductByID(ctx, pID)
		if err != nil {
			t.Fatalf("failed to reload product: %v", err)
		}
		if len(retrieved2.Variants()) != 1 {
			t.Errorf("expected 1 variant after update, got %d", len(retrieved2.Variants()))
		}
		if retrieved2.Variants()[0].PriceCents() != 3500 {
			t.Errorf("expected updated price 3500, got %d", retrieved2.Variants()[0].PriceCents())
		}

		// Delete product
		if err := repo.DeleteProduct(ctx, pID); err != nil {
			t.Fatalf("failed to delete product: %v", err)
		}

		// Verify deletion
		_, err = repo.ProductByID(ctx, pID)
		if err == nil {
			t.Error("expected product not found error, got nil")
		}
	})

	t.Run("Create, Save, and Retrieve Category and Brand", func(t *testing.T) {
		catUUID := domain.CategoryUUID{UUID: common.NewUUIDv7()}
		cat, err := domain.NewCategory(catUUID, nil, "Shoes", "shoes-"+catUUID.String())
		if err != nil {
			t.Fatalf("failed to construct category: %v", err)
		}

		if err := repo.SaveCategory(ctx, cat); err != nil {
			t.Fatalf("failed to save category: %v", err)
		}

		// Retrieve Category
		retCat, err := repo.CategoryByID(ctx, catUUID)
		if err != nil {
			t.Fatalf("failed to load category: %v", err)
		}
		if retCat.Name() != "Shoes" || retCat.Slug() != "shoes-"+catUUID.String() {
			t.Errorf("category properties incorrect")
		}

		// Brand
		brandUUID := domain.BrandUUID{UUID: common.NewUUIDv7()}
		brand, err := domain.NewBrand(brandUUID, "Nike", "nike-"+brandUUID.String())
		if err != nil {
			t.Fatalf("failed to construct brand: %v", err)
		}

		if err := repo.SaveBrand(ctx, brand); err != nil {
			t.Fatalf("failed to save brand: %v", err)
		}

		// Retrieve Brand
		retBrand, err := repo.BrandByID(ctx, brandUUID)
		if err != nil {
			t.Fatalf("failed to load brand: %v", err)
		}
		if retBrand.Name() != "Nike" || retBrand.Slug() != "nike-"+brandUUID.String() {
			t.Errorf("brand properties incorrect")
		}
	})
}
