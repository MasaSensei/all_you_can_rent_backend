package response

import "time"

// Auth is returned after a successful login or token refresh.
type Auth struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	User         User   `json:"user"`
}

// User is the API-facing shape of entity.User.
// PasswordHash is never included.
type User struct {
	ID          string     `json:"id"`
	TenantID    string     `json:"tenant_id"`
	Email       string     `json:"email"`
	FirstName   *string    `json:"first_name,omitempty"`
	LastName    *string    `json:"last_name,omitempty"`
	Phone       *string    `json:"phone,omitempty"`
	AvatarURL   *string    `json:"avatar_url,omitempty"`
	IsActive    bool       `json:"is_active"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
