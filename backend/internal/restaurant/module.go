package restaurant

import (
	"github.com/labstack/echo/v4"
	"lecritique/internal/restaurant/handlers"
	sharedMiddleware "lecritique/internal/shared/middleware"
	subscriptionMiddleware "lecritique/internal/subscription/middleware"
	subscriptionModels "lecritique/internal/subscription/models"
	"github.com/samber/do"
)

type Module struct {
	injector *do.Injector
}

func NewModule(i *do.Injector) *Module {
	return &Module{injector: i}
}

func (m *Module) RegisterRoutes(v1 *echo.Group) {
	// Get handlers and middleware from injector
	restaurantHandler := do.MustInvoke[*handlers.RestaurantHandler](m.injector)
	middlewareProvider := do.MustInvoke[*sharedMiddleware.MiddlewareProvider](m.injector)
	subscriptionMW := do.MustInvoke[*subscriptionMiddleware.SubscriptionMiddleware](m.injector)
	
	// Restaurant routes
	restaurants := v1.Group("/restaurants")
	restaurants.Use(middlewareProvider.AuthMiddleware())
	restaurants.Use(middlewareProvider.TeamAwareMiddleware())
	
	// Apply usage tracking middleware only to restaurant creation
	restaurants.POST("", restaurantHandler.Create,
		subscriptionMW.CheckResourceLimit(subscriptionModels.ResourceTypeRestaurant),
		subscriptionMW.TrackUsageAfterSuccess(),
	)
	
	restaurants.GET("", restaurantHandler.GetAll)
	restaurants.GET("/:id", restaurantHandler.GetByID)
	restaurants.PUT("/:id", restaurantHandler.Update)
	restaurants.DELETE("/:id", restaurantHandler.Delete)
}