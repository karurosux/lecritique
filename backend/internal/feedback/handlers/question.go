package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"lecritique/internal/feedback/models"
	"lecritique/internal/feedback/services"
	"lecritique/internal/shared/middleware"
	"github.com/samber/do"
)

type QuestionHandler struct {
	questionService services.QuestionService
}

func NewQuestionHandler(i *do.Injector) (*QuestionHandler, error) {
	return &QuestionHandler{
		questionService: do.MustInvoke[services.QuestionService](i),
	}, nil
}

// CreateQuestion creates a new question for a dish
// @Summary Add a question to a dish
// @Description Add a new feedback question to a specific dish
// @Tags questions
// @Accept json
// @Produce json
// @Param restaurantId path string true "Restaurant ID"
// @Param dishId path string true "Dish ID"
// @Param question body models.CreateQuestionRequest true "Question data"
// @Success 201 {object} map[string]interface{} "Question created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 404 {object} map[string]interface{} "Dish not found"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /api/v1/restaurants/{restaurantId}/dishes/{dishId}/questions [post]
// @Security BearerAuth
func (h *QuestionHandler) CreateQuestion(c echo.Context) error {
	accountID := middleware.GetResourceAccountID(c)

	dishID, err := uuid.Parse(c.Param("dishId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid dish ID")
	}

	var request models.CreateQuestionRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	question, err := h.questionService.CreateQuestion(c.Request().Context(), accountID, dishID, &request)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Question created successfully",
		"data":    question,
	})
}

// GetQuestionsByDish gets all questions for a specific dish
// @Summary Get questions for a dish
// @Description Get all feedback questions for a specific dish
// @Tags questions
// @Produce json
// @Param restaurantId path string true "Restaurant ID"
// @Param dishId path string true "Dish ID"
// @Success 200 {object} map[string]interface{} "Questions retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 404 {object} map[string]interface{} "Dish not found"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /api/v1/restaurants/{restaurantId}/dishes/{dishId}/questions [get]
// @Security BearerAuth
func (h *QuestionHandler) GetQuestionsByDish(c echo.Context) error {
	accountID := middleware.GetResourceAccountID(c)

	dishID, err := uuid.Parse(c.Param("dishId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid dish ID")
	}

	questions, err := h.questionService.GetQuestionsByDish(c.Request().Context(), accountID, dishID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    questions,
	})
}

// GetQuestion gets a specific question
// @Summary Get a specific question
// @Description Get details of a specific question
// @Tags questions
// @Produce json
// @Param restaurantId path string true "Restaurant ID"
// @Param dishId path string true "Dish ID"
// @Param questionId path string true "Question ID"
// @Success 200 {object} map[string]interface{} "Question retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 404 {object} map[string]interface{} "Question not found"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /api/v1/restaurants/{restaurantId}/dishes/{dishId}/questions/{questionId} [get]
// @Security BearerAuth
func (h *QuestionHandler) GetQuestion(c echo.Context) error {
	accountID := middleware.GetResourceAccountID(c)

	questionID, err := uuid.Parse(c.Param("questionId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid question ID")
	}

	question, err := h.questionService.GetQuestion(c.Request().Context(), accountID, questionID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    question,
	})
}

// UpdateQuestion updates an existing question
// @Summary Update a question
// @Description Update an existing question for a dish
// @Tags questions
// @Accept json
// @Produce json
// @Param restaurantId path string true "Restaurant ID"
// @Param dishId path string true "Dish ID"
// @Param questionId path string true "Question ID"
// @Param question body models.UpdateQuestionRequest true "Question data"
// @Success 200 {object} map[string]interface{} "Question updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 404 {object} map[string]interface{} "Question not found"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /api/v1/restaurants/{restaurantId}/dishes/{dishId}/questions/{questionId} [put]
// @Security BearerAuth
func (h *QuestionHandler) UpdateQuestion(c echo.Context) error {
	accountID := middleware.GetResourceAccountID(c)

	questionID, err := uuid.Parse(c.Param("questionId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid question ID")
	}

	var request models.UpdateQuestionRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	question, err := h.questionService.UpdateQuestion(c.Request().Context(), accountID, questionID, &request)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Question updated successfully",
		"data":    question,
	})
}

// DeleteQuestion deletes a question
// @Summary Delete a question
// @Description Delete a feedback question from a dish
// @Tags questions
// @Produce json
// @Param restaurantId path string true "Restaurant ID"
// @Param dishId path string true "Dish ID"
// @Param questionId path string true "Question ID"
// @Success 200 {object} map[string]interface{} "Question deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 404 {object} map[string]interface{} "Question not found"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /api/v1/restaurants/{restaurantId}/dishes/{dishId}/questions/{questionId} [delete]
// @Security BearerAuth
func (h *QuestionHandler) DeleteQuestion(c echo.Context) error {
	accountID := middleware.GetResourceAccountID(c)

	questionID, err := uuid.Parse(c.Param("questionId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid question ID")
	}

	if err := h.questionService.DeleteQuestion(c.Request().Context(), accountID, questionID); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Question deleted successfully",
	})
}

// ReorderQuestions reorders questions for a dish
// @Summary Reorder questions
// @Description Reorder questions for a specific dish
// @Tags questions
// @Accept json
// @Produce json
// @Param restaurantId path string true "Restaurant ID"
// @Param dishId path string true "Dish ID"
// @Param order body []string true "Question IDs in new order"
// @Success 200 {object} map[string]interface{} "Questions reordered successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 404 {object} map[string]interface{} "Dish not found"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /api/v1/restaurants/{restaurantId}/dishes/{dishId}/questions/reorder [post]
// @Security BearerAuth
func (h *QuestionHandler) ReorderQuestions(c echo.Context) error {
	accountID := middleware.GetResourceAccountID(c)

	dishID, err := uuid.Parse(c.Param("dishId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid dish ID")
	}

	var questionIDs []uuid.UUID
	if err := c.Bind(&questionIDs); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := h.questionService.ReorderQuestions(c.Request().Context(), accountID, dishID, questionIDs); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Questions reordered successfully",
	})
}

// GetDishesWithQuestions returns dish IDs that have questions for a restaurant
// @Summary Get dishes that have questions
// @Description Get list of dish IDs that have questions for a restaurant
// @Tags questions
// @Produce json
// @Param restaurantId path string true "Restaurant ID"
// @Success 200 {object} map[string]interface{} "Dishes with questions retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /api/v1/restaurants/{restaurantId}/questions/dishes-with-questions [get]
// @Security BearerAuth
func (h *QuestionHandler) GetDishesWithQuestions(c echo.Context) error {
	accountID := middleware.GetResourceAccountID(c)

	restaurantID, err := uuid.Parse(c.Param("restaurantId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid restaurant ID")
	}

	dishesWithQuestions, err := h.questionService.GetDishesWithQuestions(c.Request().Context(), accountID, restaurantID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    dishesWithQuestions,
	})
}