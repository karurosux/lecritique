package providers

import (
	// Auth
	authHandlers "lecritique/internal/auth/handlers"
	authRepos "lecritique/internal/auth/repositories"
	authServices "lecritique/internal/auth/services"
	
	// Restaurant
	restaurantHandlers "lecritique/internal/restaurant/handlers"
	restaurantRepos "lecritique/internal/restaurant/repositories"
	restaurantServices "lecritique/internal/restaurant/services"
	
	// Menu
	menuHandlers "lecritique/internal/menu/handlers"
	menuRepos "lecritique/internal/menu/repositories"
	menuServices "lecritique/internal/menu/services"
	
	// Feedback
	feedbackHandlers "lecritique/internal/feedback/handlers"
	feedbackRepos "lecritique/internal/feedback/repositories"
	feedbackServices "lecritique/internal/feedback/services"
	
	// QR Code
	qrcodeHandlers "lecritique/internal/qrcode/handlers"
	qrcodeRepos "lecritique/internal/qrcode/repositories"
	qrcodeServices "lecritique/internal/qrcode/services"
	
	// Analytics
	analyticsHandlers "lecritique/internal/analytics/handlers"
	analyticsServices "lecritique/internal/analytics/services"
	
	// Subscription
	subscriptionHandlers "lecritique/internal/subscription/handlers"
	subscriptionMiddleware "lecritique/internal/subscription/middleware"
	subscriptionRepos "lecritique/internal/subscription/repositories"
	subscriptionServices "lecritique/internal/subscription/services"
	
	// Shared
	"lecritique/internal/shared/config"
	"lecritique/internal/shared/middleware"
	sharedServices "lecritique/internal/shared/services"
	
	"github.com/samber/do"
	"gorm.io/gorm"
)

// RegisterAll registers all services and repositories in the DI container
func RegisterAll(i *do.Injector, cfg *config.Config, db *gorm.DB) {
	// Core dependencies
	do.ProvideValue(i, cfg)
	do.ProvideValue(i, db)
	
	// Shared services
	do.Provide(i, sharedServices.NewEmailService)
	do.Provide(i, middleware.NewMiddlewareProvider)
	
	// Auth domain
	do.Provide(i, authRepos.NewAccountRepository)
	do.Provide(i, authRepos.NewTokenRepository)
	do.Provide(i, authRepos.NewTeamMemberRepository)
	do.Provide(i, authRepos.NewTeamInvitationRepository)
	do.Provide(i, authServices.NewAuthService)
	do.Provide(i, authServices.NewTeamMemberServiceV2)
	do.Provide(i, authHandlers.NewAuthHandler)
	do.Provide(i, authHandlers.NewTeamMemberHandler)
	
	// Restaurant domain
	do.Provide(i, restaurantRepos.NewRestaurantRepository)
	do.Provide(i, restaurantServices.NewRestaurantService)
	do.Provide(i, restaurantHandlers.NewRestaurantHandler)
	
	// Menu domain
	do.Provide(i, menuRepos.NewDishRepository)
	do.Provide(i, menuServices.NewDishService)
	do.Provide(i, menuHandlers.NewDishHandler)
	do.Provide(i, menuHandlers.NewMenuPublicHandler)
	
	// Feedback domain
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
	
	// QR Code domain
	do.Provide(i, qrcodeRepos.NewQRCodeRepository)
	do.Provide(i, qrcodeServices.NewQRCodeService)
	do.Provide(i, qrcodeHandlers.NewQRCodeHandler)
	do.Provide(i, qrcodeHandlers.NewQRCodePublicHandler)
	
	// Analytics domain
	do.Provide(i, analyticsServices.NewAnalyticsService)
	do.Provide(i, analyticsHandlers.NewAnalyticsHandler)
	
	// Subscription domain
	do.Provide(i, subscriptionRepos.NewSubscriptionRepository)
	do.Provide(i, subscriptionRepos.NewSubscriptionPlanRepository)
	do.Provide(i, subscriptionRepos.NewUsageRepository)
	do.Provide(i, subscriptionServices.NewSubscriptionService)
	do.Provide(i, subscriptionServices.NewUsageService)
	do.Provide(i, subscriptionServices.NewPaymentService)
	do.Provide(i, subscriptionHandlers.NewSubscriptionHandler)
	do.Provide(i, subscriptionHandlers.NewPaymentHandler)
	
	// Subscription middleware - depends on subscription services
	do.Provide(i, func(i *do.Injector) (*subscriptionMiddleware.SubscriptionMiddleware, error) {
		return subscriptionMiddleware.NewSubscriptionMiddleware(
			do.MustInvoke[subscriptionServices.SubscriptionService](i),
			do.MustInvoke[subscriptionServices.UsageService](i),
		), nil
	})
}