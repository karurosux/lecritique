package handlers

import (
	"github.com/labstack/echo/v4"
	authServices "github.com/lecritique/api/internal/auth/services"
	"github.com/lecritique/api/internal/feedback/repositories"
	"github.com/lecritique/api/internal/feedback/services"
	menuRepos "github.com/lecritique/api/internal/menu/repositories"
	menuServices "github.com/lecritique/api/internal/menu/services"
	qrcodeRepos "github.com/lecritique/api/internal/qrcode/repositories"
	qrcodeServices "github.com/lecritique/api/internal/qrcode/services"
	restaurantRepos "github.com/lecritique/api/internal/restaurant/repositories"
	"github.com/lecritique/api/internal/shared/config"
	"github.com/lecritique/api/internal/shared/middleware"
	"gorm.io/gorm"
)

func RegisterRoutes(public *echo.Group, protected *echo.Group, db *gorm.DB, authService authServices.AuthService, cfg *config.Config) {
	// Initialize repositories
	feedbackRepo := repositories.NewFeedbackRepository(db)
	questionnaireRepo := repositories.NewQuestionnaireRepository(db)
	questionRepo := repositories.NewQuestionRepository(db)
	qrCodeRepo := qrcodeRepos.NewQRCodeRepository(db)
	restaurantRepo := restaurantRepos.NewRestaurantRepository(db)
	dishRepo := menuRepos.NewDishRepository(db)
	
	// Initialize services
	qrCodeService := qrcodeServices.NewQRCodeService(qrCodeRepo, restaurantRepo)
	feedbackService := services.NewFeedbackService(feedbackRepo, restaurantRepo, qrCodeRepo)
	questionnaireService, _ := services.NewQuestionnaireService(db, cfg)
	questionService := services.NewQuestionService(questionRepo, dishRepo, restaurantRepo)
	dishService := menuServices.NewDishService(dishRepo, restaurantRepo)
	
	// Initialize handlers
	publicHandler := NewPublicHandler(qrCodeService, feedbackService, dishRepo, questionnaireRepo, questionRepo)
	feedbackHandler := NewFeedbackHandler(feedbackService)
	questionnaireHandler := NewQuestionnaireHandler(questionnaireService, dishService)
	questionHandler := NewQuestionHandler(questionService)
	
	// Public feedback routes (no auth required)
	public.GET("/qr/:code", publicHandler.ValidateQRCode)
	public.GET("/restaurant/:id/menu", publicHandler.GetRestaurantMenu)
	public.GET("/questionnaire/:restaurantId/:dishId", publicHandler.GetQuestionnaire)
	public.GET("/restaurant/:restaurantId/dishes/:dishId/questions", publicHandler.GetDishQuestions)
	public.GET("/restaurant/:restaurantId/questions/dishes-with-questions", publicHandler.GetDishesWithQuestions)
	public.POST("/feedback", publicHandler.SubmitFeedback)
	
	// Protected feedback routes (requires auth)
	restaurants := protected.Group("/restaurants")
	restaurants.Use(middleware.JWTAuth(authService))
	restaurants.Use(middleware.TeamAware(db)) // Add team-aware middleware
	restaurants.GET("/:restaurantId/feedback", feedbackHandler.GetByRestaurant)
	restaurants.GET("/:restaurantId/analytics", feedbackHandler.GetStats)
	
	// Protected questionnaire routes
	questionnaires := protected.Group("/restaurants/:restaurantId/questionnaires")
	questionnaires.Use(middleware.JWTAuth(authService))
	questionnaires.Use(middleware.TeamAware(db)) // Add team-aware middleware
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
	dishRoutes := protected.Group("/restaurants/:restaurantId/dishes/:dishId")
	dishRoutes.Use(middleware.JWTAuth(authService))
	dishRoutes.Use(middleware.TeamAware(db)) // Add team-aware middleware
	dishRoutes.POST("/questions", questionHandler.CreateQuestion)
	dishRoutes.GET("/questions", questionHandler.GetQuestionsByDish)
	dishRoutes.GET("/questions/:questionId", questionHandler.GetQuestion)
	dishRoutes.PUT("/questions/:questionId", questionHandler.UpdateQuestion)
	dishRoutes.DELETE("/questions/:questionId", questionHandler.DeleteQuestion)
	dishRoutes.POST("/questions/reorder", questionHandler.ReorderQuestions)

	// Bulk question routes
	restaurantQuestionRoutes := protected.Group("/restaurants/:restaurantId")
	restaurantQuestionRoutes.Use(middleware.JWTAuth(authService))
	restaurantQuestionRoutes.Use(middleware.TeamAware(db)) // Add team-aware middleware
	restaurantQuestionRoutes.GET("/questions/dishes-with-questions", questionHandler.GetDishesWithQuestions)

	// AI question generation routes
	aiRoutes := protected.Group("/ai")
	aiRoutes.Use(middleware.JWTAuth(authService))
	aiRoutes.Use(middleware.TeamAware(db)) // Add team-aware middleware
	aiRoutes.POST("/generate-questions/:dishId", questionnaireHandler.GenerateQuestions)
	aiRoutes.POST("/generate-questionnaire/:dishId", questionnaireHandler.GenerateAndSaveQuestionnaire)
}