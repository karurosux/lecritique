package product

import (
	"github.com/labstack/echo/v4"
	"kyooar/internal/product/handlers"
	"kyooar/internal/product/repositories"
	"kyooar/internal/product/services"
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
	publicHandler := do.MustInvoke[*handlers.ProductPublicHandler](m.injector)
	
	middlewareProvider := do.MustInvoke[*sharedMiddleware.MiddlewareProvider](m.injector)
	v1.GET("/public/organization/:id/products", publicHandler.GetOrganizationProducts)
	
	products := v1.Group("/products")
	products.Use(middlewareProvider.AuthMiddleware())
	products.Use(middlewareProvider.TeamAwareMiddleware())
	products.GET("/:id", productHandler.GetByID)
	products.PUT("/:id", productHandler.Update)
	products.DELETE("/:id", productHandler.Delete)
}

func RegisterNewModule(container *do.Injector) {
	do.Provide(container, repositories.NewProductRepository)
	do.Provide(container, services.NewProductService)
	do.Provide(container, handlers.NewProductHandler)
	do.Provide(container, handlers.NewProductPublicHandler)
}
