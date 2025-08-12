package organization

import (
	"github.com/labstack/echo/v4"
	"kyooar/internal/organization/handlers"
	feedbackcontroller "kyooar/internal/feedback/controller"
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
	organizationHandler := do.MustInvoke[*handlers.OrganizationHandler](m.injector)
	productHandler := do.MustInvoke[*menuHandlers.ProductHandler](m.injector)
	qrCodeHandler := do.MustInvoke[*qrcodeHandlers.QRCodeHandler](m.injector)
	feedbackController := do.MustInvoke[*feedbackcontroller.FeedbackController](m.injector)
	questionnaireController := do.MustInvoke[*feedbackcontroller.QuestionnaireController](m.injector)
	questionController := do.MustInvoke[*feedbackcontroller.QuestionController](m.injector)
	middlewareProvider := do.MustInvoke[*sharedMiddleware.MiddlewareProvider](m.injector)
	subscriptionMW := do.MustInvoke[*subscriptionMiddleware.SubscriptionMiddleware](m.injector)
	organizations := v1.Group("/organizations")
	organizations.Use(middlewareProvider.AuthMiddleware())
	organizations.Use(middlewareProvider.TeamAwareMiddleware())
	organizations.POST("", organizationHandler.Create,
		subscriptionMW.CheckResourceLimit(subscriptionModels.ResourceTypeOrganization),
		subscriptionMW.TrackUsageAfterSuccess(),
	)
	
	organizations.GET("", organizationHandler.GetAll)
	organizations.GET("/:id", organizationHandler.GetByID)
	organizations.PUT("/:id", organizationHandler.Update)
	organizations.DELETE("/:id", organizationHandler.Delete)
	organizations.GET("/:organizationId/products", productHandler.GetByOrganization)
	organizations.POST("/:organizationId/products", productHandler.Create)
	organizations.POST("/:organizationId/qr-codes", qrCodeHandler.Generate)
	organizations.GET("/:organizationId/qr-codes", qrCodeHandler.GetByOrganization)
	organizations.GET("/:organizationId/feedback", feedbackController.GetByOrganization)
	organizations.GET("/:organizationId/analytics", feedbackController.GetStats)
	organizations.POST("/:organizationId/questionnaires", questionnaireController.CreateQuestionnaire)
	organizations.GET("/:organizationId/questionnaires", questionnaireController.ListQuestionnaires)
	organizations.GET("/:organizationId/questionnaires/:id", questionnaireController.GetQuestionnaire)
	organizations.PUT("/:organizationId/questionnaires/:id", questionnaireController.UpdateQuestionnaire)
	organizations.DELETE("/:organizationId/questionnaires/:id", questionnaireController.DeleteQuestionnaire)
	organizations.POST("/:organizationId/questionnaires/:id/questions", questionnaireController.AddQuestion)
	organizations.PUT("/:organizationId/questionnaires/:id/questions/:questionId", questionnaireController.UpdateQuestion)
	organizations.DELETE("/:organizationId/questionnaires/:id/questions/:questionId", questionnaireController.DeleteQuestion)
	organizations.POST("/:organizationId/questionnaires/:id/reorder", questionnaireController.ReorderQuestions)
	organizations.POST("/:organizationId/products/:productId/questions", questionController.CreateQuestion)
	organizations.GET("/:organizationId/products/:productId/questions", questionController.GetQuestionsByProduct)
	organizations.GET("/:organizationId/products/:productId/questions/:questionId", questionController.GetQuestion)
	organizations.PUT("/:organizationId/products/:productId/questions/:questionId", questionController.UpdateQuestion)
	organizations.DELETE("/:organizationId/products/:productId/questions/:questionId", questionController.DeleteQuestion)
	organizations.POST("/:organizationId/products/:productId/questions/reorder", questionController.ReorderQuestions)
	organizations.GET("/:organizationId/questions/products-with-questions", questionController.GetProductsWithQuestions)
	organizations.POST("/:organizationId/questions/batch", questionController.GetQuestionsByProducts)
	organizations.POST("/:organizationId/products/:productId/ai/generate-questions", questionnaireController.GenerateQuestions)
	organizations.POST("/:organizationId/products/:productId/ai/generate-questionnaire", questionnaireController.GenerateAndSaveQuestionnaire)
}
