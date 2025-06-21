package handlers

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lecritique/api/internal/feedback/services"
	"github.com/lecritique/api/internal/qrcode/services"
	"github.com/lecritique/api/internal/shared/errors"
	"github.com/lecritique/api/internal/shared/models"
	"github.com/lecritique/api/internal/shared/repositories"
	"github.com/lecritique/api/internal/shared/response"
	"github.com/lecritique/api/internal/shared/validator"
)

type PublicHandler struct {
	qrCodeService     qrcodeServices.QRCodeService
	feedbackService   feedbackServices.FeedbackService
	dishRepo          repositories.Repository[models.Dish]
	questionnaireRepo repositories.Repository[models.Questionnaire]
	validator         *validator.Validator
}

func NewPublicHandler(
	qrCodeService qrcodeServices.QRCodeService,
	feedbackService feedbackServices.FeedbackService,
	dishRepo repositories.Repository[models.Dish],
	questionnaireRepo repositories.Repository[models.Questionnaire],
) *PublicHandler {
	return &PublicHandler{
		qrCodeService:     qrCodeService,
		feedbackService:   feedbackService,
		dishRepo:          dishRepo,
		questionnaireRepo: questionnaireRepo,
		validator:         validator.New(),
	}
}

func (h *PublicHandler) ValidateQRCode(c echo.Context) error {
	code := c.Param("code")
	if code == "" {
		return response.Error(c, errors.ErrBadRequest)
	}

	qrCode, err := h.qrCodeService.ValidateCode(code)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, qrCode)
}

func (h *PublicHandler) GetRestaurantMenu(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	// Implementation would get restaurant menu
	// For now, return placeholder
	return response.Success(c, map[string]interface{}{
		"restaurant_id": id,
		"message":       "Menu endpoint - to be implemented",
	})
}

func (h *PublicHandler) GetQuestionnaire(c echo.Context) error {
	restaurantIDStr := c.Param("restaurantId")
	dishIDStr := c.Param("dishId")

	restaurantID, err := uuid.Parse(restaurantIDStr)
	if err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	dishID, err := uuid.Parse(dishIDStr)
	if err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	// Implementation would get questionnaire
	// For now, return placeholder
	return response.Success(c, map[string]interface{}{
		"restaurant_id": restaurantID,
		"dish_id":       dishID,
		"message":       "Questionnaire endpoint - to be implemented",
	})
}

func (h *PublicHandler) SubmitFeedback(c echo.Context) error {
	var feedback models.Feedback
	if err := c.Bind(&feedback); err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	if err := h.validator.Validate(feedback); err != nil {
		return response.Error(c, errors.NewWithDetails("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, h.validator.FormatErrors(err)))
	}

	if err := h.feedbackService.Submit(&feedback); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]string{
		"message": "Feedback submitted successfully",
	})
}