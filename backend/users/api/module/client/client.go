package client

import (
	"context"
	"time"

	"jolly/backend/users/domain"
)

type Users interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (CreateUserResponse, error)
	GetUser(ctx context.Context, req GetUserRequest) (GetUserResponse, error)
}

type CreateUserRequest struct {
	Email string      `json:"email"`
	Name  string      `json:"name"`
	Role  domain.Role `json:"role"`
}

type CreateUserResponse struct {
	UserUUID domain.UserUUID `json:"user_uuid"`
	Role     domain.Role     `json:"role"`
}

type GetUserRequest struct {
	UserUUID domain.UserUUID `json:"user_uuid"`
}

type GetUserResponse struct {
	UserUUID  domain.UserUUID `json:"user_uuid"`
	Email     string          `json:"email"`
	Name      string          `json:"name"`
	Role      domain.Role     `json:"role"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
