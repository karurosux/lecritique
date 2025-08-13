package handlers

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"kyooar/internal/shared/errors"
	"kyooar/internal/shared/response"
	"github.com/samber/do"
)

type ProductPublicHandler struct {
}

func NewProductPublicHandler(i *do.Injector) (*ProductPublicHandler, error) {
	return &ProductPublicHandler{}, nil
}

// @Summary Get organization products
// @Description Get public products for a organization
// @Tags public
// @Accept json
// @Produce json
// @Param id path string true "Organization ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/public/organization/{id}/products [get]
func (h *ProductPublicHandler) GetOrganizationProducts(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.Error(c, errors.ErrInvalidUUID)
	}

	return response.Success(c, map[string]interface{}{
		"organization_id": id,
		"message":       "Products endpoint - to be implemented",
	})
}
