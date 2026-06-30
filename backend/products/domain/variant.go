package domain

import (
	"errors"
	"time"

	"jolly/backend/common"
)

type VariantUUID struct {
	common.UUID
}

type Variant struct {
	id         VariantUUID
	productID  ProductUUID
	sku        string
	name       string
	priceCents int64
	images     []VariantImage
	createdAt  time.Time
	updatedAt  time.Time
}

func (v Variant) ID() VariantUUID {
	return v.id
}

func (v Variant) ProductID() ProductUUID {
	return v.productID
}

func (v Variant) SKU() string {
	return v.sku
}

func (v Variant) Name() string {
	return v.name
}

func (v Variant) PriceCents() int64 {
	return v.priceCents
}

func (v Variant) Images() []VariantImage {
	return v.images
}

func (v Variant) CreatedAt() time.Time {
	return v.createdAt
}

func (v Variant) UpdatedAt() time.Time {
	return v.updatedAt
}

func (v *Variant) AddImage(img VariantImage) error {
	for _, existing := range v.images {
		if existing.id == img.id {
			return errors.New("variant image already exists")
		}
	}
	v.images = append(v.images, img)
	v.updatedAt = common.NowUTC()
	return nil
}

func (v *Variant) RemoveImage(imageID VariantImageUUID) error {
	foundIdx := -1
	for idx, img := range v.images {
		if img.id == imageID {
			foundIdx = idx
			break
		}
	}
	if foundIdx == -1 {
		return errors.New("variant image not found")
	}

	v.images = append(v.images[:foundIdx], v.images[foundIdx+1:]...)
	v.updatedAt = common.NowUTC()
	return nil
}
