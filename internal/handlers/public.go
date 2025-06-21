package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lecritique/api/internal/models"
	"github.com/lecritique/api/internal/repositories"
	"github.com/lecritique/api/internal/services"
	"github.com/lecritique/api/pkg/errors"
	"github.com/lecritique/api/pkg/response"
	"github.com/lecritique/api/pkg/validator"
)

type PublicHandler struct {
	qrCodeService services.QRCodeService
	feedbackService services.FeedbackService
	dishRepo repositories.DishRepository
	questionnaireRepo repositories.QuestionnaireRepository
	validator *validator.Validator
}

func NewPublicHandler(qrCodeService services.QRCodeService, feedbackService services.FeedbackService, dishRepo repositories.DishRepository, questionnaireRepo repositories.QuestionnaireRepository) *PublicHandler {
	return &PublicHandler{
		qrCodeService: qrCodeService,
		feedbackService: feedbackService,
		dishRepo: dishRepo,
		questionnaireRepo: questionnaireRepo,
		validator: validator.New(),
	}
}

type QRCodeValidationResponse struct {
	Restaurant interface{} `json:"restaurant"`
	QRCode     interface{} `json:"qr_code"`
}

func (h *PublicHandler) ValidateQRCode(c echo.Context) error {
	code := c.Param("code")
	if code == "" {
		return response.Error(c, errors.ErrBadRequest)
	}

	qrCode, err := h.qrCodeService.GetByCode(code)
	if err != nil {
		return response.Error(c, errors.ErrNotFound)
	}

	// Record scan
	_ = h.qrCodeService.RecordScan(code)

	return response.Success(c, QRCodeValidationResponse{
		Restaurant: qrCode.Restaurant,
		QRCode: map[string]interface{}{
			"id":    qrCode.ID,
			"label": qrCode.Label,
			"type":  qrCode.Type,
		},
	})
}

func (h *PublicHandler) GetRestaurantMenu(c echo.Context) error {
	restaurantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	dishes, err := h.dishRepo.FindByRestaurantID(restaurantID)
	if err != nil {
		return response.Error(c, err)
	}

	// Filter only available dishes
	availableDishes := make([]models.Dish, 0)
	for _, dish := range dishes {
		if dish.IsAvailable && dish.IsActive {
			availableDishes = append(availableDishes, dish)
		}
	}

	return response.Success(c, availableDishes)
}

func (h *PublicHandler) GetQuestionnaire(c echo.Context) error {
	dishID, err := uuid.Parse(c.Param("dishId"))
	if err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	questionnaire, err := h.questionnaireRepo.FindByDishID(dishID)
	if err != nil {
		// Return default questionnaire if none exists
		return response.Success(c, map[string]interface{}{
			"questions": getDefaultQuestions(),
		})
	}

	return response.Success(c, questionnaire)
}

type SubmitFeedbackRequest struct {
	QRCodeID      uuid.UUID           `json:"qr_code_id" validate:"required"`
	DishID        uuid.UUID           `json:"dish_id" validate:"required"`
	CustomerName  string              `json:"customer_name"`
	CustomerEmail string              `json:"customer_email" validate:"omitempty,email"`
	CustomerPhone string              `json:"customer_phone"`
	OverallRating int                 `json:"overall_rating" validate:"required,min=1,max=5"`
	Responses     []models.Response   `json:"responses"`
}

func (h *PublicHandler) SubmitFeedback(c echo.Context) error {
	var req SubmitFeedbackRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	if err := h.validator.Validate(req); err != nil {
		return response.Error(c, errors.NewWithDetails("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, h.validator.FormatErrors(err)))
	}

	// Get device info
	deviceInfo := models.DeviceInfo{
		UserAgent: c.Request().UserAgent(),
		IP:        c.RealIP(),
	}

	feedback := &models.Feedback{
		QRCodeID:      req.QRCodeID,
		DishID:        req.DishID,
		CustomerName:  req.CustomerName,
		CustomerEmail: req.CustomerEmail,
		CustomerPhone: req.CustomerPhone,
		OverallRating: req.OverallRating,
		Responses:     req.Responses,
		DeviceInfo:    deviceInfo,
		IsComplete:    true,
	}

	if err := h.feedbackService.Submit(feedback); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]interface{}{
		"id": feedback.ID,
		"message": "Thank you for your feedback!",
	})
}

func getDefaultQuestions() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"text": "How would you rate your overall experience?",
			"type": "rating",
			"is_required": true,
		},
		{
			"text": "How would you rate the taste of your dish?",
			"type": "rating",
			"is_required": true,
		},
		{
			"text": "Was your food served at the right temperature?",
			"type": "single_choice",
			"options": []string{"Too cold", "Just right", "Too hot"},
			"is_required": false,
		},
		{
			"text": "Any additional comments or suggestions?",
			"type": "text",
			"is_required": false,
		},
	}
}
