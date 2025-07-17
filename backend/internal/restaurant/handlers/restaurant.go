package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lecritique/api/internal/restaurant/models"
	"github.com/lecritique/api/internal/restaurant/services"
	"github.com/lecritique/api/internal/shared/errors"
	"github.com/lecritique/api/internal/shared/middleware"
	"github.com/lecritique/api/internal/shared/response"
	"github.com/lecritique/api/internal/shared/validator"
)

type RestaurantHandler struct {
	restaurantService services.RestaurantService
	validator         *validator.Validator
}

func NewRestaurantHandler(restaurantService services.RestaurantService) *RestaurantHandler {
	return &RestaurantHandler{
		restaurantService: restaurantService,
		validator:         validator.New(),
	}
}

type CreateRestaurantRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Phone       string `json:"phone"`
	Email       string `json:"email" validate:"omitempty,email"`
	Website     string `json:"website"`
}

// Create godoc
// @Summary Create a new restaurant
// @Description Create a new restaurant for the authenticated account
// @Tags restaurants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateRestaurantRequest true "Restaurant details"
// @Success 200 {object} response.Response{data=models.Restaurant}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/restaurants [post]
func (h *RestaurantHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	accountID := middleware.GetResourceAccountID(c)

	var req CreateRestaurantRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	if err := h.validator.Validate(req); err != nil {
		return response.Error(c, errors.NewWithDetails("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, h.validator.FormatErrors(err)))
	}

	restaurant := &models.Restaurant{
		Name:        req.Name,
		Description: req.Description,
		Phone:       req.Phone,
		Email:       req.Email,
		Website:     req.Website,
		IsActive:    true,
	}

	if err := h.restaurantService.Create(ctx, accountID, restaurant); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, restaurant)
}

// GetAll godoc
// @Summary Get all restaurants
// @Description Get all restaurants for the authenticated account
// @Tags restaurants
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]models.Restaurant}
// @Failure 401 {object} response.Response
// @Router /api/v1/restaurants [get]
func (h *RestaurantHandler) GetAll(c echo.Context) error {
	ctx := c.Request().Context()
	accountID := middleware.GetResourceAccountID(c)

	restaurants, err := h.restaurantService.GetByAccountID(ctx, accountID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, restaurants)
}

// GetByID gets a specific restaurant by ID
// @Summary Get restaurant by ID
// @Description Get a specific restaurant by its ID
// @Tags restaurants
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Restaurant ID"
// @Success 200 {object} response.Response{data=models.Restaurant}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/restaurants/{id} [get]
func (h *RestaurantHandler) GetByID(c echo.Context) error {
	ctx := c.Request().Context()
	accountID := middleware.GetResourceAccountID(c)
	
	restaurantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	restaurant, err := h.restaurantService.GetByID(ctx, accountID, restaurantID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, restaurant)
}

// Update updates a restaurant
// @Summary Update restaurant
// @Description Update a restaurant's information
// @Tags restaurants
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Restaurant ID"
// @Param updates body map[string]interface{} true "Fields to update"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/restaurants/{id} [put]
func (h *RestaurantHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	accountID := middleware.GetResourceAccountID(c)
	
	restaurantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	var updates map[string]interface{}
	if err := c.Bind(&updates); err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	if err := h.restaurantService.Update(ctx, accountID, restaurantID, updates); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]string{
		"message": "Restaurant updated successfully",
	})
}

// Delete deletes a restaurant
// @Summary Delete restaurant
// @Description Delete a restaurant from the system
// @Tags restaurants
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Restaurant ID"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/restaurants/{id} [delete]
func (h *RestaurantHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	accountID := middleware.GetResourceAccountID(c)
	
	restaurantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	if err := h.restaurantService.Delete(ctx, accountID, restaurantID); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]string{
		"message": "Restaurant deleted successfully",
	})
}
