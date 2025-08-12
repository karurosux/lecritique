package subscription

import (
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"gorm.io/gorm"

	authinterface "kyooar/internal/auth/interface"
	"kyooar/internal/shared/config"
	organizationinterface "kyooar/internal/organization/interface"
	subscriptioncontroller "kyooar/internal/subscription/controller"
	subscriptioninterface "kyooar/internal/subscription/interface"
	gormrepo "kyooar/internal/subscription/repository/gorm"
	subscriptionservice "kyooar/internal/subscription/service"
	sharedMiddleware "kyooar/internal/shared/middleware"
	subscriptionMiddleware "kyooar/internal/subscription/middleware"
)

func ProvideSubscriptionRepository(i *do.Injector) (subscriptioninterface.SubscriptionRepository, error) {
	db := do.MustInvoke[*gorm.DB](i)
	return gormrepo.NewSubscriptionRepository(db), nil
}

func ProvideSubscriptionPlanRepository(i *do.Injector) (subscriptioninterface.SubscriptionPlanRepository, error) {
	db := do.MustInvoke[*gorm.DB](i)
	return gormrepo.NewSubscriptionPlanRepository(db), nil
}

func ProvideUsageRepository(i *do.Injector) (subscriptioninterface.UsageRepository, error) {
	db := do.MustInvoke[*gorm.DB](i)
	return gormrepo.NewUsageRepository(db), nil
}

func ProvideSubscriptionService(i *do.Injector) (subscriptioninterface.SubscriptionService, error) {
	subscriptionRepo := do.MustInvoke[subscriptioninterface.SubscriptionRepository](i)
	planRepo := do.MustInvoke[subscriptioninterface.SubscriptionPlanRepository](i)
	organizationRepo := do.MustInvoke[organizationinterface.OrganizationRepository](i)

	return subscriptionservice.NewSubscriptionService(
		subscriptionRepo,
		planRepo,
		organizationRepo,
	), nil
}

func ProvideUsageService(i *do.Injector) (subscriptioninterface.UsageService, error) {
	usageRepo := do.MustInvoke[subscriptioninterface.UsageRepository](i)
	subscriptionRepo := do.MustInvoke[subscriptioninterface.SubscriptionRepository](i)

	return subscriptionservice.NewUsageService(
		usageRepo,
		subscriptionRepo,
	), nil
}

func ProvidePaymentProvider(i *do.Injector) (subscriptioninterface.PaymentProvider, error) {
	return subscriptionservice.NewStripeProvider(), nil
}

func ProvidePaymentService(i *do.Injector) (subscriptioninterface.PaymentService, error) {
	provider := do.MustInvoke[subscriptioninterface.PaymentProvider](i)
	subscriptionRepo := do.MustInvoke[subscriptioninterface.SubscriptionRepository](i)
	planRepo := do.MustInvoke[subscriptioninterface.SubscriptionPlanRepository](i)
	config := do.MustInvoke[*config.Config](i)

	return subscriptionservice.NewPaymentService(
		provider,
		subscriptionRepo,
		planRepo,
		config,
	), nil
}

func ProvideSubscriptionController(i *do.Injector) (*subscriptioncontroller.SubscriptionController, error) {
	subscriptionService := do.MustInvoke[subscriptioninterface.SubscriptionService](i)
	usageService := do.MustInvoke[subscriptioninterface.UsageService](i)
	teamMemberService := do.MustInvoke[authinterface.TeamMemberService](i)

	return subscriptioncontroller.NewSubscriptionController(
		subscriptionService,
		usageService,
		teamMemberService,
	), nil
}

func ProvidePaymentController(i *do.Injector) (*subscriptioncontroller.PaymentController, error) {
	paymentService := do.MustInvoke[subscriptioninterface.PaymentService](i)

	return subscriptioncontroller.NewPaymentController(
		paymentService,
	), nil
}

type SubscriptionModule struct {
	injector *do.Injector
}

func NewSubscriptionModule(i *do.Injector) *SubscriptionModule {
	return &SubscriptionModule{injector: i}
}

func (m *SubscriptionModule) RegisterRoutes(v1 *echo.Group) {
	subscriptionController := do.MustInvoke[*subscriptioncontroller.SubscriptionController](m.injector)
	paymentController := do.MustInvoke[*subscriptioncontroller.PaymentController](m.injector)
	middlewareProvider := do.MustInvoke[*sharedMiddleware.MiddlewareProvider](m.injector)
	subscriptionMW := do.MustInvoke[*subscriptionMiddleware.SubscriptionMiddleware](m.injector)

	subscriptionController.RegisterRoutes(v1, middlewareProvider, subscriptionMW)
	paymentController.RegisterRoutes(v1, middlewareProvider)
}

func ProvideSubscriptionMiddleware(i *do.Injector) (*subscriptionMiddleware.SubscriptionMiddleware, error) {
	subscriptionService := do.MustInvoke[subscriptioninterface.SubscriptionService](i)
	usageService := do.MustInvoke[subscriptioninterface.UsageService](i)
	organizationService := do.MustInvoke[organizationinterface.OrganizationService](i)

	return subscriptionMiddleware.NewSubscriptionMiddleware(
		subscriptionService,
		usageService,
		organizationService,
	), nil
}

func RegisterNewModule(container *do.Injector) error {
	do.Provide(container, ProvideSubscriptionRepository)
	do.Provide(container, ProvideSubscriptionPlanRepository)
	do.Provide(container, ProvideUsageRepository)
	do.Provide(container, ProvideSubscriptionService)
	do.Provide(container, ProvideUsageService)
	do.Provide(container, ProvidePaymentProvider)
	do.Provide(container, ProvidePaymentService)
	do.Provide(container, ProvideSubscriptionController)
	do.Provide(container, ProvidePaymentController)
	do.Provide(container, ProvideSubscriptionMiddleware)

	return nil
}