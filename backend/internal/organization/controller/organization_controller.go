package organizationcontroller

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	feedbackcontroller "kyooar/internal/feedback/controller"
	menuHandlers "kyooar/internal/product/handlers"
	organizationinterface "kyooar/internal/organization/interface"
	organizationmodel "kyooar/internal/organization/model"
	qrcodecontroller "kyooar/internal/qrcode/controller"
	"kyooar/internal/shared/errors"
	sharedMiddleware "kyooar/internal/shared/middleware"
	"kyooar/internal/shared/response"
	"kyooar/internal/shared/validator"
	subscriptionMiddleware "kyooar/internal/subscription/middleware"
	subscriptionmodel "kyooar/internal/subscription/model"
)

type OrganizationController struct {
	organizationService     organizationinterface.OrganizationService
	productHandler         *menuHandlers.ProductHandler
	qrCodeHandler          *qrcodecontroller.QRCodeController
	feedbackController     *feedbackcontroller.FeedbackController
	questionnaireController *feedbackcontroller.QuestionnaireController
	questionController     *feedbackcontroller.QuestionController
	validator              *validator.Validator
}

func NewOrganizationController(
	organizationService organizationinterface.OrganizationService,
	productHandler *menuHandlers.ProductHandler,
	qrCodeHandler *qrcodecontroller.QRCodeController,
	feedbackController *feedbackcontroller.FeedbackController,
	questionnaireController *feedbackcontroller.QuestionnaireController,
	questionController *feedbackcontroller.QuestionController,
) *OrganizationController {
	return &OrganizationController{
		organizationService:     organizationService,
		productHandler:         productHandler,
		qrCodeHandler:          qrCodeHandler,
		feedbackController:     feedbackController,
		questionnaireController: questionnaireController,
		questionController:     questionController,
		validator:              validator.New(),
	}
}

func (c *OrganizationController) RegisterRoutes(v1 *echo.Group, middlewareProvider *sharedMiddleware.MiddlewareProvider, subscriptionMW *subscriptionMiddleware.SubscriptionMiddleware) {
	organizations := v1.Group("/organizations")
	organizations.Use(middlewareProvider.AuthMiddleware())
	organizations.Use(middlewareProvider.TeamAwareMiddleware())
	
	// Organization CRUD routes
	organizations.POST("", c.Create,
		subscriptionMW.CheckResourceLimit(subscriptionmodel.ResourceTypeOrganization),
		subscriptionMW.TrackUsageAfterSuccess(),
	)
	organizations.GET("", c.GetAll)
	organizations.GET("/:id", c.GetByID)
	organizations.PUT("/:id", c.Update)
	organizations.DELETE("/:id", c.Delete)
	
	// Organization-scoped product routes
	organizations.GET("/:organizationId/products", c.productHandler.GetByOrganization)
	organizations.POST("/:organizationId/products", c.productHandler.Create)
	
	// Organization-scoped QR code routes
	organizations.POST("/:organizationId/qr-codes", c.qrCodeHandler.Generate)
	organizations.GET("/:organizationId/qr-codes", c.qrCodeHandler.GetByOrganization)
	
	// Organization-scoped feedback routes
	organizations.GET("/:organizationId/feedback", c.feedbackController.GetByOrganization)
	organizations.GET("/:organizationId/analytics", c.feedbackController.GetStats)
	
	// Organization-scoped questionnaire routes
	organizations.POST("/:organizationId/questionnaires", c.questionnaireController.CreateQuestionnaire)
	organizations.GET("/:organizationId/questionnaires", c.questionnaireController.ListQuestionnaires)
	organizations.GET("/:organizationId/questionnaires/:id", c.questionnaireController.GetQuestionnaire)
	organizations.PUT("/:organizationId/questionnaires/:id", c.questionnaireController.UpdateQuestionnaire)
	organizations.DELETE("/:organizationId/questionnaires/:id", c.questionnaireController.DeleteQuestionnaire)
	organizations.POST("/:organizationId/questionnaires/:id/questions", c.questionnaireController.AddQuestion)
	organizations.PUT("/:organizationId/questionnaires/:id/questions/:questionId", c.questionnaireController.UpdateQuestion)
	organizations.DELETE("/:organizationId/questionnaires/:id/questions/:questionId", c.questionnaireController.DeleteQuestion)
	organizations.POST("/:organizationId/questionnaires/:id/reorder", c.questionnaireController.ReorderQuestions)
	
	// Organization-scoped question routes
	organizations.POST("/:organizationId/products/:productId/questions", c.questionController.CreateQuestion)
	organizations.GET("/:organizationId/products/:productId/questions", c.questionController.GetQuestionsByProduct)
	organizations.GET("/:organizationId/questions/batch", c.questionController.GetQuestionsByProducts)
	organizations.GET("/:organizationId/products/:productId/questions/:questionId", c.questionController.GetQuestion)
	organizations.PUT("/:organizationId/products/:productId/questions/:questionId", c.questionController.UpdateQuestion)
	organizations.DELETE("/:organizationId/products/:productId/questions/:questionId", c.questionController.DeleteQuestion)
	organizations.POST("/:organizationId/products/:productId/questions/reorder", c.questionController.ReorderQuestions)
	organizations.GET("/:organizationId/questions/products-with-questions", c.questionController.GetProductsWithQuestions)
	
	// AI-powered routes
	organizations.POST("/:organizationId/products/:productId/ai/generate-questions", c.questionnaireController.GenerateQuestions)
	organizations.POST("/:organizationId/products/:productId/ai/generate-questionnaire", c.questionnaireController.GenerateAndSaveQuestionnaire)
}

// @Summary Create a new organization
// @Description Create a new organization for the authenticated account
// @Tags organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body organizationmodel.CreateOrganizationRequest true "Organization details"
// @Success 200 {object} response.Response{data=organizationmodel.Organization}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/organizations [post]
func (h *OrganizationController) Create(c echo.Context) error {
	ctx := c.Request().Context()
	accountID := sharedMiddleware.GetResourceAccountID(c)

	var req organizationmodel.CreateOrganizationRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	if err := h.validator.Validate(req); err != nil {
		return response.Error(c, errors.NewWithDetails("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, h.validator.FormatErrors(err)))
	}

	organization := &organizationmodel.Organization{
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address,
		Phone:       req.Phone,
		Email:       req.Email,
		Website:     req.Website,
		IsActive:    true,
	}

	if err := h.organizationService.Create(ctx, accountID, organization); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, organization)
}

// @Summary Get all organizations
// @Description Get all organizations for the authenticated account
// @Tags organizations
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]organizationmodel.Organization}
// @Failure 401 {object} response.Response
// @Router /api/v1/organizations [get]
func (h *OrganizationController) GetAll(c echo.Context) error {
	ctx := c.Request().Context()
	accountID := sharedMiddleware.GetResourceAccountID(c)

	organizations, err := h.organizationService.GetByAccountID(ctx, accountID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, organizations)
}

// @Summary Get organization by ID
// @Description Get a specific organization by its ID
// @Tags organizations
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Organization ID"
// @Success 200 {object} response.Response{data=organizationmodel.Organization}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/organizations/{id} [get]
func (h *OrganizationController) GetByID(c echo.Context) error {
	ctx := c.Request().Context()
	accountID := sharedMiddleware.GetResourceAccountID(c)
	
	organizationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	organization, err := h.organizationService.GetByID(ctx, accountID, organizationID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, organization)
}

// @Summary Update organization
// @Description Update a organization's information
// @Tags organizations
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Organization ID"
// @Param updates body map[string]interface{} true "Fields to update"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/organizations/{id} [put]
func (h *OrganizationController) Update(c echo.Context) error {
	ctx := c.Request().Context()
	accountID := sharedMiddleware.GetResourceAccountID(c)
	
	organizationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	var updates map[string]interface{}
	if err := c.Bind(&updates); err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	if err := h.organizationService.Update(ctx, accountID, organizationID, updates); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]string{
		"message": "Organization updated successfully",
	})
}

// @Summary Delete organization
// @Description Delete a organization from the system
// @Tags organizations
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Organization ID"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/organizations/{id} [delete]
func (h *OrganizationController) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	accountID := sharedMiddleware.GetResourceAccountID(c)
	
	organizationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	if err := h.organizationService.Delete(ctx, accountID, organizationID); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]string{
		"message": "Organization deleted successfully",
	})
}
