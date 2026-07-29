package service

import (
	"context"

	"rentos-backend/internal/modules/notification/dto/request"
	"rentos-backend/internal/modules/notification/dto/response"
)

// TemplateService manages notification templates.
type TemplateService interface {
	Create(ctx context.Context, tenantID, actorID string, req request.CreateTemplate) (*response.NotificationTemplate, error)
	GetByID(ctx context.Context, id, tenantID string) (*response.NotificationTemplate, error)
	List(ctx context.Context, tenantID string) ([]response.NotificationTemplate, error)
	Update(ctx context.Context, id, tenantID, actorID string, req request.UpdateTemplate) (*response.NotificationTemplate, error)
	Delete(ctx context.Context, id, tenantID string) error
}

// NotificationService manages in-app notifications and dispatch.
type NotificationService interface {
	// Send resolves the template for the given event_trigger + channel,
	// renders it with data, persists the log, and enqueues delivery.
	// Other modules call this — it is the single dispatch entrypoint.
	Send(ctx context.Context, req request.SendNotification) error

	// CreateInApp creates a direct in-app notification (no template needed).
	CreateInApp(ctx context.Context, tenantID string, req request.CreateInAppNotification) (*response.Notification, error)

	// List returns in-app notifications for a user or customer feed.
	List(ctx context.Context, tenantID string, filter request.ListNotificationsFilter, userID, customerID *string) ([]response.Notification, error)

	MarkRead(ctx context.Context, id, tenantID string) error
	MarkAllRead(ctx context.Context, tenantID string, userID, customerID *string) error
}

// LogService manages notification delivery logs.
type LogService interface {
	List(ctx context.Context, tenantID string, page, perPage int) ([]response.NotificationLog, error)
	// UpdateStatus is called by delivery workers (email/whatsapp/sms)
	// after attempting dispatch.
	UpdateStatus(ctx context.Context, id, deliveryStatus string, errMsg *string) error
}
