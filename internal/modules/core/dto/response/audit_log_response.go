package response

import "time"

// AuditLog is the API-facing shape of entity.AuditLog.
type AuditLog struct {
	ID         string    `json:"id"`
	UserID     *string   `json:"user_id,omitempty"`
	EntityType string    `json:"entity_type"`
	EntityID   string    `json:"entity_id"`
	Action     string    `json:"action"`
	CreatedAt  time.Time `json:"created_at"`
}
