package organization

import (
	"github.com/labstack/echo/v4"
	"kyooar/internal/organization/handlers"
	feedbackHandlers "kyooar/internal/feedback/handlers"
	menuHandlers "kyooar/internal/menu/handlers"
	qrcodeHandlers "kyooar/internal/qrcode/handlers"
	sharedMiddleware "kyooar/internal/shared/middleware"
	subscriptionMiddleware "kyooar/internal/subscription/middleware"
	subscriptionModels "kyooar/internal/subscription/models"
	"github.com/samber/do"
)

type Module struct {
	injector *do.Injector
}

func NewModule(i *do.Injector) *Module {
	return &Module{injector: i}
}

func (m *Module) RegisterRoutes(v1 *echo.Group) {
	// Get handlers and middleware from injector
	organizationHandler := do.MustInvoke[*handlers.OrganizationHandler](m.injector)
	productHandler := do.MustInvoke[*menuHandlers.ProductHandler](m.injector)
	qrCodeHandler := do.MustInvoke[*qrcodeHandlers.QRCodeHandler](m.injector)
	feedbackHandler := do.MustInvoke[*feedbackHandlers.FeedbackHandler](m.injector)
	questionnaireHandler := do.MustInvoke[*feedbackHandlers.QuestionnaireHandler](m.injector)
	questionHandler := do.MustInvoke[*feedbackHandlers.QuestionHandler](m.injector)
	middlewareProvider := do.MustInvoke[*sharedMiddleware.MiddlewareProvider](m.injector)
	subscriptionMW := do.MustInvoke[*subscriptionMiddleware.SubscriptionMiddleware](m.injector)
	
	// Organization routes
	organizations := v1.Group("/organizations")
	organizations.Use(middlewareProvider.AuthMiddleware())
	organizations.Use(middlewareProvider.TeamAwareMiddleware())
	
	// Apply usage tracking middleware only to organization creation
	organizations.POST("", organizationHandler.Create,
		subscriptionMW.CheckResourceLimit(subscriptionModels.ResourceTypeOrganization),
		subscriptionMW.TrackUsageAfterSuccess(),
	)
	
	organizations.GET("", organizationHandler.GetAll)
	organizations.GET("/:id", organizationHandler.GetByID)
	organizations.PUT("/:id", organizationHandler.Update)
	organizations.DELETE("/:id", organizationHandler.Delete)
	
	// Product routes under organizations (moved from menu module)
	organizations.GET("/:organizationId/products", productHandler.GetByOrganization)
	organizations.POST("/:organizationId/products", productHandler.Create)
	
	// QR Code routes under organizations (moved from qrcode module)
	organizations.POST("/:organizationId/qr-codes", qrCodeHandler.Generate)
	organizations.GET("/:organizationId/qr-codes", qrCodeHandler.GetByOrganization)
	
	// Feedback routes under organizations (moved from feedback module)
	organizations.GET("/:organizationId/feedback", feedbackHandler.GetByOrganization)
	organizations.GET("/:organizationId/analytics", feedbackHandler.GetStats)
	
	// Questionnaire routes under organizations (moved from feedback module)
	organizations.POST("/:organizationId/questionnaires", questionnaireHandler.CreateQuestionnaire)
	organizations.GET("/:organizationId/questionnaires", questionnaireHandler.ListQuestionnaires)
	organizations.GET("/:organizationId/questionnaires/:id", questionnaireHandler.GetQuestionnaire)
	organizations.PUT("/:organizationId/questionnaires/:id", questionnaireHandler.UpdateQuestionnaire)
	organizations.DELETE("/:organizationId/questionnaires/:id", questionnaireHandler.DeleteQuestionnaire)
	organizations.POST("/:organizationId/questionnaires/:id/questions", questionnaireHandler.AddQuestion)
	organizations.PUT("/:organizationId/questionnaires/:id/questions/:questionId", questionnaireHandler.UpdateQuestion)
	organizations.DELETE("/:organizationId/questionnaires/:id/questions/:questionId", questionnaireHandler.DeleteQuestion)
	organizations.POST("/:organizationId/questionnaires/:id/reorder", questionnaireHandler.ReorderQuestions)
	
	// Question routes under organizations (moved from feedback module)
	organizations.POST("/:organizationId/products/:productId/questions", questionHandler.CreateQuestion)
	organizations.GET("/:organizationId/products/:productId/questions", questionHandler.GetQuestionsByProduct)
	organizations.GET("/:organizationId/products/:productId/questions/:questionId", questionHandler.GetQuestion)
	organizations.PUT("/:organizationId/products/:productId/questions/:questionId", questionHandler.UpdateQuestion)
	organizations.DELETE("/:organizationId/products/:productId/questions/:questionId", questionHandler.DeleteQuestion)
	organizations.POST("/:organizationId/products/:productId/questions/reorder", questionHandler.ReorderQuestions)
	organizations.GET("/:organizationId/questions/products-with-questions", questionHandler.GetProductsWithQuestions)
	organizations.POST("/:organizationId/questions/batch", questionHandler.GetQuestionsByProducts)
	
	// AI question generation routes under organizations (moved from feedback module)
	organizations.POST("/:organizationId/products/:productId/ai/generate-questions", questionnaireHandler.GenerateQuestions)
	organizations.POST("/:organizationId/products/:productId/ai/generate-questionnaire", questionnaireHandler.GenerateAndSaveQuestionnaire)
}
