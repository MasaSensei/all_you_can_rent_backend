package routes

import (
	"github.com/gofiber/fiber/v2"

	"rentos/internal/modules/auth/handler"
)

// Register mounts auth and user routes under /api/v1.
// Protected routes will gain AuthMiddleware and RBACMiddleware guards
// once middleware/auth.go is wired in cmd/api/main.go (done at the end
// of Phase 2).
func Register(router fiber.Router, h *handler.Handler) {
	auth := router.Group("/auth")
	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)
	auth.Post("/refresh", h.Refresh)
	auth.Post("/logout", h.Logout)
	auth.Post("/forgot-password", h.ForgotPassword)
	auth.Post("/reset-password", h.ResetPassword)

	users := router.Group("/users")
	users.Get("/", h.ListUsers)
	users.Get("/:id", h.GetUser)
	users.Put("/:id", h.UpdateUser)
	users.Delete("/:id", h.DeleteUser)
}
