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

func RegisterTeamRoutes(v1 *echo.Group, db *gorm.DB, cfg *config.Config, authService services.AuthService) {
	// Initialize repositories
	teamMemberRepo := repositories.NewTeamMemberRepository(db)
	invitationRepo := repositories.NewTeamInvitationRepository(db)
	accountRepo := repositories.NewAccountRepository(db)
	
	// Initialize email service
	emailService := sharedServices.NewEmailService(cfg)
	
	// Initialize team member service
	teamMemberService := services.NewTeamMemberServiceV2(
		teamMemberRepo,
		invitationRepo,
		accountRepo,
		emailService,
	)
	
	// Initialize handler
	teamHandler := NewTeamMemberHandler(teamMemberService, authService)
	
	// Team routes
	team := v1.Group("/team")
	
	// Public route for accepting invitations
	team.POST("/accept-invite", teamHandler.AcceptInvitation)
	
	// Protected routes - require JWT auth
	teamProtected := team.Group("")
	teamProtected.Use(middleware.JWTAuth(authService))
	teamProtected.Use(middleware.TeamAuthMiddleware(db))
	
	// Team member management
	teamProtected.GET("/members", teamHandler.ListMembers)
	teamProtected.POST("/members/invite", teamHandler.InviteMember)
	teamProtected.POST("/members/:id/resend-invitation", teamHandler.ResendInvitation)
	teamProtected.PUT("/members/:id/role", teamHandler.UpdateRole)
	teamProtected.DELETE("/members/:id", teamHandler.RemoveMember)
}