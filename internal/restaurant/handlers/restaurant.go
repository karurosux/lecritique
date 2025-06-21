package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lecritique/api/internal/restaurant/models"
	"github.com/lecritique/api/internal/restaurant/services"
	"github.com/lecritique/api/internal/shared/errors"
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

func (h *RestaurantHandler) Create(c echo.Context) error {
	accountID := c.Get("account_id").(uuid.UUID)

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

	if err := h.restaurantService.Create(accountID, restaurant); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, restaurant)
}

func (h *RestaurantHandler) GetAll(c echo.Context) error {
	accountID := c.Get("account_id").(uuid.UUID)

	restaurants, err := h.restaurantService.GetByAccountID(accountID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, restaurants)
}

func (h *RestaurantHandler) GetByID(c echo.Context) error {
	accountID := c.Get("account_id").(uuid.UUID)
	
	restaurantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	restaurant, err := h.restaurantService.GetByID(accountID, restaurantID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, restaurant)
}

func (h *RestaurantHandler) Update(c echo.Context) error {
	accountID := c.Get("account_id").(uuid.UUID)
	
	restaurantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	var updates map[string]interface{}
	if err := c.Bind(&updates); err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	if err := h.restaurantService.Update(accountID, restaurantID, updates); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]string{
		"message": "Restaurant updated successfully",
	})
}

func (h *RestaurantHandler) Delete(c echo.Context) error {
	accountID := c.Get("account_id").(uuid.UUID)
	
	restaurantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	if err := h.restaurantService.Delete(accountID, restaurantID); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]string{
		"message": "Restaurant deleted successfully",
	})
}
