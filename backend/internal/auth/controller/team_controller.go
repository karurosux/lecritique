package authcontroller

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	authinterface "kyooar/internal/auth/interface"
	authmodel "kyooar/internal/auth/model"
	"kyooar/internal/auth/models"
	"kyooar/internal/shared/errors"
	"kyooar/internal/shared/middleware"
	"kyooar/internal/shared/response"
	"kyooar/internal/shared/validator"
)

type TeamController struct {
	teamMemberService authinterface.TeamMemberService
	validator         *validator.Validator
}

func NewTeamController(
	teamMemberService authinterface.TeamMemberService,
	validator *validator.Validator,
) *TeamController {
	return &TeamController{
		teamMemberService: teamMemberService,
		validator:         validator,
	}
}

// @Summary List team members
// @Description Get list of team members for the authenticated account
// @Tags team
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=authmodel.MemberListResponse}
// @Failure 401 {object} response.Response
// @Router /api/v1/team/members [get]
func (c *TeamController) ListMembers(ctx echo.Context) error {
	accountID, err := middleware.GetAccountID(ctx)
	if err != nil {
		return response.Error(ctx, err)
	}

	members, err := c.teamMemberService.ListMembers(ctx.Request().Context(), accountID)
	if err != nil {
		return response.Error(ctx, err)
	}

	return response.Success(ctx, authmodel.MemberListResponse{
		Members: members,
	})
}

// @Summary Invite team member
// @Description Invite a new member to the team
// @Tags team
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body authmodel.InviteMemberRequest true "Member invitation details"
// @Success 201 {object} response.Response{data=authmodel.InvitationResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/team/members/invite [post]
func (c *TeamController) InviteMember(ctx echo.Context) error {
	accountID, err := middleware.GetAccountID(ctx)
	if err != nil {
		return response.Error(ctx, err)
	}

	memberID, err := middleware.GetMemberID(ctx)
	if err != nil {
		return response.Error(ctx, err)
	}

	var req authmodel.InviteMemberRequest
	if err := ctx.Bind(&req); err != nil {
		return response.Error(ctx, errors.ErrBadRequest)
	}

	if err := c.validator.Validate(req); err != nil {
		return response.Error(ctx, errors.NewWithDetails("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, c.validator.FormatErrors(err)))
	}

	role := models.MemberRole(req.Role)
	invitation, err := c.teamMemberService.InviteMember(ctx.Request().Context(), accountID, req.Email, role, memberID)
	if err != nil {
		return response.Error(ctx, err)
	}

	return response.Success(ctx, authmodel.InvitationResponse{
		Invitation: invitation,
		Message:    "Team member invitation sent successfully",
	})
}

// @Summary Resend invitation
// @Description Resend invitation to a team member
// @Tags team
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Invitation ID"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/team/members/{id}/resend-invitation [post]
func (c *TeamController) ResendInvitation(ctx echo.Context) error {
	accountID, err := middleware.GetAccountID(ctx)
	if err != nil {
		return response.Error(ctx, err)
	}

	invitationID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return response.Error(ctx, errors.ErrBadRequest)
	}

	if err := c.teamMemberService.ResendInvitation(ctx.Request().Context(), accountID, invitationID); err != nil {
		return response.Error(ctx, err)
	}

	return response.Success(ctx, map[string]string{
		"message": "Invitation resent successfully",
	})
}

// @Summary Update member role
// @Description Update the role of a team member
// @Tags team
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Member ID"
// @Param request body authmodel.UpdateRoleRequest true "Role update details"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/team/members/{id}/role [put]
func (c *TeamController) UpdateRole(ctx echo.Context) error {
	memberID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return response.Error(ctx, errors.ErrBadRequest)
	}

	var req authmodel.UpdateRoleRequest
	if err := ctx.Bind(&req); err != nil {
		return response.Error(ctx, errors.ErrBadRequest)
	}

	if err := c.validator.Validate(req); err != nil {
		return response.Error(ctx, errors.NewWithDetails("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, c.validator.FormatErrors(err)))
	}

	role := models.MemberRole(req.Role)
	if err := c.teamMemberService.UpdateRoleByID(ctx.Request().Context(), memberID, role); err != nil {
		return response.Error(ctx, err)
	}

	return response.Success(ctx, map[string]string{
		"message": "Member role updated successfully",
	})
}

// @Summary Remove team member
// @Description Remove a member from the team
// @Tags team
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Member ID"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/team/members/{id} [delete]
func (c *TeamController) RemoveMember(ctx echo.Context) error {
	memberID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return response.Error(ctx, errors.ErrBadRequest)
	}

	if err := c.teamMemberService.RemoveMemberByID(ctx.Request().Context(), memberID); err != nil {
		return response.Error(ctx, err)
	}

	return response.Success(ctx, map[string]string{
		"message": "Team member removed successfully",
	})
}

// @Summary Accept team invitation
// @Description Accept a team invitation using the invitation token
// @Tags team
// @Accept json
// @Produce json
// @Param request body authmodel.AcceptInvitationRequest true "Invitation token"
// @Success 200 {object} response.Response{data=interface{}}
// @Failure 400 {object} response.Response
// @Router /api/v1/team/accept-invitation [post]
func (c *TeamController) AcceptInvitation(ctx echo.Context) error {
	var req authmodel.AcceptInvitationRequest
	if err := ctx.Bind(&req); err != nil {
		return response.Error(ctx, errors.ErrBadRequest)
	}

	if err := c.validator.Validate(req); err != nil {
		return response.Error(ctx, errors.NewWithDetails("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, c.validator.FormatErrors(err)))
	}

	member, err := c.teamMemberService.AcceptInvitation(ctx.Request().Context(), req.Token)
	if err != nil {
		return response.Error(ctx, err)
	}

	return response.Success(ctx, map[string]interface{}{
		"member":  member,
		"message": "Team invitation accepted successfully",
	})
}

// Handler methods for /teams/:teamId/... endpoints
func (c *TeamController) ListMembersWithTeamID(ctx echo.Context) error {
	// For now, just extract teamId and use the same logic
	// In the future, this can be used for multi-team support
	teamID, err := uuid.Parse(ctx.Param("teamId"))
	if err != nil {
		return response.Error(ctx, errors.ErrBadRequest)
	}

	members, err := c.teamMemberService.ListMembers(ctx.Request().Context(), teamID)
	if err != nil {
		return response.Error(ctx, err)
	}

	return response.Success(ctx, authmodel.MemberListResponse{
		Members: members,
	})
}

func (c *TeamController) InviteMemberWithTeamID(ctx echo.Context) error {
	teamID, err := uuid.Parse(ctx.Param("teamId"))
	if err != nil {
		return response.Error(ctx, errors.ErrBadRequest)
	}

	memberID, err := middleware.GetMemberID(ctx)
	if err != nil {
		return response.Error(ctx, err)
	}

	var req authmodel.InviteMemberRequest
	if err := ctx.Bind(&req); err != nil {
		return response.Error(ctx, errors.ErrBadRequest)
	}

	if err := c.validator.Validate(req); err != nil {
		return response.Error(ctx, errors.NewWithDetails("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, c.validator.FormatErrors(err)))
	}

	role := models.MemberRole(req.Role)
	invitation, err := c.teamMemberService.InviteMember(ctx.Request().Context(), teamID, req.Email, role, memberID)
	if err != nil {
		return response.Error(ctx, err)
	}

	return response.Success(ctx, authmodel.InvitationResponse{
		Invitation: invitation,
		Message:    "Team member invitation sent successfully",
	})
}

func (c *TeamController) ResendInvitationWithTeamID(ctx echo.Context) error {
	teamID, err := uuid.Parse(ctx.Param("teamId"))
	if err != nil {
		return response.Error(ctx, errors.ErrBadRequest)
	}

	invitationID, err := uuid.Parse(ctx.Param("invitationId"))
	if err != nil {
		return response.Error(ctx, errors.ErrBadRequest)
	}

	if err := c.teamMemberService.ResendInvitation(ctx.Request().Context(), teamID, invitationID); err != nil {
		return response.Error(ctx, err)
	}

	return response.Success(ctx, map[string]string{
		"message": "Invitation resent successfully",
	})
}

func (c *TeamController) UpdateRoleWithTeamID(ctx echo.Context) error {
	teamID, err := uuid.Parse(ctx.Param("teamId"))
	if err != nil {
		return response.Error(ctx, errors.ErrBadRequest)
	}

	memberID, err := uuid.Parse(ctx.Param("memberId"))
	if err != nil {
		return response.Error(ctx, errors.ErrBadRequest)
	}

	var req authmodel.UpdateRoleRequest
	if err := ctx.Bind(&req); err != nil {
		return response.Error(ctx, errors.ErrBadRequest)
	}

	if err := c.validator.Validate(req); err != nil {
		return response.Error(ctx, errors.NewWithDetails("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, c.validator.FormatErrors(err)))
	}

	role := models.MemberRole(req.Role)
	if err := c.teamMemberService.UpdateRole(ctx.Request().Context(), teamID, memberID, role); err != nil {
		return response.Error(ctx, err)
	}

	return response.Success(ctx, map[string]string{
		"message": "Member role updated successfully",
	})
}

func (c *TeamController) RemoveMemberWithTeamID(ctx echo.Context) error {
	teamID, err := uuid.Parse(ctx.Param("teamId"))
	if err != nil {
		return response.Error(ctx, errors.ErrBadRequest)
	}

	memberID, err := uuid.Parse(ctx.Param("memberId"))
	if err != nil {
		return response.Error(ctx, errors.ErrBadRequest)
	}

	if err := c.teamMemberService.RemoveMember(ctx.Request().Context(), teamID, memberID); err != nil {
		return response.Error(ctx, err)
	}

	return response.Success(ctx, map[string]string{
		"message": "Team member removed successfully",
	})
}

func (c *TeamController) RegisterRoutes(v1 *echo.Group, authMiddleware echo.MiddlewareFunc) {
	// /team endpoints
	team := v1.Group("/team")
	team.Use(authMiddleware)
	team.GET("/members", c.ListMembers)
	team.POST("/members/invite", c.InviteMember)
	team.POST("/members/:id/resend-invitation", c.ResendInvitation)
	team.PUT("/members/:id/role", c.UpdateRole)
	team.DELETE("/members/:id", c.RemoveMember)

	// /teams endpoints with teamId parameter
	teams := v1.Group("/teams")
	teams.Use(authMiddleware)
	teams.GET("/:teamId/members", c.ListMembersWithTeamID)
	teams.POST("/:teamId/invitations", c.InviteMemberWithTeamID)
	teams.POST("/:teamId/invitations/:invitationId/resend", c.ResendInvitationWithTeamID)
	teams.PUT("/:teamId/members/:memberId", c.UpdateRoleWithTeamID)
	teams.DELETE("/:teamId/members/:memberId", c.RemoveMemberWithTeamID)
	teams.POST("/accept-invitation", c.AcceptInvitation)

	// Duplicate endpoint for backward compatibility
	v1.POST("/team/accept-invite", c.AcceptInvitation)
}

