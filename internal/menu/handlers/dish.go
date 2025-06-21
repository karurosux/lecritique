package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lecritique/api/internal/menu/models"
	"github.com/lecritique/api/internal/menu/services"
	"github.com/lecritique/api/internal/shared/errors"
	"github.com/lecritique/api/internal/shared/response"
	"github.com/lecritique/api/internal/shared/validator"
)

type DishHandler struct {
	dishService services.DishService
	validator   *validator.Validator
}

func NewDishHandler(dishService services.DishService) *DishHandler {
	return &DishHandler{
		dishService: dishService,
		validator:   validator.New(),
	}
}

type CreateDishRequest struct {
	RestaurantID uuid.UUID `json:"restaurant_id" validate:"required"`
	Name         string    `json:"name" validate:"required"`
	Description  string    `json:"description"`
	Category     string    `json:"category"`
	Price        float64   `json:"price" validate:"min=0"`
	Currency     string    `json:"currency"`
}

func (h *DishHandler) Create(c echo.Context) error {
	accountID := c.Get("account_id").(uuid.UUID)

	var req CreateDishRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	if err := h.validator.Validate(req); err != nil {
		return response.Error(c, errors.NewWithDetails("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, h.validator.FormatErrors(err)))
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

	if err := h.dishService.Create(accountID, dish); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, dish)
}

func (h *DishHandler) GetByRestaurant(c echo.Context) error {
	accountID := c.Get("account_id").(uuid.UUID)
	
	restaurantID, err := uuid.Parse(c.Param("restaurantId"))
	if err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	dishes, err := h.dishService.GetByRestaurantID(accountID, restaurantID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, dishes)
}

func (h *DishHandler) GetByID(c echo.Context) error {
	accountID := c.Get("account_id").(uuid.UUID)
	
	dishID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	dish, err := h.dishService.GetByID(accountID, dishID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, dish)
}

func (h *DishHandler) Update(c echo.Context) error {
	accountID := c.Get("account_id").(uuid.UUID)
	
	dishID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	var updates map[string]interface{}
	if err := c.Bind(&updates); err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	if err := h.dishService.Update(accountID, dishID, updates); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]string{
		"message": "Dish updated successfully",
	})
}

func (h *DishHandler) Delete(c echo.Context) error {
	accountID := c.Get("account_id").(uuid.UUID)
	
	dishID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	if err := h.dishService.Delete(accountID, dishID); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]string{
		"message": "Dish deleted successfully",
	})
}
