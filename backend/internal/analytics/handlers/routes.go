package handlers

import (
	"github.com/labstack/echo/v4"
	analyticsServices "github.com/lecritique/api/internal/analytics/services"
	authServices "github.com/lecritique/api/internal/auth/services"
	feedbackRepos "github.com/lecritique/api/internal/feedback/repositories"
	menuRepos "github.com/lecritique/api/internal/menu/repositories"
	qrcodeRepos "github.com/lecritique/api/internal/qrcode/repositories"
	restaurantRepos "github.com/lecritique/api/internal/restaurant/repositories"
	"github.com/lecritique/api/internal/shared/middleware"
	"gorm.io/gorm"
)

func RegisterRoutes(protected *echo.Group, db *gorm.DB, authService authServices.AuthService) {
	// Initialize repositories
	feedbackRepo := feedbackRepos.NewFeedbackRepository(db)
	dishRepo := menuRepos.NewDishRepository(db)
	qrCodeRepo := qrcodeRepos.NewQRCodeRepository(db)
	restaurantRepo := restaurantRepos.NewRestaurantRepository(db)
	
	// Initialize analytics service
	analyticsService := analyticsServices.NewAnalyticsService(feedbackRepo, dishRepo, qrCodeRepo, restaurantRepo)
	
	// Initialize handler
	analyticsHandler := NewAnalyticsHandler(feedbackRepo, dishRepo, restaurantRepo, analyticsService)
	
	// Analytics routes (protected)
	analytics := protected.Group("/analytics")
	analytics.Use(middleware.JWTAuth(authService))
	analytics.GET("/restaurants/:restaurantId", analyticsHandler.GetRestaurantAnalytics)
	analytics.GET("/dishes/:dishId", analyticsHandler.GetDishAnalytics)
	analytics.GET("/dashboard/:restaurantId", analyticsHandler.GetDashboardMetrics)
	analytics.GET("/dishes/:dishId/insights", analyticsHandler.GetDishInsights)
}