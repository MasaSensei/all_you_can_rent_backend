package repository

import (
	"context"
	"errors"

	"rentos-backend/internal/modules/notification/entity"
	"rentos-backend/pkg/database"
)

var ErrNotFound = errors.New("repository: record not found")

// TemplateRepository manages the notification_templates table.
type TemplateRepository interface {
	Create(ctx context.Context, q database.Querier, t *entity.NotificationTemplate) error
	FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.NotificationTemplate, error)
	FindByTrigger(ctx context.Context, q database.Querier, tenantID, trigger, channel string) (*entity.NotificationTemplate, error)
	List(ctx context.Context, q database.Querier, tenantID string) ([]entity.NotificationTemplate, error)
	Update(ctx context.Context, q database.Querier, t *entity.NotificationTemplate) error
	Delete(ctx context.Context, q database.Querier, id, tenantID string) error
}

// NotificationRepository manages the notifications (in-app) table.
type NotificationRepository interface {
	Create(ctx context.Context, q database.Querier, n *entity.Notification) error
	List(ctx context.Context, q database.Querier, tenantID string, userID, customerID *string, isRead *bool, limit, offset int) ([]entity.Notification, error)
	MarkRead(ctx context.Context, q database.Querier, id, tenantID string) error
	MarkAllRead(ctx context.Context, q database.Querier, tenantID string, userID, customerID *string) error
}

// LogRepository manages the notification_logs table.
type LogRepository interface {
	Create(ctx context.Context, q database.Querier, l *entity.NotificationLog) error
	UpdateStatus(ctx context.Context, q database.Querier, id, deliveryStatus string, errMsg *string) error
	List(ctx context.Context, q database.Querier, tenantID string, limit, offset int) ([]entity.NotificationLog, error)
}
