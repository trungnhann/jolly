package domain

import (
	"errors"
	"strings"
	"time"

	"jolly/backend/common"
)

type VariantImageUUID struct {
	common.UUID
}

type VariantImage struct {
	id        VariantImageUUID
	variantID VariantUUID
	url       string
	position  int
	createdAt time.Time
}

var (
	ErrImageIDEmpty  = errors.New("image id cannot be empty")
	ErrImageURLEmpty = errors.New("image url cannot be empty")
)

func (v VariantImage) ID() VariantImageUUID {
	return v.id
}

func (v VariantImage) VariantID() VariantUUID {
	return v.variantID
}

func (v VariantImage) URL() string {
	return v.url
}

func (v VariantImage) Position() int {
	return v.position
}

func (v VariantImage) CreatedAt() time.Time {
	return v.createdAt
}

func NewVariantImage(id VariantImageUUID, variantID VariantUUID, url string, position int) (VariantImage, error) {
	if id.IsZero() {
		return VariantImage{}, ErrImageIDEmpty
	}
	if variantID.IsZero() {
		return VariantImage{}, ErrVariantIDEmpty
	}
	url = strings.TrimSpace(url)
	if url == "" {
		return VariantImage{}, ErrImageURLEmpty
	}

	return VariantImage{
		id:        id,
		variantID: variantID,
		url:       url,
		position:  position,
		createdAt: common.NowUTC(),
	}, nil
}
