package request

// Register is the input for creating a new user account.
type Register struct {
	Email     string  `json:"email" validate:"required,email,max=255"`
	Password  string  `json:"password" validate:"required,min=8,max=72"`
	FirstName string  `json:"first_name" validate:"required,max=100"`
	LastName  string  `json:"last_name" validate:"required,max=100"`
	Phone     *string `json:"phone" validate:"omitempty,max=30"`
}

// Login is the input for authenticating an existing user.
type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RefreshToken is the input for rotating an access token.
type RefreshToken struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ForgotPassword initiates the password-reset flow.
type ForgotPassword struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPassword completes the password-reset flow.
type ResetPassword struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

// UpdateUser is the input for updating a user's profile.
type UpdateUser struct {
	FirstName *string `json:"first_name" validate:"omitempty,max=100"`
	LastName  *string `json:"last_name" validate:"omitempty,max=100"`
	Phone     *string `json:"phone" validate:"omitempty,max=30"`
}
