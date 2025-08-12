package subscription

import (
	"github.com/labstack/echo/v4"
	"kyooar/internal/subscription/handlers"
	sharedMiddleware "kyooar/internal/shared/middleware"
	"github.com/samber/do"
)

type Module struct {
	injector *do.Injector
}

func NewModule(i *do.Injector) *Module {
	return &Module{injector: i}
}

func (m *Module) RegisterRoutes(v1 *echo.Group) {
	subscriptionHandler := do.MustInvoke[*handlers.SubscriptionHandler](m.injector)
	paymentHandler := do.MustInvoke[*handlers.PaymentHandler](m.injector)
	
	middlewareProvider := do.MustInvoke[*sharedMiddleware.MiddlewareProvider](m.injector)
	v1.GET("/plans", subscriptionHandler.GetAvailablePlans)
	user := v1.Group("/user")
	user.Use(middlewareProvider.AuthMiddleware())
	user.GET("/subscription", subscriptionHandler.GetUserSubscription)
	user.GET("/subscription/usage", subscriptionHandler.GetUserUsage)
	user.POST("/subscription", subscriptionHandler.CreateSubscription)
	user.DELETE("/subscription", subscriptionHandler.CancelSubscription)
	user.GET("/can-create-organization", subscriptionHandler.CanUserCreateOrganization)
	payment := v1.Group("/payment")
	payment.Use(middlewareProvider.AuthMiddleware())
	payment.POST("/create-session", paymentHandler.CreateCheckoutSession)
	payment.POST("/webhook", paymentHandler.HandleWebhook)
}
