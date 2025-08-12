package providers

import (
	authController "kyooar/internal/auth/controller"
	authRepos "kyooar/internal/auth/repository/gorm"
	authServices "kyooar/internal/auth/service"
	authinterface "kyooar/internal/auth/interface"
	authmiddleware "kyooar/internal/auth/middleware"
	
	organization "kyooar/internal/organization"
	organizationinterface "kyooar/internal/organization/interface"
	
	menuHandlers "kyooar/internal/menu/handlers"
	menuRepos "kyooar/internal/menu/repositories"
	menuServices "kyooar/internal/menu/services"
	
	
	
	analytics "kyooar/internal/analytics"
	feedback "kyooar/internal/feedback"
	
	subscriptionHandlers "kyooar/internal/subscription/handlers"
	subscriptionMiddleware "kyooar/internal/subscription/middleware"
	subscriptionRepos "kyooar/internal/subscription/repositories"
	subscriptionServices "kyooar/internal/subscription/services"
	
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
		subscriptionService := do.MustInvoke[subscriptionServices.SubscriptionService](i)
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
	
	do.Provide(i, menuRepos.NewProductRepository)
	do.Provide(i, menuServices.NewProductService)
	do.Provide(i, menuHandlers.NewProductHandler)
	do.Provide(i, menuHandlers.NewMenuPublicHandler)
	
	// Feedback repositories, services and handlers are now provided by the new feedback module
	
	// QR code providers are now in the qrcode module
	
	// Analytics module registration
	analytics.RegisterModule(i)
	
	// Feedback module registration
	feedback.RegisterNewModule(i)
	
	do.Provide(i, subscriptionRepos.NewSubscriptionRepository)
	do.Provide(i, subscriptionRepos.NewSubscriptionPlanRepository)
	do.Provide(i, subscriptionRepos.NewUsageRepository)
	do.Provide(i, subscriptionServices.NewSubscriptionService)
	do.Provide(i, subscriptionServices.NewUsageService)
	do.Provide(i, subscriptionServices.NewPaymentService)
	do.Provide(i, subscriptionHandlers.NewSubscriptionHandler)
	do.Provide(i, subscriptionHandlers.NewPaymentHandler)
	
	do.Provide(i, func(i *do.Injector) (*subscriptionMiddleware.SubscriptionMiddleware, error) {
		return subscriptionMiddleware.NewSubscriptionMiddleware(
			do.MustInvoke[subscriptionServices.SubscriptionService](i),
			do.MustInvoke[subscriptionServices.UsageService](i),
			do.MustInvoke[organizationinterface.OrganizationService](i),
		), nil
	})
}
