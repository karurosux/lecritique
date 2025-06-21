package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/lecritique/api/internal/auth/repositories"
	"github.com/lecritique/api/internal/auth/services"
	"github.com/lecritique/api/internal/shared/config"
	"gorm.io/gorm"
)

func RegisterRoutes(v1 *echo.Group, db *gorm.DB, cfg *config.Config) {
	// Initialize repository
	accountRepo := repositories.NewAccountRepository(db)
	
	// Initialize service
	authService := services.NewAuthService(accountRepo, cfg)
	
	// Initialize handler
	authHandler := NewAuthHandler(authService)
	
	// Auth routes
	auth := v1.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/refresh", authHandler.RefreshToken)
}