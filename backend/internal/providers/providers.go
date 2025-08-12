package providers

import (
	authHandlers "kyooar/internal/auth/handlers"
	authRepos "kyooar/internal/auth/repositories"
	authServices "kyooar/internal/auth/services"
	
	organizationHandlers "kyooar/internal/organization/handlers"
	organizationRepos "kyooar/internal/organization/repositories"
	organizationServices "kyooar/internal/organization/services"
	
	menuHandlers "kyooar/internal/menu/handlers"
	menuRepos "kyooar/internal/menu/repositories"
	menuServices "kyooar/internal/menu/services"
	
	feedbackHandlers "kyooar/internal/feedback/handlers"
	feedbackRepos "kyooar/internal/feedback/repositories"
	feedbackServices "kyooar/internal/feedback/services"
	
	qrcodeHandlers "kyooar/internal/qrcode/handlers"
	qrcodeRepos "kyooar/internal/qrcode/repositories"
	qrcodeServices "kyooar/internal/qrcode/services"
	
	analyticsHandlers "kyooar/internal/analytics/handlers"
	analyticsServices "kyooar/internal/analytics/services"
	
	subscriptionHandlers "kyooar/internal/subscription/handlers"
	subscriptionMiddleware "kyooar/internal/subscription/middleware"
	subscriptionRepos "kyooar/internal/subscription/repositories"
	subscriptionServices "kyooar/internal/subscription/services"
	
	"kyooar/internal/shared/config"
	"kyooar/internal/shared/middleware"
	sharedServices "kyooar/internal/shared/services"
	
	"github.com/samber/do"
	"gorm.io/gorm"
)

func RegisterAll(i *do.Injector, cfg *config.Config, db *gorm.DB) {
	do.ProvideValue(i, cfg)
	do.ProvideValue(i, db)
	
	do.Provide(i, sharedServices.NewEmailService)
	do.Provide(i, middleware.NewMiddlewareProvider)
	
	do.Provide(i, authRepos.NewAccountRepository)
	do.Provide(i, authRepos.NewTokenRepository)
	do.Provide(i, authRepos.NewTeamMemberRepository)
	do.Provide(i, authRepos.NewTeamInvitationRepository)
	do.Provide(i, authServices.NewAuthService)
	do.Provide(i, authServices.NewTeamMemberServiceV2)
	do.Provide(i, authHandlers.NewAuthHandler)
	do.Provide(i, authHandlers.NewTeamMemberHandler)
	
	do.Provide(i, organizationRepos.NewOrganizationRepository)
	do.Provide(i, organizationServices.NewOrganizationService)
	do.Provide(i, organizationHandlers.NewOrganizationHandler)
	
	do.Provide(i, menuRepos.NewProductRepository)
	do.Provide(i, menuServices.NewProductService)
	do.Provide(i, menuHandlers.NewProductHandler)
	do.Provide(i, menuHandlers.NewMenuPublicHandler)
	
	do.Provide(i, feedbackRepos.NewFeedbackRepository)
	do.Provide(i, feedbackRepos.NewQuestionnaireRepository)
	do.Provide(i, feedbackRepos.NewQuestionRepository)
	do.Provide(i, feedbackRepos.NewQuestionTemplateRepository)
	do.Provide(i, feedbackServices.NewFeedbackService)
	do.Provide(i, feedbackServices.NewQuestionnaireService)
	do.Provide(i, feedbackServices.NewQuestionService)
	do.Provide(i, feedbackHandlers.NewFeedbackHandler)
	do.Provide(i, feedbackHandlers.NewQuestionnaireHandler)
	do.Provide(i, feedbackHandlers.NewQuestionHandler)
	do.Provide(i, feedbackHandlers.NewFeedbackPublicHandler)
	
	do.Provide(i, qrcodeRepos.NewQRCodeRepository)
	do.Provide(i, qrcodeServices.NewQRCodeService)
	do.Provide(i, qrcodeHandlers.NewQRCodeHandler)
	do.Provide(i, qrcodeHandlers.NewQRCodePublicHandler)
	
	do.Provide(i, analyticsServices.NewAnalyticsService)
	do.Provide(i, analyticsHandlers.NewAnalyticsHandler)
	
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
			do.MustInvoke[organizationServices.OrganizationService](i),
		), nil
	})
}
