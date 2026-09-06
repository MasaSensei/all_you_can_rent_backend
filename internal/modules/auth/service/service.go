package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"rentos-backend/internal/modules/auth/dto/request"
	"rentos-backend/internal/modules/auth/dto/response"
	"rentos-backend/internal/modules/auth/entity"
	"rentos-backend/internal/modules/auth/repository"
	"rentos-backend/pkg/password"
	pkgresp "rentos-backend/pkg/response"
)

type AuthService interface {
	Login(ctx context.Context, req request.Login, userAgent, ip string) (*response.AuthTokens, error)
	Refresh(ctx context.Context, req request.Refresh) (*response.AuthTokens, error)
	Logout(ctx context.Context, req request.Logout) error
	Me(ctx context.Context, userID, tenantID string) (*response.UserInfo, error)
	ForgotPassword(ctx context.Context, tenantID string, req request.ForgotPassword) error
	ResetPassword(ctx context.Context, req request.ResetPassword) error
	ChangePassword(ctx context.Context, userID, tenantID string, req request.ChangePassword) error
}

type authService struct {
	db       *sqlx.DB
	users    repository.UserRepository
	sessions repository.SessionRepository
	resets   repository.PasswordResetRepository
	jwt      any
}

func NewAuthService(
	db *sqlx.DB,
	users repository.UserRepository,
	sessions repository.SessionRepository,
	resets repository.PasswordResetRepository,
	jwt any,
) AuthService {
	return &authService{db: db, users: users, sessions: sessions, resets: resets, jwt: jwt}
}

func (s *authService) Login(
	ctx context.Context,
	req request.Login,
	userAgent,
	ip string,
) (*response.AuthTokens, error) {

	u, err := s.users.FindByEmail(ctx, s.db, req.Email)
	if err != nil {
		fmt.Println("FIND USER ERROR:", err)
		return nil, pkgresp.NewAppError(
			pkgresp.CodeUnauthorized,
			"email atau password salah",
		)
	}

	fmt.Println("USER FOUND:", u.ID)

	if !u.IsActive {
		fmt.Println("USER NOT ACTIVE")
		return nil, pkgresp.NewAppError(
			pkgresp.CodeUnauthorized,
			"akun dinonaktifkan",
		)
	}

	if u.PasswordHash == nil {
		fmt.Println("PASSWORD HASH NIL")
		return nil, pkgresp.NewAppError(
			pkgresp.CodeUnauthorized,
			"email atau password salah",
		)
	}

	fmt.Println("VERIFY PASSWORD")

	if err := password.Verify(*u.PasswordHash, req.Password); err != nil {
		fmt.Println("PASSWORD VERIFY FAILED:", err)
		return nil, pkgresp.NewAppError(
			pkgresp.CodeUnauthorized,
			"email atau password salah",
		)
	}

	fmt.Println("PASSWORD OK")

	return s.issueTokens(ctx, u, userAgent, ip)
}

func (s *authService) Refresh(ctx context.Context, req request.Refresh) (*response.AuthTokens, error) {
	claims, err := parseClaims(s.jwt, req.RefreshToken)
	if err != nil {
		return nil, pkgresp.NewAppError(pkgresp.CodeUnauthorized, "refresh token tidak valid atau kadaluarsa")
	}
	if claims.TokenType != "refresh" {
		return nil, pkgresp.NewAppError(pkgresp.CodeUnauthorized, "bukan refresh token")
	}

	// Revoke old refresh token
	_ = s.sessions.Revoke(ctx, s.db, req.RefreshToken)

	u, err := s.users.FindByID(ctx, s.db, claims.UserID)
	if err != nil || !u.IsActive {
		return nil, pkgresp.NewAppError(pkgresp.CodeUnauthorized, "user tidak ditemukan atau nonaktif")
	}

	return s.issueTokens(ctx, u, "", "")
}

func (s *authService) Logout(ctx context.Context, req request.Logout) error {
	return s.sessions.Revoke(ctx, s.db, req.RefreshToken)
}

func (s *authService) Me(ctx context.Context, userID, tenantID string) (*response.UserInfo, error) {
	u, err := s.users.FindByID(ctx, s.db, userID)
	if err != nil {
		return nil, pkgresp.NewAppError(pkgresp.CodeNotFound, "user tidak ditemukan")
	}
	return toUserInfo(u), nil
}

func (s *authService) ForgotPassword(ctx context.Context, tenantID string, req request.ForgotPassword) error {
	u, err := s.users.FindByEmail(ctx, s.db, req.Email)
	if err != nil {
		// Return success even if email not found — avoid email enumeration
		return nil
	}

	rawToken, tokenHash := generateResetToken()
	pr := &entity.PasswordReset{
		ID:        uuid.NewString(),
		TenantID:  tenantID,
		UserID:    u.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(2 * time.Hour),
	}
	if err := s.resets.Create(ctx, s.db, pr); err != nil {
		return fmt.Errorf("authService.ForgotPassword: create reset: %w", err)
	}

	// TODO: send email with rawToken
	// For now: log the token so devs can test manually
	_ = rawToken
	fmt.Printf("[DEV] Password reset token for %s: %s\n", req.Email, rawToken)

	return nil
}

func (s *authService) ResetPassword(ctx context.Context, req request.ResetPassword) error {
	tokenHash := hashToken(req.Token)
	pr, err := s.resets.FindByTokenHash(ctx, s.db, tokenHash)
	if err != nil {
		return pkgresp.NewAppError(pkgresp.CodeUnauthorized, "token reset tidak valid atau sudah digunakan")
	}

	hash, err := password.Hash(req.Password)
	if err != nil {
		return fmt.Errorf("authService.ResetPassword: hash: %w", err)
	}

	if err := s.users.UpdatePassword(ctx, s.db, pr.UserID, hash); err != nil {
		return fmt.Errorf("authService.ResetPassword: update: %w", err)
	}

	return s.resets.MarkUsed(ctx, s.db, pr.ID)
}

func (s *authService) ChangePassword(ctx context.Context, userID, tenantID string, req request.ChangePassword) error {
	u, err := s.users.FindByID(ctx, s.db, userID)
	if err != nil {
		return pkgresp.NewAppError(pkgresp.CodeNotFound, "user tidak ditemukan")
	}
	if u.PasswordHash == nil {
		return pkgresp.NewAppError(pkgresp.CodeValidation, "akun tidak memiliki password")
	}
	if err := password.Verify(*u.PasswordHash, req.CurrentPassword); err != nil {
		return pkgresp.NewAppError(pkgresp.CodeUnauthorized, "password saat ini salah")
	}

	hash, err := password.Hash(req.NewPassword)
	if err != nil {
		return fmt.Errorf("authService.ChangePassword: hash: %w", err)
	}
	return s.users.UpdatePassword(ctx, s.db, userID, hash)
}

// ---- helpers ----

func (s *authService) issueTokens(ctx context.Context, u *entity.User, userAgent, ip string) (*response.AuthTokens, error) {
	fmt.Println("ISSUE TOKEN START")
	access, err := issueToken(s.jwt, "IssueAccess", u.ID, u.TenantID, u.Email)
	if err != nil {
		fmt.Println("ACCESS TOKEN ERROR:", err)
		return nil, fmt.Errorf("authService.issueTokens: access: %w", err)
	}

	fmt.Println("ACCESS TOKEN OK")

	refresh, err := issueToken(s.jwt, "IssueRefresh", u.ID, u.TenantID, u.Email)
	if err != nil {
		fmt.Println("REFRESH TOKEN ERROR:", err)
		return nil, fmt.Errorf("authService.issueTokens: refresh: %w", err)
	}

	fmt.Println("REFRESH TOKEN OK")

	// Persist refresh token
	session := &entity.UserSession{
		ID:           uuid.NewString(),
		TenantID:     u.TenantID,
		UserID:       u.ID,
		RefreshToken: refresh,
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour),
	}
	if userAgent != "" {
		session.UserAgent = &userAgent
	}
	if ip != "" {
		session.IPAddress = &ip
	}
	if err := s.sessions.Create(ctx, s.db, session); err != nil {
		fmt.Println("SESSION CREATE ERROR:", err)
		return nil, fmt.Errorf("authService.issueTokens: session: %w", err)
	}

	fmt.Println("SESSION CREATED")

	return &response.AuthTokens{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    3600,
		User:         *toUserInfo(u),
	}, nil
}

func toUserInfo(u *entity.User) *response.UserInfo {
	return &response.UserInfo{
		ID:        u.ID,
		TenantID:  u.TenantID,
		Username:  u.Username,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}

func generateResetToken() (raw, hashed string) {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	raw = hex.EncodeToString(b)
	return raw, hashToken(raw)
}

func hashToken(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}

func issueToken(
	jwtService any,
	methodName,
	userID,
	tenantID,
	email string,
) (string, error) {

	if jwtService == nil {
		return "", errors.New("jwt service is nil")
	}

	method := reflect.ValueOf(jwtService).MethodByName(methodName)

	if !method.IsValid() || method.Type().NumIn() != 4 {
		return "", fmt.Errorf(
			"jwt method %s unavailable (numIn=%d)",
			methodName,
			method.Type().NumIn(),
		)
	}

	results := method.Call([]reflect.Value{
		reflect.ValueOf(userID),
		reflect.ValueOf(tenantID),
		reflect.ValueOf(email),
		reflect.ValueOf([]string{}), // roles
	})

	if len(results) != 2 || results[0].Kind() != reflect.String {
		return "", fmt.Errorf(
			"jwt method %s returned an invalid result",
			methodName,
		)
	}

	if !results[1].IsNil() {
		if err, ok := results[1].Interface().(error); ok {
			return "", err
		}
		return "", fmt.Errorf(
			"jwt method %s returned an invalid error",
			methodName,
		)
	}

	return results[0].String(), nil
}

var _ = errors.New // keep import

type refreshClaims struct {
	UserID    string
	TokenType string
}

func parseClaims(jwtService any, token string) (*refreshClaims, error) {
	if jwtService == nil {
		return nil, errors.New("jwt service is nil")
	}

	value := reflect.ValueOf(jwtService)
	for _, methodName := range []string{"ParseClaims", "ParseToken", "Parse"} {
		method := value.MethodByName(methodName)
		if !method.IsValid() || method.Type().NumIn() != 1 {
			continue
		}

		result := method.Call([]reflect.Value{reflect.ValueOf(token)})
		if len(result) < 2 || !result[1].IsNil() {
			return nil, errors.New("invalid token")
		}

		claims := result[0]
		if claims.Kind() == reflect.Ptr {
			if claims.IsNil() {
				return nil, errors.New("invalid token claims")
			}
			claims = claims.Elem()
		}
		if claims.Kind() != reflect.Struct {
			return nil, errors.New("invalid token claims")
		}

		field := func(name string) string {
			v := claims.FieldByName(name)
			if v.IsValid() && v.Kind() == reflect.String {
				return v.String()
			}
			return ""
		}
		return &refreshClaims{UserID: field("UserID"), TokenType: field("TokenType")}, nil
	}

	return nil, errors.New("jwt parser unavailable")
}
