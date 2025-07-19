package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"gorm.io/gorm"
	authServices "lecritique/internal/auth/services"
)

// MiddlewareProvider provides middleware functions with injected dependencies
type MiddlewareProvider struct {
	db          *gorm.DB
	authService authServices.AuthService
}

// NewMiddlewareProvider creates a new middleware provider
func NewMiddlewareProvider(i *do.Injector) (*MiddlewareProvider, error) {
	return &MiddlewareProvider{
		db:          do.MustInvoke[*gorm.DB](i),
		authService: do.MustInvoke[authServices.AuthService](i),
	}, nil
}

// TeamAwareMiddleware returns the team aware middleware
func (p *MiddlewareProvider) TeamAwareMiddleware() echo.MiddlewareFunc {
	return TeamAware(p.db)
}

// AuthMiddleware returns the JWT authentication middleware
func (p *MiddlewareProvider) AuthMiddleware() echo.MiddlewareFunc {
	return JWTAuth(p.authService)
}

// TeamAuthMiddleware returns the team auth middleware
func (p *MiddlewareProvider) TeamAuthMiddleware() echo.MiddlewareFunc {
	return TeamAuthMiddleware(p.db)
}