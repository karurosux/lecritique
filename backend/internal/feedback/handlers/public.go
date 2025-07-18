package handlers

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	feedbackModels "lecritique/internal/feedback/models"
	feedbackRepos "lecritique/internal/feedback/repositories"
	feedbackServices "lecritique/internal/feedback/services"
	menuRepos "lecritique/internal/menu/repositories"
	"lecritique/internal/shared/errors"
	"lecritique/internal/shared/logger"
	"lecritique/internal/shared/response"
	"lecritique/internal/shared/utils"
	"github.com/samber/do"
	"github.com/sirupsen/logrus"
)

type FeedbackPublicHandler struct {
	feedbackService   feedbackServices.FeedbackService
	dishRepo          menuRepos.DishRepository
	questionnaireRepo feedbackRepos.QuestionnaireRepository
	questionRepo      feedbackRepos.QuestionRepository
}

func NewFeedbackPublicHandler(i *do.Injector) (*FeedbackPublicHandler, error) {
	return &FeedbackPublicHandler{
		feedbackService:   do.MustInvoke[feedbackServices.FeedbackService](i),
		dishRepo:          do.MustInvoke[menuRepos.DishRepository](i),
		questionnaireRepo: do.MustInvoke[feedbackRepos.QuestionnaireRepository](i),
		questionRepo:      do.MustInvoke[feedbackRepos.QuestionRepository](i),
	}, nil
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
func (h *FeedbackPublicHandler) GetQuestionnaire(c echo.Context) error {
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
func (h *FeedbackPublicHandler) SubmitFeedback(c echo.Context) error {
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
func (h *FeedbackPublicHandler) GetDishQuestions(c echo.Context) error {
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
func (h *FeedbackPublicHandler) GetDishesWithQuestions(c echo.Context) error {
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