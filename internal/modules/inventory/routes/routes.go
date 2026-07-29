package routes

import (
	"github.com/gofiber/fiber/v2"

	"rentos-backend/internal/modules/inventory/handler"
)

// Register mounts all inventory routes under /api/v1.
func Register(router fiber.Router, h *handler.Handler) {
	// Categories
	categories := router.Group("/categories")
	categories.Post("/", h.CreateCategory)
	categories.Get("/", h.ListCategories)
	categories.Get("/:id", h.GetCategory)
	categories.Put("/:id", h.UpdateCategory)
	categories.Delete("/:id", h.DeleteCategory)

	// Asset Templates
	templates := router.Group("/asset-templates")
	templates.Post("/", h.CreateAssetTemplate)
	templates.Get("/", h.ListAssetTemplates)
	templates.Get("/:id", h.GetAssetTemplate)
	templates.Put("/:id", h.UpdateAssetTemplate)
	templates.Delete("/:id", h.DeleteAssetTemplate)

	// Assets
	assets := router.Group("/assets")
	assets.Post("/", h.CreateAsset)
	assets.Get("/", h.ListAssets)
	assets.Get("/:id", h.GetAsset)
	assets.Put("/:id", h.UpdateAsset)
	assets.Delete("/:id", h.DeleteAsset)

	// Asset Images
	assets.Post("/:id/images", h.AddAssetImage)
	assets.Delete("/:id/images/:image_id", h.DeleteAssetImage)

	// Asset Documents
	assets.Post("/:id/documents", h.AddAssetDocument)
	assets.Delete("/:id/documents/:doc_id", h.DeleteAssetDocument)

	// Asset Availability
	assets.Post("/:id/availability", h.BlockAvailability)
	assets.Get("/:id/availability", h.ListAvailability)
}
