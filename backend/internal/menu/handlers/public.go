package handlers

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"lecritique/internal/shared/errors"
	"lecritique/internal/shared/response"
	"github.com/samber/do"
)

type MenuPublicHandler struct {
}

func NewMenuPublicHandler(i *do.Injector) (*MenuPublicHandler, error) {
	return &MenuPublicHandler{}, nil
}

// GetRestaurantMenu gets public restaurant menu
// @Summary Get restaurant menu
// @Description Get public menu for a restaurant
// @Tags public
// @Accept json
// @Produce json
// @Param id path string true "Restaurant ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/public/restaurant/{id}/menu [get]
func (h *MenuPublicHandler) GetRestaurantMenu(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.Error(c, errors.ErrInvalidUUID)
	}

	// Implementation would get restaurant menu
	// For now, return placeholder
	return response.Success(c, map[string]interface{}{
		"restaurant_id": id,
		"message":       "Menu endpoint - to be implemented",
	})
}