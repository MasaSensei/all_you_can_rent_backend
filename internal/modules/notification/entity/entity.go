package entity

import "time"

// NotificationTemplate mirrors the notification_templates table.
// Templates are resolved by event_trigger at send time.
type NotificationTemplate struct {
	ID           string     `db:"id"`
	TenantID     string     `db:"tenant_id"`
	Name         string     `db:"name"`
	Channel      string     `db:"channel"`      // email, whatsapp, in_app, sms
	Subject      *string    `db:"subject"`
	Body         string     `db:"body"`
	EventTrigger string     `db:"event_trigger"` // booking.confirmed, payment.succeeded, etc.
	Status       string     `db:"status"`
	CreatedBy    *string    `db:"created_by"`
	UpdatedBy    *string    `db:"updated_by"`
	DeletedBy    *string    `db:"deleted_by"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
	Version      int        `db:"version"`
}

// Notification mirrors the notifications table (in-app only).
type Notification struct {
	ID         string     `db:"id"`
	TenantID   string     `db:"tenant_id"`
	UserID     *string    `db:"user_id"`
	CustomerID *string    `db:"customer_id"`
	Channel    string     `db:"channel"`
	Title      string     `db:"title"`
	Message    string     `db:"message"`
	IsRead     bool       `db:"is_read"`
	ReadAt     *time.Time `db:"read_at"`
	Status     string     `db:"status"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at"`
	Version    int        `db:"version"`
}

// NotificationLog mirrors the notification_logs table.
type NotificationLog struct {
	ID                     string     `db:"id"`
	TenantID               string     `db:"tenant_id"`
	NotificationTemplateID *string    `db:"notification_template_id"`
	RecipientID            string     `db:"recipient_id"`
	RecipientType          string     `db:"recipient_type"` // user, customer
	Channel                string     `db:"channel"`
	DeliveryStatus         string     `db:"delivery_status"` // pending, sent, failed
	ErrorMessage           *string    `db:"error_message"`
	SentAt                 *time.Time `db:"sent_at"`
	Status                 string     `db:"status"`
	CreatedAt              time.Time  `db:"created_at"`
	UpdatedAt              time.Time  `db:"updated_at"`
	DeletedAt              *time.Time `db:"deleted_at"`
	Version                int        `db:"version"`
}

const (
	ChannelEmail    = "email"
	ChannelWhatsApp = "whatsapp"
	ChannelInApp    = "in_app"
	ChannelSMS      = "sms"

	DeliveryStatusPending = "pending"
	DeliveryStatusSent    = "sent"
	DeliveryStatusFailed  = "failed"

	RecipientTypeUser     = "user"
	RecipientTypeCustomer = "customer"
)
