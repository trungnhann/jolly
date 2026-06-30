package domain

import (
	"time"
)

func UnmarshalProduct(
	id ProductUUID,
	name string,
	description string,
	status ProductStatus,
	categoryUUID *CategoryUUID,
	brandUUID *BrandUUID,
	variants []Variant,
	createdAt time.Time,
	updatedAt time.Time,
) Product {
	return Product{
		id:           id,
		name:         name,
		description:  description,
		status:       status,
		categoryUUID: categoryUUID,
		brandUUID:    brandUUID,
		variants:     variants,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

func UnmarshalVariant(
	id VariantUUID,
	productID ProductUUID,
	sku string,
	name string,
	priceCents int64,
	images []VariantImage,
	createdAt time.Time,
	updatedAt time.Time,
) Variant {
	return Variant{
		id:         id,
		productID:  productID,
		sku:        sku,
		name:       name,
		priceCents: priceCents,
		images:     images,
		createdAt:  createdAt,
		updatedAt:  updatedAt,
	}
}

func UnmarshalVariantImage(
	id VariantImageUUID,
	variantID VariantUUID,
	url string,
	position int,
	createdAt time.Time,
) VariantImage {
	return VariantImage{
		id:        id,
		variantID: variantID,
		url:       url,
		position:  position,
		createdAt: createdAt,
	}
}
