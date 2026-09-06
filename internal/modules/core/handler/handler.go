package handler

import (
	"runtime"
	"github.com/gofiber/fiber/v2"
	"rentos-backend/pkg/response"
)

type Handler struct{}
func New() *Handler { return &Handler{} }

func (h *Handler) Version(c *fiber.Ctx) error {
	return response.Success(c, fiber.Map{"version": "1.0.0", "go": runtime.Version(), "api": "RentOS"})
}
func (h *Handler) Ping(c *fiber.Ctx) error { return c.SendString("pong") }
