package menu

import (
	"github.com/labstack/echo/v4"
	"kyooar/internal/product/handlers"
	sharedMiddleware "kyooar/internal/shared/middleware"
	"github.com/samber/do"
)

type Module struct {
	injector *do.Injector
}

func NewModule(i *do.Injector) *Module {
	return &Module{injector: i}
}

func (m *Module) RegisterRoutes(v1 *echo.Group) {
	productHandler := do.MustInvoke[*handlers.ProductHandler](m.injector)
	publicHandler := do.MustInvoke[*handlers.MenuPublicHandler](m.injector)
	
	middlewareProvider := do.MustInvoke[*sharedMiddleware.MiddlewareProvider](m.injector)
	v1.GET("/organization/:id/menu", publicHandler.GetOrganizationMenu)
	
	products := v1.Group("/products")
	products.Use(middlewareProvider.AuthMiddleware())
	products.Use(middlewareProvider.TeamAwareMiddleware())
	products.GET("/:id", productHandler.GetByID)
	products.PUT("/:id", productHandler.Update)
	products.DELETE("/:id", productHandler.Delete)
}
