package routes

import (
	"github.com/gofiber/fiber/v2"

	"rentos-backend/internal/modules/auth/handler"
)

// RegisterPublic mounts public auth routes (no JWT required).
func RegisterPublic(router fiber.Router, h *handler.Handler) {
	auth := router.Group("/auth")
	auth.Post("/login",          h.Login)
	auth.Post("/refresh",        h.Refresh)
	auth.Post("/logout",         h.Logout)
	auth.Post("/forgot-password", h.ForgotPassword)
	auth.Post("/reset-password",  h.ResetPassword)
}

// RegisterProtected mounts auth routes that require JWT.
func RegisterProtected(router fiber.Router, h *handler.Handler) {
	auth := router.Group("/auth")
	auth.Get("/me",              h.Me)
	auth.Post("/change-password", h.ChangePassword)
}
