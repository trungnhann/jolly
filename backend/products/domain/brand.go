package domain

import (
	"errors"
	"strings"
	"time"

	"jolly/backend/common"
)

type BrandUUID struct {
	common.UUID
}

type Brand struct {
	id        BrandUUID
	name      string
	slug      string
	createdAt time.Time
	updatedAt time.Time
}

var (
	ErrBrandIDEmpty   = errors.New("brand id cannot be empty")
	ErrBrandNameEmpty = errors.New("brand name cannot be empty")
	ErrBrandSlugEmpty = errors.New("brand slug cannot be empty")
)

func (b Brand) ID() BrandUUID {
	return b.id
}

func (b Brand) Name() string {
	return b.name
}

func (b Brand) Slug() string {
	return b.slug
}

func (b Brand) CreatedAt() time.Time {
	return b.createdAt
}

func (b Brand) UpdatedAt() time.Time {
	return b.updatedAt
}

func NewBrand(id BrandUUID, name string, slug string) (Brand, error) {
	if id.IsZero() {
		return Brand{}, ErrBrandIDEmpty
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return Brand{}, ErrBrandNameEmpty
	}
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return Brand{}, ErrBrandSlugEmpty
	}
	return Brand{
		id:        id,
		name:      name,
		slug:      slug,
		createdAt: time.Now().UTC(),
		updatedAt: time.Now().UTC(),
	}, nil
}

func UnmarshalBrand(id BrandUUID, name string, slug string, createdAt time.Time, updatedAt time.Time) Brand {
	return Brand{
		id:        id,
		name:      name,
		slug:      slug,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}
