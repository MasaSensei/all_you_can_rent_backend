package response

import "time"

// APIKey is the API-facing shape of entity.APIKey.
// RawKey is only populated on creation; never stored or returned again.
type APIKey struct {
	ID         string     `json:"id"`
	TenantID   string     `json:"tenant_id"`
	Name       string     `json:"name"`
	KeyPrefix  string     `json:"key_prefix"`
	Scopes     []string   `json:"scopes"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
}

// APIKeyCreated extends APIKey with the one-time raw key.
type APIKeyCreated struct {
	APIKey
	RawKey string `json:"key"`
}

// Webhook is the API-facing shape of entity.Webhook.
// Secret is never returned after creation.
type Webhook struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	URL       string    `json:"url"`
	Events    []string  `json:"events"`
	IsActive  bool      `json:"is_active"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// WebhookLog is the API-facing shape of entity.WebhookLog.
type WebhookLog struct {
	ID           string     `json:"id"`
	WebhookID    string     `json:"webhook_id"`
	EventType    string     `json:"event_type"`
	ResponseCode *int       `json:"response_code,omitempty"`
	TriggeredAt  time.Time  `json:"triggered_at"`
	Status       string     `json:"status"`
}
