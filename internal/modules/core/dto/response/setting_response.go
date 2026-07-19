package response

import "time"

// Setting is the API-facing shape of entity.Setting.
type Setting struct {
	ID        string    `json:"id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	Type      string    `json:"type"`
	UpdatedAt time.Time `json:"updated_at"`
}
