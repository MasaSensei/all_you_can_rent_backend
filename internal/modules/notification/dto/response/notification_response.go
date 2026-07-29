package response

import "time"

// NotificationTemplate is the API-facing shape of entity.NotificationTemplate.
type NotificationTemplate struct {
	ID           string    `json:"id"`
	TenantID     string    `json:"tenant_id"`
	Name         string    `json:"name"`
	Channel      string    `json:"channel"`
	Subject      *string   `json:"subject,omitempty"`
	Body         string    `json:"body"`
	EventTrigger string    `json:"event_trigger"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Notification is the API-facing shape of entity.Notification.
type Notification struct {
	ID         string     `json:"id"`
	TenantID   string     `json:"tenant_id"`
	UserID     *string    `json:"user_id,omitempty"`
	CustomerID *string    `json:"customer_id,omitempty"`
	Channel    string     `json:"channel"`
	Title      string     `json:"title"`
	Message    string     `json:"message"`
	IsRead     bool       `json:"is_read"`
	ReadAt     *time.Time `json:"read_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// NotificationLog is the API-facing shape of entity.NotificationLog.
type NotificationLog struct {
	ID             string     `json:"id"`
	TenantID       string     `json:"tenant_id"`
	RecipientID    string     `json:"recipient_id"`
	RecipientType  string     `json:"recipient_type"`
	Channel        string     `json:"channel"`
	DeliveryStatus string     `json:"delivery_status"`
	ErrorMessage   *string    `json:"error_message,omitempty"`
	SentAt         *time.Time `json:"sent_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}
