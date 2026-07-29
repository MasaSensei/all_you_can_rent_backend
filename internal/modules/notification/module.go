package notification

import (
	"github.com/gofiber/fiber/v2"

	"rentos-backend/internal/bootstrap"
	"rentos-backend/internal/modules/notification/handler"
	"rentos-backend/internal/modules/notification/repository/postgres"
	"rentos-backend/internal/modules/notification/routes"
	"rentos-backend/internal/modules/notification/service"
)

// Module holds the notification module's wired handler and services.
type Module struct {
	handler      *handler.Handler
	notifService service.NotificationService
	logService   service.LogService
}

// New builds the notification module: repositories → services → handler.
func New(c *bootstrap.Container) *Module {
	templateRepo := postgres.NewTemplateRepository(
		query("create_template.sql"),
		query("find_template_by_id.sql"),
		query("find_template_by_trigger.sql"),
		query("list_templates.sql"),
		query("update_template.sql"),
		query("delete_template.sql"),
	)
	notifRepo := postgres.NewNotificationRepository(
		query("create_notification.sql"),
		query("list_notifications.sql"),
		query("mark_notification_read.sql"),
		query("mark_all_read.sql"),
	)
	logRepo := postgres.NewLogRepository(
		query("create_notification_log.sql"),
		query("update_log_status.sql"),
		query("list_notification_logs.sql"),
	)

	templateSvc := service.NewTemplateService(c.DB, templateRepo)
	notifSvc := service.NewNotificationService(c.DB, templateRepo, notifRepo, logRepo)
	logSvc := service.NewLogService(c.DB, logRepo)

	h := handler.New(templateSvc, notifSvc, logSvc, c.Validator)
	return &Module{handler: h, notifService: notifSvc, logService: logSvc}
}

// RegisterRoutes mounts the module's routes onto /api/v1.
func (m *Module) RegisterRoutes(router fiber.Router) {
	routes.Register(router, m.handler)
}

// NotificationService exposes the service so other modules can call
// Send() to dispatch event-driven notifications without importing
// the notification package's internal packages.
func (m *Module) NotificationService() service.NotificationService {
	return m.notifService
}

// LogService exposes the log service for delivery workers to call
// UpdateStatus() after dispatching emails/whatsapp/sms.
func (m *Module) LogService() service.LogService {
	return m.logService
}
