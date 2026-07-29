package request

// CreateTemplate creates a notification template.
type CreateTemplate struct {
	Name         string  `json:"name" validate:"required,max=150"`
	Channel      string  `json:"channel" validate:"required,oneof=email whatsapp in_app sms"`
	Subject      *string `json:"subject" validate:"omitempty,max=255"`
	Body         string  `json:"body" validate:"required"`
	EventTrigger string  `json:"event_trigger" validate:"required,max=100"`
}

// UpdateTemplate updates an existing template.
type UpdateTemplate struct {
	Name    *string `json:"name" validate:"omitempty,max=150"`
	Subject *string `json:"subject" validate:"omitempty,max=255"`
	Body    *string `json:"body" validate:"omitempty"`
}

// SendNotification sends a notification using a template.
// Used internally by other services — not a public HTTP endpoint.
type SendNotification struct {
	TenantID     string            `json:"tenant_id" validate:"required,uuid"`
	EventTrigger string            `json:"event_trigger" validate:"required"`
	RecipientID  string            `json:"recipient_id" validate:"required,uuid"`
	RecipientType string           `json:"recipient_type" validate:"required,oneof=user customer"`
	Data         map[string]string `json:"data"` // template variable substitution
}

// CreateInAppNotification creates a direct in-app notification.
type CreateInAppNotification struct {
	UserID     *string `json:"user_id" validate:"omitempty,uuid"`
	CustomerID *string `json:"customer_id" validate:"omitempty,uuid"`
	Title      string  `json:"title" validate:"required,max=255"`
	Message    string  `json:"message" validate:"required"`
}

// ListNotificationsFilter holds whitelisted query-param filters.
type ListNotificationsFilter struct {
	IsRead *bool
	Page   int
	PerPage int
}
