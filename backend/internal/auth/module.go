package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"gorm.io/gorm"

	authcontroller "kyooar/internal/auth/controller"
	authinterface "kyooar/internal/auth/interface"
	authmiddleware "kyooar/internal/auth/middleware"
	gormrepo "kyooar/internal/auth/repository/gorm"
	authservice "kyooar/internal/auth/service"
	"kyooar/internal/shared/config"
	"kyooar/internal/shared/services"
	"kyooar/internal/shared/validator"
	subscriptionServices "kyooar/internal/subscription/services"
)

func ProvideAccountRepository(i *do.Injector) (authinterface.AccountRepository, error) {
	db := do.MustInvoke[*gorm.DB](i)
	return gormrepo.NewAccountRepository(db), nil
}

func ProvideTokenRepository(i *do.Injector) (authinterface.TokenRepository, error) {
	db := do.MustInvoke[*gorm.DB](i)
	return gormrepo.NewTokenRepository(db), nil
}

func ProvideTeamMemberRepository(i *do.Injector) (authinterface.TeamMemberRepository, error) {
	db := do.MustInvoke[*gorm.DB](i)
	return gormrepo.NewTeamMemberRepository(db), nil
}

func ProvideTeamInvitationRepository(i *do.Injector) (authinterface.TeamInvitationRepository, error) {
	db := do.MustInvoke[*gorm.DB](i)
	return gormrepo.NewTeamInvitationRepository(db), nil
}

func ProvidePasswordHasher(i *do.Injector) (authinterface.PasswordHasher, error) {
	return authservice.NewBcryptPasswordHasher(12), nil
}

func ProvideTokenGenerator(i *do.Injector) (authinterface.TokenGenerator, error) {
	return authservice.NewTokenGenerator(), nil
}

func ProvideAuthService(i *do.Injector) (authinterface.AuthService, error) {
	accountRepo := do.MustInvoke[authinterface.AccountRepository](i)
	tokenRepo := do.MustInvoke[authinterface.TokenRepository](i)
	emailService := do.MustInvoke[services.EmailService](i)
	teamService := do.MustInvoke[authinterface.TeamMemberService](i)
	subscriptionService := do.MustInvoke[subscriptionServices.SubscriptionService](i)
	config := do.MustInvoke[*config.Config](i)

	return authservice.NewAuthService(
		accountRepo,
		tokenRepo,
		emailService,
		teamService,
		subscriptionService,
		config,
	), nil
}

func ProvideAuthMiddleware(i *do.Injector) (*authmiddleware.AuthMiddleware, error) {
	authService := do.MustInvoke[authinterface.AuthService](i)
	return authmiddleware.NewAuthMiddleware(authService), nil
}

func ProvideAuthController(i *do.Injector) (*authcontroller.AuthController, error) {
	authService := do.MustInvoke[authinterface.AuthService](i)
	teamMemberService := do.MustInvoke[authinterface.TeamMemberService](i)
	validator := validator.New()
	config := do.MustInvoke[*config.Config](i)

	return authcontroller.NewAuthController(
		authService,
		teamMemberService,
		validator,
		config,
	), nil
}

func ProvideTeamController(i *do.Injector) (*authcontroller.TeamController, error) {
	teamMemberService := do.MustInvoke[authinterface.TeamMemberService](i)
	validator := validator.New()

	return authcontroller.NewTeamController(
		teamMemberService,
		validator,
	), nil
}

type NewModule struct {
	injector *do.Injector
}

func NewAuthModule(i *do.Injector) *NewModule {
	return &NewModule{injector: i}
}

func (m *NewModule) RegisterRoutes(v1 *echo.Group) {
	authController := do.MustInvoke[*authcontroller.AuthController](m.injector)
	teamController := do.MustInvoke[*authcontroller.TeamController](m.injector)
	authMiddleware := do.MustInvoke[*authmiddleware.AuthMiddleware](m.injector)
	
	// Register routes with the middleware
	authController.RegisterRoutes(v1, authMiddleware.RequireAuth())
	teamController.RegisterRoutes(v1, authMiddleware.RequireAuth())
}

func RegisterNewModule(container *do.Injector) error {
	do.Provide(container, ProvideAccountRepository)
	do.Provide(container, ProvideTokenRepository)
	do.Provide(container, ProvideTeamMemberRepository)
	do.Provide(container, ProvideTeamInvitationRepository)
	do.Provide(container, ProvidePasswordHasher)
	do.Provide(container, ProvideTokenGenerator)
	do.Provide(container, ProvideAuthService)
	do.Provide(container, ProvideAuthMiddleware)
	do.Provide(container, ProvideAuthController)
	do.Provide(container, ProvideTeamController)

	return nil
}