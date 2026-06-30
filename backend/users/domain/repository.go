package domain

import (
	"context"
	"time"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user User) error
	UserByID(ctx context.Context, userID UserUUID) (User, error)
	UserByEmail(ctx context.Context, email string) (User, error)
	UpdateUserAvatar(ctx context.Context, userID UserUUID, avatarURL string) error
	UpdateUserPassword(ctx context.Context, userID UserUUID, passwordHash string) error
	CreateResetToken(ctx context.Context, tokenHash string, userID UserUUID, expiresAt time.Time) error
	GetResetToken(ctx context.Context, tokenHash string) (UserUUID, time.Time, error)
	DeleteResetToken(ctx context.Context, tokenHash string) error
}
