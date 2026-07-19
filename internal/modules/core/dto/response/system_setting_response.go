package response

import "time"

// SystemSetting is the API-facing shape of entity.SystemSetting.
type SystemSetting struct {
	ID        string    `json:"id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	Type      string    `json:"type"`
	UpdatedAt time.Time `json:"updated_at"`
}
