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

	return domain.UnmarshalUser(
		dbUser.UserUuid,
		dbUser.Email,
		dbUser.Name,
		dbUser.PasswordHash,
		dbUser.Role,
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

	return domain.UnmarshalUser(
		dbUser.UserUuid,
		dbUser.Email,
		dbUser.Name,
		dbUser.PasswordHash,
		dbUser.Role,
		dbUser.CreatedAt,
		dbUser.UpdatedAt,
	), nil
}
