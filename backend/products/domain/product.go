package domain

import (
	"errors"
	"strings"
	"time"

	"jolly/backend/common"
)

type ProductStatus struct {
	common.Enum[ProductStatusType]
}

type ProductStatusType string

func (p ProductStatusType) Values() []string {
	return []string{"draft", "active", "archived"}
}

func ProductStatusDraft() ProductStatus {
	return common.MustEnum[ProductStatus]("draft")
}

func ProductStatusActive() ProductStatus {
	return common.MustEnum[ProductStatus]("active")
}

func ProductStatusArchived() ProductStatus {
	return common.MustEnum[ProductStatus]("archived")
}

type ProductUUID struct {
	common.UUID
}

type Product struct {
	id           ProductUUID
	name         string
	description  string
	status       ProductStatus
	categoryUUID *CategoryUUID
	brandUUID    *BrandUUID
	variants     []Variant
	createdAt    time.Time
	updatedAt    time.Time
}

var (
	ErrProductIDEmpty       = errors.New("product id cannot be empty")
	ErrProductNameEmpty     = errors.New("product name cannot be empty")
	ErrProductStatusEmpty   = errors.New("product status cannot be empty")
	ErrProductStatusInvalid = errors.New("invalid product status")
	ErrVariantIDEmpty       = errors.New("variant id cannot be empty")
	ErrVariantSKUEmpty      = errors.New("variant sku cannot be empty")
	ErrVariantPriceNegative = errors.New("variant price cannot be negative")
	ErrVariantNotFound      = errors.New("variant not found")
	ErrDuplicateVariantSKU  = errors.New("variant sku already exists in this product")
)

func (p Product) ID() ProductUUID {
	return p.id
}

func (p Product) Name() string {
	return p.name
}

func (p Product) Description() string {
	return p.description
}

func (p Product) Status() ProductStatus {
	return p.status
}

func (p Product) Variants() []Variant {
	return p.variants
}

func (p Product) CategoryUUID() *CategoryUUID {
	return p.categoryUUID
}

func (p Product) BrandUUID() *BrandUUID {
	return p.brandUUID
}

func (p Product) CreatedAt() time.Time {
	return p.createdAt
}

func (p Product) UpdatedAt() time.Time {
	return p.updatedAt
}

func NewProduct(id ProductUUID, name string, description string, status ProductStatus, categoryUUID *CategoryUUID, brandUUID *BrandUUID) (Product, error) {
	if id.IsZero() {
		return Product{}, ErrProductIDEmpty
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return Product{}, ErrProductNameEmpty
	}
	if status.IsZero() {
		return Product{}, ErrProductStatusEmpty
	}
	if err := status.UnmarshalText([]byte(status.String())); err != nil {
		return Product{}, ErrProductStatusInvalid
	}

	now := common.NowUTC()
	return Product{
		id:           id,
		name:         name,
		description:  description,
		status:       status,
		categoryUUID: categoryUUID,
		brandUUID:    brandUUID,
		variants:     []Variant{},
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

func (p *Product) UpdateDetails(name string, description string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrProductNameEmpty
	}
	p.name = name
	p.description = description
	p.updatedAt = common.NowUTC()
	return nil
}

func (p *Product) UpdateCategoryAndBrand(categoryUUID *CategoryUUID, brandUUID *BrandUUID) {
	p.categoryUUID = categoryUUID
	p.brandUUID = brandUUID
	p.updatedAt = common.NowUTC()
}

func (p *Product) ChangeStatus(status ProductStatus) error {
	if status.IsZero() {
		return ErrProductStatusEmpty
	}
	if err := status.UnmarshalText([]byte(status.String())); err != nil {
		return ErrProductStatusInvalid
	}
	p.status = status
	p.updatedAt = common.NowUTC()
	return nil
}

func (p *Product) AddVariant(id VariantUUID, sku string, name string, priceCents int64) error {
	if id.IsZero() {
		return ErrVariantIDEmpty
	}
	sku = strings.TrimSpace(sku)
	if sku == "" {
		return ErrVariantSKUEmpty
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("variant name cannot be empty")
	}
	if priceCents < 0 {
		return ErrVariantPriceNegative
	}

	for _, v := range p.variants {
		if strings.EqualFold(v.sku, sku) {
			return ErrDuplicateVariantSKU
		}
	}

	now := common.NowUTC()
	variant := Variant{
		id:         id,
		productID:  p.id,
		sku:        sku,
		name:       name,
		priceCents: priceCents,
		createdAt:  now,
		updatedAt:  now,
	}

	p.variants = append(p.variants, variant)
	p.updatedAt = common.NowUTC()
	return nil
}

func (p *Product) UpdateVariant(id VariantUUID, sku string, name string, priceCents int64) error {
	if id.IsZero() {
		return ErrVariantIDEmpty
	}
	sku = strings.TrimSpace(sku)
	if sku == "" {
		return ErrVariantSKUEmpty
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("variant name cannot be empty")
	}
	if priceCents < 0 {
		return ErrVariantPriceNegative
	}

	foundIdx := -1
	for idx, v := range p.variants {
		if v.id == id {
			foundIdx = idx
		} else if strings.EqualFold(v.sku, sku) {
			return ErrDuplicateVariantSKU
		}
	}

	if foundIdx == -1 {
		return ErrVariantNotFound
	}

	p.variants[foundIdx].sku = sku
	p.variants[foundIdx].name = name
	p.variants[foundIdx].priceCents = priceCents
	p.variants[foundIdx].updatedAt = common.NowUTC()
	p.updatedAt = common.NowUTC()
	return nil
}

func (p *Product) RemoveVariant(id VariantUUID) error {
	if id.IsZero() {
		return ErrVariantIDEmpty
	}

	foundIdx := -1
	for idx, v := range p.variants {
		if v.id == id {
			foundIdx = idx
			break
		}
	}

	if foundIdx == -1 {
		return ErrVariantNotFound
	}

	p.variants = append(p.variants[:foundIdx], p.variants[foundIdx+1:]...)
	p.updatedAt = common.NowUTC()
	return nil
}

func (p *Product) AddVariantImage(variantID VariantUUID, imageID VariantImageUUID, url string, position int) error {
	foundIdx := -1
	for idx, v := range p.variants {
		if v.id == variantID {
			foundIdx = idx
			break
		}
	}
	if foundIdx == -1 {
		return ErrVariantNotFound
	}

	img, err := NewVariantImage(imageID, variantID, url, position)
	if err != nil {
		return err
	}

	if err := p.variants[foundIdx].AddImage(img); err != nil {
		return err
	}

	p.updatedAt = common.NowUTC()
	return nil
}

func (p *Product) RemoveVariantImage(variantID VariantUUID, imageID VariantImageUUID) error {
	foundIdx := -1
	for idx, v := range p.variants {
		if v.id == variantID {
			foundIdx = idx
			break
		}
	}
	if foundIdx == -1 {
		return ErrVariantNotFound
	}

	if err := p.variants[foundIdx].RemoveImage(imageID); err != nil {
		return err
	}

	p.updatedAt = common.NowUTC()
	return nil
}
