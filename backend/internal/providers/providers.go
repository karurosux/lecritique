package providers

import (
	authController "kyooar/internal/auth/controller"
	authRepos "kyooar/internal/auth/repository/gorm"
	authServices "kyooar/internal/auth/service"
	authinterface "kyooar/internal/auth/interface"
	authmiddleware "kyooar/internal/auth/middleware"
	
	organization "kyooar/internal/organization"
	
	productHandlers "kyooar/internal/product/handlers"
	productRepos "kyooar/internal/product/repositories"
	productServices "kyooar/internal/product/services"
	
	
	
	analytics "kyooar/internal/analytics"
	feedback "kyooar/internal/feedback"
	subscriptioninterface "kyooar/internal/subscription/interface"
	
	
	"kyooar/internal/shared/config"
	"kyooar/internal/shared/middleware"
	sharedServices "kyooar/internal/shared/services"
	"kyooar/internal/shared/validator"
	
	"github.com/samber/do"
	"gorm.io/gorm"
)

func RegisterAll(i *do.Injector, cfg *config.Config, db *gorm.DB) {
	do.ProvideValue(i, cfg)
	do.ProvideValue(i, db)
	
	do.Provide(i, sharedServices.NewEmailService)
	do.Provide(i, middleware.NewMiddlewareProvider)
	
	// Auth repositories
	do.Provide(i, func(i *do.Injector) (authinterface.AccountRepository, error) {
		db := do.MustInvoke[*gorm.DB](i)
		return authRepos.NewAccountRepository(db), nil
	})
	do.Provide(i, func(i *do.Injector) (authinterface.TokenRepository, error) {
		db := do.MustInvoke[*gorm.DB](i)
		return authRepos.NewTokenRepository(db), nil
	})
	do.Provide(i, func(i *do.Injector) (authinterface.TeamMemberRepository, error) {
		db := do.MustInvoke[*gorm.DB](i)
		return authRepos.NewTeamMemberRepository(db), nil
	})
	do.Provide(i, func(i *do.Injector) (authinterface.TeamInvitationRepository, error) {
		db := do.MustInvoke[*gorm.DB](i)
		return authRepos.NewTeamInvitationRepository(db), nil
	})
	
	// Auth utility services
	do.Provide(i, func(i *do.Injector) (authinterface.PasswordHasher, error) {
		return authServices.NewBcryptPasswordHasher(12), nil
	})
	do.Provide(i, func(i *do.Injector) (authinterface.TokenGenerator, error) {
		return authServices.NewTokenGenerator(), nil
	})
	
	// Auth services - put TeamMemberService before AuthService to avoid dependency cycle
	do.Provide(i, func(i *do.Injector) (authinterface.TeamMemberService, error) {
		teamMemberRepo := do.MustInvoke[authinterface.TeamMemberRepository](i)
		invitationRepo := do.MustInvoke[authinterface.TeamInvitationRepository](i)
		accountRepo := do.MustInvoke[authinterface.AccountRepository](i)
		tokenGenerator := do.MustInvoke[authinterface.TokenGenerator](i)
		// emailSender will be nil for now
		
		return authServices.NewTeamMemberService(
			teamMemberRepo,
			invitationRepo,
			accountRepo,
			tokenGenerator,
			nil, // emailSender - temporarily nil
		), nil
	})
	
	do.Provide(i, func(i *do.Injector) (authinterface.AuthService, error) {
		accountRepo := do.MustInvoke[authinterface.AccountRepository](i)
		tokenRepo := do.MustInvoke[authinterface.TokenRepository](i)
		emailService := do.MustInvoke[sharedServices.EmailService](i)
		teamService := do.MustInvoke[authinterface.TeamMemberService](i)
		subscriptionService := do.MustInvoke[subscriptioninterface.SubscriptionService](i)
		config := do.MustInvoke[*config.Config](i)
		
		return authServices.NewAuthService(
			accountRepo,
			tokenRepo,
			emailService,
			teamService,
			subscriptionService,
			config,
		), nil
	})
	
	// Auth middleware
	do.Provide(i, func(i *do.Injector) (*authmiddleware.AuthMiddleware, error) {
		authService := do.MustInvoke[authinterface.AuthService](i)
		return authmiddleware.NewAuthMiddleware(authService), nil
	})
	
	// Auth controllers
	do.Provide(i, func(i *do.Injector) (*authController.AuthController, error) {
		authService := do.MustInvoke[authinterface.AuthService](i)
		teamMemberService := do.MustInvoke[authinterface.TeamMemberService](i)
		config := do.MustInvoke[*config.Config](i)
		
		return authController.NewAuthController(
			authService,
			teamMemberService,
			validator.New(),
			config,
		), nil
	})
	
	do.Provide(i, func(i *do.Injector) (*authController.TeamController, error) {
		teamMemberService := do.MustInvoke[authinterface.TeamMemberService](i)
		
		return authController.NewTeamController(
			teamMemberService,
			validator.New(),
		), nil
	})
	
	// Organization module registration
	organization.RegisterNewModule(i)
	
	do.Provide(i, productRepos.NewProductRepository)
	do.Provide(i, productServices.NewProductService)
	do.Provide(i, productHandlers.NewProductHandler)
	do.Provide(i, productHandlers.NewMenuPublicHandler)
	
	// Feedback repositories, services and handlers are now provided by the new feedback module
	
	// QR code providers are now in the qrcode module
	
	// Analytics module registration
	analytics.RegisterModule(i)
	
	// Feedback module registration
	feedback.RegisterNewModule(i)
	
	// Subscription providers are now in the subscription module
}
