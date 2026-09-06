package entity

import "time"

type User struct {
	ID              string     `db:"id"`
	TenantID        string     `db:"tenant_id"`
	Username        string     `db:"username"`
	Email           string     `db:"email"`
	PasswordHash    *string    `db:"password_hash"`
	FirstName       *string    `db:"first_name"`
	LastName        *string    `db:"last_name"`
	Phone           *string    `db:"phone"`
	AvatarURL       *string    `db:"avatar_url"`
	EmailVerifiedAt *time.Time `db:"email_verified_at"`
	IsActive        bool       `db:"is_active"`
	LastLoginAt     *time.Time `db:"last_login_at"`
	Status          string     `db:"status"`
	CreatedBy       *string    `db:"created_by"`
	UpdatedBy       *string    `db:"updated_by"`
	DeletedBy       *string    `db:"deleted_by"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at"`
	Version         int        `db:"version"`
}

type UserSession struct {
	ID           string     `db:"id"`
	TenantID     string     `db:"tenant_id"`
	UserID       string     `db:"user_id"`
	RefreshToken string     `db:"refresh_token"`
	UserAgent    *string    `db:"user_agent"`
	IPAddress    *string    `db:"ip_address"`
	ExpiresAt    time.Time  `db:"expires_at"`
	RevokedAt    *time.Time `db:"revoked_at"`
	Status       string     `db:"status"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
}

type PasswordReset struct {
	ID        string     `db:"id"`
	TenantID  string     `db:"tenant_id"`
	UserID    string     `db:"user_id"`
	TokenHash string     `db:"token_hash"`
	ExpiresAt time.Time  `db:"expires_at"`
	UsedAt    *time.Time `db:"used_at"`
	Status    string     `db:"status"`
	CreatedAt time.Time  `db:"created_at"`
}
