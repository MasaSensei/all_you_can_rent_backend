package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("jwt: invalid or expired token")

type Config struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

type Claims struct {
	UserID   string   `json:"user_id"`
	TenantID string   `json:"tenant_id"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

type Service struct {
	cfg Config
}

func New(cfg Config) *Service {
	return &Service{cfg: cfg}
}

func (s *Service) IssueAccess(userID, tenantID, email string, roles []string) (string, error) {
	return s.sign(userID, tenantID, email, roles, s.cfg.AccessTTL, s.cfg.AccessSecret)
}

func (s *Service) IssueRefresh(userID, tenantID, email string, roles []string) (string, error) {
	return s.sign(userID, tenantID, email, roles, s.cfg.RefreshTTL, s.cfg.RefreshSecret)
}

func (s *Service) ParseAccess(tokenStr string) (*Claims, error) {
	return s.parse(tokenStr, s.cfg.AccessSecret)
}

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
