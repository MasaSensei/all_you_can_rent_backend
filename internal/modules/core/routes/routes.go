package routes

import (
	"github.com/gofiber/fiber/v2"
	"rentos-backend/internal/modules/core/handler"
)

func Register(r fiber.Router, h *handler.Handler) {
	r.Get("/ping", h.Ping)
	r.Get("/version", h.Version)
}
