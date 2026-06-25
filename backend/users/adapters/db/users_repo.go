package db

import (
	"context"
	"errors"
	"fmt"

	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"jolly/backend/common"
	"jolly/backend/users/adapters/db/dbmodels"
	"jolly/backend/users/domain"
)

const usersEmailUniqueConstraint = "users_users_email_unique"

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	if db == nil {
		panic("db connection pool cannot be nil")
	}

	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateUser(ctx context.Context, user domain.User) error {
	queries := dbmodels.New(r.db)

	err := queries.CreateUser(ctx, dbmodels.CreateUserParams{
		UserUuid:     user.ID(),
		Email:        user.Email(),
		Name:         user.Name(),
		Role:         user.Role(),
		PasswordHash: user.PasswordHash(),
		AvatarUrl:    nil, // Empty by default upon registration
		CreatedAt:    user.CreatedAt(),
		UpdatedAt:    user.UpdatedAt(),
	})
	if err != nil {
		if common.IsUniqueViolationError(err, usersEmailUniqueConstraint) {
			return common.NewConflictError("email_already_exists", "email already exists")
		}
		return fmt.Errorf("failed to create user %s: %w", user.ID(), err)
	}

	return nil
}

func (r *PostgresRepository) UserByID(ctx context.Context, userID domain.UserUUID) (domain.User, error) {
	queries := dbmodels.New(r.db)

	dbUser, err := queries.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, common.NewNotFoundError("user_not_found", "user not found")
		}
		return domain.User{}, fmt.Errorf("failed to get user %s: %w", userID, err)
	}

	var avatarURL string
	if dbUser.AvatarUrl != nil {
		avatarURL = *dbUser.AvatarUrl
	}

	return domain.UnmarshalUser(
		dbUser.UserUuid,
		dbUser.Email,
		dbUser.Name,
		dbUser.PasswordHash,
		dbUser.Role,
		avatarURL,
		dbUser.CreatedAt,
		dbUser.UpdatedAt,
	), nil
}

func (r *PostgresRepository) UserByEmail(ctx context.Context, email string) (domain.User, error) {
	queries := dbmodels.New(r.db)

	dbUser, err := queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, common.NewNotFoundError("user_not_found", "user not found")
		}
		return domain.User{}, fmt.Errorf("failed to get user by email %s: %w", email, err)
	}

	var avatarURL string
	if dbUser.AvatarUrl != nil {
		avatarURL = *dbUser.AvatarUrl
	}

	return domain.UnmarshalUser(
		dbUser.UserUuid,
		dbUser.Email,
		dbUser.Name,
		dbUser.PasswordHash,
		dbUser.Role,
		avatarURL,
		dbUser.CreatedAt,
		dbUser.UpdatedAt,
	), nil
}

func (r *PostgresRepository) UpdateUserAvatar(ctx context.Context, userID domain.UserUUID, avatarURL string) error {
	queries := dbmodels.New(r.db)

	var avatarVal *string
	if avatarURL != "" {
		avatarVal = &avatarURL
	}

	err := queries.UpdateUserAvatar(ctx, dbmodels.UpdateUserAvatarParams{
		UserUuid:  userID,
		AvatarUrl: avatarVal,
		UpdatedAt: common.NowUTC(),
	})
	if err != nil {
		return fmt.Errorf("failed to update avatar for user %s: %w", userID, err)
	}

	return nil
}
