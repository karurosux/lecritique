package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"kyooar/internal/feedback/models"
	"kyooar/internal/feedback/services"
	menuServices "kyooar/internal/menu/services"
	"kyooar/internal/shared/middleware"
	"github.com/samber/do"
)

type QuestionnaireHandler struct {
	questionnaireService *services.QuestionnaireService
	productService         menuServices.ProductService
}

func NewQuestionnaireHandler(i *do.Injector) (*QuestionnaireHandler, error) {
	return &QuestionnaireHandler{
		questionnaireService: do.MustInvoke[*services.QuestionnaireService](i),
		productService:         do.MustInvoke[menuServices.ProductService](i),
	}, nil
}

// CreateQuestionnaire creates a new questionnaire
// @Summary Create questionnaire
// @Description Create a new questionnaire for a organization
// @Tags questionnaires
// @Accept json
// @Produce json
// @Param organizationId path string true "Organization ID"
// @Param questionnaire body models.CreateQuestionnaireRequest true "Questionnaire data"
// @Success 201 {object} response.Response{data=models.Questionnaire}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/organizations/{organizationId}/questionnaires [post]
// @Security Bearer
func (h *QuestionnaireHandler) CreateQuestionnaire(c echo.Context) error {
	accountID := middleware.GetResourceAccountID(c)
	organizationID, err := uuid.Parse(c.Param("organizationId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization ID")
	}

	var input models.Questionnaire
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	questionnaire, err := h.questionnaireService.Create(c.Request().Context(), accountID, organizationID, &input)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create questionnaire")
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Questionnaire created successfully",
		"data":    questionnaire,
	})
}

// GetQuestionnaire retrieves a questionnaire by ID
// @Summary Get questionnaire
// @Description Get a specific questionnaire by ID
// @Tags questionnaires
// @Accept json
// @Produce json
// @Param organizationId path string true "Organization ID"
// @Param id path string true "Questionnaire ID"
// @Success 200 {object} response.Response{data=models.Questionnaire}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/organizations/{organizationId}/questionnaires/{id} [get]
// @Security Bearer
func (h *QuestionnaireHandler) GetQuestionnaire(c echo.Context) error {
	accountID := middleware.GetResourceAccountID(c)
	questionnaireID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid questionnaire ID")
	}

	questionnaire, err := h.questionnaireService.GetByID(c.Request().Context(), accountID, questionnaireID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Questionnaire not found")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Questionnaire retrieved successfully",
		"data":    questionnaire,
	})
}

// ListQuestionnaires lists all questionnaires for a organization
// @Summary List questionnaires
// @Description Get all questionnaires for a organization
// @Tags questionnaires
// @Accept json
// @Produce json
// @Param organizationId path string true "Organization ID"
// @Success 200 {object} response.Response{data=[]models.Questionnaire}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/organizations/{organizationId}/questionnaires [get]
// @Security Bearer
func (h *QuestionnaireHandler) ListQuestionnaires(c echo.Context) error {
	accountID := middleware.GetResourceAccountID(c)
	organizationID, err := uuid.Parse(c.Param("organizationId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization ID")
	}

	questionnaires, err := h.questionnaireService.ListByOrganization(c.Request().Context(), accountID, organizationID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list questionnaires")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Questionnaires retrieved successfully",
		"data":    questionnaires,
	})
}

// UpdateQuestionnaire updates a questionnaire
// @Summary Update questionnaire
// @Description Update an existing questionnaire
// @Tags questionnaires
// @Accept json
// @Produce json
// @Param organizationId path string true "Organization ID"
// @Param id path string true "Questionnaire ID"
// @Param questionnaire body models.Questionnaire true "Questionnaire data"
// @Success 200 {object} response.Response{data=models.Questionnaire}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/organizations/{organizationId}/questionnaires/{id} [put]
// @Security Bearer
func (h *QuestionnaireHandler) UpdateQuestionnaire(c echo.Context) error {
	accountID := middleware.GetResourceAccountID(c)
	questionnaireID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid questionnaire ID")
	}

	var input models.Questionnaire
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	questionnaire, err := h.questionnaireService.Update(c.Request().Context(), accountID, questionnaireID, &input)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update questionnaire")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Questionnaire updated successfully",
		"data":    questionnaire,
	})
}

// DeleteQuestionnaire deletes a questionnaire
// @Summary Delete questionnaire
// @Description Delete a questionnaire
// @Tags questionnaires
// @Accept json
// @Produce json
// @Param organizationId path string true "Organization ID"
// @Param id path string true "Questionnaire ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/organizations/{organizationId}/questionnaires/{id} [delete]
// @Security Bearer
func (h *QuestionnaireHandler) DeleteQuestionnaire(c echo.Context) error {
	accountID := middleware.GetResourceAccountID(c)
	questionnaireID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid questionnaire ID")
	}

	if err := h.questionnaireService.Delete(c.Request().Context(), accountID, questionnaireID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete questionnaire")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Questionnaire deleted successfully",
	})
}

// AddQuestion adds a question to a questionnaire
// @Summary Add a question to a questionnaire
// @Description Add a new question to an existing questionnaire
// @Tags questionnaires
// @Accept json
// @Produce json
// @Param id path string true "Questionnaire ID"
// @Param question body models.Question true "Question data"
// @Success 201 {object} map[string]interface{} "Question added successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /organizations/{organizationId}/questionnaires/{id}/questions [post]
// @Security BearerAuth
func (h *QuestionnaireHandler) AddQuestion(c echo.Context) error {
	accountID := middleware.GetResourceAccountID(c)
	questionnaireID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid questionnaire ID")
	}

	var input models.Question
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	question, err := h.questionnaireService.AddQuestion(c.Request().Context(), accountID, questionnaireID, &input)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to add question")
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Question added successfully",
		"data":    question,
	})
}

// UpdateQuestion updates a question
// @Summary Update a question
// @Description Update an existing question in a questionnaire
// @Tags questionnaires
// @Accept json
// @Produce json
// @Param id path string true "Questionnaire ID"
// @Param questionId path string true "Question ID"
// @Param question body models.Question true "Question data"
// @Success 200 {object} map[string]interface{} "Question updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /organizations/{organizationId}/questionnaires/{id}/questions/{questionId} [put]
// @Security BearerAuth
func (h *QuestionnaireHandler) UpdateQuestion(c echo.Context) error {
	accountID := middleware.GetResourceAccountID(c)
	questionID, err := uuid.Parse(c.Param("questionId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid question ID")
	}

	var input models.Question
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	question, err := h.questionnaireService.UpdateQuestion(c.Request().Context(), accountID, questionID, &input)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update question")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Question updated successfully",
		"data":    question,
	})
}

// DeleteQuestion deletes a question
// @Summary Delete a question
// @Description Delete a question from a questionnaire
// @Tags questionnaires
// @Accept json
// @Produce json
// @Param id path string true "Questionnaire ID"
// @Param questionId path string true "Question ID"
// @Success 200 {object} map[string]interface{} "Question deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /organizations/{organizationId}/questionnaires/{id}/questions/{questionId} [delete]
// @Security BearerAuth
func (h *QuestionnaireHandler) DeleteQuestion(c echo.Context) error {
	accountID := middleware.GetResourceAccountID(c)
	questionID, err := uuid.Parse(c.Param("questionId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid question ID")
	}

	if err := h.questionnaireService.DeleteQuestion(c.Request().Context(), accountID, questionID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete question")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Question deleted successfully",
	})
}

// ReorderQuestions reorders questions in a questionnaire
// @Summary Reorder questions
// @Description Reorder questions in a questionnaire
// @Tags questionnaires
// @Accept json
// @Produce json
// @Param id path string true "Questionnaire ID"
// @Param order body []uuid.UUID true "Question IDs in new order"
// @Success 200 {object} map[string]interface{} "Questions reordered successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /organizations/{organizationId}/questionnaires/{id}/reorder [post]
// @Security BearerAuth
func (h *QuestionnaireHandler) ReorderQuestions(c echo.Context) error {
	accountID := middleware.GetResourceAccountID(c)
	questionnaireID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid questionnaire ID")
	}

	var input struct {
		QuestionIDs []uuid.UUID `json:"question_ids"`
	}
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := h.questionnaireService.ReorderQuestions(c.Request().Context(), accountID, questionnaireID, input.QuestionIDs); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to reorder questions")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Questions reordered successfully",
	})
}

// GenerateQuestions generates AI-powered questions for a product
// @Summary Generate AI questions
// @Description Generate AI-powered questions for a specific product
// @Tags questionnaires,ai
// @Accept json
// @Produce json
// @Param productId path string true "Product ID"
// @Success 200 {object} response.Response{data=[]models.GeneratedQuestion}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/ai/generate-questions/{productId} [post]
// @Security Bearer
func (h *QuestionnaireHandler) GenerateQuestions(c echo.Context) error {
	accountID := middleware.GetResourceAccountID(c)
	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid product ID")
	}

	// Get the product details
	product, err := h.productService.GetByID(c.Request().Context(), accountID, productID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	// Generate questions using AI
	questions, err := h.questionnaireService.GenerateQuestionsForProduct(c.Request().Context(), accountID, product)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate questions")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Questions generated successfully",
		"data":    questions,
	})
}

// GenerateAndSaveQuestionnaire generates AI questions and creates a complete questionnaire for a product
// @Summary Generate and save AI questionnaire
// @Description Generate AI questions and create a complete questionnaire for a product
// @Tags questionnaires,ai
// @Accept json
// @Produce json
// @Param productId path string true "Product ID"
// @Param questionnaire body models.GenerateQuestionnaireRequest true "Questionnaire generation data"
// @Success 201 {object} response.Response{data=models.Questionnaire}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/ai/generate-questionnaire/{productId} [post]
// @Security Bearer
func (h *QuestionnaireHandler) GenerateAndSaveQuestionnaire(c echo.Context) error {
	accountID := middleware.GetResourceAccountID(c)
	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid product ID")
	}

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
		IsDefault   bool   `json:"is_default,omitempty"`
	}
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Get the product details
	product, err := h.productService.GetByID(c.Request().Context(), accountID, productID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	// Generate and save questionnaire with AI questions
	questionnaire, err := h.questionnaireService.GenerateAndSaveQuestionnaireForProduct(c.Request().Context(), accountID, product, input.Name, input.Description, input.IsDefault)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate and save questionnaire")
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Questionnaire generated and saved successfully",
		"data":    questionnaire,
	})
}
