package domain_test

import (
	"testing"

	"jolly/backend/common"
	"jolly/backend/products/domain"
)

func TestProduct_CreationAndMutations(t *testing.T) {
	// Success case
	id := domain.ProductUUID{UUID: common.NewUUIDv7()}
	p, err := domain.NewProduct(id, "Test Product", "Description", domain.ProductStatusDraft())
	if err != nil {
		t.Fatalf("unexpected error creating product: %v", err)
	}

	if p.ID() != id {
		t.Errorf("expected product ID %v, got %v", id, p.ID())
	}
	if p.Name() != "Test Product" {
		t.Errorf("expected product name 'Test Product', got '%s'", p.Name())
	}
	if p.Description() != "Description" {
		t.Errorf("expected product description 'Description', got '%s'", p.Description())
	}
	if p.Status() != domain.ProductStatusDraft() {
		t.Errorf("expected product status draft, got %s", p.Status())
	}

	// Update details
	err = p.UpdateDetails("Updated Name", "Updated Desc")
	if err != nil {
		t.Fatalf("unexpected error updating details: %v", err)
	}
	if p.Name() != "Updated Name" || p.Description() != "Updated Desc" {
		t.Errorf("details not updated properly")
	}

	// Empty name validation
	err = p.UpdateDetails("", "Desc")
	if err == nil {
		t.Error("expected error for empty name, got nil")
	}

	// Add Variant
	vID := domain.VariantUUID{UUID: common.NewUUIDv7()}
	err = p.AddVariant(vID, "SKU-123", "Red - M", 1500)
	if err != nil {
		t.Fatalf("unexpected error adding variant: %v", err)
	}

	if len(p.Variants()) != 1 {
		t.Errorf("expected 1 variant, got %d", len(p.Variants()))
	}
	v := p.Variants()[0]
	if v.ID() != vID || v.SKU() != "SKU-123" || v.Name() != "Red - M" || v.PriceCents() != 1500 {
		t.Errorf("variant properties incorrect: %+v", v)
	}

	// Add Duplicate Variant SKU
	vID2 := domain.VariantUUID{UUID: common.NewUUIDv7()}
	err = p.AddVariant(vID2, "SKU-123", "Blue - M", 1500)
	if err == nil {
		t.Error("expected duplicate SKU error, got nil")
	}

	// Update Variant
	err = p.UpdateVariant(vID, "SKU-123-UPDATED", "Red - M - Updated", 1800)
	if err != nil {
		t.Fatalf("unexpected error updating variant: %v", err)
	}
	v = p.Variants()[0]
	if v.SKU() != "SKU-123-UPDATED" || v.Name() != "Red - M - Updated" || v.PriceCents() != 1800 {
		t.Errorf("updated variant properties incorrect: %+v", v)
	}

	// Add Variant Image
	imgID := domain.VariantImageUUID{UUID: common.NewUUIDv7()}
	err = p.AddVariantImage(vID, imgID, "http://test.com/img.png", 0)
	if err != nil {
		t.Fatalf("unexpected error adding variant image: %v", err)
	}

	v = p.Variants()[0]
	if len(v.Images()) != 1 {
		t.Errorf("expected 1 variant image, got %d", len(v.Images()))
	}
	img := v.Images()[0]
	if img.ID() != imgID || img.URL() != "http://test.com/img.png" || img.Position() != 0 {
		t.Errorf("variant image properties incorrect: %+v", img)
	}

	// Remove Variant Image
	err = p.RemoveVariantImage(vID, imgID)
	if err != nil {
		t.Fatalf("unexpected error removing variant image: %v", err)
	}
	v = p.Variants()[0]
	if len(v.Images()) != 0 {
		t.Errorf("expected 0 variant images after removal, got %d", len(v.Images()))
	}

	// Remove Variant
	err = p.RemoveVariant(vID)
	if err != nil {
		t.Fatalf("unexpected error removing variant: %v", err)
	}
	if len(p.Variants()) != 0 {
		t.Errorf("expected 0 variants, got %d", len(p.Variants()))
	}
}
