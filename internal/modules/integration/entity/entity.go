package entity

import "time"

// APIKey mirrors the api_keys table.
// key_hash is stored; the raw key is shown once on creation only.
type APIKey struct {
	ID          string     `db:"id"`
	TenantID    string     `db:"tenant_id"`
	Name        string     `db:"name"`
	KeyPrefix   string     `db:"key_prefix"`
	KeyHash     string     `db:"key_hash"`
	Scopes      *string    `db:"scopes"` // JSON array
	LastUsedAt  *time.Time `db:"last_used_at"`
	ExpiresAt   *time.Time `db:"expires_at"`
	Status      string     `db:"status"`
	CreatedBy   *string    `db:"created_by"`
	UpdatedBy   *string    `db:"updated_by"`
	DeletedBy   *string    `db:"deleted_by"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
	Version     int        `db:"version"`
}

// Webhook mirrors the webhooks table.
type Webhook struct {
	ID        string     `db:"id"`
	TenantID  string     `db:"tenant_id"`
	URL       string     `db:"url"`
	Events    string     `db:"events"` // JSON array of event names
	Secret    string     `db:"secret"`
	IsActive  bool       `db:"is_active"`
	Status    string     `db:"status"`
	CreatedBy *string    `db:"created_by"`
	UpdatedBy *string    `db:"updated_by"`
	DeletedBy *string    `db:"deleted_by"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
	Version   int        `db:"version"`
}

// WebhookLog mirrors the webhook_logs table.
type WebhookLog struct {
	ID           string     `db:"id"`
	TenantID     string     `db:"tenant_id"`
	WebhookID    string     `db:"webhook_id"`
	EventType    string     `db:"event_type"`
	Payload      *string    `db:"payload"` // JSON
	ResponseCode *int       `db:"response_code"`
	ResponseBody *string    `db:"response_body"`
	TriggeredAt  time.Time  `db:"triggered_at"`
	Status       string     `db:"status"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
	Version      int        `db:"version"`
}
