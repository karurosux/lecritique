package handlers

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	feedbackModels "github.com/lecritique/api/internal/feedback/models"
	feedbackRepos "github.com/lecritique/api/internal/feedback/repositories"
	feedbackServices "github.com/lecritique/api/internal/feedback/services"
	menuRepos "github.com/lecritique/api/internal/menu/repositories"
	qrcodeServices "github.com/lecritique/api/internal/qrcode/services"
	"github.com/lecritique/api/internal/shared/errors"
	"github.com/lecritique/api/internal/shared/logger"
	"github.com/lecritique/api/internal/shared/response"
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

// ValidateQRCode validates a QR code
// @Summary Validate QR code
// @Description Validate a QR code and return associated data
// @Tags public
// @Accept json
// @Produce json
// @Param code path string true "QR Code"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/public/qr/{code} [get]
func (h *PublicHandler) ValidateQRCode(c echo.Context) error {
	ctx := c.Request().Context()
	code := c.Param("code")
	if code == "" {
		return response.Error(c, errors.BadRequest("QR code parameter is required"))
	}

	qrCode, err := h.qrCodeService.GetByCode(ctx, code)
	if err != nil {
		return response.Error(c, errors.NotFound("QR code"))
	}

	return response.Success(c, qrCode)
}

// GetRestaurantMenu gets public restaurant menu
// @Summary Get restaurant menu
// @Description Get public menu for a restaurant
// @Tags public
// @Accept json
// @Produce json
// @Param id path string true "Restaurant ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/public/restaurant/{id}/menu [get]
func (h *PublicHandler) GetRestaurantMenu(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.Error(c, errors.ErrInvalidUUID)
	}

	// Implementation would get restaurant menu
	// For now, return placeholder
	return response.Success(c, map[string]interface{}{
		"restaurant_id": id,
		"message":       "Menu endpoint - to be implemented",
	})
}

// GetQuestionnaire gets questionnaire for a dish
// @Summary Get questionnaire
// @Description Get questionnaire for a specific dish
// @Tags public
// @Accept json
// @Produce json
// @Param restaurantId path string true "Restaurant ID"
// @Param dishId path string true "Dish ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/public/questionnaire/{restaurantId}/{dishId} [get]
func (h *PublicHandler) GetQuestionnaire(c echo.Context) error {
	restaurantIDStr := c.Param("restaurantId")
	dishIDStr := c.Param("dishId")

	restaurantID, err := uuid.Parse(restaurantIDStr)
	if err != nil {
		return response.Error(c, errors.BadRequest("Invalid restaurant ID format"))
	}

	dishID, err := uuid.Parse(dishIDStr)
	if err != nil {
		return response.Error(c, errors.BadRequest("Invalid dish ID format"))
	}

	// Implementation would get questionnaire
	// For now, return placeholder
	return response.Success(c, map[string]interface{}{
		"restaurant_id": restaurantID,
		"dish_id":       dishID,
		"message":       "Questionnaire endpoint - to be implemented",
	})
}

// SubmitFeedback submits customer feedback
// @Summary Submit feedback
// @Description Submit customer feedback for a dish
// @Tags public
// @Accept json
// @Produce json
// @Param feedback body feedbackModels.Feedback true "Feedback data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/public/feedback [post]
func (h *PublicHandler) SubmitFeedback(c echo.Context) error {
	ctx := c.Request().Context()
	var feedback feedbackModels.Feedback
	if err := c.Bind(&feedback); err != nil {
		return response.Error(c, errors.BadRequest("Invalid feedback data provided"))
	}

	if err := h.feedbackService.Submit(ctx, &feedback); err != nil {
		logger.Error("Failed to submit feedback", err, logrus.Fields{
			"feedback": feedback,
		})
		return response.Error(c, errors.Internal("Failed to process feedback submission"))
	}

	return response.Success(c, map[string]string{
		"message": "Thank you for your feedback!",
	})
}