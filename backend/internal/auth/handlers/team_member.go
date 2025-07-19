package handlers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"lecritique/internal/auth/models"
	"lecritique/internal/auth/services"
	"lecritique/internal/shared/errors"
	"lecritique/internal/shared/logger"
	"lecritique/internal/shared/middleware"
	"lecritique/internal/shared/response"
	"lecritique/internal/shared/validator"
	"github.com/samber/do"
	"github.com/sirupsen/logrus"
)

type TeamMemberHandler struct {
	teamMemberService services.TeamMemberServiceV2
	authService       services.AuthService
	validator         *validator.Validator
}

func NewTeamMemberHandler(i *do.Injector) (*TeamMemberHandler, error) {
	return &TeamMemberHandler{
		teamMemberService: do.MustInvoke[services.TeamMemberServiceV2](i),
		authService:       do.MustInvoke[services.AuthService](i),
		validator:         validator.New(),
	}, nil
}

type InviteMemberRequest struct {
	Email string            `json:"email" validate:"required,email"`
	Role  models.MemberRole `json:"role" validate:"required,oneof=ADMIN MANAGER VIEWER"`
}

type UpdateRoleRequest struct {
	Role models.MemberRole `json:"role" validate:"required,oneof=ADMIN MANAGER VIEWER"`
}

type AcceptInviteRequest struct {
	Token string `json:"token" validate:"required"`
}

// ListMembers godoc
// @Summary List team members
// @Description Get all team members for the account
// @Tags team
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]models.TeamMember}
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/team/members [get]
// @Security BearerAuth
func (h *TeamMemberHandler) ListMembers(c echo.Context) error {
	ctx := c.Request().Context()
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		return response.Error(c, err)
	}

	members, err := h.teamMemberService.ListMembers(ctx, accountID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, members)
}

// InviteMember godoc
// @Summary Invite team member
// @Description Invite a new team member to the account
// @Tags team
// @Accept json
// @Produce json
// @Param request body InviteMemberRequest true "Invitation details"
// @Success 201 {object} response.Response{data=models.TeamMember}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/team/members/invite [post]
// @Security BearerAuth
func (h *TeamMemberHandler) InviteMember(c echo.Context) error {
	ctx := c.Request().Context()
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		return response.Error(c, err)
	}

	// Get current user ID (inviter)
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return response.Error(c, errors.ErrUnauthorized)
	}

	// Check user role - only owners and admins can invite
	userRole, ok := c.Get("user_role").(string)
	if !ok || (userRole != string(models.RoleOwner) && userRole != string(models.RoleAdmin)) {
		return response.Error(c, errors.Forbidden("invite team members"))
	}

	var req InviteMemberRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, errors.BadRequest("Invalid invitation data provided"))
	}

	if err := h.validator.Validate(&req); err != nil {
		return response.Error(c, err)
	}

	member, err := h.teamMemberService.InviteMember(ctx, accountID, userID, req.Email, req.Role)
	if err != nil {
		return response.Error(c, err)
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusCreated)
	return c.JSON(http.StatusCreated, response.Response{
		Success: true,
		Data:    member,
	})
}

// UpdateRole godoc
// @Summary Update team member role
// @Description Update the role of a team member
// @Tags team
// @Accept json
// @Produce json
// @Param id path string true "Member ID"
// @Param request body UpdateRoleRequest true "New role"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/team/members/{id}/role [put]
// @Security BearerAuth
func (h *TeamMemberHandler) UpdateRole(c echo.Context) error {
	ctx := c.Request().Context()
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		return response.Error(c, err)
	}

	// Check user role - only owners and admins can update roles
	userRole, ok := c.Get("user_role").(string)
	if !ok || (userRole != string(models.RoleOwner) && userRole != string(models.RoleAdmin)) {
		return response.Error(c, errors.Forbidden("update team member roles"))
	}

	memberID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, errors.ErrInvalidUUID)
	}

	var req UpdateRoleRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, errors.BadRequest("Invalid role update data provided"))
	}

	if err := h.validator.Validate(&req); err != nil {
		return response.Error(c, err)
	}

	if err := h.teamMemberService.UpdateRole(ctx, accountID, memberID, req.Role); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]string{
		"message": "Role updated successfully",
	})
}

// RemoveMember godoc
// @Summary Remove team member
// @Description Remove a team member from the account
// @Tags team
// @Accept json
// @Produce json
// @Param id path string true "Member ID"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/team/members/{id} [delete]
// @Security BearerAuth
func (h *TeamMemberHandler) RemoveMember(c echo.Context) error {
	ctx := c.Request().Context()
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		return response.Error(c, err)
	}

	// Check user role - only owners and admins can remove members
	userRole, ok := c.Get("user_role").(string)
	if !ok || (userRole != string(models.RoleOwner) && userRole != string(models.RoleAdmin)) {
		return response.Error(c, errors.Forbidden("remove team members"))
	}

	memberID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, errors.ErrInvalidUUID)
	}

	if err := h.teamMemberService.RemoveMember(ctx, accountID, memberID); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]string{
		"message": "Team member removed successfully",
	})
}

// ResendInvitation godoc
// @Summary Resend team invitation
// @Description Resend an invitation email to a pending team member
// @Tags team
// @Accept json
// @Produce json
// @Param id path string true "Invitation ID"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/team/members/{id}/resend-invitation [post]
// @Security BearerAuth
func (h *TeamMemberHandler) ResendInvitation(c echo.Context) error {
	ctx := c.Request().Context()
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		return response.Error(c, err)
	}

	// Check user role - only owners and admins can resend invitations
	userRole, ok := c.Get("user_role").(string)
	if !ok || (userRole != string(models.RoleOwner) && userRole != string(models.RoleAdmin)) {
		return response.Error(c, errors.Forbidden("resend invitations"))
	}

	invitationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, errors.ErrInvalidUUID)
	}

	if err := h.teamMemberService.ResendInvitation(ctx, accountID, invitationID); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]string{
		"message": "Invitation resent successfully",
	})
}

// AcceptInvitation godoc
// @Summary Accept team invitation
// @Description Accept a team invitation using the invitation token
// @Tags team
// @Accept json
// @Produce json
// @Param request body AcceptInviteRequest true "Invitation token"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/team/accept-invite [post]
func (h *TeamMemberHandler) AcceptInvitation(c echo.Context) error {
	ctx := c.Request().Context()

	var req AcceptInviteRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, errors.BadRequest("Invalid invitation token data"))
	}

	if err := h.validator.Validate(&req); err != nil {
		return response.Error(c, err)
	}

	// Get the invitation details
	invitation, err := h.teamMemberService.GetInvitationByToken(ctx, req.Token)
	if err != nil {
		return response.Error(c, err)
	}

	// Get the account for the invitation email
	account, err := h.authService.GetAccountByEmail(ctx, invitation.Email)
	if err != nil {
		// Account doesn't exist - mark invitation as email_accepted
		// When they register with this email, they'll be auto-added to the team
		now := time.Now()
		invitation.EmailAcceptedAt = &now
		if err := h.teamMemberService.UpdateInvitation(ctx, invitation); err != nil {
			logger.Warn("Failed to update invitation", logrus.Fields{
				"error": err.Error(),
			})
		}

		return response.Success(c, map[string]interface{}{
			"invitation": invitation,
			"message":    "Please register with this email address to join the team.",
			"status":     "needs_registration",
		})
	}

	// Account exists - accept the invitation
	if err := h.teamMemberService.AcceptInvitation(ctx, req.Token, account.ID); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]string{
		"message": "Invitation accepted successfully. You have been added to the team.",
		"status":  "accepted",
	})
}

