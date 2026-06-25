package command

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

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
	Email    string
	Name     string
	Password string
	Role     domain.Role
}

func (h *Handlers) CreateUser(ctx context.Context, cmd CreateUser) (domain.UserUUID, error) {
	if cmd.Password == "" {
		return domain.UserUUID{}, common.NewInvalidInputError("password_empty", "password cannot be empty")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(cmd.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.UserUUID{}, err
	}

	userUUID := domain.UserUUID{UUID: common.NewUUIDv7()}

	user, err := domain.NewUser(
		userUUID,
		cmd.Email,
		cmd.Name,
		string(hashed),
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

type UpdateUserAvatar struct {
	UserID    domain.UserUUID
	AvatarURL string
}

func (h *Handlers) UpdateUserAvatar(ctx context.Context, cmd UpdateUserAvatar) error {
	if cmd.UserID.IsZero() {
		return domain.ErrUserIDEmpty
	}

	// Verify user exists
	_, err := h.userRepository.UserByID(ctx, cmd.UserID)
	if err != nil {
		return err
	}

	return h.userRepository.UpdateUserAvatar(ctx, cmd.UserID, cmd.AvatarURL)
}
