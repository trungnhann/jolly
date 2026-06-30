package command

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

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

	err = h.userRepository.UpdateUserAvatar(ctx, cmd.UserID, cmd.AvatarURL)
	if err != nil {
		return err
	}

	return nil
}

type RequestPasswordReset struct {
	Email string
}

func (h *Handlers) RequestPasswordReset(ctx context.Context, cmd RequestPasswordReset) error {
	user, err := h.userRepository.UserByEmail(ctx, cmd.Email)
	if err != nil {
		var commonErr common.Error
		if errors.As(err, &commonErr) && commonErr.HttpErrorCode == 404 {
			// Neutral discovery check: do not expose email existence
			fmt.Printf("[PasswordReset] Request received for non-existent email: %s\n", cmd.Email)
			return nil
		}
		return err
	}

	token, err := generateSecureToken()
	if err != nil {
		return fmt.Errorf("failed to generate secure token: %w", err)
	}

	hashedToken := hashToken(token)
	expiresAt := common.NowUTC().Add(15 * time.Minute)

	err = h.userRepository.CreateResetToken(ctx, hashedToken, user.ID(), expiresAt)
	if err != nil {
		return err
	}

	// In a real application, this sends an email.
	// For local testing and monolithic development, we output to console/logs so the user can copy the token.
	fmt.Printf("\n==================================================\n")
	fmt.Printf("[PASSWORD RESET LINK] Sent to: %s\n", user.Email())
	fmt.Printf("Reset Link: http://localhost:3000/reset-password?token=%s\n", token)
	fmt.Printf("==================================================\n\n")

	return nil
}

type ResetPassword struct {
	Token       string
	NewPassword string
}

func (h *Handlers) ResetPassword(ctx context.Context, cmd ResetPassword) error {
	if cmd.Token == "" {
		return common.NewInvalidInputError("token_empty", "reset token is required")
	}
	if cmd.NewPassword == "" {
		return common.NewInvalidInputError("password_empty", "new password cannot be empty")
	}
	if len(cmd.NewPassword) < 8 {
		return common.NewInvalidInputError("password_too_short", "password must be at least 8 characters long")
	}

	hashedToken := hashToken(cmd.Token)

	userUUID, expiresAt, err := h.userRepository.GetResetToken(ctx, hashedToken)
	if err != nil {
		return err
	}

	if common.NowUTC().After(expiresAt) {
		_ = h.userRepository.DeleteResetToken(ctx, hashedToken) // Invalidate expired token
		return common.NewInvalidInputError("token_expired", "reset token has expired")
	}

	user, err := h.userRepository.UserByID(ctx, userUUID)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cmd.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := user.UpdatePassword(string(hashedPassword)); err != nil {
		return err
	}

	err = h.userRepository.UpdateUserPassword(ctx, user.ID(), user.PasswordHash())
	if err != nil {
		return err
	}

	_ = h.userRepository.DeleteResetToken(ctx, hashedToken)

	return nil
}

type ChangePassword struct {
	UserID          domain.UserUUID
	CurrentPassword string
	NewPassword     string
}

func (h *Handlers) ChangePassword(ctx context.Context, cmd ChangePassword) error {
	if cmd.UserID.IsZero() {
		return domain.ErrUserIDEmpty
	}
	if cmd.CurrentPassword == "" {
		return common.NewInvalidInputError("current_password_empty", "current password is required")
	}
	if cmd.NewPassword == "" {
		return common.NewInvalidInputError("new_password_empty", "new password is required")
	}
	if len(cmd.NewPassword) < 8 {
		return common.NewInvalidInputError("new_password_too_short", "new password must be at least 8 characters long")
	}

	user, err := h.userRepository.UserByID(ctx, cmd.UserID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash()), []byte(cmd.CurrentPassword)); err != nil {
		return common.NewInvalidInputError("incorrect_password", "incorrect current password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cmd.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := user.UpdatePassword(string(hashedPassword)); err != nil {
		return err
	}

	return h.userRepository.UpdateUserPassword(ctx, user.ID(), user.PasswordHash())
}

func generateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
