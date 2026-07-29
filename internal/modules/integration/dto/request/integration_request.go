package request

import "time"

// CreateAPIKey creates a new API key.
type CreateAPIKey struct {
	Name      string     `json:"name" validate:"required,max=150"`
	Scopes    []string   `json:"scopes" validate:"required,min=1,dive,max=100"`
	ExpiresAt *time.Time `json:"expires_at"`
}

// CreateWebhook registers a new webhook endpoint.
type CreateWebhook struct {
	URL    string   `json:"url" validate:"required,url,max=500"`
	Events []string `json:"events" validate:"required,min=1,dive,max=100"`
	Secret string   `json:"secret" validate:"required,min=16,max=255"`
}

// UpdateWebhook updates an existing webhook.
type UpdateWebhook struct {
	URL      *string  `json:"url" validate:"omitempty,url,max=500"`
	Events   []string `json:"events" validate:"omitempty,min=1,dive,max=100"`
	IsActive *bool    `json:"is_active"`
}

// DispatchWebhook is used internally to trigger webhook delivery.
// Not exposed over HTTP — called by domain event handlers.
type DispatchWebhook struct {
	TenantID  string            `json:"tenant_id"`
	EventType string            `json:"event_type"`
	Payload   map[string]any    `json:"payload"`
}
