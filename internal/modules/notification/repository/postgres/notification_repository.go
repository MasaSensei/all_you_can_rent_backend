package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"rentos-backend/internal/modules/notification/entity"
	"rentos-backend/internal/modules/notification/repository"
	"rentos-backend/pkg/database"
)

// ============================================================
// templateRepository
// ============================================================

type templateRepository struct {
	qCreate         string
	qFindByID       string
	qFindByTrigger  string
	qList           string
	qUpdate         string
	qDelete         string
}

func NewTemplateRepository(qCreate, qFindByID, qFindByTrigger, qList, qUpdate, qDelete string) repository.TemplateRepository {
	return &templateRepository{
		qCreate: qCreate, qFindByID: qFindByID, qFindByTrigger: qFindByTrigger,
		qList: qList, qUpdate: qUpdate, qDelete: qDelete,
	}
}

func (r *templateRepository) Create(ctx context.Context, q database.Querier, t *entity.NotificationTemplate) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		t.ID, t.TenantID, t.Name, t.Channel, t.Subject, t.Body, t.EventTrigger, t.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("templateRepository.Create: %w", err)
	}
	return nil
}

func (r *templateRepository) FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.NotificationTemplate, error) {
	var t entity.NotificationTemplate
	if err := q.GetContext(ctx, &t, r.qFindByID, id, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("templateRepository.FindByID: %w", err)
	}
	return &t, nil
}

func (r *templateRepository) FindByTrigger(ctx context.Context, q database.Querier, tenantID, trigger, channel string) (*entity.NotificationTemplate, error) {
	var t entity.NotificationTemplate
	if err := q.GetContext(ctx, &t, r.qFindByTrigger, tenantID, trigger, channel); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("templateRepository.FindByTrigger: %w", err)
	}
	return &t, nil
}

func (r *templateRepository) List(ctx context.Context, q database.Querier, tenantID string) ([]entity.NotificationTemplate, error) {
	var out []entity.NotificationTemplate
	if err := q.SelectContext(ctx, &out, r.qList, tenantID); err != nil {
		return nil, fmt.Errorf("templateRepository.List: %w", err)
	}
	return out, nil
}

func (r *templateRepository) Update(ctx context.Context, q database.Querier, t *entity.NotificationTemplate) error {
	res, err := q.ExecContext(ctx, r.qUpdate, t.ID, t.TenantID, t.Name, t.Subject, t.Body)
	if err != nil {
		return fmt.Errorf("templateRepository.Update: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *templateRepository) Delete(ctx context.Context, q database.Querier, id, tenantID string) error {
	res, err := q.ExecContext(ctx, r.qDelete, id, tenantID)
	if err != nil {
		return fmt.Errorf("templateRepository.Delete: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// ============================================================
// notificationRepository
// ============================================================

type notificationRepository struct {
	qCreate      string
	qList        string
	qMarkRead    string
	qMarkAllRead string
}

func NewNotificationRepository(qCreate, qList, qMarkRead, qMarkAllRead string) repository.NotificationRepository {
	return &notificationRepository{
		qCreate: qCreate, qList: qList, qMarkRead: qMarkRead, qMarkAllRead: qMarkAllRead,
	}
}

func (r *notificationRepository) Create(ctx context.Context, q database.Querier, n *entity.Notification) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		n.ID, n.TenantID, n.UserID, n.CustomerID, n.Channel, n.Title, n.Message,
	)
	if err != nil {
		return fmt.Errorf("notificationRepository.Create: %w", err)
	}
	return nil
}

func (r *notificationRepository) List(ctx context.Context, q database.Querier, tenantID string, userID, customerID *string, isRead *bool, limit, offset int) ([]entity.Notification, error) {
	var out []entity.Notification
	if err := q.SelectContext(ctx, &out, r.qList, tenantID, userID, customerID, isRead, limit, offset); err != nil {
		return nil, fmt.Errorf("notificationRepository.List: %w", err)
	}
	return out, nil
}

func (r *notificationRepository) MarkRead(ctx context.Context, q database.Querier, id, tenantID string) error {
	_, err := q.ExecContext(ctx, r.qMarkRead, id, tenantID)
	if err != nil {
		return fmt.Errorf("notificationRepository.MarkRead: %w", err)
	}
	return nil
}

func (r *notificationRepository) MarkAllRead(ctx context.Context, q database.Querier, tenantID string, userID, customerID *string) error {
	_, err := q.ExecContext(ctx, r.qMarkAllRead, tenantID, userID, customerID)
	if err != nil {
		return fmt.Errorf("notificationRepository.MarkAllRead: %w", err)
	}
	return nil
}

// ============================================================
// logRepository
// ============================================================

type logRepository struct {
	qCreate       string
	qUpdateStatus string
	qList         string
}

func NewLogRepository(qCreate, qUpdateStatus, qList string) repository.LogRepository {
	return &logRepository{qCreate: qCreate, qUpdateStatus: qUpdateStatus, qList: qList}
}

func (r *logRepository) Create(ctx context.Context, q database.Querier, l *entity.NotificationLog) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		l.ID, l.TenantID, l.NotificationTemplateID, l.RecipientID, l.RecipientType,
		l.Channel, l.DeliveryStatus, l.ErrorMessage, l.SentAt,
	)
	if err != nil {
		return fmt.Errorf("logRepository.Create: %w", err)
	}
	return nil
}

func (r *logRepository) UpdateStatus(ctx context.Context, q database.Querier, id, deliveryStatus string, errMsg *string) error {
	_, err := q.ExecContext(ctx, r.qUpdateStatus, id, deliveryStatus, errMsg)
	if err != nil {
		return fmt.Errorf("logRepository.UpdateStatus: %w", err)
	}
	return nil
}

func (r *logRepository) List(ctx context.Context, q database.Querier, tenantID string, limit, offset int) ([]entity.NotificationLog, error) {
	var out []entity.NotificationLog
	if err := q.SelectContext(ctx, &out, r.qList, tenantID, limit, offset); err != nil {
		return nil, fmt.Errorf("logRepository.List: %w", err)
	}
	return out, nil
}
