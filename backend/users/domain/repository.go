package domain

import "context"

type UserRepository interface {
	CreateUser(ctx context.Context, user User) error
	UserByID(ctx context.Context, userID UserUUID) (User, error)
}
