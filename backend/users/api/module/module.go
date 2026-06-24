package module

import (
	"context"

	"jolly/backend/users/api/module/client"
	"jolly/backend/users/app/command"
	"jolly/backend/users/app/query"
	"jolly/backend/users/domain"
)

type Users struct {
	commandHandlers *command.Handlers
	queryHandlers   *query.Handlers
}

func New(commandHandlers *command.Handlers, queryHandlers *query.Handlers) *Users {
	if commandHandlers == nil {
		panic("users command handlers cannot be nil")
	}
	if queryHandlers == nil {
		panic("users query handlers cannot be nil")
	}

	return &Users{
		commandHandlers: commandHandlers,
		queryHandlers:   queryHandlers,
	}
}

func (u *Users) CreateUser(ctx context.Context, req client.CreateUserRequest) (client.CreateUserResponse, error) {
	role := req.Role
	if role.IsZero() {
		role = domain.RoleCustomer()
	}

	userUUID, err := u.commandHandlers.CreateUser(ctx, command.CreateUser{
		Email: req.Email,
		Name:  req.Name,
		Role:  role,
	})
	if err != nil {
		return client.CreateUserResponse{}, err
	}

	return client.CreateUserResponse{
		UserUUID: userUUID,
		Role:     role,
	}, nil
}

func (u *Users) GetUser(ctx context.Context, req client.GetUserRequest) (client.GetUserResponse, error) {
	user, err := u.queryHandlers.GetUser(ctx, query.GetUser{
		UserID: req.UserUUID,
	})
	if err != nil {
		return client.GetUserResponse{}, err
	}

	return client.GetUserResponse{
		UserUUID:  user.ID(),
		Email:     user.Email(),
		Name:      user.Name(),
		Role:      user.Role(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}, nil
}
