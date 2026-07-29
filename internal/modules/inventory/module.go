package inventory

import (
	"github.com/gofiber/fiber/v2"

	"rentos-backend/internal/bootstrap"
	"rentos-backend/internal/modules/inventory/handler"
	"rentos-backend/internal/modules/inventory/repository/postgres"
	"rentos-backend/internal/modules/inventory/routes"
	"rentos-backend/internal/modules/inventory/service"
)

// Module holds the inventory module's wired handler and services.
type Module struct {
	handler   *handler.Handler
	assetsSvc service.AssetService
}

// New builds the inventory module: repositories → services → handler.
func New(c *bootstrap.Container) *Module {
	categoryRepo := postgres.NewCategoryRepository(
		query("create_category.sql"),
		query("find_category_by_id.sql"),
		query("list_categories.sql"),
		query("update_category.sql"),
		query("delete_category.sql"),
	)
	templateRepo := postgres.NewAssetTemplateRepository(
		query("create_asset_template.sql"),
		query("find_asset_template_by_id.sql"),
		query("list_asset_templates.sql"),
		query("update_asset_template.sql"),
		query("delete_asset_template.sql"),
	)
	fieldRepo := postgres.NewTemplateFieldRepository(
		query("create_template_field.sql"),
		query("list_template_fields.sql"),
	)
	assetRepo := postgres.NewAssetRepository(
		query("create_asset.sql"),
		query("find_asset_by_id.sql"),
		query("list_assets.sql"),
		query("update_asset.sql"),
		query("delete_asset.sql"),
	)
	valueRepo := postgres.NewAssetValueRepository(
		query("upsert_asset_value.sql"),
		query("list_asset_values.sql"),
	)
	imageRepo := postgres.NewAssetImageRepository(
		query("create_asset_image.sql"),
		query("list_asset_images.sql"),
		query("delete_asset_image.sql"),
	)
	docRepo := postgres.NewAssetDocumentRepository(
		query("create_asset_document.sql"),
		query("list_asset_documents.sql"),
		query("delete_asset_document.sql"),
	)
	availRepo := postgres.NewAssetAvailabilityRepository(
		query("create_asset_availability.sql"),
		query("find_availability_conflicts.sql"),
		query("list_asset_availability.sql"),
	)

	categorySvc := service.NewCategoryService(c.DB, categoryRepo)
	templateSvc := service.NewAssetTemplateService(c.DB, templateRepo, fieldRepo)
	assetSvc := service.NewAssetService(c.DB, assetRepo, valueRepo, imageRepo, docRepo, availRepo, fieldRepo)

	h := handler.New(categorySvc, templateSvc, assetSvc, c.Validator)
	return &Module{handler: h, assetsSvc: assetSvc}
}

// RegisterRoutes mounts the module's routes onto /api/v1.
func (m *Module) RegisterRoutes(router fiber.Router) {
	routes.Register(router, m.handler)
}

// AssetService exposes AssetService so the booking module can call
// CheckAvailability via the service interface without importing the
// inventory postgres package.
func (m *Module) AssetService() service.AssetService {
	return m.assetsSvc
}
