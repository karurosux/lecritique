package menu

import (
	"github.com/labstack/echo/v4"
	"lecritique/internal/menu/handlers"
	sharedMiddleware "lecritique/internal/shared/middleware"
	"github.com/samber/do"
)

type Module struct {
	injector *do.Injector
}

func NewModule(i *do.Injector) *Module {
	return &Module{injector: i}
}

func (m *Module) RegisterRoutes(v1 *echo.Group) {
	// Get handlers from injector
	productHandler := do.MustInvoke[*handlers.ProductHandler](m.injector)
	publicHandler := do.MustInvoke[*handlers.MenuPublicHandler](m.injector)
	
	// Get middleware provider
	middlewareProvider := do.MustInvoke[*sharedMiddleware.MiddlewareProvider](m.injector)
	
	// Public menu routes (no auth required)
	v1.GET("/organization/:id/menu", publicHandler.GetOrganizationMenu)
	
	// Menu routes under organizations (moved to organization module)
	
	// Direct product routes
	products := v1.Group("/products")
	products.Use(middlewareProvider.AuthMiddleware())
	products.Use(middlewareProvider.TeamAwareMiddleware())
	products.GET("/:id", productHandler.GetByID)
	products.PUT("/:id", productHandler.Update)
	products.DELETE("/:id", productHandler.Delete)
}
