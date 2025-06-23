package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lecritique/api/internal/shared/errors"
	"github.com/lecritique/api/internal/shared/response"
	"github.com/lecritique/api/internal/shared/validator"
	"github.com/lecritique/api/internal/subscription/services"
)

type SubscriptionHandler struct {
	subscriptionService services.SubscriptionService
	validator          *validator.Validator
}

func NewSubscriptionHandler(subscriptionService services.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
		validator:          validator.New(),
	}
}

type CreateSubscriptionRequest struct {
	PlanID string `json:"plan_id" validate:"required,uuid"`
}

// GetAvailablePlans godoc
// @Summary Get available subscription plans
// @Description Retrieve all available subscription plans with their features and pricing
// @Tags subscription
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]models.SubscriptionPlan}
// @Failure 500 {object} response.Response
// @Router /api/v1/plans [get]
func (h *SubscriptionHandler) GetAvailablePlans(c echo.Context) error {
	ctx := c.Request().Context()

	plans, err := h.subscriptionService.GetAvailablePlans(ctx)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, plans)
}

// GetUserSubscription godoc
// @Summary Get user's current subscription
// @Description Retrieve the current subscription details for the authenticated user
// @Tags subscription
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=models.Subscription}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/user/subscription [get]
func (h *SubscriptionHandler) GetUserSubscription(c echo.Context) error {
	ctx := c.Request().Context()
	accountID := c.Get("account_id").(uuid.UUID)

	subscription, err := h.subscriptionService.GetUserSubscription(ctx, accountID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, subscription)
}

// CanUserCreateRestaurant godoc
// @Summary Check if user can create more restaurants
// @Description Check if the authenticated user can create more restaurants based on their subscription plan
// @Tags subscription
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=services.PermissionResponse}
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/user/can-create-restaurant [get]
func (h *SubscriptionHandler) CanUserCreateRestaurant(c echo.Context) error {
	ctx := c.Request().Context()
	accountID := c.Get("account_id").(uuid.UUID)

	permission, err := h.subscriptionService.CanUserCreateRestaurant(ctx, accountID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, permission)
}

// CreateSubscription godoc
// @Summary Create a new subscription
// @Description Create a new subscription for the authenticated user
// @Tags subscription
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateSubscriptionRequest true "Subscription details"
// @Success 201 {object} response.Response{data=models.Subscription}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/user/subscription [post]
func (h *SubscriptionHandler) CreateSubscription(c echo.Context) error {
	ctx := c.Request().Context()
	accountID := c.Get("account_id").(uuid.UUID)

	var req CreateSubscriptionRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	if err := h.validator.Validate(req); err != nil {
		return response.Error(c, errors.NewWithDetails("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, h.validator.FormatErrors(err)))
	}

	planID, err := uuid.Parse(req.PlanID)
	if err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	createReq := &services.CreateSubscriptionRequest{
		AccountID: accountID,
		PlanID:    planID,
	}

	subscription, err := h.subscriptionService.CreateSubscription(ctx, createReq)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, subscription)
}

// CancelSubscription godoc
// @Summary Cancel user's subscription
// @Description Cancel the current subscription for the authenticated user
// @Tags subscription
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/user/subscription [delete]
func (h *SubscriptionHandler) CancelSubscription(c echo.Context) error {
	ctx := c.Request().Context()
	accountID := c.Get("account_id").(uuid.UUID)

	err := h.subscriptionService.CancelSubscription(ctx, accountID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]string{"message": "Subscription cancelled successfully"})
}