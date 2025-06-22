package handlers

import (
	"github.com/labstack/echo/v4"
	authServices "github.com/lecritique/api/internal/auth/services"
	"github.com/lecritique/api/internal/feedback/repositories"
	"github.com/lecritique/api/internal/feedback/services"
	menuRepos "github.com/lecritique/api/internal/menu/repositories"
	qrcodeRepos "github.com/lecritique/api/internal/qrcode/repositories"
	qrcodeServices "github.com/lecritique/api/internal/qrcode/services"
	restaurantRepos "github.com/lecritique/api/internal/restaurant/repositories"
	"github.com/lecritique/api/internal/shared/middleware"
	"gorm.io/gorm"
)

func RegisterRoutes(public *echo.Group, protected *echo.Group, db *gorm.DB, authService authServices.AuthService) {
	// Initialize repositories
	feedbackRepo := repositories.NewFeedbackRepository(db)
	questionnaireRepo := repositories.NewQuestionnaireRepository(db)
	qrCodeRepo := qrcodeRepos.NewQRCodeRepository(db)
	restaurantRepo := restaurantRepos.NewRestaurantRepository(db)
	dishRepo := menuRepos.NewDishRepository(db)
	
	// Initialize services
	qrCodeService := qrcodeServices.NewQRCodeService(qrCodeRepo, restaurantRepo)
	feedbackService := services.NewFeedbackService(feedbackRepo, restaurantRepo, qrCodeRepo)
	
	// Initialize handlers
	publicHandler := NewPublicHandler(qrCodeService, feedbackService, dishRepo, questionnaireRepo)
	feedbackHandler := NewFeedbackHandler(feedbackService)
	
	// Public feedback routes (no auth required)
	public.GET("/qr/:code", publicHandler.ValidateQRCode)
	public.GET("/restaurant/:id/menu", publicHandler.GetRestaurantMenu)
	public.GET("/questionnaire/:restaurantId/:dishId", publicHandler.GetQuestionnaire)
	public.POST("/feedback", publicHandler.SubmitFeedback)
	
	// Protected feedback routes (requires auth)
	restaurants := protected.Group("/restaurants")
	restaurants.Use(middleware.JWTAuth(authService))
	restaurants.GET("/:restaurantId/feedback", feedbackHandler.GetByRestaurant)
	restaurants.GET("/:restaurantId/analytics", feedbackHandler.GetStats)
}