package controller

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	feedbackinterface "kyooar/internal/feedback/interface"
	feedbackmodel "kyooar/internal/feedback/model"
	menuRepos "kyooar/internal/product/repositories"
	"kyooar/internal/shared/errors"
	"kyooar/internal/shared/logger"
	"kyooar/internal/shared/response"
	"kyooar/internal/shared/utils"
	"github.com/sirupsen/logrus"
)

type PublicController struct {
	feedbackService   feedbackinterface.FeedbackService
	productRepo          menuRepos.ProductRepository
	questionnaireRepo feedbackinterface.QuestionnaireRepository
	questionRepo      feedbackinterface.QuestionRepository
}

func NewPublicController(
	feedbackService feedbackinterface.FeedbackService,
	productRepo menuRepos.ProductRepository,
	questionnaireRepo feedbackinterface.QuestionnaireRepository,
	questionRepo feedbackinterface.QuestionRepository,
) *PublicController {
	return &PublicController{
		feedbackService:   feedbackService,
		productRepo:       productRepo,
		questionnaireRepo: questionnaireRepo,
		questionRepo:      questionRepo,
	}
}

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
func (h *PublicController) GetQuestionnaire(c echo.Context) error {
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

	return response.Success(c, map[string]interface{}{
		"organization_id": organizationID,
		"product_id":      productID,
		"message":         "Questionnaire endpoint - to be implemented",
	})
}

// @Summary Submit feedback
// @Description Submit customer feedback for a product
// @Tags public
// @Accept json
// @Produce json
// @Param feedback body feedbackmodel.Feedback true "Feedback data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/public/feedback [post]
func (h *PublicController) SubmitFeedback(c echo.Context) error {
	ctx := c.Request().Context()
	var feedback feedbackmodel.Feedback
	if err := c.Bind(&feedback); err != nil {
		return response.Error(c, errors.BadRequest("Invalid feedback data provided"))
	}

	deviceInfo := utils.ExtractDeviceInfo(c.Request())
	feedback.DeviceInfo = feedbackmodel.DeviceInfo{
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
func (h *PublicController) GetProductQuestions(c echo.Context) error {
	ctx := c.Request().Context()

	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		return response.Error(c, errors.BadRequest("Invalid product ID"))
	}

	product, err := h.productRepo.FindByID(ctx, productID)
	if err != nil {
		return response.Error(c, errors.NotFound("Product not found"))
	}

	questions, err := h.questionRepo.GetQuestionsByProductID(ctx, productID)
	if err != nil {
		logger.Error("Failed to get questions for product", err, logrus.Fields{
			"product_id": productID,
		})
		return response.Error(c, errors.Internal("Failed to get questions"))
	}

	return response.Success(c, map[string]interface{}{
		"product":   product,
		"questions": questions,
	})
}

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
func (h *PublicController) GetProductsWithQuestions(c echo.Context) error {
	ctx := c.Request().Context()

	organizationID, err := uuid.Parse(c.Param("organizationId"))
	if err != nil {
		return response.Error(c, errors.BadRequest("Invalid organization ID"))
	}

	productIDs, err := h.questionRepo.GetProductsWithQuestions(ctx, organizationID)
	if err != nil {
		logger.Error("Failed to get products with questions", err, logrus.Fields{
			"organization_id": organizationID,
		})
		return response.Error(c, errors.Internal("Failed to get products with questions"))
	}

	products := make([]interface{}, 0, len(productIDs))
	for _, productID := range productIDs {
		product, err := h.productRepo.FindByID(ctx, productID)
		if err != nil {
			logger.Warn("Failed to get product details", logrus.Fields{
				"product_id": productID,
				"error":      err,
			})
			continue
		}
		products = append(products, product)
	}

	return response.Success(c, products)
}