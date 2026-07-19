package entity

import "time"

// User mirrors the users table.
type User struct {
	ID           string     `db:"id"`
	TenantID     string     `db:"tenant_id"`
	Email        string     `db:"email"`
	Username     *string    `db:"username"`
	PasswordHash string     `db:"password_hash"`
	FirstName    *string    `db:"first_name"`
	LastName     *string    `db:"last_name"`
	Phone        *string    `db:"phone"`
	AvatarURL    *string    `db:"avatar_url"`
	IsActive     bool       `db:"is_active"`
	LastLoginAt  *time.Time `db:"last_login_at"`
	Status       string     `db:"status"`
	CreatedBy    *string    `db:"created_by"`
	UpdatedBy    *string    `db:"updated_by"`
	DeletedBy    *string    `db:"deleted_by"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
	Version      int        `db:"version"`
}
