package command

import (
	"context"
	"errors"

	"jolly/backend/common"
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

type CreateUser struct {
	Email string
	Name  string
	Role  domain.Role
}

func (h *Handlers) CreateUser(ctx context.Context, cmd CreateUser) (domain.UserUUID, error) {
	userUUID := domain.UserUUID{UUID: common.NewUUIDv7()}

	user, err := domain.NewUser(
		userUUID,
		cmd.Email,
		cmd.Name,
		cmd.Role,
	)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidEmail):
			return domain.UserUUID{}, common.NewInvalidInputError("invalid_email", "invalid email")
		case errors.Is(err, domain.ErrNameEmpty):
			return domain.UserUUID{}, common.NewInvalidInputError("name_empty", "name cannot be empty")
		case errors.Is(err, domain.ErrRoleEmpty):
			return domain.UserUUID{}, common.NewInvalidInputError("role_empty", "role cannot be empty")
		case errors.Is(err, domain.ErrRoleInvalid):
			return domain.UserUUID{}, common.NewInvalidInputError("invalid_role", "invalid role")
		default:
			return domain.UserUUID{}, err
		}
	}

	if err := h.userRepository.CreateUser(ctx, user); err != nil {
		return domain.UserUUID{}, err
	}

	return userUUID, nil
}
