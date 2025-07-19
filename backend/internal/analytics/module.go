package analytics

import (
	"github.com/labstack/echo/v4"
	"lecritique/internal/analytics/handlers"
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
	analyticsHandler := do.MustInvoke[*handlers.AnalyticsHandler](m.injector)
	
	// Get middleware provider
	middlewareProvider := do.MustInvoke[*sharedMiddleware.MiddlewareProvider](m.injector)
	
	// Analytics routes
	analytics := v1.Group("/analytics")
	analytics.Use(middlewareProvider.AuthMiddleware())
	analytics.Use(middlewareProvider.TeamAwareMiddleware())
	analytics.GET("/restaurants/:restaurantId", analyticsHandler.GetRestaurantAnalytics)
	analytics.GET("/restaurants/:restaurantId/charts", analyticsHandler.GetRestaurantChartData)
	analytics.GET("/dashboard/:restaurantId", analyticsHandler.GetDashboardMetrics)
	analytics.GET("/dishes/:dishId", analyticsHandler.GetDishAnalytics)
	analytics.GET("/dishes/:dishId/insights", analyticsHandler.GetDishInsights)
}