package middleware

import (
	"log"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	authModels "kyooar/internal/auth/models"
	"kyooar/internal/auth/services"
)

func TeamAware(teamMemberService services.TeamMemberServiceV2) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			accountID, err := GetAccountID(c)
			if err != nil {
				return next(c)
			}

			c.Set("personal_account_id", accountID)

			teamMember, err := teamMemberService.GetMemberByMemberID(ctx, accountID)
			if err == nil && teamMember != nil {
				orgAccountID := teamMember.AccountID
				log.Printf("TeamAware: Account %s is a team member of org %s", accountID, orgAccountID)
				c.Set("resource_account_id", orgAccountID)
				c.Set("is_team_member", true)
				c.Set("team_role", teamMember.Role)
				c.Set("user_role", teamMember.Role)
			} else {
				c.Set("resource_account_id", accountID)
				c.Set("is_team_member", false)
			}

			return next(c)
		}
	}
}

func GetResourceAccountID(c echo.Context) uuid.UUID {
	if id, ok := c.Get("resource_account_id").(uuid.UUID); ok {
		return id
	}

	if id, ok := c.Get("account_id").(uuid.UUID); ok {
		return id
	}

	return uuid.Nil
}

func GetPersonalAccountID(c echo.Context) uuid.UUID {
	if id, ok := c.Get("personal_account_id").(uuid.UUID); ok {
		return id
	}
	return c.Get("account_id").(uuid.UUID)
}

func IsTeamMember(c echo.Context) bool {
	if isTeam, ok := c.Get("is_team_member").(bool); ok {
		return isTeam
	}
	return false
}

func GetTeamRole(c echo.Context) (authModels.MemberRole, bool) {
	if role, ok := c.Get("team_role").(authModels.MemberRole); ok {
		return role, true
	}
	return "", false
}
