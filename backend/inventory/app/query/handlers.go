package query

import (
	"jolly/backend/inventory/domain"
)

type Handlers struct {
	repo domain.Repository
}

func NewHandlers(repo domain.Repository) *Handlers {
	if repo == nil {
		panic("inventory repository cannot be nil")
	}
	return &Handlers{repo: repo}
}
