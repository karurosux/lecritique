package handlers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"kyooar/internal/auth/models"
	"kyooar/internal/auth/services"
	"kyooar/internal/shared/errors"
	"kyooar/internal/shared/logger"
	"kyooar/internal/shared/middleware"
	"kyooar/internal/shared/response"
	"kyooar/internal/shared/validator"
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

func (h *TeamMemberHandler) InviteMember(c echo.Context) error {
	ctx := c.Request().Context()
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		return response.Error(c, err)
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, err)
	}

	userRole, err := middleware.GetRole(c)
	if err != nil {
		return response.Error(c, err)
	}
	if userRole != models.RoleOwner && userRole != models.RoleAdmin {
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

func (h *TeamMemberHandler) UpdateRole(c echo.Context) error {
	ctx := c.Request().Context()
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		return response.Error(c, err)
	}

	userRole, err := middleware.GetRole(c)
	if err != nil {
		return response.Error(c, err)
	}
	if userRole != models.RoleOwner && userRole != models.RoleAdmin {
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

func (h *TeamMemberHandler) RemoveMember(c echo.Context) error {
	ctx := c.Request().Context()
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		return response.Error(c, err)
	}

	userRole, err := middleware.GetRole(c)
	if err != nil {
		return response.Error(c, err)
	}
	if userRole != models.RoleOwner && userRole != models.RoleAdmin {
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

func (h *TeamMemberHandler) ResendInvitation(c echo.Context) error {
	ctx := c.Request().Context()
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		return response.Error(c, err)
	}

	userRole, err := middleware.GetRole(c)
	if err != nil {
		return response.Error(c, err)
	}
	if userRole != models.RoleOwner && userRole != models.RoleAdmin {
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

func (h *TeamMemberHandler) AcceptInvitation(c echo.Context) error {
	ctx := c.Request().Context()

	var req AcceptInviteRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, errors.BadRequest("Invalid invitation token data"))
	}

	if err := h.validator.Validate(&req); err != nil {
		return response.Error(c, err)
	}

	invitation, err := h.teamMemberService.GetInvitationByToken(ctx, req.Token)
	if err != nil {
		return response.Error(c, err)
	}

	account, err := h.authService.GetAccountByEmail(ctx, invitation.Email)
	if err != nil {
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

	if err := h.teamMemberService.AcceptInvitation(ctx, req.Token, account.ID); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]string{
		"message": "Invitation accepted successfully. You have been added to the team.",
		"status":  "accepted",
	})
}

