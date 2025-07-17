package handlers

import (
	"github.com/labstack/echo/v4"
	authServices "github.com/lecritique/api/internal/auth/services"
	"github.com/lecritique/api/internal/qrcode/repositories"
	"github.com/lecritique/api/internal/qrcode/services"
	"github.com/lecritique/api/internal/shared/middleware"
	restaurantRepos "github.com/lecritique/api/internal/restaurant/repositories"
	"gorm.io/gorm"
)

func RegisterRoutes(protected *echo.Group, db *gorm.DB, authService authServices.AuthService) {
	// Initialize repositories
	qrCodeRepo := repositories.NewQRCodeRepository(db)
	restaurantRepo := restaurantRepos.NewRestaurantRepository(db)
	
	// Initialize service
	qrCodeService := services.NewQRCodeService(qrCodeRepo, restaurantRepo)
	
	// Initialize handler
	qrCodeHandler := NewQRCodeHandler(qrCodeService)
	
	// QR Code routes (protected)
	restaurants := protected.Group("/restaurants")
	restaurants.Use(middleware.JWTAuth(authService))
	restaurants.Use(middleware.TeamAware(db)) // Add team-aware middleware
	restaurants.POST("/:restaurantId/qr-codes", qrCodeHandler.Generate)
	restaurants.GET("/:restaurantId/qr-codes", qrCodeHandler.GetByRestaurant)
	
	qrCodes := protected.Group("/qr-codes")
	qrCodes.Use(middleware.JWTAuth(authService))
	qrCodes.Use(middleware.TeamAware(db)) // Add team-aware middleware
	qrCodes.PATCH("/:id", qrCodeHandler.Update)
	qrCodes.DELETE("/:id", qrCodeHandler.Delete)
}