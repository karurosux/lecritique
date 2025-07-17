package middleware

import (
	"log"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	authModels "github.com/lecritique/api/internal/auth/models"
	"gorm.io/gorm"
)

// TeamAware middleware checks if the user is a team member and sets the appropriate account IDs in context
// It sets:
// - resource_account_id: The account ID to use for accessing resources (org ID for team members)
// - personal_account_id: The user's personal account ID (always their own)
// - is_team_member: Boolean indicating if accessing as a team member
func TeamAware(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			accountID := c.Get("account_id").(uuid.UUID)
			
			// Always set personal account ID
			c.Set("personal_account_id", accountID)
			
			// Check if this user is a team member of another organization
			var teamMemberships []authModels.TeamMember
			db.WithContext(ctx).
				Where("member_id = ? AND account_id != ? AND accepted_at IS NOT NULL", accountID, accountID).
				Find(&teamMemberships)
			
			if len(teamMemberships) > 0 {
				// Use the organization's account ID for resources
				orgAccountID := teamMemberships[0].AccountID
				log.Printf("TeamAware: Account %s is a team member of org %s", accountID, orgAccountID)
				c.Set("resource_account_id", orgAccountID)
				c.Set("is_team_member", true)
				c.Set("team_role", teamMemberships[0].Role)
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
	// Fallback to regular account_id if middleware wasn't applied
	return c.Get("account_id").(uuid.UUID)
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