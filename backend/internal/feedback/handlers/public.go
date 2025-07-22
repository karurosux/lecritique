package handlers

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	feedbackModels "kyooar/internal/feedback/models"
	feedbackRepos "kyooar/internal/feedback/repositories"
	feedbackServices "kyooar/internal/feedback/services"
	menuRepos "kyooar/internal/menu/repositories"
	"kyooar/internal/shared/errors"
	"kyooar/internal/shared/logger"
	"kyooar/internal/shared/response"
	"kyooar/internal/shared/utils"
	"github.com/samber/do"
	"github.com/sirupsen/logrus"
)

type FeedbackPublicHandler struct {
	feedbackService   feedbackServices.FeedbackService
	productRepo          menuRepos.ProductRepository
	questionnaireRepo feedbackRepos.QuestionnaireRepository
	questionRepo      feedbackRepos.QuestionRepository
}

func NewFeedbackPublicHandler(i *do.Injector) (*FeedbackPublicHandler, error) {
	return &FeedbackPublicHandler{
		feedbackService:   do.MustInvoke[feedbackServices.FeedbackService](i),
		productRepo:          do.MustInvoke[menuRepos.ProductRepository](i),
		questionnaireRepo: do.MustInvoke[feedbackRepos.QuestionnaireRepository](i),
		questionRepo:      do.MustInvoke[feedbackRepos.QuestionRepository](i),
	}, nil
}


// GetQuestionnaire gets questionnaire for a product
// @Summary Get questionnaire
// @Description Get questionnaire for a specific product
// @Tags public
// @Accept json
// @Produce json
// @Param organizationId path string true "Organization ID"
// @Param productId path string true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/public/questionnaire/{organizationId}/{productId} [get]
func (h *FeedbackPublicHandler) GetQuestionnaire(c echo.Context) error {
	organizationIDStr := c.Param("organizationId")
	productIDStr := c.Param("productId")

	organizationID, err := uuid.Parse(organizationIDStr)
	if err != nil {
		return response.Error(c, errors.BadRequest("Invalid organization ID format"))
	}

	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		return response.Error(c, errors.BadRequest("Invalid product ID format"))
	}

	// Implementation would get questionnaire
	// For now, return placeholder
	return response.Success(c, map[string]interface{}{
		"organization_id": organizationID,
		"product_id":       productID,
		"message":       "Questionnaire endpoint - to be implemented",
	})
}

// SubmitFeedback submits customer feedback
// @Summary Submit feedback
// @Description Submit customer feedback for a product
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

// GetProductQuestions gets questions for a specific product (public endpoint)
// @Summary Get questions for a product
// @Description Get all feedback questions for a specific product (public access for customer feedback)
// @Tags public
// @Produce json
// @Param organizationId path string true "Organization ID"
// @Param productId path string true "Product ID"
// @Success 200 {object} map[string]interface{} "Questions retrieved successfully"
// @Failure 400 {object} response.Response "Invalid request"
// @Failure 404 {object} response.Response "Product not found"
// @Failure 500 {object} response.Response "Server error"
// @Router /api/v1/public/organization/{organizationId}/products/{productId}/questions [get]
func (h *FeedbackPublicHandler) GetProductQuestions(c echo.Context) error {
	ctx := c.Request().Context()

	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		return response.Error(c, errors.BadRequest("Invalid product ID"))
	}

	// Verify product exists
	product, err := h.productRepo.FindByID(ctx, productID)
	if err != nil {
		return response.Error(c, errors.NotFound("Product not found"))
	}

	// Get questions for this product
	questions, err := h.questionRepo.GetQuestionsByProductID(ctx, productID)
	if err != nil {
		logger.Error("Failed to get questions for product", err, logrus.Fields{
			"product_id": productID,
		})
		return response.Error(c, errors.Internal("Failed to get questions"))
	}

	return response.Success(c, map[string]interface{}{
		"product":      product,
		"questions": questions,
	})
}

// GetProductsWithQuestions gets all products that have questions for a organization (public endpoint)
// @Summary Get products with questions
// @Description Get all products that have feedback questions for a organization (public access for QR code scans)
// @Tags public
// @Produce json
// @Param organizationId path string true "Organization ID"
// @Success 200 {object} map[string]interface{} "Products with questions retrieved successfully"
// @Failure 400 {object} response.Response "Invalid request"
// @Failure 404 {object} response.Response "Organization not found"
// @Failure 500 {object} response.Response "Server error"
// @Router /api/v1/public/organization/{organizationId}/questions/products-with-questions [get]
func (h *FeedbackPublicHandler) GetProductsWithQuestions(c echo.Context) error {
	ctx := c.Request().Context()

	organizationID, err := uuid.Parse(c.Param("organizationId"))
	if err != nil {
		return response.Error(c, errors.BadRequest("Invalid organization ID"))
	}

	// Get products with questions for this organization
	productIDs, err := h.questionRepo.GetProductsWithQuestions(ctx, organizationID)
	if err != nil {
		logger.Error("Failed to get products with questions", err, logrus.Fields{
			"organization_id": organizationID,
		})
		return response.Error(c, errors.Internal("Failed to get products with questions"))
	}

	// Get full product details for each product ID
	products := make([]interface{}, 0, len(productIDs))
	for _, productID := range productIDs {
		product, err := h.productRepo.FindByID(ctx, productID)
		if err != nil {
			logger.Warn("Failed to get product details", logrus.Fields{
				"product_id": productID,
				"error":   err,
			})
			continue
		}
		products = append(products, product)
	}

	return response.Success(c, products)
}
