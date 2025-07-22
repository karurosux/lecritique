package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"kyooar/internal/menu/models"
	"kyooar/internal/menu/services"
	"kyooar/internal/shared/errors"
	"kyooar/internal/shared/middleware"
	"kyooar/internal/shared/response"
	"kyooar/internal/shared/validator"
	"github.com/samber/do"
)

type ProductHandler struct {
	productService services.ProductService
	validator   *validator.Validator
}

func NewProductHandler(i *do.Injector) (*ProductHandler, error) {
	return &ProductHandler{
		productService: do.MustInvoke[services.ProductService](i),
		validator:   validator.New(),
	}, nil
}

type CreateProductRequest struct {
	OrganizationID uuid.UUID `json:"organization_id" validate:"required"`
	Name         string    `json:"name" validate:"required"`
	Description  string    `json:"description"`
	Category     string    `json:"category"`
	Price        float64   `json:"price" validate:"min=0"`
	Currency     string    `json:"currency"`
}

// Create creates a new product
// @Summary Create a new product
// @Description Create a new product for a organization
// @Tags products
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param organizationId path string true "Organization ID"
// @Param product body CreateProductRequest true "Product information"
// @Success 200 {object} response.Response{data=models.Product}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/organizations/{organizationId}/products [post]
func (h *ProductHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	accountID := middleware.GetResourceAccountID(c)

	var req CreateProductRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, errors.BadRequest("Invalid product data provided"))
	}

	if err := h.validator.Validate(req); err != nil {
		return response.Error(c, errors.NewWithDetails("VALIDATION_ERROR", "Please check the provided product information", http.StatusBadRequest, h.validator.FormatErrors(err)))
	}

	product := &models.Product{
		OrganizationID: req.OrganizationID,
		Name:         req.Name,
		Description:  req.Description,
		Category:     req.Category,
		Price:        req.Price,
		Currency:     req.Currency,
		IsAvailable:  true,
		IsActive:     true,
	}

	if product.Currency == "" {
		product.Currency = "USD"
	}

	if err := h.productService.Create(ctx, accountID, product); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, product)
}

// GetByOrganization gets all products for a specific organization
// @Summary Get products by organization
// @Description Get all products for a specific organization
// @Tags products
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param organizationId path string true "Organization ID"
// @Success 200 {object} response.Response{data=[]models.Product}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/organizations/{organizationId}/products [get]
func (h *ProductHandler) GetByOrganization(c echo.Context) error {
	ctx := c.Request().Context()
	accountID := middleware.GetResourceAccountID(c)

	organizationID, err := uuid.Parse(c.Param("organizationId"))
	if err != nil {
		return response.Error(c, errors.ErrInvalidUUID)
	}

	products, err := h.productService.GetByOrganizationID(ctx, accountID, organizationID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, products)
}

// GetByID gets a specific product by ID
// @Summary Get product by ID
// @Description Get a specific product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Product ID"
// @Success 200 {object} response.Response{data=models.Product}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/products/{id} [get]
func (h *ProductHandler) GetByID(c echo.Context) error {
	ctx := c.Request().Context()
	accountID := middleware.GetResourceAccountID(c)

	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, errors.ErrInvalidUUID)
	}

	product, err := h.productService.GetByID(ctx, accountID, productID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, product)
}

// Update updates a product
// @Summary Update a product
// @Description Update a product's information
// @Tags products
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Product ID"
// @Param updates body map[string]interface{} true "Fields to update"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/products/{id} [put]
func (h *ProductHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	accountID := middleware.GetResourceAccountID(c)

	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, errors.ErrInvalidUUID)
	}

	var updates map[string]interface{}
	if err := c.Bind(&updates); err != nil {
		return response.Error(c, errors.BadRequest("Invalid update data provided"))
	}

	if err := h.productService.Update(ctx, accountID, productID, updates); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]string{
		"message": "Product updated successfully",
	})
}

// Delete deletes a product
// @Summary Delete a product
// @Description Delete a product from the system
// @Tags products
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Product ID"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/products/{id} [delete]
func (h *ProductHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	accountID := middleware.GetResourceAccountID(c)

	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, errors.ErrInvalidUUID)
	}

	if err := h.productService.Delete(ctx, accountID, productID); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]string{
		"message": "Product deleted successfully",
	})
}
