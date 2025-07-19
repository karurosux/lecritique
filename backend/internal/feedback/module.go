package feedback

import (
	"github.com/labstack/echo/v4"
	"lecritique/internal/feedback/handlers"
	sharedMiddleware "lecritique/internal/shared/middleware"
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
	publicHandler := do.MustInvoke[*handlers.FeedbackPublicHandler](m.injector)
	feedbackHandler := do.MustInvoke[*handlers.FeedbackHandler](m.injector)
	questionnaireHandler := do.MustInvoke[*handlers.QuestionnaireHandler](m.injector)
	questionHandler := do.MustInvoke[*handlers.QuestionHandler](m.injector)
	
	// Get middleware provider
	middlewareProvider := do.MustInvoke[*sharedMiddleware.MiddlewareProvider](m.injector)
	
	// Public feedback routes (no auth required)
	v1.GET("/questionnaire/:restaurantId/:dishId", publicHandler.GetQuestionnaire)
	v1.GET("/restaurant/:restaurantId/dishes/:dishId/questions", publicHandler.GetDishQuestions)
	v1.GET("/restaurant/:restaurantId/questions/dishes-with-questions", publicHandler.GetDishesWithQuestions)
	v1.POST("/feedback", publicHandler.SubmitFeedback)
	
	// Protected feedback routes (requires auth)
	restaurants := v1.Group("/restaurants")
	restaurants.Use(middlewareProvider.AuthMiddleware())
	restaurants.Use(middlewareProvider.TeamAwareMiddleware())
	restaurants.GET("/:restaurantId/feedback", feedbackHandler.GetByRestaurant)
	restaurants.GET("/:restaurantId/analytics", feedbackHandler.GetStats)
	
	// Protected questionnaire routes
	questionnaires := v1.Group("/restaurants/:restaurantId/questionnaires")
	questionnaires.Use(middlewareProvider.AuthMiddleware())
	questionnaires.Use(middlewareProvider.TeamAwareMiddleware())
	questionnaires.POST("", questionnaireHandler.CreateQuestionnaire)
	questionnaires.GET("", questionnaireHandler.ListQuestionnaires)
	questionnaires.GET("/:id", questionnaireHandler.GetQuestionnaire)
	questionnaires.PUT("/:id", questionnaireHandler.UpdateQuestionnaire)
	questionnaires.DELETE("/:id", questionnaireHandler.DeleteQuestionnaire)
	questionnaires.POST("/:id/questions", questionnaireHandler.AddQuestion)
	questionnaires.PUT("/:id/questions/:questionId", questionnaireHandler.UpdateQuestion)
	questionnaires.DELETE("/:id/questions/:questionId", questionnaireHandler.DeleteQuestion)
	questionnaires.POST("/:id/reorder", questionnaireHandler.ReorderQuestions)
	
	// Question routes (new simplified structure)
	dishRoutes := v1.Group("/restaurants/:restaurantId/dishes/:dishId")
	dishRoutes.Use(middlewareProvider.AuthMiddleware())
	dishRoutes.Use(middlewareProvider.TeamAwareMiddleware())
	dishRoutes.POST("/questions", questionHandler.CreateQuestion)
	dishRoutes.GET("/questions", questionHandler.GetQuestionsByDish)
	dishRoutes.GET("/questions/:questionId", questionHandler.GetQuestion)
	dishRoutes.PUT("/questions/:questionId", questionHandler.UpdateQuestion)
	dishRoutes.DELETE("/questions/:questionId", questionHandler.DeleteQuestion)
	dishRoutes.POST("/questions/reorder", questionHandler.ReorderQuestions)

	// Bulk question routes
	restaurantQuestionRoutes := v1.Group("/restaurants/:restaurantId")
	restaurantQuestionRoutes.Use(middlewareProvider.AuthMiddleware())
	restaurantQuestionRoutes.Use(middlewareProvider.TeamAwareMiddleware())
	restaurantQuestionRoutes.GET("/questions/dishes-with-questions", questionHandler.GetDishesWithQuestions)

	// AI question generation routes
	aiRoutes := v1.Group("/ai")
	aiRoutes.Use(middlewareProvider.AuthMiddleware())
	aiRoutes.Use(middlewareProvider.TeamAwareMiddleware())
	aiRoutes.POST("/generate-questions/:dishId", questionnaireHandler.GenerateQuestions)
	aiRoutes.POST("/generate-questionnaire/:dishId", questionnaireHandler.GenerateAndSaveQuestionnaire)
}