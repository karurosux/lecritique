package subscriptioncontroller

import (
	"github.com/labstack/echo/v4"
	sharedMiddleware "kyooar/internal/shared/middleware"
	subscriptionMiddleware "kyooar/internal/subscription/middleware"
)

func (sc *SubscriptionController) RegisterRoutes(v1 *echo.Group, middlewareProvider *sharedMiddleware.MiddlewareProvider, subscriptionMW *subscriptionMiddleware.SubscriptionMiddleware) {
	// Public plans endpoint
	v1.GET("/plans", sc.GetAvailablePlans)
	
	userSubscriptions := v1.Group("/user", middlewareProvider.AuthMiddleware())
	userSubscriptions.GET("/subscription", sc.GetUserSubscription)
	userSubscriptions.POST("/subscription", sc.CreateSubscription)
	userSubscriptions.DELETE("/subscription", sc.CancelSubscription)
	userSubscriptions.GET("/subscription/usage", sc.GetUserUsage)
	userSubscriptions.GET("/can-create-organization", sc.CanUserCreateOrganization)
}

func (pc *PaymentController) RegisterRoutes(v1 *echo.Group, middlewareProvider *sharedMiddleware.MiddlewareProvider) {
	payments := v1.Group("/payment")
	
	payments.POST("/webhook", pc.HandleWebhook)
	
	userPayments := payments.Group("", middlewareProvider.AuthMiddleware())
	userPayments.POST("/checkout", pc.CreateCheckoutSession)
	userPayments.POST("/checkout/complete", pc.CompleteCheckout)
	userPayments.POST("/portal", pc.CreatePortalSession)
	userPayments.GET("/methods", pc.GetPaymentMethods)
	userPayments.POST("/methods/default", pc.SetDefaultPaymentMethod)
	userPayments.GET("/invoices", pc.GetInvoices)
}