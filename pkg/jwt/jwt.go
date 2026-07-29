// Package jwt issues and parses signed JWTs for the RentOS API.
// Every other package that needs token data receives a *Claims value —
// nothing outside this package imports a JWT library directly.
package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ErrInvalidToken is returned when a token cannot be parsed or its
// signature is invalid.
var ErrInvalidToken = errors.New("jwt: invalid or expired token")

// Config holds signing configuration. Use separate secrets for access
// and refresh tokens so a leaked refresh secret cannot forge access
// tokens.
type Config struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

// Claims is the payload embedded in every token. It is the only type
// middleware and services work with — they never touch raw jwt.MapClaims.
type Claims struct {
	UserID   string   `json:"user_id"`
	TenantID string   `json:"tenant_id"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// Service issues and validates JWTs.
type Service struct {
	cfg Config
}

// New builds a JWT Service from the given Config.
func New(cfg Config) *Service {
	return &Service{cfg: cfg}
}

// IssueAccess mints a short-lived access token.
func (s *Service) IssueAccess(userID, tenantID, email string, roles []string) (string, error) {
	return s.sign(userID, tenantID, email, roles, s.cfg.AccessTTL, s.cfg.AccessSecret)
}

// IssueRefresh mints a long-lived refresh token.
func (s *Service) IssueRefresh(userID, tenantID, email string, roles []string) (string, error) {
	return s.sign(userID, tenantID, email, roles, s.cfg.RefreshTTL, s.cfg.RefreshSecret)
}

// ParseAccess validates an access token and returns its claims.
func (s *Service) ParseAccess(tokenStr string) (*Claims, error) {
	return s.parse(tokenStr, s.cfg.AccessSecret)
}

// ParseRefresh validates a refresh token and returns its claims.
func (s *Service) ParseRefresh(tokenStr string) (*Claims, error) {
	return s.parse(tokenStr, s.cfg.RefreshSecret)
}

func (s *Service) sign(userID, tenantID, email string, roles []string, ttl time.Duration, secret string) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID:   userID,
		TenantID: tenantID,
		Email:    email,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("jwt: sign: %w", err)
	}
	return signed, nil
}

func (s *Service) parse(tokenStr, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
