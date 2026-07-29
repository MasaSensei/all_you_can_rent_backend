package auth

import (
	"github.com/gofiber/fiber/v2"

	"rentos-backend/internal/bootstrap"
	"rentos-backend/internal/modules/auth/handler"
	"rentos-backend/internal/modules/auth/repository/postgres"
	"rentos-backend/internal/modules/auth/routes"
	"rentos-backend/internal/modules/auth/service"
)

// Module holds the auth module's wired handler.
type Module struct {
	handler *handler.Handler
}

// New builds the auth module: repositories → services → handler.
func New(c *bootstrap.Container) *Module {
	userRepo := postgres.NewUserRepository(
		query("create_user.sql"),
		query("find_user_by_id.sql"),
		query("find_user_by_email.sql"),
		query("list_users.sql"),
		query("update_user.sql"),
		query("update_last_login.sql"),
		query("update_password.sql"),
		query("delete_user.sql"),
	)
	sessionRepo := postgres.NewSessionRepository(
		query("create_session.sql"),
		query("find_session_by_refresh_token.sql"),
		query("revoke_session.sql"),
		query("revoke_all_user_sessions.sql"),
	)
	passwordResetRepo := postgres.NewPasswordResetRepository(
		query("create_password_reset.sql"),
		query("find_password_reset_by_token.sql"),
		query("consume_password_reset.sql"),
	)

	authSvc := service.NewAuthService(c.DB, userRepo, sessionRepo, c.JWT)
	userSvc := service.NewUserService(c.DB, userRepo)
	passwordSvc := service.NewPasswordService(c.DB, userRepo, passwordResetRepo)

	h := handler.New(authSvc, userSvc, passwordSvc, c.Validator)
	return &Module{handler: h}
}

// RegisterRoutes mounts the module's routes onto /api/v1.
func (m *Module) RegisterRoutes(router fiber.Router) {
	routes.Register(router, m.handler)
}
