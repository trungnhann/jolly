package http

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"jolly/backend/common"
	"jolly/backend/common/file"
	"jolly/backend/users/app/command"
	"jolly/backend/users/app/query"
	"jolly/backend/users/domain"
)

type Handler struct {
	commands *command.Handlers
	queries  *query.Handlers
	storage  file.Storage
}

func NewHandler(commands *command.Handlers, queries *query.Handlers, storage file.Storage) *Handler {
	if commands == nil {
		panic("users command handlers cannot be nil")
	}
	if queries == nil {
		panic("users query handlers cannot be nil")
	}
	if storage == nil {
		panic("storage cannot be nil")
	}

	return &Handler{
		commands: commands,
		queries:  queries,
		storage:  storage,
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
		Email:    request.Body.Email,
		Name:     request.Body.Name,
		Password: request.Body.Password,
		Role:     role,
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

	var avatarURL *string
	if user.AvatarURL() != "" {
		s := user.AvatarURL()
		avatarURL = &s
	}

	return GetUser200JSONResponse{
		UserUuid:  UserUUID{UUID: user.ID().UUID},
		Email:     user.Email(),
		Name:      user.Name(),
		Role:      user.Role(),
		AvatarUrl: avatarURL,
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}, nil
}

func (h Handler) LoginUser(ctx context.Context, request LoginUserRequestObject) (LoginUserResponseObject, error) {
	if request.Body == nil {
		return nil, common.NewInvalidInputError("empty-body", "request body is required")
	}

	result, err := h.queries.Login(ctx, query.LoginUser{
		Email:    request.Body.Email,
		Password: request.Body.Password,
	})
	if err != nil {
		return nil, err
	}

	return LoginUser200JSONResponse{
		Token:    result.Token,
		UserUuid: UserUUID{UUID: result.UserUUID.UUID},
	}, nil
}

func (h Handler) UploadAvatar(ctx context.Context, request UploadAvatarRequestObject) (UploadAvatarResponseObject, error) {
	if request.Body == nil {
		return nil, common.NewInvalidInputError("empty-body", "multipart body is required")
	}

	fileContent, filename, err := parseMultipartFile(request.Body)
	if err != nil {
		return nil, common.NewInvalidInputError("invalid-multipart", "failed to parse multipart file: %v", err)
	}

	userUUID := domain.UserUUID{UUID: request.UserUuid.UUID}
	ext := filepath.Ext(filename)
	storagePath := fmt.Sprintf("users/%s/avatar%s", userUUID.String(), ext)

	url, err := h.storage.StoreFile(ctx, storagePath, fileContent)
	if err != nil {
		return nil, err
	}

	err = h.commands.UpdateUserAvatar(ctx, command.UpdateUserAvatar{
		UserID:    userUUID,
		AvatarURL: url,
	})
	if err != nil {
		return nil, err
	}

	user, err := h.queries.GetUser(ctx, query.GetUser{UserID: userUUID})
	if err != nil {
		return nil, err
	}

	var avatarURL *string
	if user.AvatarURL() != "" {
		s := user.AvatarURL()
		avatarURL = &s
	}

	return UploadAvatar200JSONResponse{
		UserUuid:  UserUUID{UUID: user.ID().UUID},
		Email:     user.Email(),
		Name:      user.Name(),
		Role:      user.Role(),
		AvatarUrl: avatarURL,
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}, nil
}

func (h Handler) ForgotPassword(ctx context.Context, request ForgotPasswordRequestObject) (ForgotPasswordResponseObject, error) {
	if request.Body == nil {
		return nil, common.NewInvalidInputError("empty-body", "request body is required")
	}

	err := h.commands.RequestPasswordReset(ctx, command.RequestPasswordReset{
		Email: request.Body.Email,
	})
	if err != nil {
		return nil, err
	}

	return ForgotPassword200Response{}, nil
}

func (h Handler) ResetPassword(ctx context.Context, request ResetPasswordRequestObject) (ResetPasswordResponseObject, error) {
	if request.Body == nil {
		return nil, common.NewInvalidInputError("empty-body", "request body is required")
	}

	err := h.commands.ResetPassword(ctx, command.ResetPassword{
		Token:       request.Body.Token,
		NewPassword: request.Body.NewPassword,
	})
	if err != nil {
		return nil, err
	}

	return ResetPassword200Response{}, nil
}

func (h Handler) ChangePassword(ctx context.Context, request ChangePasswordRequestObject) (ChangePasswordResponseObject, error) {
	if request.Body == nil {
		return nil, common.NewInvalidInputError("empty-body", "request body is required")
	}

	authHeader := request.Params.Authorization
	if len(authHeader) <= 7 || !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, common.NewUnauthorizedError("invalid-auth", "invalid authorization header format")
	}
	tokenString := authHeader[7:]

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "jolly-secret-key-development"
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, common.NewUnauthorizedError("invalid-token", "invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, common.NewUnauthorizedError("invalid-claims", "invalid token claims")
	}

	userIDStr, ok := claims["sub"].(string)
	if !ok {
		return nil, common.NewUnauthorizedError("invalid-subject", "invalid token subject")
	}

	var rawUUID common.UUID
	err = rawUUID.UnmarshalText([]byte(userIDStr))
	if err != nil {
		return nil, common.NewUnauthorizedError("invalid-user-uuid", "failed to parse user uuid from token")
	}

	err = h.commands.ChangePassword(ctx, command.ChangePassword{
		UserID:          domain.UserUUID{UUID: rawUUID},
		CurrentPassword: request.Body.CurrentPassword,
		NewPassword:     request.Body.NewPassword,
	})
	if err != nil {
		return nil, err
	}

	return ChangePassword200Response{}, nil
}

func Register(ctx context.Context, e common.EchoRouter, commands *command.Handlers, queries *query.Handlers, storage file.Storage) error {
	_ = ctx

	handler := Handler{
		commands: commands,
		queries:  queries,
		storage:  storage,
	}

	RegisterHandlers(e, NewStrictHandler(handler, nil))
	return nil
}

func parseMultipartFile(reader *multipart.Reader) ([]byte, string, error) {
	for {
		part, err := reader.NextPart()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, "", err
		}

		if part.FormName() == "file" {
			defer part.Close()
			content, err := io.ReadAll(part)
			if err != nil {
				return nil, "", err
			}
			return content, part.FileName(), nil
		}
	}
	return nil, "", errors.New("no file part named 'file' found")
}
