package response

// AuthTokens is returned on successful login or refresh.
type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // seconds
	User         UserInfo `json:"user"`
}

type UserInfo struct {
	ID        string  `json:"id"`
	TenantID  string  `json:"tenant_id"`
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
}
