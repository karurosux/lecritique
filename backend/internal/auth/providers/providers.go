package providers

import (
	"lecritique/internal/auth/repositories"
	"lecritique/internal/auth/services"
	"github.com/samber/do"
)

// ProvideAuthServices registers all auth-related services and repositories
func ProvideAuthServices(i *do.Injector) {
	// Repositories
	do.Provide(i, repositories.NewAccountRepository)
	do.Provide(i, repositories.NewTeamMemberRepository)
	do.Provide(i, repositories.NewTeamInvitationRepository)
	do.Provide(i, repositories.NewTokenRepository)
	
	// Services
	do.Provide(i, services.NewAuthService)
	do.Provide(i, services.NewTeamMemberServiceV2)
}