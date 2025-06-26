package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lecritique/api/internal/auth/models"
	"github.com/lecritique/api/internal/auth/repositories"
	"gorm.io/gorm"
)

// TeamAuthMiddleware adds team member information to the request context
func TeamAuthMiddleware(db *gorm.DB) echo.MiddlewareFunc {
	teamMemberRepo := repositories.NewTeamMemberRepository(db)
	
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
			members, err := teamMemberRepo.FindByAccountID(c.Request().Context(), accountID)
			if err != nil {
				// Continue without team info
				return next(c)
			}

			// Find the owner
			for _, member := range members {
				if member.Role == models.RoleOwner {
					c.Set("user_id", member.UserID)
					c.Set("user_role", string(member.Role))
					break
				}
			}

			// If no owner found, set as owner by default (for backward compatibility)
			if c.Get("user_role") == nil {
				c.Set("user_role", string(models.RoleOwner))
			}

			return next(c)
		}
	}
}