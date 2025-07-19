package auth

import (
	"github.com/labstack/echo/v4"
	"lecritique/internal/auth/handlers"
	sharedMiddleware "lecritique/internal/shared/middleware"
	"github.com/samber/do"
)

type Module struct {
	injector *do.Injector
}

func NewModule(i *do.Injector) *Module {
	return &Module{injector: i}
}

func (m *Module) RegisterRoutes(v1 *echo.Group) {
	// Get handlers from injector
	authHandler := do.MustInvoke[*handlers.AuthHandler](m.injector)
	teamMemberHandler := do.MustInvoke[*handlers.TeamMemberHandler](m.injector)
	
	// Get middleware provider
	middlewareProvider := do.MustInvoke[*sharedMiddleware.MiddlewareProvider](m.injector)
	
	// Public auth routes
	auth := v1.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/refresh", authHandler.RefreshToken)
	auth.GET("/verify-email", authHandler.VerifyEmail)
	auth.POST("/resend-verification", authHandler.ResendVerificationEmail)
	auth.POST("/password-reset", authHandler.SendPasswordReset)
	auth.POST("/password-reset/confirm", authHandler.ResetPassword)
	
	// Protected auth routes
	authProtected := v1.Group("/auth")
	authProtected.Use(middlewareProvider.AuthMiddleware())
	// Profile route handled by UpdateProfile
	authProtected.PUT("/profile", authHandler.UpdateProfile)
	authProtected.POST("/deactivate", authHandler.RequestDeactivation)
	authProtected.POST("/deactivate/cancel", authHandler.CancelDeactivation)
	authProtected.POST("/email-change", authHandler.ChangeEmail)
	authProtected.POST("/email-change/confirm", authHandler.ConfirmEmailChange)
	authProtected.POST("/send-verification", authHandler.SendEmailVerification)
	
	// Team member routes (current user's team context)
	team := v1.Group("/team")
	team.Use(middlewareProvider.AuthMiddleware())
	team.GET("/members", teamMemberHandler.ListMembers)
	
	// Team member routes (specific team ID) - kept for compatibility
	teams := v1.Group("/teams")
	teams.Use(middlewareProvider.AuthMiddleware())
	teams.GET("/:teamId/members", teamMemberHandler.ListMembers)
	teams.POST("/:teamId/invitations", teamMemberHandler.InviteMember)
	// Invitation management handled through other endpoints
	teams.POST("/:teamId/invitations/:invitationId/resend", teamMemberHandler.ResendInvitation)
	teams.PUT("/:teamId/members/:memberId", teamMemberHandler.UpdateRole)
	teams.DELETE("/:teamId/members/:memberId", teamMemberHandler.RemoveMember)
	teams.POST("/accept-invitation", teamMemberHandler.AcceptInvitation)
}