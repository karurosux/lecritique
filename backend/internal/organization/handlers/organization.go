package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"kyooar/internal/organization/models"
	"kyooar/internal/organization/services"
	"kyooar/internal/shared/errors"
	"kyooar/internal/shared/middleware"
	"kyooar/internal/shared/response"
	"kyooar/internal/shared/validator"
	"github.com/samber/do"
)

type OrganizationHandler struct {
	organizationService services.OrganizationService
	validator         *validator.Validator
}

func NewOrganizationHandler(i *do.Injector) (*OrganizationHandler, error) {
	return &OrganizationHandler{
		organizationService: do.MustInvoke[services.OrganizationService](i),
		validator:         validator.New(),
	}, nil
}

type CreateOrganizationRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Phone       string `json:"phone"`
	Email       string `json:"email" validate:"omitempty,email"`
	Website     string `json:"website"`
}

// Create godoc
// @Summary Create a new organization
// @Description Create a new organization for the authenticated account
// @Tags organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateOrganizationRequest true "Organization details"
// @Success 200 {object} response.Response{data=models.Organization}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/organizations [post]
func (h *OrganizationHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	accountID := middleware.GetResourceAccountID(c)

	var req CreateOrganizationRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	if err := h.validator.Validate(req); err != nil {
		return response.Error(c, errors.NewWithDetails("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, h.validator.FormatErrors(err)))
	}

	organization := &models.Organization{
		Name:        req.Name,
		Description: req.Description,
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

// GetAll godoc
// @Summary Get all organizations
// @Description Get all organizations for the authenticated account
// @Tags organizations
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]models.Organization}
// @Failure 401 {object} response.Response
// @Router /api/v1/organizations [get]
func (h *OrganizationHandler) GetAll(c echo.Context) error {
	ctx := c.Request().Context()
	accountID := middleware.GetResourceAccountID(c)

	organizations, err := h.organizationService.GetByAccountID(ctx, accountID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, organizations)
}

// GetByID gets a specific organization by ID
// @Summary Get organization by ID
// @Description Get a specific organization by its ID
// @Tags organizations
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Organization ID"
// @Success 200 {object} response.Response{data=models.Organization}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/organizations/{id} [get]
func (h *OrganizationHandler) GetByID(c echo.Context) error {
	ctx := c.Request().Context()
	accountID := middleware.GetResourceAccountID(c)
	
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

// Update updates a organization
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
func (h *OrganizationHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	accountID := middleware.GetResourceAccountID(c)
	
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

// Delete deletes a organization
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
func (h *OrganizationHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	accountID := middleware.GetResourceAccountID(c)
	
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
