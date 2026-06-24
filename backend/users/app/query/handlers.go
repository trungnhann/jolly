package query

import (
	"context"

	"jolly/backend/users/domain"
)

type Handlers struct {
	userRepository domain.UserRepository
}

func NewHandlers(userRepository domain.UserRepository) *Handlers {
	if userRepository == nil {
		panic("user repository cannot be nil")
	}

	return &Handlers{userRepository: userRepository}
}

type GetUser struct {
	UserID domain.UserUUID
}

func (h *Handlers) GetUser(ctx context.Context, q GetUser) (domain.User, error) {
	return h.userRepository.UserByID(ctx, q.UserID)
}
