package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"rentos/internal/modules/auth/dto/request"
	"rentos/internal/modules/auth/dto/response"
	"rentos/internal/modules/auth/entity"
	"rentos/internal/modules/auth/repository"
	"rentos/pkg/jwt"
	"rentos/pkg/password"
	pkgresponse "rentos/pkg/response"
)

const (
	passwordResetTTL = 1 * time.Hour
	defaultPageSize  = 20
	maxPageSize      = 100
)

// ============================================================
// authService
// ============================================================

type authService struct {
	db       *sqlx.DB
	users    repository.UserRepository
	sessions repository.SessionRepository
	jwt      *jwt.Service
}

func NewAuthService(
	db *sqlx.DB,
	users repository.UserRepository,
	sessions repository.SessionRepository,
	jwtSvc *jwt.Service,
) AuthService {
	return &authService{db: db, users: users, sessions: sessions, jwt: jwtSvc}
}

func (s *authService) Register(ctx context.Context, tenantID string, req request.Register) (*response.Auth, error) {
	// Guard: email must be unique within the tenant.
	if _, err := s.users.FindByEmail(ctx, s.db, req.Email, tenantID); !errors.Is(err, repository.ErrNotFound) {
		if err == nil {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeConflict, "email already registered")
		}
		return nil, err
	}

	hash, err := password.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	u := &entity.User{
		ID:           uuid.NewString(),
		TenantID:     tenantID,
		Email:        req.Email,
		PasswordHash: hash,
		FirstName:    &req.FirstName,
		LastName:     &req.LastName,
		Phone:        req.Phone,
	}
	if err := s.users.Create(ctx, s.db, u); err != nil {
		return nil, err
	}

	return s.issueTokenPair(ctx, u, "", "")
}

func (s *authService) Login(ctx context.Context, tenantID string, req request.Login, ip, ua string) (*response.Auth, error) {
	u, err := s.users.FindByEmail(ctx, s.db, req.Email, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeUnauthorized, "invalid credentials")
		}
		return nil, err
	}

	if !u.IsActive {
		return nil, pkgresponse.NewAppError(pkgresponse.CodeForbidden, "account is inactive")
	}

	if err := password.Verify(req.Password, u.PasswordHash); err != nil {
		if errors.Is(err, password.ErrMismatch) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeUnauthorized, "invalid credentials")
		}
		return nil, err
	}

	// Fire-and-forget last login update — failure is non-fatal.
	_ = s.users.UpdateLastLogin(ctx, s.db, u.ID)

	return s.issueTokenPair(ctx, u, ip, ua)
}

func (s *authService) Refresh(ctx context.Context, req request.RefreshToken) (*response.Auth, error) {
	claims, err := s.jwt.ParseRefresh(req.RefreshToken)
	if err != nil {
		return nil, pkgresponse.NewAppError(pkgresponse.CodeUnauthorized, "invalid refresh token")
	}

	sess, err := s.sessions.FindByRefreshToken(ctx, s.db, req.RefreshToken)
	if err != nil {
		return nil, pkgresponse.NewAppError(pkgresponse.CodeUnauthorized, "session not found or expired")
	}

	// Rotate: revoke old session.
	if err := s.sessions.Revoke(ctx, s.db, sess.ID); err != nil {
		return nil, err
	}

	u, err := s.users.FindByID(ctx, s.db, claims.UserID, claims.TenantID)
	if err != nil {
		return nil, pkgresponse.NewAppError(pkgresponse.CodeUnauthorized, "user not found")
	}

	return s.issueTokenPair(ctx, u, "", "")
}

func (s *authService) Logout(ctx context.Context, sessionID string) error {
	return s.sessions.Revoke(ctx, s.db, sessionID)
}

// issueTokenPair mints access + refresh tokens and persists the session.
func (s *authService) issueTokenPair(ctx context.Context, u *entity.User, ip, ua string) (*response.Auth, error) {
	roles := []string{} // populated by RBAC module once it exists

	accessToken, err := s.jwt.IssueAccess(u.ID, u.TenantID, u.Email, roles)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.jwt.IssueRefresh(u.ID, u.TenantID, u.Email, roles)
	if err != nil {
		return nil, err
	}

	sess := &entity.UserSession{
		ID:           uuid.NewString(),
		TenantID:     u.TenantID,
		UserID:       u.ID,
		Token:        accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour),
	}
	if ip != "" {
		sess.IPAddress = &ip
	}
	if ua != "" {
		sess.UserAgent = &ua
	}

	if err := s.sessions.Create(ctx, s.db, sess); err != nil {
		return nil, err
	}

	return &response.Auth{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		User:         toUserResponse(u),
	}, nil
}

// ============================================================
// userService
// ============================================================

type userService struct {
	db    *sqlx.DB
	users repository.UserRepository
}

func NewUserService(db *sqlx.DB, users repository.UserRepository) UserService {
	return &userService{db: db, users: users}
}

func (s *userService) GetByID(ctx context.Context, id, tenantID string) (*entity.User, error) {
	u, err := s.users.FindByID(ctx, s.db, id, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "user not found")
		}
		return nil, err
	}
	return u, nil
}

func (s *userService) List(ctx context.Context, tenantID string, page, perPage int) ([]entity.User, error) {
	if perPage <= 0 {
		perPage = defaultPageSize
	}
	if perPage > maxPageSize {
		perPage = maxPageSize
	}
	if page <= 0 {
		page = 1
	}
	return s.users.List(ctx, s.db, tenantID, perPage, (page-1)*perPage)
}

func (s *userService) Update(ctx context.Context, id, tenantID string, req request.UpdateUser) (*entity.User, error) {
	u := &entity.User{ID: id, TenantID: tenantID, FirstName: req.FirstName, LastName: req.LastName, Phone: req.Phone}
	if err := s.users.Update(ctx, s.db, u); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "user not found")
		}
		return nil, err
	}
	return s.users.FindByID(ctx, s.db, id, tenantID)
}

func (s *userService) Delete(ctx context.Context, id, tenantID string) error {
	if err := s.users.Delete(ctx, s.db, id, tenantID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return pkgresponse.NewAppError(pkgresponse.CodeNotFound, "user not found")
		}
		return err
	}
	return nil
}

// ============================================================
// passwordService
// ============================================================

type passwordService struct {
	db     *sqlx.DB
	users  repository.UserRepository
	resets repository.PasswordResetRepository
}

func NewPasswordService(
	db *sqlx.DB,
	users repository.UserRepository,
	resets repository.PasswordResetRepository,
) PasswordService {
	return &passwordService{db: db, users: users, resets: resets}
}

func (s *passwordService) ForgotPassword(ctx context.Context, tenantID string, req request.ForgotPassword) error {
	u, err := s.users.FindByEmail(ctx, s.db, req.Email, tenantID)
	if err != nil {
		// Return no error even when user not found — prevents email enumeration.
		return nil
	}

	token, err := generateSecureToken()
	if err != nil {
		return err
	}

	pr := &entity.PasswordReset{
		ID:        uuid.NewString(),
		TenantID:  tenantID,
		UserID:    u.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(passwordResetTTL),
	}
	if err := s.resets.Create(ctx, s.db, pr); err != nil {
		return err
	}

	// TODO: enqueue email job (notification module, Phase 5)
	return nil
}

func (s *passwordService) ResetPassword(ctx context.Context, tenantID string, req request.ResetPassword) error {
	pr, err := s.resets.FindByToken(ctx, s.db, req.Token)
	if err != nil {
		return pkgresponse.NewAppError(pkgresponse.CodeNotFound, "invalid or expired reset token")
	}

	hash, err := password.Hash(req.Password)
	if err != nil {
		return err
	}

	if err := s.users.UpdatePassword(ctx, s.db, pr.UserID, hash); err != nil {
		return err
	}

	return s.resets.Consume(ctx, s.db, pr.ID)
}

// ============================================================
// helpers
// ============================================================

func generateSecureToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func toUserResponse(u *entity.User) response.User {
	return response.User{
		ID:          u.ID,
		TenantID:    u.TenantID,
		Email:       u.Email,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Phone:       u.Phone,
		AvatarURL:   u.AvatarURL,
		IsActive:    u.IsActive,
		LastLoginAt: u.LastLoginAt,
		Status:      u.Status,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}
