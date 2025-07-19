package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"lecritique/internal/auth/models"
	"lecritique/internal/auth/services"
)

// TeamAuthMiddleware adds team member information to the request context
// Note: This should be called through MiddlewareProvider which provides proper DI
func TeamAuthMiddleware(teamMemberService services.TeamMemberServiceV2) echo.MiddlewareFunc {
	
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get account ID from JWT middleware
			accountID, ok := c.Get("account_id").(uuid.UUID)
			if !ok {
				return next(c)
			}

			// For now, since we're using account-based auth, we'll assume the logged-in account is the owner
			// In a real multi-user system, we'd need to track which user is logged in separately
			
			// Get team members for this account
			members, err := teamMemberService.ListMembers(c.Request().Context(), accountID)
			if err != nil {
				// Continue without team info
				return next(c)
			}

			// Find the owner
			for _, member := range members {
				if member.Role == models.RoleOwner {
					// For owner, use account_id as user_id since member_id might be invalid
					c.Set("user_id", accountID)
					c.Set("user_role", string(member.Role))
					break
				}
			}

			// If no owner found, set as owner by default (for backward compatibility)
			if c.Get("user_role") == nil {
				c.Set("user_id", accountID)
				c.Set("user_role", string(models.RoleOwner))
			}

			return next(c)
		}
	}
}