package restaurant

import (
	"github.com/labstack/echo/v4"
	"lecritique/internal/restaurant/handlers"
	feedbackHandlers "lecritique/internal/feedback/handlers"
	menuHandlers "lecritique/internal/menu/handlers"
	qrcodeHandlers "lecritique/internal/qrcode/handlers"
	sharedMiddleware "lecritique/internal/shared/middleware"
	subscriptionMiddleware "lecritique/internal/subscription/middleware"
	subscriptionModels "lecritique/internal/subscription/models"
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
	restaurantHandler := do.MustInvoke[*handlers.RestaurantHandler](m.injector)
	dishHandler := do.MustInvoke[*menuHandlers.DishHandler](m.injector)
	qrCodeHandler := do.MustInvoke[*qrcodeHandlers.QRCodeHandler](m.injector)
	feedbackHandler := do.MustInvoke[*feedbackHandlers.FeedbackHandler](m.injector)
	questionnaireHandler := do.MustInvoke[*feedbackHandlers.QuestionnaireHandler](m.injector)
	questionHandler := do.MustInvoke[*feedbackHandlers.QuestionHandler](m.injector)
	middlewareProvider := do.MustInvoke[*sharedMiddleware.MiddlewareProvider](m.injector)
	subscriptionMW := do.MustInvoke[*subscriptionMiddleware.SubscriptionMiddleware](m.injector)
	
	// Restaurant routes
	restaurants := v1.Group("/restaurants")
	restaurants.Use(middlewareProvider.AuthMiddleware())
	restaurants.Use(middlewareProvider.TeamAwareMiddleware())
	
	// Apply usage tracking middleware only to restaurant creation
	restaurants.POST("", restaurantHandler.Create,
		subscriptionMW.CheckResourceLimit(subscriptionModels.ResourceTypeRestaurant),
		subscriptionMW.TrackUsageAfterSuccess(),
	)
	
	restaurants.GET("", restaurantHandler.GetAll)
	restaurants.GET("/:id", restaurantHandler.GetByID)
	restaurants.PUT("/:id", restaurantHandler.Update)
	restaurants.DELETE("/:id", restaurantHandler.Delete)
	
	// Dish routes under restaurants (moved from menu module)
	restaurants.GET("/:restaurantId/dishes", dishHandler.GetByRestaurant)
	restaurants.POST("/:restaurantId/dishes", dishHandler.Create)
	
	// QR Code routes under restaurants (moved from qrcode module)
	restaurants.POST("/:restaurantId/qr-codes", qrCodeHandler.Generate)
	restaurants.GET("/:restaurantId/qr-codes", qrCodeHandler.GetByRestaurant)
	
	// Feedback routes under restaurants (moved from feedback module)
	restaurants.GET("/:restaurantId/feedback", feedbackHandler.GetByRestaurant)
	restaurants.GET("/:restaurantId/analytics", feedbackHandler.GetStats)
	
	// Questionnaire routes under restaurants (moved from feedback module)
	restaurants.POST("/:restaurantId/questionnaires", questionnaireHandler.CreateQuestionnaire)
	restaurants.GET("/:restaurantId/questionnaires", questionnaireHandler.ListQuestionnaires)
	restaurants.GET("/:restaurantId/questionnaires/:id", questionnaireHandler.GetQuestionnaire)
	restaurants.PUT("/:restaurantId/questionnaires/:id", questionnaireHandler.UpdateQuestionnaire)
	restaurants.DELETE("/:restaurantId/questionnaires/:id", questionnaireHandler.DeleteQuestionnaire)
	restaurants.POST("/:restaurantId/questionnaires/:id/questions", questionnaireHandler.AddQuestion)
	restaurants.PUT("/:restaurantId/questionnaires/:id/questions/:questionId", questionnaireHandler.UpdateQuestion)
	restaurants.DELETE("/:restaurantId/questionnaires/:id/questions/:questionId", questionnaireHandler.DeleteQuestion)
	restaurants.POST("/:restaurantId/questionnaires/:id/reorder", questionnaireHandler.ReorderQuestions)
	
	// Question routes under restaurants (moved from feedback module)
	restaurants.POST("/:restaurantId/dishes/:dishId/questions", questionHandler.CreateQuestion)
	restaurants.GET("/:restaurantId/dishes/:dishId/questions", questionHandler.GetQuestionsByDish)
	restaurants.GET("/:restaurantId/dishes/:dishId/questions/:questionId", questionHandler.GetQuestion)
	restaurants.PUT("/:restaurantId/dishes/:dishId/questions/:questionId", questionHandler.UpdateQuestion)
	restaurants.DELETE("/:restaurantId/dishes/:dishId/questions/:questionId", questionHandler.DeleteQuestion)
	restaurants.POST("/:restaurantId/dishes/:dishId/questions/reorder", questionHandler.ReorderQuestions)
	restaurants.GET("/:restaurantId/questions/dishes-with-questions", questionHandler.GetDishesWithQuestions)
	
	// AI question generation routes under restaurants (moved from feedback module)
	restaurants.POST("/:restaurantId/dishes/:dishId/ai/generate-questions", questionnaireHandler.GenerateQuestions)
	restaurants.POST("/:restaurantId/dishes/:dishId/ai/generate-questionnaire", questionnaireHandler.GenerateAndSaveQuestionnaire)
}