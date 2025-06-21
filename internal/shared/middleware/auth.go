package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lecritique/api/internal/auth/services"
	"github.com/lecritique/api/internal/shared/errors"
	"github.com/lecritique/api/internal/shared/response"
)

func JWTAuth(authService services.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get token from header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return response.Error(c, errors.ErrUnauthorized)
			}

			// Check Bearer prefix
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return response.Error(c, errors.ErrUnauthorized)
			}

			token := parts[1]

			// Validate token
			claims, err := authService.ValidateToken(token)
			if err != nil {
				return response.Error(c, err)
			}

			// Set account ID in context
			c.Set("account_id", claims.AccountID)
			c.Set("email", claims.Email)

			return next(c)
		}
	}
}
