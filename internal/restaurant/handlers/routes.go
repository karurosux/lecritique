package handlers

import (
	"github.com/labstack/echo/v4"
	authServices "github.com/lecritique/api/internal/auth/services"
	"github.com/lecritique/api/internal/restaurant/repositories"
	"github.com/lecritique/api/internal/restaurant/services"
	"github.com/lecritique/api/internal/shared/middleware"
	subscriptionRepos "github.com/lecritique/api/internal/subscription/repositories"
	"gorm.io/gorm"
)

func RegisterRoutes(protected *echo.Group, db *gorm.DB, authService authServices.AuthService) {
	// Initialize repositories
	restaurantRepo := repositories.NewRestaurantRepository(db)
	subscriptionRepo := subscriptionRepos.NewSubscriptionRepository(db)
	
	// Initialize service
	restaurantService := services.NewRestaurantService(restaurantRepo, subscriptionRepo)
	
	// Initialize handler
	restaurantHandler := NewRestaurantHandler(restaurantService)
	
	// Restaurant routes (protected)
	restaurants := protected.Group("/restaurants")
	restaurants.Use(middleware.JWTAuth(authService))
	restaurants.POST("", restaurantHandler.Create)
	restaurants.GET("", restaurantHandler.GetAll)
	restaurants.GET("/:id", restaurantHandler.GetByID)
	restaurants.PUT("/:id", restaurantHandler.Update)
	restaurants.DELETE("/:id", restaurantHandler.Delete)
}