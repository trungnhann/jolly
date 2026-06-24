package http

import (
	"context"

	"jolly/backend/common"
	"jolly/backend/users/app/command"
	"jolly/backend/users/app/query"
	"jolly/backend/users/domain"
)

type Handler struct {
	commands *command.Handlers
	queries  *query.Handlers
}

func NewHandler(commands *command.Handlers, queries *query.Handlers) *Handler {
	if commands == nil {
		panic("users command handlers cannot be nil")
	}
	if queries == nil {
		panic("users query handlers cannot be nil")
	}

	return &Handler{
		commands: commands,
		queries:  queries,
	}
}

func (h Handler) CreateUser(ctx context.Context, request CreateUserRequestObject) (CreateUserResponseObject, error) {
	if request.Body == nil {
		return nil, common.NewInvalidInputError("empty-body", "request body is required")
	}

	role := domain.RoleCustomer()
	if request.Body.Role != nil && !request.Body.Role.IsZero() {
		role = *request.Body.Role
	}

	userUUID, err := h.commands.CreateUser(ctx, command.CreateUser{
		Email: request.Body.Email,
		Name:  request.Body.Name,
		Role:  role,
	})
	if err != nil {
		return nil, err
	}

	return CreateUser201JSONResponse{
		UserUuid: userUUID,
		Role:     role,
	}, nil
}

func (h Handler) GetUser(ctx context.Context, request GetUserRequestObject) (GetUserResponseObject, error) {
	user, err := h.queries.GetUser(ctx, query.GetUser{
		UserID: domain.UserUUID{UUID: request.UserUuid.UUID},
	})
	if err != nil {
		return nil, err
	}

	return GetUser200JSONResponse{
		UserUuid:  UserUUID{UUID: user.ID().UUID},
		Email:     user.Email(),
		Name:      user.Name(),
		Role:      user.Role(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}, nil
}

func Register(ctx context.Context, e common.EchoRouter, commands *command.Handlers, queries *query.Handlers) error {
	_ = ctx

	handler := Handler{
		commands: commands,
		queries:  queries,
	}

	RegisterHandlers(e, NewStrictHandler(handler, nil))
	return nil
}
