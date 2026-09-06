package core

import (
	"github.com/gofiber/fiber/v2"

	"rentos-backend/internal/bootstrap"
	"rentos-backend/internal/modules/core/handler"
	"rentos-backend/internal/modules/core/routes"
)

type Module struct{ handler *handler.Handler }

func New(_ *bootstrap.Container) *Module { return &Module{handler: handler.New()} }
func (m *Module) RegisterRoutes(r fiber.Router) { routes.Register(r, m.handler) }
func (m *Module) RegisterPublic(r fiber.Router) { routes.Register(r, m.handler) }
