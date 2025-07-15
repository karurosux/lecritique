package handlers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	sharedErrors "github.com/lecritique/api/internal/shared/errors"
	sharedRepos "github.com/lecritique/api/internal/shared/repositories"
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

type PlanFeaturesResponse struct {
	PlanName           string                 `json:"plan_name"`
	PlanCode           string                 `json:"plan_code"`
	Features           map[string]interface{} `json:"features"`
	SubscriptionStatus string                 `json:"subscription_status,omitempty"`
	IsActive           bool                   `json:"is_active"`
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
		// If no subscription found, return null instead of error
		if errors.Is(err, sharedRepos.ErrRecordNotFound) {
			return response.Success(c, nil)
		}
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
		return response.Error(c, sharedErrors.ErrBadRequest)
	}

	if err := h.validator.Validate(req); err != nil {
		return response.Error(c, sharedErrors.NewWithDetails("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, h.validator.FormatErrors(err)))
	}

	planID, err := uuid.Parse(req.PlanID)
	if err != nil {
		return response.Error(c, sharedErrors.ErrBadRequest)
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

// GetCurrentPlanFeatures godoc
// @Summary Get current plan features
// @Description Get the current subscription plan features in a structured format
// @Tags subscription
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=PlanFeaturesResponse}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/user/subscription/features [get]
func (h *SubscriptionHandler) GetCurrentPlanFeatures(c echo.Context) error {
	ctx := c.Request().Context()
	accountID := c.Get("account_id").(uuid.UUID)

	subscription, err := h.subscriptionService.GetUserSubscription(ctx, accountID)
	if err != nil {
		if errors.Is(err, sharedRepos.ErrRecordNotFound) {
			// Return default/free tier features if no subscription
			return response.Success(c, &PlanFeaturesResponse{
				PlanName: "Free",
				PlanCode: "free",
				Features: map[string]interface{}{
					"max_restaurants":         1,
					"max_qr_codes":           1,
					"max_feedbacks_per_month": 100,
					"max_team_members":        1,
					"basic_analytics":         true,
					"advanced_analytics":      false,
					"feedback_explorer":       true,
					"custom_branding":         false,
					"priority_support":        false,
				},
			})
		}
		return response.Error(c, err)
	}

	// Build response with structured features
	featuresResponse := &PlanFeaturesResponse{
		PlanName: subscription.Plan.Name,
		PlanCode: subscription.Plan.Code,
		Features: map[string]interface{}{
			"limits": map[string]int{
				"max_restaurants":       subscription.Plan.MaxRestaurants,
				"max_qr_codes":         subscription.Plan.MaxQRCodes,
				"max_feedbacks_per_month": subscription.Plan.MaxFeedbacksPerMonth,
				"max_team_members":      subscription.Plan.MaxTeamMembers,
			},
			"flags": map[string]bool{
				"basic_analytics":    subscription.Plan.HasBasicAnalytics,
				"advanced_analytics": subscription.Plan.HasAdvancedAnalytics,
				"feedback_explorer":  subscription.Plan.HasFeedbackExplorer,
				"custom_branding":    subscription.Plan.HasCustomBranding,
				"priority_support":   subscription.Plan.HasPrioritySupport,
			},
		},
		SubscriptionStatus: string(subscription.Status),
		IsActive:          subscription.IsActive(),
	}

	return response.Success(c, featuresResponse)
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