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
	"github.com/lecritique/api/internal/shared/utils"
	"github.com/sirupsen/logrus"
)

type PublicHandler struct {
	qrCodeService     qrcodeServices.QRCodeService
	feedbackService   feedbackServices.FeedbackService
	dishRepo          menuRepos.DishRepository
	questionnaireRepo feedbackRepos.QuestionnaireRepository
	questionRepo      feedbackRepos.QuestionRepository
}

func NewPublicHandler(
	qrCodeService qrcodeServices.QRCodeService,
	feedbackService feedbackServices.FeedbackService,
	dishRepo menuRepos.DishRepository,
	questionnaireRepo feedbackRepos.QuestionnaireRepository,
	questionRepo feedbackRepos.QuestionRepository,
) *PublicHandler {
	return &PublicHandler{
		qrCodeService:     qrCodeService,
		feedbackService:   feedbackService,
		dishRepo:          dishRepo,
		questionnaireRepo: questionnaireRepo,
		questionRepo:      questionRepo,
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

	// Record the scan event for analytics
	if err := h.qrCodeService.RecordScan(ctx, code); err != nil {
		logger.Error("Failed to record QR scan", err, logrus.Fields{
			"qr_code_id": qrCode.ID,
			"code":       code,
		})
		// Don't fail the request if scan recording fails
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

	// Extract device information from request
	deviceInfo := utils.ExtractDeviceInfo(c.Request())
	feedback.DeviceInfo = feedbackModels.DeviceInfo{
		UserAgent: deviceInfo.UserAgent,
		IP:        deviceInfo.IP,
		Platform:  deviceInfo.Platform,
		Browser:   deviceInfo.Browser,
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

// GetDishQuestions gets questions for a specific dish (public endpoint)
// @Summary Get questions for a dish
// @Description Get all feedback questions for a specific dish (public access for customer feedback)
// @Tags public
// @Produce json
// @Param restaurantId path string true "Restaurant ID"
// @Param dishId path string true "Dish ID"
// @Success 200 {object} map[string]interface{} "Questions retrieved successfully"
// @Failure 400 {object} response.Response "Invalid request"
// @Failure 404 {object} response.Response "Dish not found"
// @Failure 500 {object} response.Response "Server error"
// @Router /api/v1/public/restaurant/{restaurantId}/dishes/{dishId}/questions [get]
func (h *PublicHandler) GetDishQuestions(c echo.Context) error {
	ctx := c.Request().Context()

	dishID, err := uuid.Parse(c.Param("dishId"))
	if err != nil {
		return response.Error(c, errors.BadRequest("Invalid dish ID"))
	}

	// Verify dish exists
	dish, err := h.dishRepo.FindByID(ctx, dishID)
	if err != nil {
		return response.Error(c, errors.NotFound("Dish not found"))
	}

	// Get questions for this dish
	questions, err := h.questionRepo.GetQuestionsByDishID(ctx, dishID)
	if err != nil {
		logger.Error("Failed to get questions for dish", err, logrus.Fields{
			"dish_id": dishID,
		})
		return response.Error(c, errors.Internal("Failed to get questions"))
	}

	return response.Success(c, map[string]interface{}{
		"dish":      dish,
		"questions": questions,
	})
}

// GetDishesWithQuestions gets all dishes that have questions for a restaurant (public endpoint)
// @Summary Get dishes with questions
// @Description Get all dishes that have feedback questions for a restaurant (public access for QR code scans)
// @Tags public
// @Produce json
// @Param restaurantId path string true "Restaurant ID"
// @Success 200 {object} map[string]interface{} "Dishes with questions retrieved successfully"
// @Failure 400 {object} response.Response "Invalid request"
// @Failure 404 {object} response.Response "Restaurant not found"
// @Failure 500 {object} response.Response "Server error"
// @Router /api/v1/public/restaurant/{restaurantId}/questions/dishes-with-questions [get]
func (h *PublicHandler) GetDishesWithQuestions(c echo.Context) error {
	ctx := c.Request().Context()

	restaurantID, err := uuid.Parse(c.Param("restaurantId"))
	if err != nil {
		return response.Error(c, errors.BadRequest("Invalid restaurant ID"))
	}

	// Get dishes with questions for this restaurant
	dishIDs, err := h.questionRepo.GetDishesWithQuestions(ctx, restaurantID)
	if err != nil {
		logger.Error("Failed to get dishes with questions", err, logrus.Fields{
			"restaurant_id": restaurantID,
		})
		return response.Error(c, errors.Internal("Failed to get dishes with questions"))
	}

	// Get full dish details for each dish ID
	dishes := make([]interface{}, 0, len(dishIDs))
	for _, dishID := range dishIDs {
		dish, err := h.dishRepo.FindByID(ctx, dishID)
		if err != nil {
			logger.Warn("Failed to get dish details", logrus.Fields{
				"dish_id": dishID,
				"error":   err,
			})
			continue
		}
		dishes = append(dishes, dish)
	}

	return response.Success(c, dishes)
}