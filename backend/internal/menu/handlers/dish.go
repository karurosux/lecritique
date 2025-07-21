package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"lecritique/internal/menu/models"
	"lecritique/internal/menu/services"
	"lecritique/internal/shared/errors"
	"lecritique/internal/shared/middleware"
	"lecritique/internal/shared/response"
	"lecritique/internal/shared/validator"
	"github.com/samber/do"
)

type DishHandler struct {
	dishService services.DishService
	validator   *validator.Validator
}

func NewDishHandler(i *do.Injector) (*DishHandler, error) {
	return &DishHandler{
		dishService: do.MustInvoke[services.DishService](i),
		validator:   validator.New(),
	}, nil
}

type CreateDishRequest struct {
	RestaurantID uuid.UUID `json:"restaurant_id" validate:"required"`
	Name         string    `json:"name" validate:"required"`
	Description  string    `json:"description"`
	Category     string    `json:"category"`
	Price        float64   `json:"price" validate:"min=0"`
	Currency     string    `json:"currency"`
}

// Create creates a new dish
// @Summary Create a new dish
// @Description Create a new dish for a restaurant
// @Tags dishes
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param restaurantId path string true "Restaurant ID"
// @Param dish body CreateDishRequest true "Dish information"
// @Success 200 {object} response.Response{data=models.Dish}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/restaurants/{restaurantId}/dishes [post]
func (h *DishHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	accountID := middleware.GetResourceAccountID(c)

	var req CreateDishRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, errors.BadRequest("Invalid dish data provided"))
	}

	if err := h.validator.Validate(req); err != nil {
		return response.Error(c, errors.NewWithDetails("VALIDATION_ERROR", "Please check the provided dish information", http.StatusBadRequest, h.validator.FormatErrors(err)))
	}

	dish := &models.Dish{
		RestaurantID: req.RestaurantID,
		Name:         req.Name,
		Description:  req.Description,
		Category:     req.Category,
		Price:        req.Price,
		Currency:     req.Currency,
		IsAvailable:  true,
		IsActive:     true,
	}

	if dish.Currency == "" {
		dish.Currency = "USD"
	}

	if err := h.dishService.Create(ctx, accountID, dish); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, dish)
}

// GetByRestaurant gets all dishes for a specific restaurant
// @Summary Get dishes by restaurant
// @Description Get all dishes for a specific restaurant
// @Tags dishes
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param restaurantId path string true "Restaurant ID"
// @Success 200 {object} response.Response{data=[]models.Dish}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/restaurants/{restaurantId}/dishes [get]
func (h *DishHandler) GetByRestaurant(c echo.Context) error {
	ctx := c.Request().Context()
	accountID := middleware.GetResourceAccountID(c)

	restaurantID, err := uuid.Parse(c.Param("restaurantId"))
	if err != nil {
		return response.Error(c, errors.ErrInvalidUUID)
	}

	dishes, err := h.dishService.GetByRestaurantID(ctx, accountID, restaurantID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, dishes)
}

// GetByID gets a specific dish by ID
// @Summary Get dish by ID
// @Description Get a specific dish by its ID
// @Tags dishes
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Dish ID"
// @Success 200 {object} response.Response{data=models.Dish}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/dishes/{id} [get]
func (h *DishHandler) GetByID(c echo.Context) error {
	ctx := c.Request().Context()
	accountID := middleware.GetResourceAccountID(c)

	dishID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, errors.ErrInvalidUUID)
	}

	dish, err := h.dishService.GetByID(ctx, accountID, dishID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, dish)
}

// Update updates a dish
// @Summary Update a dish
// @Description Update a dish's information
// @Tags dishes
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Dish ID"
// @Param updates body map[string]interface{} true "Fields to update"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/dishes/{id} [put]
func (h *DishHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	accountID := middleware.GetResourceAccountID(c)

	dishID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, errors.ErrInvalidUUID)
	}

	var updates map[string]interface{}
	if err := c.Bind(&updates); err != nil {
		return response.Error(c, errors.BadRequest("Invalid update data provided"))
	}

	if err := h.dishService.Update(ctx, accountID, dishID, updates); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]string{
		"message": "Dish updated successfully",
	})
}

// Delete deletes a dish
// @Summary Delete a dish
// @Description Delete a dish from the system
// @Tags dishes
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Dish ID"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/dishes/{id} [delete]
func (h *DishHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	accountID := middleware.GetResourceAccountID(c)

	dishID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, errors.ErrInvalidUUID)
	}

	if err := h.dishService.Delete(ctx, accountID, dishID); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]string{
		"message": "Dish deleted successfully",
	})
}
