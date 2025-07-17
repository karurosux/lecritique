package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/lecritique/api/internal/auth/repositories"
	"github.com/lecritique/api/internal/auth/services"
	"github.com/lecritique/api/internal/shared/config"
	"github.com/lecritique/api/internal/shared/middleware"
	sharedServices "github.com/lecritique/api/internal/shared/services"
	"gorm.io/gorm"
)

func RegisterRoutes(v1 *echo.Group, db *gorm.DB, cfg *config.Config) services.AuthService {
	// Initialize repositories
	accountRepo := repositories.NewAccountRepository(db)
	tokenRepo := repositories.NewTokenRepository(db)
	teamMemberRepo := repositories.NewTeamMemberRepository(db)
	invitationRepo := repositories.NewTeamInvitationRepository(db)
	
	// Initialize email service
	emailService := sharedServices.NewEmailService(cfg)
	
	// Initialize services
	authService := services.NewAuthService(accountRepo, tokenRepo, emailService, cfg)
	teamMemberService := services.NewTeamMemberServiceV2(teamMemberRepo, invitationRepo, accountRepo, emailService)
	
	// Initialize handler
	authHandler := NewAuthHandler(authService, teamMemberService, cfg, db)
	
	// Auth routes
	auth := v1.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	
	// Public routes
	auth.GET("/verify-email", authHandler.VerifyEmail)
	auth.POST("/resend-verification", authHandler.ResendVerificationEmail)
	auth.POST("/forgot-password", authHandler.SendPasswordReset)
	auth.POST("/reset-password", authHandler.ResetPassword)
	auth.POST("/confirm-email-change", authHandler.ConfirmEmailChange)
	
	// Protected routes - require JWT auth
	authProtected := auth.Group("")
	authProtected.Use(middleware.JWTAuth(authService))
	authProtected.POST("/refresh", authHandler.RefreshToken)
	authProtected.POST("/send-verification", authHandler.SendEmailVerification)
	authProtected.POST("/change-email", authHandler.ChangeEmail)
	authProtected.POST("/deactivate", authHandler.RequestDeactivation)
	authProtected.POST("/cancel-deactivation", authHandler.CancelDeactivation)
	authProtected.PUT("/profile", authHandler.UpdateProfile)
	
	return authService
}