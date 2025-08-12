package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"gorm.io/gorm"
	authinterface "kyooar/internal/auth/interface"
)

type MiddlewareProvider struct {
	db                  *gorm.DB
	authService         authinterface.AuthService
	teamMemberService   authinterface.TeamMemberService
}

func NewMiddlewareProvider(i *do.Injector) (*MiddlewareProvider, error) {
	return &MiddlewareProvider{
		db:                  do.MustInvoke[*gorm.DB](i),
		authService:         do.MustInvoke[authinterface.AuthService](i),
		teamMemberService:   do.MustInvoke[authinterface.TeamMemberService](i),
	}, nil
}

func (p *MiddlewareProvider) TeamAwareMiddleware() echo.MiddlewareFunc {
	return TeamAware(p.teamMemberService)
}

func (p *MiddlewareProvider) AuthMiddleware() echo.MiddlewareFunc {
	return JWTAuth(p.authService)
}

func (p *MiddlewareProvider) TeamAuthMiddleware() echo.MiddlewareFunc {
	return TeamAuthMiddleware(p.teamMemberService)
}