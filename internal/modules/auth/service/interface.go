package service

import (
	"context"

	"rentos-backend/internal/modules/auth/dto/request"
	"rentos-backend/internal/modules/auth/dto/response"
	"rentos-backend/internal/modules/auth/entity"
)

// AuthService handles registration, login, token refresh, and logout.
type AuthService interface {
	Register(ctx context.Context, tenantID string, req request.Register) (*response.Auth, error)
	Login(ctx context.Context, tenantID string, req request.Login, ip, ua string) (*response.Auth, error)
	Refresh(ctx context.Context, req request.RefreshToken) (*response.Auth, error)
	Logout(ctx context.Context, sessionID string) error
}

// UserService handles user profile management.
type UserService interface {
	GetByID(ctx context.Context, id, tenantID string) (*entity.User, error)
	List(ctx context.Context, tenantID string, page, perPage int) ([]entity.User, error)
	Update(ctx context.Context, id, tenantID string, req request.UpdateUser) (*entity.User, error)
	Delete(ctx context.Context, id, tenantID string) error
}

// PasswordService handles password-reset flows.
type PasswordService interface {
	ForgotPassword(ctx context.Context, tenantID string, req request.ForgotPassword) error
	ResetPassword(ctx context.Context, tenantID string, req request.ResetPassword) error
}
