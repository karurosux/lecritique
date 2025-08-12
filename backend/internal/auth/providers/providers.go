package providers

import (
	"kyooar/internal/auth/repositories"
	"kyooar/internal/auth/services"
	"github.com/samber/do"
)

func ProvideAuthServices(i *do.Injector) {
	do.Provide(i, repositories.NewAccountRepository)
	do.Provide(i, repositories.NewTeamMemberRepository)
	do.Provide(i, repositories.NewTeamInvitationRepository)
	do.Provide(i, repositories.NewTokenRepository)
	
	do.Provide(i, services.NewAuthService)
	do.Provide(i, services.NewTeamMemberServiceV2)
}