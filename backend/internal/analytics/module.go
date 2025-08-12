package analytics

import (
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"gorm.io/gorm"

	analyticscontroller "kyooar/internal/analytics/controller"
	analyticsinterface "kyooar/internal/analytics/interface"
	gormrepo "kyooar/internal/analytics/repository/gorm"
	analyticsservice "kyooar/internal/analytics/services"
	feedbackRepos "kyooar/internal/feedback/repositories"
	feedbackServices "kyooar/internal/feedback/services"
	menuRepos "kyooar/internal/menu/repositories"
	organizationRepos "kyooar/internal/organization/repositories"
	organizationServices "kyooar/internal/organization/services"
	qrcodeRepos "kyooar/internal/qrcode/repositories"
	sharedMiddleware "kyooar/internal/shared/middleware"
)

func ProvideAnalyticsRepository(i *do.Injector) (analyticsinterface.AnalyticsRepository, error) {
	db := do.MustInvoke[*gorm.DB](i)
	return gormrepo.NewAnalyticsRepository(db), nil
}

func ProvideTimeSeriesRepository(i *do.Injector) (analyticsinterface.TimeSeriesRepository, error) {
	db := do.MustInvoke[*gorm.DB](i)
	return gormrepo.NewTimeSeriesRepository(db), nil
}

func ProvideAnalyticsService(i *do.Injector) (analyticsinterface.AnalyticsService, error) {
	analyticsRepo := do.MustInvoke[analyticsinterface.AnalyticsRepository](i)
	feedbackRepo := do.MustInvoke[feedbackRepos.FeedbackRepository](i)
	productRepo := do.MustInvoke[menuRepos.ProductRepository](i)
	qrCodeRepo := do.MustInvoke[qrcodeRepos.QRCodeRepository](i)
	organizationRepo := do.MustInvoke[organizationRepos.OrganizationRepository](i)

	return analyticsservice.NewAnalyticsService(
		analyticsRepo,
		feedbackRepo,
		productRepo,
		qrCodeRepo,
		organizationRepo,
	), nil
}

func ProvideTimeSeriesService(i *do.Injector) (analyticsinterface.TimeSeriesService, error) {
	timeSeriesRepo := do.MustInvoke[analyticsinterface.TimeSeriesRepository](i)
	feedbackService := do.MustInvoke[feedbackServices.FeedbackService](i)
	organizationService := do.MustInvoke[organizationServices.OrganizationService](i)
	analyticsService := do.MustInvoke[analyticsinterface.AnalyticsService](i)
	questionService := do.MustInvoke[feedbackServices.QuestionService](i)

	return analyticsservice.NewTimeSeriesService(
		timeSeriesRepo,
		feedbackService,
		organizationService,
		analyticsService,
		questionService,
	), nil
}

func ProvideAnalyticsController(i *do.Injector) (*analyticscontroller.AnalyticsController, error) {
	feedbackRepo := do.MustInvoke[feedbackRepos.FeedbackRepository](i)
	productRepo := do.MustInvoke[menuRepos.ProductRepository](i)
	organizationRepo := do.MustInvoke[organizationRepos.OrganizationRepository](i)
	analyticsService := do.MustInvoke[analyticsinterface.AnalyticsService](i)

	return analyticscontroller.NewAnalyticsController(
		feedbackRepo,
		productRepo,
		organizationRepo,
		analyticsService,
	), nil
}

func ProvideTimeSeriesController(i *do.Injector) (*analyticscontroller.TimeSeriesController, error) {
	timeSeriesService := do.MustInvoke[analyticsinterface.TimeSeriesService](i)
	organizationRepo := do.MustInvoke[organizationRepos.OrganizationRepository](i)

	return analyticscontroller.NewTimeSeriesController(
		timeSeriesService,
		organizationRepo,
	), nil
}

type Module struct {
	injector *do.Injector
}

func NewModule(i *do.Injector) *Module {
	return &Module{injector: i}
}

func (m *Module) RegisterRoutes(v1 *echo.Group) {
	analyticsController := do.MustInvoke[*analyticscontroller.AnalyticsController](m.injector)
	timeSeriesController := do.MustInvoke[*analyticscontroller.TimeSeriesController](m.injector)
	
	middlewareProvider := do.MustInvoke[*sharedMiddleware.MiddlewareProvider](m.injector)
	analytics := v1.Group("/analytics")
	analytics.Use(middlewareProvider.AuthMiddleware())
	analytics.Use(middlewareProvider.TeamAwareMiddleware())
	analytics.GET("/organizations/:organizationId", analyticsController.GetOrganizationAnalytics)
	analytics.GET("/organizations/:organizationId/charts", analyticsController.GetOrganizationChartData)
	analytics.GET("/dashboard/:organizationId", analyticsController.GetDashboardMetrics)
	analytics.GET("/products/:productId", analyticsController.GetProductAnalytics)
	analytics.GET("/products/:productId/insights", analyticsController.GetProductInsights)
	analytics.GET("/organizations/:organizationId/time-series", timeSeriesController.GetTimeSeries)
	analytics.POST("/organizations/:organizationId/compare", timeSeriesController.CompareTimePeriods)
	analytics.POST("/organizations/:organizationId/collect-metrics", timeSeriesController.CollectMetrics)
}

func RegisterModule(container *do.Injector) error {
	do.Provide(container, ProvideAnalyticsRepository)
	do.Provide(container, ProvideTimeSeriesRepository)
	do.Provide(container, ProvideAnalyticsService)
	do.Provide(container, ProvideTimeSeriesService)
	do.Provide(container, ProvideAnalyticsController)
	do.Provide(container, ProvideTimeSeriesController)

	return nil
}