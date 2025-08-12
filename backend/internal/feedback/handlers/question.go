package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"kyooar/internal/feedback/models"
	"kyooar/internal/feedback/services"
	"kyooar/internal/shared/middleware"
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

// @Summary Add a question to a product
// @Description Add a new feedback question to a specific product
// @Tags questions
// @Accept json
// @Produce json
// @Param organizationId path string true "Organization ID"
// @Param productId path string true "Product ID"
// @Param question body models.CreateQuestionRequest true "Question data"
// @Success 201 {object} map[string]interface{} "Question created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 404 {object} map[string]interface{} "Product not found"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /api/v1/organizations/{organizationId}/products/{productId}/questions [post]
// @Security BearerAuth
func (h *QuestionHandler) CreateQuestion(c echo.Context) error {
	accountID := middleware.GetResourceAccountID(c)

	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid product ID")
	}

	var request models.CreateQuestionRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	question, err := h.questionService.CreateQuestion(c.Request().Context(), accountID, productID, &request)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Question created successfully",
		"data":    question,
	})
}

// @Summary Get questions for a product
// @Description Get all feedback questions for a specific product
// @Tags questions
// @Produce json
// @Param organizationId path string true "Organization ID"
// @Param productId path string true "Product ID"
// @Success 200 {object} map[string]interface{} "Questions retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 404 {object} map[string]interface{} "Product not found"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /api/v1/organizations/{organizationId}/products/{productId}/questions [get]
// @Security BearerAuth
func (h *QuestionHandler) GetQuestionsByProduct(c echo.Context) error {
	accountID := middleware.GetResourceAccountID(c)

	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid product ID")
	}

	questions, err := h.questionService.GetQuestionsByProduct(c.Request().Context(), accountID, productID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    questions,
	})
}

// @Summary Get questions for multiple products (optimized payload)
// @Description Get essential question fields for multiple products in a single request - returns only ID, ProductID, Text, and Type
// @Tags questions
// @Accept json
// @Produce json
// @Param organizationId path string true "Organization ID"
// @Param request body models.BatchQuestionsRequest true "Product IDs"
// @Success 200 {object} map[string]interface{} "Questions retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 404 {object} map[string]interface{} "Organization not found"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /api/v1/organizations/{organizationId}/questions/batch [post]
// @Security BearerAuth
func (h *QuestionHandler) GetQuestionsByProducts(c echo.Context) error {
	accountID := middleware.GetResourceAccountID(c)

	organizationID, err := uuid.Parse(c.Param("organizationId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization ID")
	}

	var request models.BatchQuestionsRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	questions, err := h.questionService.GetQuestionsByProductsOptimized(c.Request().Context(), accountID, organizationID, request.ProductIDs)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    questions,
	})
}

// @Summary Get a specific question
// @Description Get details of a specific question
// @Tags questions
// @Produce json
// @Param organizationId path string true "Organization ID"
// @Param productId path string true "Product ID"
// @Param questionId path string true "Question ID"
// @Success 200 {object} map[string]interface{} "Question retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 404 {object} map[string]interface{} "Question not found"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /api/v1/organizations/{organizationId}/products/{productId}/questions/{questionId} [get]
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

// @Summary Update a question
// @Description Update an existing question for a product
// @Tags questions
// @Accept json
// @Produce json
// @Param organizationId path string true "Organization ID"
// @Param productId path string true "Product ID"
// @Param questionId path string true "Question ID"
// @Param question body models.UpdateQuestionRequest true "Question data"
// @Success 200 {object} map[string]interface{} "Question updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 404 {object} map[string]interface{} "Question not found"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /api/v1/organizations/{organizationId}/products/{productId}/questions/{questionId} [put]
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

// @Summary Delete a question
// @Description Delete a feedback question from a product
// @Tags questions
// @Produce json
// @Param organizationId path string true "Organization ID"
// @Param productId path string true "Product ID"
// @Param questionId path string true "Question ID"
// @Success 200 {object} map[string]interface{} "Question deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 404 {object} map[string]interface{} "Question not found"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /api/v1/organizations/{organizationId}/products/{productId}/questions/{questionId} [delete]
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

// @Summary Reorder questions
// @Description Reorder questions for a specific product
// @Tags questions
// @Accept json
// @Produce json
// @Param organizationId path string true "Organization ID"
// @Param productId path string true "Product ID"
// @Param order body []string true "Question IDs in new order"
// @Success 200 {object} map[string]interface{} "Questions reordered successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 404 {object} map[string]interface{} "Product not found"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /api/v1/organizations/{organizationId}/products/{productId}/questions/reorder [post]
// @Security BearerAuth
func (h *QuestionHandler) ReorderQuestions(c echo.Context) error {
	accountID := middleware.GetResourceAccountID(c)

	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid product ID")
	}

	var questionIDs []uuid.UUID
	if err := c.Bind(&questionIDs); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := h.questionService.ReorderQuestions(c.Request().Context(), accountID, productID, questionIDs); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Questions reordered successfully",
	})
}

// @Summary Get products that have questions
// @Description Get list of product IDs that have questions for a organization
// @Tags questions
// @Produce json
// @Param organizationId path string true "Organization ID"
// @Success 200 {object} map[string]interface{} "Products with questions retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /api/v1/organizations/{organizationId}/questions/products-with-questions [get]
// @Security BearerAuth
func (h *QuestionHandler) GetProductsWithQuestions(c echo.Context) error {
	accountID := middleware.GetResourceAccountID(c)

	organizationID, err := uuid.Parse(c.Param("organizationId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization ID")
	}

	productsWithQuestions, err := h.questionService.GetProductsWithQuestions(c.Request().Context(), accountID, organizationID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    productsWithQuestions,
	})
}
