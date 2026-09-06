package auth

import (
	"github.com/gofiber/fiber/v2"

	"rentos-backend/internal/bootstrap"
	"rentos-backend/internal/modules/auth/handler"
	authpostgres "rentos-backend/internal/modules/auth/repository/postgres"
	"rentos-backend/internal/modules/auth/routes"
	"rentos-backend/internal/modules/auth/service"
)

type Module struct {
	handler *handler.Handler
	svc     service.AuthService
}

func New(c *bootstrap.Container) *Module {
	userRepo := authpostgres.NewUserRepository(
		query("create_user.sql"),
		query("find_user_by_id.sql"),
		query("find_user_by_email.sql"),
		query("find_user_by_username.sql"),
		query("update_password.sql"),
	)
	sessionRepo := authpostgres.NewSessionRepository(
		query("create_session.sql"),
		query("refresh_token.sql"),
		query("revoke_session.sql"),
	)
	resetRepo := authpostgres.NewPasswordResetRepository(
		query("create_password_reset.sql"),
		query("find_password_reset.sql"),
		query("use_password_reset.sql"),
	)

	svc := service.NewAuthService(c.DB, userRepo, sessionRepo, resetRepo, c.JWT)
	h := handler.New(svc, c.Validator)

	return &Module{handler: h, svc: svc}
}

func (m *Module) RegisterPublic(router fiber.Router) {
	routes.RegisterPublic(router, m.handler)
}

func (m *Module) RegisterProtected(router fiber.Router) {
	routes.RegisterProtected(router, m.handler)
}

// UserRepository exposes the user repo so superadmin can reuse it.
func (m *Module) AuthService() service.AuthService { return m.svc }
