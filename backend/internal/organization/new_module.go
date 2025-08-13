package organization

import (
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"gorm.io/gorm"

	feedbackcontroller "kyooar/internal/feedback/controller"
	productHandlers "kyooar/internal/product/handlers"
	organizationcontroller "kyooar/internal/organization/controller"
	organizationinterface "kyooar/internal/organization/interface"
	gormrepo "kyooar/internal/organization/repository/gorm"
	organizationservice "kyooar/internal/organization/service"
	qrcodecontroller "kyooar/internal/qrcode/controller"
	sharedMiddleware "kyooar/internal/shared/middleware"
	subscriptionMiddleware "kyooar/internal/subscription/middleware"
	subscriptioninterface "kyooar/internal/subscription/interface"
)

func ProvideOrganizationRepository(i *do.Injector) (organizationinterface.OrganizationRepository, error) {
	db := do.MustInvoke[*gorm.DB](i)
	return gormrepo.NewOrganizationRepository(db), nil
}

func ProvideOrganizationService(i *do.Injector) (organizationinterface.OrganizationService, error) {
	organizationRepo := do.MustInvoke[organizationinterface.OrganizationRepository](i)
	subscriptionRepo := do.MustInvoke[subscriptioninterface.SubscriptionRepository](i)

	return organizationservice.NewOrganizationService(
		organizationRepo,
		subscriptionRepo,
	), nil
}

func ProvideOrganizationController(i *do.Injector) (*organizationcontroller.OrganizationController, error) {
	organizationService := do.MustInvoke[organizationinterface.OrganizationService](i)
	productHandler := do.MustInvoke[*productHandlers.ProductHandler](i)
	qrCodeHandler := do.MustInvoke[*qrcodecontroller.QRCodeController](i)
	feedbackController := do.MustInvoke[*feedbackcontroller.FeedbackController](i)
	questionnaireController := do.MustInvoke[*feedbackcontroller.QuestionnaireController](i)
	questionController := do.MustInvoke[*feedbackcontroller.QuestionController](i)
	
	return organizationcontroller.NewOrganizationController(
		organizationService,
		productHandler,
		qrCodeHandler,
		feedbackController,
		questionnaireController,
		questionController,
	), nil
}

type OrganizationModule struct {
	injector *do.Injector
}

func NewOrganizationModule(i *do.Injector) *OrganizationModule {
	return &OrganizationModule{injector: i}
}

func (m *OrganizationModule) RegisterRoutes(v1 *echo.Group) {
	organizationController := do.MustInvoke[*organizationcontroller.OrganizationController](m.injector)
	middlewareProvider := do.MustInvoke[*sharedMiddleware.MiddlewareProvider](m.injector)
	subscriptionMW := do.MustInvoke[*subscriptionMiddleware.SubscriptionMiddleware](m.injector)

	organizationController.RegisterRoutes(v1, middlewareProvider, subscriptionMW)
}

func RegisterNewModule(container *do.Injector) error {
	do.Provide(container, ProvideOrganizationRepository)
	do.Provide(container, ProvideOrganizationService)
	do.Provide(container, ProvideOrganizationController)

	return nil
}