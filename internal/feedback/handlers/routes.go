package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/lecritique/api/internal/feedback/repositories"
	"github.com/lecritique/api/internal/feedback/services"
	menuRepos "github.com/lecritique/api/internal/menu/repositories"
	qrcodeRepos "github.com/lecritique/api/internal/qrcode/repositories"
	qrcodeServices "github.com/lecritique/api/internal/qrcode/services"
	restaurantRepos "github.com/lecritique/api/internal/restaurant/repositories"
	"gorm.io/gorm"
)

func RegisterRoutes(public *echo.Group, db *gorm.DB) {
	// Initialize repositories
	feedbackRepo := repositories.NewFeedbackRepository(db)
	questionnaireRepo := repositories.NewQuestionnaireRepository(db)
	qrCodeRepo := qrcodeRepos.NewQRCodeRepository(db)
	restaurantRepo := restaurantRepos.NewRestaurantRepository(db)
	dishRepo := menuRepos.NewDishRepository(db)
	
	// Initialize services
	qrCodeService := qrcodeServices.NewQRCodeService(qrCodeRepo, restaurantRepo)
	feedbackService := services.NewFeedbackService(feedbackRepo, restaurantRepo, qrCodeRepo)
	
	// Initialize handler
	publicHandler := NewPublicHandler(qrCodeService, feedbackService, dishRepo, questionnaireRepo)
	
	// Public feedback routes (no auth required)
	public.GET("/qr/:code", publicHandler.ValidateQRCode)
	public.GET("/restaurant/:id/menu", publicHandler.GetRestaurantMenu)
	public.GET("/questionnaire/:restaurantId/:dishId", publicHandler.GetQuestionnaire)
	public.POST("/feedback", publicHandler.SubmitFeedback)
}