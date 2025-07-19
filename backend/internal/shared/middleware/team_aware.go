package middleware

import (
	"log"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	authModels "lecritique/internal/auth/models"
	"lecritique/internal/auth/services"
)

// TeamAware middleware checks if the user is a team member and sets the appropriate account IDs in context
// It sets:
// - resource_account_id: The account ID to use for accessing resources (org ID for team members)
// - personal_account_id: The user's personal account ID (always their own)
// - is_team_member: Boolean indicating if accessing as a team member
func TeamAware(teamMemberService services.TeamMemberServiceV2) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			accountID, err := GetAccountID(c)
			if err != nil {
				return next(c)
			}

			// Always set personal account ID
			c.Set("personal_account_id", accountID)

			// Check if this user is a team member of another organization
			teamMember, err := teamMemberService.GetMemberByMemberID(ctx, accountID)
			if err == nil && teamMember != nil {
				// Use the organization's account ID for resources
				orgAccountID := teamMember.AccountID
				log.Printf("TeamAware: Account %s is a team member of org %s", accountID, orgAccountID)
				c.Set("resource_account_id", orgAccountID)
				c.Set("is_team_member", true)
				c.Set("team_role", teamMember.Role)
				c.Set("user_role", teamMember.Role)
			} else {
				// Use their own account ID for resources
				c.Set("resource_account_id", accountID)
				c.Set("is_team_member", false)
			}

			return next(c)
		}
	}
}

// GetResourceAccountID is a helper to get the account ID for resource access
func GetResourceAccountID(c echo.Context) uuid.UUID {
	if id, ok := c.Get("resource_account_id").(uuid.UUID); ok {
		return id
	}

	if id, ok := c.Get("account_id").(uuid.UUID); ok {
		return id
	}

	// Fallback to regular account_id if middleware wasn't applied
	return uuid.Nil
}

// GetPersonalAccountID is a helper to get the user's personal account ID
func GetPersonalAccountID(c echo.Context) uuid.UUID {
	if id, ok := c.Get("personal_account_id").(uuid.UUID); ok {
		return id
	}
	// Fallback to regular account_id
	return c.Get("account_id").(uuid.UUID)
}

// IsTeamMember checks if the current user is accessing as a team member
func IsTeamMember(c echo.Context) bool {
	if isTeam, ok := c.Get("is_team_member").(bool); ok {
		return isTeam
	}
	return false
}

// GetTeamRole gets the team member's role if they are a team member
func GetTeamRole(c echo.Context) (authModels.MemberRole, bool) {
	if role, ok := c.Get("team_role").(authModels.MemberRole); ok {
		return role, true
	}
	return "", false
}
