package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	feedbackModels "github.com/lecritique/api/internal/feedback/models"
	feedbackRepos "github.com/lecritique/api/internal/feedback/repositories"
	feedbackServices "github.com/lecritique/api/internal/feedback/services"
	menuRepos "github.com/lecritique/api/internal/menu/repositories"
	qrcodeServices "github.com/lecritique/api/internal/qrcode/services"
	"github.com/lecritique/api/internal/shared/logger"
	"github.com/sirupsen/logrus"
)

type PublicHandler struct {
	qrCodeService     qrcodeServices.QRCodeService
	feedbackService   feedbackServices.FeedbackService
	dishRepo          menuRepos.DishRepository
	questionnaireRepo feedbackRepos.QuestionnaireRepository
}

func NewPublicHandler(
	qrCodeService qrcodeServices.QRCodeService,
	feedbackService feedbackServices.FeedbackService,
	dishRepo menuRepos.DishRepository,
	questionnaireRepo feedbackRepos.QuestionnaireRepository,
) *PublicHandler {
	return &PublicHandler{
		qrCodeService:     qrCodeService,
		feedbackService:   feedbackService,
		dishRepo:          dishRepo,
		questionnaireRepo: questionnaireRepo,
	}
}

func (h *PublicHandler) ValidateQRCode(c echo.Context) error {
	code := c.Param("code")
	if code == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	qrCode, err := h.qrCodeService.GetByCode(code)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "QR code not found")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"success": true, "data": qrCode})
}

func (h *PublicHandler) GetRestaurantMenu(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	// Implementation would get restaurant menu
	// For now, return placeholder
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"restaurant_id": id,
			"message":       "Menu endpoint - to be implemented",
		},
	})
}

func (h *PublicHandler) GetQuestionnaire(c echo.Context) error {
	restaurantIDStr := c.Param("restaurantId")
	dishIDStr := c.Param("dishId")

	restaurantID, err := uuid.Parse(restaurantIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	dishID, err := uuid.Parse(dishIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	// Implementation would get questionnaire
	// For now, return placeholder
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"restaurant_id": restaurantID,
			"dish_id":       dishID,
			"message":       "Questionnaire endpoint - to be implemented",
		},
	})
}

func (h *PublicHandler) SubmitFeedback(c echo.Context) error {
	var feedback feedbackModels.Feedback
	if err := c.Bind(&feedback); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	if err := h.feedbackService.Submit(&feedback); err != nil {
		logger.Error("Failed to submit feedback", err, logrus.Fields{
			"feedback": feedback,
		})
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to submit feedback")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]string{
			"message": "Feedback submitted successfully",
		},
	})
}