package domain

import (
	"errors"
	"strings"
	"time"

	"jolly/backend/common"
)

type Role struct {
	common.Enum[RoleType]
}

type RoleType string

func (r RoleType) Values() []string {
	return []string{"customer", "admin"}
}

func RoleCustomer() Role {
	return common.MustEnum[Role]("customer")
}

func RoleAdmin() Role {
	return common.MustEnum[Role]("admin")
}

type UserUUID struct {
	common.UUID
}

type User struct {
	id           UserUUID
	email        string
	name         string
	passwordHash string
	role         Role
	createdAt    time.Time
	updatedAt    time.Time
}

var (
	ErrUserIDEmpty       = errors.New("user id cannot be empty")
	ErrInvalidEmail      = errors.New("invalid email")
	ErrNameEmpty         = errors.New("name cannot be empty")
	ErrPasswordHashEmpty = errors.New("password hash cannot be empty")
	ErrRoleEmpty         = errors.New("role cannot be empty")
	ErrRoleInvalid       = errors.New("invalid role")
)

func (u User) ID() UserUUID {
	return u.id
}

func (u User) Email() string {
	return u.email
}

func (u User) Name() string {
	return u.name
}

func (u User) PasswordHash() string {
	return u.passwordHash
}

func (u User) Role() Role {
	return u.role
}

func (u User) CreatedAt() time.Time {
	return u.createdAt
}

func (u User) UpdatedAt() time.Time {
	return u.updatedAt
}

func NewUser(id UserUUID, email string, name string, passwordHash string, role Role) (User, error) {
	if id.IsZero() {
		return User{}, ErrUserIDEmpty
	}

	email = strings.TrimSpace(email)
	if email == "" || !strings.Contains(email, "@") {
		return User{}, ErrInvalidEmail
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return User{}, ErrNameEmpty
	}

	if passwordHash == "" {
		return User{}, ErrPasswordHashEmpty
	}

	if role.IsZero() {
		return User{}, ErrRoleEmpty
	}
	if err := role.UnmarshalText([]byte(role.String())); err != nil {
		return User{}, ErrRoleInvalid
	}

	now := common.NowUTC()

	return User{
		id:           id,
		email:        email,
		name:         name,
		passwordHash: passwordHash,
		role:         role,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}
