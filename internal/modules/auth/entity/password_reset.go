package entity

import "time"

// PasswordReset mirrors the password_resets table.
type PasswordReset struct {
	ID        string     `db:"id"`
	TenantID  string     `db:"tenant_id"`
	UserID    string     `db:"user_id"`
	Token     string     `db:"token"`
	ExpiresAt time.Time  `db:"expires_at"`
	Used      bool       `db:"used"`
	Status    string     `db:"status"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
	Version   int        `db:"version"`
}
