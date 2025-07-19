package restaurant

import (
	"github.com/labstack/echo/v4"
	"lecritique/internal/restaurant/handlers"
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
	// Get handler from injector
	restaurantHandler := do.MustInvoke[*handlers.RestaurantHandler](m.injector)
	
	// Get middleware provider
	middlewareProvider := do.MustInvoke[*sharedMiddleware.MiddlewareProvider](m.injector)
	
	// Restaurant routes
	restaurants := v1.Group("/restaurants")
	restaurants.Use(middlewareProvider.AuthMiddleware())
	restaurants.Use(middlewareProvider.TeamAwareMiddleware())
	restaurants.POST("", restaurantHandler.Create)
	restaurants.GET("", restaurantHandler.GetAll)
	restaurants.GET("/:id", restaurantHandler.GetByID)
	restaurants.PUT("/:id", restaurantHandler.Update)
	restaurants.DELETE("/:id", restaurantHandler.Delete)
}