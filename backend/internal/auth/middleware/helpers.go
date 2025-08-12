package authmiddleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	authconstants "kyooar/internal/auth/constants"
	"kyooar/internal/auth/models"
	"kyooar/internal/shared/errors"
)

func GetUserIDFromContext(c echo.Context) (uuid.UUID, error) {
	return GetMemberIDFromContext(c)
}

func GetAccountFromContext(c echo.Context) *models.Account {
	account, ok := c.Get(string(authconstants.UserKey)).(*models.Account)
	if !ok {
		return nil
	}
	return account
}

func SetAccountInContext(c echo.Context, account *models.Account) {
	c.Set(string(authconstants.UserKey), account)
}

func RequireRole(requiredRole models.MemberRole) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, err := GetClaimsFromContext(c)
			if err != nil {
				return err
			}

			if !hasRole(claims.Role, requiredRole) {
				return errors.NewWithDetails("INSUFFICIENT_PRIVILEGES", "Insufficient privileges for this operation", 403, nil)
			}

			return next(c)
		}
	}
}

func hasRole(userRole, requiredRole models.MemberRole) bool {
	roleHierarchy := map[models.MemberRole]int{
		models.RoleViewer:  1,
		models.RoleManager: 2,
		models.RoleAdmin:   3,
		models.RoleOwner:   4,
	}

	userLevel, userExists := roleHierarchy[userRole]
	requiredLevel, requiredExists := roleHierarchy[requiredRole]

	if !userExists || !requiredExists {
		return false
	}

	return userLevel >= requiredLevel
}