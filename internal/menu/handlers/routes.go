package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/lecritique/api/internal/menu/repositories"
	"github.com/lecritique/api/internal/menu/services"
	"github.com/lecritique/api/internal/shared/middleware"
	restaurantRepos "github.com/lecritique/api/internal/restaurant/repositories"
	"gorm.io/gorm"
)

func RegisterRoutes(protected *echo.Group, db *gorm.DB, authService interface{}) {
	// Initialize repositories
	dishRepo := repositories.NewDishRepository(db)
	restaurantRepo := restaurantRepos.NewRestaurantRepository(db)
	
	// Initialize service
	dishService := services.NewDishService(dishRepo, restaurantRepo)
	
	// Initialize handler
	dishHandler := NewDishHandler(dishService)
	
	// Dish routes (protected)
	dishes := protected.Group("/dishes")
	dishes.Use(middleware.JWTAuth(authService))
	dishes.POST("", dishHandler.Create)
	dishes.GET("/:id", dishHandler.GetByID)
	dishes.PUT("/:id", dishHandler.Update)
	dishes.DELETE("/:id", dishHandler.Delete)
	
	// Restaurant-specific dish routes
	restaurants := protected.Group("/restaurants")
	restaurants.Use(middleware.JWTAuth(authService))
	restaurants.GET("/:restaurantId/dishes", dishHandler.GetByRestaurant)
}