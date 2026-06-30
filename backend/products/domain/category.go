package domain

import (
	"errors"
	"strings"
	"time"

	"jolly/backend/common"
)

type CategoryUUID struct {
	common.UUID
}

type Category struct {
	id                 CategoryUUID
	parentCategoryUUID *CategoryUUID
	name               string
	slug               string
	createdAt          time.Time
	updatedAt          time.Time
}

var (
	ErrCategoryIDEmpty   = errors.New("category id cannot be empty")
	ErrCategoryNameEmpty = errors.New("category name cannot be empty")
	ErrCategorySlugEmpty = errors.New("category slug cannot be empty")
)

func (c Category) ID() CategoryUUID {
	return c.id
}

func (c Category) ParentCategoryUUID() *CategoryUUID {
	return c.parentCategoryUUID
}

func (c Category) Name() string {
	return c.name
}

func (c Category) Slug() string {
	return c.slug
}

func (c Category) CreatedAt() time.Time {
	return c.createdAt
}

func (c Category) UpdatedAt() time.Time {
	return c.updatedAt
}

func NewCategory(id CategoryUUID, parentCategoryUUID *CategoryUUID, name string, slug string) (Category, error) {
	if id.IsZero() {
		return Category{}, ErrCategoryIDEmpty
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return Category{}, ErrCategoryNameEmpty
	}
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return Category{}, ErrCategorySlugEmpty
	}

	// Prevent a category from being its own parent
	if parentCategoryUUID != nil && parentCategoryUUID.String() == id.String() {
		return Category{}, errors.New("category cannot be its own parent")
	}

	return Category{
		id:                 id,
		parentCategoryUUID: parentCategoryUUID,
		name:               name,
		slug:               slug,
		createdAt:          time.Now().UTC(),
		updatedAt:          time.Now().UTC(),
	}, nil
}

func UnmarshalCategory(id CategoryUUID, parentCategoryUUID *CategoryUUID, name string, slug string, createdAt time.Time, updatedAt time.Time) Category {
	return Category{
		id:                 id,
		parentCategoryUUID: parentCategoryUUID,
		name:               name,
		slug:               slug,
		createdAt:          createdAt,
		updatedAt:          updatedAt,
	}
}
