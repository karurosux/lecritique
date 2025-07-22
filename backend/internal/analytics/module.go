package analytics

import (
	"github.com/labstack/echo/v4"
	"kyooar/internal/analytics/handlers"
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
	// Get handler from injector
	analyticsHandler := do.MustInvoke[*handlers.AnalyticsHandler](m.injector)
	
	// Get middleware provider
	middlewareProvider := do.MustInvoke[*sharedMiddleware.MiddlewareProvider](m.injector)
	
	// Analytics routes
	analytics := v1.Group("/analytics")
	analytics.Use(middlewareProvider.AuthMiddleware())
	analytics.Use(middlewareProvider.TeamAwareMiddleware())
	analytics.GET("/organizations/:organizationId", analyticsHandler.GetOrganizationAnalytics)
	analytics.GET("/organizations/:organizationId/charts", analyticsHandler.GetOrganizationChartData)
	analytics.GET("/dashboard/:organizationId", analyticsHandler.GetDashboardMetrics)
	analytics.GET("/products/:productId", analyticsHandler.GetProductAnalytics)
	analytics.GET("/products/:productId/insights", analyticsHandler.GetProductInsights)
}
