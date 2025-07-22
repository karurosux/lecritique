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
	// Get handlers from injector
	subscriptionHandler := do.MustInvoke[*handlers.SubscriptionHandler](m.injector)
	paymentHandler := do.MustInvoke[*handlers.PaymentHandler](m.injector)
	
	// Get middleware provider
	middlewareProvider := do.MustInvoke[*sharedMiddleware.MiddlewareProvider](m.injector)
	
	// Public routes (no authentication required)
	v1.GET("/plans", subscriptionHandler.GetAvailablePlans)
	
	// Protected routes (authentication required)
	user := v1.Group("/user")
	user.Use(middlewareProvider.AuthMiddleware())
	
	// User subscription routes
	user.GET("/subscription", subscriptionHandler.GetUserSubscription)
	user.GET("/subscription/usage", subscriptionHandler.GetUserUsage)
	user.POST("/subscription", subscriptionHandler.CreateSubscription)
	user.DELETE("/subscription", subscriptionHandler.CancelSubscription)
	
	// Permission checking routes
	user.GET("/can-create-organization", subscriptionHandler.CanUserCreateOrganization)
	
	// Payment routes
	payment := v1.Group("/payment")
	payment.Use(middlewareProvider.AuthMiddleware())
	payment.POST("/create-session", paymentHandler.CreateCheckoutSession)
	payment.POST("/webhook", paymentHandler.HandleWebhook)
}
