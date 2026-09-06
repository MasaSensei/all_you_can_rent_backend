package request

// Login authenticates a tenant user.
type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// Refresh exchanges a refresh token for a new access token.
type Refresh struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// Logout revokes a refresh token.
type Logout struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ForgotPassword initiates the reset flow.
type ForgotPassword struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPassword completes the reset flow.
type ResetPassword struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

// ChangePassword changes password for authenticated user.
type ChangePassword struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}
