package service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"rentos-backend/internal/modules/notification/dto/request"
	"rentos-backend/internal/modules/notification/dto/response"
	"rentos-backend/internal/modules/notification/entity"
	"rentos-backend/internal/modules/notification/repository"
	pkgresponse "rentos-backend/pkg/response"
)

// ============================================================
// templateService
// ============================================================

type templateService struct {
	db   *sqlx.DB
	repo repository.TemplateRepository
}

func NewTemplateService(db *sqlx.DB, repo repository.TemplateRepository) TemplateService {
	return &templateService{db: db, repo: repo}
}

func (s *templateService) Create(ctx context.Context, tenantID, actorID string, req request.CreateTemplate) (*response.NotificationTemplate, error) {
	t := &entity.NotificationTemplate{
		ID:           uuid.NewString(),
		TenantID:     tenantID,
		Name:         req.Name,
		Channel:      req.Channel,
		Subject:      req.Subject,
		Body:         req.Body,
		EventTrigger: req.EventTrigger,
		CreatedBy:    &actorID,
	}
	if err := s.repo.Create(ctx, s.db, t); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, t.ID, tenantID)
}

func (s *templateService) GetByID(ctx context.Context, id, tenantID string) (*response.NotificationTemplate, error) {
	t, err := s.repo.FindByID(ctx, s.db, id, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "notification template not found")
		}
		return nil, err
	}
	return toTemplateResponse(t), nil
}

func (s *templateService) List(ctx context.Context, tenantID string) ([]response.NotificationTemplate, error) {
	templates, err := s.repo.List(ctx, s.db, tenantID)
	if err != nil {
		return nil, err
	}
	out := make([]response.NotificationTemplate, 0, len(templates))
	for _, t := range templates {
		out = append(out, *toTemplateResponse(&t))
	}
	return out, nil
}

func (s *templateService) Update(ctx context.Context, id, tenantID, actorID string, req request.UpdateTemplate) (*response.NotificationTemplate, error) {
	t := &entity.NotificationTemplate{
		ID:        id,
		TenantID:  tenantID,
		Name:      derefStr(req.Name),
		Subject:   req.Subject,
		Body:      derefStr(req.Body),
		UpdatedBy: &actorID,
	}
	if err := s.repo.Update(ctx, s.db, t); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "notification template not found")
		}
		return nil, err
	}
	return s.GetByID(ctx, id, tenantID)
}

func (s *templateService) Delete(ctx context.Context, id, tenantID string) error {
	if err := s.repo.Delete(ctx, s.db, id, tenantID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return pkgresponse.NewAppError(pkgresponse.CodeNotFound, "notification template not found")
		}
		return err
	}
	return nil
}

// ============================================================
// notificationService
// ============================================================

type notificationService struct {
	db        *sqlx.DB
	templates repository.TemplateRepository
	notifs    repository.NotificationRepository
	logs      repository.LogRepository
}

func NewNotificationService(
	db *sqlx.DB,
	templates repository.TemplateRepository,
	notifs repository.NotificationRepository,
	logs repository.LogRepository,
) NotificationService {
	return &notificationService{db: db, templates: templates, notifs: notifs, logs: logs}
}

// Send resolves the template, renders it, writes a log entry with
// status=pending, then returns. Actual delivery (email/whatsapp/sms)
// is handled by background workers that update the log via UpdateStatus.
// in_app channel writes directly to the notifications table.
func (s *notificationService) Send(ctx context.Context, req request.SendNotification) error {
	channels := []string{entity.ChannelEmail, entity.ChannelWhatsApp, entity.ChannelInApp, entity.ChannelSMS}

	for _, ch := range channels {
		tmpl, err := s.templates.FindByTrigger(ctx, s.db, req.TenantID, req.EventTrigger, ch)
		if errors.Is(err, repository.ErrNotFound) {
			continue // no template configured for this channel — skip silently
		}
		if err != nil {
			return err
		}

		rendered := renderTemplate(tmpl.Body, req.Data)

		if ch == entity.ChannelInApp {
			subject := derefStr(tmpl.Subject)
			n := &entity.Notification{
				ID:       uuid.NewString(),
				TenantID: req.TenantID,
				Channel:  entity.ChannelInApp,
				Title:    subject,
				Message:  rendered,
			}
			if req.RecipientType == entity.RecipientTypeUser {
				n.UserID = &req.RecipientID
			} else {
				n.CustomerID = &req.RecipientID
			}
			if err := s.notifs.Create(ctx, s.db, n); err != nil {
				return err
			}
			continue
		}

		// Async channels: write a pending log — worker will pick it up.
		log := &entity.NotificationLog{
			ID:                     uuid.NewString(),
			TenantID:               req.TenantID,
			NotificationTemplateID: &tmpl.ID,
			RecipientID:            req.RecipientID,
			RecipientType:          req.RecipientType,
			Channel:                ch,
			DeliveryStatus:         entity.DeliveryStatusPending,
		}
		if err := s.logs.Create(ctx, s.db, log); err != nil {
			return err
		}
		// TODO: enqueue log.ID to the channel-specific worker queue (Phase 12 workers)
	}
	return nil
}

func (s *notificationService) CreateInApp(ctx context.Context, tenantID string, req request.CreateInAppNotification) (*response.Notification, error) {
	n := &entity.Notification{
		ID:         uuid.NewString(),
		TenantID:   tenantID,
		UserID:     req.UserID,
		CustomerID: req.CustomerID,
		Channel:    entity.ChannelInApp,
		Title:      req.Title,
		Message:    req.Message,
	}
	if err := s.notifs.Create(ctx, s.db, n); err != nil {
		return nil, err
	}
	return toNotificationResponse(n), nil
}

func (s *notificationService) List(ctx context.Context, tenantID string, filter request.ListNotificationsFilter, userID, customerID *string) ([]response.Notification, error) {
	perPage, page := normPage(filter.PerPage, filter.Page)
	items, err := s.notifs.List(ctx, s.db, tenantID, userID, customerID, filter.IsRead, perPage, (page-1)*perPage)
	if err != nil {
		return nil, err
	}
	out := make([]response.Notification, 0, len(items))
	for _, n := range items {
		out = append(out, *toNotificationResponse(&n))
	}
	return out, nil
}

func (s *notificationService) MarkRead(ctx context.Context, id, tenantID string) error {
	return s.notifs.MarkRead(ctx, s.db, id, tenantID)
}

func (s *notificationService) MarkAllRead(ctx context.Context, tenantID string, userID, customerID *string) error {
	return s.notifs.MarkAllRead(ctx, s.db, tenantID, userID, customerID)
}

// ============================================================
// logService
// ============================================================

type logService struct {
	db   *sqlx.DB
	repo repository.LogRepository
}

func NewLogService(db *sqlx.DB, repo repository.LogRepository) LogService {
	return &logService{db: db, repo: repo}
}

func (s *logService) List(ctx context.Context, tenantID string, page, perPage int) ([]response.NotificationLog, error) {
	perPage, page = normPage(perPage, page)
	logs, err := s.repo.List(ctx, s.db, tenantID, perPage, (page-1)*perPage)
	if err != nil {
		return nil, err
	}
	out := make([]response.NotificationLog, 0, len(logs))
	for _, l := range logs {
		out = append(out, *toLogResponse(&l))
	}
	return out, nil
}

func (s *logService) UpdateStatus(ctx context.Context, id, deliveryStatus string, errMsg *string) error {
	return s.repo.UpdateStatus(ctx, s.db, id, deliveryStatus, errMsg)
}

// ============================================================
// helpers
// ============================================================

// renderTemplate performs simple {{key}} substitution from data map.
func renderTemplate(body string, data map[string]string) string {
	for k, v := range data {
		body = strings.ReplaceAll(body, "{{"+k+"}}", v)
	}
	return body
}

func toTemplateResponse(t *entity.NotificationTemplate) *response.NotificationTemplate {
	return &response.NotificationTemplate{
		ID: t.ID, TenantID: t.TenantID, Name: t.Name, Channel: t.Channel,
		Subject: t.Subject, Body: t.Body, EventTrigger: t.EventTrigger,
		Status: t.Status, CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt,
	}
}

func toNotificationResponse(n *entity.Notification) *response.Notification {
	return &response.Notification{
		ID: n.ID, TenantID: n.TenantID, UserID: n.UserID, CustomerID: n.CustomerID,
		Channel: n.Channel, Title: n.Title, Message: n.Message,
		IsRead: n.IsRead, ReadAt: n.ReadAt, CreatedAt: n.CreatedAt,
	}
}

func toLogResponse(l *entity.NotificationLog) *response.NotificationLog {
	return &response.NotificationLog{
		ID: l.ID, TenantID: l.TenantID, RecipientID: l.RecipientID,
		RecipientType: l.RecipientType, Channel: l.Channel,
		DeliveryStatus: l.DeliveryStatus, ErrorMessage: l.ErrorMessage,
		SentAt: l.SentAt, CreatedAt: l.CreatedAt,
	}
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func normPage(perPage, page int) (int, int) {
	if perPage <= 0 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}
	if page <= 0 {
		page = 1
	}
	return perPage, page
}
