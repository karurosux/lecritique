package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"gorm.io/gorm"
	authServices "kyooar/internal/auth/services"
)

type MiddlewareProvider struct {
	db                  *gorm.DB
	authService         authServices.AuthService
	teamMemberService   authServices.TeamMemberServiceV2
}

func NewMiddlewareProvider(i *do.Injector) (*MiddlewareProvider, error) {
	return &MiddlewareProvider{
		db:                  do.MustInvoke[*gorm.DB](i),
		authService:         do.MustInvoke[authServices.AuthService](i),
		teamMemberService:   do.MustInvoke[authServices.TeamMemberServiceV2](i),
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