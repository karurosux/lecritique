package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/lecritique/api/internal/auth/repositories"
	"github.com/lecritique/api/internal/auth/services"
	"github.com/lecritique/api/internal/shared/config"
	sharedServices "github.com/lecritique/api/internal/shared/services"
	"gorm.io/gorm"
)

func RegisterRoutes(v1 *echo.Group, db *gorm.DB, cfg *config.Config) services.AuthService {
	// Initialize repositories
	accountRepo := repositories.NewAccountRepository(db)
	tokenRepo := repositories.NewTokenRepository(db)
	
	// Initialize email service
	emailService := sharedServices.NewEmailService(cfg)
	
	// Initialize service
	authService := services.NewAuthService(accountRepo, tokenRepo, emailService, cfg)
	
	// Initialize handler
	authHandler := NewAuthHandler(authService)
	
	// Auth routes
	auth := v1.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/refresh", authHandler.RefreshToken)
	
	// Email verification routes
	auth.POST("/send-verification", authHandler.SendEmailVerification) // Protected
	auth.GET("/verify-email", authHandler.VerifyEmail)                 // Public
	
	// Password reset routes
	auth.POST("/forgot-password", authHandler.SendPasswordReset) // Public
	auth.POST("/reset-password", authHandler.ResetPassword)      // Public
	
	return authService
}