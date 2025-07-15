package handlers

import (
	"github.com/labstack/echo/v4"
	authServices "github.com/lecritique/api/internal/auth/services"
	restaurantRepos "github.com/lecritique/api/internal/restaurant/repositories"
	"github.com/lecritique/api/internal/shared/middleware"
	"github.com/lecritique/api/internal/subscription/repositories"
	"github.com/lecritique/api/internal/subscription/services"
	"gorm.io/gorm"
)

func RegisterRoutes(v1 *echo.Group, db *gorm.DB, authService authServices.AuthService) {
	// Initialize repositories
	subscriptionRepo := repositories.NewSubscriptionRepository(db)
	planRepo := repositories.NewSubscriptionPlanRepository(db)
	restaurantRepo := restaurantRepos.NewRestaurantRepository(db)
	
	// Initialize service
	subscriptionService := services.NewSubscriptionService(subscriptionRepo, planRepo, restaurantRepo)
	
	// Initialize handler
	subscriptionHandler := NewSubscriptionHandler(subscriptionService)
	
	// Public routes (no authentication required)
	v1.GET("/plans", subscriptionHandler.GetAvailablePlans)
	
	// Protected routes (authentication required)
	user := v1.Group("/user")
	user.Use(middleware.JWTAuth(authService))
	
	// User subscription routes
	user.GET("/subscription", subscriptionHandler.GetUserSubscription)
	user.POST("/subscription", subscriptionHandler.CreateSubscription)
	user.DELETE("/subscription", subscriptionHandler.CancelSubscription)
	
	// Permission checking routes
	user.GET("/can-create-restaurant", subscriptionHandler.CanUserCreateRestaurant)
}