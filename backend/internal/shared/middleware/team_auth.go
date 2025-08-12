package middleware

import (
	"github.com/labstack/echo/v4"
	"kyooar/internal/auth/models"
	"kyooar/internal/auth/services"
)

func TeamAuthMiddleware(teamMemberService services.TeamMemberServiceV2) echo.MiddlewareFunc {
	
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			accountID, err := GetAccountID(c)
			if err != nil {
				return next(c)
			}

			members, err := teamMemberService.ListMembers(c.Request().Context(), accountID)
			if err != nil {
				return next(c)
			}

			for _, member := range members {
				if member.Role == models.RoleOwner {
					c.Set("user_id", accountID)
					c.Set("user_role", string(member.Role))
					break
				}
			}

			if c.Get("user_role") == nil {
				c.Set("user_id", accountID)
				c.Set("user_role", string(models.RoleOwner))
			}

			return next(c)
		}
	}
}