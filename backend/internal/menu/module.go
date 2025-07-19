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
	dishHandler := do.MustInvoke[*handlers.DishHandler](m.injector)
	publicHandler := do.MustInvoke[*handlers.MenuPublicHandler](m.injector)
	
	// Get middleware provider
	middlewareProvider := do.MustInvoke[*sharedMiddleware.MiddlewareProvider](m.injector)
	
	// Public menu routes (no auth required)
	v1.GET("/restaurant/:id/menu", publicHandler.GetRestaurantMenu)
	
	// Menu routes under restaurants
	restaurants := v1.Group("/restaurants")
	restaurants.Use(middlewareProvider.AuthMiddleware())
	restaurants.Use(middlewareProvider.TeamAwareMiddleware())
	restaurants.GET("/:restaurantId/dishes", dishHandler.GetByRestaurant)
	restaurants.POST("/:restaurantId/dishes", dishHandler.Create)
	
	// Direct dish routes
	dishes := v1.Group("/dishes")
	dishes.Use(middlewareProvider.AuthMiddleware())
	dishes.Use(middlewareProvider.TeamAwareMiddleware())
	dishes.GET("/:id", dishHandler.GetByID)
	dishes.PUT("/:id", dishHandler.Update)
	dishes.DELETE("/:id", dishHandler.Delete)
}