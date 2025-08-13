package handlers

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"kyooar/internal/shared/errors"
	"kyooar/internal/shared/response"
	"github.com/samber/do"
)

type MenuPublicHandler struct {
}

func NewMenuPublicHandler(i *do.Injector) (*MenuPublicHandler, error) {
	return &MenuPublicHandler{}, nil
}

// @Summary Get organization menu
// @Description Get public menu for a organization
// @Tags public
// @Accept json
// @Produce json
// @Param id path string true "Organization ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/public/organization/{id}/menu [get]
func (h *MenuPublicHandler) GetOrganizationMenu(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.Error(c, errors.ErrInvalidUUID)
	}

	return response.Success(c, map[string]interface{}{
		"organization_id": id,
		"message":       "Menu endpoint - to be implemented",
	})
}
