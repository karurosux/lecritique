package providers

import (
	"lecritique/internal/auth/repositories"
	"lecritique/internal/auth/services"
	"github.com/samber/do"
)

// ProvideAuthServices registers all auth-related services and repositories
func ProvideAuthServices(i *do.Injector) {
	// Repositories
	do.Provide(i, repositories.NewUserRepository)
	do.Provide(i, repositories.NewTeamRepository)
	do.Provide(i, repositories.NewTeamMemberRepository)
	do.Provide(i, repositories.NewInvitationRepository)
	
	// Services
	do.Provide(i, services.NewAuthService)
	do.Provide(i, services.NewTeamService)
	do.Provide(i, services.NewTeamMemberService)
}