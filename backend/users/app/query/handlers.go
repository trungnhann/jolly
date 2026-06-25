package query

import (
	"context"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

type GetUser struct {
	UserID domain.UserUUID
}

func (h *Handlers) GetUser(ctx context.Context, q GetUser) (domain.User, error) {
	return h.userRepository.UserByID(ctx, q.UserID)
}

type LoginUser struct {
	Email    string
	Password string
}

type LoginResult struct {
	Token    string
	UserUUID domain.UserUUID
}

func (h *Handlers) Login(ctx context.Context, q LoginUser) (LoginResult, error) {
	user, err := h.userRepository.UserByEmail(ctx, q.Email)
	if err != nil {
		return LoginResult{}, common.NewUnauthorizedError("invalid_credentials", "invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash()), []byte(q.Password)); err != nil {
		return LoginResult{}, common.NewUnauthorizedError("invalid_credentials", "invalid email or password")
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID().String(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	// Use JWT_SECRET from env, fallback to hardcoded default for local dev
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "jolly-secret-key-development"
	}

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{
		Token:    tokenString,
		UserUUID: user.ID(),
	}, nil
}
