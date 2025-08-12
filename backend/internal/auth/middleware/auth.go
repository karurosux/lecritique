package authmiddleware

import (
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	authconstants "kyooar/internal/auth/constants"
	authinterface "kyooar/internal/auth/interface"
	"kyooar/internal/auth/models"
	"kyooar/internal/shared/errors"
	"kyooar/internal/shared/response"
)

type AuthMiddleware struct {
	authService authinterface.AuthService
}

func NewAuthMiddleware(authService authinterface.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func (m *AuthMiddleware) RequireAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := extractToken(c)
			if token == "" {
				return response.Error(c, errors.ErrUnauthorized)
			}

			claims, err := m.authService.ValidateToken(token)
			if err != nil {
				return response.Error(c, errors.ErrUnauthorized)
			}

			c.Set(string(authconstants.AccountIDKey), claims.AccountID)
			c.Set(string(authconstants.MemberIDKey), claims.MemberID)
			c.Set(string(authconstants.ClaimsKey), claims)

			return next(c)
		}
	}
}

func RequireVerifiedEmail() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, ok := c.Get(string(authconstants.ClaimsKey)).(*authinterface.Claims)
			if !ok {
				return response.Error(c, errors.ErrUnauthorized)
			}

			account, err := GetAccountFromClaims(c, claims)
			if err != nil {
				return response.Error(c, err)
			}

			if !account.EmailVerified {
				return response.Error(c, errors.NewWithDetails("EMAIL_NOT_VERIFIED", "Email verification required", 403, nil))
			}

			return next(c)
		}
	}
}

func extractToken(c echo.Context) string {
	auth := c.Request().Header.Get("Authorization")
	if auth == "" {
		return ""
	}

	parts := strings.Split(auth, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}

func GetAccountIDFromContext(c echo.Context) (uuid.UUID, error) {
	accountID, ok := c.Get(string(authconstants.AccountIDKey)).(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.ErrUnauthorized
	}
	return accountID, nil
}

func GetMemberIDFromContext(c echo.Context) (uuid.UUID, error) {
	memberID, ok := c.Get(string(authconstants.MemberIDKey)).(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.ErrUnauthorized
	}
	return memberID, nil
}

func GetClaimsFromContext(c echo.Context) (*authinterface.Claims, error) {
	claims, ok := c.Get(string(authconstants.ClaimsKey)).(*authinterface.Claims)
	if !ok {
		return nil, errors.ErrUnauthorized
	}
	return claims, nil
}

func GetAccountFromClaims(c echo.Context, claims *authinterface.Claims) (*models.Account, error) {
	return nil, errors.NewWithDetails("NOT_IMPLEMENTED", "Account retrieval from claims not implemented", 500, nil)
}