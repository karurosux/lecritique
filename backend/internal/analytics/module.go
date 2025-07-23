package analytics

import (
	"github.com/labstack/echo/v4"
	"kyooar/internal/analytics/handlers"
	"kyooar/internal/analytics/services"
	"kyooar/internal/analytics/repositories"
	sharedMiddleware "kyooar/internal/shared/middleware"
	"github.com/samber/do"
)

type Module struct {
	injector *do.Injector
}

func NewModule(i *do.Injector) *Module {
	// Register time series services
	do.Provide(i, repositories.NewTimeSeriesRepository)
	do.Provide(i, services.NewTimeSeriesService)
	do.Provide(i, handlers.NewTimeSeriesHandler)
	
	return &Module{injector: i}
}

func (m *Module) RegisterRoutes(v1 *echo.Group) {
	// Get handlers from injector
	analyticsHandler := do.MustInvoke[*handlers.AnalyticsHandler](m.injector)
	timeSeriesHandler := do.MustInvoke[*handlers.TimeSeriesHandler](m.injector)
	
	// Get middleware provider
	middlewareProvider := do.MustInvoke[*sharedMiddleware.MiddlewareProvider](m.injector)
	
	// Analytics routes
	analytics := v1.Group("/analytics")
	analytics.Use(middlewareProvider.AuthMiddleware())
	analytics.Use(middlewareProvider.TeamAwareMiddleware())
	
	// Existing analytics endpoints
	analytics.GET("/organizations/:organizationId", analyticsHandler.GetOrganizationAnalytics)
	analytics.GET("/organizations/:organizationId/charts", analyticsHandler.GetOrganizationChartData)
	analytics.GET("/dashboard/:organizationId", analyticsHandler.GetDashboardMetrics)
	analytics.GET("/products/:productId", analyticsHandler.GetProductAnalytics)
	analytics.GET("/products/:productId/insights", analyticsHandler.GetProductInsights)
	
	// New time series endpoints
	analytics.GET("/organizations/:organizationId/time-series", timeSeriesHandler.GetTimeSeries)
	analytics.POST("/organizations/:organizationId/compare", timeSeriesHandler.CompareTimePeriods)
	analytics.POST("/organizations/:organizationId/collect-metrics", timeSeriesHandler.CollectMetrics)
}
